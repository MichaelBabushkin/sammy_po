package main

import (
    "fmt"
    "log"

    "github.com/MichaelBabushkin/sammy-po/internal/scraper"
)

func main() {
    events, err := scraper.ScrapeEvents()
    if err != nil {
        log.Fatal(err)
    }

    for _, e := range events {
        fmt.Printf("Event: %s on %v at %s\n", e.Name, e.DateTime, e.Location)
    }
}
