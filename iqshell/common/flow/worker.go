package flow

type Worker interface {
	DoWork(work Work) (Result, error)
}