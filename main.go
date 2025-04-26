package main

import (
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

	"github.com/MichaelBabushkin/sammy_po/api"
	"github.com/MichaelBabushkin/sammy_po/pkg/scraper" // Import the new scraper package
	"github.com/joho/godotenv"
)

// Simplified types section
type Match struct {
	ID             int     `json:"id"`
	HomeTeam       string  `json:"homeTeam"`
	HomeTeamLogo   string  `json:"homeTeamLogo"`
	AwayTeam       string  `json:"awayTeam"`
	AwayTeamLogo   string  `json:"awayTeamLogo"`
	HomeScore      *int    `json:"homeScore"`
	AwayScore      *int    `json:"awayScore"`
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	Competition    string  `json:"competition"`
	CompetitionLogo string `json:"competitionLogo"`
	Status         string  `json:"status"`
	Round          string  `json:"round"`
	Venue          string  `json:"venue"`
}

// SammyOferInfo represents data about Sammy Ofer Stadium
type SammyOferInfo struct {
	Name        string   `json:"name"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Capacity    int      `json:"capacity"`
	ImageURL    string   `json:"imageUrl"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Teams       []string `json:"teams"`
}

type FotmobClient struct {
	client *http.Client
}

func NewFotmobClient() *FotmobClient {
	return &FotmobClient{
		client: &http.Client{},
	}
}

func (c *FotmobClient) makeRequest(url string) ([]byte, error) {
	log.Printf("Making Fotmob request to: %s", url)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Get the headers from our API package
	headers := api.GetFotmobHeaders()
	
	// Log the token info including timestamp
	log.Printf("Using x-mas token: %s... (scraped at: %s)", 
		truncateToken(headers.XMasToken), 
		time.Unix(headers.ScrapedAt, 0).Format(time.RFC3339))
	
	// Add the headers to the request
	req.Header.Add("x-mas", headers.XMasToken)
	req.Header.Add("User-Agent", headers.UserAgent)
	if headers.Accept != "" {
		req.Header.Add("Accept", headers.Accept)
	}
	
	// Add any other important headers we might have discovered
	for k, v := range headers.AllHeaders {
		if k != "x-mas" && k != "User-Agent" && k != "Accept" {
			if !strings.HasPrefix(k, ":") && k != "Connection" && k != "Host" {
				req.Header.Add(k, v)
			}
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return ioutil.ReadAll(resp.Body)
}

// Helper function to truncate token for logging
func truncateToken(token string) string {
	if len(token) > 30 {
		return token[:30] + "..."
	}
	return token
}

func (c *FotmobClient) FetchIsraeliLeagueData() (map[string]interface{}, error) {
	url := "https://www.fotmob.com/api/leagues?id=127&ccode3=ISR"
	body, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *FotmobClient) FetchIsraeliLeagueMatches() (interface{}, error) {
	fullData, err := c.FetchIsraeliLeagueData()
	if err != nil {
		return nil, err
	}
	
	matches, ok := fullData["matches"]
	if !ok {
		return nil, fmt.Errorf("matches field not found in API response")
	}
	
	matchesMap, ok := matches.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("matches is not a proper object")
	}
	
	allMatches, ok := matchesMap["allMatches"]
	if !ok {
		return nil, fmt.Errorf("allMatches field not found in matches object")
	}
	
	return allMatches, nil
}

// FilterMatches filters matches based on provided criteria
func FilterMatches(matches []interface{}, teamName string, isHome bool) []interface{} {
	filtered := []interface{}{}
	
	for _, match := range matches {
		matchMap, ok := match.(map[string]interface{})
		if !ok {
			continue
		}
		
		var teamToCheck map[string]interface{}
		var teamKey string
		
		if isHome {
			teamKey = "home"
		} else {
			teamKey = "away"
		}
		
		teamData, ok := matchMap[teamKey]
		if !ok {
			continue
		}
		
		teamToCheck, ok = teamData.(map[string]interface{})
		if !ok {
			continue
		}
		
		name, ok := teamToCheck["name"].(string)
		if !ok {
			continue
		}
		
		if strings.Contains(name, teamName) {
			filtered = append(filtered, match)
		}
	}
	
	return filtered
}

// Get Sammy Ofer Stadium info
func GetSammyOferInfo() SammyOferInfo {
	return SammyOferInfo{
		Name:        "Sammy Ofer Stadium",
		City:        "Haifa",
		Country:     "Israel",
		Capacity:    30858,
		ImageURL:    "https://stadiumdb.com/pictures/stadiums/isr/sammy_ofer_stadium/sammy_ofer_stadium21.jpg",
		Description: "Sammy Ofer Stadium is a football stadium in Haifa, Israel. It serves as a venue for home matches of both Maccabi Haifa and Hapoel Haifa football clubs. The stadium is named after shipping magnate and philanthropist Sammy Ofer, who donated $20 million to help build the stadium.",
		Address:     "32 Haim Weizmann St., Haifa, Israel",
		Teams:       []string{"Maccabi Haifa", "Hapoel Haifa"},
	}
}

// FilterHaifaHomeMatches filters matches for Haifa teams playing at home
func FilterHaifaHomeMatches(matches []interface{}) []interface{} {
	haifaMatches := []interface{}{}
	now := time.Now().UTC()
	
	// First: collect all Haifa home matches
	for _, match := range matches {
		matchMap, ok := match.(map[string]interface{})
		if !ok {
			continue
		}
		
		homeTeamData, ok := matchMap["home"]
		if !ok {
			continue
		}
		
		homeTeam, ok := homeTeamData.(map[string]interface{})
		if !ok {
			continue
		}
		
		homeTeamName, ok := homeTeam["name"].(string)
		if !ok {
			continue
		}
		
		if strings.Contains(homeTeamName, "Maccabi Haifa") || 
		   strings.Contains(homeTeamName, "Hapoel Haifa") {
			haifaMatches = append(haifaMatches, match)
		}
	}
	
	// Second: separate upcoming and past matches
	upcomingMatches := []interface{}{}
	pastMatches := []interface{}{}
	
	for _, match := range haifaMatches {
		matchMap := match.(map[string]interface{})
		
		statusData, ok := matchMap["status"]
		if !ok {
			continue
		}
		
		status, ok := statusData.(map[string]interface{})
		if !ok {
			continue
		}
		
		utcTimeStr, ok := status["utcTime"].(string)
		if !ok {
			if timeTS, ok := matchMap["timeTS"]; ok {
				if tsFloat, ok := timeTS.(float64); ok {
					matchTime := time.Unix(int64(tsFloat), 0)
					if matchTime.After(now) {
						upcomingMatches = append(upcomingMatches, match)
					} else {
						pastMatches = append(pastMatches, match)
					}
				}
			}
			continue
		}
		
		var matchTime time.Time
		var err error
		
		if strings.Contains(utcTimeStr, ".") {
			matchTime, err = time.Parse("2006-01-02T15:04:05.000Z", utcTimeStr)
		} else {
			matchTime, err = time.Parse("2006-01-02T15:04:05Z", utcTimeStr)
		}
		
		if err != nil {
			continue
		}
		
		if matchTime.After(now) {
			upcomingMatches = append(upcomingMatches, match)
		} else {
			pastMatches = append(pastMatches, match)
		}
	}
	
	// Return upcoming matches first, then past matches
	return append(upcomingMatches, pastMatches...)
}

// Updated function for refreshing the token (uses scraper package directly)
func refreshToken() {
	log.Println("Refreshing token...")

	// Ensure responses directory exists
	os.MkdirAll("responses", 0755)

	// --- Attempt 1: Run scraper package (chromedp version) ---
	log.Println("Attempting token refresh via scraper package (chromedp)...")
	token, err := scraper.RunTokenScraper(true) // Run silently

	if err != nil {
		log.Printf("Scraper package error: %v", err)
		// Don't immediately fail, try HTTP next
	} else if token != "" {
		// Verify token file exists (scraper should have saved it)
		tokenFile := filepath.Join("responses", "x-mas-token.txt")
		if _, statErr := os.Stat(tokenFile); os.IsNotExist(statErr) {
			log.Println("Warning: Scraper ran successfully but token file was not created.")
			// Proceed to HTTP method
		} else {
			log.Println("Token refreshed successfully via scraper package.")
			return // Success!
		}
	}

	// --- Attempt 2: Direct HTTP Method (Fallback) ---
	log.Println("Scraper package failed or didn't produce token. Attempting token refresh via direct HTTP...")
	_, httpErr := getTokenDirectHTTP() // Assumes getTokenDirectHTTP still exists
	if httpErr != nil {
		log.Printf("Direct HTTP token refresh also failed: %v", httpErr)
		log.Println("All automatic token refresh methods failed.")
	} else {
		log.Println("Token refreshed successfully via direct HTTP.")
	}
}

// Updated function for manually refreshing the token via endpoint (uses scraper package)
func manualRefreshToken() (string, error) {
	log.Println("Manually refreshing token via scraper package...")

	// Ensure responses directory exists
	os.MkdirAll("responses", 0755)

	// Run the scraper package (not silent)
	_, err := scraper.RunTokenScraper(false) // Discard the token value since it's not used

	if err != nil {
		log.Printf("Error running scraper package for manual refresh: %v", err)
		// Try direct HTTP as fallback
		log.Println("Scraper package failed, trying direct HTTP for manual refresh...")
		return getTokenDirectHTTP()
	}

	// Check if token file was created (scraper should save it)
	tokenFile := filepath.Join("responses", "x-mas-token.txt")
	if _, statErr := os.Stat(tokenFile); os.IsNotExist(statErr) {
		log.Println("Manual scraper ran but token file not found, trying direct HTTP...")
		return getTokenDirectHTTP() // Fallback if file not created
	}

	// Read the token (optional, scraper already saved it)
	tokenData, readErr := ioutil.ReadFile(tokenFile)
	if readErr != nil {
		return "", fmt.Errorf("failed to read token file after manual refresh: %v", readErr)
	}

	tokenFromFile := string(tokenData)
	if len(tokenFromFile) < 20 {
		return "", fmt.Errorf("invalid token in file after manual refresh")
	}

	log.Printf("Successfully refreshed token manually: %s...", truncateToken(tokenFromFile))
	return tokenFromFile, nil
}

// Add a function to check if token is fresh
func isTokenFresh() bool {
	// Check if the token file exists and is fresh
	tokenFile := filepath.Join("responses", "x-mas-token.txt")
	info, err := os.Stat(tokenFile)
	if err != nil {
		return false // File doesn't exist or can't be accessed
	}
	
	// Check if the file is less than 2 minutes old
	return time.Since(info.ModTime()) < 2*time.Minute
}

// Implement direct HTTP token fetch for the endpoint
func getTokenDirectHTTP() (string, error) {
	log.Println("Getting token via direct HTTP request...")

	client := &http.Client{
		Timeout: 15 * time.Second, // Increased timeout slightly
	}
	req, err := http.NewRequest("GET", "https://www.fotmob.com", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Add more realistic headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

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
	// log.Printf("HTML Content Length: %d", len(content)) // Optional: Log content length for debugging

	var token string

	// --- Strategy 1: Look for specific script tag content ---
	scriptRegex := regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	match := scriptRegex.FindStringSubmatch(content)
	if len(match) > 1 {
		jsonData := match[1]
		// Try parsing the JSON to find the token - this is more robust if structure is known
		var nextData map[string]interface{}
		if json.Unmarshal([]byte(jsonData), &nextData) == nil {
			// Navigate through the JSON structure if the path to the token is known
			// Example: if props, ok := nextData["props"].(map[string]interface{}); ok { ... }
			// For now, just search within the raw JSON string
			jsonStr := string(jsonData)
			patterns := []string{`"x-mas":"(.*?)"`} // More specific pattern within JSON
			for _, pattern := range patterns {
				r := regexp.MustCompile(pattern)
				tokenMatch := r.FindStringSubmatch(jsonStr)
				if len(tokenMatch) > 1 && len(tokenMatch[1]) > 20 {
					token = tokenMatch[1]
					log.Println("Found token via __NEXT_DATA__ JSON search.")
					break
				}
			}
		}
	}

	// --- Strategy 2: Regex for JWT patterns in the whole HTML ---
	if token == "" {
		log.Println("Token not found in __NEXT_DATA__, searching entire HTML for JWT patterns...")
		jwtPattern := `eyJ[a-zA-Z0-9_-]{10,}\.eyJ[a-zA-Z0-9_-]{50,}\.[a-zA-Z0-9_-]+` // More specific JWT pattern
		jwtRegex := regexp.MustCompile(jwtPattern)
		matches := jwtRegex.FindAllString(content, -1)

		if len(matches) > 0 {
			// Often the longest JWT is the one needed
			longestToken := ""
			for _, match := range matches {
				if len(match) > len(longestToken) {
					longestToken = match
				}
			}
			if len(longestToken) > 100 { // Basic validation
				token = longestToken
				log.Println("Found potential token via general JWT regex search.")
			}
		}
	}

	// --- Strategy 3: Original simple string search (fallback) ---
	if token == "" {
		log.Println("JWT regex failed, trying simple string search...")
		patterns := []string{`"x-mas":"`, `'x-mas':'`, `"x-mas"\s*:\s*"`}
		for _, pattern := range patterns {
			idx := strings.Index(content, pattern)
			if idx >= 0 {
				tokenStart := idx + len(pattern)
				quoteChar := pattern[len(pattern)-1:] // Get the quote character
				tokenEnd := strings.Index(content[tokenStart:], quoteChar)

				if tokenEnd > 0 {
					potentialToken := content[tokenStart : tokenStart+tokenEnd]
					if len(potentialToken) > 20 { // Basic validation
						token = potentialToken
						log.Println("Found token via simple string search.")
						break
					}
				}
			}
		}
	}

	if token == "" {
		// Optional: Save HTML for inspection if token not found
		// ioutil.WriteFile("fotmob_debug.html", []byte(content), 0644)
		return "", fmt.Errorf("token not found in HTML response using multiple strategies")
	}

	log.Printf("Found token: %s...", truncateToken(token))

	// Create headers file
	headers := map[string]interface{}{
		"x-mas":      token,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Accept":     "application/json, text/plain, */*",
		"_timestamp": time.Now().Format(time.RFC3339),
		"_scrapedAt": time.Now().Unix(),
	}

	// Save headers file (Corrected path)
	headersFilePath := filepath.Join("responses", "currency_api_headers.json")
	tokenFilePath := filepath.Join("responses", "x-mas-token.txt")

	os.MkdirAll("responses", 0755)
	headersJSON, _ := json.MarshalIndent(headers, "", "  ")
	err = ioutil.WriteFile(headersFilePath, headersJSON, 0644)
	if err != nil {
		log.Printf("Warning: Failed to write headers file %s: %v", headersFilePath, err)
		// Don't return error, just log it
	}

	// Also save token to separate file
	err = ioutil.WriteFile(tokenFilePath, []byte(token), 0644)
	if err != nil {
		log.Printf("Warning: Failed to write token file %s: %v", tokenFilePath, err)
		// Don't return error, just log it
	}

	return token, nil
}

func init() {
	// Load .env file
	godotenv.Load()
	
	// Ensure the responses directory exists
	os.MkdirAll("responses", 0755)
}

func main() {
	// Add specialized endpoint for Sammy Ofer matches (Haifa home games)
	http.HandleFunc("/api/fotmob/sammyofer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request for Sammy Ofer matches")
		
			// Check if the token needs refreshing
		if !isTokenFresh() {
			log.Println("Token is stale, refreshing...")
			refreshToken()
		}
		
		// Record the start time for performance tracking
		startTime := time.Now()
		
		fotmobClient := NewFotmobClient()
		matchesData, err := fotmobClient.FetchIsraeliLeagueMatches()

		if err != nil {
			log.Printf("Error fetching Fotmob matches: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		log.Printf("Successfully fetched matches data in %v", time.Since(startTime))
		
		matchesArr, ok := matchesData.([]interface{})
		if !ok {
			http.Error(w, "Invalid matches data format", http.StatusInternalServerError)
			return
		}
		
		log.Printf("Processing %d total matches from API", len(matchesArr))
		filteredMatches := FilterHaifaHomeMatches(matchesArr)
		
		// Get all upcoming matches
		now := time.Now().UTC()
		log.Printf("Current time (UTC): %s", now.Format(time.RFC3339))
		upcomingMatches := []interface{}{}
		
		for _, match := range filteredMatches {
			matchMap, ok := match.(map[string]interface{})
			if !ok {
				continue
			}
			
			statusData, ok := matchMap["status"]
			if !ok {
				continue
			}
			
			status, ok := statusData.(map[string]interface{})
			if !ok {
				continue
			}
			
			utcTimeStr, ok := status["utcTime"].(string)
			if !ok {
				continue
			}
			
			var matchTime time.Time
			var err error
			
			if strings.Contains(utcTimeStr, ".") {
				matchTime, err = time.Parse("2006-01-02T15:04:05.000Z", utcTimeStr)
			} else {
				matchTime, err = time.Parse("2006-01-02T15:04:05Z", utcTimeStr)
			}
			
			if err != nil {
				continue
			}
			
			// Include only future matches
			if matchTime.After(now) {
				log.Printf("Found upcoming match: %s vs %s on %s", 
					getTeamName(matchMap, "home"),
					getTeamName(matchMap, "away"),
					matchTime.Format(time.RFC1123))
				upcomingMatches = append(upcomingMatches, match)
			}
		}
		
		log.Printf("Found %d upcoming Sammy Ofer matches (request took %v)", 
			len(upcomingMatches), time.Since(startTime))
		
		w.Header().Set("Cache-Control", "max-age=3600") // Cache for 1 hour
		json.NewEncoder(w).Encode(upcomingMatches)
	})

	// Add endpoint for Sammy Ofer Stadium info
	http.HandleFunc("/api/stadium/sammyofer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		stadiumInfo := GetSammyOferInfo()
		w.Header().Set("Cache-Control", "max-age=86400") // Cache for 24 hours
		json.NewEncoder(w).Encode(stadiumInfo)
	})

	// Add a new endpoint for manually refreshing the token
	http.HandleFunc("/api/refresh-token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request to manually refresh token")
		
		// First try to run the scraper as a subprocess
		token, err := manualRefreshToken()
		
		// If that fails, try direct HTTP method
		if err != nil {
			log.Printf("Scraper method failed: %v. Trying direct HTTP method...", err)
			token, err = getTokenDirectHTTP()
			
			if err != nil {
				log.Printf("All token refresh methods failed: %v", err)
				http.Error(w, fmt.Sprintf("Failed to refresh token: %v", err), http.StatusInternalServerError)
				return
			}
		}
		
		// Return success response
		response := map[string]interface{}{
			"success": true,
			"message": "Token refreshed successfully",
			"tokenPreview": truncateToken(token),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	// Serve static files from the frontend/build directory
	fs := http.FileServer(http.Dir("frontend/build"))
	http.Handle("/", fs)

	// Determine the port to listen on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting server on :%s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Helper function to get team name from match object
func getTeamName(matchMap map[string]interface{}, side string) string {
	if teamData, ok := matchMap[side]; ok {
		if team, ok := teamData.(map[string]interface{}); ok {
			if name, ok := team["name"].(string); ok {
				return name
			}
		}
	}
	return side
}
