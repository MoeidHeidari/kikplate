import { http } from "./httpClient"
import type { IAuthRepository } from "@/src/domain/repositories/IAuthRepository"
import type { AuthResult, LoginInput, MeResult, RegisterInput } from "@/src/domain/entities/User"

const BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001"

class AuthRepository implements IAuthRepository {
  register(input: RegisterInput): Promise<{ message: string }> {
    return http.post("/auth/register", input)
  }
  login(input: LoginInput): Promise<AuthResult> {
    return http.post("/auth/login", input)
  }
  verifyEmail(token: string): Promise<AuthResult> {
    return http.get("/auth/verify-email", { token })
  }
  me(): Promise<MeResult> {
    return http.get("/me")
  }
  updateProfile(input: { display_name?: string; avatar_url?: string }): Promise<MeResult> {
    return http.patch("/me/profile", input)
  }
  setUsername(username: string): Promise<{ message: string }> {
    return http.patch("/me/username", { username })
  }
  oauthRedirectURL(provider: string): string {
    return `${BASE}/auth/${provider}/redirect`
  }
  providers(): Promise<{ providers: string[] }> {
  return http.get("/auth/providers")
  }
}

export const authRepository = new AuthRepository()
