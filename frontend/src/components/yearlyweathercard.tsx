import React from 'react';

interface YearlyWeatherCardProps {
    year: string;
    tempUnit: string;
    precipitation: number;
    precipitationUnit: string;
    temperature: number;
    weather_code: number;
}

const YearlyWeatherCard: React.FC<YearlyWeatherCardProps> = ({
    year,
    tempUnit,
    precipitation,
    precipitationUnit,
    temperature,
    weather_code
}) => {
    // Helper to get weather icon based on weather code
    const getWeatherIcon = () => {
        switch (weather_code) {
            case 0: // Clear sky
                return "☀️";
            case 1: // Mainly clear
                return "⛅";
            case 2: // Partly cloudy
                return "🌥️";
            case 3: // Overcast
                return "☁️";
            case 45: // Fog
            case 48: // Depositing rime fog
                return "🌫️";
            case 51: // Drizzle: Light
            case 53: // Drizzle: Moderate
            case 55: // Drizzle: Dense intensity
                return "🌦️";
            case 61: // Rain: Slight
            case 63: // Rain: Moderate
            case 65: // Rain: Heavy intensity
                return "🌧️";
            case 71: // Snow fall: Slight
            case 73: // Snow fall: Moderate
            case 75: // Snow fall: Heavy intensity
                return "❄️";
            case 95: // Thunderstorm: Slight or moderate
                return "⛈️";
            case 96: // Thunderstorm with slight hail
            case 99: // Thunderstorm with heavy hail
                return "🌩️";
            default:
                return "❓"; // Unknown weather code
        }
    };


    return (
        <div className="flex-shrink-0 bg-white/20 backdrop-blur-sm rounded-lg p-3 min-w-[100px] hover:bg-white/30 transition-colors hover:cursor-default">
            <div className="text-center font-medium">{year}</div>

            <div className="text-center my-2 text-3xl">
                {getWeatherIcon()}
            </div>

            <div className="text-center font-bold text-2xl">
                {Math.round(temperature)}{tempUnit}
            </div>

        </div>
    );
};

export default YearlyWeatherCard;