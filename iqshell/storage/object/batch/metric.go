package batch

import (
	"sync"
	"time"
)

type Metric struct {
	mu    sync.Mutex
	start time.Time

	Duration     int64 `json:"duration"`
	TotalCount   int64 `json:"total_count"`
	CurrentCount int64 `json:"-"`
	SuccessCount int64 `json:"success_count"`
	FailureCount int64 `json:"failure_count"`
	SkippedCount int64 `json:"skipped_count"`
}

func (m *Metric) Start() {
	if m == nil {
		return
	}
	m.start = time.Now()
}

func (m *Metric) End() {
	if m == nil {
		return
	}
	sUnix :=  m.start.Unix()
	eUnix := time.Now().Unix()
	m.Duration = eUnix - sUnix
}

func (m *Metric) AddTotalCount(count int64) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.TotalCount += count
	m.mu.Unlock()
}

func (m *Metric) AddCurrentCount(count int64) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.CurrentCount += count
	m.mu.Unlock()
}


func (m *Metric) AddSuccessCount(count int64) {
	m.mu.Lock()
	m.SuccessCount += count
	m.mu.Unlock()
}

func (m *Metric) AddFailureCount(count int64) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.FailureCount += count
	m.mu.Unlock()
}

func (m *Metric) AddSkippedCount(count int64) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.SkippedCount += count
	m.mu.Unlock()
}

func (m *Metric) Lock() {
	if m == nil {
		return
	}
	m.mu.Lock()
}

func (m *Metric) Unlock() {
	if m == nil {
		return
	}
	m.mu.Unlock()
}
