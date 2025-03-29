package main

// datatype for current weather data to send to open-meteo
type historyWeatherQuery struct {
	Latitude          float64  `json:"latitude"`
	Longitude         float64  `json:"longitude"`
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	Daily             []string `json:"daily"`
	Hourly            []string `json:"hourly"`
	Timezone          string   `json:"timezone"`
	WindSpeedUnit     string   `json:"wind_speed_unit"`
	TemperatureUnit   string   `json:"temperature_unit"`
	PrecipitationUnit string   `json:"precipitation_unit"`
}

// datatype for current weather data received from the API
type historyWeatherObject struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`

	HourlyUnits struct {
		Time                string `json:"time"`
		Temperature2m       string `json:"temperature_2m"`
		RelativeHumidity2m  string `json:"relative_humidity_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		Precipitation       string `json:"precipitation"`
		WeatherCode         string `json:"weather_code"`
	} `json:"hourly_units"`

	Hourly struct {
		Time                []string  `json:"time"`
		Temperature2m       []float64 `json:"temperature_2m"`
		RelativeHumidity2m  []int     `json:"relative_humidity_2m"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
		Precipitation       []float64 `json:"precipitation"`
		WeatherCode         []int     `json:"weather_code"`
	} `json:"hourly"`

	DailyUnits struct {
		Time                    string `json:"time"`
		WeatherCode             string `json:"weather_code"`
		Temperature2mMean       string `json:"temperature_2m_mean"`
		Temperature2mMax        string `json:"temperature_2m_max"`
		Temperature2mMin        string `json:"temperature_2m_min"`
		ApparentTemperatureMean string `json:"apparent_temperature_mean"`
		ApparentTemperatureMax  string `json:"apparent_temperature_max"`
		ApparentTemperatureMin  string `json:"apparent_temperature_min"`
		Sunrise                 string `json:"sunrise"`
		Sunset                  string `json:"sunset"`
		PrecipitationSum        string `json:"precipitation_sum"`
		PrecipitationHours      string `json:"precipitation_hours"`
	} `json:"daily_units"`

	Daily struct {
		Time                    []string  `json:"time"`
		WeatherCode             []int     `json:"weather_code"`
		Temperature2mMean       []float64 `json:"temperature_2m_mean"`
		Temperature2mMax        []float64 `json:"temperature_2m_max"`
		Temperature2mMin        []float64 `json:"temperature_2m_min"`
		ApparentTemperatureMean []float64 `json:"apparent_temperature_mean"`
		ApparentTemperatureMax  []float64 `json:"apparent_temperature_max"`
		ApparentTemperatureMin  []float64 `json:"apparent_temperature_min"`
		Sunrise                 []string  `json:"sunrise"`
		Sunset                  []string  `json:"sunset"`
		PrecipitationSum        []float64 `json:"precipitation_sum"`
		PrecipitationHours      []int     `json:"precipitation_hours"`
	} `json:"daily"`
}
