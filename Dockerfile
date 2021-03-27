FROM golang:1.16 AS builder

# Copy all relevant files
COPY ${PWD}/cmd /app/cmd
COPY ${PWD}/internal /app/internal
COPY ${PWD}/pkg /app/pkg
COPY ${PWD}/go* /app/

WORKDIR /app

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:3
LABEL MAINTAINER Jannis Lehmann <cludch@gmail.com>

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

# Add new user 'csgo'
RUN adduser -D csgo
USER csgo

COPY --from=builder /app /app

WORKDIR /app

VOLUME ["/config", "/demos"]
