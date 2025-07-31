package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func BuildPrompt(city string, area string, userType string, preferences []string, halal int, prayer int, bars int, vibe string) string {
	return fmt.Sprintf(`
You are a smart, inclusive travel assistant. Help users decide whether to visit %s in %s based on their needs.

User Profile:
- Traveler type: %s
- Preferences: %s

Area Data:
- Halal restaurants nearby: %d
- Prayer spaces nearby: %d
- Bars/nightlife spots: %d
- General vibe: %s

Instructions:
1. Write a friendly, travel-style summary of the area.
2. Highlight any concerns, mismatches, or red flags (e.g., loud nightlife, lack of religious facilities).
3. End with a recommendation on whether this area fits the traveler's style or not — and why.
Be concise, helpful, and objective.
`, area, city, userType, strings.Join(preferences, ", "), halal, prayer, bars, vibe)
}

func (cfg *apiConfig) RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		City        string   `json:"city"`
		Area        string   `json:"area"`
		UserType    string   `json:"user_type"`
		Preferences []string `json:"preferences"`
		Halal       int      `json:"halal"`
		Prayer      int      `json:"prayer"`
		Bars        int      `json:"bars"`
		Vibe        string   `json:"vibe"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate required fields
	if req.City == "" || req.Area == "" || req.UserType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields: city, area, or user_type", nil)
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		respondWithError(w, http.StatusInternalServerError, "Gemini API key not configured", nil)
		return
	}

	prompt := BuildPrompt(req.City, req.Area, req.UserType, req.Preferences, req.Halal, req.Prayer, req.Bars, req.Vibe)

	geminiReq := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create request", err)
		return
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gemini API request failed", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondWithError(w, http.StatusBadRequest, "Gemini API returned error", nil)
		return
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode Gemini response", err)
		return
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		respondWithError(w, http.StatusInternalServerError, "No recommendation generated", nil)
		return
	}

	recommendation := geminiResp.Candidates[0].Content.Parts[0].Text

	response := map[string]interface{}{
		"city":           req.City,
		"area":           req.Area,
		"user_type":      req.UserType,
		"recommendation": recommendation,
	}

	respondWithJSON(w, http.StatusOK, response)
}
