package xsync

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTripleBuffer_Write(t *testing.T) {
	bufferSize := 64
	tb := NewTripleBuffer(bufferSize)
	data := []byte("hello world")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	assert.True(t, tb.Write(length, writeFunc), "Write failed when it should have succeeded")

	writeIdx := atomic.LoadInt32(&tb.writeIndex)
	offset := atomic.LoadInt32(&tb.writeOffsets[writeIdx]) - length
	assert.Equal(t, string(data), string(tb.buffers[writeIdx][offset:offset+length]), "Data mismatch after write")
}

func TestTripleBuffer_Write_Overflow(t *testing.T) {
	bufferSize := 8
	tb := NewTripleBuffer(bufferSize)
	data := []byte("12345678")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	assert.True(t, tb.Write(length, writeFunc), "Write failed when it should have succeeded")
	assert.False(t, tb.Write(length, writeFunc), "Write succeeded when it should have failed due to overflow")
}

func TestTripleBuffer_Flush(t *testing.T) {
	bufferSize := 64
	tb := NewTripleBuffer(bufferSize)
	data := []byte("flush me")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	flushFunc := func(dst []byte) {
		assert.Equal(t, string(data), string(dst), "Data mismatch during flush")
	}

	assert.True(t, tb.Write(length, writeFunc), "Write failed when it should have succeeded")
	tb.Flush(flushFunc)
}

func TestTripleBuffer_WriteAndFlushCycle(t *testing.T) {
	bufferSize := 16
	tb := NewTripleBuffer(bufferSize)

	data := [][]byte{
		[]byte("first"),
		[]byte("second"),
		[]byte("third"),
	}

	writeFunc := func(idx int) BufferFunc {
		return func(dst []byte) {
			copy(dst, data[idx])
		}
	}

	flushFunc := func(expected []byte) BufferFunc {
		return func(dst []byte) {
			assert.Equal(t, string(expected), string(dst), "Data mismatch during flush")
		}
	}

	for i := 0; i < len(data); i++ {
		assert.True(t, tb.Write(int32(len(data[i])), writeFunc(i)), "Write failed at index %d", i)
		tb.Flush(flushFunc(data[i]))
	}
}

func TestTripleBuffer_ConcurrentWrite(t *testing.T) {
	bufferSize := 64
	tb := NewTripleBuffer(bufferSize)
	data := []byte("concurrent write")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	wg := sync.WaitGroup{}
	writeCount := 100

	wg.Add(writeCount)
	for i := 0; i < writeCount; i++ {
		go func() {
			defer wg.Done()
			tb.Write(length, writeFunc)
		}()
	}

	wg.Wait()
}

func TestTripleBuffer_ConcurrentFlush(t *testing.T) {
	bufferSize := 64
	tb := NewTripleBuffer(bufferSize)
	data := []byte("flush concurrently")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	flushFunc := func(dst []byte) {
		assert.Equal(t, string(data), string(dst), "Data mismatch during concurrent flush")
	}

	assert.True(t, tb.Write(length, writeFunc), "Write failed when it should have succeeded")

	wg := sync.WaitGroup{}
	flushCount := 10

	wg.Add(flushCount)
	for i := 0; i < flushCount; i++ {
		go func() {
			defer wg.Done()
			tb.Flush(flushFunc)
		}()
	}

	wg.Wait()
}

func TestTripleBuffer_ConcurrentWriteAndFlush(t *testing.T) {
	bufferSize := 64
	tb := NewTripleBuffer(bufferSize)
	data := []byte("write and flush")
	length := int32(len(data))

	writeFunc := func(dst []byte) {
		copy(dst, data)
	}

	flushFunc := func(dst []byte) {
		assert.Equal(t, string(data), string(dst), "Data mismatch during concurrent write and flush")
	}

	wg := sync.WaitGroup{}
	operationCount := 100

	wg.Add(operationCount * 2)
	for i := 0; i < operationCount; i++ {
		go func() {
			defer wg.Done()
			tb.Write(length, writeFunc)
		}()
		go func() {
			defer wg.Done()
			tb.Flush(flushFunc)
		}()
	}

	wg.Wait()
}
