package batch

type Metric struct {
	TotalCount   int64 `json:"total_count"`
	SuccessCount int64 `json:"success_count"`
	SkippedCount int64 `json:"skipped_count"`
	FailureCount int64 `json:"failure_count"`
}
