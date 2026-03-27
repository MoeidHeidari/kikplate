export interface Organization {
  id: string
  name: string
  description: string
  logo_url?: string
  owner_id: string
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
}

export interface UpdateOrganizationInput {
  name?: string
  description?: string
  logo_url?: string
}

export interface OrganizationListResponse {
  organizations: Organization[]
  total: number
  limit: number
  offset: number
}
