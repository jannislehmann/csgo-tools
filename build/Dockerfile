FROM golang:1.15-alpine

LABEL author="Jannis Lehmman"

WORKDIR /go/src/app
COPY . .

# GCC is required to build go-sqlite3
RUN apk add --no-cache build-base=0.5-r2

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME ["/config", "/demos"]

CMD ["apiclient"]
