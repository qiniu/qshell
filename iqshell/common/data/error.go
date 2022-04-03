package data

import (
	"fmt"
)

var (
	ErrorCodeUnknown       = 10000
	ErrorCodeParamNotExist = 11000
	ErrorCodeParamMissing  = 11001
	ErrorCodeAlreadyDone   = 15000
)

func ConvertError(err error) *CodeError {
	if err == nil {
		return nil
	}
	return NewEmptyError().AppendError(err)
}

type CodeError struct {
	Code int
	Desc string
}

func NewAlreadyDoneError(desc string) *CodeError {
	return &CodeError{
		Code: ErrorCodeAlreadyDone,
		Desc: desc,
	}
}

func NewError(code int, desc string) *CodeError {
	return &CodeError{
		Code: code,
		Desc: desc,
	}
}

func NewEmptyError() *CodeError {
	return &CodeError{}
}

func (e *CodeError) SetCode(code int) *CodeError {
	e.Code = code
	return e
}

func (e *CodeError) AppendDesc(desc string) *CodeError {
	if len(e.Desc) > 0 {
		e.Desc += " "
	}
	e.Desc += desc
	return e
}

func (e *CodeError) AppendDescF(f string, a ...interface{}) *CodeError {
	if len(e.Desc) > 0 {
		e.Desc += " "
	}
	e.Desc += fmt.Sprintf(f, a...)
	return e
}

func (e *CodeError) AppendError(err error) *CodeError {
	if err != nil {
		if len(e.Desc) > 0 {
			e.Desc += " "
		}
		e.Desc += "error:" + err.Error()
	}
	return e
}

func NewErrorWithError(code int, desc string, err error) *CodeError {
	e := &CodeError{}
	e.Code = code
	e.Desc = desc
	if err != nil {
		e.Desc += ":" + err.Error()
	}
	return e
}

func NewErrorWithCode(code int) *CodeError {
	return &CodeError{
		Code: code,
	}
}

func (c *CodeError) Error() string {
	if c == nil {
		return ""
	}
	return c.Desc
}

func ErrorCode(err error) *Int {
	if e, ok := err.(*CodeError); ok {
		return NewInt(e.Code)
	} else {
		return nil
	}
}
