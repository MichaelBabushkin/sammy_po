package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// TokenStatus tracks information about the x-mas token
type TokenStatus struct {
	Token     string    `json:"token"`
	UpdatedAt time.Time `json:"updatedAt"`
	Timestamp int64     `json:"timestamp"`
}

// FotmobHeaders holds the headers needed for Fotmob API requests
type FotmobHeaders struct {
	XMasToken   string            `json:"x-mas"`
	UserAgent   string            `json:"User-Agent"`
	Accept      string            `json:"Accept"`
	AllHeaders  map[string]string `json:"-"`
	LastUpdated time.Time         `json:"-"`
	ScrapedAt   int64             `json:"_scrapedAt"`
}

// Global settings for token expiration
var tokenExpirationTime = 24 * time.Hour // Default to 24 hours

// GetFotmobHeaders loads the headers from the currency_api_headers.json file
func GetFotmobHeaders() *FotmobHeaders {
	// Path to headers file (Corrected: relative to project root)
	headersFile := filepath.Join("responses", "currency_api_headers.json")

	// Check if file exists
	if _, err := os.Stat(headersFile); os.IsNotExist(err) {
		log.Println("Headers file not found, using defaults")
		return getDefaultHeaders()
	}

	// Read file
	data, err := ioutil.ReadFile(headersFile)
	if err != nil {
		log.Printf("Error reading headers file: %v", err)
		return getDefaultHeaders()
	}

	// Parse headers
	var headers map[string]interface{}
	if err := json.Unmarshal(data, &headers); err != nil {
		log.Printf("Error parsing headers JSON: %v", err)
		return getDefaultHeaders()
	}

	// Check if timestamp exists and is fresh
	if timestamp, ok := headers["_scrapedAt"]; ok {
		if ts, ok := timestamp.(float64); ok {
			scrapedAt := time.Unix(int64(ts), 0)
			if time.Since(scrapedAt) > tokenExpirationTime {
				log.Printf("Token is too old (scraped at %v), using default", scrapedAt)
				return getDefaultHeaders()
			}
			log.Printf("Using token scraped at %v", scrapedAt)
		}
	}

	// Create headers object
	result := &FotmobHeaders{
		AllHeaders:  make(map[string]string),
		LastUpdated: time.Now(),
	}

	// Extract scraped timestamp
	if scrapedAt, ok := headers["_scrapedAt"]; ok {
		if ts, ok := scrapedAt.(float64); ok {
			result.ScrapedAt = int64(ts)
		}
	}

	// Check if we have the x-mas token
	if xmas, ok := headers["x-mas"]; ok {
		if xmasStr, ok := xmas.(string); ok {
			result.XMasToken = xmasStr
			result.AllHeaders["x-mas"] = xmasStr
		}
	}

	// Process all other headers
	for k, v := range headers {
		if strVal, ok := v.(string); ok {
			result.AllHeaders[k] = strVal

			switch k {
			case "User-Agent":
				result.UserAgent = strVal
			case "Accept":
				result.Accept = strVal
			}
		}
	}

	// Fall back to defaults if critical headers are missing
	if result.XMasToken == "" {
		defaultHeaders := getDefaultHeaders()
		result.XMasToken = defaultHeaders.XMasToken
		result.AllHeaders["x-mas"] = defaultHeaders.XMasToken
	}

	if result.UserAgent == "" {
		result.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
		result.AllHeaders["User-Agent"] = result.UserAgent
	}

	return result
}

// SetTokenExpirationTime sets how long tokens are considered valid
func SetTokenExpirationTime(duration time.Duration) {
	if duration > 0 {
		tokenExpirationTime = duration
		log.Printf("Token expiration time set to %v", duration)
	}
}

// getDefaultHeaders returns a default set of headers
func getDefaultHeaders() *FotmobHeaders {
	return &FotmobHeaders{
		XMasToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiI2N2YxNGNkOTZlNTA0MzM2Y2M5MDA2ZWQiLCJJc0Fub255bW91cyI6IlRydWUiLCJFeHRlcm5hbElkIjoiYW5vbnltb3VzZTA4ZDg1YmItZmZhNi00M2JhLThkNmYtMDEwZWYyOWUzYjQwIiwiSXNUZW1wIjoiVHJ1ZSIsIk9yaWdpbmFsQW5vbnltb3VzVXNlcklkIjoiNjdmMTRjZDk2ZTUwNDMzNmNjOTAwNmVkIiwibmJmIjoxNzQzODY3MDk3LCJleHAiOjE3NDM5NTM0OTcsImlhdCI6MTc0Mzg2NzA5NywiaXNzIjoiV1NDIn0.Crw7k6iJMCAGLxfedRGFdwTjKR32uN3qy6UwUFIL5L4",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		Accept: "application/json, text/plain, */*",
		AllHeaders: map[string]string{
			"x-mas": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiI2N2YxNGNkOTZlNTA0MzM2Y2M5MDA2ZWQiLCJJc0Fub255bW91cyI6IlRydWUiLCJFeHRlcm5hbElkIjoiYW5vbnltb3VzZTA4ZDg1YmItZmZhNi00M2JhLThkNmYtMDEwZWYyOWUzYjQwIiwiSXNUZW1wIjoiVHJ1ZSIsIk9yaWdpbmFsQW5vbnltb3VzVXNlcklkIjoiNjdmMTRjZDk2ZTUwNDMzNmNjOTAwNmVkIiwibmJmIjoxNzQzODY3MDk3LCJleHAiOjE3NDM5NTM0OTcsImlhdCI6MTc0Mzg2NzA5NywiaXNzIjoiV1NDIn0.Crw7k6iJMCAGLxfedRGFdwTjKR32uN3qy6UwUFIL5L4",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept": "application/json, text/plain, */*",
		},
		LastUpdated: time.Now(),
		ScrapedAt: time.Now().Unix(),
	}
}
