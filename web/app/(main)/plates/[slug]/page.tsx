import { notFound } from "next/navigation"
import Link from "next/link"
import {
  GitBranch, FileText, Download, Star,
  Shield, Tag, CheckCircle2, Calendar,
  ExternalLink, ArrowLeft, Clock
} from "lucide-react"
import { plateRepository } from "@/src/data/repositories/PlateRepository"
import { fetchRepoFile } from "@/src/data/repositories/githubClient"
import { formatCount, tierColour, relativeTime } from "@/src/presentation/utils/plateUtils"
import { UseButtonClient } from "@/src/presentation/components/plates/UseButtonClient"
import { PlateContentTabs } from "@/src/presentation/components/plates/PlateContentTabs"

interface Props {
  params: Promise<{ slug: string }>
}

export default async function PlateDetailPage({ params }: Props) {
  const { slug } = await params
  const plate = await plateRepository.getBySlug(slug).catch(() => null)
  if (!plate) notFound()

  let readme: string | null = null
  let license: string | null = null

  if (plate.type === "repository" && plate.repo_url) {
    const branch = plate.branch ?? "main"
    ;[readme, license] = await Promise.all([
      fetchRepoFile(plate.repo_url, branch, "README.md"),
      fetchRepoFile(plate.repo_url, branch, "LICENSE"),
    ])
  }

  return (
    <div className="min-h-screen bg-background">

      {/* ── full width header ── */}
      <div className="border-b border-border bg-muted/20">
        <div className="container mx-auto px-4 py-10">

          {/* back */}
          <Link
            href="/explore"
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6"
          >
            <ArrowLeft className="h-3.5 w-3.5" />
            Back to explore
          </Link>

          {/* type + category + verified */}
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-3">
            {plate.type === "file"
              ? <FileText className="h-4 w-4" />
              : <GitBranch className="h-4 w-4" />
            }
            <span className="capitalize">{plate.type}</span>
            <span>·</span>
            <span className="capitalize">{plate.category}</span>
            {plate.is_verified && (
              <>
                <span>·</span>
                <span className="flex items-center gap-1 text-green-600">
                  <CheckCircle2 className="h-3.5 w-3.5" />
                  Verified
                </span>
              </>
            )}
            {plate.published_at && (
              <>
                <span>·</span>
                <span>Published {relativeTime(plate.published_at)}</span>
              </>
            )}
          </div>

          {/* name */}
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground mb-3">
            {plate.name}
          </h1>

          {/* description */}
          {plate.description && (
            <p className="text-muted-foreground leading-relaxed max-w-3xl mb-6">
              {plate.description}
            </p>
          )}

          {/* tags + badges row */}
          <div className="flex flex-wrap items-center gap-4">

            {/* tags */}
            {plate.tags && plate.tags.length > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {plate.tags.map((t) => (
                  <Link
                    key={t.id}
                    href={`/explore?tag=${t.tag}`}
                    className="inline-flex items-center gap-1 px-2 py-0.5 text-xs border border-border text-muted-foreground hover:text-foreground hover:border-foreground/30 transition-colors"
                  >
                    <Tag className="h-2.5 w-2.5" />
                    {t.tag}
                  </Link>
                ))}
              </div>
            )}

            {/* divider */}
            {plate.tags && plate.tags.length > 0 && plate.badges && plate.badges.length > 0 && (
              <div className="h-4 w-px bg-border" />
            )}

            {/* badges */}
            {plate.badges && plate.badges.length > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {plate.badges.map((pb) => pb.badge && (
                  <span
                    key={pb.id}
                    className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs border font-medium ${tierColour(pb.badge.tier)}`}
                  >
                    <Shield className="h-2.5 w-2.5" />
                    {pb.badge.name}
                  </span>
                ))}
              </div>
            )}

          </div>

          {/* stats row */}
          <div className="flex items-center gap-6 mt-6 text-sm text-muted-foreground">
            <div className="flex items-center gap-1.5">
              <Download className="h-4 w-4" />
              <span><span className="font-semibold text-foreground">{formatCount(plate.use_count)}</span> uses</span>
            </div>
            {plate.avg_rating > 0 && (
              <div className="flex items-center gap-1.5">
                <Star className="h-4 w-4" />
                <span><span className="font-semibold text-foreground">{plate.avg_rating.toFixed(1)}</span> / 5</span>
              </div>
            )}
            <div className="flex items-center gap-1.5">
              <Calendar className="h-4 w-4" />
              <span>Updated {relativeTime(plate.updated_at)}</span>
            </div>
            {plate.type === "repository" && plate.last_synced_at && (
              <div className="flex items-center gap-1.5">
                <Clock className="h-4 w-4" />
                <span>Synced {relativeTime(plate.last_synced_at)}</span>
              </div>
            )}
          </div>

        </div>
      </div>

      {/* ── two column content ── */}
      <div className="container mx-auto px-4 py-10">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-10">

          {/* left — main content */}
          <div className="lg:col-span-2">
            {plate.type === "repository" ? (
              <PlateContentTabs readme={readme} license={license} />
            ) : plate.content ? (
              <div className="space-y-2">
                <p className="text-xs uppercase tracking-widest text-muted-foreground font-medium mb-3">
                  {plate.filename ?? "Content"}
                </p>
                <div className="border border-border">
                  <div className="flex items-center px-4 py-2 border-b border-border bg-muted/50">
                    <span className="text-xs text-muted-foreground font-mono">{plate.filename}</span>
                  </div>
                  <pre className="p-4 text-xs text-foreground font-mono overflow-x-auto leading-relaxed whitespace-pre">
                    {plate.content}
                  </pre>
                </div>
              </div>
            ) : null}
          </div>

          {/* right — sidebar */}
          <div className="space-y-6">

            {/* use button */}
            <UseButtonClient
              plateId={plate.id}
              slug={plate.slug}
              repoUrl={plate.repo_url}
            />

            {/* repository info */}
            {plate.type === "repository" && plate.repo_url && (
              <div className="border border-border p-4 space-y-3">
                <p className="text-xs uppercase tracking-widest text-muted-foreground font-medium">Repository</p>
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-mono text-foreground truncate">
                    {plate.repo_url.replace("https://github.com/", "")}
                  </span>
                  <Link
                    href={plate.repo_url}
                    target="_blank"
                    className="shrink-0 text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                  </Link>
                </div>
                {plate.branch && (
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <GitBranch className="h-3 w-3" />
                    <span className="font-mono">{plate.branch}</span>
                  </div>
                )}
                {plate.sync_status && (
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-muted-foreground">Sync</span>
                    <span className="capitalize font-medium text-foreground">{plate.sync_status}</span>
                  </div>
                )}
              </div>
            )}

            {/* info */}
            <div className="border border-border p-4 space-y-3">
              <p className="text-xs uppercase tracking-widest text-muted-foreground font-medium">Info</p>
              <div className="space-y-2.5">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Status</span>
                  <span className="capitalize font-medium text-foreground">{plate.status}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Visibility</span>
                  <span className="capitalize font-medium text-foreground">{plate.visibility}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Type</span>
                  <span className="capitalize font-medium text-foreground">{plate.type}</span>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>

    </div>
  )
}