package work

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"sync"
)

type FlowHandler interface {
	ReadWork(func() (work Work, hasMore bool)) FlowHandler
	OnWillWork(func(work Work) (shouldContinue bool, err error)) FlowHandler
	DoWork(func(work Work) (Result, error)) FlowHandler
	OnWorkError(func(work Work, err error)) FlowHandler
	OnWorkResult(func(work Work, result Result)) FlowHandler
	OnWorksComplete(func()) FlowHandler
	Start()
}

func NewFlowHandler(info FlowInfo) FlowHandler {
	return &flowHandler{
		info: &info,
	}
}

type flowHandler struct {
	info                *FlowInfo
	workReader          func() (work Work, hasMore bool)
	workHandler         func(work Work) (Result, error)
	willWorkHandler     func(work Work) (shouldContinue bool, err error)
	workErrorHandler    func(action Work, err error)
	workResultHandler   func(action Work, result Result)
	workCompleteHandler func()
}

func (b *flowHandler) ReadWork(reader func() (work Work, hasMore bool)) FlowHandler {
	b.workReader = reader
	return b
}

func (b *flowHandler) OnWillWork(handler func(work Work) (shouldContinue bool, err error)) FlowHandler {
	b.willWorkHandler = handler
	return b
}

func (b *flowHandler) DoWork(handler func(work Work) (Result, error)) FlowHandler {
	b.workHandler = handler
	return b
}

func (b *flowHandler) OnWorkError(handler func(worker Work, err error)) FlowHandler {
	b.workErrorHandler = handler
	return b
}

func (b *flowHandler) OnWorkResult(handler func(worker Work, result Result)) FlowHandler {
	b.workResultHandler = handler
	return b
}

func (b *flowHandler) OnWorksComplete(handler func()) FlowHandler {
	b.workCompleteHandler = handler
	return b
}

func (b *flowHandler) readWork() (Work, bool) {
	if b.workReader != nil {
		return b.workReader()
	}
	return nil, true
}

func (b *flowHandler) willWork(work Work) (bool, error) {
	if b.info.WorkOverseer != nil {
		if b.info.WorkOverseer.HasDone(work) {
			return false, data.NewAlreadyDoneError("")
		}
		b.info.WorkOverseer.WillWork(work)
	}

	if b.willWorkHandler != nil {
		return b.willWorkHandler(work)
	} else {
		return true, nil
	}
}

func (b *flowHandler) doWork(work Work) (Result, error) {
	if b.workHandler != nil {
		return b.workHandler(work)
	}
	return nil, errors.New("no worker")
}

func (b *flowHandler) handlerActionResult(work Work, result Result) {
	if b.info.WorkOverseer != nil {
		b.info.WorkOverseer.WorkDone(work, result, nil)
	}

	if b.workResultHandler != nil {
		b.workResultHandler(work, result)
	}
}

func (b *flowHandler) handleActionError(work Work, err error) {
	if b.info.WorkOverseer != nil {
		b.info.WorkOverseer.WorkDone(work, nil, err)
	}

	if b.workErrorHandler != nil {
		b.workErrorHandler(work, err)
	}
}

func (b *flowHandler) handlerComplete() {
	if b.workCompleteHandler != nil {
		b.workCompleteHandler()
	}
}

func (b *flowHandler) Start() {
	log.Debug("work flow did start")

	workChan := make(chan Work, b.info.WorkerCount)
	// 生产者
	go func() {
		log.DebugF("work producer start")
		for {
			work, hasMore := b.readWork()
			if work == nil {
				if !hasMore {
					break
				} else {
					continue
				}
			}

			// 检测 work 是否有问题
			if shouldContinue, err := b.willWork(work); !shouldContinue {
				b.handleActionError(work, err)
				continue
			}

			workChan <- work
		}
		close(workChan)
		log.DebugF("work producer   end")
	}()

	// 消费者
	wait := &sync.WaitGroup{}
	wait.Add(b.info.WorkerCount)
	for i := 0; i < b.info.WorkerCount; i++ {
		go func(index int) {
			log.DebugF("work consumer %d start", index)
			for work := range workChan {
				if workspace.IsCmdInterrupt() {
					break
				}

				result, err := b.doWork(work)
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
