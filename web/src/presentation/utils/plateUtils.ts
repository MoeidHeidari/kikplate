import type { PlateStatus } from "@/src/domain/entities/Plate"
import type { BadgeTier } from "@/src/domain/entities/Badge"

export function statusLabel(status: PlateStatus): string {
  const map: Record<PlateStatus, string> = {
    pending: "Pending", approved: "Approved",
    rejected: "Rejected", archived: "Archived",
  }
  return map[status] ?? status
}

export function statusVariant(status: PlateStatus): "default" | "secondary" | "destructive" | "outline" {
  switch (status) {
    case "approved": return "default"
    case "pending":  return "secondary"
    case "rejected": return "destructive"
    default:         return "outline"
  }
}

export function tierColour(tier: BadgeTier): string {
  const map: Record<BadgeTier, string> = {
    official:  "text-blue-600 border-blue-200 bg-blue-50",
    verified:  "text-green-600 border-green-200 bg-green-50",
    sponsored: "text-purple-600 border-purple-200 bg-purple-50",
    community: "text-gray-600 border-gray-200 bg-gray-50",
  }
  return map[tier]
}

export function formatCount(n: number): string {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return String(n)
}

export function buildPlateHref(slug: string): string {
  const encodedSlug = slug
    .split("/")
    .map((segment) => encodeURIComponent(segment))
    .join("/")
  return `/plates/${encodedSlug}`
}

export function relativeTime(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const m = Math.floor(diff / 60_000)
  if (m < 60)  return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24)  return `${h}h ago`
  const d = Math.floor(h / 24)
  if (d < 30)  return `${d}d ago`
  return `${Math.floor(d / 30)}mo ago`
}
