package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

func New(info Info) *WorkProvideBuilder {
	return &WorkProvideBuilder{
		flow: &Flow{
			Info: info,
		},
	}
}

type WorkProvideBuilder struct {
	flow *Flow
	err  error
}

func (b *WorkProvideBuilder) WorkProvider(provider WorkProvider) *WorkerProvideBuilder {
	b.flow.WorkProvider = provider
	return &WorkerProvideBuilder{
		flow: b.flow,
		err:  nil,
	}
}

func (b *WorkProvideBuilder) WorkProviderWithArray(workList []Work) *WorkerProvideBuilder {
	if provider, err := NewArrayWorkProvider(workList); err != nil {
		return &WorkerProvideBuilder{
			flow: b.flow,
			err:  err,
		}
	} else {
		b.flow.WorkProvider = provider
		return &WorkerProvideBuilder{
			flow: b.flow,
			err:  b.err,
		}
	}
}

func (b *WorkProvideBuilder) WorkProviderWithFile(filePath string, enableStdin bool, creator WorkCreator) *WorkerProvideBuilder {
	if provider, err := NewWorkProviderOfFile(filePath, enableStdin, creator); err != nil {
		return &WorkerProvideBuilder{
			flow: b.flow,
			err:  err,
		}
	} else {
		b.flow.WorkProvider = provider
		return &WorkerProvideBuilder{
			flow: b.flow,
			err:  b.err,
		}
	}
}

type WorkerProvideBuilder struct {
	flow *Flow
	err  error
}

func (b *WorkProvideBuilder) WorkerProvider(provider WorkerProvider) *FlowBuilder {
	b.flow.WorkerProvider = provider
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

func (b *WorkProvideBuilder) OnWillWork(f func(work Work) (shouldContinue bool, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.WillWorkFunc = f
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

func (b *WorkProvideBuilder) OnWorkSkip(f func(work Work, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.OnWorkSkipFunc = f
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

func (b *WorkProvideBuilder) OnWorkSuccess(f func(work Work, result Result)) *FlowBuilder {
	b.flow.EventListener.OnWorkSuccessFunc = f
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

func (b *WorkProvideBuilder) OnWorkFail(f func(work Work, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.OnWorkFailFunc = f
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

type FlowBuilder struct {
	flow *Flow
	err  error
}

func (b *FlowBuilder) Builder() *Flow {
	return b.flow
}
