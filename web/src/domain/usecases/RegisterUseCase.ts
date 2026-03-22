import type { IAuthRepository } from "@/src/domain/repositories/IAuthRepository"
import type { RegisterInput } from "@/src/domain/entities/User"

export class RegisterUseCase {
  constructor(private readonly repo: IAuthRepository) {}
  execute(input: RegisterInput) { return this.repo.register(input) }
}
