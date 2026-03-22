const TOKEN_KEY = "kp_token"

export const AuthService = {
  getToken(): string | null {
    if (typeof window === "undefined") return null
    return localStorage.getItem(TOKEN_KEY)
  },
  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  },
  clearToken(): void {
    localStorage.removeItem(TOKEN_KEY)
  },
  isAuthenticated(): boolean {
    return AuthService.getToken() !== null
  },
}
