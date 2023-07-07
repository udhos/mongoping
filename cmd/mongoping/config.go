package main

import "time"

type config struct {
	secretRoleArn         string
	targets               string
	interval              time.Duration
	timeout               time.Duration
	metricsAddr           string
	metricsPath           string
	metricsNamespace      string
	metricsLatencyBuckets []float64
	healthAddr            string
	healthPath            string
	debug                 bool
}

func getConfig() config {
	return config{
		secretRoleArn:         envString("SECRET_ROLE_ARN", ""),
		targets:               envString("TARGETS", "targets.yaml"),
		interval:              envDuration("INTERVAL", 10*time.Second),
		timeout:               envDuration("TIMEOUT", 5*time.Second),
		metricsAddr:           envString("METRICS_ADDR", ":3000"),
		metricsPath:           envString("METRICS_PATH", "/metrics"),
		metricsNamespace:      envString("METRICS_NAMESPACE", ""),
		metricsLatencyBuckets: envFloat64Slice("METRICS_BUCKETS_LATENCY", []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, .5, 1}),
		healthAddr:            envString("HEALTH_ADDR", ":8888"),
		healthPath:            envString("HEALTH_PATH", "/health"),
		debug:                 envBool("DEBUG", false),
	}
}
