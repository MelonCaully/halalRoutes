package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func testAPIKeys() {
	fmt.Println("🔍 Testing API Keys...")
	fmt.Println("========================")

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ Error loading .env file:", err)
		return
	}

	// Test Weather API
	weatherKey := os.Getenv("WEATHER_API_KEY")
	if weatherKey != "" {
		fmt.Print("🌤️  Testing Weather API... ")
		resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=London&appid=%s", weatherKey))
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("❌ FAILED")
		} else {
			fmt.Println("✅ SUCCESS")
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	// Test Google Maps API
	mapsKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if mapsKey != "" {
		fmt.Print("🗺️  Testing Google Maps API... ")
		resp, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=51.5074,-0.1278&radius=1000&type=restaurant&key=%s", mapsKey))
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("❌ FAILED")
		} else {
			fmt.Println("✅ SUCCESS")
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	// Test Gemini API
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey != "" {
		fmt.Print("🤖 Testing Gemini AI API... ")
		// Simple test request
		testData := `{"contents":[{"parts":[{"text":"Say hello in one word"}]}]}`
		resp, err := http.Post(
			fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", geminiKey),
			"application/json",
			strings.NewReader(testData),
		)
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("❌ FAILED")
		} else {
			fmt.Println("✅ SUCCESS")
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	fmt.Println("========================")
	fmt.Println("✨ API Key testing complete!")
}

func TestApiKeys(t *testing.T) {
	testAPIKeys()
}
