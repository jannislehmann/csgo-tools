FROM golang:1.16 AS builder

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -v -tags netgo -ldflags '-w -extldflags "-static"' -o bin/ ./...

FROM alpine:3
LABEL MAINTAINER Jannis Lehmann <cludch@gmail.com>

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update --no-cache add ca-certificates libc6-compat  && \
    rm -rf /var/cache/apk/*

# Add new user 'csgo'
RUN adduser -S -D -H -h /app csgo
USER csgo

WORKDIR /app

COPY --from=builder /app/bin /app

VOLUME ["/config", "/demos"]