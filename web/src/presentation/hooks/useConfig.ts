"use client"

import { useQuery } from "@tanstack/react-query"
import { configRepository } from "@/src/data/repositories/ConfigRepository"
import { GetConfigUseCase } from "@/src/domain/usecases/GetConfigUseCase"

const getConfig = new GetConfigUseCase(configRepository)

export function useConfig() {
  return useQuery({
    queryKey: ["config"],
    queryFn: () => getConfig.execute(),
    staleTime: 10 * 60_000,
  })
}
