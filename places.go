package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (cfg *apiConfig) PlacesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	location := r.URL.Query().Get("location") // e.g. "35.6895,139.6917" for Tokyo

	if query == "" || location == "" {
		respondWithError(w, http.StatusBadRequest, "Missing query or location parameters", nil)
		return
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		respondWithError(w, http.StatusInternalServerError, "Google Maps API key not configured", nil)
		return
	}

	endpoint := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?keyword=%s&location=%s&radius=3000&type=restaurant&key=%s",
		url.QueryEscape(query), location, apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Places API request failed", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondWithError(w, http.StatusBadRequest, "Places API returned error", nil)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read Places API response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
