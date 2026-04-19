"use client"

import { useState, useEffect, useMemo, useCallback } from "react"
import { ChevronLeft, ChevronRight, ListFilter } from "lucide-react"
import { useBadges } from "@/src/presentation/hooks/useBadges"
import { usePlateFilterOptions, usePlates } from "@/src/presentation/hooks/usePlates"
import { PlateGrid } from "./PlateGrid"
import { PlateFilters } from "./PlateFilters"
import { ExploreFilterChips } from "./ExploreFilterChips"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"

const PAGE_SIZE = 24

interface Props {
  limit?: number
  initialSearch?: string
  initialTag?: string
  initialCategory?: string
  initialBadge?: string
}

function parseInitialListParam(raw: string): string[] {
  if (!raw.trim()) return []
  return raw.split(",").map((s) => s.trim()).filter(Boolean)
}

export function PlateGridClient({
  initialSearch = "",
  initialTag = "",
  initialCategory = "",
  initialBadge = "",
}: Props) {
  const [search, setSearch] = useState(initialSearch)
  const [categories, setCategories] = useState<string[]>(initialCategory ? [initialCategory] : [])
  const [tags, setTags] = useState<string[]>(initialTag ? [initialTag] : [])
  const [badges, setBadges] = useState<string[]>(() => parseInitialListParam(initialBadge))
  const [page, setPage] = useState(1)
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false)

  const { data: filterOptions } = usePlateFilterOptions()
  const { data: catalogBadges } = useBadges()

  const badgeFilterOptions = useMemo(() => {
    if (filterOptions?.badges?.length) {
      return filterOptions.badges
    }
    return (catalogBadges ?? []).map((b) => ({
      slug: b.slug,
      name: b.name,
      count: 0,
    }))
  }, [filterOptions?.badges?.length, catalogBadges?.length])

  const badgeLabel = useCallback(
    (slug: string) => badgeFilterOptions.find((b) => b.slug === slug)?.name ?? slug,
    [badgeFilterOptions.length]
  )

  useEffect(() => {
    setPage(1)
  }, [search, categories, tags, badges])

  const { data, isLoading, isError } = usePlates({
    search,
    tags: tags.length > 0 ? tags : undefined,
    categories: categories.length > 0 ? categories : undefined,
    badges: badges.length > 0 ? badges : undefined,
    page,
    limit: PAGE_SIZE,
  })

  const total = data?.total ?? 0
  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  const clearAllFilters = useCallback(() => {
    setSearch("")
    setCategories([])
    setTags([])
    setBadges([])
  }, [])

  const filterProps = {
    search,
    onSearch: setSearch,
    onClearAll: clearAllFilters,
    activeCategories: categories,
    onCategories: setCategories,
    activeTags: tags,
    onTags: setTags,
    activeBadges: badges,
    onBadges: setBadges,
    categories: filterOptions?.categories ?? [],
    tags: filterOptions?.tags ?? [],
    badges: badgeFilterOptions,
  }

  const activeFilterCount =
    (search.trim() ? 1 : 0) + categories.length + tags.length + badges.length

  const hasActiveFilters = activeFilterCount > 0

  const showPagination = total > 0 && totalPages > 1
  const showEmpty = !isLoading && !isError && data?.data?.length === 0

  return (
    <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:gap-10">
      <div className="lg:hidden">
        <Button
          type="button"
          variant="outline"
          className="h-11 w-full justify-between gap-3 px-4"
          onClick={() => setMobileFiltersOpen(true)}
        >
          <span className="flex items-center gap-2 font-medium">
            <ListFilter className="h-4 w-4 shrink-0 text-muted-foreground" />
            Filters
          </span>
          {activeFilterCount > 0 ? (
            <span className="flex h-6 min-w-6 items-center justify-center border border-primary bg-primary px-2 text-xs font-semibold text-primary-foreground">
              {activeFilterCount}
            </span>
          ) : (
            <span className="text-xs text-muted-foreground">Refine results</span>
          )}
        </Button>

        <Sheet open={mobileFiltersOpen} onOpenChange={setMobileFiltersOpen}>
          <SheetContent
            side="left"
            className="flex h-full max-h-dvh w-[min(100%,22rem)] flex-col gap-0 bg-card p-0 text-card-foreground sm:max-w-md"
          >
            <SheetHeader className="shrink-0 border-b border-border px-4 py-4 text-left">
              <SheetTitle className="font-heading text-lg">Filters</SheetTitle>
              <SheetDescription className="text-xs leading-relaxed">
                Search and narrow plates by badge, tag, or category.
              </SheetDescription>
            </SheetHeader>
            <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-4 py-4">
              <PlateFilters {...filterProps} />
            </div>
          </SheetContent>
        </Sheet>
      </div>

      <aside className="hidden w-70 shrink-0 lg:sticky lg:top-24 lg:block lg:max-h-[calc(100dvh-7rem)] lg:overflow-y-auto lg:pr-1 lg:pt-0.5">
        <div className="border border-border bg-card p-4">
          <p className="mb-4 font-heading text-sm font-semibold text-foreground">Filter plates</p>
          <PlateFilters {...filterProps} />
        </div>
      </aside>

      <div className="min-w-0 flex-1 space-y-4">
        <ExploreFilterChips
          search={search}
          onClearSearch={() => setSearch("")}
          categories={categories}
          onRemoveCategory={(c) => setCategories((prev) => prev.filter((x) => x !== c))}
          tags={tags}
          onRemoveTag={(t) => setTags((prev) => prev.filter((x) => x !== t))}
          badgeSlugs={badges}
          badgeLabel={badgeLabel}
          onRemoveBadge={(slug) => setBadges((prev) => prev.filter((x) => x !== slug))}
          onClearAll={clearAllFilters}
        />

        <div className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 border-b border-border pb-3">
          <p className="text-sm text-muted-foreground">
            {isLoading ? (
              <span className="inline-flex items-center gap-2">
                <span className="inline-block h-3.5 w-28 animate-pulse rounded-sm bg-muted" />
              </span>
            ) : (
              <>
                <span className="font-medium tabular-nums text-foreground">{total}</span>
                {" "}
                {total === 1 ? "plate" : "plates"}
                {hasActiveFilters ? " match" : " found"}
              </>
            )}
          </p>
          {showPagination && (
            <span className="text-xs text-muted-foreground tabular-nums">
              Page {page} of {totalPages}
            </span>
          )}
        </div>

        {isLoading && (
          <div className="flex min-h-[240px] items-center justify-center py-12">
            <LoadingSpinner />
          </div>
        )}

        {isError && (
          <p className="rounded-none border border-destructive/40 bg-destructive/5 px-4 py-3 text-sm text-destructive">
            Something went wrong while loading plates. Try again in a moment.
          </p>
        )}

        {showEmpty && (
          <div className="rounded-none border border-dashed border-border bg-muted/10 px-6 py-16 text-center">
            <p className="font-medium text-foreground">No plates match these criteria</p>
            <p className="mt-2 text-sm text-muted-foreground">
              {hasActiveFilters
                ? "Remove a filter or broaden your search."
                : "Check back later for new templates."}
            </p>
            {hasActiveFilters && (
              <Button type="button" variant="outline" size="sm" className="mt-6" onClick={clearAllFilters}>
                Clear all filters
              </Button>
            )}
          </div>
        )}

        {!isLoading && !isError && data?.data && data.data.length > 0 && <PlateGrid plates={data.data} />}

        {showPagination && (
          <div className="flex items-center justify-center gap-1 pt-4">
            <button
              type="button"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="flex h-9 w-9 items-center justify-center border border-border text-muted-foreground transition-colors hover:border-foreground/40 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
              aria-label="Previous page"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>

            {Array.from({ length: totalPages }, (_, i) => i + 1)
              .filter((p) => p === 1 || p === totalPages || Math.abs(p - page) <= 2)
              .reduce<(number | "…")[]>((acc, p, idx, arr) => {
                if (idx > 0 && p - (arr[idx - 1] as number) > 1) acc.push("…")
                acc.push(p)
                return acc
              }, [])
              .map((p, i) =>
                p === "…" ? (
                  <span
                    key={`ellipsis-${i}`}
                    className="flex h-9 w-9 items-center justify-center text-xs text-muted-foreground"
                  >
                    …
                  </span>
                ) : (
                  <button
                    key={p}
                    type="button"
                    onClick={() => setPage(p as number)}
                    className={`flex h-9 min-w-9 items-center justify-center border px-2 text-xs transition-colors ${
                      p === page
                        ? "border-foreground bg-foreground text-background"
                        : "border-border text-muted-foreground hover:border-foreground/40 hover:text-foreground"
                    }`}
                  >
                    {p}
                  </button>
                )
              )}

            <button
              type="button"
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="flex h-9 w-9 items-center justify-center border border-border text-muted-foreground transition-colors hover:border-foreground/40 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
              aria-label="Next page"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
