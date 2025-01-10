package xlog

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/xakepp35/pkg/xslice"
)

// Log — основной логгер
type Log struct {
	entries     chan []byte // Канал для сообщений
	stop        chan struct{}
	file        io.WriteCloser
	flushTicker *time.Ticker
	cfg         LogConfig
	ctx         context.Context
	cancel      context.CancelFunc
}

type LogConfig struct {
	OutputFile         string        `json:"output_file"`
	InitialLineBufSize int           `json:"initial_line_buf_size"`
	FlushInterval      time.Duration `json:"flush_interval"`
	FlushChanSize      uint64        `json:"flush_chan_size"`
	FlushBatchSize     uint64        `json:"flush_batch_size"`
}

// New создает новый логгер
func New(ctx context.Context, cfg *LogConfig) (*Log, error) {
	file, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	l := &Log{
		entries:     make(chan []byte, cfg.FlushChanSize),
		stop:        make(chan struct{}),
		file:        file,
		flushTicker: time.NewTicker(cfg.FlushInterval),
		ctx:         ctx,
		cancel:      cancel,
	}
	return l, nil
}

// New создает новую строку лога
func (l *Log) New() *Line {
	if l == nil {
		return nil
	}
	line := &Line{
		buf: *xslice.NewBufferSized(l.cfg.InitialLineBufSize),
		log: l,
	}
	line.buf.Byte('{')
	return line
}

// Run обрабатывает записи и пишет их в файл
func (l *Log) Run() {
	if l == nil {
		return
	}
	var batch [][]byte
	for {
		select {
		case entry := <-l.entries:
			batch = append(batch, entry)
			if uint64(len(batch)) >= l.cfg.FlushBatchSize {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-l.flushTicker.C:
			if len(batch) > 0 {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-l.stop:
			if len(batch) > 0 {
				l.flush(batch)
			}
			_ = l.file.Close()
			return
		}
	}
}

// flush записывает батч в файл
func (l *Log) flush(batch [][]byte) {
	if len(batch) == 0 {
		return
	}
	var buf bytes.Buffer
	for _, entry := range batch {
		buf.Write(entry)
	}
	_, _ = l.file.Write(buf.Bytes())
}

// Close завершает работу логгера
func (l *Log) Close() {
	close(l.stop)
	l.flushTicker.Stop()
}
