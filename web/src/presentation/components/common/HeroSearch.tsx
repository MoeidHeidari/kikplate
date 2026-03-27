"use client"

import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Search, GitBranch, Loader2, Github, Linkedin } from "lucide-react"
import { usePlates, useStats } from "@/src/presentation/hooks/usePlates"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import Link from "next/link"

const SAMPLE_QUERIES = [
  "Golang starter",
  "Clean architecture boilerplate for Nodejs",
  "Java spring-boot starter",
  "Python http server",
  "Next.js",
  "Gin framework",
  "Postgresql docker-compose",
  "Nginx",
]

export function HeroSearch() {
  const router = useRouter()
  const [query, setQuery] = useState("")
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const { data, isLoading } = usePlates({ search: query, limit: 6 })
  const { data: stats } = useStats()
  const { data: appConfig } = useConfig()

  const results = data?.data ?? []
  const showDropdown = open && query.trim().length > 1
  const titleLines = (appConfig?.banner_title ?? "The biggest library of\nstarter boilerplates").split("\n")
  const socialItems = (appConfig?.social_media ?? []).slice(0, 4)
  const sampleQueries = appConfig?.prepared_queries ?? SAMPLE_QUERIES

  const SlackIcon = () => (
    <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor" aria-hidden="true">
      <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z" />
    </svg>
  )

  function socialLabel(type: string) {
    const t = type.toLowerCase()
    if (t === "x") return "X"
    return t.charAt(0).toUpperCase() + t.slice(1)
  }

  function socialIcon(type: string) {
    const t = type.toLowerCase()
    if (t === "github") return <Github className="h-4 w-4" />
    if (t === "linkedin") return <Linkedin className="h-4 w-4" />
    if (t === "slack") return <SlackIcon />
    if (t === "x") return <span className="text-xs font-bold">X</span>
    return null
  }

  function handleSearch(value: string) {
    if (!value.trim()) return
    setOpen(false)
    router.push(`/explore?search=${encodeURIComponent(value.trim())}`)
  }

  function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Enter") handleSearch(query)
    if (e.key === "Escape") setOpen(false)
  }

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  return (
    <div className="dark relative min-h-[calc(100vh-3.5rem)] bg-background flex flex-col items-center justify-center gap-8 px-4 text-center">

      <h1 className="text-4xl sm:text-6xl font-bold tracking-tight text-foreground max-w-3xl leading-tight">
        {titleLines.map((line, idx) => (
          <span key={`${line}-${idx}`}>
            {line}
            {idx < titleLines.length - 1 && <br />}
          </span>
        ))}
      </h1>

      <div className="w-full max-w-2xl" ref={containerRef}>
        <div className="relative">
          <div className="flex items-center border border-border bg-card px-4 gap-3 focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/20 transition-all">
            <Search className="h-4 w-4 text-muted-foreground shrink-0" />
            <input
              className="h-12 w-full bg-transparent text-sm text-foreground outline-none placeholder:text-muted-foreground"
              placeholder="Search Plates..."
              value={query}
              onChange={(e) => {
                setQuery(e.target.value)
                setOpen(true)
              }}
              onKeyDown={handleKeyDown}
              onFocus={() => setOpen(true)}
            />
            {isLoading && query.trim().length > 1 && (
              <Loader2 className="h-4 w-4 text-muted-foreground animate-spin shrink-0" />
            )}
          </div>

          {showDropdown && (
            <div className="absolute top-full left-0 right-0 z-50 border border-border border-t-0 bg-card shadow-none">
              {isLoading ? (
                <div className="flex items-center justify-center py-6">
                  <Loader2 className="h-4 w-4 text-muted-foreground animate-spin" />
                </div>
              ) : results.length === 0 ? (
                <div className="px-4 py-6 text-center">
                  <p className="text-sm text-muted-foreground">No plates found for &quot;{query}&quot;</p>
                  <button
                    onClick={() => handleSearch(query)}
                    className="mt-2 text-xs text-foreground underline underline-offset-4"
                  >
                    Search all plates →
                  </button>
                </div>
              ) : (
                <>
                  <div className="py-1">
                    {results.map((plate) => (
                      <Link
                        key={plate.id}
                        href={`/plates/${plate.slug}`}
                        onClick={() => setOpen(false)}
                        className="flex items-start gap-3 px-4 py-3 hover:bg-muted transition-colors"
                      >
                        <div className="mt-0.5 text-muted-foreground shrink-0">
                          <GitBranch className="h-4 w-4" />
                        </div>
                        <div className="text-left min-w-0">
                          <p className="text-sm font-medium text-foreground truncate">{plate.name}</p>
                          {plate.description && (
                            <p className="text-xs text-muted-foreground truncate mt-0.5">{plate.description}</p>
                          )}
                          <p className="text-xs text-muted-foreground/60 mt-0.5 capitalize">{plate.category}</p>
                        </div>
                      </Link>
                    ))}
                  </div>
                  <div className="border-t border-border px-4 py-2.5">
                    <button
                      onClick={() => handleSearch(query)}
                      className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                    >
                      See all results for &quot;{query}&quot; →
                    </button>
                  </div>
                </>
              )}
            </div>
          )}
        </div>

        <p className="mt-2 text-xs text-muted-foreground">
          Tip: Use <code className="font-mono">-</code> to exclude words from your search. Example:{" "}
          <span className="font-semibold text-foreground">nodejs-starter</span>
        </p>
      </div>

      <div className="flex flex-wrap justify-center gap-2 max-w-2xl">
        {sampleQueries.map((q) => (
          <button
            key={q}
            onClick={() => handleSearch(q)}
            className="px-3 py-1 text-xs border border-border text-muted-foreground hover:text-foreground hover:border-foreground/40 transition-colors"
          >
            {q}
          </button>
        ))}
      </div>

      <div className="w-full max-w-xl grid grid-cols-4 gap-2">
        {socialItems.map((item, idx) => (
          <Link
            key={`${item.type}-${idx}`}
            href={item.link}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center justify-center gap-1.5 border border-border bg-card px-3 py-2 text-xs text-foreground transition-colors hover:bg-muted"
          >
            {socialIcon(item.type)}
            <span>{socialLabel(item.type)}</span>
          </Link>
        ))}
      </div>

      <div className="flex items-center gap-8 text-center">
        <div>
          <p className="text-3xl font-bold text-foreground">{stats ? formatCount(stats.total_plates) : "-"}</p>
          <p className="text-sm text-muted-foreground">plates</p>
        </div>
        <div className="h-10 w-px bg-border" />
        <div>
          <p className="text-3xl font-bold text-foreground">{stats ? formatCount(stats.total_contributors) : "-"}</p>
          <p className="text-sm text-muted-foreground">contributors</p>
        </div>
      </div>

      <p className="text-sm text-muted-foreground">
        KikPlate Hub is an Open Source project
      </p>

      <div className="absolute bottom-8 flex flex-col items-center gap-1 text-muted-foreground/40">
        <p className="text-xs">scroll to explore</p>
        <div className="h-4 w-px bg-border" />
      </div>

    </div>
  )
}