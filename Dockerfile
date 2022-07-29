# Builder
FROM golang:1.18-alpine as builder

RUN         apk add --no-cache gcc g++ make git

WORKDIR     /app

COPY        go.mod ./

COPY        go.sum ./

RUN         go mod download

ARG         VERSION=unknown

ARG         BUILD_TIME=unknown

ARG         COMMIT_HASH=unknown

COPY        . ./

RUN         CGO_ENABLED=0 \
            GOOS=linux \
            GOARCH=amd64 \
            go build \
              -trimpath \
              -ldflags '\
                -X "github.com/go-zoox/serve/constants.Version=${VERSION}" \
                -X "github.com/go-zoox/serve/constants.BuildTime=${BUILD_TIME}" \
                -X "github.com/go-zoox/serve/constants.CommitHash=${COMMIT_HASH}" \
                -w -s -buildid= \
              ' \
              -v -o serve

# Product
FROM alpine:latest

LABEL       MAINTAINER="Zero<tobewhatwewant@gmail.com>"

COPY        --from=builder /app/serve /bin

LABEL       org.opencontainers.image.source="https://github.com/go-zoox/serve"

ARG         VERSION=v1.0.0

ENV         VERSION=${VERSION}

COPY        entrypoint.sh /entrypoint.sh

CMD         /entrypoint.sh
