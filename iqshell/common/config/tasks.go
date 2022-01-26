package config

import "github.com/qiniu/qshell/v2/iqshell/common/utils"

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed,omitempty"`
}

func (t *Tasks) merge(from *Tasks) {
	if from == nil {
		return
	}

	t.ConcurrentCount = utils.GetNotZeroIntIfExist(t.ConcurrentCount, from.ConcurrentCount)
	t.StopWhenOneTaskFailed = utils.GetNotEmptyStringIfExist(t.StopWhenOneTaskFailed, from.StopWhenOneTaskFailed)
}
