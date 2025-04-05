import React from 'react';

interface HistoricalAverages {
    avgTemperature: number;
    avgFeelsLike: number;
    avgHumidity: number;
    avgPrecipitation: number;
    mostCommonWeatherCode: number;
    yearsOfData: number;
}

interface CurrentWeatherProps {
    temperature: number;
    condition: string;
    isDay: boolean;
    prec_type: string;
    prec_chance: string;
    historicalAverages?: HistoricalAverages | null;
    tempUnit?: string;
}

const CurrentWeather: React.FC<CurrentWeatherProps> = ({
    temperature,
    condition,
    isDay,
    prec_type,
    prec_chance,
    historicalAverages,
    tempUnit = '°F',
}) => {
    // Function to get weather icon based on condition and time of day
    const getWeatherIcon = () => {
        if (condition.includes("Clear") || condition.includes("Mainly clear")) {
            return isDay ? "☀️" : "🌙";
        } else if (condition.includes("cloud")) {
            return isDay ? "⛅" : "☁️";
        } else if (condition.includes("rain") || condition.includes("drizzle") || condition.includes("shower")) {
            return "🌧️";
        } else if (condition.includes("snow")) {
            return "❄️";
        } else if (condition.includes("fog")) {
            return "🌫️";
        } else if (condition.includes("Thunderstorm")) {
            return "⛈️";
        }
        return isDay ? "🌤️" : "🌃";
    };

    return (
        <div className="bg-white/20 backdrop-blur-md rounded-2xl p-6 w-full max-w-md text-center">
            <div className="flex justify-between items-center 2xl:mb-4 md:mb-2">
                <div className="text-4xl">{getWeatherIcon()}</div>
                <h2 className="text-2xl font-bold">{isDay ? "Day" : "Night"}</h2>
            </div>

            {/* Current temperature */}
            <div className="text-8xl font-bold mb-2">{Math.round(temperature)}{tempUnit}</div>
            <div className="text-xl mb-4">{condition}</div>

            {/* Precipitation info */}
            {prec_type !== "None" && (
                <div className="mb-4">
                    <span className="text-lg">{prec_type}</span>
                    <span className="text-lg ml-2">({prec_chance} chance)</span>
                </div>
            )}

            {/* Historical average section */}
            {historicalAverages && (
                <div className="2xl:mt-6 md:mt-2 2xl:pt-4 border-t border-white/30">
                    <h3 className="text-lg font-semibold mb-2">
                        Historical Average ({historicalAverages.yearsOfData} years)
                    </h3>

                    <div className="flex justify-between items-center">
                        <div>
                            <div className="text-4xl font-bold">
                                {Math.round(historicalAverages.avgTemperature)}{tempUnit}
                            </div>
                            <div className="text-sm">Average</div>
                        </div>

                        <div className="text-right">
                            <div className="text-sm mb-1">
                                Feels like: {Math.round(historicalAverages.avgFeelsLike)}{tempUnit}
                            </div>
                            <div className="text-sm mb-1">
                                Humidity: {Math.round(historicalAverages.avgHumidity)}%
                            </div>
                        </div>
                    </div>

                    <div className="2xl:mt-2 text-sm">
                        {historicalAverages.avgTemperature > temperature ?
                            "Colder than average" :
                            historicalAverages.avgTemperature < temperature ?
                                "Warmer than average" :
                                "Average temperature"}
                    </div>
                </div>
            )}
        </div>
    );
};

export default CurrentWeather;