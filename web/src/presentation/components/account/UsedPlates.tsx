"use client"

import { Heart } from "lucide-react"
import Link from "next/link"

export function BookmarkedPlates() {
  return (
    <div className="flex flex-col items-center gap-3 py-20 text-muted-foreground">
      <Heart className="h-8 w-8 opacity-30" />
      <p className="text-sm">Your bookmarked plates will appear here.</p>
      <p className="text-xs text-muted-foreground/60">Bookmark plates to save them for later.</p>
      <Link
        href="/explore"
        className="text-sm text-foreground underline underline-offset-4"
      >
        Explore plates →
      </Link>
    </div>
  )
}