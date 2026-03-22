import { http } from "./httpClient"
import type { IBadgeRepository } from "@/src/domain/repositories/IBadgeRepository"
import type { Badge } from "@/src/domain/entities/Badge"

class BadgeRepository implements IBadgeRepository {
  list(): Promise<Badge[]> { return http.get("/badges") }
}

export const badgeRepository = new BadgeRepository()
