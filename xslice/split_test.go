package xslice_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xakepp35/pkg/xslice"
)

const (
	testSplit    = "1,22,,"
	testSplitSep = ','
)

func TestSplitBytes(t *testing.T) {
	res := xslice.SplitBytes([]byte(testSplit), testSplitSep)
	assert.Equal(t, 4, len(res))
}

func BenchmarkSplitBytes(b *testing.B) {
	ts := []byte(testSplit)
	b.ReportAllocs()
	b.SetBytes(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xslice.SplitBytes(ts, testSplitSep)
	}
}

func BenchmarkBytesSplit(b *testing.B) {
	ts := []byte(testSplit)
	sep := []byte{testSplitSep}
	b.ReportAllocs()
	b.SetBytes(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes.Split(ts, sep)
	}
}
