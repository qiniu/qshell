package limit

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"sync"
	"time"
)

// BlockLimit 并发 + 限流
type BlockLimit interface {
	Limit
	AddLimitCount(limitCount int)
}

func NewBlockList(limitCount int) BlockLimit {
	if limitCount <= 0 {
		limitCount = 0
	}

	return &blockLimit{
		mu:               sync.RWMutex{},
		limitCount:       limitCount,
		leftCount:        limitCount,
		indexOfRound:     0,
		startTimeOfRound: time.Now().Add(time.Second * -1),
	}
}

type blockLimit struct {
	mu               sync.RWMutex //
	limitCount       int          // qps 及并发限制数
	leftCount        int          // 可消费的数量
	indexOfRound     int          // 当前轮中已消费的号
	startTimeOfRound time.Time    // 当前轮开始时间
}

func (l *blockLimit) AddLimitCount(limitCount int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limitCount += limitCount
	if l.limitCount < 1 {
		l.limitCount = 1
	}
}

func (l *blockLimit) Acquire(count int) *data.CodeError {
	if count <= 0 {
		return nil
	}

	leftCount := count
	for {
		if leftCount <= 0 {
			break
		}

		if c := l.tryAcquire(leftCount); c <= 0 {
			// 触及限制
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			leftCount -= c
		}
	}
	return nil
}

func (l *blockLimit) tryAcquire(count int) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	if count > l.limitCount {
		count = l.limitCount
	}

	if l.leftCount < count {
		// 并发耗尽
		return 0
	}

	current := time.Now()
	if current.Add(time.Second * -1).Before(l.startTimeOfRound) {
		// qps 耗尽
		return 0
	}

	l.leftCount -= count
	if l.leftCount < 0 {
		l.leftCount = 0
	}

	l.indexOfRound += 1
	if l.indexOfRound >= l.limitCount {
		l.indexOfRound = 0
		l.startTimeOfRound = current
	}
	return count
}

func (l *blockLimit) Release(count int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.leftCount += count
	if l.leftCount > l.limitCount {
		l.leftCount = l.limitCount
	}
}
