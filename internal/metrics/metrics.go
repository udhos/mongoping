// Package metrics abstracts metrics.
package metrics

import "time"

// Metrics is interface for recording metrics.
type Metrics interface {
	RecordLatency(target, uri, outcome string, latency time.Duration)
}
