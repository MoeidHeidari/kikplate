import type { IPlateRepository } from "@/src/domain/repositories/IPlateRepository"

export class SetBookmarkUseCase {
  constructor(private readonly repo: IPlateRepository) {}
  execute(id: string, bookmarked: boolean) { return this.repo.setBookmark(id, bookmarked) }
}
