package lib

import (
	"errors"
	"fmt"
	"strconv"
)

type FoobarError struct {
	msg      string
	original error
}

func (err *FoobarError) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.original.Error())
}

func (err *FoobarError) Unwrap() error {
	return err.original
}

func (err *FoobarError) Is(target error) bool {
	_, ok := target.(*FoobarError)
	return ok
}

func IsFoobar(err error) bool {
	return errors.Is(err, &FoobarError{})
}

func SomeFunc() error {
	// strconv.ErrSyntax is used as a dummy error here for the error
	// that might be returned by strconv.Atoi or any other operation.
	err := strconv.ErrSyntax
	return &FoobarError{"foobar", err}
}
