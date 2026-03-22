"use client"

import { useState } from "react"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { PlateGrid } from "./PlateGrid"
import { PlateFilters } from "./PlateFilters"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"

interface Props {
  limit?: number
  initialSearch?: string
  initialTag?: string
  initialCategory?: string
}

export function PlateGridClient({
  limit = 20,
  initialSearch = "",
  initialTag = "",
  initialCategory = "",
}: Props) {
  const [search, setSearch]     = useState(initialSearch)
  const [tag, setTag]           = useState<string | undefined>(initialTag || undefined)
  const [category, setCategory] = useState<string | undefined>(initialCategory || undefined)

  const { data, isLoading, isError } = usePlates({ search, tag, category, limit })

  return (
    <div className="flex gap-8">

      {/* left sidebar */}
      <aside className="hidden lg:block w-48 shrink-0">
        <PlateFilters
          search={search}
          onSearch={setSearch}
          activeTag={tag}
          onTag={setTag}
          activeCategory={category}
          onCategory={setCategory}
        />
      </aside>

      {/* grid */}
      <div className="flex-1 min-w-0">
        {isLoading && <LoadingSpinner />}
        {isError && (
          <p className="text-sm text-destructive">Failed to load plates.</p>
        )}
        {data && (
          <>
            <p className="text-xs text-muted-foreground mb-4">
              {data.total} plate{data.total !== 1 ? "s" : ""} found
            </p>
            <PlateGrid plates={data.data} />
          </>
        )}
      </div>

    </div>
  )
}