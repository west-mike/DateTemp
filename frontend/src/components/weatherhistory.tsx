'use client'

import { useState, useEffect } from 'react';
import YearlySideScroll from './yearlysidescroll';
import { HistoricalWeatherEntry } from './yearlysidescroll';

interface WeatherHistoryProps {
    initialData: {
        weather_data: HistoricalWeatherEntry[]
    };
    locationName?: string;
    onLocationChange: (lat: number, lon: number, name: string) => void;
}

const WeatherHistory = ({ initialData, locationName = "Ann Arbor" }: WeatherHistoryProps) => {
    // State to hold the historical data
    const [historicalData, setHistoricalData] = useState(initialData);

    // Update state when initialData changes
    useEffect(() => {
        setHistoricalData(initialData);
    }, [initialData]);

    return (
        <YearlySideScroll
            historicalData={historicalData}
            locationName={locationName}
        />
    );
};

export default WeatherHistory;