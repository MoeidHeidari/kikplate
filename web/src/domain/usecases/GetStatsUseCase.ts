import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"

export class GetStatsUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute() { return this.repo.getStats() }
}