package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/limit"
	"sync"
)

type AutoLimit interface {
	limit.BlockLimit

	IsLimitError(code int, err *data.CodeError) bool
}

func NewAutoLimit(limitCount, maxLimitCount, minLimitCount int64) AutoLimit {
	if limitCount < 1 {
		limitCount = 1
	}
	if maxLimitCount < limitCount {
		maxLimitCount = 0
	}
	if minLimitCount > limitCount {
		minLimitCount = 0
	}
	return &autoLimit{
		mu:            sync.RWMutex{},
		blockLimit:    limit.NewBlockList(limitCount),
		limitCount:    limitCount,
		maxLimitCount: maxLimitCount,
		minLimitCount: minLimitCount,
	}
}

type autoLimit struct {
	mu            sync.RWMutex
	blockLimit    limit.BlockLimit //
	limitCount    int64            // qps 及并发限制数
	maxLimitCount int64            // 最大限制数
	minLimitCount int64            // 做小限制数
}

func (l *autoLimit) Acquire(count int64) *data.CodeError {
	return l.blockLimit.Acquire(count)
}

func (l *autoLimit) Release(count int64) {
	l.blockLimit.Release(count)
}

func (l *autoLimit) AddLimitCount(count int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if count == 0 {
		return
	}

	if l.maxLimitCount > 0 && l.limitCount+count > l.maxLimitCount {
		count = l.maxLimitCount - l.limitCount
	}

	if l.minLimitCount > 0 && l.limitCount+count < l.minLimitCount {
		count = l.limitCount - l.minLimitCount
	}

	l.limitCount += count
	l.blockLimit.AddLimitCount(count)
}

func (l *autoLimit) IsLimitError(code int, err *data.CodeError) bool {
	if code == 573 {
		return true
	}
	return false
}
