import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
    // Get query parameters
    const searchParams = request.nextUrl.searchParams;
    const latitude = searchParams.get('latitude');
    const longitude = searchParams.get('longitude');

    // Log environment variables and their values
    console.log('======================= WEATHER API DEBUG =======================');
    console.log('Environment Variables:');
    console.log('API_BASE_URL:', process.env.API_BASE_URL);
    console.log('NODE_ENV:', process.env.NODE_ENV);

    // Get the API base URL from environment variables
    const baseUrl = process.env.API_BASE_URL || 'http://localhost:8080';
    console.log('Using baseUrl:', baseUrl);

    if (!latitude || !longitude) {
        console.log('Missing required parameters');
        return NextResponse.json(
            { error: 'Missing required parameters' },
            { status: 400 }
        );
    }

    try {
        // Log the full URL being used
        const fullUrl = `${baseUrl}/currentweather?latitude=${latitude}&longitude=${longitude}`;
        console.log('Fetching from URL:', fullUrl);

        // Forward the request to the backend
        const response = await fetch(
            fullUrl,
            { cache: 'no-store' }
        );

        console.log('Response status:', response.status);

        if (!response.ok) {
            console.log('Response not OK:', response.statusText);
            return NextResponse.json(
                { error: 'Failed to fetch weather data' },
                { status: response.status }
            );
        }

        const data = await response.json();
        console.log('Successfully received data');
        console.log('==============================================================');
        return NextResponse.json(data);
    } catch (error) {
        console.error('Error fetching weather data:', error);
        console.log('==============================================================');
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        );
    }
}