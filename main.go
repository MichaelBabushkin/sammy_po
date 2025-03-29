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

	"github.com/joho/godotenv"
)

type Message struct {
	Text string `json:"text"`
}

type APIResponse struct {
	Get        string      `json:"get"`
	Parameters interface{} `json:"parameters"`
	Errors     interface{} `json:"errors"`
	Results    int         `json:"results"`
	Paging     struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []Fixture `json:"response"`
}

type Fixture struct {
	Fixture struct {
		ID        int       `json:"id"`
		Referee   string    `json:"referee"`
		Timezone  string    `json:"timezone"`
		Date      time.Time `json:"date"`
		Timestamp int       `json:"timestamp"`
		Venue     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"venue"`
		Status struct {
			Long    string `json:"long"`
			Short   string `json:"short"`
			Elapsed int    `json:"elapsed"`
		} `json:"status"`
	} `json:"fixture"`
	League struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Logo    string `json:"logo"`
		Flag    string `json:"flag"`
		Season  int    `json:"season"`
		Round   string `json:"round"`
	} `json:"league"`
	Teams struct {
		Home struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"home"`
		Away struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home int `json:"home"`
		Away int `json:"away"`
	} `json:"goals"`
	Score struct {
		Halftime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"halftime"`
		Fulltime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fulltime"`
		Extratime struct {
			Home interface{} `json:"home"`
			Away interface{} `json:"away"`
		} `json:"extratime"`
		Penalty struct {
			Home interface{} `json:"home"`
			Away interface{} `json:"away"`
		} `json:"penalty"`
	} `json:"score"`
}

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
	LeagueID       int     `json:"leagueID"`
}

type VenueResponse struct {
	Get        string  `json:"get"`
	Parameters struct {
		Name string `json:"name"`
	} `json:"parameters"`
	Results int     `json:"results"`
	Venues  []Venue `json:"response"`
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

type LeagueResponse struct {
	Get        string      `json:"get"`
	Parameters interface{} `json:"parameters"`
	Errors     interface{} `json:"errors"`
	Results    int         `json:"results"`
	Paging     struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []LeagueData `json:"response"`
}

type LeagueData struct {
	League struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Logo    string `json:"logo"`
		Country string `json:"country"`
	} `json:"league"`
	Country struct {
		Name string `json:"name"`
		Code string `json:"code"`
		Flag string `json:"flag"`
	} `json:"country"`
	Seasons []struct {
		Year     int  `json:"year"`
		Start    string `json:"start"`
		End      string `json:"end"`
		Current  bool `json:"current"`
		Coverage struct {
			Fixtures struct {
				Events            bool `json:"events"`
				Lineups           bool `json:"lineups"`
				StatisticsFixtures bool `json:"statistics_fixtures"`
				StatisticsPlayers bool `json:"statistics_players"`
			} `json:"fixtures"`
			Standings   bool `json:"standings"`
			Players     bool `json:"players"`
			TopScorers  bool `json:"top_scorers"`
			TopAssists  bool `json:"top_assists"`
			TopCards    bool `json:"top_cards"`
			Injuries    bool `json:"injuries"`
			Predictions bool `json:"predictions"`
			Odds        bool `json:"odds"`
		} `json:"coverage"`
	} `json:"seasons"`
}

type LeagueInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Current bool   `json:"current"`
	Season  int    `json:"season"`
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

	req.Header.Add("x-mas", "eyJib2R5Ijp7InVybCI6Ii9hcGkvbGVhZ3Vlcz9pZD0xMjcmY2NvZGUzPUlTUiIsImNvZGUiOjE3NDMyNjY2MDA5MzgsImZvbyI6InByb2R1Y3Rpb246NDJlZWVlNmVlM2UzNmNhNDgyNmExMzkyYWIzMWE4ODk1YzNjODc0Yi11bmRlZmluZWQifSwic2lnbmF0dXJlIjoiMjRCMTJBMzg5MjRDRjdGMDMwRTA3QjQ1QkFDMjZFMDIifQ==")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	
	log.Printf("Request headers: %v", req.Header)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	log.Printf("Response status: %s", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	preview := string(body)
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	log.Printf("Response preview: %s", preview)

	return body, nil
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
// and returns upcoming matches first
func FilterHaifaHomeMatches(matches []interface{}) []interface{} {
	haifaMatches := []interface{}{}
	now := time.Now().UTC()
	
	// First pass: collect all Haifa home matches
	for _, match := range matches {
		matchMap, ok := match.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Get home team
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
		
		// Check if it's a Haifa team
		if strings.Contains(homeTeamName, "Maccabi Haifa") || 
		   strings.Contains(homeTeamName, "Hapoel Haifa") {
			haifaMatches = append(haifaMatches, match)
		}
	}
	
	// Second pass: separate upcoming and past matches
	upcomingMatches := []interface{}{}
	pastMatches := []interface{}{}
	
	for _, match := range haifaMatches {
		matchMap := match.(map[string]interface{})
		
		// Get status data
		statusData, ok := matchMap["status"]
		if !ok {
			continue
		}
		
		status, ok := statusData.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Check if it has utcTime
		utcTimeStr, ok := status["utcTime"].(string)
		if !ok {
			// Fallback to timeTS if utcTime is not available
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
		
		// Parse the UTC time string
		// Handle both formats: with and without milliseconds
		var matchTime time.Time
		var err error
		
		if strings.Contains(utcTimeStr, ".") {
			matchTime, err = time.Parse("2006-01-02T15:04:05.000Z", utcTimeStr)
		} else {
			matchTime, err = time.Parse("2006-01-02T15:04:05Z", utcTimeStr)
		}
		
		if err != nil {
			log.Printf("Error parsing UTC time '%s': %v", utcTimeStr, err)
			continue
		}
		
		// Check if the match is in the future
		if matchTime.After(now) {
			upcomingMatches = append(upcomingMatches, match)
		} else {
			pastMatches = append(pastMatches, match)
		}
	}
	
	// Sort upcoming matches by date (earliest first)
	for i := 0; i < len(upcomingMatches)-1; i++ {
		for j := i + 1; j < len(upcomingMatches); j++ {
			matchA := upcomingMatches[i].(map[string]interface{})
			matchB := upcomingMatches[j].(map[string]interface{})
			
			statusA, okA := matchA["status"].(map[string]interface{})
			statusB, okB := matchB["status"].(map[string]interface{})
			
			if !okA || !okB {
				continue
			}
			
			utcTimeStrA, okA := statusA["utcTime"].(string)
			utcTimeStrB, okB := statusB["utcTime"].(string)
			
			if !okA || !okB {
				continue
			}
			
			timeA, errA := time.Parse(time.RFC3339, utcTimeStrA)
			timeB, errB := time.Parse(time.RFC3339, utcTimeStrB)
			
			if errA == nil && errB == nil && timeA.After(timeB) {
				upcomingMatches[i], upcomingMatches[j] = upcomingMatches[j], upcomingMatches[i]
			}
		}
	}
	
	// Return upcoming matches first, then past matches
	return append(upcomingMatches, pastMatches...)
}

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}
}

func main() {
	http.HandleFunc("/api/greeting", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		json.NewEncoder(w).Encode(Message{Text: "Hello from Go Backend!"})
	})

	http.HandleFunc("/api/fotmob/league", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request for Fotmob league data")
		
		fotmobClient := NewFotmobClient()
		
		onlyMatches := r.URL.Query().Get("matches_only") == "true"
		
		var responseData interface{}
		var err error
		
		if onlyMatches {
			responseData, err = fotmobClient.FetchIsraeliLeagueMatches()
			log.Println("Fetching only matches data")
		} else {
			responseData, err = fotmobClient.FetchIsraeliLeagueData()
			log.Println("Fetching full league data")
		}

		if err != nil {
			log.Printf("Error fetching Fotmob data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully fetched Fotmob data")
		
		w.Header().Set("Cache-Control", "max-age=3600")
		
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/fotmob/matches", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request for Fotmob matches")
		
		// Parse query parameters for filtering
		teamFilter := r.URL.Query().Get("team")
		homeOnly := r.URL.Query().Get("home") == "true"
		
		fotmobClient := NewFotmobClient()
		matchesData, err := fotmobClient.FetchIsraeliLeagueMatches()

		if err != nil {
			log.Printf("Error fetching Fotmob matches: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Apply filters if specified
		if teamFilter != "" {
			log.Printf("Filtering matches for team: %s, homeOnly: %v", teamFilter, homeOnly)
			
			matchesArr, ok := matchesData.([]interface{})
			if ok {
				// Filter for Maccabi Haifa or Hapoel Haifa if requested
				if teamFilter == "haifa" {
					maccabiHaifa := FilterMatches(matchesArr, "Maccabi Haifa", homeOnly)
					hapoelHaifa := FilterMatches(matchesArr, "Hapoel Haifa", homeOnly)
					
					// Combine both filtered results
					combined := append(maccabiHaifa, hapoelHaifa...)
					matchesData = combined
					
					log.Printf("Found %d matches for Haifa teams", len(combined))
				} else {
					filtered := FilterMatches(matchesArr, teamFilter, homeOnly)
					matchesData = filtered
					log.Printf("Found %d matches for team: %s", len(filtered), teamFilter)
				}
			}
		}

		log.Printf("Successfully fetched and processed Fotmob matches")
		
		w.Header().Set("Cache-Control", "max-age=3600")
		
		if err := json.NewEncoder(w).Encode(matchesData); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})

	// Add a specialized endpoint for Haifa teams
	http.HandleFunc("/api/fotmob/haifa", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request for Haifa teams matches")
		
		// Parse query parameter for home matches only
		homeOnly := r.URL.Query().Get("home") == "true"
		
		fotmobClient := NewFotmobClient()
		matchesData, err := fotmobClient.FetchIsraeliLeagueMatches()

		if err != nil {
			log.Printf("Error fetching Fotmob matches for Haifa teams: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Filter for Haifa teams
		matchesArr, ok := matchesData.([]interface{})
		if !ok {
			http.Error(w, "Invalid matches data format", http.StatusInternalServerError)
			return
		}
		
		// Get matches for both Haifa teams
		maccabiHaifa := FilterMatches(matchesArr, "Maccabi Haifa", homeOnly)
		hapoelHaifa := FilterMatches(matchesArr, "Hapoel Haifa", homeOnly)
		
		// Combine both filtered results
		combined := append(maccabiHaifa, hapoelHaifa...)
		
		log.Printf("Found %d matches for Haifa teams", len(combined))
		
		w.Header().Set("Cache-Control", "max-age=3600")
		
		if err := json.NewEncoder(w).Encode(combined); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
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
		if err := json.NewEncoder(w).Encode(stadiumInfo); err != nil {
			log.Printf("Error encoding stadium info: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})

	// Add specialized endpoint for Sammy Ofer matches (Haifa home games)
	// This endpoint now returns only upcoming matches by default
	http.HandleFunc("/api/fotmob/sammyofer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request for Sammy Ofer (Haifa home) matches")
		
		fotmobClient := NewFotmobClient()
		matchesData, err := fotmobClient.FetchIsraeliLeagueMatches()

		if err != nil {
			log.Printf("Error fetching Fotmob matches: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Filter for Haifa home games and sort by date
		matchesArr, ok := matchesData.([]interface{})
		if !ok {
			http.Error(w, "Invalid matches data format", http.StatusInternalServerError)
			return
		}
		
		filteredMatches := FilterHaifaHomeMatches(matchesArr)
		
		// Get all upcoming matches first
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
				log.Printf("Error parsing UTC time '%s': %v", utcTimeStr, err)
				continue
			}
			
			// Include only future matches
			if matchTime.After(now) {
				upcomingMatches = append(upcomingMatches, match)
			}
		}
		
		log.Printf("Found %d upcoming Sammy Ofer matches", len(upcomingMatches))
		
		w.Header().Set("Cache-Control", "max-age=3600") // Cache for 1 hour
		
		if err := json.NewEncoder(w).Encode(upcomingMatches); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})

	fs := http.FileServer(http.Dir("frontend/build"))
	http.Handle("/", fs)

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
