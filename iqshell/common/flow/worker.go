package flow

import "github.com/qiniu/qshell/v2/iqshell/common/alert"

type Worker interface {
	DoWork(work Work) (Result, error)
}

func NewWorker(doFunc func(work Work) (Result, error)) Worker {
	return &worker{
		DoFunc: doFunc,
	}
}

type worker struct {
	DoFunc func(work Work) (Result, error)
}

func (w *worker) DoWork(work Work) (Result, error) {
	if w == nil || w.DoFunc == nil {
		return nil, alert.Error("worker: no worker func", "")
	}
	return w.DoFunc(work)
}
