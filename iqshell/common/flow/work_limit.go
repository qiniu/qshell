package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/limit"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type AutoLimitOption func(l *autoLimit)

func MaxLimitCount(count int) AutoLimitOption {
	return func(l *autoLimit) {
		l.maxLimitCount = count
	}
}

func MinLimitCount(count int) AutoLimitOption {
	return func(l *autoLimit) {
		l.minLimitCount = count
	}
}

func IncreaseLimitCount(count int) AutoLimitOption {
	return func(l *autoLimit) {
		l.increaseLimitCount = count
	}
}

func IncreaseLimitCountPeriod(period time.Duration) AutoLimitOption {
	return func(l *autoLimit) {
		l.increaseLimitCountPeriod = period
	}
}

func NewBlockLimit(limitCount int, options ...AutoLimitOption) limit.BlockLimit {
	l := &autoLimit{
		mu:                       sync.RWMutex{},
		blockLimit:               limit.NewBlockList(limitCount),
		limitCount:               limitCount,
		maxLimitCount:            0,
		minLimitCount:            0,
		increaseLimitCountPeriod: 30 * time.Second,
		lastLimitCountChangeTime: time.Now(),
		increaseLimitCount:       10,
	}
	for _, option := range options {
		option(l)
	}
	l.check()
	return l
}

type autoLimit struct {
	mu                       sync.RWMutex
	blockLimit               limit.BlockLimit //
	limitCount               int              // qps 及并发限制数
	maxLimitCount            int              // 最大限制数
	minLimitCount            int              // 做小限制数
	increaseLimitCountPeriod time.Duration    // 增长检测周期
	lastLimitCountChangeTime time.Time        // 上次减小限制数的时间
	increaseLimitCount       int              // 增加幅度
	shouldWait               bool             //
	notReleaseCount          int64            //
}

func (l *autoLimit) check() {
	if l.limitCount < 1 {
		l.limitCount = 1
	}

	if l.maxLimitCount > 0 && l.maxLimitCount < l.limitCount {
		// 上限尽可能小
		l.limitCount = l.maxLimitCount
	}

	if l.minLimitCount > 0 && l.minLimitCount > l.limitCount {
		// 下限尽可能不设置
		l.minLimitCount = 0
	}
}

func (l *autoLimit) Acquire(count int) *data.CodeError {
	// 在 acquire 的时候尝试增加一次 limit count
	if l.shouldAutoIncreaseLimitCount() {
		l.AddLimitCount(l.increaseLimitCount)
	}

	l.waitIfNeeded()

	atomic.AddInt64(&l.notReleaseCount, int64(count))

	return l.blockLimit.Acquire(count)
}

func (l *autoLimit) Release(count int) {
	atomic.AddInt64(&l.notReleaseCount, int64(-count))
	l.blockLimit.Release(count)
}

func (l *autoLimit) AddLimitCount(count int) {
	if count == 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if count < 0 {
		l.shouldWait = true
	}

	l.lastLimitCountChangeTime = time.Now()

	if l.maxLimitCount > 0 && l.limitCount+count > l.maxLimitCount {
		count = l.maxLimitCount - l.limitCount
	}
	if l.minLimitCount > 0 && l.limitCount+count < l.minLimitCount {
		count = l.limitCount - l.minLimitCount
	}
	l.limitCount += count

	l.blockLimit.AddLimitCount(count)
}

func (l *autoLimit) shouldAutoIncreaseLimitCount() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.maxLimitCount > 0 && l.limitCount >= l.maxLimitCount {
		return false
	}

	if l.shouldWait {
		return false
	}

	return l.lastLimitCountChangeTime.Before(time.Now().Add(-1 * l.increaseLimitCountPeriod))
}

func (l *autoLimit) waitIfNeeded() {
	waitTime := time.Millisecond * time.Duration(rand.Int31n(1000)+500)
	for {
		if !l.shouldWait {
			break
		}

		if l.notReleaseCount <= (int64(l.limitCount) / 3) {
			l.shouldWait = false
		}

		time.Sleep(waitTime)
	}
}
