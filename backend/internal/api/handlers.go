// internal/api/handlers.go
package api

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/MichaelBabushkin/sammy-po/internal/models"
)

// Temporary in-memory mock data:
var mockEvents = []models.Event{
    {
        ID:          1,
        Name:        "Mock Event 1",
        Description: "This is a mock event",
        DateTime:    time.Now().Add(24 * time.Hour),
        Location:    "Main Stadium",
    },
    {
        ID:          2,
        Name:        "Mock Event 2",
        Description: "Another mock event",
        DateTime:    time.Now().Add(48 * time.Hour),
        Location:    "Main Stadium",
    },
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mockEvents)
}
