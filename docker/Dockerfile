# STEP 1 build executable binary

FROM golang:1.22.4-alpine3.20 as builder

RUN apk update
RUN apk add git

# Create appuser on builder image
RUN adduser -D -g '' appuser

COPY ./cmd/ /tmp/app/cmd
COPY ./internal/ /tmp/app/internal
COPY go.* /tmp/app/
WORKDIR /tmp/app
RUN go mod tidy
RUN go env -w CGO_ENABLED=0
RUN go build -o /tmp/mongoping github.com/udhos/mongoping/cmd/mongoping

#
# STEP 2 build a small image from alpine
#
# curl: https://security.alpinelinux.org/srcpkg/curl
#
FROM alpine:3.20.0
COPY --from=builder /tmp/mongoping /bin/mongoping
#RUN apk add curl=8.1.2-r0 libcrypto3=3.1.0-r4 libssl3=3.1.0-r4
RUN apk update
RUN apk add curl
RUN apk upgrade libssl3 libcrypto3 busybox busybox-binsh ssl_client libcurl nghttp2-libs
RUN adduser -D -g '' user
USER user
ENTRYPOINT ["/bin/mongoping"]
