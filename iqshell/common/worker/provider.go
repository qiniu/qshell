package worker

import "errors"

// WorkProvider work 提供者
type WorkProvider interface {
	Provide() (work Work, hasMore bool, err error)
}

func NewWorkProvider(provider func() (work Work, hasMore bool, err error)) WorkProvider {
	return &workProvider{
		provider: provider,
	}
}

type workProvider struct {
	provider func() (work Work, hasMore bool, err error)
}

func (w *workProvider) Provide() (work Work, hasMore bool, err error) {
	if w == nil || w.provider == nil {
		return nil, false, errors.New("no work provider found")
	} else {
		return w.provider()
	}
}

// WorkerProvider worker 提供者
type WorkerProvider interface {
	Provide() (Worker, error)
}

func NewWorkerProvider(provider func() (Worker, error)) WorkerProvider {
	return &workerProvider{
		provider: provider,
	}
}

type workerProvider struct {
	provider func() (Worker, error)
}

func (w *workerProvider) Provide() (Worker, error) {
	if w == nil || w.provider == nil {
		return nil, errors.New("no worker provider found")
	} else {
		return w.provider()
	}
}
