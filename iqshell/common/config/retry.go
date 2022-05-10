package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Retry struct {
	Max      *data.Int `json:"max,omitempty"`
	Interval *data.Int `json:"interval,omitempty"`
}

func (r *Retry) merge(from *Retry) {
	if from == nil {
		return
	}

	r.Max = data.GetNotEmptyIntIfExist(r.Max, from.Max)
	r.Interval = data.GetNotEmptyIntIfExist(r.Interval, from.Interval)
}
