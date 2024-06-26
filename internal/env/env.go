// Package env loads configuration from env vars.
package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// String extracts string from env var.
// It returns the provided defaultValue if the env var is empty.
// The string returned is also recorded in logs.
func String(name string, defaultValue string) string {
	str := os.Getenv(name)
	if str != "" {
		log.Printf("%s=[%s] using %s=%s default=%s", name, str, name, str, defaultValue)
		return str
	}
	log.Printf("%s=[%s] using %s=%s default=%s", name, str, name, defaultValue, defaultValue)
	return defaultValue
}

// Duration extracts time.Duration value from env var.
// It returns the provided defaultValue if the env var is empty.
// The value returned is also recorded in logs.
func Duration(name string, defaultValue time.Duration) time.Duration {
	str := os.Getenv(name)
	if str != "" {
		value, errConv := time.ParseDuration(str)
		if errConv == nil {
			log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, value, defaultValue)
			return value
		}
		log.Printf("bad %s=[%s]: error: %v", name, str, errConv)
	}
	log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, defaultValue, defaultValue)
	return defaultValue
}

// Float64Slice extracts []float64 from env var.
// It returns the provided defaultValue if the env var is empty.
// The value returned is also recorded in logs.
func Float64Slice(name string, defaultValue []float64) []float64 {
	str := os.Getenv(name)
	if str == "" {
		log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, defaultValue, defaultValue)
		return defaultValue
	}

	var value []float64
	items := strings.FieldsFunc(str, func(sep rune) bool { return sep == ',' })
	for i, field := range items {
		field = strings.TrimSpace(field)
		f, errConv := strconv.ParseFloat(field, 64)
		if errConv != nil {
			log.Printf("bad %s=[%s] error parsing item %d='%s': %v: using %s=%v default=%v",
				name, str, i, field, errConv, name, value, defaultValue)
			return defaultValue
		}
		value = append(value, f)
	}

	log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, value, defaultValue)

	return value
}

// Bool extracts boolean value from env var.
// It returns the provided defaultValue if the env var is empty.
// The value returned is also recorded in logs.
func Bool(name string, defaultValue bool) bool {
	str := os.Getenv(name)
	if str != "" {
		value, errConv := strconv.ParseBool(str)
		if errConv == nil {
			log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, value, defaultValue)
			return value
		}
		log.Printf("bad %s=[%s]: error: %v", name, str, errConv)
	}
	log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, defaultValue, defaultValue)
	return defaultValue
}
