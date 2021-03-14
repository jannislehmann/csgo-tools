FROM golang:1.16-alpine

LABEL author="Jannis Lehmman"

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME ["/config", "/demos"]
