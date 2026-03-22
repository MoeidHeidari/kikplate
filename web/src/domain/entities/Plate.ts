export type PlateType       = "repository" | "file"
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

export interface Plate {
  id: string
  owner_id: string
  type: PlateType
  slug: string
  name: string
  description?: string
  category: string
  status: PlateStatus
  visibility: PlateVisibility
  metadata?: Record<string, unknown>
  use_count: number
  star_count: number
  avg_rating: number
  is_verified: boolean
  verified_at?: string
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
  content?: string
  filename?: string
  tags?: PlateTag[]
  badges?: PlateBadge[]
}

export interface PlateListResponse {
  data: Plate[]
  total: number
}

export interface PlateFilter {
  type?: PlateType
  category?: string
  tag?: string
  search?: string
  page?: number
  limit?: number
}

export interface SubmitFileInput {
  name: string
  description?: string
  category: string
  visibility: PlateVisibility
  filename: string
  content: string
  tags?: string[]
}

export interface SubmitRepositoryInput {
  repo_url: string
  branch?: string
}

export const isRepositoryPlate = (p: Plate) => p.type === "repository"
export const isFilePlate        = (p: Plate) => p.type === "file"
export const isApproved         = (p: Plate) => p.status === "approved"
export const isPending          = (p: Plate) => p.status === "pending"
export const isOwnedBy          = (p: Plate, accountId: string) => p.owner_id === accountId
