import type {
  Plate,
  PlateFilter,
  PlateFilterOptions,
  PlateListResponse,
  RatePlateInput,
  SubmitRepositoryInput,
} from "@/src/domain/entities/Plate"
import { PlateStats } from "../entities/Stats"

export interface IPlateRepository {
  list(filter: PlateFilter): Promise<PlateListResponse>
  listBookmarked(limit?: number): Promise<PlateListResponse>
  getFilterOptions(): Promise<PlateFilterOptions>
  getBySlug(slug: string): Promise<Plate>
  submitRepository(input: SubmitRepositoryInput): Promise<Plate>
  moveToOrganization(id: string, organizationId?: string): Promise<Plate>
  remove(id: string): Promise<void>
  setBookmark(id: string, bookmarked: boolean): Promise<void>
  rate(id: string, input: RatePlateInput): Promise<void>
  getStats(): Promise<PlateStats>
}
