import { AuthService } from "@/src/domain/services/AuthService"

const BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001"

export class ApiError extends Error {
  constructor(public readonly status: number, message: string) {
    super(message)
    this.name = "ApiError"
  }
  isNotFound()      { return this.status === 404 }
  isUnauthorized()  { return this.status === 401 }
  isForbidden()     { return this.status === 403 }
  isConflict()      { return this.status === 409 }
  isUnprocessable() { return this.status === 422 }
}

function buildQuery(params: Record<string, unknown>): string {
  const q = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== null && v !== "") q.set(k, String(v))
  }
  const s = q.toString()
  return s ? `?${s}` : ""
}

async function httpFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = AuthService.getToken()
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers as Record<string, string> ?? {}),
  }
  if (token) headers["Authorization"] = `Bearer ${token}`

  const res = await fetch(`${BASE}${path}`, { ...init, headers })

  if (!res.ok) {
    let message = res.statusText
    try { const d = await res.json(); message = d.error ?? message } catch { /* noop */ }
    throw new ApiError(res.status, message)
  }

  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

export const http = {
  get<T>(path: string, params: Record<string, unknown> = {}) {
    return httpFetch<T>(`${path}${buildQuery(params)}`)
  },
  post<T>(path: string, body?: unknown) {
    return httpFetch<T>(path, {
      method: "POST",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    })
  },
  patch<T>(path: string, body?: unknown) {
    return httpFetch<T>(path, {
      method: "PATCH",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    })
  },
  delete<T = void>(path: string) {
    return httpFetch<T>(path, { method: "DELETE" })
  },
}
