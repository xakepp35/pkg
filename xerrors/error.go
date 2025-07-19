package xerrors

type messageError struct {
	err     error
	message string

	output string
}

func (m *messageError) Unwrap() error {
	return m.err
}

func (m *messageError) Error() string {
	return m.output
}

func New(err error, message string) error {
	errStr := err.Error()
	output := make([]byte, len(errStr)+len(message)+2)

	copy(output, message)
	copy(output[len(message):len(message)+2], ": ")
	copy(output[len(message)+2:], errStr)

	return &messageError{
		err:     err,
		message: message,
		output:  string(output),
	}
}
