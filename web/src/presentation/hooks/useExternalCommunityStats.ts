"use client"

import { useQuery } from "@tanstack/react-query"

export interface ExternalCommunityStats {
  githubStars: number | null
  slackMembers: number | null
}

async function fetchExternalStats(): Promise<ExternalCommunityStats> {
  const res = await fetch("/api/stats/external", { cache: "no-store" })
  if (!res.ok) {
    return { githubStars: null, slackMembers: null }
  }
  return res.json()
}

export function useExternalCommunityStats() {
  return useQuery({
    queryKey: ["externalCommunityStats"],
    queryFn: fetchExternalStats,
    staleTime: 5 * 60_000,
  })
}
