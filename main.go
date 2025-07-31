package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MelonCaully/halalRoutes/internal/database"
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
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
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
	mux.HandleFunc("GET /api/healthz", handlerHealthz)                      // healthz endpoint
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)            // users endpoint
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)                  // login endpoint
	mux.HandleFunc("GET /api/weather", apiCfg.WeatherHandler)               // weather endpoint
	mux.HandleFunc("GET /api/places", apiCfg.PlacesHandler)                 // places endpoint
	mux.HandleFunc("POST /api/recommendations", apiCfg.RecommendationsHandler) // AI recommendations endpoint

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
