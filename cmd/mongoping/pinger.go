package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/udhos/mongodbclient/mongodbclient"
)

func pinger(app *application) {
	const me = "pinger"
	for {
		for i, t := range app.targets {
			go pingTarget(i+1, len(app.targets), t, app.met, app.conf.timeout, app.conf.debug)
		}
		log.Printf("%s: sleeping for %v", me, app.conf.interval)
		time.Sleep(app.conf.interval)
	}
}

func pingTarget(i, max int, t target, met *metrics, timeout time.Duration, debug bool) {

	me := fmt.Sprintf("pingTarget[%d/%d]", i, max)

	log.Printf("%s: name=%s URL=%s timeout=%v", me, t.Name, t.URI, timeout)

	outcome := "unknown"
	var errPing error

	begin := time.Now()

	defer func() {
		elap := time.Since(begin)
		if errPing == nil {
			outcome = "success"
		} else {
			outcome = "error"
		}
		log.Printf("%s: name=%s URL=%s elapsed=%v outcome=%s error:%v",
			me, t.Name, t.URI, elap, outcome, errPing)
		met.recordLatency(t.Name, t.URI, outcome, elap)
	}()

	clientOptions := mongodbclient.Options{
		Debug:     debug,
		URI:       t.URI,
		Username:  t.User,
		Password:  t.Pass,
		TLSCAFile: t.TLSCaFile,
		Timeout:   timeout,
	}

	client, errClient := mongodbclient.New(clientOptions)
	if errClient != nil {
		errPing = errClient
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errPing = client.Ping(ctx, nil)
}
