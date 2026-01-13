# Builder
FROM --platform=$BUILDPLATFORM whatwewant/builder-go:v1.24-1 as builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

ARG VERSION=unknown

ARG BUILD_TIME=unknown

ARG COMMIT_HASH=unknown

COPY . ./

ARG TARGETOS

ARG TARGETARCH

RUN CGO_ENABLED=0 \
  GOOS=${TARGETOS} \
  GOARCH=${TARGETARCH} \
  go build \
  -trimpath \
  -ldflags '\
      -X "github.com/go-zoox/serve/constants.Version=${VERSION}" \
      -X "github.com/go-zoox/serve/constants.BuildTime=${BUILD_TIME}" \
      -X "github.com/go-zoox/serve/constants.CommitHash=${COMMIT_HASH}" \
      -w -s -buildid= \
    ' \
  -v -o serve ./cmd/serve

# Product
FROM whatwewant/alpine:v3.17-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-zoox/serve"

ENV MODE=production

COPY --from=builder /build/serve /bin

ARG VERSION=latest

ENV VERSION=${VERSION}

COPY entrypoint.sh /entrypoint.sh

CMD /entrypoint.sh
