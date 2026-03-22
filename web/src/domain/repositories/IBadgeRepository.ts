import type { Badge } from "@/src/domain/entities/Badge"

export interface IBadgeRepository {
  list(): Promise<Badge[]>
}
