import React from 'react';
import HourlyWeatherCard from './hourlyweathercard';

// Define the interface that matches the hourly data format
interface CurrentHourlySideScrollProps {
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
    units?: {
        temperature_2m: string;
        precipitation: string;
        wind_speed_10m: string;
    };
}

const CurrentHourlySideScroll: React.FC<CurrentHourlySideScrollProps> = ({ hourly, units = { temperature_2m: '°F', precipitation: 'in', wind_speed_10m: 'mp/h' } }) => {
    // Get current time
    const currentTime = new Date();

    // Filter to only include future hours
    const futureHours = hourly.time
        .map((timeString, index) => {
            const hourTime = new Date(timeString);
            return {
                time: hourTime,
                index
            };
        })
        .filter(item => item.time > currentTime)
        .slice(0, 24); // Limit to next 24 hours

    return (
        <div className="w-full 2xl:mt-6 md:mt-2">
            <h2 className="text-xl font-semibold mb-3">Hourly Forecast</h2>

            {/* Horizontal scrollable container */}
            <div className="flex gap-4 overflow-x-auto pb-4 forecast-container" >
                {futureHours.map(({ index, time }) => (
                    <HourlyWeatherCard
                        key={index}
                        time={time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        temperature={hourly.temperature_2m[index]}
                        tempUnit={units.temperature_2m}
                        precipitationProbability={hourly.precipitation_probability[index]}
                        precipitation={hourly.precipitation[index]}
                        cloudCover={hourly.cloud_cover[index]}
                    />
                ))}
            </div>
        </div>
    );
};

export default CurrentHourlySideScroll;