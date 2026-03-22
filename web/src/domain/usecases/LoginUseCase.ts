import type { IAuthRepository } from "@/src/domain/repositories/IAuthRepository"
import type { LoginInput } from "@/src/domain/entities/User"
import { AuthService } from "@/src/domain/services/AuthService"

export class LoginUseCase {
  constructor(private readonly repo: IAuthRepository) {}
  async execute(input: LoginInput) {
    const result = await this.repo.login(input)
    AuthService.setToken(result.token)
    return result
  }
}
