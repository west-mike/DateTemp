import React from 'react';
import YearlyWeatherCard from './yearlyweathercard';

// Define the interface for historical weather data
export interface HistoricalWeatherEntry {
    apparent_temperature: number;
    date: string;
    hour: {
        Microseconds: number;
        Valid: boolean;
    } | string;  // Update to accept string too
    id?: number;  // Make ID optional
    latitude: number;
    longitude: number;
    precipitation: number;
    relative_humidity: number;
    temperature: number;
    weather_code: number;
    year?: number;  // Some entries might have this
}

interface YearlySideScrollProps {
    historicalData: HistoricalWeatherEntry[] | { weather_data: HistoricalWeatherEntry[] };
    units?: {
        temperature: string;
        precipitation: string;
    };
    locationName?: string;
}

const YearlySideScroll: React.FC<YearlySideScrollProps> = ({
    historicalData,
    units = { temperature: '°F', precipitation: 'in' },
    locationName = "Historical Data"
}) => {
    // Extract the weather data correctly
    const weatherData = Array.isArray(historicalData)
        ? historicalData
        : historicalData?.weather_data || [];

    // Safety check to prevent errors
    if (!weatherData || weatherData.length === 0) {
        return (
            <div className="w-full 2xl:mt-6 m:mt-2">
                <h2 className="text-xl font-semibold mb-3">
                    Historical Weather on This Day
                    {locationName && <span className="text-sm font-normal ml-2">({locationName})</span>}
                </h2>
                <p>No historical data available</p>
            </div>
        );
    }

    // Sort data by year from oldest to newest
    const sortedData = [...weatherData].sort((a, b) => {
        return new Date(a.date).getTime() - new Date(b.date).getTime();
    });

    return (
        <div className="w-full 2xl:mt-6 m:mt-2">
            <h2 className="text-xl font-semibold mb-3">
                Historical Weather on This Day

            </h2>

            {/* Horizontal scrollable container */}
            <div className="flex gap-4 overflow-x-auto pb-4 forecast-container">
                {sortedData.map((entry, index) => {
                    // Generate a year from the date field
                    const year = entry.year || new Date(entry.date).getFullYear();

                    return (
                        <YearlyWeatherCard
                            // Create a reliable key using multiple fields or index as fallback
                            key={entry.id || `${entry.date}-${year}-${index}`}
                            year={year.toString()}
                            temperature={entry.temperature}
                            tempUnit={units.temperature}
                            precipitation={entry.precipitation}
                            precipitationUnit={units.precipitation}
                            weather_code={entry.weather_code}
                        />
                    );
                })}
            </div>
        </div>
    );
};

export default YearlySideScroll;