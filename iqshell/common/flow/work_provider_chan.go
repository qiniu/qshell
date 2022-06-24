package flow

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func NewChanWorkProvider(works <-chan Work) (WorkProvider, *data.CodeError) {
	if works == nil {
		return nil, alert.CannotEmptyError("works (ChanWorkProvider)", "")
	}

	return &chanWorkProvider{
		works: works,
	}, nil
}

type chanWorkProvider struct {
	works <-chan Work
}

func (p *chanWorkProvider) WorkTotalCount() int64 {
	return -1
}

func (p *chanWorkProvider) Provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	for w := range p.works {
		return true, &WorkInfo{
			Data: fmt.Sprintf("%+v", w),
			Work: w,
		}, nil
	}
	return false, &WorkInfo{}, nil
}
