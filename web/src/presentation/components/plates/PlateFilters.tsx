"use client"

import { Search } from "lucide-react"

const CATEGORIES = [
  "backend",
  "frontend",
  "fullstack",
  "devops",
  "mobile",
  "cli",
  "library",
  "other",
]

interface Props {
  search: string
  onSearch: (v: string) => void
  activeTag?: string
  onTag: (tag: string | undefined) => void
  activeCategory?: string
  onCategory: (cat: string | undefined) => void
}

export function PlateFilters({
  search,
  onSearch,
  activeCategory,
  onCategory,
}: Props) {
  return (
    <div className="space-y-6">

      {/* search */}
      <div className="flex items-center border border-border px-3 gap-2 focus-within:border-ring transition-colors">
        <Search className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
        <input
          className="h-9 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          placeholder="Filter plates..."
          value={search}
          onChange={(e) => onSearch(e.target.value)}
        />
      </div>

      {/* categories */}
      <div>
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
          Category
        </p>
        <div className="space-y-0.5">
          <button
            onClick={() => onCategory(undefined)}
            className={`w-full text-left px-2 py-1.5 text-sm transition-colors ${
              !activeCategory
                ? "bg-primary text-primary-foreground"
                : "text-muted-foreground hover:text-foreground hover:bg-muted"
            }`}
          >
            All
          </button>
          {CATEGORIES.map((cat) => (
            <button
              key={cat}
              onClick={() => onCategory(activeCategory === cat ? undefined : cat)}
              className={`w-full text-left px-2 py-1.5 text-sm capitalize transition-colors ${
                activeCategory === cat
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground hover:text-foreground hover:bg-muted"
              }`}
            >
              {cat}
            </button>
          ))}
        </div>
      </div>

      {/* type */}
      <div>
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
          Type
        </p>
        <div className="space-y-0.5">
          {["repository", "file"].map((type) => (
            <button
              key={type}
              className="w-full text-left px-2 py-1.5 text-sm capitalize text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
            >
              {type}
            </button>
          ))}
        </div>
      </div>

    </div>
  )
}