package scraper

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/MichaelBabushkin/sammy-po/internal/models"
)

// reverseLetters reverses the letters within each word without changing the word order.
func reverseLetters(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		// Check if the word consists only of digits
		if _, err := strconv.Atoi(word); err == nil {
			continue // Skip reversing if it's a digit-only word
		}

		// Reverse letters in the word
		runes := []rune(word)
		for j, k := 0, len(runes)-1; j < k; j, k = j+1, k-1 {
			runes[j], runes[k] = runes[k], runes[j]
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// normalizeHebrewText ensures proper Unicode rendering for Hebrew text.
func normalizeHebrewText(input string) string {
	normalized, _, _ := transform.String(norm.NFC, input)
	return normalized
}

// ScrapeEvents fetches the stadium's event page, parses the events, and returns them.
func ScrapeEvents() ([]models.Event, error) {
	// The URL you want to scrape
	url := "https://www.haifa-stadium.co.il/%d7%9c%d7%95%d7%97_%d7%94%d7%9e%d7%a9%d7%97%d7%a7%d7%99%d7%9d_%d7%91%d7%90%d7%a6%d7%98%d7%93%d7\x99%d7\x95%d7\x9f/"

	// 1. Fetch the page
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	// 2. Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []models.Event

	doc.Find("section.elementor-inner-section").Each(func(i int, section *goquery.Selection) {
		// Use the <p> tags within the columns to extract the data
		columns := section.Find("p")

		if columns.Length() < 4 {
			return
		}

		// Extract event data from <p> tags
		competition := strings.TrimSpace(columns.Eq(0).Text())
		team1 := strings.TrimSpace(columns.Eq(1).Text())
		dateStr := strings.TrimSpace(columns.Eq(2).Text())
		team2 := strings.TrimSpace(columns.Eq(3).Text())

		// Fix Hebrew text by reversing letters in each word
		competition = normalizeHebrewText(reverseLetters(competition))
		team1 = normalizeHebrewText(reverseLetters(team1))
		team2 = normalizeHebrewText(reverseLetters(team2))

		// Remove Hebrew day name using regex
		re := regexp.MustCompile(`^[\p{Hebrew}\s]+`)
		dateStr = re.ReplaceAllString(dateStr, " ")

		// Remove non-breaking spaces and other invisible characters
		dateStr = strings.ReplaceAll(dateStr, "\u00a0", "")
		dateStr = strings.ReplaceAll(dateStr, "\xc2\xa0", "")
		dateStr = strings.TrimSpace(dateStr)

		// Try parsing the date using multiple layouts
		var dateTime time.Time
		var err error
		layouts := []string{"02/01/2006 15:04", "02/01/06 15:04"}
		for _, layout := range layouts {
			dateTime, err = time.Parse(layout, dateStr)
			if err == nil {
				break // Exit loop if parsing succeeds
			}
		}

		if err != nil {
			log.Printf("Error parsing date %q: %v", dateStr, err)
			return
		}

		// Create an Event object
		event := models.Event{
			Name:        competition + ": " + team1 + " vs " + team2,
			Description: competition,
			DateTime:    dateTime,
			Location:    "Unknown", // Adjust if location is included elsewhere
		}

		events = append(events, event)
	})

	return events, nil
}
