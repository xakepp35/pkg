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

	return &messageError{
		err:     err,
		message: message,
		output:  string(output),
	}
}
