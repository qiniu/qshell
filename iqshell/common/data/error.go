package data

import (
	"fmt"
)

var (
	ErrorCodeUnknown       = -10000
	ErrorCodeParamNotExist = -11000
	ErrorCodeParamMissing  = -11001
	ErrorCodeLineHeader    = -11002
	ErrorCodeAlreadyDone   = -15000
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

func (c *CodeError) SetCode(code int) *CodeError {
	c.Code = code
	return c
}

func (c *CodeError) HeaderInsertDesc(desc string) *CodeError {
	if len(desc) == 0 {
		return c
	}
	c.Desc = desc + " " + c.Desc
	return c
}

func (c *CodeError) AppendDesc(desc string) *CodeError {
	if len(c.Desc) > 0 {
		c.Desc += ", "
	}
	c.Desc += desc
	return c
}

func (c *CodeError) HeaderInsertDescF(f string, a ...interface{}) *CodeError {
	c.Desc = fmt.Sprintf(f, a...) + ", " + c.Desc
	return c
}

func (c *CodeError) AppendDescF(f string, a ...interface{}) *CodeError {
	if len(c.Desc) > 0 {
		c.Desc += ", "
	}
	c.Desc += fmt.Sprintf(f, a...)
	return c
}

func (c *CodeError) AppendError(err error) *CodeError {
	if err != nil {
		if len(c.Desc) > 0 {
			c.Desc += " => "
		}
		c.Desc += err.Error()
	}
	return c
}

func NewErrorWithError(code int, desc string, err error) *CodeError {
	e := &CodeError{}
	e.Code = code
	e.Desc = desc
	if err != nil {
		e.Desc += " => " + err.Error()
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
	if c.Code == 0 {
		return c.Desc
	}
	return fmt.Sprintf("【%d】%s", c.Code, c.Desc)
}

func ErrorCode(err error) *Int {
	if e, ok := err.(*CodeError); ok {
		return NewInt(e.Code)
	} else {
		return nil
	}
}
