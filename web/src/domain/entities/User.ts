export type UserRole = "member" | "admin"

export interface Account {
  id: string
  user_id?: string
  provider: string
  provider_user_id: string
  display_name?: string
  avatar_url?: string
  created_at: string
}

export interface MeResult {
  account_id: string
  provider: string
  display_name?: string
  avatar_url?: string
  username?: string
  email?: string
  role?: UserRole
  is_active?: boolean
}

export interface AuthResult {
  token: string
  account: Account
}

export interface RegisterInput {
  username: string
  email: string
  password: string
}

export interface LoginInput {
  email: string
  password: string
}
