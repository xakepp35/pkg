package xerrors

import (
	"bytes"
	"errors"
	"strconv"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var buildersPool = sync.Pool{
	New: func() interface{} {
		return &errorBuilder{
			buffer: bytes.Buffer{},
		}
	},
}

var _ = ErrBuilder(&errorBuilder{})

type ErrBuilder interface {
	Err() error
	ProtoErr(code codes.Code) error

	Msg(msg string) ErrBuilder
	Str(field, value string) ErrBuilder
	Int64(field string, value int64) ErrBuilder
}

type errorBuilder struct {
	err        error
	buffer     bytes.Buffer
	hasFields  bool
	hasMessage bool
}

func (e *errorBuilder) resetSelf() {
	e.buffer.Reset()
	e.hasFields = false
	e.hasMessage = false
	buildersPool.Put(e)
}

// ProtoErr *status.Status proto error
func (e *errorBuilder) ProtoErr(code codes.Code) error {
	defer e.resetSelf()
	return status.Error(code, e.renderErr().Error())
}

func (e *errorBuilder) renderErr() error {
	if e.err == nil {
		return errors.New(e.buffer.String())
	}

	return New(e.err, e.buffer.String())
}

// Err return compiled err wrapped in message
func (e *errorBuilder) Err() error {
	defer e.resetSelf()
	return e.renderErr()
}

func (e *errorBuilder) Msg(msg string) ErrBuilder {
	e.buffer.WriteString(msg)
	e.hasMessage = true
	return e
}

func (e *errorBuilder) Str(field, value string) ErrBuilder {
	if e.hasFields || e.hasMessage {
		e.buffer.WriteByte(' ') // ставить ли пробел, чтобы не было пробела перед ": "
	}
	e.buffer.WriteString(field)
	e.buffer.WriteRune('=')
	e.buffer.WriteString(value)
	e.hasFields = true

	return e
}

func (e *errorBuilder) Int64(field string, value int64) ErrBuilder {
	if e.hasFields || e.hasMessage {
		e.buffer.WriteByte(' ') // ставить ли пробел, чтобы не было пробела перед ": "
	}
	e.buffer.WriteString(field)
	e.buffer.WriteRune('=')
	e.buffer.WriteString(strconv.FormatInt(value, 10))
	e.hasFields = true

	return e
}

func Err(err error) ErrBuilder {
	builder := buildersPool.Get().(*errorBuilder)
	builder.err = err

	return builder
}
