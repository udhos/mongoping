[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/mongoping/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/mongoping)](https://goreportcard.com/report/github.com/udhos/mongoping)
[![Go Reference](https://pkg.go.dev/badge/github.com/udhos/mongoping.svg)](https://pkg.go.dev/github.com/udhos/mongoping)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/mongoping)](https://artifacthub.io/packages/search?repo=mongoping)
[![Docker Pulls mongoping](https://img.shields.io/docker/pulls/udhos/mongoping)](https://hub.docker.com/r/udhos/mongoping)

# mongoping

mongoping

# Env vars

```
export SECRET_ROLE_ARN=""
export TARGETS=targets.yaml
export INTERVAL=10s
export TIMEOUT=5s
export METRICS_ADDR=:3000
export METRICS_PATH=/metrics
export METRICS_NAMESPACE=""
export METRICS_BUCKETS_LATENCY="0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, .5, 1"
export HEALTH_ADDR=:8888
export HEALTH_PATH=/health
```

# Config file

Config file is defined by env var `TARGETS=targets.yaml`.

```
$ cat targets.yaml
- name: "mongo1"
  uri: "mongodb://localhost:27017"
  #tls_ca_file: ca-bundle.pem
- name: "mongo2"
  uri: "mongodb://localhost:27018"
  user: user2
  pass: aws-parameterstore:us-east-1:mongo_pass_user2 # Retrieve from parameter store
  role_arn: arn:aws:iam::100010001000:role/admin
  #tls_ca_file: ca-bundle.pem
```

# Docker

Docker hub:

https://hub.docker.com/r/udhos/mongoping

Run from docker hub:

```
docker run -p 8080:8080 --rm udhos/mongoping:0.1.0
```

Build recipe:

```
./docker/build.sh

docker push udhos/mongoping:0.1.0
```

# Helm chart

You can use the provided helm charts to install mongoping in your Kubernetes cluster.

See: https://udhos.github.io/mongoping/

## Lint

    helm lint ./charts/mongoping --values charts/mongoping/values.yaml

## Debug

    helm template ./charts/mongoping --values charts/mongoping/values.yaml --debug

## Render at server

    helm install my-mongoping ./charts/mongoping --values charts/mongoping/values.yaml --dry-run

## Install

    helm install my-mongoping ./charts/mongoping --values charts/mongoping/values.yaml

    helm list -A
