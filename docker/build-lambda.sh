#!/bin/bash

app=mongoping-lambda

version=$(go run ./cmd/mongoping -version | awk '{ print $2 }' | awk -F= '{ print $2 }')

echo version=$version

rm -f global-bundle.pem
wget https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem

docker build --no-cache \
    -t udhos/$app:latest \
    -t udhos/$app:$version \
    -f docker/Dockerfile.lambda .

echo push:
push=docker-push-lambda.sh
echo "docker push udhos/$app:$version; docker push udhos/$app:latest" > $push
chmod a+rx $push
echo $push:
cat $push
