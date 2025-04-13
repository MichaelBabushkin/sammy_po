package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// FotmobHeaders holds the headers needed for Fotmob API requests
type FotmobHeaders struct {
	XMasToken   string            `json:"x-mas"`
	UserAgent   string            `json:"User-Agent"`
	Accept      string            `json:"Accept"`
	AllHeaders  map[string]string `json:"-"`
	LastUpdated time.Time         `json:"-"`
}

// GetFotmobHeaders loads the headers from the currency_api_headers.json file
func GetFotmobHeaders() *FotmobHeaders {
	// Path to headers file
	headersFile := filepath.Join("tools", "responses/currency_api_headers.json")
	
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
	
	// Create headers object
	result := &FotmobHeaders{
		AllHeaders:  make(map[string]string),
		LastUpdated: time.Now(),
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
	}
}
