package data

type Checker interface {
	Check() *CodeError
}
