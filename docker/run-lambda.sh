
echo EXAMPLE:
echo "EXAMPLE: curl http://localhost:8080/2015-03-31/functions/function/invocations --data-binary @targets-lambda.json"
echo EXAMPLE:
echo

docker run --rm -p 8080:8080 \
    --entrypoint /usr/local/bin/aws-lambda-rie \
    -e TARGETS=/dev/null \
    udhos/mongoping-lambda:latest /bin/mongoping-lambda

