package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/udhos/mongodbclient/mongodbclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func pinger(app *application) {
	const me = "pinger"

	clients := make([]*mongo.Client, len(app.targets))

	for {
		for i, t := range app.targets {
			go pingTarget(clients, i, len(app.targets), t, app.met, app.conf.timeout, app.conf.debug)
		}
		if app.conf.debug {
			log.Printf("%s: sleeping for %v", me, app.conf.interval)
		}
		time.Sleep(app.conf.interval)
	}
}

func pingTarget(clients []*mongo.Client, i, max int, t target, met *metrics, timeout time.Duration, debug bool) {

	me := fmt.Sprintf("pingTarget[%d/%d] cmd=[%s]", i+1, max, t.Cmd)

	if debug {
		log.Printf("%s: name=%s URL=%s timeout=%v", me, t.Name, t.URI, timeout)
	}

	var errPing error

	begin := time.Now()

	defer func() {
		elap := time.Since(begin)
		var outcome string
		if errPing == nil {
			outcome = "success"
			if debug {
				log.Printf("%s: name=%s URL=%s elapsed=%v outcome=%s",
					me, t.Name, t.URI, elap, outcome)
			}
		} else {
			outcome = "error"
			log.Printf("%s: name=%s URL=%s elapsed=%v outcome=%s error:%v",
				me, t.Name, t.URI, elap, outcome, errPing)
		}
		met.recordLatency(t.Name, t.URI, outcome, elap)
	}()

	if clients[i] == nil {
		//
		// create new client
		//
		clientOptions := mongodbclient.Options{
			Debug:     debug,
			URI:       t.URI,
			Username:  t.User,
			Password:  t.Pass,
			TLSCAFile: t.TLSCaFile,
			Timeout:   timeout,
		}
		c, errClient := mongodbclient.New(clientOptions)
		if errClient != nil {
			errPing = errClient
			return
		}
		clients[i] = c // save client
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if t.Cmd == "hello" {
		//
		// special case "hello"
		//
		db := clients[i].Database(t.Database)
		command := bson.D{{Key: "hello"}}
		opts := options.RunCmd()
		var result bson.M
		errPing = db.RunCommand(ctx, command, opts).Decode(&result)
		if debug && errPing == nil {
			log.Printf("%s: name=%s URL=%s hello result: %v",
				me, t.Name, t.URI, result)
		}
		return
	}

	// default to "ping"

	errPing = clients[i].Ping(ctx, nil)
}
