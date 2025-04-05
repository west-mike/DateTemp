# About
Ever been fooled thinking spring is starting based on one warm day? Ever gotten snow in early October and been worried its winter already? DateTemp aims to help with this by letting you see the weather on the same day for the past 25 years, hopefully helping you avoid false optimism/pessimism.

# Backend Code
The backend is a RESTful API written in Go using a Supabase PostreSQL DB and the OpenMeteo Weather API to get weather data. It's deployed on GCP (Google Cloud Platform).

# Frontend Code
Frontend is written in TypeScript using the Next.js app router framework and TailwindCSS. It's deployed on Vercel and accessible [here](https://datetemp.westmike.com)