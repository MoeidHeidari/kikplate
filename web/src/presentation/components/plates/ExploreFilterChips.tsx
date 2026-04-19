"use client"

import { X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { formatExplorerCategoryLabel } from "@/src/presentation/utils/exploreLabels"

interface Props {
  search: string
  onClearSearch: () => void
  categories: string[]
  onRemoveCategory: (c: string) => void
  tags: string[]
  onRemoveTag: (t: string) => void
  badgeSlugs: string[]
  badgeLabel: (slug: string) => string
  onRemoveBadge: (slug: string) => void
  onClearAll: () => void
}

function Chip({
  label,
  onRemove,
  preserveCase,
}: {
  label: string
  onRemove: () => void
  preserveCase?: boolean
}) {
  return (
    <span className="inline-flex max-w-full items-center gap-1 border border-border bg-muted/50 pl-2.5 pr-1 py-1 text-xs text-foreground">
      <span className={cn("truncate", !preserveCase && "capitalize")}>{label}</span>
      <button
        type="button"
        onClick={onRemove}
        className="shrink-0 rounded-none p-0.5 text-muted-foreground hover:bg-muted hover:text-foreground"
        aria-label={`Remove ${label}`}
      >
        <X className="h-3 w-3" />
      </button>
    </span>
  )
}

export function ExploreFilterChips({
  search,
  onClearSearch,
  categories,
  onRemoveCategory,
  tags,
  onRemoveTag,
  badgeSlugs,
  badgeLabel,
  onRemoveBadge,
  onClearAll,
}: Props) {
  const hasSearch = search.trim().length > 0
  const hasListFilters =
    categories.length > 0 || tags.length > 0 || badgeSlugs.length > 0
  const hasAny = hasSearch || hasListFilters

  if (!hasAny) return null

  return (
    <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4">
      <div className="flex min-w-0 flex-wrap items-center gap-2">
        <span className="text-xs font-medium uppercase tracking-wider text-muted-foreground shrink-0">
          Active
        </span>
        {hasSearch && (
          <Chip
            preserveCase
            label={`"${search.trim().slice(0, 48)}${search.trim().length > 48 ? "…" : ""}"`}
            onRemove={onClearSearch}
          />
        )}
        {categories.map((c) => (
          <Chip key={`cat-${c}`} label={formatExplorerCategoryLabel(c)} onRemove={() => onRemoveCategory(c)} />
        ))}
        {tags.map((t) => (
          <Chip key={`tag-${t}`} label={t} onRemove={() => onRemoveTag(t)} />
        ))}
        {badgeSlugs.map((slug) => (
          <Chip key={`badge-${slug}`} label={badgeLabel(slug)} onRemove={() => onRemoveBadge(slug)} />
        ))}
      </div>
      <Button type="button" variant="ghost" size="xs" className="shrink-0 text-muted-foreground hover:text-foreground" onClick={onClearAll}>
        Clear all
      </Button>
    </div>
  )
}
