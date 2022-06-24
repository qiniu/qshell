package operations

import "github.com/qiniu/qshell/v2/iqshell/storage/object/batch"

type Metric struct {
	batch.Metric

	OverwriteCount    int64 `json:"overwrite_count"`
	NotOverwriteCount int64 `json:"not_overwrite_count"`
}

func (m *Metric) AddOverwriteCount(count int64) {
	if m == nil {
		return
	}
	m.Lock()
	m.OverwriteCount += count
	m.Unlock()
}

func (m *Metric) AddNotOverwriteCount(count int64) {
	m.Lock()
	m.NotOverwriteCount += count
	m.Unlock()
}
