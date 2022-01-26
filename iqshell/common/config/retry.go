package config

import "github.com/qiniu/qshell/v2/iqshell/common/utils"

type Retry struct {
	Max      int `json:"max,omitempty"`
	Interval int `json:"interval,omitempty"`
}

func (r *Retry) merge(from *Retry) {
	if from == nil {
		return
	}

	r.Max = utils.GetNotZeroIntIfExist(r.Max, from.Max)
	r.Interval = utils.GetNotZeroIntIfExist(r.Interval, from.Interval)
}
