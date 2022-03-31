package flow

import "sync"

func NewArrayWorkProvider(works []Work) WorkProvider {
	return &arrayWorkProvider{
		readOffset: 0,
		works:      works,
	}
}

type arrayWorkProvider struct {
	mu         sync.Mutex
	readOffset int
	works      []Work
}

func (a *arrayWorkProvider) Provide() (hasMore bool, work Work, err error) {
	a.mu.Lock()
	hasMore, work, err = a.provide()
	a.mu.Unlock()
	return
}

func (a *arrayWorkProvider) provide() (hasMore bool, work Work, err error) {
	if a.readOffset > len(a.works)-1 {
		return false, nil, nil
	}
	hasMore = true
	work = a.works[a.readOffset]
	a.readOffset ++
	return
}
