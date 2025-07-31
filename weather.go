package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}

func (cfg *apiConfig) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		respondWithError(w, http.StatusBadRequest, "City parameter is required", nil)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		respondWithError(w, http.StatusInternalServerError, "Weather API key not configured", nil)
		return
	}

	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s",
		city, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch weather data", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondWithError(w, http.StatusBadRequest, "Invalid city or weather service error", nil)
		return
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode weather response", err)
		return
	}

	response := map[string]interface{}{
		"city":    city,
		"temp":    weather.Main.Temp,
		"summary": weather.Weather[0].Main,
	}

	respondWithJSON(w, http.StatusOK, response)
}
