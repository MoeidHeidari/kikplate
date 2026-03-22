import { http } from "./httpClient"
import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"
import type {
  Plate, PlateFilter, PlateListResponse,
  SubmitFileInput, SubmitRepositoryInput,
} from "@/src/domain/entities/Plate"
import { PlateStats } from "@/src/domain/entities/Stats"

class PlateRepository implements IPlateRepository {
  list(filter: PlateFilter): Promise<PlateListResponse> {
    return http.get("/plates", filter as Record<string, unknown>)
  }
  getBySlug(slug: string): Promise<Plate> {
    return http.get(`/plates/${slug}`)
  }
  submitFile(input: SubmitFileInput): Promise<Plate> {
    return http.post("/plates/file", input)
  }
  submitRepository(input: SubmitRepositoryInput): Promise<Plate> {
    return http.post("/plates/repository", input)
  }
  recordUse(id: string): Promise<void> {
    return http.post(`/plates/${id}/use`)
  }
  getStats(): Promise<PlateStats> {
  return http.get("/plates/stats")
}
}


export const plateRepository = new PlateRepository()
