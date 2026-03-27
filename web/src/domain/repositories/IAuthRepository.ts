import type { AuthResult, LoginInput, MeResult, RegisterInput } from "@/src/domain/entities/User"

export interface IAuthRepository {
  register(input: RegisterInput): Promise<{ message: string }>
  login(input: LoginInput): Promise<AuthResult>
  verifyEmail(token: string): Promise<AuthResult>
  me(): Promise<MeResult>
  updateProfile(input: { display_name?: string; avatar_url?: string }): Promise<MeResult>
  setUsername(username: string): Promise<{ message: string }>
  oauthRedirectURL(provider: string): string
  providers(): Promise<{ providers: string[] }>
}
