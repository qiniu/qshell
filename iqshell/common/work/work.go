package work

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"sync"
)

type FlowHandler interface {
	ReadWork(func() (work Work, hasMore bool)) FlowHandler
	DoWork(func(work Work) (Result, error)) FlowHandler
	OnWorkError(func(work Work, err error)) FlowHandler
	OnWorkResult(func(work Work, result Result)) FlowHandler
	OnWorksComplete(func()) FlowHandler
	Start()
}

func NewFlowHandler(info Info) FlowHandler {
	return &flowHandler{
		info:   &info,
		worker: &worker{},
	}
}

type flowHandler struct {
	info                *Info
	worker              *worker
	workReader          func() (work Work, hasMore bool)
	workHandler         func(work Work) (Result, error)
	workErrorHandler    func(action Work, err error)
	workResultHandler   func(action Work, result Result)
	workCompleteHandler func()
}

func (b *flowHandler) ReadWork(reader func() (work Work, hasMore bool)) FlowHandler {
	b.worker.reader = reader
	return b
}

func (b *flowHandler) DoWork(handler func(work Work) (Result, error)) FlowHandler {
	b.worker.handler = handler
	return b
}

func (b *flowHandler) OnWorkError(handler func(worker Work, err error)) FlowHandler {
	b.workErrorHandler = handler
	return b
}
func (b *flowHandler) handleActionError(worker Work, err error) {
	if b.workErrorHandler != nil {
		b.workErrorHandler(worker, err)
	}
}
func (b *flowHandler) OnWorkResult(handler func(worker Work, result Result)) FlowHandler {
	b.workResultHandler = handler
	return b
}
func (b *flowHandler) handlerActionResult(worker Work, result Result) {
	if b.workResultHandler != nil {
		b.workResultHandler(worker, result)
	}
}
func (b *flowHandler) OnWorksComplete(handler func()) FlowHandler {
	b.workCompleteHandler = handler
	return b
}
func (b *flowHandler) handlerComplete() {
	if b.workCompleteHandler != nil {
		b.workCompleteHandler()
	}
}

func (b *flowHandler) Start() {
	log.Debug("work flow did start")

	workChan := make(chan Work, b.info.WorkCount)
	// 生产者
	go func() {
		log.DebugF("work producer start")
		for {
			if workspace.IsCmdInterrupt() {
				break
			}

			work, hasMore := b.worker.ReadWork()
			if !hasMore {
				break
			}
			if work != nil {
				workChan <- work
			}
		}
		close(workChan)
		log.DebugF("work producer   end")
	}()

	// 消费者
	wait := &sync.WaitGroup{}
	wait.Add(b.info.WorkCount)
	for i := 0; i < b.info.WorkCount; i++ {
		go func(index int) {
			log.DebugF("work consumer %d start", index)
			for work := range workChan {
				result, err := b.worker.DoWork(work)
				if err != nil {
					b.handleActionError(work, err)
					b.info.workErrorHappened = true
				} else {
					b.handlerActionResult(work, result)
				}

				// 检测是否需要停止
				if b.info.workErrorHappened && b.info.StopWhenWorkError {
					break
				}
			}
			wait.Done()
			log.DebugF("work consumer %d   end", index)
		}(i)
	}
	wait.Wait()

	log.Debug("work flow did end")

	b.handlerComplete()
}
