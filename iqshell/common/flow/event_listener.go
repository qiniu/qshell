package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type EventListener struct {
	FlowWillStartFunc func(flow *Flow) (err *data.CodeError)
	FlowWillEndFunc   func(flow *Flow) (err *data.CodeError)
	WillWorkFunc      func(work *WorkInfo) (shouldContinue bool, err *data.CodeError)
	OnWorkSkipFunc    func(work *WorkInfo, err *data.CodeError)
	OnWorkSuccessFunc func(work *WorkInfo, result Result)
	OnWorkFailFunc    func(work *WorkInfo, err *data.CodeError)
}

func (e *EventListener) FlowWillStart(flow *Flow) (err *data.CodeError) {
	if e.FlowWillStartFunc == nil {
		return nil
	}
	return e.FlowWillStartFunc(flow)
}

func (e *EventListener) FlowWillEnd(flow *Flow) (err *data.CodeError) {
	if e.FlowWillEndFunc == nil {
		return nil
	}
	return e.FlowWillEndFunc(flow)
}

func (e *EventListener) WillWork(work *WorkInfo) (shouldContinue bool, err *data.CodeError) {
	if e.WillWorkFunc == nil {
		return true, nil
	}
	return e.WillWorkFunc(work)
}

func (e *EventListener) OnWorkSkip(work *WorkInfo, err *data.CodeError) {
	if e.OnWorkSkipFunc == nil {
		return
	}
	e.OnWorkSkipFunc(work, err)
}

func (e *EventListener) OnWorkSuccess(work *WorkInfo, result Result) {
	if e.OnWorkSuccessFunc == nil {
		return
	}
	e.OnWorkSuccessFunc(work, result)
}

func (e *EventListener) OnWorkFail(work *WorkInfo, err *data.CodeError) {
	if e.OnWorkFailFunc == nil {
		return
	}
	e.OnWorkFailFunc(work, err)
}
