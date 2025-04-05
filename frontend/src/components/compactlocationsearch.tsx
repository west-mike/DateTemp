'use client'

import { useState } from 'react';

interface CompactLocationSearchProps {
    onLocationChange: (latitude: number, longitude: number, locationName: string) => void;
}

const CompactLocationSearch: React.FC<CompactLocationSearchProps> = ({ onLocationChange }) => {
    const [latitude, setLatitude] = useState<string>('42.28471');
    const [longitude, setLongitude] = useState<string>('-83.67496');
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [isExpanded, setIsExpanded] = useState<boolean>(false);

    // Sanitize input to ensure only valid numerical values are accepted
    const sanitizeCoordinateInput = (input: string): string => {
        // Allow only digits, decimal point, and minus sign
        // Remove any other characters that could be used for XSS
        return input.replace(/[^\d.-]/g, '');
    };

    const handleLatitudeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const sanitizedValue = sanitizeCoordinateInput(e.target.value);
        setLatitude(sanitizedValue);
    };

    const handleLongitudeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const sanitizedValue = sanitizeCoordinateInput(e.target.value);
        setLongitude(sanitizedValue);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Sanitize the inputs again before parsing
        const sanitizedLat = sanitizeCoordinateInput(latitude);
        const sanitizedLong = sanitizeCoordinateInput(longitude);

        // Parse the sanitized values
        const lat = parseFloat(sanitizedLat);
        const long = parseFloat(sanitizedLong);

        // Validate parsed values
        if (isNaN(lat) || isNaN(long)) {
            setError('Please enter valid numbers');
            return;
        }

        // Validate ranges
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

        // Set new location in parent component
        onLocationChange(lat, long, `${lat.toFixed(4)}, ${long.toFixed(4)}`);

        setIsLoading(false);
        setIsExpanded(false); // Collapse after submit
    };

    return (
        <div className="absolute top-14 left-4 z-10">
            {isExpanded ? (
                <div className="bg-white/20 backdrop-blur-sm rounded-lg p-3 shadow-lg w-64">
                    <form onSubmit={handleSubmit} className="space-y-2">
                        <div className="flex justify-between items-center mb-1">
                            <span className="text-sm font-medium">Search Location</span>
                            <button
                                type="button"
                                onClick={() => setIsExpanded(false)}
                                className="text-xs text-white/80 hover:text-white"
                            >
                                Close
                            </button>
                        </div>

                        <div className="grid grid-cols-2 gap-2">
                            <div>
                                <input
                                    type="text"
                                    value={latitude}
                                    onChange={handleLatitudeChange}
                                    className="w-full px-2 py-1 text-sm bg-white/20 rounded text-white focus:outline-none focus:ring-1 focus:ring-white/50"
                                    placeholder="Latitude"
                                    aria-label="Latitude"
                                />
                            </div>

                            <div>
                                <input
                                    type="text"
                                    value={longitude}
                                    onChange={handleLongitudeChange}
                                    className="w-full px-2 py-1 text-sm bg-white/20 rounded text-white focus:outline-none focus:ring-1 focus:ring-white/50"
                                    placeholder="Longitude"
                                    aria-label="Longitude"
                                />
                            </div>
                        </div>

                        {error && (
                            <div className="text-red-300 text-xs">
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            disabled={isLoading}
                            className="w-full py-1 text-sm bg-white/30 hover:bg-white/40 rounded font-medium transition-colors disabled:opacity-50"
                        >
                            {isLoading ? 'Loading...' : 'Search'}
                        </button>
                    </form>
                </div>
            ) : (
                <button
                    onClick={() => setIsExpanded(true)}
                    className="bg-white/20 backdrop-blur-sm px-3 py-1 rounded-lg text-sm hover:bg-white/30 transition-colors"
                >
                    Search Location
                </button>
            )}
        </div>
    );
};

export default CompactLocationSearch;