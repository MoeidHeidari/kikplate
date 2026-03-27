"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { plateRepository } from "@/src/data/repositories/PlateRepository"
import { GetPlatesUseCase }                          from "@/src/domain/usecases/GetPlatesUseCase"
import { GetPlateUseCase }                           from "@/src/domain/usecases/GetPlateUseCase"
import { SubmitRepositoryUseCase } from "@/src/domain/usecases/SubmitPlateUseCase"
import { SetBookmarkUseCase }                          from "@/src/domain/usecases/SetBookmarkUseCase"
import { RatePlateUseCase }                          from "@/src/domain/usecases/RatePlateUseCase"
import type { PlateFilter, SubmitRepositoryInput } from "@/src/domain/entities/Plate"
import { GetStatsUseCase } from "@/src/domain/usecases/GetStatsUseCase"

const getPlates  = new GetPlatesUseCase(plateRepository)
const getPlate   = new GetPlateUseCase(plateRepository)
const submitRepo = new SubmitRepositoryUseCase(plateRepository)
const setBookmark  = new SetBookmarkUseCase(plateRepository)
const ratePlate  = new RatePlateUseCase(plateRepository)
const getStats = new GetStatsUseCase(plateRepository)

export function useStats() {
  return useQuery({
    queryKey: ["stats"],
    queryFn: () => getStats.execute(),
    staleTime: 5 * 60_000,
  })
}
  
export function usePlates(filter: PlateFilter & { bookmarked?: boolean } = {}) {
  const isBookmarkQuery = filter.bookmarked
  const filterWithoutBookmarked = { ...filter }
  delete filterWithoutBookmarked.bookmarked
  
  return useQuery({
    queryKey: ["plates", filter],
    queryFn: () => 
      isBookmarkQuery 
        ? plateRepository.listBookmarked(filter.limit)
        : getPlates.execute(filterWithoutBookmarked),
    staleTime: 30_000,
  })
}

export function usePlateFilterOptions() {
  return useQuery({
    queryKey: ["plate-filter-options"],
    queryFn: () => plateRepository.getFilterOptions(),
    staleTime: 60_000,
  })
}

export function usePlate(slug: string) {
  return useQuery({
    queryKey: ["plate", slug],
    queryFn: () => getPlate.execute(slug),
    enabled: Boolean(slug),
  })
}

export function useSubmitRepository() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: SubmitRepositoryInput) => submitRepo.execute(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plates"] }),
  })
}

export function useVerifyRepository() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => plateRepository.verifyRepository(id),
    onSuccess: (plate) => {
      qc.invalidateQueries({ queryKey: ["plates"] })
      qc.invalidateQueries({ queryKey: ["plate", plate.slug] })
    },
  })
}

export function useRemovePlate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => plateRepository.remove(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["plates"] })
    },
  })
}

export function useMovePlateOrganization() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, organizationId }: { id: string; organizationId?: string }) =>
      plateRepository.moveToOrganization(id, organizationId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["plates"] })
    },
  })
}

export function useSetBookmark() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, bookmarked }: { id: string; bookmarked: boolean }) => setBookmark.execute(id, bookmarked),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["plates"] })
    },
  })
}

export function useRatePlate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, rating }: { id: string; slug: string; rating: number }) => ratePlate.execute(id, rating),
    onSuccess: (_data, variables) => {
      qc.invalidateQueries({ queryKey: ["plates"] })
      qc.invalidateQueries({ queryKey: ["plate", variables.slug] })
    },
  })
}
