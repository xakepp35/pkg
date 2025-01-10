package xsync

import "sync/atomic"

type RingBuffer struct {
	data   []byte
	size   int64
	head   int64 // куда пишем
	tail   int64 // откуда читаем
	writer int64
}

// NewRingBuffer создаёт кольцевой буфер с указанным размером
func NewRingBuffer(size int64) *RingBuffer {
	return &RingBuffer{
		data: make([]byte, size),
		size: size,
	}
}

// Write пытается записать данные в буфер
func (rb *RingBuffer) Write(entry []byte) bool {
	entryLen := int64(len(entry))
	if entryLen > rb.size {
		return false // слишком большой entry
	}

	for {
		head := atomic.LoadInt64(&rb.head)
		tail := atomic.LoadInt64(&rb.tail)

		used := (head - tail + rb.size) % rb.size
		free := rb.size - used
		if free <= entryLen {
			return false // недостаточно места
		}

		newHead := (head + entryLen) % rb.size
		if atomic.CompareAndSwapInt64(&rb.head, head, newHead) {
			// Пишем данные в буфер
			if head+entryLen <= rb.size {
				copy(rb.data[head:head+entryLen], entry)
			} else {
				// Обработка wrap-around
				firstPart := rb.size - head
				copy(rb.data[head:], entry[:firstPart])
				copy(rb.data[0:], entry[firstPart:])
			}
			return true
		}
	}
}

// Read читает из буфера до batchSize
func (rb *RingBuffer) Read(batchSize int64) []byte {
	tail := atomic.LoadInt64(&rb.tail)
	head := atomic.LoadInt64(&rb.head)

	if tail == head {
		return nil // буфер пуст
	}

	available := (head - tail + rb.size) % rb.size
	readSize := batchSize
	if available < batchSize {
		readSize = available
	}

	batch := make([]byte, readSize)
	if tail+readSize <= rb.size {
		copy(batch, rb.data[tail:tail+readSize])
	} else {
		firstPart := rb.size - tail
		copy(batch, rb.data[tail:])
		copy(batch[firstPart:], rb.data[:readSize-firstPart])
	}

	atomic.StoreInt64(&rb.tail, (tail+readSize)%rb.size)
	return batch
}
