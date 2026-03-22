import type { IAuthRepository } from "@/src/domain/repositories/IAuthRepository"

export class GetMeUseCase {
  constructor(private readonly repo: IAuthRepository) {}
  execute() { return this.repo.me() }
}
