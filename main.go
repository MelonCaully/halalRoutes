package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MelonCaully/halalRoutes/internal/database"
	"github.com/MelonCaully/halalRoutes/internal/handlers"
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

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database using '%s': %v", dbURL, err)
	}
	defer dbConn.Close()

	//dbQueries := database.New(dbConn)
	//apiCfg := apiConfig{
	//	db: dbQueries,
	//}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealthz) // healthz endpoint

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
