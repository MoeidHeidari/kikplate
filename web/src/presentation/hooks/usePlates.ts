"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { plateRepository } from "@/src/data/repositories/PlateRepository"
import { GetPlatesUseCase }                          from "@/src/domain/usecases/GetPlatesUseCase"
import { GetPlateUseCase }                           from "@/src/domain/usecases/GetPlateUseCase"
import { SubmitFileUseCase, SubmitRepositoryUseCase } from "@/src/domain/usecases/SubmitPlateUseCase"
import { RecordUseUseCase }                          from "@/src/domain/usecases/RecordUseUseCase"
import type { PlateFilter, SubmitFileInput, SubmitRepositoryInput } from "@/src/domain/entities/Plate"
import { GetStatsUseCase } from "@/src/domain/usecases/GetStatsUseCase"

const getPlates  = new GetPlatesUseCase(plateRepository)
const getPlate   = new GetPlateUseCase(plateRepository)
const submitFile = new SubmitFileUseCase(plateRepository)
const submitRepo = new SubmitRepositoryUseCase(plateRepository)
const recordUse  = new RecordUseUseCase(plateRepository)
const getStats = new GetStatsUseCase(plateRepository)

export function useStats() {
  return useQuery({
    queryKey: ["stats"],
    queryFn: () => getStats.execute(),
    staleTime: 5 * 60_000,
  })
}
  
export function usePlates(filter: PlateFilter = {}) {
  return useQuery({
    queryKey: ["plates", filter],
    queryFn: () => getPlates.execute(filter),
    staleTime: 30_000,
  })
}

export function usePlate(slug: string) {
  return useQuery({
    queryKey: ["plate", slug],
    queryFn: () => getPlate.execute(slug),
    enabled: Boolean(slug),
  })
}

export function useSubmitFile() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: SubmitFileInput) => submitFile.execute(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plates"] }),
  })
}

export function useSubmitRepository() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: SubmitRepositoryInput) => submitRepo.execute(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plates"] }),
  })
}

export function useRecordUse() {
  return useMutation({
    mutationFn: (id: string) => recordUse.execute(id),
  })
}
