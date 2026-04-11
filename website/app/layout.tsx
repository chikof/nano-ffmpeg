import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "nano-ffmpeg - A beautiful terminal UI for ffmpeg",
  description:
    "Every ffmpeg feature, zero flags to remember. A beginner-friendly TUI dashboard for ffmpeg built in Go.",
  keywords: ["ffmpeg", "terminal", "tui", "cli", "video", "converter", "go"],
  openGraph: {
    title: "nano-ffmpeg",
    description: "Every ffmpeg feature. Zero flags to remember.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased dark`}
    >
      <body className="min-h-full flex flex-col bg-[#0a0a0f] text-[#f9fafb]">
        <Navbar />
        <main className="flex-1">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
