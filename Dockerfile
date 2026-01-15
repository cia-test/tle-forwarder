FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY src/ ./src/

# Build both servers
RUN go build -o tle-forwarder src/main.go
RUN go build -o tle-forwarder-coap src/coap_server.go

FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binaries from builder
COPY --from=builder /app/tle-forwarder .
COPY --from=builder /app/tle-forwarder-coap .

# Expose ports
EXPOSE 8000 5683/udp

# Default to HTTP server
CMD ["./tle-forwarder"]
