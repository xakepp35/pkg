package xpgx

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xakepp35/pkg/xlog"
)

type tracerBatch struct {
	pool sync.Pool
}

func NewTracerBatch() pgx.BatchTracer {
	return &tracerBatch{
		pool: sync.Pool{
			New: func() any { return new(tracerBatchData) },
		},
	}
}

type tracerBatchData struct {
	Batch     *pgx.Batch
	StartedAt time.Time
}

func (s *tracerBatch) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	traceData := s.pool.Get().(*tracerBatchData)
	traceData.Batch = data.Batch
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerBatch) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	traceData, _ := ctx.Value(s).(*tracerBatchData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("sql", data.SQL).
		Any("args", data.Args).
		Str("command_tag", data.CommandTag.String()).
		Dur("cost", duration).
		Uint32("pid", conn.PgConn().PID()).
		Msg("batch query")
}

func (s *tracerBatch) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	traceData, _ := ctx.Value(s).(*tracerBatchData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Int("len", traceData.Batch.Len()).
		Dur("cost", duration).
		Uint32("pid", conn.PgConn().PID()).
		Msg("batch end")
}
