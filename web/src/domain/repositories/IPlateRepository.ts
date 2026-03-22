import type {
  Plate,
  PlateFilter,
  PlateListResponse,
  SubmitFileInput,
  SubmitRepositoryInput,
} from "@/src/domain/entities/Plate"
import { PlateStats } from "../entities/Stats"

export interface IPlateRepository {
  list(filter: PlateFilter): Promise<PlateListResponse>
  getBySlug(slug: string): Promise<Plate>
  submitFile(input: SubmitFileInput): Promise<Plate>
  submitRepository(input: SubmitRepositoryInput): Promise<Plate>
  recordUse(id: string): Promise<void>
  getStats(): Promise<PlateStats>
}
