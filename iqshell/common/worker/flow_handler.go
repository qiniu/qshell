package worker

type FlowHandler interface {
	OnWorkSuccess(work Work, result Result)
	OnWorkError(work Work, err error)
	OnComplete()
}

type FlowHandlerBuilder interface {
	WorkSuccessHandler(handler func(work Work, err error))
	WorkErrorHandler(handler func(work Work, result Result))
	WorkCompleteHandler(handler func())
	Build() FlowHandler
}

func NewFlowHandlerBuilder() FlowHandlerBuilder {
	return &flowHandlerBuilder{
		handler: &flowHandler{},
	}
}

type flowHandlerBuilder struct {
	handler *flowHandler
}

func (b *flowHandlerBuilder) WorkSuccessHandler(handler func(work Work, err error)) {
	b.handler.workSuccessHandler = handler
}

func (b *flowHandlerBuilder) WorkErrorHandler(handler func(work Work, result Result)) {
	b.handler.workErrorHandler = handler
}

func (b *flowHandlerBuilder) WorkCompleteHandler(handler func()) {
	b.handler.workCompleteHandler = handler
}

func (b *flowHandlerBuilder) Build() FlowHandler {
	return b.handler
}

type flowHandler struct {
	workSuccessHandler  func(work Work, err error)
	workErrorHandler    func(work Work, result Result)
	workCompleteHandler func()
}

func (g *flowHandler) OnWorkSuccess(work Work, result Result) {
	if g.workSuccessHandler != nil {
		g.OnWorkSuccess(work, result)
	}
}

func (g *flowHandler) OnWorkError(work Work, err error) {
	if g.workErrorHandler != nil {
		g.OnWorkError(work, err)
	}
}

func (g *flowHandler) OnComplete() {
	if g.workSuccessHandler != nil {
		g.OnComplete()
	}
}
