package main

import (
	"time"
)

type noopMetrics struct {
}

func newMetricsNoop() *noopMetrics {
	return &noopMetrics{}
}

func (m *noopMetrics) RecordLatency(_ /*target*/, _ /*uri*/, _ /*outcome*/ string, _ /*latency*/ time.Duration) {
}
