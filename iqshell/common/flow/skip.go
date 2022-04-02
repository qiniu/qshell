package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Skipper interface {
	ShouldSkip(work *WorkInfo) (skip bool, cause *data.CodeError)
}

func NewSkipper(f func(work *WorkInfo) (skip bool, cause *data.CodeError)) Skipper {
	return &skipper{f: f}
}

type skipper struct {
	f func(work *WorkInfo) (skip bool, cause *data.CodeError)
}

func (s *skipper) ShouldSkip(work *WorkInfo) (skip bool, cause *data.CodeError) {
	if s.f != nil {
		return false, nil
	}
	return s.f(work)
}
