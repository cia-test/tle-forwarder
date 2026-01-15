# TLE Forwarder

A web application that acts as a TLE (Two-Line Element) file forwarder, fetching orbital data from CelesTrak. Supports both HTTP and CoAP protocols.

## Features

- Fetch TLE data by NORAD satellite catalog number
- Search TLE data by satellite name
- Retrieve TLE data for satellite groups
- RESTful HTTP API (with Gin framework)
- Native CoAP protocol support
- Written in Go for high performance and easy deployment

## Installation

```bash
# Install Go dependencies
go mod download

# Build both servers
make build

# Or build individually
go build -o tle-forwarder main.go
go build -o tle-forwarder-coap coap_server.go
```

## Usage

### Option 1: HTTP Server

Start the HTTP server:

```bash
# Run directly
go run main.go

# Or run the compiled binary
./tle-forwarder

# Or use make
make run-http
```

The HTTP server runs on `http://localhost:8000`

### Option 2: CoAP Server

Start the CoAP server:

```bash
# Run directly
go run coap_server.go

# Or run the compiled binary
./tle-forwarder-coap

# Or use make
make run-coap
```

The CoAP server runs on `coap://localhost:5683`

### Testing CoAP

Use a CoAP client tool like `coap-client`:

```bash
coap-client -m get "coap://localhost:5683/tle?satellite_id=25544"
```

## API Endpoints

### HTTP Endpoints

#### GET /tle

Fetch TLE data from CelesTrak.

**Query Parameters:**
- `satellite_id` (int): NORAD catalog number (e.g., 25544 for ISS)
- `name` (string): Satellite name search (e.g., "ISS", "STARLINK")
- `group` (string): Pre-defined satellite group (e.g., "stations", "visual", "active")

**Examples:**

```bash
# By satellite ID (ISS)
curl "http://localhost:8000/tle?satellite_id=25544"

# By name
curl "http://localhost:8000/tle?name=ISS"

# By group
curl "http://localhost:8000/tle?group=stations"
```

**Response:** Returns TLE data in plain text format

#### GET /

Returns service information and usage examples in JSON format.

#### GET /health

Health check endpoint. Returns `{"status": "healthy"}`

### CoAP Endpoints

#### GET coap://localhost:5683/tle

Fetch TLE data via CoAP protocol.

**Query Parameters:** Same as HTTP endpoint (satellite_id, name, group)

**Examples:**

```bash
# Using coap-client tool
coap-client -m get "coap://localhost:5683/tle?satellite_id=25544"
coap-client -m get "coap://localhost:5683/tle?name=ISS"
coap-client -m get "coap://localhost:5683/tle?group=stations"
```

#### GET coap://localhost:5683/

Returns service information.

## TLE Data Source

This application fetches data from [CelesTrak](https://celestrak.org), a leading source for orbital element sets.

## Common Satellite IDs

- ISS (International Space Station): 25544
- Hubble Space Telescope: 20580
- Starlink satellites: Various IDs (search by name "STARLINK")

## Common Groups

- `stations` - Space stations
- `visual` - Bright satellites visible to the naked eye
- `active` - Active satellites
- `analyst` - Analyst satellites
- `2024-launch` - Satellites launched in 2024

## Response Format

The response is in standard TLE format:

```
SATELLITE NAME
1 NNNNNU NNNNNAAA NNNNN.NNNNNNNN +.NNNNNNNN +NNNNN-N +NNNNN-N N NNNNN
2 NNNNN NNN.NNNN NNN.NNNN NNNNNNN NNN.NNNN NNN.NNNN NN.NNNNNNNNNNNNNN
```

Each satellite has three lines:
1. Satellite name
2. Line 1 of TLE data
3. Line 2 of TLE data
