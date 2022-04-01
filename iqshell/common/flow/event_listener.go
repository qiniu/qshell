package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type EventListener struct {
	WillWorkFunc      func(work Work) (shouldContinue bool, err *data.CodeError)
	OnWorkSkipFunc    func(work Work, err *data.CodeError)
	OnWorkSuccessFunc func(work Work, result Result)
	OnWorkFailFunc    func(work Work, err *data.CodeError)
}

func (e *EventListener) WillWork(work Work) (shouldContinue bool, err *data.CodeError) {
	if e.WillWorkFunc == nil {
		return true, nil
	}
	return e.WillWorkFunc(work)
}

func (e *EventListener) OnWorkSkip(work Work, err *data.CodeError) {
	if e.OnWorkSkipFunc == nil {
		return
	}
	e.OnWorkSkipFunc(work, err)
}

func (e *EventListener) OnWorkSuccess(work Work, result Result) {
	if e.OnWorkSuccessFunc == nil {
		return
	}
	e.OnWorkSuccessFunc(work, result)
}

func (e *EventListener) OnWorkFail(work Work, err *data.CodeError) {
	if e.OnWorkFailFunc == nil {
		return
	}
	e.OnWorkFailFunc(work, err)
}
