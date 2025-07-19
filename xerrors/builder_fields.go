package xerrors

import (
	"github.com/xakepp35/pkg/src/pkg/types"
	"time"
)

type Fielder interface {
	Str(field, value string) ErrBuilder
	Float64(field string, value float64) ErrBuilder
	Float32(field string, value float32) ErrBuilder
	Int64(field string, value int64) ErrBuilder
	Uint64(field string, value uint64) ErrBuilder
	Int(field string, value int) ErrBuilder
	Int8(field string, value int8) ErrBuilder
	Int16(field string, value int16) ErrBuilder
	Int32(field string, value int32) ErrBuilder
	Uint(field string, value uint) ErrBuilder
	Uint8(field string, value uint8) ErrBuilder
	Uint16(field string, value uint16) ErrBuilder
	Uint32(field string, value uint32) ErrBuilder
	Bool(field string, value bool) ErrBuilder
	Strs(field string, values []string) ErrBuilder
	Float64s(field string, values []float64) ErrBuilder
	Float32s(field string, values []float32) ErrBuilder
	Int64s(field string, values []int64) ErrBuilder
	Uint64s(field string, values []uint64) ErrBuilder
	Ints(field string, values []int) ErrBuilder
	Int8s(field string, values []int8) ErrBuilder
	Int16s(field string, values []int16) ErrBuilder
	Int32s(field string, values []int32) ErrBuilder
	Uints(field string, values []uint) ErrBuilder
	Uint8s(field string, values []uint8) ErrBuilder
	Uint16s(field string, values []uint16) ErrBuilder
	Uint32s(field string, values []uint32) ErrBuilder
	Bools(field string, values []bool) ErrBuilder
	Time(field string, value time.Time) ErrBuilder
	Times(field string, values []time.Time) ErrBuilder
	XTime(field string, value *types.Time) ErrBuilder
	XTimes(field string, values []*types.Time) ErrBuilder
}

func (e *errorBuilder) Bool(field string, value bool) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendBool(e.argsBuffer, value)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Str(field, value string) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = append(e.argsBuffer, value...)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Float64(field string, value float64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendFloat64(e.argsBuffer, value, -1)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Float32(field string, value float32) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendFloat32(e.argsBuffer, value, -1)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int64(field string, value int64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInt64(e.argsBuffer, value)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uint64(field string, value uint64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUint64(e.argsBuffer, value)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int(field string, value int) ErrBuilder {
	return e.Int64(field, int64(value))
}

func (e *errorBuilder) Int8(field string, value int8) ErrBuilder {
	return e.Int64(field, int64(value))
}

func (e *errorBuilder) Int16(field string, value int16) ErrBuilder {
	return e.Int64(field, int64(value))
}

func (e *errorBuilder) Int32(field string, value int32) ErrBuilder {
	return e.Int64(field, int64(value))
}

func (e *errorBuilder) Uint(field string, value uint) ErrBuilder {
	return e.Uint64(field, uint64(value))
}

func (e *errorBuilder) Uint8(field string, value uint8) ErrBuilder {
	return e.Uint64(field, uint64(value))
}

func (e *errorBuilder) Uint16(field string, value uint16) ErrBuilder {
	return e.Uint64(field, uint64(value))
}

func (e *errorBuilder) Uint32(field string, value uint32) ErrBuilder {
	return e.Uint64(field, uint64(value))
}

func (e *errorBuilder) Strs(field string, values []string) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendStrings(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Float64s(field string, values []float64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendFloats64(e.argsBuffer, values, -1)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Float32s(field string, values []float32) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendFloats32(e.argsBuffer, values, -1)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int64s(field string, values []int64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInts64(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uint64s(field string, values []uint64) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUints64(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Ints(field string, values []int) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInts(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int8s(field string, values []int8) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInts8(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int16s(field string, values []int16) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInts16(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Int32s(field string, values []int32) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendInts32(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uints(field string, values []uint) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUints(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uint8s(field string, values []uint8) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUints8(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uint16s(field string, values []uint16) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUints16(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Uint32s(field string, values []uint32) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendUints32(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Bools(field string, values []bool) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendBools(e.argsBuffer, values)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Time(field string, value time.Time) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendTime(e.argsBuffer, value, time.RFC3339)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) Times(field string, values []time.Time) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendTimes(e.argsBuffer, values, time.RFC3339)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) XTime(field string, value *types.Time) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	e.argsBuffer = enc.AppendTime(e.argsBuffer, value.AsTime(), time.RFC3339)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}

func (e *errorBuilder) XTimes(field string, values []*types.Time) ErrBuilder {
	e.argsBuffer = append(e.argsBuffer, field...)
	e.argsBuffer = append(e.argsBuffer, '=')
	times := make([]time.Time, len(values))
	for i, v := range values {
		times[i] = v.AsTime()
	}
	e.argsBuffer = enc.AppendTimes(e.argsBuffer, times, time.RFC3339)
	e.argsBuffer = append(e.argsBuffer, ' ')
	return e
}
