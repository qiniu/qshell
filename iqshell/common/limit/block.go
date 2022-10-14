package limit

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
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
		startTimeOfRound: time.Now().Add(time.Second * -10),
	}
}

type blockLimit struct {
	mu               sync.RWMutex //
	limitCount       int          // qps 及并发限制数
	leftCount        int          // 可消费的数量
	indexOfRound     int          // 当前轮中已消费的号
	startTimeOfRound time.Time    // 当前轮开始时间
}

func (l *blockLimit) AddLimitCount(count int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if (l.limitCount + count) < 1 {
		count = l.limitCount - 1
	}
	l.limitCount += count
	l.leftCount = l.limitCount + count
	log.ErrorF("===== limitCount:%d", l.limitCount)
}

func (l *blockLimit) Acquire(count int) *data.CodeError {
	if count <= 0 {
		return nil
	}

	lCount := count
	for {
		if lCount <= 0 {
			break
		}

		if c := l.tryAcquire(lCount); c <= 0 {
			// 触及限制
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			lCount -= c
		}
	}
	return nil
}

func (l *blockLimit) tryAcquire(count int) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if count > l.limitCount {
		count = l.limitCount
	}

	// 没有余量，并发耗尽
	if l.leftCount < count {
		return 0
	}

	// 并发满足，查看 QPS 是否超标
	current := time.Now()
	if l.indexOfRound >= l.limitCount {
		if current.Add(time.Second * -1).Before(l.startTimeOfRound) {
			return 0
		}
		l.indexOfRound = 0
		l.startTimeOfRound = current
	}

	l.leftCount -= count
	l.indexOfRound += count

	return count
}

func (l *blockLimit) Release(count int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.leftCount += count
}
