package worker

import "errors"

type Work interface {
}

type Result interface {
}

type Worker interface {
	DoWork(work Work) (Result, error)
}

func NewWorker(workHandler func(work Work) (Result, error)) Worker {
	return &worker{
		workHandler: workHandler,
	}
}

type worker struct {
	workHandler func(work Work) (Result, error)
}

func (w *worker) DoWork(work Work) (Result, error) {
	if w == nil || w.workHandler == nil {
		return nil, errors.New("no work handler found")
	} else {
		return w.workHandler(work)
	}
}
