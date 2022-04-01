package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type WorkerProvider interface {
	Provide() (worker Worker, err *data.CodeError)
}

func NewWorkerProvider(builder func() (Worker, *data.CodeError)) WorkerProvider {
	return &workerProvider{
		Builder: builder,
	}
}

type workerProvider struct {
	Builder func() (Worker, *data.CodeError)
}

func (w *workerProvider) Provide() (Worker, *data.CodeError) {
	if w == nil || w.Builder == nil {
		return nil, alert.Error("worker: no workerProvider Builder", "")
	}
	return w.Builder()
}
