FROM golang:1.16 AS builder

WORKDIR /app

# Copy all relevant files
COPY ${PWD}/cmd /app/cmd
COPY ${PWD}/internal /app/internal
COPY ${PWD}/pkg /app/pkg
COPY ${PWD}/go* /app/

RUN go get -d -v ./...
RUN go build  -o bin/ ./...

FROM alpine:3
LABEL MAINTAINER Jannis Lehmann <cludch@gmail.com>

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

# Add new user 'csgo'
RUN adduser -S -D -H -h /app csgo
USER csgo

COPY --from=builder /app/bin /app

RUN chmod +x *

WORKDIR /app

VOLUME ["/config", "/demos"]
