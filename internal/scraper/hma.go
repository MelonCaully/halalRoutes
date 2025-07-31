package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/MelonCaully/halalRoutes/internal/database"
	"github.com/gocolly/colly/v2"
)

type Restaurant struct {
	Name      string
	Region    string
	Address   string
	Website   string
	Source    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ScraperHMA(ctx context.Context, db *database.Queries) error {
	startURL := "https://hmacanada.org/hma-certified-restaurants/"

	c := colly.NewCollector(
		colly.AllowedDomains("hmacanada.org"),
		colly.UserAgent("Mozilla/5.0"),
	)

	var currentRegion string

	c.OnHTML("h4.wp-block-heading", func(e *colly.HTMLElement) {
		header := strings.TrimSpace(e.Text)
		if header != "" {
			currentRegion = header
		}
	})

	c.OnHTML("ul li", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.Text)
		if name == "" || currentRegion == "" {
			return
		}

		err := db.CreateRestaurant(ctx, database.CreateRestaurantParams{
			Name:      name,
			Region:    currentRegion,
			Address:   sql.NullString{Valid: false},
			Website:   "",
			Source:    "HMA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			log.Fatal("unable to create parameters for resteraunt")
		}
	})

	return c.Visit(startURL)
}
