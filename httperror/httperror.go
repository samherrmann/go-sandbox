package httperror

import (
	"errors"
	"fmt"
)

type Error struct {
	Message string
	Source  error
}

func New(msg string, source error) *Error {
	return &Error{
		Message: msg,
		Source:  source,
	}
}

func (e *Error) Error() string {
	if e.Source != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Source.Error())
	}
	return e.Message
}

func String(err error) string {
	if err == nil {
		return ""
	}
	var clientErr *Error
	if errors.As(err, &clientErr) {
		if clientErr.Message != "" {
			return clientErr.Message
		}
	}
	return "internal server error"
}
