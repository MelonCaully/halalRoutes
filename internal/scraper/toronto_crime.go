package scraper

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MelonCaully/halalRoutes/internal/database"
)

func ScraperCrimeToronto(ctx context.Context, db *database.Queries) error {
	url := "https://ckan0.cf.opendata.inter.prod-toronto.ca/dataset/21db0f45-1828-4fa3-94de-db92f454314c/resource/3c3925de-3a85-476a-85ca-b3cdff91b47f/download/neighbourhood-crime-rates%20-%204326.csv"
	year := "2024"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	reader := csv.NewReader(resp.Body)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return err
	}

	// Map header to index
	colIndex := map[string]int{}
	for i, h := range headers {
		colIndex[strings.TrimSpace(h)] = i
	}

	log.Println("Starting crime scraping...")

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("CSV read error:", err)
			continue
		}

		neighborhood := record[colIndex["AREA_NAME"]]

		assaultCrimeStr := record[colIndex["ASSAULT_"+year]]
		autoTheftCrimeStr := record[colIndex["AUTOTHEFT_"+year]]
		bikeTheftCrimeStr := record[colIndex["BIKETHEFT_"+year]]
		breakEnterCrimeStr := record[colIndex["BREAKENTER_"+year]]
		homicideCrimeStr := record[colIndex["HOMICIDE_"+year]]
		robberyCrimeStr := record[colIndex["ROBBERY_"+year]]
		shootingCrimeStr := record[colIndex["SHOOTING_"+year]]
		theftRommvCrimeStr := record[colIndex["THEFTROMMV_"+year]]
		theftOverCrimeStr := record[colIndex["THEFTOVER_"+year]]

		violentCrime, err := parseAndSum(assaultCrimeStr, homicideCrimeStr, shootingCrimeStr, robberyCrimeStr)
		if err != nil {
			log.Printf("Skipping invalid violent crime count for %s: %v", neighborhood, err)
			continue
		}

		propertyCrime, err := parseAndSum(autoTheftCrimeStr, bikeTheftCrimeStr, breakEnterCrimeStr, theftOverCrimeStr, theftRommvCrimeStr)
		if err != nil {
			log.Printf("Skipping invalid property crime count for %s: %v", neighborhood, err)
			continue
		}

		totalCrime := violentCrime + propertyCrime

		_, err = db.CreateCrimeStats(ctx, database.CreateCrimeStatsParams{
			Neighborhood:  neighborhood,
			TotalCrime:    int32(totalCrime),
			ViolentCrime:  int32(violentCrime),
			PropertyCrime: int32(propertyCrime),
			Source:        "Toronto Open Data",
		})
		if err != nil {
			log.Printf("Failed to insert stats for %s: %v", neighborhood, err)
		}
	}

	log.Println("✅ Crime scraping completed.")
	return nil
}

func parseAndSum(crimeStrs ...string) (int, error) {
	sum := 0
	for _, s := range crimeStrs {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		sum += v
	}
	return sum, nil
}
