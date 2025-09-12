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
	var offset int
	sep := ": "
	if len(message) != 0 {
		offset = len([]rune(sep))
	}

	errStr := err.Error()
	output := make([]byte, len(errStr)+len(message)+offset)

	copy(output, message)
	copy(output[len(message):len(message)+offset], sep)
	copy(output[len(message)+offset:], errStr)

	return &messageError{
		err:     err,
		message: message,
		output:  string(output),
	}
}
