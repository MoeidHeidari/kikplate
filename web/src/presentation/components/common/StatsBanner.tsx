"use client"

import Link from "next/link"
import { useStats } from "@/src/presentation/hooks/usePlates"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { useExternalCommunityStats } from "@/src/presentation/hooks/useExternalCommunityStats"
import { getSocialLink } from "@/src/lib/socialLinks"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import { Github, Users, Layers, Download, Package, Sparkles, ArrowRight, CheckCircle2 } from "lucide-react"

const SlackIcon = () => (
  <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor">
    <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z" />
  </svg>
)

interface StatCardProps {
  icon: React.ReactNode
  value: string
  label: string
}

function StatCard({ icon, value, label }: StatCardProps) {
  return (
    <div className="flex flex-col items-center gap-2 border border-border bg-card p-5 text-center">
      <div className="text-muted-foreground">{icon}</div>
      <p className="text-3xl font-bold tabular-nums text-foreground">{value}</p>
      <p className="text-xs text-muted-foreground">{label}</p>
    </div>
  )
}

export function StatsBanner() {
  const { data } = useStats()
  const { data: appConfig } = useConfig()
  const { data: externalStats, isLoading: externalLoading } = useExternalCommunityStats()
  const githubUrl = getSocialLink(appConfig?.social_media, "github") ?? "https://github.com/kikplate/kikplate"
  const slackUrl = getSocialLink(appConfig?.social_media, "slack")

  const githubStarsDisplay =
    externalLoading ? "…" : externalStats?.githubStars != null ? formatCount(externalStats.githubStars) : "—"
  const slackMembersDisplay =
    externalLoading ? "…" : externalStats?.slackMembers != null ? formatCount(externalStats.slackMembers) : "—"

  const allStats = [
    { icon: <Package className="h-5 w-5" />, value: data ? formatCount(data.total_plates) : "—", label: "Plates" },
    { icon: <Users className="h-5 w-5" />, value: data ? formatCount(data.total_contributors) : "—", label: "Contributors" },
    { icon: <Layers className="h-5 w-5" />, value: data ? formatCount(data.total_categories) : "—", label: "Categories" },
    { icon: <Download className="h-5 w-5" />, value: data ? formatCount(data.total_bookmarks) : "—", label: "Bookmarks" },
    { icon: <Github className="h-5 w-5" />, value: githubStarsDisplay, label: "GitHub Stars" },
    { icon: <SlackIcon />, value: slackMembersDisplay, label: "Slack Members" },
  ]

  return (
    <section className="bg-background py-20">
      <div className="container mx-auto px-4">

        <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">

          <div className="flex flex-col gap-6">
            <div>
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
            </div>

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="border border-border bg-card p-4">
                <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Open model</p>
                <p className="mt-1.5 text-sm text-foreground">Anyone can contribute quality plates.</p>
              </div>
              <div className="border border-border bg-card p-4">
                <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Clear metadata</p>
                <p className="mt-1.5 text-sm text-foreground">Owner, tags, category, and usage context.</p>
              </div>
            </div>

            <div className="space-y-2.5 border-t border-border pt-5 text-sm text-muted-foreground">
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 shrink-0" />Production-ready plates</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 shrink-0" />Transparent and auditable plates</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 shrink-0" />Community driven improvement loop</p>
            </div>

            <div className="flex flex-wrap items-center gap-3">
              <Link
                href={githubUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="group inline-flex items-center gap-1.5 border border-border px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-muted"
              >
                <Github className="h-4 w-4" />
                View on GitHub
                <ArrowRight className="h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" />
              </Link>
              {slackUrl && (
                <Link
                  href={slackUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="group inline-flex items-center gap-1.5 border border-border px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-muted"
                >
                  <SlackIcon />
                  Join Slack
                  <ArrowRight className="h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" />
                </Link>
              )}
            </div>
          </div>

          <div>
            <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              By the numbers
            </p>
            <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
              {allStats.map((s) => (
                <StatCard key={s.label} icon={s.icon} value={s.value} label={s.label} />
              ))}
            </div>
          </div>

        </div>

      </div>
    </section>
  )
}