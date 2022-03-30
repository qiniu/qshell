package data

import "fmt"

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
