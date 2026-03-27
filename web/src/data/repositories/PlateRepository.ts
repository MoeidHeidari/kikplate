import { http } from "./httpClient"
import { ApiError } from "./httpClient"
import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"
import type {
  Plate, PlateFilter, PlateFilterOptions, PlateListResponse,
  RatePlateInput, SubmitRepositoryInput,
} from "@/src/domain/entities/Plate"
import { PlateStats } from "@/src/domain/entities/Stats"

class PlateRepository implements IPlateRepository {
  list(filter: PlateFilter): Promise<PlateListResponse> {
    return http.get("/plates", filter as Record<string, unknown>)
  }
  listBookmarked(limit?: number): Promise<PlateListResponse> {
    const params = limit ? { limit } : {}
    return http.get("/plates/bookmarked", params)
  }
  getFilterOptions(): Promise<PlateFilterOptions> {
    return http.get("/plates/filters")
  }
  getBySlug(slug: string): Promise<Plate> {
    return http.get(`/plates/${slug}`)
  }
  submitRepository(input: SubmitRepositoryInput): Promise<Plate> {
    return http.post("/plates/repository", input)
  }
  moveToOrganization(id: string, organizationId?: string): Promise<Plate> {
    return http.put(`/plates/${id}/organization`, {
      organization_id: organizationId && organizationId.trim() !== "" ? organizationId : null,
    })
  }
  remove(id: string): Promise<void> {
    return http.delete(`/plates/${id}/remove`)
  }
  verifyRepository(id: string): Promise<Plate> {
    return http.post(`/plates/${id}/verify`, {})
  }
  setBookmark(id: string, bookmarked: boolean): Promise<void> {
    return http.post(`/plates/${id}/bookmark`, { bookmarked })
  }
  async rate(id: string, input: RatePlateInput): Promise<void> {
    try {
      await http.post(`/plates/${id}/reviews`, input)
      return
    } catch (error) {
      if (error instanceof ApiError && error.isNotFound()) {
        await http.post(`/plates/${id}/rate`, input)
        return
      }
      throw error
    }
  }
  getStats(): Promise<PlateStats> {
  return http.get("/plates/stats")
}
}


export const plateRepository = new PlateRepository()
