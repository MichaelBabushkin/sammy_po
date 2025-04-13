package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MichaelBabushkin/sammy_po/api"
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
		
		fotmobClient := NewFotmobClient()
		matchesData, err := fotmobClient.FetchIsraeliLeagueMatches()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		matchesArr, ok := matchesData.([]interface{})
		if !ok {
			http.Error(w, "Invalid matches data format", http.StatusInternalServerError)
			return
		}
		
		filteredMatches := FilterHaifaHomeMatches(matchesArr)
		
		// Get all upcoming matches
		now := time.Now().UTC()
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
				upcomingMatches = append(upcomingMatches, match)
			}
		}
		
		log.Printf("Found %d upcoming Sammy Ofer matches", len(upcomingMatches))
		
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
