package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const CelesTrakGPURL = "https://celestrak.org/NORAD/elements/gp.php"

type TLEResponse struct {
	Service     string                 `json:"service"`
	Description string                 `json:"description"`
	Endpoints   map[string]interface{} `json:"endpoints"`
	DataSource  string                 `json:"data_source"`
}

func fetchTLE(satelliteID, name, group string) (string, int, error) {
	params := url.Values{}
	params.Add("FORMAT", "TLE")

	if satelliteID != "" {
		params.Set("CATNR", satelliteID)
	} else if name != "" {
		params.Set("NAME", name)
	} else if group != "" {
		params.Set("GROUP", group)
	}

	reqURL := fmt.Sprintf("%s?%s", CelesTrakGPURL, params.Encode())
	log.Printf("Fetching TLE data from CelesTrak: %s", reqURL)

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", http.StatusServiceUnavailable, fmt.Errorf("unable to reach CelesTrak: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, fmt.Errorf("CelesTrak returned status %d", resp.StatusCode)
	}

	tleData := strings.TrimSpace(string(body))

	if tleData == "" || strings.Contains(tleData, "No GP data found") {
		return "", http.StatusNotFound, fmt.Errorf("no TLE data found for the specified parameters")
	}

	log.Printf("Successfully fetched TLE data (%d bytes)", len(tleData))
	return tleData, http.StatusOK, nil
}

func getTLE(c *gin.Context) {
	satelliteID := c.Query("satellite_id")
	name := c.Query("name")
	group := c.Query("group")

	if satelliteID == "" && name == "" && group == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide at least one parameter: satellite_id, name, or group",
		})
		return
	}

	tleData, statusCode, err := fetchTLE(satelliteID, name, group)
	if err != nil {
		c.String(statusCode, err.Error())
		return
	}

	c.String(http.StatusOK, tleData)
}

func root(c *gin.Context) {
	response := TLEResponse{
		Service:     "TLE Forwarder",
		Description: "Fetch TLE (Two-Line Element) data from CelesTrak",
		Endpoints: map[string]interface{}{
			"/tle": map[string]interface{}{
				"method": "GET",
				"parameters": map[string]string{
					"satellite_id": "NORAD catalog number (e.g., 25544 for ISS)",
					"name":         "Satellite name search (e.g., ISS, STARLINK)",
					"group":        "Satellite group (e.g., stations, visual, active)",
				},
				"examples": []string{
					"/tle?satellite_id=25544",
					"/tle?name=ISS",
					"/tle?group=stations",
				},
			},
		},
		DataSource: "CelesTrak (https://celestrak.org)",
	}
	c.JSON(http.StatusOK, response)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", root)
	r.GET("/tle", getTLE)
	r.GET("/health", healthCheck)

	log.Println("Starting HTTP server on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
