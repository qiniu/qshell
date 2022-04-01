package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Worker interface {
	DoWork(work Work) (Result, *data.CodeError)
}

func NewWorker(doFunc func(work Work) (Result, *data.CodeError)) Worker {
	return &worker{
		DoFunc: doFunc,
	}
}

type worker struct {
	DoFunc func(work Work) (Result, *data.CodeError)
}

func (w *worker) DoWork(work Work) (Result, *data.CodeError) {
	if w == nil || w.DoFunc == nil {
		return nil, alert.Error("worker: no worker func", "")
	}
	return w.DoFunc(work)
}
