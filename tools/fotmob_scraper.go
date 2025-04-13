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

// RunScraper is the function that can be called from other parts of the code
func RunScraper(silent bool) (string, error) {
	// Create a context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-sandbox", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create directory for saving data
	os.MkdirAll("responses", 0755)

	// Flag to track if we found the currency API request
	currencyApiFound := false
	var xmasToken string

	// Listen for network events to capture request headers
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			// Get the URL
			url := e.Request.URL
			
			// Extract all headers for this request
			headers := make(map[string]interface{})
			for key, value := range e.Request.Headers {
				headers[key] = value
			}
			
			// Special handling for currency API
			if strings.Contains(url, "/api/currency") {
				currencyApiFound = true
				if !silent {
					fmt.Printf("\n🔍 CURRENCY API REQUEST DETECTED: %s\n", url)
				}
				
				// Add timestamp to headers
				headers["_timestamp"] = time.Now().Format(time.RFC3339)
				headers["_scrapedAt"] = time.Now().Unix()
				
				// Save currency API headers to a separate file
				currencyHeaders, _ := json.MarshalIndent(headers, "", "  ")
				os.WriteFile("responses/currency_api_headers.json", currencyHeaders, 0644)
				
				// Save x-mas token if present
				if xmas, ok := headers["x-mas"]; ok {
					if str, ok := xmas.(string); ok {
						xmasToken = str
						os.WriteFile("responses/x-mas-token.txt", []byte(str), 0644)
					}
				}
			}
		}
	})

	if !silent {
		fmt.Println("Navigating to fotmob.com...")
	}

	// Run the browser automation
	err := chromedp.Run(ctx, 
		network.Enable(),
		chromedp.Navigate("https://www.fotmob.com"),
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if !currencyApiFound {
				return nil
			}
			return nil
		}),
		chromedp.Navigate("https://www.fotmob.com/leagues/127/overview/ligat-ha'al"),
		chromedp.Sleep(5*time.Second),
	)

	if err != nil {
		return "", fmt.Errorf("error running browser: %v", err)
	}

	if !silent && currencyApiFound {
		fmt.Println("\n✅ Currency API request was detected and saved!")
	}

	return xmasToken, nil
}

// main function is only called when run as a standalone program
func main() {
	// Check for silent flag
	silent := false
	for _, arg := range os.Args {
		if arg == "--silent" {
			silent = true
			break
		}
	}
	
	token, err := RunScraper(silent)
	if err != nil {
		log.Fatalf("Error running scraper: %v", err)
	}
	
	if token != "" {
		if !silent {
			fmt.Printf("Found token: %s...\n", token[:20])
		}
	} else {
		if !silent {
			fmt.Println("No token found")
		}
	}
}
