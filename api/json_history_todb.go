package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OpenMeteo historical weather data structure
type HistoricalHourlyWeatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Hourly    struct {
		Time                []string  `json:"time"`
		Temperature         []float64 `json:"temperature_2m"`
		RelativeHumidity    []float64 `json:"relative_humidity_2m"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
		Precipitation       []float64 `json:"precipitation"`
		WeatherCode         []int     `json:"weather_code"`
	} `json:"hourly"`
}

type HistoricalDailyWeatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Daily     struct {
		Time                    []string  `json:"time"`
		WeatherCode             []int     `json:"weather_code"`
		TemperatureMean         []float64 `json:"temperature_2m_mean"`
		TemperatureMax          []float64 `json:"temperature_2m_max"`
		TemperatureMin          []float64 `json:"temperature_2m_min"`
		ApparentTemperatureMean []float64 `json:"apparent_temperature_mean"`
		ApparentTemperatureMax  []float64 `json:"apparent_temperature_max"`
		ApparentTemperatureMin  []float64 `json:"apparent_temperature_min"`
		Sunrise                 []string  `json:"sunrise"`
		Sunset                  []string  `json:"sunset"`
		PrecipitationSum        []float64 `json:"precipitation_sum"`
		PrecipitationHours      []float64 `json:"precipitation_hours"`
	} `json:"daily"`
}

// Structure for each record to be sent to Supabase
type HourlyWeatherRecord struct {
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Hour                string  `json:"hour"` // Changed from int to string
	Date                string  `json:"date"`
	Temperature         float64 `json:"temperature"`
	RelativeHumidity    float64 `json:"relative_humidity"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	Precipitation       float64 `json:"precipitation"`
	WeatherCode         int     `json:"weather_code"`
}

// Structure for each record to be sent to Supabase
type DailyWeatherRecord struct {
	Latitude                float64 `json:"latitude"`
	Longitude               float64 `json:"longitude"`
	Date                    string  `json:"date"` // Format as "2006-01-02"
	WeatherCode             int     `json:"weather_code"`
	TemperatureMean         float64 `json:"mean_temperature"`
	TemperatureMax          float64 `json:"max_temperature"`
	TemperatureMin          float64 `json:"min_temperature"`
	ApparentTemperatureMean float64 `json:"mean_apparent_temperature"`
	ApparentTemperatureMax  float64 `json:"max_apparent_temperature"`
	ApparentTemperatureMin  float64 `json:"min_apparent_temperature"`
	Sunrise                 string  `json:"sunrise_time"` // Format as "15:04:05"
	Sunset                  string  `json:"sunset_time"`  // Format as "15:04:05"
	PrecipitationSum        float64 `json:"precipitation_sum"`
	PrecipitationHours      float64 `json:"precipitation_hours"`
}

// send the history data to the database
func SendHourlyHistoryToDB(filename string) {
	// Supabase REST API configuration
	supabaseURL := "https://bsztzmwbkzkmzhkepkwb.supabase.co/rest/v1/Hourly%20Historical%20Weather%20Since%202000"
	supabaseAPIKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJzenR6bXdia3prbXpoa2Vwa3diIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDMyODE2NDcsImV4cCI6MjA1ODg1NzY0N30.Q_SLkbdb-zGW1KOw8Uh0tfOfSOqZaeZfXWEDKMiW3nQ"

	// Read JSON file
	fmt.Printf("Reading file: %s\n", filename)
	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Parse the JSON data
	var weatherData HistoricalHourlyWeatherData
	if err := json.Unmarshal(fileData, &weatherData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed JSON data with %d hourly records\n", len(weatherData.Hourly.Time))

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create headers for Supabase request
	headers := map[string]string{
		"apikey":        supabaseAPIKey,
		"Authorization": "Bearer " + supabaseAPIKey,
		"Content-Type":  "application/json",
		"Prefer":        "return=minimal", // Don't return the created records
	}

	// Process records in batches to avoid overwhelming the API
	batchSize := 100
	totalRecords := len(weatherData.Hourly.Time)
	successCount := 0

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := make([]HourlyWeatherRecord, 0, end-i)

		// Create records for this batch
		for j := i; j < end; j++ {
			// Parse time from the API response
			timeStr := weatherData.Hourly.Time[j]
			t, err := time.Parse("2006-01-02T15:04", timeStr)
			if err != nil {
				fmt.Printf("Error parsing time %s: %v\n", timeStr, err)
				continue
			}

			// Extract date and hour
			date := t.Format("2006-01-02")
			hour := t.Format("15:04:05") // Format hour as proper time string HH:MM:SS

			// Create record
			record := HourlyWeatherRecord{
				Latitude:            weatherData.Latitude,
				Longitude:           weatherData.Longitude,
				Hour:                hour,
				Date:                date,
				Temperature:         weatherData.Hourly.Temperature[j],
				RelativeHumidity:    weatherData.Hourly.RelativeHumidity[j],
				ApparentTemperature: weatherData.Hourly.ApparentTemperature[j],
				Precipitation:       weatherData.Hourly.Precipitation[j],
				WeatherCode:         weatherData.Hourly.WeatherCode[j],
			}

			batch = append(batch, record)
		}

		// Convert batch to JSON
		batchJSON, err := json.Marshal(batch)
		if err != nil {
			fmt.Printf("Error marshaling batch to JSON: %v\n", err)
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", supabaseURL, bytes.NewBuffer(batchJSON))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		// Add headers
		for key, value := range headers {
			req.Header.Add(key, value)
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}

		// Check response
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			successCount += len(batch)
			fmt.Printf("Successfully inserted batch %d-%d of %d\n", i, end-1, totalRecords)
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error response (status %d): %s\n", resp.StatusCode, string(body))
		}

		// Add a small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Completed processing with %d/%d records successfully inserted\n", successCount, totalRecords)
}

func SendDailyHistoryToDB(filename string) {
	// Supabase REST API configuration
	supabaseURL := "https://bsztzmwbkzkmzhkepkwb.supabase.co/rest/v1/Daily%20Historical%20Weather%202000" // Change URL to point to daily table
	supabaseAPIKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJzenR6bXdia3prbXpoa2Vwa3diIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDMyODE2NDcsImV4cCI6MjA1ODg1NzY0N30.Q_SLkbdb-zGW1KOw8Uh0tfOfSOqZaeZfXWEDKMiW3nQ"

	// Read JSON file
	fmt.Printf("Reading file: %s\n", filename)
	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Parse the JSON data
	var weatherData HistoricalDailyWeatherData // Changed from HistoricalWeatherData
	if err := json.Unmarshal(fileData, &weatherData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed JSON data with %d daily records\n", len(weatherData.Daily.Time))

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create headers for Supabase request
	headers := map[string]string{
		"apikey":        supabaseAPIKey,
		"Authorization": "Bearer " + supabaseAPIKey,
		"Content-Type":  "application/json",
		"Prefer":        "return=minimal", // Don't return the created records
	}

	// Process records in batches to avoid overwhelming the API
	batchSize := 100
	totalRecords := len(weatherData.Daily.Time)
	successCount := 0

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := make([]DailyWeatherRecord, 0, end-i) // Changed to DailyWeatherRecord

		// Create records for this batch
		for j := i; j < end; j++ {
			// Parse date from the API response
			dateStr := weatherData.Daily.Time[j]
			date := dateStr // Daily times are already in "2025-03-23" format

			// Parse sunrise and sunset times
			sunriseStr := weatherData.Daily.Sunrise[j]
			sunsetStr := weatherData.Daily.Sunset[j]

			// Extract just the time portion from sunrise/sunset timestamps
			sunriseTime, err := time.Parse("2006-01-02T15:04", sunriseStr)
			if err != nil {
				fmt.Printf("Error parsing sunrise time %s: %v\n", sunriseStr, err)
				continue
			}
			sunrise := sunriseTime.Format("15:04:05")

			sunsetTime, err := time.Parse("2006-01-02T15:04", sunsetStr)
			if err != nil {
				fmt.Printf("Error parsing sunset time %s: %v\n", sunsetStr, err)
				continue
			}
			sunset := sunsetTime.Format("15:04:05")

			// Create record
			record := DailyWeatherRecord{
				Latitude:                weatherData.Latitude,
				Longitude:               weatherData.Longitude,
				Date:                    date,
				WeatherCode:             weatherData.Daily.WeatherCode[j],
				TemperatureMean:         weatherData.Daily.TemperatureMean[j],
				TemperatureMax:          weatherData.Daily.TemperatureMax[j],
				TemperatureMin:          weatherData.Daily.TemperatureMin[j],
				ApparentTemperatureMean: weatherData.Daily.ApparentTemperatureMean[j],
				ApparentTemperatureMax:  weatherData.Daily.ApparentTemperatureMax[j],
				ApparentTemperatureMin:  weatherData.Daily.ApparentTemperatureMin[j],
				Sunrise:                 sunrise,
				Sunset:                  sunset,
				PrecipitationSum:        weatherData.Daily.PrecipitationSum[j],
				PrecipitationHours:      weatherData.Daily.PrecipitationHours[j],
			}

			batch = append(batch, record)
		}

		// Convert batch to JSON
		batchJSON, err := json.Marshal(batch)
		if err != nil {
			fmt.Printf("Error marshaling batch to JSON: %v\n", err)
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", supabaseURL, bytes.NewBuffer(batchJSON))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		// Add headers
		for key, value := range headers {
			req.Header.Add(key, value)
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}

		// Check response
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			successCount += len(batch)
			fmt.Printf("Successfully inserted batch %d-%d of %d\n", i, end-1, totalRecords)
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error response (status %d): %s\n", resp.StatusCode, string(body))
		}

		// Add a small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Completed processing with %d/%d records successfully inserted\n", successCount, totalRecords)
}
