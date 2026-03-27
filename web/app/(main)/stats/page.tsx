"use client"

import Link from "next/link"
import { useMemo } from "react"
import { BarChart3, Bookmark, Layers, Package, ShieldCheck, Users } from "lucide-react"
import { useStats, usePlates } from "@/src/presentation/hooks/usePlates"
import { useBadges } from "@/src/presentation/hooks/useBadges"
import { buildPlateHref, formatCount } from "@/src/presentation/utils/plateUtils"

type Datum = {
  label: string
  value: number
  href?: string
}

type TrendDatum = {
  label: string
  value: number
}

function KpiCard({
  label,
  value,
  icon,
}: {
  label: string
  value: string
  icon: React.ReactNode
}) {
  return (
    <div className="border border-border bg-card p-4 sm:p-5">
      <div className="mb-2 flex items-center justify-between">
        <p className="text-xs uppercase tracking-widest text-muted-foreground">{label}</p>
        <span className="text-muted-foreground">{icon}</span>
      </div>
      <p className="text-3xl font-black tabular-nums text-foreground">{value}</p>
    </div>
  )
}

function LinePlot({
  title,
  subtitle,
  points,
}: {
  title: string
  subtitle: string
  points: TrendDatum[]
}) {
  const width = 720
  const height = 260
  const padX = 36
  const padY = 24

  const maxY = useMemo(() => {
    if (!points.length) return 1
    return Math.max(...points.map((point) => point.value), 1)
  }, [points])

  const path = useMemo(() => {
    if (!points.length) return ""

    const innerW = width - padX * 2
    const innerH = height - padY * 2
    const step = points.length > 1 ? innerW / (points.length - 1) : 0

    const coords = points.map((point, idx) => {
      const x = padX + idx * step
      const y = padY + innerH - (point.value / maxY) * innerH
      return { x, y }
    })

    return coords.map((coord, idx) => `${idx === 0 ? "M" : "L"} ${coord.x} ${coord.y}`).join(" ")
  }, [maxY, points])

  const yTicks = useMemo(() => {
    return [0, 0.25, 0.5, 0.75, 1].map((ratio) => ({
      label: Math.round(maxY * ratio),
      y: padY + (height - padY * 2) * (1 - ratio),
    }))
  }, [maxY])

  return (
    <section className="border border-border bg-card p-5 sm:p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>

      {points.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data available yet.</p>
      ) : (
        <div className="space-y-3">
          <svg viewBox={`0 0 ${width} ${height}`} className="w-full" role="img" aria-label={title}>
            {yTicks.map((tick, idx) => (
              <g key={`tick-${idx}-${tick.y}-${tick.label}`}>
                <line
                  x1={padX}
                  y1={tick.y}
                  x2={width - padX}
                  y2={tick.y}
                  stroke="currentColor"
                  strokeOpacity="0.16"
                  strokeWidth="1"
                />
                <text x={padX - 8} y={tick.y + 4} textAnchor="end" fontSize="10" fill="currentColor" opacity="0.6">
                  {tick.label}
                </text>
              </g>
            ))}

            <path d={path} fill="none" stroke="currentColor" strokeOpacity="0.9" strokeWidth="2.5" />

            {points.map((point, idx) => {
              const x = padX + (points.length > 1 ? ((width - padX * 2) / (points.length - 1)) * idx : 0)
              const y = padY + (height - padY * 2) - (point.value / maxY) * (height - padY * 2)
              return <circle key={point.label} cx={x} cy={y} r="3" fill="currentColor" />
            })}

            {points.map((point, idx) => {
              const x = padX + (points.length > 1 ? ((width - padX * 2) / (points.length - 1)) * idx : 0)
              return (
                <text
                  key={`${point.label}-x`}
                  x={x}
                  y={height - 6}
                  textAnchor="middle"
                  fontSize="10"
                  fill="currentColor"
                  opacity="0.6"
                >
                  {point.label}
                </text>
              )
            })}
          </svg>
        </div>
      )}
    </section>
  )
}

function ColumnPlot({
  title,
  subtitle,
  rows,
}: {
  title: string
  subtitle: string
  rows: Datum[]
}) {
  const width = 720
  const height = 280
  const padX = 36
  const padY = 20
  const maxY = useMemo(() => {
    if (!rows.length) return 1
    return Math.max(...rows.map((row) => row.value), 1)
  }, [rows])

  return (
    <section className="border border-border bg-card p-5 sm:p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>

      {rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data available yet.</p>
      ) : (
        <div className="space-y-3">
          <svg viewBox={`0 0 ${width} ${height}`} className="w-full" role="img" aria-label={title}>
            <line
              x1={padX}
              y1={height - padY}
              x2={width - padX}
              y2={height - padY}
              stroke="currentColor"
              strokeOpacity="0.3"
              strokeWidth="1"
            />

            {rows.map((row, idx) => {
              const innerW = width - padX * 2
              const slot = innerW / rows.length
              const barW = Math.min(56, slot * 0.72)
              const x = padX + idx * slot + (slot - barW) / 2
              const barH = ((height - padY * 2) * row.value) / maxY
              const y = height - padY - barH
              return (
                <g key={row.label}>
                  <rect x={x} y={y} width={barW} height={barH} fill="currentColor" fillOpacity="0.82" />
                  <text x={x + barW / 2} y={y - 6} textAnchor="middle" fontSize="10" fill="currentColor" opacity="0.75">
                    {formatCount(row.value)}
                  </text>
                  <text
                    x={x + barW / 2}
                    y={height - 7}
                    textAnchor="middle"
                    fontSize="10"
                    fill="currentColor"
                    opacity="0.6"
                  >
                    {row.label.length > 10 ? `${row.label.slice(0, 10)}..` : row.label}
                  </text>
                </g>
              )
            })}
          </svg>
        </div>
      )}
    </section>
  )
}

function DonutPlot({
  title,
  subtitle,
  rows,
}: {
  title: string
  subtitle: string
  rows: Datum[]
}) {
  const size = 280
  const stroke = 34
  const radius = (size - stroke) / 2
  const center = size / 2
  const circumference = 2 * Math.PI * radius
  const total = rows.reduce((sum, row) => sum + row.value, 0)
  const colors = ["#111111", "#3f3f46", "#71717a", "#a1a1aa", "#d4d4d8"]

  let offset = 0

  return (
    <section className="border border-border bg-card p-5 sm:p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>

      {rows.length === 0 || total === 0 ? (
        <p className="text-sm text-muted-foreground">No data available yet.</p>
      ) : (
        <div className="grid grid-cols-1 items-center gap-4 sm:grid-cols-[auto_1fr]">
          <svg viewBox={`0 0 ${size} ${size}`} className="mx-auto w-56" role="img" aria-label={title}>
            <circle cx={center} cy={center} r={radius} stroke="currentColor" strokeOpacity="0.12" strokeWidth={stroke} fill="none" />
            {rows.map((row, idx) => {
              const portion = row.value / total
              const dash = portion * circumference
              const style = {
                strokeDasharray: `${dash} ${circumference - dash}`,
                strokeDashoffset: -offset,
              }
              offset += dash
              return (
                <circle
                  key={row.label}
                  cx={center}
                  cy={center}
                  r={radius}
                  stroke={colors[idx % colors.length]}
                  strokeWidth={stroke}
                  strokeLinecap="butt"
                  fill="none"
                  transform={`rotate(-90 ${center} ${center})`}
                  style={style}
                />
              )
            })}
            <text x={center} y={center - 4} textAnchor="middle" fontSize="28" fontWeight="700" fill="currentColor">
              {total}
            </text>
            <text x={center} y={center + 18} textAnchor="middle" fontSize="11" fill="currentColor" opacity="0.65">
              total badges
            </text>
          </svg>

          <div className="space-y-2">
            {rows.map((row, idx) => {
              const pct = ((row.value / total) * 100).toFixed(1)
              return (
                <div key={row.label} className="flex items-center justify-between border-b border-border pb-2 text-sm last:border-b-0">
                  <div className="flex items-center gap-2">
                    <span className="h-2.5 w-2.5" style={{ backgroundColor: colors[idx % colors.length] }} />
                    <span className="capitalize text-foreground">{row.label}</span>
                  </div>
                  <span className="tabular-nums text-muted-foreground">{row.value} ({pct}%)</span>
                </div>
              )
            })}
          </div>
        </div>
      )}
    </section>
  )
}

function RankedBars({
  title,
  subtitle,
  rows,
}: {
  title: string
  subtitle: string
  rows: Datum[]
}) {
  const max = useMemo(() => {
    if (!rows.length) return 1
    return Math.max(...rows.map((row) => row.value), 1)
  }, [rows])

  return (
    <section className="border border-border bg-card p-5 sm:p-6">
      <div className="mb-5">
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>

      {rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No data available yet.</p>
      ) : (
        <div className="space-y-3">
          {rows.map((row) => {
            const width = `${Math.max((row.value / max) * 100, 8)}%`
            return (
              <div key={row.label} className="space-y-1.5">
                <div className="flex items-center justify-between gap-4">
                  {row.href ? (
                    <Link href={row.href} className="text-sm font-medium text-foreground hover:underline">
                      {row.label}
                    </Link>
                  ) : (
                    <p className="text-sm font-medium text-foreground">{row.label}</p>
                  )}
                  <p className="shrink-0 text-sm tabular-nums text-muted-foreground">{formatCount(row.value)}</p>
                </div>
                <div className="h-2 w-full bg-muted/60">
                  <div className="h-2 bg-foreground/80" style={{ width }} role="img" aria-label={`${row.label}: ${row.value}`} />
                </div>
              </div>
            )
          })}
        </div>
      )}
    </section>
  )
}

export default function StatsPage() {
  const { data: stats, isLoading: statsLoading } = useStats()
  const { data: plateData, isLoading: platesLoading } = usePlates({ limit: 120, page: 1 })
  const { data: badges } = useBadges()

  const categoryRows = useMemo<Datum[]>(() => {
    const plates = plateData?.data ?? []
    const counts = new Map<string, number>()

    for (const plate of plates) {
      const key = plate.category?.trim() || "uncategorized"
      counts.set(key, (counts.get(key) ?? 0) + 1)
    }

    return Array.from(counts.entries())
      .sort((a, b) => b[1] - a[1])
      .slice(0, 8)
      .map(([label, value]) => ({ label, value, href: `/explore?category=${encodeURIComponent(label)}` }))
  }, [plateData])

  const topPlatesRows = useMemo<Datum[]>(() => {
    const plates = [...(plateData?.data ?? [])]
    return plates
      .sort((a, b) => b.bookmark_count - a.bookmark_count)
      .slice(0, 8)
      .map((plate) => ({
        label: plate.name,
        value: plate.bookmark_count,
        href: buildPlateHref(plate.slug),
      }))
  }, [plateData])

  const badgeTierRows = useMemo<Datum[]>(() => {
    if (!badges?.length) return []

    const counts = new Map<string, number>()
    for (const badge of badges) {
      counts.set(badge.tier, (counts.get(badge.tier) ?? 0) + 1)
    }

    return Array.from(counts.entries())
      .sort((a, b) => b[1] - a[1])
      .map(([label, value]) => ({ label, value }))
  }, [badges])

  const growthRows = useMemo<TrendDatum[]>(() => {
    const plates = plateData?.data ?? []
    const now = new Date()
    const months: { key: string; label: string }[] = []

    for (let i = 7; i >= 0; i -= 1) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`
      const label = d.toLocaleString("en-US", { month: "short" })
      months.push({ key, label })
    }

    const counts = new Map<string, number>()
    for (const month of months) {
      counts.set(month.key, 0)
    }

    for (const plate of plates) {
      const d = new Date(plate.created_at)
      if (Number.isNaN(d.getTime())) continue
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`
      if (!counts.has(key)) continue
      counts.set(key, (counts.get(key) ?? 0) + 1)
    }

    return months.map((month) => ({ label: month.label, value: counts.get(month.key) ?? 0 }))
  }, [plateData])

  const isLoading = statsLoading || platesLoading

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-10 sm:py-12">
        <div className="border border-border bg-muted/20 p-5 sm:p-7">
          <p className="mb-2 text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">Public Dashboard</p>
          <h1 className="text-3xl font-black tracking-tight text-foreground sm:text-4xl">KikPlate Stats</h1>
          <p className="mt-3 max-w-3xl text-sm leading-relaxed text-muted-foreground sm:text-base">
            Live platform metrics and distribution plots from public registry data.
            This page is visible to everyone and updates automatically.
          </p>
        </div>

        <div className="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
          <KpiCard
            label="Templates"
            value={isLoading ? "..." : formatCount(stats?.total_plates ?? 0)}
            icon={<Package className="h-4 w-4" />}
          />
          <KpiCard
            label="Contributors"
            value={isLoading ? "..." : formatCount(stats?.total_contributors ?? 0)}
            icon={<Users className="h-4 w-4" />}
          />
          <KpiCard
            label="Categories"
            value={isLoading ? "..." : formatCount(stats?.total_categories ?? 0)}
            icon={<Layers className="h-4 w-4" />}
          />
          <KpiCard
            label="Bookmarks"
            value={isLoading ? "..." : formatCount(stats?.total_bookmarks ?? 0)}
            icon={<Bookmark className="h-4 w-4" />}
          />
        </div>

        <div className="mt-8 grid grid-cols-1 gap-4 xl:grid-cols-2">
          <LinePlot
            title="Plate Growth (Last 8 Months)"
            subtitle="Monthly new templates based on created_at timestamps"
            points={growthRows}
          />
          <ColumnPlot
            title="Top Categories (Column Chart)"
            subtitle="Category volume among currently loaded public templates"
            rows={categoryRows.slice(0, 6)}
          />
        </div>

        <div className="mt-4 grid grid-cols-1 gap-4 xl:grid-cols-2">
          <DonutPlot
            title="Badge Tier Mix (Donut)"
            subtitle="Share of badge catalog by tier"
            rows={badgeTierRows}
          />
          <RankedBars
            title="Most Bookmarked Plates"
            subtitle="Top templates by total bookmarks"
            rows={topPlatesRows}
          />
        </div>

        <div className="mt-6 flex items-center gap-2 text-xs text-muted-foreground">
          <BarChart3 className="h-3.5 w-3.5" />
          <ShieldCheck className="h-3.5 w-3.5" />
          <span>Data source: public API endpoints (/plates, /plates/stats, /badges)</span>
        </div>
      </div>
    </div>
  )
}