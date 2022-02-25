package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Tasks struct {
	ConcurrentCount       *data.Int  `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed *data.Bool `json:"stop_when_one_task_failed,omitempty"`
}

func (t *Tasks) merge(from *Tasks) {
	if from == nil {
		return
	}

	t.ConcurrentCount = data.GetNotEmptyIntIfExist(t.ConcurrentCount, from.ConcurrentCount)
	t.StopWhenOneTaskFailed = data.GetNotEmptyBoolIfExist(t.StopWhenOneTaskFailed, from.StopWhenOneTaskFailed)
}
