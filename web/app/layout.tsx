import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { ThemeProvider } from "next-themes"
import { QueryProvider } from "@/src/presentation/providers/QueryProvider"
import { AuthCookieSync } from "@/src/presentation/providers/AuthCookieSync"
import { Navbar } from "@/src/presentation/components/layout/Navbar"
import { Footer } from "@/src/presentation/components/layout/Footer"
import { Toaster } from "@/components/ui/sonner"
import { NavigationProgress } from "@/src/presentation/components/common/NavigationProgress"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "kikplate — template registry",
  description: "Discover, share, and scaffold production-ready project templates.",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <QueryProvider>
            <NavigationProgress />
            <AuthCookieSync />
            <Navbar />
            <main className="min-h-screen">{children}</main>
            <Footer />
            <Toaster />
          </QueryProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}