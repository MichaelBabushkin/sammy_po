package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Message struct {
    Text string `json:"text"`
}

type Match struct {
    ID          string `json:"id"`
    HomeTeam    string `json:"homeTeam"`
    AwayTeam    string `json:"awayTeam"`
    HomeScore   int    `json:"homeScore"`
    AwayScore   int    `json:"awayScore"`
    Date        string `json:"date"`
    League      string `json:"league"`
}

func getMockMatches() []Match {
    return []Match{
        {
            ID: "1", HomeTeam: "Manchester United", AwayTeam: "Liverpool",
            HomeScore: 2, AwayScore: 1, Date: "2024-03-10", League: "Premier League",
        },
        {
            ID: "2", HomeTeam: "Barcelona", AwayTeam: "Real Madrid",
            HomeScore: 3, AwayScore: 3, Date: "2024-03-09", League: "La Liga",
        },
        {
            ID: "3", HomeTeam: "Bayern Munich", AwayTeam: "Dortmund",
            HomeScore: 4, AwayScore: 0, Date: "2024-03-08", League: "Bundesliga",
        },
    }
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    // API endpoints
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

    http.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        matches := getMockMatches()
        json.NewEncoder(w).Encode(matches)
    })

    // Serve static files
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
