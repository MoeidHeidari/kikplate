"use client"

import Link from "next/link"
import { useStats } from "@/src/presentation/hooks/usePlates"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import { Github, Users, Layers, Download, Package, Sparkles, ArrowRight, CheckCircle2 } from "lucide-react"

const SlackIcon = () => (
  <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor">
    <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z" />
  </svg>
)

interface StatRowProps {
  icon: React.ReactNode
  value: string
  label: string
  inverted?: boolean
}

function StatRow({ icon, value, label, inverted = false }: StatRowProps) {
  return (
    <div
      className={`flex items-center justify-between border-b py-4 last:border-0 ${
        inverted ? "border-background/20 dark:border-secondary/30" : "border-border"
      }`}
    >
      <div
        className={`flex items-center gap-3 ${
          inverted ? "text-background/75 dark:text-secondary/70" : "text-muted-foreground"
        }`}
      >
        {icon}
        <span className={`text-sm ${inverted ? "text-background dark:text-secondary" : "text-foreground"}`}>
          {label}
        </span>
      </div>
      <span className={`text-2xl font-bold tabular-nums ${inverted ? "text-background dark:text-secondary" : "text-foreground"}`}>
        {value}
      </span>
    </div>
  )
}

export function StatsBanner() {
  const { data } = useStats()

  const leftStats = [
    {
      icon: <Package className="h-4 w-4" />,
      value: data ? formatCount(data.total_plates) : "—",
      label: "Plates",
    },
    {
      icon: <Users className="h-4 w-4" />,
      value: data ? formatCount(data.total_contributors) : "—",
      label: "Contributors",
    },
    {
      icon: <Layers className="h-4 w-4" />,
      value: data ? formatCount(data.total_categories) : "—",
      label: "Categories",
    },
  ]

  const rightStats = [
    {
      icon: <Download className="h-4 w-4" />,
      value: data ? formatCount(data.total_bookmarks) : "—",
      label: "Bookmarks",
    },
    {
      icon: <Github className="h-4 w-4" />,
      value: "—",
      label: "GitHub Stars",
    },
    {
      icon: <SlackIcon />,
      value: "—",
      label: "Slack Members",
    },
  ]

  return (
    <section className="bg-background py-20">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-1 border border-border lg:grid-cols-2">

          <div className="bg-card p-6 sm:p-8">
            <div className="mb-4 inline-flex items-center gap-2 border border-border bg-muted/40 px-2.5 py-1 text-muted-foreground">
              <Sparkles className="h-3.5 w-3.5" />
              <p className="text-xs font-semibold uppercase tracking-widest">Community</p>
            </div>

            <h2 className="text-3xl font-bold leading-tight tracking-tight text-foreground">
              Built by the community,
              <br />
              for the community
            </h2>

            <p className="mt-4 max-w-lg leading-relaxed text-muted-foreground">
              KikPlate is an open source registry of production-ready project plates.
              Developers submit real plates, maintainers review structure and metadata,
              and every plate remains free to use.
            </p>

            <div className="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="border border-border bg-background p-3">
                <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Open model</p>
                <p className="mt-1 text-sm text-foreground">Anyone can contribute quality plates.</p>
              </div>
              <div className="border border-border bg-background p-3">
                <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Clear metadata</p>
                <p className="mt-1 text-sm text-foreground">Owner, tags, category, and usage context.</p>
              </div>
            </div>

            <div className="mt-6 space-y-2 border-t border-border pt-4 text-sm text-muted-foreground">
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4" />Production-ready plates</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4" />Transparent and auditable plates</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4" />Community driven improvement loop</p>
            </div>

            <div className="mt-6 flex flex-wrap items-center gap-3">
              <Link
                href="https://github.com/kickplate"
                target="_blank"
                className="inline-flex items-center gap-1.5 border border-border px-3 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-muted"
              >
                View on GitHub
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
              <Link
                href="#"
                className="inline-flex items-center gap-1.5 border border-border px-3 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-muted"
              >
                Join Slack
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>

          <div className="border-t border-border bg-[oklch(0.22_0.02_250)] p-6 text-background dark:bg-secondary-foreground dark:text-secondary sm:p-8 lg:border-t-0 lg:border-l lg:border-border">
            <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-background/80 dark:text-secondary/80">
              By the numbers
            </p>
            <div className="grid grid-cols-1 gap-x-10 sm:grid-cols-2">
              <div className="sm:pr-4">
                {leftStats.map((s) => (
                  <StatRow key={s.label} icon={s.icon} value={s.value} label={s.label} inverted />
                ))}
              </div>
              <div className="sm:border-l sm:border-background/20 dark:sm:border-secondary/30 sm:pl-4">
                {rightStats.map((s) => (
                  <StatRow key={s.label} icon={s.icon} value={s.value} label={s.label} inverted />
                ))}
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}