const TOKEN_KEY = "kp_token"

function setTokenCookie(token: string) {
  if (typeof document === "undefined") return
  document.cookie = `${TOKEN_KEY}=${encodeURIComponent(token)}; Path=/; Max-Age=2592000; SameSite=Lax`
}

function clearTokenCookie() {
  if (typeof document === "undefined") return
  document.cookie = `${TOKEN_KEY}=; Path=/; Max-Age=0; SameSite=Lax`
}

export const AuthService = {
  getToken(): string | null {
    if (typeof window === "undefined") return null
    return localStorage.getItem(TOKEN_KEY)
  },
  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
    setTokenCookie(token)
  },
  clearToken(): void {
    localStorage.removeItem(TOKEN_KEY)
    clearTokenCookie()
  },
  isAuthenticated(): boolean {
    return AuthService.getToken() !== null
  },
}
