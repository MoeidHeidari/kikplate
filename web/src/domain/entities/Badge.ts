export type BadgeTier = "community" | "verified" | "official" | "sponsored"

export interface Badge {
  id: string
  slug: string
  name: string
  description: string
  icon: string
  tier: BadgeTier
  created_at: string
}
