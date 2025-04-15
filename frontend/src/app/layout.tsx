import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { Analytics } from "@vercel/analytics/react"
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "DateTemp - Historical Weather App",
  description: "View current and historical weather data for any location",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased  min-h-screen relative`}
      >
        <div className="absolute top-4 left-4 lg:text-3xl md:text-xl s:text-lg font-bold text-white">
          DateTemp
        </div>
        <div className="absolute top-4 right-4 lg:text-3xl md:text-xl s:text-lg font-bold text-white">
          <a href="https://www.westmike.com" target="_blank" rel="noopener noreferrer">By: Michael West</a>
        </div>
        <main className="min-h-screen">
          {children}
          <Analytics />
        </main>
      </body>
    </html>
  );
}
