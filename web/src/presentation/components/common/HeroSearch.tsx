"use client"

import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Search, GitBranch, Loader2, Github, Linkedin, ChevronDown, HelpCircle } from "lucide-react"
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
  const [showHints, setShowHints] = useState(false)

  const { data, isLoading } = usePlates({ search: query, limit: 6 })
  const { data: stats } = useStats()
  const { data: appConfig } = useConfig()

  const results = data?.data ?? []
  const showDropdown = open && query.trim().length > 1
  const titleLines = (appConfig?.banner_title ?? "The biggest library of\nstarter boilerplates").split("\n")
  const socialItems = (appConfig?.social_media ?? [])
    .filter((s) => s.link && s.link !== "#")
    .slice(0, 6)
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
        setShowHints(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  return (
    <div className="relative min-h-[calc(100vh-12rem)] bg-background flex flex-col items-center justify-center px-4 text-center">

      <div className="flex flex-col items-center gap-6 w-full max-w-6xl">

        <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight text-foreground w-full leading-[1.1]">
          {titleLines.map((line, idx) => (
            <span key={`${line}-${idx}`}>
              {line}
              {idx < titleLines.length - 1 && <br />}
            </span>
          ))}
        </h1>

        <p className="text-base sm:text-lg text-muted-foreground max-w-xl leading-relaxed">
          Discover, share and scaffold production-ready project templates. Built by the community, for the community.
        </p>

        <div className="w-full max-w-4xl mt-2" ref={containerRef}>
          <div className="relative">
            <div className="flex items-center border border-border bg-card px-4 gap-3 focus-within:border-foreground/30 focus-within:ring-1 focus-within:ring-foreground/10 transition-all">
              <Search className="h-5 w-5 text-muted-foreground shrink-0" />
              <input
                className="h-14 w-full bg-transparent text-base text-foreground outline-none placeholder:text-muted-foreground"
                placeholder="Search plates... e.g. golang clean architecture"
                value={query}
                onChange={(e) => {
                  setQuery(e.target.value)
                  setOpen(true)
                  setShowHints(false)
                }}
                onKeyDown={handleKeyDown}
                onFocus={() => setOpen(true)}
              />
              {isLoading && query.trim().length > 1 && (
                <Loader2 className="h-4 w-4 text-muted-foreground animate-spin shrink-0" />
              )}
              <button
                type="button"
                onClick={() => setShowHints(!showHints)}
                className="shrink-0 text-muted-foreground hover:text-foreground transition-colors"
                aria-label="Search tips"
              >
                <HelpCircle className="h-4 w-4" />
              </button>
            </div>

            {showHints && (
              <div className="absolute top-full left-0 right-0 z-50 border border-border border-t-0 bg-card shadow-lg shadow-black/10 p-4">
                <div className="flex items-center justify-between mb-3">
                  <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Search tips</p>
                  <button onClick={() => setShowHints(false)} className="text-xs text-muted-foreground hover:text-foreground transition-colors">Close</button>
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-left">
                  <div className="space-y-2">
                    <div>
                      <p className="text-xs font-medium text-foreground">By name or keyword</p>
                      <p className="text-xs text-muted-foreground mt-0.5"><code className="bg-muted px-1 py-0.5">golang starter</code></p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-foreground">By framework</p>
                      <p className="text-xs text-muted-foreground mt-0.5"><code className="bg-muted px-1 py-0.5">spring-boot</code> <code className="bg-muted px-1 py-0.5">nestjs</code></p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-foreground">By language</p>
                      <p className="text-xs text-muted-foreground mt-0.5"><code className="bg-muted px-1 py-0.5">python</code> <code className="bg-muted px-1 py-0.5">java</code> <code className="bg-muted px-1 py-0.5">go</code></p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <div>
                      <p className="text-xs font-medium text-foreground">Exclude words</p>
                      <p className="text-xs text-muted-foreground mt-0.5">Use <code className="bg-muted px-1 py-0.5">-</code> prefix: <code className="bg-muted px-1 py-0.5">nodejs -express</code></p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-foreground">By description</p>
                      <p className="text-xs text-muted-foreground mt-0.5"><code className="bg-muted px-1 py-0.5">clean architecture</code> <code className="bg-muted px-1 py-0.5">docker-compose</code></p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-foreground">Combine terms</p>
                      <p className="text-xs text-muted-foreground mt-0.5"><code className="bg-muted px-1 py-0.5">react typescript starter</code></p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {showDropdown && (
              <div className="absolute top-full left-0 right-0 z-50 border border-border border-t-0 bg-card shadow-lg shadow-black/20">
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
        </div>

        <p className="text-sm text-muted-foreground">
          You can also{" "}
          <Link href="/explore" className="text-foreground underline underline-offset-4 hover:text-primary transition-colors">
            browse all plates
          </Link>
          {" "}or try one of the sample queries:
        </p>

        <div className="flex flex-wrap justify-center gap-2 max-w-2xl">
          {sampleQueries.map((q) => (
            <button
              key={q}
              onClick={() => handleSearch(q)}
              className="px-3 py-1.5 text-xs border border-border text-muted-foreground hover:text-foreground hover:border-foreground/30 hover:bg-card transition-all"
            >
              {q}
            </button>
          ))}
        </div>

        <div className="flex items-center gap-6 mt-4">
          <div className="flex items-center gap-8 border border-border bg-card/50 px-8 py-4">
            <div className="text-center">
              <p className="text-2xl font-bold tabular-nums text-foreground">{stats ? formatCount(stats.total_plates) : "—"}</p>
              <p className="text-xs text-muted-foreground mt-0.5">plates</p>
            </div>
            <div className="h-8 w-px bg-border" />
            <div className="text-center">
              <p className="text-2xl font-bold tabular-nums text-foreground">{stats ? formatCount(stats.total_contributors) : "—"}</p>
              <p className="text-xs text-muted-foreground mt-0.5">contributors</p>
            </div>
            <div className="h-8 w-px bg-border" />
            <div className="text-center">
              <p className="text-2xl font-bold tabular-nums text-foreground">{stats ? formatCount(stats.total_categories) : "—"}</p>
              <p className="text-xs text-muted-foreground mt-0.5">categories</p>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4 mt-2">
          <p className="text-xs text-muted-foreground">KikPlate is an Open Source project</p>
          <span className="text-muted-foreground/30">|</span>
          <div className="flex items-center gap-2">
            {socialItems.map((item, idx) => (
              <Link
                key={`footer-${item.type}-${idx}`}
                href={item.link}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center h-8 w-8 border border-border text-muted-foreground transition-colors hover:text-foreground hover:bg-card"
                title={socialLabel(item.type)}
              >
                {socialIcon(item.type)}
              </Link>
            ))}
          </div>
        </div>

      </div>

      <div className="absolute bottom-6 flex flex-col items-center gap-1.5 text-muted-foreground/30 animate-bounce">
        <ChevronDown className="h-4 w-4" />
      </div>

    </div>
  )
}