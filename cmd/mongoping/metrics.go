package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	latency *prometheus.HistogramVec
}

func newMetricsPrometheus(namespace string, latencyBuckets []float64) *prometheusMetrics {
	const me = "newMetrics"

	//
	// latency
	//

	latency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "ping_requests_seconds",
			Help:      "Mongo ping request duration in seconds.",
			Buckets:   latencyBuckets,
		},
		[]string{"target", "uri", "outcome"},
	)

	if err := prometheus.Register(latency); err != nil {
		log.Fatalf("%s: latency was not registered: %s", me, err)
	}

	//
	// all metrics
	//

	m := &prometheusMetrics{
		latency: latency,
	}

	return m
}

func (m *prometheusMetrics) RecordLatency(target, uri, outcome string, latency time.Duration) {
	m.latency.WithLabelValues(target, uri, outcome).Observe(float64(latency) / float64(time.Second))
}
