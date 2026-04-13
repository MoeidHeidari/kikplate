import { notFound } from "next/navigation"
import { cookies } from "next/headers"
import Link from "next/link"
import Image from "next/image"
import {
  GitBranch, FileText, Heart, Star,
  Tag, CheckCircle2, Calendar, CalendarCheck, Layers,
  ExternalLink, ArrowLeft
} from "lucide-react"
import { fetchRepoFile, fetchRepoReadme, fetchRepoTree } from "@/src/data/repositories/githubClient"
import { formatCount, relativeTime } from "@/src/presentation/utils/plateUtils"
import { UseButtonClient } from "@/src/presentation/components/plates/UseButtonClient"
import { BookmarkButtonClient } from "@/src/presentation/components/plates/BookmarkButtonClient"
import { PlateContentTabs } from "@/src/presentation/components/plates/PlateContentTabs"
import { PlateHeaderTabs } from "@/src/presentation/components/plates/PlateHeaderTabs"
import { PlateRatingCard } from "@/src/presentation/components/plates/PlateRatingCard"
import { BadgeShowcase } from "@/src/presentation/components/plates/BadgeShowcase"
import { HeaderBadges } from "@/src/presentation/components/plates/HeaderBadges"
import { RelatedPlates } from "@/src/presentation/components/plates/RelatedPlates"
import type { Plate, PlateListResponse } from "@/src/domain/entities/Plate"
import type { Badge } from "@/src/domain/entities/Badge"
import type { AppConfig } from "@/src/domain/entities/Config"
import { getServerApiBaseUrl } from "@/src/lib/api"

interface Props {
  params: Promise<{ slug: string[] }>
}

export default async function PlateDetailPage({ params }: Props) {
  const { slug } = await params
  const rawSlug = Array.isArray(slug) ? slug.join("/") : slug

  let normalizedSlug = rawSlug
  try {
    normalizedSlug = decodeURIComponent(rawSlug)
  } catch {}
  const base = await getServerApiBaseUrl()
  const token = (await cookies()).get("kp_token")?.value

  const res = await fetch(`${base}/plates/${encodeURIComponent(normalizedSlug)}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    cache: "no-store",
  })

  if (!res.ok) {
    notFound()
  }

  const plate = (await res.json()) as Plate

  const relatedQs = new URLSearchParams()
  relatedQs.set("limit", "12")
  let relatedExploreHref = ""
  let relatedExploreLabel = ""
  if (plate.category?.trim()) {
    relatedQs.set("category", plate.category.trim())
    relatedExploreHref = `/explore?category=${encodeURIComponent(plate.category.trim())}`
    relatedExploreLabel = `More in ${plate.category}`
  } else if (plate.tags?.[0]?.tag) {
    relatedQs.set("tag", plate.tags[0].tag)
    relatedExploreHref = `/explore?tag=${encodeURIComponent(plate.tags[0].tag)}`
    relatedExploreLabel = `More tagged "${plate.tags[0].tag}"`
  }

  const relatedUrl =
    relatedQs.has("category") || relatedQs.has("tag")
      ? `${base}/plates?${relatedQs.toString()}`
      : null

  const [badgesRes, configRes, relatedRes] = await Promise.all([
    fetch(`${base}/badges`, { cache: "no-store" }),
    fetch(`${base}/config`, { cache: "no-store" }),
    relatedUrl
      ? fetch(relatedUrl, {
          headers: token ? { Authorization: `Bearer ${token}` } : undefined,
          cache: "no-store",
        })
      : Promise.resolve({ ok: false } as Response),
  ])
  const allBadges: Badge[] = badgesRes.ok ? await badgesRes.json() : []
  const appConfig: AppConfig | null = configRes.ok ? await configRes.json() : null

  let relatedPlates: Plate[] = []
  if (relatedRes.ok) {
    const body = (await relatedRes.json()) as PlateListResponse
    relatedPlates = (body.data ?? [])
      .filter((p) => p.slug !== plate.slug)
      .slice(0, 6)
  }

  let readme: string | null = null
  let license: string | null = null
  let tree = null

  if (plate.type === "repository" && plate.repo_url) {
    const branch = plate.branch ?? "main"
    ;[readme, license, tree] = await Promise.all([
      fetchRepoReadme(plate.repo_url, branch),
      fetchRepoFile(plate.repo_url, branch, "LICENSE"),
      fetchRepoTree(plate.repo_url, branch),
    ])
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border">
        <div className="container mx-auto px-4 pt-6">
          <div className="mb-5 flex items-center justify-between gap-4">
            <Link
              href="/explore"
              className="inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              Back to explore
            </Link>
            <div className="hidden sm:flex gap-3">
              <BookmarkButtonClient
                plateId={plate.id}
                isBookmarked={plate.is_bookmarked}
              />
              <UseButtonClient
                plateId={plate.id}
                slug={plate.slug}
                repoUrl={plate.repo_url}
              />
            </div>
          </div>

          <h1 className="max-w-4xl text-3xl font-black leading-tight tracking-tight text-foreground sm:text-4xl">
            {plate.name}
          </h1>

          <div className="mt-2 flex flex-wrap items-center gap-2">
            {plate.is_verified && (
              <span className="inline-flex items-center gap-1.5 border border-emerald-400/50 bg-emerald-500/10 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:text-emerald-400">
                <CheckCircle2 className="h-3.5 w-3.5" />
                Verified
              </span>
            )}
            <span
              className={`inline-flex items-center gap-1.5 border px-2.5 py-1 text-xs font-medium uppercase tracking-wide ${
                plate.visibility === "private"
                  ? "border-amber-400/50 bg-amber-500/10 text-amber-700 dark:text-amber-300"
                  : plate.visibility === "unlisted"
                    ? "border-blue-400/50 bg-blue-500/10 text-blue-700 dark:text-blue-300"
                    : "border-emerald-400/50 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300"
              }`}
            >
              {plate.visibility}
            </span>
            <HeaderBadges badges={plate.badges ?? []} />
          </div>

          {plate.description && (
            <p className="mt-3 max-w-3xl text-sm leading-relaxed text-muted-foreground sm:text-base">
              {plate.description}
            </p>
          )}

          <div className="pb-4" />

          <div>
            <PlateHeaderTabs
              isRepository={plate.type === "repository"}
              hasReadme={Boolean(readme)}
              hasLicense={Boolean(license)}
              hasTree={Boolean(tree?.length)}
            />
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-10 pb-24 sm:pb-10">
        <div className="grid grid-cols-1 gap-7 xl:gap-8 lg:grid-cols-12">
          <section className="lg:col-span-9">
            <div>
              <PlateContentTabs
                readme={readme}
                license={license}
                tree={tree}
                repoUrl={plate.repo_url}
                branch={plate.branch ?? undefined}
              />
            </div>
          </section>

          <aside className="space-y-6 lg:col-span-3">
            <PlateRatingCard
              plateId={plate.id}
              plateSlug={plate.slug}
              plateOwnerId={plate.owner_id}
              avgRating={plate.avg_rating}
              userRating={plate.user_rating}
            />

            <div className="border border-border bg-card p-5">
              <p className="mb-3 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">Overview</p>

              {plate.organization ? (
                <div className="mb-4 flex items-center gap-2.5 border-b border-border pb-4">
                  <Link href={`/orgs/${encodeURIComponent(plate.organization.name)}`} className="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden border border-border bg-muted hover:border-foreground/30 transition-colors">
                    {plate.organization.logo_url ? (
                      <Image
                        src={plate.organization.logo_url}
                        alt={`${plate.organization.name} logo`}
                        width={48}
                        height={48}
                        unoptimized
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <span className="text-xs font-semibold uppercase text-muted-foreground">
                        {plate.organization.name.slice(0, 2)}
                      </span>
                    )}
                  </Link>
                  <div>
                    <p className="text-xs text-muted-foreground">Organization</p>
                    <Link href={`/orgs/${encodeURIComponent(plate.organization.name)}`} className="text-sm font-semibold text-foreground hover:underline">
                      {plate.organization.name}
                    </Link>
                    {plate.organization.owner && (
                      <p className="text-xs text-muted-foreground">
                        Owner:{" "}
                        <Link href={`/users/${encodeURIComponent(plate.organization.owner.username ?? "")}`} className="font-medium text-foreground hover:underline">
                          {plate.organization.owner.username ?? plate.organization.owner.display_name ?? "Unknown"}
                        </Link>
                      </p>
                    )}
                  </div>
                </div>
              ) : plate.owner ? (
                <div className="mb-4 flex items-center gap-2.5 border-b border-border pb-4">
                  <Link href={`/users/${encodeURIComponent(plate.owner.username ?? "")}`} className="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden border border-border bg-muted hover:border-foreground/30 transition-colors">
                    {plate.owner.avatar_url ? (
                      <Image
                        src={plate.owner.avatar_url}
                        alt={plate.owner.username ?? "owner"}
                        width={48}
                        height={48}
                        unoptimized
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <span className="text-xs font-semibold uppercase text-muted-foreground">
                        {(plate.owner.username ?? plate.owner.display_name ?? "?").slice(0, 2)}
                      </span>
                    )}
                  </Link>
                  <div>
                    <p className="text-xs text-muted-foreground">Owner</p>
                    <Link href={`/users/${encodeURIComponent(plate.owner.username ?? "")}`} className="text-sm font-semibold text-foreground hover:underline">
                      {plate.owner.username ?? plate.owner.display_name ?? "Unknown"}
                    </Link>
                  </div>
                </div>
              ) : null}

              <div className="space-y-3">
                <div className="flex items-center justify-between text-sm">
                  <span className="flex items-center gap-1.5 text-muted-foreground"><FileText className="h-3.5 w-3.5" /> Type</span>
                  <span className="font-semibold capitalize text-foreground">{plate.type}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="flex items-center gap-1.5 text-muted-foreground">
                    <Layers className="h-3.5 w-3.5 shrink-0" /> Category
                  </span>
                  <span className="font-semibold capitalize text-foreground">{plate.category}</span>
                </div>
                {plate.published_at && (
                  <div className="flex items-center justify-between text-sm">
                    <span className="flex items-center gap-1.5 text-muted-foreground">
                      <CalendarCheck className="h-3.5 w-3.5 shrink-0" /> Published
                    </span>
                    <span className="font-semibold text-foreground">{relativeTime(plate.published_at)}</span>
                  </div>
                )}
                <div className="flex items-center justify-between text-sm">
                  <span className="flex items-center gap-1.5 text-muted-foreground"><Heart className="h-3.5 w-3.5" /> Bookmarks</span>
                  <span className="font-semibold text-foreground">{formatCount(plate.bookmark_count)}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="flex items-center gap-1.5 text-muted-foreground"><Star className="h-3.5 w-3.5" /> Rating</span>
                  <span className="font-semibold text-foreground">{plate.avg_rating > 0 ? plate.avg_rating.toFixed(1) : "-"}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="flex items-center gap-1.5 text-muted-foreground"><Calendar className="h-3.5 w-3.5" /> Updated</span>
                  <span className="font-semibold text-foreground">{relativeTime(plate.updated_at)}</span>
                </div>
              </div>

              {(plate.tags?.length ?? 0) > 0 ? (
                <div className="mt-4 border-t border-border pt-4">
                  <div className="flex flex-wrap gap-2">
                    {plate.tags?.map((t) => (
                      <Link
                        key={t.id}
                        href={`/explore?tag=${t.tag}`}
                        className="inline-flex items-center gap-1 border border-border bg-background px-2 py-0.5 text-xs font-medium text-muted-foreground transition-colors hover:border-foreground/30 hover:text-foreground"
                      >
                        <Tag className="h-2.5 w-2.5" />
                        {t.tag}
                      </Link>
                    ))}
                  </div>
                </div>
              ) : null}
            </div>

            <BadgeShowcase
              allBadges={allBadges}
              plateBadges={plate.badges ?? []}
              plateOwnerId={plate.owner_id}
              plateSlug={plate.slug}
              requestUrl={appConfig?.badge_request_url}
            />

            {plate.type === "repository" && plate.repo_url && (
              <div className="border border-border bg-card p-5">
                <p className="mb-3 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">Repository</p>
                <div className="flex items-center justify-between gap-2">
                  <span className="truncate font-mono text-xs text-foreground">{plate.repo_url.replace("https://github.com/", "")}</span>
                  <Link
                    href={plate.repo_url}
                    target="_blank"
                    className="shrink-0 text-muted-foreground transition-colors hover:text-foreground"
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                  </Link>
                </div>
                {plate.branch && (
                  <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
                    <GitBranch className="h-3 w-3" />
                    <span className="font-mono">{plate.branch}</span>
                  </div>
                )}
              </div>
            )}

            {relatedPlates.length > 0 && relatedExploreHref ? (
              <RelatedPlates
                plates={relatedPlates}
                exploreHref={relatedExploreHref}
                exploreLabel={relatedExploreLabel}
              />
            ) : null}

          </aside>
        </div>
      </div>

      <div className="fixed inset-x-4 bottom-4 z-40 space-y-2 sm:hidden">
        <UseButtonClient
          plateId={plate.id}
          slug={plate.slug}
          repoUrl={plate.repo_url}
          prominent
        />
        <BookmarkButtonClient
          plateId={plate.id}
          isBookmarked={plate.is_bookmarked}
          className="w-full"
        />
      </div>
    </div>
  )
}