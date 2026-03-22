import type { AuthResult, LoginInput, MeResult, RegisterInput } from "@/src/domain/entities/User"

export interface IAuthRepository {
  register(input: RegisterInput): Promise<{ message: string }>
  login(input: LoginInput): Promise<AuthResult>
  verifyEmail(token: string): Promise<AuthResult>
  me(): Promise<MeResult>
  oauthRedirectURL(provider: string): string
  providers(): Promise<{ providers: string[] }>
}
