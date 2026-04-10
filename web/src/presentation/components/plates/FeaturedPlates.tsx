"use client"

import Link from "next/link"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { GitBranch, Star, Heart, ArrowRight, Sparkles, Terminal, CheckCircle2 } from "lucide-react"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import type { Plate } from "@/src/domain/entities/Plate"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"

function FeaturedPlateCard({ plate, index }: { plate: Plate; index: number }) {
  return (
    <Link
      href={`/plates/${plate.slug}`}
      className="group relative flex h-full flex-col gap-3 border border-border bg-card p-5 transition-all hover:border-foreground/20 hover:bg-muted/20 hover:-translate-y-0.5"
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2 text-muted-foreground">
          <GitBranch className="h-4 w-4 shrink-0" />
          <span className="text-xs capitalize">{plate.type}</span>
        </div>
        <div className="flex items-center gap-2">
          {plate.avg_rating > 0 && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Star className="h-3 w-3" />
              {plate.avg_rating.toFixed(1)}
            </div>
          )}
          <span className="text-[10px] tabular-nums text-muted-foreground/40">{String(index + 1).padStart(2, "0")}</span>
        </div>
      </div>

      <div>
        <p className="truncate font-semibold text-foreground transition-colors group-hover:text-primary">
          {plate.name}
        </p>
        {plate.description && (
          <p className="mt-1.5 line-clamp-2 text-xs leading-relaxed text-muted-foreground">
            {plate.description}
          </p>
        )}
      </div>

      <div className="mt-auto flex items-center justify-between border-t border-border pt-3">
        <span className="inline-flex items-center gap-1.5 text-xs capitalize text-muted-foreground">
          {plate.category}
        </span>
        <div className="flex items-center gap-1 text-xs text-muted-foreground">
          <Heart className="h-3 w-3" />
          {formatCount(plate.bookmark_count)}
        </div>
      </div>
    </Link>
  )
}

export function FeaturedPlates() {
  const { data, isLoading } = usePlates({ limit: 8 })
  const plates = data?.data ?? []

  return (
    <section className="bg-background py-20">
      <div className="container mx-auto px-4">

        <div className="mb-10 flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Discover
            </p>
            <h2 className="mt-2 text-3xl font-bold tracking-tight text-foreground">Most Used Plates</h2>
            <p className="mt-2 max-w-md text-sm leading-relaxed text-muted-foreground">
              Templates developers repeatedly trust in production.
            </p>
          </div>
          <Link
            href="/explore"
            className="group inline-flex items-center gap-1.5 border border-border px-4 py-2 text-sm text-foreground transition-colors hover:bg-muted"
          >
            Explore all
            <ArrowRight className="h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" />
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="h-44 animate-pulse border border-border bg-muted/20" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {plates.map((plate, i) => (
              <FeaturedPlateCard key={plate.id} plate={plate} index={i} />
            ))}
          </div>
        )}

        <div className="mt-8 grid grid-cols-1 gap-3 lg:grid-cols-2">
          <div className="border border-border bg-card p-6">
            <div className="mb-4 flex items-center gap-2 text-muted-foreground">
              <Sparkles className="h-4 w-4" />
              <p className="text-xs font-semibold uppercase tracking-widest">Contribute</p>
            </div>
            <h3 className="text-xl font-bold text-foreground">Publish your own plate</h3>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Share your starter with the community and help other developers ship faster.
            </p>
            <div className="mt-5 space-y-2 border-t border-border pt-4 text-sm text-muted-foreground">
              <p className="flex items-center gap-2"><CheckCircle2 className="h-3.5 w-3.5" />Repository templates with clear ownership</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-3.5 w-3.5" />Discoverable in explore and search</p>
              <p className="flex items-center gap-2"><CheckCircle2 className="h-3.5 w-3.5" />Community reviews and badges</p>
            </div>
            <div className="mt-6">
              <Link
                href="/submit"
                className={cn(buttonVariants({ size: "sm" }), "gap-1.5")}
              >
                Submit a plate
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>

          <div className="border border-border bg-card p-6">
            <div className="mb-4 flex items-center gap-2 text-muted-foreground">
              <Terminal className="h-4 w-4" />
              <p className="text-xs font-semibold uppercase tracking-widest">Quick start</p>
            </div>
            <h3 className="text-xl font-bold text-foreground">Scaffold in seconds</h3>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Use the KikPlate CLI to search, scaffold and manage plates directly from your terminal.
            </p>
            <div className="mt-5 border border-border bg-background p-4 font-mono text-sm text-muted-foreground">
              <p><span className="text-foreground/50">$</span> kikplate search golang</p>
              <p className="mt-1"><span className="text-foreground/50">$</span> kikplate scaf go-clean-arch</p>
              <p className="mt-1"><span className="text-foreground/50">$</span> kikplate describe my-plate</p>
            </div>
            <div className="mt-6">
              <Link
                href="https://github.com/kikplate/kikplate/tree/main/cli"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 border border-border px-3 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-muted"
              >
                Install CLI
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>
        </div>

      </div>
    </section>
  )
}