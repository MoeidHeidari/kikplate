import type { IBadgeRepository } from "@/src/domain/repositories/IBadgeRepository"

export class GetBadgesUseCase {
  constructor(private readonly repo: IBadgeRepository) {}
  execute() { return this.repo.list() }
}
