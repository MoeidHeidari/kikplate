"use client"

import NextTopLoader from "nextjs-toploader"

/**
 * Thin top bar during client-side navigations (Next.js App Router).
 * Theme primary color; no corner spinner so it stays a clean bar.
 */
export function NavigationProgress() {
  return (
    <NextTopLoader
      color="var(--primary)"
      height={3}
      showSpinner={false}
      crawlSpeed={180}
      speed={280}
      shadow={false}
      zIndex={99999}
    />
  )
}
