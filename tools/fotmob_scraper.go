package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create directory for saving data
	os.MkdirAll("responses", 0755)

	// Flag to track if we found the currency API request
	currencyApiFound := false

	// Listen for network events to capture request headers
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			// Get the URL and method
			url := e.Request.URL
			method := e.Request.Method
			
			// Extract all headers for this request
			headers := make(map[string]interface{})
			for key, value := range e.Request.Headers {
				headers[key] = value
			}
			
			// Save interesting API requests for analysis
			if strings.Contains(url, "fotmob.com/api") {
				fmt.Printf("\n=== REQUEST: %s %s ===\n", method, url)
				for key, value := range headers {
					fmt.Printf("  %s: %v\n", key, value)
				}
				
				// Special handling for currency API
				if strings.Contains(url, "/api/currency") {
					currencyApiFound = true
					fmt.Printf("\n🔍 CURRENCY API REQUEST DETECTED: %s\n", url)
					
					// Save currency API headers to a separate file
					currencyHeaders, _ := json.MarshalIndent(headers, "", "  ")
					err := os.WriteFile("responses/currency_api_headers.json", currencyHeaders, 0644)
					if err == nil {
						fmt.Println("✅ Saved currency API headers to responses/currency_api_headers.json")
					}
					
					// Save x-mas token if present
					if xmas, ok := headers["x-mas"]; ok {
						fmt.Printf("🔑 X-MAS TOKEN: %v\n", xmas)
						if str, ok := xmas.(string); ok {
							os.WriteFile("responses/x-mas-token.txt", []byte(str), 0644)
						}
					}
				}
			}
		}
	})

	fmt.Println("Navigating to fotmob.com...")

	// Enable network monitoring and navigate to the site
	err := chromedp.Run(ctx, 
		network.Enable(),
		chromedp.Navigate("https://www.fotmob.com"),
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if !currencyApiFound {
				fmt.Println("\nCurrency API not detected, trying league page...")
				return nil
			}
			return nil
		}),
		chromedp.Navigate("https://www.fotmob.com/leagues/127/overview/ligat-ha'al"),
		chromedp.Sleep(5*time.Second),
	)

	if err != nil {
		log.Fatalf("Error running browser: %v", err)
	}

	// Summarize what we found
	if currencyApiFound {
		fmt.Println("\n✅ Currency API request was detected and saved!")
		fmt.Println("Check responses/currency_api_headers.json for headers")
		fmt.Println("Check responses/x-mas-token.txt for the token")
	} else {
		fmt.Println("\n❌ Currency API request was not detected")
	}

	fmt.Println("\nDone! Check the responses directory for captured data.")
}
