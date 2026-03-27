import type { IConfigRepository } from "@/src/domain/repositories/IConfigRepository"

export class GetConfigUseCase {
  constructor(private readonly repo: IConfigRepository) {}
  execute() { return this.repo.getConfig() }
}
