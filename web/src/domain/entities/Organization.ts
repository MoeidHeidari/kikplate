export interface Organization {
  id: string
  name: string
  visibility: "public" | "private"
  description: string
  logo_url?: string
  owner_id: string
  github_connected?: boolean
  github_installation_account?: string
  membership_role?: "owner" | "admin" | "member"
  created_at: string
  updated_at: string
  owner?: {
    id: string
    username?: string
    display_name?: string
    avatar_url?: string
  }
}

export interface CreateOrganizationInput {
  name: string
  description: string
  logo_url?: string
  visibility?: "public" | "private"
}

export interface UpdateOrganizationInput {
  name?: string
  description?: string
  logo_url?: string
  visibility?: "public" | "private"
}

export interface OrganizationInvitation {
  id: string
  organization_id: string
  invited_by: string
  email: string
  role: "admin" | "member"
  status: "pending" | "accepted" | "declined"
  expires_at: string
  created_at: string
  updated_at: string
}

export interface OrganizationMember {
  id: string
  account_id: string
  role: "owner" | "admin" | "member"
  status: "pending" | "accepted" | "declined"
  joined_at: string
  profile?: {
    id: string
    username?: string
    display_name?: string
    avatar_url?: string
  }
}

export interface OrganizationListResponse {
  organizations: Organization[]
  total: number
  limit: number
  offset: number
}
