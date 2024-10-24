// Package ping implements the pinger.
package ping

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/udhos/mongodbclient/mongodbclient"
	"github.com/udhos/mongoping/internal/config"
	"github.com/udhos/mongoping/internal/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ping pings one target.
func Ping(clients []*mongo.Client, i, targets int, t config.Target, met metrics.Metrics, timeout time.Duration, debug bool) {

	me := fmt.Sprintf("Ping[%d/%d] cmd=[%s]", i+1, targets, t.Cmd)

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
		met.RecordLatency(t.Name, t.URI, outcome, elap)
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
