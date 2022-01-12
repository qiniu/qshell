package config

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed,omitempty"`
}

func (t *Tasks) merge(from *Tasks) {
	if from == nil {
		return
	}

	if t.ConcurrentCount == 0 {
		t.ConcurrentCount = from.ConcurrentCount
	}

	if len(t.StopWhenOneTaskFailed) == 0 {
		t.StopWhenOneTaskFailed = from.StopWhenOneTaskFailed
	}
}