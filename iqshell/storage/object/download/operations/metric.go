package operations

import "github.com/qiniu/qshell/v2/iqshell/storage/object/batch"

type Metric struct {
	batch.Metric

	ExistCount  int64 `json:"exist_count"`
	UpdateCount int64 `json:"update_count"`
}

func (m *Metric) AddExistCount(count int64) {
	if m == nil {
		return
	}
	m.Lock()
	m.ExistCount += count
	m.Unlock()
}

func (m *Metric) AddUpdateCount(count int64) {
	m.Lock()
	m.UpdateCount += count
	m.Unlock()
}