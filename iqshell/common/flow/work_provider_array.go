package flow

import (
	"sync"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func NewArrayWorkProvider(works []Work) (WorkProvider, *data.CodeError) {
	if works == nil {
		return nil, alert.CannotEmptyError("works (ArrayWorkProvider)", "")
	}

	return &arrayWorkProvider{
		readOffset: 0,
		works:      createWorkInfoListWithWorkList(works),
	}, nil
}

type arrayWorkProvider struct {
	mu         sync.Mutex
	readOffset int
	works      []*WorkInfo
}

func (p *arrayWorkProvider) WorkTotalCount() int64 {
	return int64(len(p.works))
}

func (p *arrayWorkProvider) Provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	p.mu.Lock()
	hasMore, work, err = p.provide()
	p.mu.Unlock()
	return
}

func (p *arrayWorkProvider) provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	if p.readOffset > len(p.works)-1 {
		return false, &WorkInfo{}, nil
	}
	hasMore = true
	work = p.works[p.readOffset]
	p.readOffset++
	return
}

func createWorkInfoListWithWorkList(works []Work) []*WorkInfo {
	if works == nil {
		return nil
	}

	infos := make([]*WorkInfo, 0, len(works))
	for _, work := range works {
		infos = append(infos, &WorkInfo{
			Data: "",
			Work: work,
		})
	}
	return infos
}
