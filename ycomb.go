package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
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

	// Write to file
	err = os.WriteFile("articles.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing JSON file:", err)
	}

	fmt.Println("Data saved to articles.json")
	fmt.Println("Scraping Done!")
}
