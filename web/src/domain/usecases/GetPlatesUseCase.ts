import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"
import type { PlateFilter } from "@/src/domain/entities/Plate"

export class GetPlatesUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(filter: PlateFilter = {}) { return this.repo.list(filter) }
}
