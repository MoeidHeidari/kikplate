"use client"

import { useState, useEffect } from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { usePlateFilterOptions, usePlates } from "@/src/presentation/hooks/usePlates"
import { PlateGrid } from "./PlateGrid"
import { PlateFilters } from "./PlateFilters"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"

const PAGE_SIZE = 24

interface Props {
  limit?: number
  initialSearch?: string
  initialTag?: string
  initialCategory?: string
}

export function PlateGridClient({
  initialSearch = "",
  initialTag = "",
  initialCategory = "",
}: Props) {
  const [search, setSearch]           = useState(initialSearch)
  const [categories, setCategories]   = useState<string[]>(initialCategory ? [initialCategory] : [])
  const [tags, setTags]               = useState<string[]>(initialTag ? [initialTag] : [])
  const [page, setPage]               = useState(1)
  const { data: filterOptions } = usePlateFilterOptions()

  useEffect(() => { setPage(1) }, [search, categories, tags])

  const { data, isLoading, isError } = usePlates({
    search,
    tags: tags.length > 0 ? tags : undefined,
    categories: categories.length > 0 ? categories : undefined,
    page,
    limit: PAGE_SIZE,
  })

  const total = data?.total ?? 0
  const totalPages = Math.ceil(total / PAGE_SIZE)

  return (
    <div className="flex gap-8">
      <aside className="hidden lg:block w-52 shrink-0">
        <PlateFilters
          search={search}
          onSearch={setSearch}
          activeCategories={categories}
          onCategories={setCategories}
          activeTags={tags}
          onTags={setTags}
          categories={filterOptions?.categories ?? []}
          tags={filterOptions?.tags ?? []}
        />
      </aside>

      <div className="flex-1 min-w-0">
        <div className="mb-4 flex items-center justify-between gap-4">
          <p className="text-xs text-muted-foreground">
            {isLoading ? "Loading…" : `${total} plate${total !== 1 ? "s" : ""} found`}
          </p>
          {totalPages > 1 && (
            <span className="text-xs text-muted-foreground">
              Page {page} of {totalPages}
            </span>
          )}
        </div>

        {isLoading && <LoadingSpinner />}
        {isError && <p className="text-sm text-destructive">Failed to load plates.</p>}
        {!isLoading && !isError && data?.data?.length === 0 && (
          <p className="py-20 text-center text-sm text-muted-foreground">No plates found.</p>
        )}
        {!isLoading && !isError && data?.data && data.data.length > 0 && <PlateGrid plates={data.data} />}

        <div className="mt-8 flex items-center justify-center gap-1">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="flex h-8 w-8 items-center justify-center border border-border text-muted-foreground transition-colors hover:border-foreground/40 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
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
                <span key={`ellipsis-${i}`} className="flex h-8 w-8 items-center justify-center text-xs text-muted-foreground">…</span>
              ) : (
                <button
                  key={p}
                  onClick={() => setPage(p as number)}
                  className={`flex h-8 w-8 items-center justify-center border text-xs transition-colors ${
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
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="flex h-8 w-8 items-center justify-center border border-border text-muted-foreground transition-colors hover:border-foreground/40 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
          >
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      </div>

    </div>
  )
}
