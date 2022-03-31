package flow

type EventListener struct {
	WillWorkFunc      func(work Work) (shouldContinue bool, err error)
	OnWorkSkipFunc    func(work Work, err error)
	OnWorkSuccessFunc func(work Work, result Result)
	OnWorkFailFunc    func(work Work, err error)
}

func (e *EventListener) WillWork(work Work) (shouldContinue bool, err error) {
	if e.WillWorkFunc == nil {
		return true, nil
	}
	return e.WillWorkFunc(work)
}

func (e *EventListener) OnWorkSkip(work Work, err error) {
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

func (e *EventListener) OnWorkFail(work Work, err error) {
	if e.OnWorkFailFunc == nil {
		return
	}
	e.OnWorkFailFunc(work, err)
}
