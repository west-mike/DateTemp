'use client'

import React, { useState, useEffect } from "react";
import CurrentWeather from "@/components/currentweather";
import CurrentHourlySideScroll from "@/components/currenthourlysidescroll";
import WeatherHistory from "@/components/weatherhistory";
import { HistoricalWeatherEntry } from "@/components/yearlysidescroll";
import CompactLocationSearch from "@/components/compactlocationsearch";

// Interfaces remain the same
interface WeatherData {
  current: {
    time: string;
    temperature_2m: number;
    is_day: number;
    relative_humidity_2m: number;
    precipitation: number;
    rain: number;
    showers: number;
    snowfall: number;
    weather_code: number;
  };
  current_units: {
    temperature_2m: string;
    precipitation: string;
    relative_humidity_2m: string;
  };
  hourly: {
    time: string[];
    temperature_2m: number[];
    precipitation: number[];
    precipitation_probability: number[];
    cloud_cover: number[];
    visibility: number[];
    wind_gusts_10m: number[];
    relative_humidity_2m: number[];
    apparent_temperature: number[];
    wind_direction_10m: number[];
    wind_speed_10m: number[];
  };
}

interface YearlyHourData {
  weather_data: HistoricalWeatherEntry[]
}

// Move all the utility functions outside the component
// Helper function to determine precipitation type
function getPrecipitationType(data: WeatherData): string {
  if (data.current.snowfall > 0) return "Snow";
  if (data.current.rain > 0) return "Rain";
  if (data.current.showers > 0) return "Showers";
  if (data.current.precipitation > 0) return "Precipitation";
  return "None";
}

function getPrecipitationChance(data: WeatherData): string {
  const currentHour = new Date().getHours();
  const currentHourIndex = currentHour % 24;
  return `${data.hourly.precipitation_probability[currentHourIndex]}%`;
}

function getWeatherCondition(code: number): string {
  const weatherCodes: { [key: number]: string } = {
    0: "Clear sky",
    1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
    45: "Fog", 48: "Depositing rime fog",
    51: "Light drizzle", 53: "Moderate drizzle", 55: "Dense drizzle",
    56: "Light freezing drizzle", 57: "Dense freezing drizzle",
    61: "Slight rain", 63: "Moderate rain", 65: "Heavy rain",
    66: "Light freezing rain", 67: "Heavy freezing rain",
    71: "Slight snow fall", 73: "Moderate snow fall", 75: "Heavy snow fall",
    77: "Snow grains",
    80: "Slight rain showers", 81: "Moderate rain showers", 82: "Violent rain showers",
    85: "Slight snow showers", 86: "Heavy snow showers",
    95: "Thunderstorm", 96: "Thunderstorm with slight hail", 99: "Thunderstorm with heavy hail"
  };

  return weatherCodes[code] || "Unknown";
}

function calculateHistoricalAverages(historicalData: HistoricalWeatherEntry[]) {
  if (!historicalData || historicalData.length === 0) {
    return null;
  }

  let totalTemp = 0;
  let totalFeelsLike = 0;
  let totalHumidity = 0;
  let totalPrecip = 0;
  let weatherCodeCounts: { [key: number]: number } = {};

  historicalData.forEach(entry => {
    totalTemp += entry.temperature;
    totalFeelsLike += entry.apparent_temperature;
    totalHumidity += entry.relative_humidity;
    totalPrecip += entry.precipitation;

    if (weatherCodeCounts[entry.weather_code]) {
      weatherCodeCounts[entry.weather_code]++;
    } else {
      weatherCodeCounts[entry.weather_code] = 1;
    }
  });

  let mostCommonCode = 0;
  let highestCount = 0;

  Object.entries(weatherCodeCounts).forEach(([code, count]) => {
    if (count > highestCount) {
      highestCount = count;
      mostCommonCode = parseInt(code);
    }
  });

  return {
    avgTemperature: totalTemp / historicalData.length,
    avgFeelsLike: totalFeelsLike / historicalData.length,
    avgHumidity: totalHumidity / historicalData.length,
    avgPrecipitation: totalPrecip / historicalData.length,
    mostCommonWeatherCode: mostCommonCode,
    yearsOfData: historicalData.length
  };
}

export default function Home() {
  // State for coordinates
  const [coordinates, setCoordinates] = useState({
    latitude: 42.28471,
    longitude: -83.67496,
    locationName: "Ann Arbor"
  });

  // State for weather data
  const [weatherData, setWeatherData] = useState<WeatherData | null>(null);
  const [yearlyHourData, setYearlyHourData] = useState<YearlyHourData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCustomLocation, setIsCustomLocation] = useState(false);

  // Function to fetch current weather
  async function fetchWeatherData(lat: number, lon: number) {
    try {
      // Use the Next.js API route instead of direct backend call
      const res = await fetch(`/api/weather?latitude=${lat}&longitude=${lon}`, {
        cache: 'no-store'
      });

      if (!res.ok) {
        throw new Error('Failed to fetch weather data');
      }

      return await res.json();
    } catch (error) {
      console.error('Error fetching current weather:', error);
      throw error;
    }
  }

  // Function to fetch historical data
  async function fetchYearlyHourData(lat: number, lon: number, useDB: boolean = true) {
    try {
      const currentHour = new Date().toLocaleTimeString('en-US', { hour12: false }).split(':')[0] + ":00:00";
      const curDate = new Date().toISOString().split('T')[0];

      // Determine which endpoint to use based on whether it's Ann Arbor (use DB) or custom location (use nonDB)
      const endpoint = useDB ? 'hourlyhistory' : 'hourlyhistory/nonDB';

      // Use the Next.js API route 
      const res = await fetch(
        `/api/historical?endpoint=${endpoint}&latitude=${lat}&longitude=${lon}&hour=${currentHour}&date=${curDate}`,
        { cache: 'no-store' }
      );

      if (!res.ok) {
        throw new Error('Failed to fetch historical data');
      }

      return await res.json();
    } catch (error) {
      console.error('Error fetching historical data:', error);
      throw error;
    }
  }

  // Handler for location changes
  const handleLocationChange = (lat: number, lon: number, name: string) => {
    // Additional validation as a security measure
    if (isNaN(lat) || isNaN(lon) ||
      lat < -90 || lat > 90 ||
      lon < -180 || lon > 180) {
      setError('Invalid coordinates provided');
      return;
    }

    // Sanitize the location name to prevent XSS
    const sanitizedName = name.replace(/[<>]/g, '');

    // Set loading to true immediately when location changes
    setIsLoading(true);

    // Clear existing data to force loading screen
    setWeatherData(null);
    setYearlyHourData(null);

    // Update coordinates
    setCoordinates({
      latitude: lat,
      longitude: lon,
      locationName: sanitizedName || `${lat.toFixed(4)}, ${lon.toFixed(4)}`
    });

    // Flag this as a custom location (should use nonDB)
    setIsCustomLocation(true);
  };

  // Effect to fetch data when coordinates change
  useEffect(() => {
    const fetchAllData = async () => {
      // No need to set loading here since we already did it in handleLocationChange
      setError(null);

      try {
        // Fetch both data types in parallel
        const [newWeatherData, newYearlyData] = await Promise.all([
          fetchWeatherData(coordinates.latitude, coordinates.longitude),
          fetchYearlyHourData(coordinates.latitude, coordinates.longitude, !isCustomLocation)
        ]);

        setWeatherData(newWeatherData);
        setYearlyHourData(newYearlyData);
      } catch (err) {
        setError('Failed to fetch weather data. Please try again.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAllData();
  }, [coordinates.latitude, coordinates.longitude, isCustomLocation]);

  // Show loading state - simplified to just check isLoading
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-6 bg-gradient-to-b from-blue-400 to-blue-600 text-white">
        <div className="text-2xl">Loading weather data...</div>
      </div>
    );
  }

  // Show error state
  if (error || !weatherData || !yearlyHourData) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-6 bg-gradient-to-b from-blue-400 to-blue-600 text-white">
        <div className="text-2xl text-red-300">
          {error || "Failed to load weather data"}
        </div>
      </div>
    );
  }

  // Calculate historical averages
  const historicalAverages = calculateHistoricalAverages(yearlyHourData.weather_data);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-6 bg-gradient-to-b from-blue-400 to-blue-600 text-white">
      {/* Location display */}
      <div className="text-xl mb-4">
        {coordinates.locationName}
      </div>

      {/* Add the compact search component */}
      <CompactLocationSearch onLocationChange={handleLocationChange} />

      <CurrentWeather
        temperature={weatherData.current.temperature_2m}
        condition={getWeatherCondition(weatherData.current.weather_code)}
        isDay={weatherData.current.is_day === 1 ? true : false}
        prec_type={getPrecipitationType(weatherData)}
        prec_chance={getPrecipitationChance(weatherData)}
        historicalAverages={historicalAverages}
        tempUnit={weatherData.current_units.temperature_2m}
      />

      <CurrentHourlySideScroll
        hourly={weatherData.hourly}
      />

      <WeatherHistory
        initialData={yearlyHourData}
        onLocationChange={handleLocationChange}
      />
    </div>
  );
}
