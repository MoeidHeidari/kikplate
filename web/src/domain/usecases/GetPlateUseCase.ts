import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"

export class GetPlateUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(slug: string) { return this.repo.getBySlug(slug) }
}
