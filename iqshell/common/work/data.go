package work

import "errors"

type Info struct {
	WorkCount         int  // work 数量
	StopWhenWorkError bool // 当某个 action 遇到执行错误是否结束 batch 任务
	workErrorHappened bool // 执行中是否出现错误
}

func (i *Info) initData() {
	if i.WorkCount <= 0 {
		i.WorkCount = 1
	}
	i.workErrorHappened = false
}

type Work interface{}
type Result interface{}
type Worker interface {
	ReadWork() (work Work, hasMore bool)
	DoWork(work Work) (Result, error)
}

func NewWorker(reader func() (Work, bool), handler func(Work) (Result, error)) Worker {
	return &worker{
		reader:  reader,
		handler: handler,
	}
}

type worker struct {
	reader  func() (Work, bool)
	handler func(Work) (Result, error)
}

func (w *worker) ReadWork() (Work, bool) {
	if w.reader != nil {
		return w.reader()
	}
	return nil, true
}

func (w *worker) DoWork(action Work) (Result, error) {
	if w.handler != nil {
		return w.handler(action)
	}
	return nil, errors.New("no worker")
}
