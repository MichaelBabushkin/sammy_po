package scraper

import (
	"testing"
)

func TestScrapeEvents(t *testing.T) {
	events, err := ScrapeEvents()
	if err != nil {
		t.Fatalf("Failed to scrape events: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("Expected events, but got none")
	}

	// Example: Validate the first event's fields
	firstEvent := events[0]
	if firstEvent.Name == "" {
		t.Errorf("Event Name is empty")
	}

	if firstEvent.DateTime.IsZero() {
		t.Errorf("Event DateTime is not set")
	}
}
