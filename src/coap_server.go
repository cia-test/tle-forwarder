package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	coapNet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/net/responsewriter"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/options"
	udpClient "github.com/plgd-dev/go-coap/v3/udp/client"
	"github.com/plgd-dev/go-coap/v3/udp/server"
)

const celestrakURL = "https://celestrak.org/NORAD/elements/gp.php"

func fetchTLEData(satelliteID, name, group string) (string, codes.Code, error) {
	params := url.Values{}
	params.Add("FORMAT", "TLE")

	if satelliteID != "" {
		params.Set("CATNR", satelliteID)
	} else if name != "" {
		params.Set("NAME", name)
	} else if group != "" {
		params.Set("GROUP", group)
	}

	reqURL := fmt.Sprintf("%s?%s", celestrakURL, params.Encode())
	log.Printf("Fetching TLE data from CelesTrak: %s", reqURL)

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", codes.ServiceUnavailable, fmt.Errorf("unable to reach CelesTrak: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", codes.InternalServerError, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", codes.BadGateway, fmt.Errorf("CelesTrak returned status %d", resp.StatusCode)
	}

	tleData := strings.TrimSpace(string(body))

	if tleData == "" || strings.Contains(tleData, "No GP data found") {
		return "", codes.NotFound, fmt.Errorf("no TLE data found for the specified parameters")
	}

	log.Printf("Successfully fetched TLE data (%d bytes)", len(tleData))
	return tleData, codes.Content, nil
}

func parseQuery(queries []string) (satelliteID, name, group string) {
	for _, q := range queries {
		parts := strings.SplitN(q, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "satellite_id":
			satelliteID = value
		case "name":
			name = value
		case "group":
			group = value
		}
	}
	return
}

func coapHandler(w *responsewriter.ResponseWriter[*udpClient.Conn], r *pool.Message) {
	path, err := r.Options().Path()
	if err != nil {
		log.Printf("Error getting path: %v", err)
		w.SetResponse(codes.BadRequest, message.TextPlain, strings.NewReader("Invalid request"))
		return
	}
	
	queries, _ := r.Options().Queries()
	satelliteID, name, group := parseQuery(queries)
	
	if path == "/tle" || path == "tle" {
		if satelliteID == "" && name == "" && group == "" {
			w.SetResponse(codes.BadRequest, message.TextPlain, strings.NewReader("Please provide satellite_id, name, or group parameter"))
			return
		}

		tleData, code, err := fetchTLEData(satelliteID, name, group)
		if err != nil {
			log.Printf("Error fetching TLE: %v", err)
			w.SetResponse(code, message.TextPlain, strings.NewReader(err.Error()))
			return
		}

		w.SetResponse(codes.Content, message.TextPlain, strings.NewReader(tleData))
	} else {
		info := "TLE Forwarder CoAP Service\n" +
			"Usage: coap://localhost:5683/tle?satellite_id=25544\n" +
			"Parameters: satellite_id, name, or group"
		w.SetResponse(codes.Content, message.TextPlain, strings.NewReader(info))
	}
}

func main() {
	log.Println("Starting CoAP server on :5683")
	log.Println("Try: coap://localhost:5683/tle?satellite_id=25544")

	conn, err := coapNet.NewListenUDP("udp4", ":5683")
	if err != nil {
		log.Fatalf("Failed to create UDP listener: %v", err)
	}
	defer conn.Close()

	srv := server.New(options.WithHandlerFunc(coapHandler))
	
	if err := srv.Serve(conn); err != nil {
		log.Fatalf("Failed to start CoAP server: %v", err)
	}
}
