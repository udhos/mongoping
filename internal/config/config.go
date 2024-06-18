// Package config loads configuration.
package config

import (
	"log"
	"os"
	"time"

	"github.com/udhos/boilerplate/secret"
	"github.com/udhos/mongoping/internal/env"
	"gopkg.in/yaml.v3"
)

// Version is program version.
const Version = "1.2.2"

// Config holds program configuration.
type Config struct {
	SecretRoleArn         string
	Targets               string
	Interval              time.Duration
	Timeout               time.Duration
	MetricsAddr           string
	MetricsPath           string
	MetricsNamespace      string
	MetricsLatencyBuckets []float64
	HealthAddr            string
	HealthPath            string
	Debug                 bool
}

// GetConfig loads configuration.
func GetConfig() Config {
	return Config{
		SecretRoleArn:         env.String("SECRET_ROLE_ARN", ""),
		Targets:               env.String("TARGETS", "targets.yaml"),
		Interval:              env.Duration("INTERVAL", 10*time.Second),
		Timeout:               env.Duration("TIMEOUT", 5*time.Second),
		MetricsAddr:           env.String("METRICS_ADDR", ":3000"),
		MetricsPath:           env.String("METRICS_PATH", "/metrics"),
		MetricsNamespace:      env.String("METRICS_NAMESPACE", ""),
		MetricsLatencyBuckets: env.Float64Slice("METRICS_BUCKETS_LATENCY", []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, .5, 1}),
		HealthAddr:            env.String("HEALTH_ADDR", ":8888"),
		HealthPath:            env.String("HEALTH_PATH", "/health"),
		Debug:                 env.Bool("DEBUG", false),
	}
}

// Target holds ping target.
type Target struct {
	Name      string `yaml:"name"`
	URI       string `yaml:"uri"`
	Cmd       string `yaml:"cmd"`
	Database  string `yaml:"database"` // command hello requires database
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	TLSCaFile string `yaml:"tls_ca_file"`
	RoleArn   string `yaml:"role_arn"`
}

// LoadTargets load targets from file.
func LoadTargets(targetsFile, sessionName, secretRoleArn string) []Target {
	const me = "loadTargets"

	log.Printf("%s: file=%s session=%s role=%s", me,
		targetsFile, sessionName, secretRoleArn)

	var targets []Target

	buf, errRead := os.ReadFile(targetsFile)
	if errRead != nil {
		log.Fatalf("%s: load targets: %s: %v",
			me, targetsFile, errRead)
	}

	errYaml := yaml.Unmarshal(buf, &targets)
	if errYaml != nil {
		log.Fatalf("%s: parse targets yaml: %s: %v",
			me, targetsFile, errYaml)
	}

	// get secret using global role
	sec := secret.New(secret.Options{
		RoleSessionName: sessionName,
		RoleArn:         secretRoleArn,
	})

	for _, t := range targets {

		if t.RoleArn != "" {
			//
			// non-empty per-target role overrides global role
			//
			s := secret.New(secret.Options{
				RoleSessionName: sessionName,
				RoleArn:         t.RoleArn,
			})
			t.Pass = s.Retrieve(t.Pass)
			continue
		}

		t.Pass = sec.Retrieve(t.Pass) // get secret using global role
	}

	return targets
}
