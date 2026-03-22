import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"

export class RecordUseUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(id: string) { return this.repo.recordUse(id) }
}
