package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MelonCaully/halalRoutes/internal/database"
	"github.com/MelonCaully/halalRoutes/internal/scraper"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	port := "8080"

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on Render-provided environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in enviroment")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM not set in enviroment")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database using '%s': %v", dbURL, err)
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handlerHealthz)                         // healthz endpoint
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)               // users endpoint
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)                     // login endpoint
	mux.HandleFunc("GET /api/weather", apiCfg.WeatherHandler)                  // weather endpoint
	mux.HandleFunc("GET /api/places", apiCfg.PlacesHandler)                    // places endpoint
	mux.HandleFunc("POST /api/recommendations", apiCfg.RecommendationsHandler) // AI recommendations endpoint
	mux.HandleFunc("GET /api/crime", apiCfg.handlerCrimeStats)                 // crime endpiont

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	if platform == "dev" {
		err = scraper.ScraperHMA(context.Background(), database.New(dbConn))
		if err != nil {
			log.Fatalf("Scraping failed for restaurants: %v", err)
		}
	}

	if platform == "dev" {
		err = scraper.ScraperCrimeToronto(context.Background(), database.New(dbConn))
		if err != nil {
			log.Fatalf("Scraping failed for crime: %v", err)
		}
	}

	log.Printf("Serving files on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
