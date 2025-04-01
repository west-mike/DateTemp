package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// base api url for open-meteo current forecast api
var current_weather_url = "https://api.open-meteo.com/v1/forecast"

// hard-code lat and long for ann arbor for now, maybe others in future
var a2_lat float64 = 42.2808
var a2_long float64 = -83.732124

var refreshCurWeatherQuery = currWeatherQuery{
	Latitude:          a2_lat,
	Longitude:         a2_long,
	Hourly:            []string{"temperature_2m", "precipitation", "precipitation_probability", "cloud_cover", "visibility", "wind_gusts_10m", "relative_humidity_2m", "apparent_temperature", "wind_direction_10m", "wind_speed_10m"},
	Current:           []string{"temperature_2m", "is_day", "relative_humidity_2m", "apparent_temperature", "precipitation", "rain", "showers", "snowfall", "weather_code", "cloud_cover", "wind_speed_10m", "wind_direction_10m", "wind_gusts_10m"},
	Timezone:          "America/New_York",
	PastDays:          3,
	WindSpeedUnit:     "mph",
	TemperatureUnit:   "fahrenheit",
	PrecipitationUnit: "inch",
}
var history2000Query = historyWeatherQuery{
	Latitude:  a2_lat,
	Longitude: a2_long,
	StartDate: "2000-01-01",
	// in future make this a variable thats today -5 days
	EndDate:           "2025-03-27",
	Daily:             []string{"weather_code", "temperature_2m_mean", "temperature_2m_max", "temperature_2m_min", "apparent_temperature_mean", "apparent_temperature_max", "apparent_temperature_min", "sunrise", "sunset", "precipitation_sum", "precipitation_hours"},
	Hourly:            []string{"temperature_2m", "relative_humidity_2m", "apparent_temperature", "precipitation", "weather_code"},
	WindSpeedUnit:     "mph",
	TemperatureUnit:   "fahrenheit",
	PrecipitationUnit: "inch",
}
var currentWeatherData = currentWeatherObject{}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	// Test DB Connection
	//testConnection()
	router := gin.Default()
	router.GET("/currentweather", getCurrentWeather)
	router.GET("/dailyhistory", grabDailyHistory)
	router.GET("/populatehistory", populateHistory)
	router.GET("/history", populateHistory)
	// history migration routes, probably never need to use again
	router.GET("/history/migrate/hourly/:filename", migrateHourlyHistoryFile)
	router.GET("/history/migrate/daily/:filename", migrateDailyHistoryFile)
	router.Run("localhost:8080")
}

// Test the database connection and close it
func testConnection() {
	// Code from Supabase site to connect to DB
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test query
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		conn.Close(context.Background())
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	// Close the connection after test
	conn.Close(context.Background())
}
func getCurrentWeather(c *gin.Context) {
	// want to query open-meteo for the weather right now
	// theoretically, this call runs every 5 minutes or so
	// we can also package in the updates for the weather forecast for 24 hours to save a call

	// create http client
	client := &http.Client{}
	// establish base request headers and target
	req, err := http.NewRequest("GET", current_weather_url, nil)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Add query parameters to the request
	query := req.URL.Query()
	query.Add("latitude", fmt.Sprintf("%f", refreshCurWeatherQuery.Latitude))
	query.Add("longitude", fmt.Sprintf("%f", refreshCurWeatherQuery.Longitude))
	query.Add("hourly", strings.Join(refreshCurWeatherQuery.Hourly, ","))
	query.Add("current", strings.Join(refreshCurWeatherQuery.Current, ","))
	query.Add("timezone", refreshCurWeatherQuery.Timezone)
	query.Add("past_days", fmt.Sprintf("%d", refreshCurWeatherQuery.PastDays))
	query.Add("wind_speed_unit", refreshCurWeatherQuery.WindSpeedUnit) // Changed from "windspeed_unit" to "wind_speed_unit"
	query.Add("temperature_unit", refreshCurWeatherQuery.TemperatureUnit)
	query.Add("precipitation_unit", refreshCurWeatherQuery.PrecipitationUnit)
	req.URL.RawQuery = query.Encode()

	// Debug: Print the complete URL being requested
	if gin.IsDebugging() {
		fmt.Printf("[DEBUG] Request URL: %s\n", req.URL.String())
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch weather data"})
		return
	}
	// close the response body when done, presents resource leak
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		// Read the error response body
		var apiError struct {
			Error  bool   `json:"error"`
			Reason string `json:"reason"`
		}

		// Try to decode the error response
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			// Log the error to console in debug mode
			if gin.IsDebugging() {
				fmt.Printf("[ERROR] Failed to decode API error: %v\n", err)
				fmt.Printf("[ERROR] Status code: %d\n", resp.StatusCode)
			}

			// If we can't decode the error, return a generic one with the status code
			c.IndentedJSON(http.StatusBadGateway, gin.H{
				"error":       "unexpected response from weather API",
				"status_code": resp.StatusCode,
			})
			return
		}

		// Log the API error to console in debug mode
		if gin.IsDebugging() {
			fmt.Printf("[ERROR] Weather API error: %s\n", apiError.Reason)
			fmt.Printf("[ERROR] Full details: %+v\n", apiError)
		}

		// Return the actual error from the API with the same status code
		c.IndentedJSON(resp.StatusCode, gin.H{
			"error":   apiError.Reason,
			"details": apiError,
		})
		return
	} else {
		// Log successful response in debug mode
		if gin.IsDebugging() {
			// Read the response body
			fmt.Printf("[DEBUG] Successful response from weather API: %s\n", resp.Status)
		}
	}

	// Decode the response body into currentWeatherData
	if err := json.NewDecoder(resp.Body).Decode(&currentWeatherData); err != nil {
		if gin.IsDebugging() {
			fmt.Printf("[ERROR] Failed to decode response: %v\n", err)
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to decode weather data"})
		return
	}
	// Return the current weather data as JSON
	c.IndentedJSON(http.StatusOK, currentWeatherData)
	//TODO: Reformat the hourly data in a json array with each hour and then each category within that hour
	///TODO: Convert Weather codes to human-readable strings
}

// this function grabs each daily weather record from the daily db
func grabDailyHistory(c *gin.Context) {
	// Get the database connection
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}
	defer conn.Close(context.Background())

	// Get current month and day
	now := time.Now()
	month := now.Month()
	day := now.Day()
	// Get latitude and longitude from query parameters
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	if latitude == "" || longitude == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude parameters are required"})
		return
	}

	// Execute the SQL query using the provided latitude and longitude
	rows, err := conn.Query(context.Background(),
		`SELECT * FROM "Daily Historical Weather 2000" 
		 WHERE EXTRACT(MONTH FROM date) = $1
		 AND EXTRACT(DAY FROM date) = $2
		 AND latitude = $3
		 AND longitude = $4`,
		month, day, latitude, longitude)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("query failed: %v", err)})
		return
	}
	defer rows.Close()

	// Create a slice to hold the results
	var results []map[string]interface{}

	// Iterate over the rows
	for rows.Next() {
		// Create a map to hold the row data
		rowData := make(map[string]interface{})

		// Get column names
		columnNames := rows.FieldDescriptions()
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to scan row: %v", err)})
			return
		}

		// Map the values to their column names
		for i, col := range columnNames {
			rowData[string(col.Name)] = values[i]
		}

		// Append the row data to the results slice
		results = append(results, rowData)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("row iteration error: %v", err)})
		return
	}

	// Wrap the results in a JSON property called weather_data
	c.IndentedJSON(http.StatusOK, gin.H{"weather_data": results})
}

// grab all weather from 2000-now and populate the history in the database for a given location
func populateHistory(c *gin.Context) {
	// Log the start of the operation
	if gin.IsDebugging() {
		fmt.Println("[INFO] Starting historical weather data download for Ann Arbor...")
	}

	// Create output directory if it doesn't exist
	historyDir := "history/2000"
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create history directory"})
		return
	}

	// URL for historical weather data
	historyURL := "https://archive-api.open-meteo.com/v1/archive"

	// Create HTTP client with a longer timeout due to the large response size
	client := &http.Client{
		Timeout: 120 * time.Second, // 2 minute timeout for this large request
	}

	// Create request
	req, err := http.NewRequest("GET", historyURL, nil)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Build query parameters
	query := req.URL.Query()
	query.Add("latitude", fmt.Sprintf("%f", a2_lat))
	query.Add("longitude", fmt.Sprintf("%f", a2_long))
	query.Add("start_date", history2000Query.StartDate)
	query.Add("end_date", history2000Query.EndDate)
	query.Add("daily", strings.Join(history2000Query.Daily, ","))
	query.Add("hourly", strings.Join(history2000Query.Hourly, ","))
	query.Add("wind_speed_unit", history2000Query.WindSpeedUnit)
	query.Add("temperature_unit", history2000Query.TemperatureUnit)
	query.Add("precipitation_unit", history2000Query.PrecipitationUnit)
	req.URL.RawQuery = query.Encode()

	// Log the request URL
	if gin.IsDebugging() {
		fmt.Printf("[DEBUG] Historical data request URL: %s\n", req.URL.String())
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch historical weather data: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error  bool   `json:"error"`
			Reason string `json:"reason"`
		}

		// Try to decode the error response
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			if gin.IsDebugging() {
				fmt.Printf("[ERROR] Failed to decode API error: %v\n", err)
				fmt.Printf("[ERROR] Status code: %d\n", resp.StatusCode)
			}
			c.IndentedJSON(http.StatusBadGateway, gin.H{
				"error":       "unexpected response from weather API",
				"status_code": resp.StatusCode,
			})
			return
		}

		// Log and return the error
		if gin.IsDebugging() {
			fmt.Printf("[ERROR] Weather API error: %s\n", apiError.Reason)
		}
		c.IndentedJSON(resp.StatusCode, gin.H{
			"error":   apiError.Reason,
			"details": apiError,
		})
		return
	}

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to read response body"})
		return
	}

	// Verify it's valid JSON before saving (optional)
	var jsonCheck interface{}
	if err := json.Unmarshal(body, &jsonCheck); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "received invalid JSON from API"})
		return
	}

	// Save response to file
	filePath := filepath.Join(historyDir, "ann_arbor.json")
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to save historical data to file"})
		return
	}

	// Log success
	if gin.IsDebugging() {
		fmt.Printf("[INFO] Successfully saved historical weather data to %s\n", filePath)
		fmt.Printf("[INFO] File size: %.2f MB\n", float64(len(body))/(1024*1024))
	}

	// Return success to the client
	c.IndentedJSON(http.StatusOK, gin.H{
		"message":         "Historical weather data successfully downloaded and saved",
		"file_path":       filePath,
		"file_size_bytes": len(body),
	})
}

// migrateHistoryFile handles the /history/migrate/:filename endpoint
// It extracts the filename and calls sendHistoryToDB
func migrateHourlyHistoryFile(c *gin.Context) {
	// Get filename from URL parameter
	filename := c.Param("filename")

	if filename == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "filename parameter is required"})
		return
	}

	// Log the operation if in debug mode
	if gin.IsDebugging() {
		fmt.Printf("[INFO] Migrating history file: %s\n", filename)
	}

	// Define the base directory where history files are stored
	// Adjust this path if your files are stored elsewhere
	baseDir := "history/2000"

	// Construct the full file path
	filePath := filepath.Join(baseDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("file not found: %s", filename)})
		return
	}

	// Call sendHistoryToDB with the filename
	SendHourlyHistoryToDB(filePath)

	// Return success response
	c.IndentedJSON(http.StatusOK, gin.H{
		"message":   "Historical weather data successfully migrated to database",
		"file_path": filePath,
	})
}

func migrateDailyHistoryFile(c *gin.Context) {
	// Get filename from URL parameter
	filename := c.Param("filename")

	if filename == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "filename parameter is required"})
		return
	}

	// Log the operation if in debug mode
	if gin.IsDebugging() {
		fmt.Printf("[INFO] Migrating daily history file: %s\n", filename)
	}

	// Define the base directory where history files are stored
	// Adjust this path if your files are stored elsewhere
	baseDir := "history/2000"

	// Construct the full file path
	filePath := filepath.Join(baseDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("file not found: %s", filename)})
		return
	}

	// Call sendHistoryToDB with the filename
	SendDailyHistoryToDB(filePath)

	// Return success response
	c.IndentedJSON(http.StatusOK, gin.H{
		"message":   "Historical Daily weather data successfully migrated to database",
		"file_path": filePath,
	})
}
