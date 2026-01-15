.PHONY: all build run-http run-coap test clean

all: build

build:
	go mod download
	go build -o tle-forwarder main.go
	go build -o tle-forwarder-coap coap_server.go

run-http:
	go run main.go

run-coap:
	go run coap_server.go

test:
	go test -v ./...

clean:
	rm -f tle-forwarder tle-forwarder-coap
	go clean
