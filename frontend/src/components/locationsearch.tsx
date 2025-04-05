'use client'

import { useState } from 'react';
import YearlySideScroll from './yearlysidescroll';
import { HistoricalWeatherEntry } from './yearlysidescroll';

interface YearlyHourData {
    weather_data: HistoricalWeatherEntry[]
}

export default function LocationSearch() {
    const [latitude, setLatitude] = useState<string>('42.28471');
    const [longitude, setLongitude] = useState<string>('-83.67496');
    const [customData, setCustomData] = useState<YearlyHourData | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validate inputs
        const lat = parseFloat(latitude);
        const long = parseFloat(longitude);

        if (isNaN(lat) || isNaN(long)) {
            setError('Please enter valid numbers for latitude and longitude');
            return;
        }

        if (lat < -90 || lat > 90) {
            setError('Latitude must be between -90 and 90');
            return;
        }

        if (long < -180 || long > 180) {
            setError('Longitude must be between -180 and 180');
            return;
        }

        setError(null);
        setIsLoading(true);

        try {
            // Get current hour and date for the API call
            const currentHour = new Date().toLocaleTimeString('en-US', { hour12: false }).split(':')[0] + ":00:00";
            const curDate = new Date().toISOString().split('T')[0];

            // Use the correct Next.js API route
            const response = await fetch(
                `/api/weather-history?latitude=${lat}&longitude=${long}&hour=${currentHour}&date=${curDate}`,
                { cache: 'no-store' }
            );

            if (!response.ok) {
                throw new Error('Failed to fetch data');
            }

            const data = await response.json();
            console.log(data);
            setCustomData(data);
        } catch (err) {
            setError('Error fetching weather data. Please try again.');
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="w-full max-w-md mb-6">
            <div className="p-4 bg-white/20 backdrop-blur-sm rounded-lg">
                <h2 className="text-xl font-semibold mb-3">Check Weather History for Another Location</h2>

                <form onSubmit={handleSubmit} className="space-y-3">
                    <div className="flex flex-col sm:flex-row gap-3">
                        <div className="flex-1">
                            <label htmlFor="latitude" className="block text-sm mb-1">Latitude</label>
                            <input
                                id="latitude"
                                type="text"
                                value={latitude}
                                onChange={(e) => setLatitude(e.target.value)}
                                className="w-full px-3 py-2 bg-white/20 rounded text-white focus:outline-none focus:ring-2 focus:ring-white/50"
                                placeholder="42.28471"
                            />
                        </div>

                        <div className="flex-1">
                            <label htmlFor="longitude" className="block text-sm mb-1">Longitude</label>
                            <input
                                id="longitude"
                                type="text"
                                value={longitude}
                                onChange={(e) => setLongitude(e.target.value)}
                                className="w-full px-3 py-2 bg-white/20 rounded text-white focus:outline-none focus:ring-2 focus:ring-white/50"
                                placeholder="-83.67496"
                            />
                        </div>
                    </div>

                    {error && (
                        <div className="text-red-300 text-sm">
                            {error}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full py-2 bg-white/30 hover:bg-white/40 rounded font-medium transition-colors disabled:opacity-50"
                    >
                        {isLoading ? 'Loading...' : 'Get Weather History'}
                    </button>
                </form>
            </div>

            {customData && (
                <div className="mt-6">
                    <h3 className="text-lg font-semibold">
                        Custom Location: {parseFloat(latitude).toFixed(4)}, {parseFloat(longitude).toFixed(4)}
                    </h3>
                    <YearlySideScroll historicalData={customData} />
                </div>
            )}
        </div>
    );
}