import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"
import type { SubmitRepositoryInput } from "@/src/domain/entities/Plate"

export class SubmitRepositoryUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(input: SubmitRepositoryInput) { return this.repo.submitRepository(input) }
}
