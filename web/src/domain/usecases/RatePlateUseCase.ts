import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"

export class RatePlateUseCase {
  constructor(private readonly repo: IPlateRepository) {}

  execute(id: string, rating: number) {
    return this.repo.rate(id, { rating })
  }
}
