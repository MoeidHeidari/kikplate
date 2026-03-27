export type PlateType       = "repository"
export type PlateStatus     = "pending" | "approved" | "rejected" | "archived"
export type PlateVisibility = "public" | "private" | "unlisted"
export type SyncStatus      = "pending" | "syncing" | "synced" | "failed" | "unverified"

export interface PlateTag {
  id: string
  plate_id: string
  tag: string
}

export interface PlateBadge {
  id: string
  plate_id: string
  badge_id: string
  granted_by: string
  reason?: string
  granted_at: string
  badge?: import("./Badge").Badge
}

export interface PlateOwner {
  id: string
  username?: string
  display_name?: string
  avatar_url?: string
}

export interface PlateOrganization {
  id: string
  name: string
  description?: string
  logo_url?: string
  owner_id: string
  owner?: {
    id: string
    username?: string
    display_name?: string
    avatar_url?: string
  }
}

export interface Plate {
  id: string
  owner_id: string
  organization_id?: string
  type: PlateType
  slug: string
  name: string
  description?: string
  category: string
  status: PlateStatus
  visibility: PlateVisibility
  metadata?: Record<string, unknown>
  bookmark_count: number
  star_count: number
  avg_rating: number
  user_rating?: number
  is_verified: boolean
  verified_at?: string
  verification_token?: string
  verification_token_set_at?: string
  published_at?: string
  created_at: string
  updated_at: string
  repo_url?: string
  branch?: string
  sync_status?: SyncStatus
  sync_error?: string
  sync_interval?: string
  next_sync_at?: string
  last_synced_at?: string
  consecutive_failures?: number
  is_bookmarked?: boolean
  tags?: PlateTag[]
  badges?: PlateBadge[]
  owner?: PlateOwner
  organization?: PlateOrganization
}

export interface PlateListResponse {
  data: Plate[]
  total: number
}

export interface PlateFilterOptions {
  categories: string[]
  tags: string[]
}

export interface PlateFilter {
  type?: PlateType | PlateType[]
  types?: PlateType[]
  category?: string
  categories?: string[]
  tag?: string
  tags?: string[]
  search?: string
  owner_id?: string
  page?: number
  limit?: number
}

export interface SubmitRepositoryInput {
  repo_url: string
  branch?: string
  organization_id?: string
}

export interface RatePlateInput {
  rating: number
}

export const isRepositoryPlate = (p: Plate) => p.type === "repository"
export const isApproved         = (p: Plate) => p.status === "approved"
export const isPending          = (p: Plate) => p.status === "pending"
export const isOwnedBy          = (p: Plate, accountId: string) => p.owner_id === accountId
