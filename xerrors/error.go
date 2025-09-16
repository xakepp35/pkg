package xerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type messageError struct {
	code    codes.Code
	err     error
	message string
	output  string
}

func (m *messageError) Unwrap() error {
	return m.err
}

func (m *messageError) Error() string {
	return m.output
}

func (m *messageError) GRPCStatus() *status.Status {
	return status.New(m.code, m.message)
}

func New(err error, message string) error {
	return &messageError{
		err:     err,
		message: message,
		output:  outputBuild(err, message),
	}
}

func NewProto(code codes.Code, err error, message string) error {
	return &messageError{
		err:     err,
		message: message,
		output:  outputBuild(err, message),
		code:    code,
	}
}

func outputBuild(err error, message string) string{
	var output []byte

	errStr := err.Error()

	if len(message) != 0 {
		sep := ": "
		offset := len([]rune(sep))

		output = make([]byte, len(errStr)+len(message)+offset)
		copy(output[len(message):len(message)+offset], sep)
		copy(output[len(message)+offset:], errStr)
		copy(output, message)
	} else {

		output = make([]byte, len(errStr))
		copy(output[len(message):], errStr)
		copy(output, message)
	}
	return string(output)
}
