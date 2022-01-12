package config

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed,omitempty"`
}
