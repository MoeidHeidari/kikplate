import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"
import type { SubmitFileInput, SubmitRepositoryInput } from "@/src/domain/entities/Plate"

export class SubmitFileUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(input: SubmitFileInput) { return this.repo.submitFile(input) }
}

export class SubmitRepositoryUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(input: SubmitRepositoryInput) { return this.repo.submitRepository(input) }
}
