package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/udhos/dogstatsdclient/dogstatsdclient"
)

type appMetrics struct {
	latency         *prometheus.HistogramVec
	dogstatsdClient *dogstatsdclient.Client
}

func newMetrics(namespace string, latencyBuckets []float64,
	prometheusEnabled, dogstatsdEnabled, dogstatsdDebug bool) *appMetrics {
	const me = "newMetrics"

	m := &appMetrics{}

	if prometheusEnabled {
		//
		// prometheus
		//

		latency := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "requests_seconds",
				Help:      "Mongo ping request duration in seconds.",
				Buckets:   latencyBuckets,
			},
			[]string{"target", "uri", "outcome"},
		)

		if err := prometheus.Register(latency); err != nil {
			log.Fatalf("%s: latency was not registered: %s", me, err)
		}

		m.latency = latency
	}

	if dogstatsdEnabled {
		//
		// dogstatsd
		//

		options := dogstatsdclient.Options{
			Namespace: namespace,
			Debug:     dogstatsdDebug,
		}

		client, errClient := dogstatsdclient.New(options)
		if errClient != nil {
			log.Fatalf("%s: dogstatsd client error: %s", me, errClient)
		}

		m.dogstatsdClient = client
	}

	return m
}

func (m *appMetrics) RecordLatency(target, uri, outcome string, latency time.Duration) {
	if m.latency != nil {
		//
		// prometheus
		//
		m.latency.WithLabelValues(target, uri, outcome).Observe(latency.Seconds())
	}
	if m.dogstatsdClient != nil {
		//
		// dogstatsd
		//
		tags := []string{"target:" + target, "uri:" + uri, "outcome:" + outcome}
		m.dogstatsdClient.TimeInMilliseconds("requests_milliseconds", float64(latency.Milliseconds()), tags, 1)
	}
}
