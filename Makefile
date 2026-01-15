.PHONY: all build run-http run-coap test clean docker-build docker-up docker-down

all: build

build:
	go mod download
	go build -o tle-forwarder src/main.go
	go build -o tle-forwarder-coap src/coap_server.go

run-http:
	go run src/main.go

run-coap:
	go run src/coap_server.go

test:
	go test -v ./...

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -f tle-forwarder tle-forwarder-coap
	go clean
