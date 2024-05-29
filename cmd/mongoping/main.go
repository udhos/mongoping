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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/udhos/boilerplate/secret"
	"gopkg.in/yaml.v3"
)

const version = "1.1.10"

type application struct {
	me            string
	conf          config
	targets       []target
	serverMetrics *http.Server
	serverHealth  *http.Server
	met           *metrics
}

type target struct {
	Name      string `yaml:"name"`
	URI       string `yaml:"uri"`
	Cmd       string `yaml:"cmd"`
	Database  string `yaml:"database"` // command hello requires database
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	TLSCaFile string `yaml:"tls_ca_file"`
	RoleArn   string `yaml:"role_arn"`
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
		v := longVersion(me + " version=" + version)
		if showVersion {
			fmt.Println(v)
			return
		}
		log.Print(v)
	}

	app := &application{
		me:   me,
		conf: getConfig(),
	}

	app.targets = loadTargets(app.conf.targets, me, app.conf.secretRoleArn)

	//
	// start metrics server
	//

	{
		app.met = newMetrics(app.conf.metricsNamespace, app.conf.metricsLatencyBuckets)

		mux := http.NewServeMux()
		app.serverMetrics = &http.Server{
			Addr:    app.conf.metricsAddr,
			Handler: mux,
		}

		mux.Handle(app.conf.metricsPath, promhttp.Handler())

		go func() {
			log.Printf("metrics server: listening on %s %s", app.conf.metricsAddr, app.conf.metricsPath)
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
			Addr:    app.conf.healthAddr,
			Handler: mux,
		}

		mux.HandleFunc(app.conf.healthPath, func(w http.ResponseWriter, _ /*r*/ *http.Request) {
			http.Error(w, "health ok", 200)
		})

		go func() {
			log.Printf("health server: listening on %s %s", app.conf.healthAddr, app.conf.healthPath)
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

func loadTargets(targetsFile, sessionName, secretRoleArn string) []target {
	const me = "loadTargets"

	var targets []target

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
