package data

import (
	"errors"
	"fmt"
)

var (
	ErrorCodeUnknown     = 10000
	ErrorCodeAlreadyDone = 10001
)

func NewAlreadyDoneError(desc string) error {
	return &codeError{
		Code: ErrorCodeAlreadyDone,
		err:  errors.New(desc),
	}
}

func NewError(code int, err error) error {
	return &codeError{
		Code: code,
		err:  err,
	}
}

type codeError struct {
	Code int
	err  error
}

func (c codeError) Error() string {
	return fmt.Sprintf("Code:%d desc:%v", c.Code, c.err)
}

func ErrorCode(err error) *Int {
	if e, ok := err.(codeError); ok {
		return NewInt(e.Code)
	} else {
		return nil
	}
}
