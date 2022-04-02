package synchronized

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"sync"
)

type Locker interface {
	Lock()
	Unlock()
}

type Synchronized interface {
	Do(fn func())
	DoError(fn func() *data.CodeError) *data.CodeError
}

func NewSynchronized(locker Locker) Synchronized {
	if locker == nil {
		locker = &sync.Mutex{}
	}

	return &synchronized{
		locker: locker,
	}
}

type synchronized struct {
	locker Locker
}

func (s *synchronized) Do(fn func()) {
	if fn == nil {
		return
	}

	s.locker.Lock()
	fn()
	s.locker.Unlock()
	return
}

func (s *synchronized) DoError(fn func() *data.CodeError) (err *data.CodeError) {
	if fn == nil {
		return
	}

	s.locker.Lock()
	err = fn()
	s.locker.Unlock()
	return
}
