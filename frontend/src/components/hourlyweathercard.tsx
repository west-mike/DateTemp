import React from 'react';

interface HourlyWeatherCardProps {
    time: string;
    temperature: number;
    tempUnit: string;
    precipitationProbability: number;
    precipitation: number;
    cloudCover: number;
}

const HourlyWeatherCard: React.FC<HourlyWeatherCardProps> = ({
    time,
    temperature,
    tempUnit,
    precipitationProbability,
    precipitation,
    cloudCover
}) => {
    // Helper to get weather icon based on cloud cover and precipitation
    const getWeatherIcon = () => {
        if (precipitation > 0) {
            return "🌧️";
        } else if (cloudCover > 80) {
            return "☁️";
        } else if (cloudCover > 30) {
            return "⛅";
        } else {
            return "☀️";
        }
    };


    return (
        <div className="flex-shrink-0 bg-white/20 backdrop-blur-sm rounded-lg p-3 min-w-[100px] hover:bg-white/30 transition-colors hover:cursor-default">
            <div className="text-center font-medium">{time}</div>

            <div className="text-center my-2 text-3xl ">
                {getWeatherIcon()}
            </div>

            <div className="text-center font-bold text-2xl">
                {Math.round(temperature)}{tempUnit}
            </div>

            <div className="text-center text-sm mt-2">
                {precipitationProbability > 0 && (
                    <div>💧 {precipitationProbability}%</div>
                )}
            </div>

        </div>
    );
};

export default HourlyWeatherCard;