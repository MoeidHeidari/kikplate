import type { IAuthRepository } from "@/src/domain/repositories/IAuthRepository"
import { AuthService } from "@/src/domain/services/AuthService"

export class VerifyEmailUseCase {
  constructor(private readonly repo: IAuthRepository) {}
  async execute(token: string) {
    const result = await this.repo.verifyEmail(token)
    AuthService.setToken(result.token)
    return result
  }
}
