package xerrors

import (
	"errors"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var buildersPool = sync.Pool{
	New: func() interface{} {
		return &errorBuilder{
			// new 2Kb
			errBuffer:  make([]byte, 0, 2048),
			argsBuffer: make([]byte, 0, 2048),
		}
	},
}

var _ = ErrBuilder(&errorBuilder{})

type ErrBuilder interface {
	Send() error
	Msg(msg string) error
	MsgProto(code codes.Code, msg string) error

	Fielder
}

type errorBuilder struct {
	err        error
	errBuffer  []byte
	argsBuffer []byte
}

func (e *errorBuilder) resetSelf() {
	e.errBuffer = e.errBuffer[:0]
	e.argsBuffer = e.argsBuffer[:0]
	buildersPool.Put(e)
}

func (e *errorBuilder) renderErr(msg string) error {
	if msg != "" {
		e.errBuffer = append(e.errBuffer, msg...)

		if len(e.argsBuffer) > 0 {
			e.errBuffer = append(e.errBuffer, ' ')
		}
	}

	if len(e.argsBuffer) > 0 {
		// отрезать последний пробел
		e.errBuffer = append(e.errBuffer, e.argsBuffer[:len(e.argsBuffer)-1]...)
	}

	if e.err == nil {
		return errors.New(string(e.errBuffer))
	}

	return New(e.err, string(e.errBuffer))
}

func (e *errorBuilder) Send() error {
	defer e.resetSelf()
	return e.renderErr("")
}

func (e *errorBuilder) Msg(msg string) error {
	defer e.resetSelf()
	return e.renderErr(msg)
}

// MsgProto *status.Status proto error
func (e *errorBuilder) MsgProto(code codes.Code, msg string) error {
	defer e.resetSelf()
	return status.Error(code, e.renderErr(msg).Error())
}

func Err(err error) ErrBuilder {
	builder := buildersPool.Get().(*errorBuilder)
	builder.err = err

	return builder
}
