"use client";

import type { Metadata } from "next";
//import type {RootLayoutProps} from "next"
//import { Inter as FontSans } from "next/font/google";
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

import { ThemeProvider } from "@/components/theme-provider"
import { SessionProvider } from '@/lib/session';
import { cn } from "@/lib/utils"
import "./globals.css";
import TopBar from "@/components/top-bar";
import CssBaseline from "@mui/material/CssBaseline";
import ScrollTop from "@/components/scroll-top";
import Paper from "@mui/material/Paper";
import { UIContextProvider } from "@/lib/ui.context";
import ErrorBoundary from "@/components/error";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </head>
      <body className={cn(
        "min-h-screen bg-background font-sans antialiased",
        /*fontSans.className*/
      )}>
      <ErrorBoundary>
        <UIContextProvider>
          <SessionProvider>
            <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
              <CssBaseline />
              <TopBar anchorId='scroll-to-top' />
              <Paper  variant="outlined" sx={{ mt: 3, mr: 4, ml: 4, mb: 3, padding: 4 }} >{children}</Paper>
              <ScrollTop anchorId='scroll-to-top' />
            </ThemeProvider>
          </SessionProvider>
        </UIContextProvider>
      </ErrorBoundary>
      </body>
    </html>
  );
}