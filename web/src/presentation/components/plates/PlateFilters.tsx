"use client"

import { useState } from "react"
import { Search, X } from "lucide-react"
import type {
  PlateBadgeFilterOption,
  PlateCategoryFilterOption,
  PlateTagFilterOption,
} from "@/src/domain/entities/Plate"
import { formatExplorerCategoryLabel } from "@/src/presentation/utils/exploreLabels"

interface Props {
  search: string
  onSearch: (v: string) => void
  onClearAll: () => void
  activeCategories: string[]
  onCategories: (cats: string[]) => void
  activeTags: string[]
  onTags: (tags: string[]) => void
  activeBadges: string[]
  onBadges: (slugs: string[]) => void
  categories: PlateCategoryFilterOption[]
  tags: PlateTagFilterOption[]
  badges: PlateBadgeFilterOption[]
}

function Checkbox({
  checked,
  onChange,
  label,
  count,
  disabled,
}: {
  checked: boolean
  onChange: () => void
  label: string
  count?: number
  disabled?: boolean
}) {
  const muted = Boolean(disabled)
  return (
    <button
      type="button"
      onClick={onChange}
      disabled={disabled}
      className={`group flex w-full items-center gap-2 px-1 py-1 text-sm text-left ${
        muted ? "cursor-not-allowed opacity-50" : ""
      }`}
      aria-pressed={checked}
    >
      <span
        className={`flex h-3.5 w-3.5 shrink-0 items-center justify-center border transition-colors ${
          checked ? "border-primary bg-primary" : "border-border bg-background group-hover:border-foreground/40"
        } ${muted ? "border-border/60" : ""}`}
      >
        {checked && (
          <svg viewBox="0 0 10 10" className="h-2.5 w-2.5 text-primary-foreground" fill="none" stroke="currentColor" strokeWidth="1.5">
            <polyline points="1.5,5 4,7.5 8.5,2.5" />
          </svg>
        )}
      </span>
      <span
        className={`min-w-0 flex-1 truncate ${
          checked ? "font-medium text-foreground" : muted ? "text-muted-foreground" : "text-muted-foreground group-hover:text-foreground"
        }`}
      >
        {label}
      </span>
      {count !== undefined && (
        <span
          className={`shrink-0 tabular-nums text-xs ${
            muted ? "text-muted-foreground/70" : "text-muted-foreground"
          }`}
        >
          {count}
        </span>
      )}
    </button>
  )
}

function BadgesDropdown({
  activeBadges,
  onBadges,
  searchTerm,
  setSearchTerm,
  badges,
}: {
  activeBadges: string[]
  onBadges: (slugs: string[]) => void
  searchTerm: string
  setSearchTerm: (s: string) => void
  badges: PlateBadgeFilterOption[]
}) {
  const filtered = badges.filter((b) =>
    b.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    b.slug.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="space-y-1.5">
      <div className="flex items-center border border-border px-2 gap-1.5 focus-within:border-ring transition-colors">
        <Search className="h-3 w-3 text-muted-foreground shrink-0" />
        <input
          type="text"
          placeholder="Search badges..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="h-7 w-full bg-transparent text-xs outline-none placeholder:text-muted-foreground"
        />
      </div>
      <div className="space-y-0">
        {filtered.length === 0 ? (
          <p className="text-xs text-muted-foreground px-1 py-1.5">No badges found</p>
        ) : (
          filtered.map((b) => {
            const selected = activeBadges.includes(b.slug)
            const noPlates = b.count === 0
            return (
              <Checkbox
                key={b.slug}
                checked={selected}
                disabled={noPlates && !selected}
                count={b.count}
                onChange={() => {
                  if (selected) {
                    onBadges(activeBadges.filter((s) => s !== b.slug))
                  } else {
                    onBadges([...activeBadges, b.slug])
                  }
                }}
                label={b.name}
              />
            )
          })
        )}
      </div>
    </div>
  )
}

function CategoriesDropdown({
  activeCategories,
  onCategories,
  searchTerm,
  setSearchTerm,
  categories,
}: {
  activeCategories: string[]
  onCategories: (cats: string[]) => void
  searchTerm: string
  setSearchTerm: (s: string) => void
  categories: PlateCategoryFilterOption[]
}) {
  const q = searchTerm.toLowerCase()
  const filtered = categories.filter((row) => {
    const label = formatExplorerCategoryLabel(row.slug)
    return (
      row.slug.toLowerCase().includes(q) ||
      label.toLowerCase().includes(q)
    )
  })

  function toggleCategory(slug: string) {
    if (activeCategories.includes(slug)) {
      onCategories(activeCategories.filter((c) => c !== slug))
    } else {
      onCategories([...activeCategories, slug])
    }
  }

  return (
    <div className="space-y-1.5">
      <div className="flex items-center border border-border px-2 gap-1.5 focus-within:border-ring transition-colors">
        <Search className="h-3 w-3 text-muted-foreground shrink-0" />
        <input
          type="text"
          placeholder="Search categories..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="h-7 w-full bg-transparent text-xs outline-none placeholder:text-muted-foreground"
        />
      </div>
      <div className="space-y-0">
        {categories.length === 0 ? (
          <p className="text-xs text-muted-foreground px-1 py-1.5">No categories available.</p>
        ) : filtered.length === 0 ? (
          <p className="text-xs text-muted-foreground px-1 py-1.5">No categories found</p>
        ) : (
          filtered.map((row) => {
            const selected = activeCategories.includes(row.slug)
            const noPlates = row.count === 0
            return (
              <Checkbox
                key={row.slug}
                checked={selected}
                disabled={noPlates && !selected}
                count={row.count}
                onChange={() => toggleCategory(row.slug)}
                label={formatExplorerCategoryLabel(row.slug)}
              />
            )
          })
        )}
      </div>
    </div>
  )
}

function TagsDropdown({
  activeTags,
  onTags,
  searchTerm,
  setSearchTerm,
  tags,
}: {
  activeTags: string[]
  onTags: (tags: string[]) => void
  searchTerm: string
  setSearchTerm: (s: string) => void
  tags: PlateTagFilterOption[]
}) {
  const filtered = tags.filter((row) => row.tag.toLowerCase().includes(searchTerm.toLowerCase()))

  return (
    <div className="space-y-1.5">
      <div className="flex items-center border border-border px-2 gap-1.5 focus-within:border-ring transition-colors">
        <Search className="h-3 w-3 text-muted-foreground shrink-0" />
        <input
          type="text"
          placeholder="Search tags..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="h-7 w-full bg-transparent text-xs outline-none placeholder:text-muted-foreground"
        />
      </div>
      <div className="space-y-0">
        {filtered.length === 0 ? (
          <p className="text-xs text-muted-foreground px-1 py-1.5">No tags found</p>
        ) : (
          filtered.map((row) => {
            const selected = activeTags.includes(row.tag)
            const noPlates = row.count === 0
            return (
              <Checkbox
                key={row.tag}
                checked={selected}
                disabled={noPlates && !selected}
                count={row.count}
                onChange={() => {
                  if (selected) {
                    onTags(activeTags.filter((at) => at !== row.tag))
                  } else {
                    onTags([...activeTags, row.tag])
                  }
                }}
                label={row.tag}
              />
            )
          })
        )}
      </div>
    </div>
  )
}

export function PlateFilters({
  search,
  onSearch,
  onClearAll,
  activeCategories,
  onCategories,
  activeTags,
  onTags,
  activeBadges,
  onBadges,
  categories,
  tags,
  badges,
}: Props) {
  const [tagSearch, setTagSearch] = useState("")
  const [badgeSearch, setBadgeSearch] = useState("")
  const [categorySearch, setCategorySearch] = useState("")

  const hasActiveFilters =
    search.trim() !== "" ||
    activeCategories.length > 0 ||
    activeTags.length > 0 ||
    activeBadges.length > 0

  return (
    <div className="space-y-5">
      <div className="flex items-center border border-border px-3 gap-2 focus-within:border-ring transition-colors">
        <Search className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
        <input
          className="h-9 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          placeholder="Search plates..."
          value={search}
          onChange={(e) => onSearch(e.target.value)}
        />
        {search && (
          <button onClick={() => onSearch("")} className="text-muted-foreground hover:text-foreground">
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      {hasActiveFilters && (
        <button
          type="button"
          onClick={onClearAll}
          className="w-full text-left text-xs text-muted-foreground hover:text-foreground flex items-center gap-1"
        >
          <X className="h-3 w-3" /> Clear all filters
        </button>
      )}

      <div>
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Badges</p>
        {badges.length === 0 ? (
          <p className="text-xs text-muted-foreground px-1 py-1.5">No badges available.</p>
        ) : (
          <BadgesDropdown
            activeBadges={activeBadges}
            onBadges={onBadges}
            searchTerm={badgeSearch}
            setSearchTerm={setBadgeSearch}
            badges={badges}
          />
        )}
      </div>

      <div>
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Tags</p>
        <TagsDropdown
          activeTags={activeTags}
          onTags={onTags}
          searchTerm={tagSearch}
          setSearchTerm={setTagSearch}
          tags={tags}
        />
      </div>

      <div>
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Category</p>
        <CategoriesDropdown
          activeCategories={activeCategories}
          onCategories={onCategories}
          searchTerm={categorySearch}
          setSearchTerm={setCategorySearch}
          categories={categories}
        />
      </div>
    </div>
  )
}

