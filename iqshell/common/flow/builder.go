package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

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

func (b *WorkProvideBuilder) WorkProviderWithChan(works <-chan Work) *WorkerProvideBuilder {
	if provider, err := NewChanWorkProvider(works); err != nil {
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

func (b *WorkerProvideBuilder) WorkerProvider(provider WorkerProvider) *FlowBuilder {
	b.flow.WorkerProvider = provider
	return &FlowBuilder{
		flow: b.flow,
		err:  b.err,
	}
}

func (b *FlowBuilder) DoWorkListMaxCount(count int) *FlowBuilder {
	b.flow.DoWorkInfoListMaxCount = count
	return b
}

func (b *FlowBuilder) SetOverseer(overseer Overseer) *FlowBuilder {
	b.flow.Overseer = overseer
	return b
}

func (b *FlowBuilder) SetOverseerEnable(enable bool) *FlowBuilder {
	b.enableOverseer = enable
	return b
}

func (b *FlowBuilder) SetDBOverseer(dbPath string, blankWorkRecordBuilder func() *WorkRecord) *FlowBuilder {
	if overseer, err := NewDBRecordOverseer(dbPath, blankWorkRecordBuilder); err != nil {
		b.err = err
		return b
	} else {
		b.flow.Overseer = overseer
		return b
	}
}

func (b *FlowBuilder) ShouldSkip(f func(workInfo *WorkInfo) (skip bool, cause *data.CodeError)) *FlowBuilder {
	b.flow.Skipper = NewSkipper(f)
	return b
}

func (b *FlowBuilder) ShouldRedo(f func(workInfo *WorkInfo, workRecord *WorkRecord) (shouldRedo bool, cause *data.CodeError)) *FlowBuilder {
	b.flow.Redo = NewRedo(f)
	return b
}

func (b *FlowBuilder) FlowWillStartFunc(f func(flow *Flow) (err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.FlowWillStartFunc = f
	return b
}

func (b *FlowBuilder) FlowWillEndFunc(f func(flow *Flow) (err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.FlowWillEndFunc = f
	return b
}

func (b *FlowBuilder) OnWillWork(f func(workInfo *WorkInfo) (shouldContinue bool, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.WillWorkFunc = f
	return b
}

func (b *FlowBuilder) OnWorkSkip(f func(workInfo *WorkInfo, result Result, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.OnWorkSkipFunc = f
	return b
}

func (b *FlowBuilder) OnWorkSuccess(f func(workInfo *WorkInfo, result Result)) *FlowBuilder {
	b.flow.EventListener.OnWorkSuccessFunc = f
	return b
}

func (b *FlowBuilder) OnWorkFail(f func(workInfo *WorkInfo, err *data.CodeError)) *FlowBuilder {
	b.flow.EventListener.OnWorkFailFunc = f
	return b
}

type FlowBuilder struct {
	enableOverseer bool
	flow           *Flow
	err            error
}

func (b *FlowBuilder) Build() *Flow {
	if !b.enableOverseer {
		b.flow.Overseer = nil
	}

	if b.err != nil {
		log.ErrorF("Flow Builder error:%s", b.err)
	}
	return b.flow
}
