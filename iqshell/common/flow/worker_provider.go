package flow

import "github.com/qiniu/qshell/v2/iqshell/common/alert"

type WorkerProvider interface {
	Provide() (worker Worker, err error)
}

func NewWorkerProvider(builder func() (Worker, error)) WorkerProvider {
	return &workerProvider{
		Builder: builder,
	}
}

type workerProvider struct {
	Builder func() (Worker, error)
}

func (w *workerProvider) Provide() (Worker, error) {
	if w == nil || w.Builder == nil {
		return nil, alert.Error("worker: no workerProvider Builder", "")
	}
	return w.Builder()
}
