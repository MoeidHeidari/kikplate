"use client"

import { useQuery } from "@tanstack/react-query"
import { badgeRepository } from "@/src/data/repositories/BadgeRepository"
import { GetBadgesUseCase } from "@/src/domain/usecases/GetBadgesUseCase"

const getBadges = new GetBadgesUseCase(badgeRepository)

export function useBadges() {
  return useQuery({
    queryKey: ["badges"],
    queryFn: () => getBadges.execute(),
    staleTime: 10 * 60_000,
  })
}
