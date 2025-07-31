package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/MelonCaully/halalRoutes/internal/database"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func ScraperHMA(ctx context.Context, db *database.Queries) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var html string
	url := "https://hmacanada.org/hma-certified-restaurants/"

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML("body", &html, chromedp.ByQuery),
	)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return err
	}

	var currentRegion string

	doc.Find("h4.wp-block-heading, p").Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "h4" {
			currentRegion = strings.TrimSpace(s.Text())
			return
		}

		if currentRegion == "" {
			return
		}

		html, _ := s.Html()
		name := s.Find("strong").Text()
		if name == "" {
			return
		}

		address := ""
		parts := strings.Split(html, "<br>")
		if len(parts) > 1 {
			address = strings.TrimSpace(stripHTML(parts[1]))
		}

		err := db.CreateRestaurant(ctx, database.CreateRestaurantParams{
			Name:      strings.TrimSpace(name),
			Region:    currentRegion,
			Address:   sql.NullString{String: address, Valid: address != ""},
			Website:   "",
			Source:    "HMA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			log.Printf("DB insert error for %s: %v", name, err)
		}
	})

	log.Println("✅ Restaurant scraping completed.")
	return nil
}

// Simple HTML stripper
func stripHTML(input string) string {
	output := ""
	skip := false
	for _, r := range input {
		switch r {
		case '<':
			skip = true
		case '>':
			skip = false
		default:
			if !skip {
				output += string(r)
			}
		}
	}
	return output
}
