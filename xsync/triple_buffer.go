package xsync

import "sync/atomic"

type BufferFunc func(dst []byte)

type TripleBuffer struct {
	buffers      [3][]byte // Фиксированные три буфера
	writeIndex   int32     // Индекс текущего буфера для записи
	flushIndex   int32     // Индекс текущего буфера для сброса
	writeOffsets [3]int32  // Текущая длина данных в каждом буфере
}

// NewTripleBuffer создаёт новый TripleBuffer с фиксированным размером буферов
func NewTripleBuffer(bufferSize int) *TripleBuffer {
	return &TripleBuffer{
		buffers: [3][]byte{
			make([]byte, 0, bufferSize),
			make([]byte, 0, bufferSize),
			make([]byte, 0, bufferSize),
		},
	}
}

// Write пытается записать данные в текущий буфер
func (tb *TripleBuffer) Write(length int32, writeFunc BufferFunc) bool {
	for {
		writeIdx := atomic.LoadInt32(&tb.writeIndex)
		offset := atomic.LoadInt32(&tb.writeOffsets[writeIdx])
		if offset+length > int32(cap(tb.buffers[writeIdx])) {
			// Если текущий буфер заполнен, переключаемся на следующий
			nextIdx := (writeIdx + 1) % 3
			if atomic.CompareAndSwapInt32(&tb.writeIndex, writeIdx, nextIdx) {
				continue
			}
		} else if atomic.CompareAndSwapInt32(&tb.writeOffsets[writeIdx], offset, offset+length) {
			writeFunc(tb.buffers[writeIdx][offset : offset+length])
			return true
		}
	}
}

// Flush сбрасывает данные из текущего буфера для сброса
func (tb *TripleBuffer) Flush(flushFunc BufferFunc) {
	flushIdx := atomic.LoadInt32(&tb.flushIndex)
	for {
		offset := atomic.LoadInt32(&tb.writeOffsets[flushIdx])
		if offset == 0 {
			// Если в буфере нет данных, переключаемся на следующий
			nextIdx := (flushIdx + 1) % 3
			if atomic.CompareAndSwapInt32(&tb.flushIndex, flushIdx, nextIdx) {
				return
			}
			continue
		}
		flushFunc(tb.buffers[flushIdx][:offset])
		atomic.StoreInt32(&tb.writeOffsets[flushIdx], 0)
		return
	}
}
