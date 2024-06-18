// Package main implements the mongoping tool.
package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/udhos/mongoping/internal/config"
	"github.com/udhos/mongoping/internal/metrics"
	"github.com/udhos/mongoping/internal/ping"
	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	me      string
	conf    config.Config
	targets []config.Target
	met     metrics.Metrics
}

func longVersion(me string) string {
	return fmt.Sprintf("%s runtime=%s GOOS=%s GOARCH=%s GOMAXPROCS=%d",
		me, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
}

const me = "mongoping-lambda"

var (
	app     *application
	clients []*mongo.Client
)

func init() {

	app = &application{
		me:   me,
		conf: config.GetConfig(),
	}

	app.targets = config.LoadTargets(app.conf.Targets, me, app.conf.SecretRoleArn)

	app.met = newMetricsNoop()

	clients = make([]*mongo.Client, len(app.targets))
}

type pingEvent struct {
	Targets []config.Target `yaml:"targets"`
}

// HandleRequest is lambda handler.
func HandleRequest(_ /*ctx*/ context.Context, event *pingEvent) (string, error) {

	//
	// show version
	//

	const me = "mongoping-lambda"

	log.Print(longVersion(me + " version=" + config.Version))

	log.Printf("targets from file: %d", len(app.targets))

	targets := app.targets // defaults to targets file from env var TARGETS

	if event == nil {
		log.Print("received nil event")
	} else {
		if len(event.Targets) == 0 {
			log.Print("event has zero targets")
		} else {
			log.Printf("using targets from event: %d", len(event.Targets))
			targets = event.Targets // event targets override env var TARGETS
		}
	}

	log.Printf("targets: %d", len(targets))

	var wg sync.WaitGroup

	for i, t := range targets {
		wg.Add(1)
		go func() {
			ping.Ping(clients, i, len(targets), t, app.met, app.conf.Timeout, app.conf.Debug)
			wg.Done()
		}()
	}

	wg.Wait()

	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
