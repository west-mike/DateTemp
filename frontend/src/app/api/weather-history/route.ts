import { NextResponse } from 'next/server';

export async function GET(request: Request) {
    const { searchParams } = new URL(request.url);
    const latitude = searchParams.get('latitude');
    const longitude = searchParams.get('longitude');
    const hour = searchParams.get('hour');
    const date = searchParams.get('date');

    if (!latitude || !longitude || !hour || !date) {
        return NextResponse.json({ error: 'Missing required parameters' }, { status: 400 });
    }

    try {
        // Server-side request to your API (no CORS issues here)
        const apiRes = await fetch(
            `http://localhost:8080/hourlyhistory/nonDB?latitude=${latitude}&longitude=${longitude}&hour=${hour}&date=${date}`,
            { cache: 'no-store' }
        );

        if (!apiRes.ok) {
            throw new Error('Failed to fetch data from API');
        }

        const data = await apiRes.json();
        return NextResponse.json(data);
    } catch (error) {
        console.error('Error fetching weather data:', error);
        return NextResponse.json({ error: 'Failed to fetch weather data' }, { status: 500 });
    }
}