"use client"

import Link from "next/link"
import { useState } from "react"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { PlateGrid } from "@/src/presentation/components/plates/PlateGrid"
import { Heart } from "lucide-react"

export function BookmarkedPlates({ accountId }: { accountId: string }) {
  const { data, isLoading, isError } = usePlates({ bookmarked: true, limit: 48 })

  if (isLoading) {
    return <LoadingSpinner />
  }

  if (isError || !data) {
    return (
      <div className="flex flex-col items-center gap-3 py-20 text-muted-foreground">
        <Heart className="h-8 w-8 opacity-30" />
        <p className="text-sm">Failed to load bookmarked plates</p>
      </div>
    )
  }

  if (!data.data || data.data.length === 0) {
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

  return <PlateGrid plates={data.data} />
}
