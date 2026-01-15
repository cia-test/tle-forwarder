package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v3/udp"
)

func main() {
	satelliteID := flag.String("satellite_id", "", "NORAD catalog number")
	name := flag.String("name", "", "Satellite name")
	group := flag.String("group", "", "Satellite group")
	root := flag.Bool("root", false, "Test root endpoint")
	flag.Parse()

	co, err := udp.Dial("localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer co.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	path := "/tle"
	queries := []string{}
	
	if *root {
		path = "/"
	} else if *satelliteID != "" {
		queries = append(queries, fmt.Sprintf("satellite_id=%s", *satelliteID))
	} else if *name != "" {
		queries = append(queries, fmt.Sprintf("name=%s", *name))
	} else if *group != "" {
		queries = append(queries, fmt.Sprintf("group=%s", *group))
	}

	req, err := co.NewGetRequest(ctx, path)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	
	for _, q := range queries {
		req.AddQuery(q)
	}

	resp, err := co.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	data, err := resp.ReadBody()
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Response Code: %v\n", resp.Code())
	fmt.Printf("Response Body:\n%s\n", string(data))
}
