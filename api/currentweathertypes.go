package main

// datatype for current weather data to send to open-meteo
type currWeatherQuery struct {
	Latitude          float64  `json:"latitude"`
	Longitude         float64  `json:"longitude"`
	Hourly            []string `json:"hourly"`
	Current           []string `json:"current"`
	Timezone          string   `json:"timezone"`
	PastDays          int      `json:"past_days"`
	WindSpeedUnit     string   `json:"wind_speed_unit"`
	TemperatureUnit   string   `json:"temperature_unit"`
	PrecipitationUnit string   `json:"precipitation_unit"`
}

// datatype for current weather data received from the API
type currentWeatherObject struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`

	CurrentUnits struct {
		Time                string `json:"time"`
		Interval            string `json:"interval"`
		Temperature2m       string `json:"temperature_2m"`
		IsDay               string `json:"is_day"`
		RelativeHumidity2m  string `json:"relative_humidity_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		Precipitation       string `json:"precipitation"`
		Rain                string `json:"rain"`
		Showers             string `json:"showers"`
		Snowfall            string `json:"snowfall"`
		WeatherCode         string `json:"weather_code"`
		CloudCover          string `json:"cloud_cover"`
		WindSpeed10m        string `json:"wind_speed_10m"`
		WindDirection10m    string `json:"wind_direction_10m"`
		WindGusts10m        string `json:"wind_gusts_10m"`
	} `json:"current_units"`

	Current struct {
		Time                string  `json:"time"`
		Interval            int     `json:"interval"`
		Temperature2m       float64 `json:"temperature_2m"`
		IsDay               int     `json:"is_day"`
		RelativeHumidity2m  int     `json:"relative_humidity_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		Precipitation       float64 `json:"precipitation"`
		Rain                float64 `json:"rain"`
		Showers             float64 `json:"showers"`
		Snowfall            float64 `json:"snowfall"`
		WeatherCode         int     `json:"weather_code"`
		CloudCover          int     `json:"cloud_cover"`
		WindSpeed10m        float64 `json:"wind_speed_10m"`
		WindDirection10m    int     `json:"wind_direction_10m"`
		WindGusts10m        float64 `json:"wind_gusts_10m"`
	} `json:"current"`

	HourlyUnits struct {
		Time                     string `json:"time"`
		Temperature2m            string `json:"temperature_2m"`
		Precipitation            string `json:"precipitation"`
		PrecipitationProbability string `json:"precipitation_probability"`
		CloudCover               string `json:"cloud_cover"`
		Visibility               string `json:"visibility"`
		WindGusts10m             string `json:"wind_gusts_10m"`
		RelativeHumidity2m       string `json:"relative_humidity_2m"`
		ApparentTemperature      string `json:"apparent_temperature"`
		WindDirection10m         string `json:"wind_direction_10m"`
		WindSpeed10m             string `json:"wind_speed_10m"`
	} `json:"hourly_units"`

	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2m            []float64 `json:"temperature_2m"`
		Precipitation            []float64 `json:"precipitation"`
		PrecipitationProbability []int     `json:"precipitation_probability"`
		CloudCover               []int     `json:"cloud_cover"`
		Visibility               []float64 `json:"visibility"`
		WindGusts10m             []float64 `json:"wind_gusts_10m"`
		RelativeHumidity2m       []int     `json:"relative_humidity_2m"`
		ApparentTemperature      []float64 `json:"apparent_temperature"`
		WindDirection10m         []int     `json:"wind_direction_10m"`
		WindSpeed10m             []float64 `json:"wind_speed_10m"`
	} `json:"hourly"`
}
