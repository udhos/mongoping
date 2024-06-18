// Package main implements the mongoping tool.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/udhos/mongoping/internal/config"
	"github.com/udhos/mongoping/internal/metrics"
	"github.com/udhos/mongoping/internal/ping"
	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	me            string
	conf          config.Config
	targets       []config.Target
	serverMetrics *http.Server
	serverHealth  *http.Server
	met           metrics.Metrics
}

func longVersion(me string) string {
	return fmt.Sprintf("%s runtime=%s GOOS=%s GOARCH=%s GOMAXPROCS=%d",
		me, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
}

func main() {

	//
	// parse cmd line
	//

	var showVersion bool
	flag.BoolVar(&showVersion, "version", showVersion, "show version")
	flag.Parse()

	//
	// show version
	//

	me := filepath.Base(os.Args[0])

	{
		v := longVersion(me + " version=" + config.Version)
		if showVersion {
			fmt.Println(v)
			return
		}
		log.Print(v)
	}

	app := &application{
		me:   me,
		conf: config.GetConfig(),
	}

	app.targets = config.LoadTargets(app.conf.Targets, me, app.conf.SecretRoleArn)

	//
	// start metrics server
	//

	{
		app.met = newMetricsPrometheus(app.conf.MetricsNamespace, app.conf.MetricsLatencyBuckets)

		mux := http.NewServeMux()
		app.serverMetrics = &http.Server{
			Addr:    app.conf.MetricsAddr,
			Handler: mux,
		}

		mux.Handle(app.conf.MetricsPath, promhttp.Handler())

		go func() {
			log.Printf("metrics server: listening on %s %s", app.conf.MetricsAddr, app.conf.MetricsPath)
			err := app.serverMetrics.ListenAndServe()
			log.Fatalf("metrics server: exited: %v", err)
		}()
	}

	//
	// start health server
	//

	{
		mux := http.NewServeMux()
		app.serverHealth = &http.Server{
			Addr:    app.conf.HealthAddr,
			Handler: mux,
		}

		mux.HandleFunc(app.conf.HealthPath, func(w http.ResponseWriter, _ /*r*/ *http.Request) {
			http.Error(w, "health ok", 200)
		})

		go func() {
			log.Printf("health server: listening on %s %s", app.conf.HealthAddr, app.conf.HealthPath)
			err := app.serverHealth.ListenAndServe()
			log.Fatalf("health server: exited: %v", err)
		}()
	}

	//
	// start pinger
	//

	go pinger(app)

	<-make(chan struct{}) // wait forever
}

func pinger(app *application) {
	const me = "pinger"

	clients := make([]*mongo.Client, len(app.targets))

	for {
		for i, t := range app.targets {
			go ping.Ping(clients, i, len(app.targets), t, app.met, app.conf.Timeout, app.conf.Debug)
		}
		if app.conf.Debug {
			log.Printf("%s: sleeping for %v", me, app.conf.Interval)
		}
		time.Sleep(app.conf.Interval)
	}
}
