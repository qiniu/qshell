package synchronized

import (
	"sync"
)

type Locker interface {
	Lock()
	Unlock()
}

type Synchronized interface {
	Do(fn func())
	DoError(fn func() error) error
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

func (s *synchronized) Do(fn func()){
	if fn == nil {
		return
	}

	s.locker.Lock()
	fn()
	s.locker.Unlock()
	return
}

func (s *synchronized) DoError(fn func() error) (err error) {
	if fn == nil {
		return
	}

	s.locker.Lock()
	err = fn()
	s.locker.Unlock()
	return
}