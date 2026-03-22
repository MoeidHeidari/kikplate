import type { Plate } from "@/src/domain/entities/Plate"
import type { BadgeTier } from "@/src/domain/entities/Badge"

export const PlateService = {
  syncStatusLabel(plate: Plate): string {
    const map: Record<string, string> = {
      synced: "Synced", syncing: "Syncing…",
      failed: "Sync failed", unverified: "Unverified", pending: "Pending sync",
    }
    return map[plate.sync_status ?? ""] ?? ""
  },

  badgeTierColour(tier: BadgeTier): string {
    const map: Record<BadgeTier, string> = {
      official:  "text-blue-600 border-blue-200 bg-blue-50",
      verified:  "text-green-600 border-green-200 bg-green-50",
      sponsored: "text-purple-600 border-purple-200 bg-purple-50",
      community: "text-gray-600 border-gray-200 bg-gray-50",
    }
    return map[tier]
  },
}
