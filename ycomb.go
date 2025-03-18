package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	// Store articles
	var articles []map[string]string

	// Scrape titles and links
	c.OnHTML(".titleline a", func(e *colly.HTMLElement) {
		articles = append(articles, map[string]string{
			"title": e.Text,
			"link":  e.Attr("href"),
		})
	})

	// Error handling
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	// Start scraping
	fmt.Println("Scraping Hacker News...")
	err := c.Visit("https://news.ycombinator.com/")
	if err != nil {
		log.Fatal(err)
	}

	// Print results
	for _, article := range articles {
		fmt.Printf("Title: %s\nLink: %s\n\n", article["title"], article["link"])
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		log.Fatal("Error converting to JSON:", err)
	}

	// Create filename with current date using relative path
	filename := filepath.Join("data", fmt.Sprintf("articles_%s.json", 
		time.Now().Format("2006-01-02")))

	// Write to file in data directory
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing JSON file:", err)
	}

	log.Printf("Successfully scraped %d articles and saved to %s\n", len(articles), filename)
}
