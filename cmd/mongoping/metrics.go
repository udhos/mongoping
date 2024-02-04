package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	latency *prometheus.HistogramVec
}

func newMetrics(namespace string, latencyBuckets []float64) *metrics {
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

	m := &metrics{
		latency: latency,
	}

	return m
}

func (m *metrics) recordLatency(target, uri, outcome string, latency time.Duration) {
	m.latency.WithLabelValues(target, uri, outcome).Observe(float64(latency) / float64(time.Second))
}
