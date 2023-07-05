#!/bin/bash

version=$(go run ./cmd/mongoping -version | awk '{ print $2 }' | awk -F= '{ print $2 }')

echo version=$version

docker build --no-cache \
    -t udhos/mongoping:latest \
    -t udhos/mongoping:$version \
    -f docker/Dockerfile .

echo "push: docker push udhos/mongoping:$version; docker push udhos/mongoping:latest"
