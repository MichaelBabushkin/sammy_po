package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// FallbackWithSimpleHTTP gets the token using a simple HTTP request (kept as fallback)
func FallbackWithSimpleHTTP() (string, error) {
	fmt.Println("Using simple HTTP request to get token...")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://www.fotmob.com", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Add realistic browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	content := string(body)
	var token string

	// Try to find the token in the HTML using regex (look for long tokens)
	patterns := []string{
		`"x-mas":"([^"]{100,})"`, // x-mas token in double quotes, at least 100 chars
		`'x-mas':'([^']{100,})'`, // x-mas token in single quotes, at least 100 chars
		`"x-mas"\s*:\s*"([^"]{100,})"`, // x-mas token with spaces, double quotes
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			token = matches[1]
			fmt.Println("Found token using x-mas pattern (HTTP).")
			break
		}
	}

	// If not found, try looking for a generic long JWT pattern
	if token == "" {
		fmt.Println("x-mas pattern failed (HTTP), trying generic JWT pattern...")
		jwtPattern := `eyJ[a-zA-Z0-9_-]{50,}\.[a-zA-Z0-9_-]{100,}\.[a-zA-Z0-9_-]+` // Look for long JWTs
		re := regexp.MustCompile(jwtPattern)
		matches := re.FindAllString(content, -1)

		// Find the longest JWT token (often the one needed)
		for _, match := range matches {
			if len(match) > len(token) {
				token = match
			}
		}
		if token != "" {
			fmt.Println("Found potential token using generic JWT pattern (HTTP).")
		}
	}

	if token == "" {
		return "", fmt.Errorf("no valid token found in HTML response (HTTP)")
	}

	// Save token and headers
	saveTokenAndHeaders(token)

	return token, nil
}

// RunTokenScraper runs the browser automation to get the Fotmob token
func RunTokenScraper(silent bool) (string, error) {
	if !silent {
		fmt.Println("Starting browser automation...")
	}
	// Create a context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Run headless
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-sandbox", true), // Often needed in containerized environments
		chromedp.Flag("disable-dev-shm-usage", true), // Overcome resource limits
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second) // Increased timeout
	defer cancel()

	var xmasToken string
	var headersJSON []byte // To store the full headers JSON

	// Listen for network events to capture request headers
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			url := e.Request.URL
			// Look for requests to the Fotmob API
			if strings.Contains(url, "/api/") {
				if !silent {
					// fmt.Printf("Detected API request: %s\n", url) // Can be noisy
				}
				// Check if this request has the x-mas token
				if xmas, ok := e.Request.Headers["x-mas"]; ok {
					if str, ok := xmas.(string); ok && len(str) > 100 { // Check length
						xmasToken = str
						if !silent {
							fmt.Printf("Found x-mas token via network listener: %s...\n", truncateToken(xmasToken))
						}

						// Capture all headers from this specific request
						headersMap := make(map[string]interface{})
						for k, v := range e.Request.Headers {
							headersMap[k] = v
						}
						headersMap["_timestamp"] = time.Now().Format(time.RFC3339)
						headersMap["_scrapedAt"] = time.Now().Unix()

						// Marshal the captured headers
						headersJSON, _ = json.MarshalIndent(headersMap, "", "  ")

						// We found the token, potentially cancel further actions
						// cancel() // Uncomment if you want to stop immediately after finding token
					}
				}
			}
		}
	})

	// Run the browser automation steps
	err := chromedp.Run(ctx,
		network.Enable(), // Enable network domain
		chromedp.Navigate("https://www.fotmob.com"),
		chromedp.Sleep(5*time.Second), // Wait for initial load and potential redirects/scripts
		// Optional: Add more actions like clicking or scrolling if needed
		// chromedp.Click("#some-element", chromedp.NodeVisible),
		// chromedp.ScrollIntoView("footer"),
		chromedp.Sleep(5*time.Second), // Wait a bit longer for background requests
	)

	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		log.Printf("Browser automation error: %v", err)
		// Don't return error yet, try fallback
	} else if err == context.DeadlineExceeded {
		log.Println("Browser automation timed out.")
	}

	// If token was found via listener, save it
	if xmasToken != "" && len(headersJSON) > 0 {
		if !silent {
			fmt.Println("Saving token and headers found via browser automation.")
		}
		saveTokenAndHeadersFromJSON(xmasToken, headersJSON)
		return xmasToken, nil
	}

	// If browser automation failed or didn't find token, try fallback
	if !silent {
		fmt.Println("Browser automation didn't find token, trying fallback HTTP method...")
	}
	token, fallbackErr := FallbackWithSimpleHTTP()
	if fallbackErr != nil {
		// Combine errors if both methods failed
		if err != nil {
			return "", fmt.Errorf("browser automation failed (%v) and fallback HTTP failed (%v)", err, fallbackErr)
		}
		return "", fmt.Errorf("browser automation succeeded but found no token, and fallback HTTP failed (%v)", fallbackErr)
	}

	return token, nil // Return token from fallback
}

// saveTokenAndHeaders saves the token and constructs headers map
func saveTokenAndHeaders(token string) {
	headers := map[string]interface{}{
		"x-mas":      token,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Accept":     "application/json, text/plain, */*",
		"_timestamp": time.Now().Format(time.RFC3339),
		"_scrapedAt": time.Now().Unix(),
	}
	headersJSON, _ := json.MarshalIndent(headers, "", "  ")
	saveTokenAndHeadersFromJSON(token, headersJSON)
}

// saveTokenAndHeadersFromJSON saves token and headers from marshaled JSON
func saveTokenAndHeadersFromJSON(token string, headersJSON []byte) {
	responsesDir := "responses"
	headersFilePath := filepath.Join(responsesDir, "currency_api_headers.json")
	tokenFilePath := filepath.Join(responsesDir, "x-mas-token.txt")

	// Ensure the responses directory exists
	if err := os.MkdirAll(responsesDir, 0755); err != nil {
		log.Printf("Warning: Failed to create responses directory: %v", err)
		return // Can't save if directory fails
	}

	// Save headers file
	if err := ioutil.WriteFile(headersFilePath, headersJSON, 0644); err != nil {
		log.Printf("Warning: Failed to write headers file '%s': %v", headersFilePath, err)
	}

	// Save token file
	if err := ioutil.WriteFile(tokenFilePath, []byte(token), 0644); err != nil {
		log.Printf("Warning: Failed to write token file '%s': %v", tokenFilePath, err)
	}
}

// Helper function to truncate token for logging
func truncateToken(token string) string {
	if len(token) > 30 {
		return token[:30] + "..."
	}
	return token
}

// No main function needed when used as a package
