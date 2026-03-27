"use client"

import Link from "next/link"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { GitBranch, Star, Heart, ArrowRight, Sparkles } from "lucide-react"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import type { Plate } from "@/src/domain/entities/Plate"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"

function FeaturedPlateCard({ plate }: { plate: Plate }) {
  return (
    <Link
      href={`/plates/${plate.slug}`}
      className="group flex h-full flex-col gap-3 border border-border bg-card p-5 transition-colors hover:border-foreground/20 hover:bg-muted/20"
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2 text-muted-foreground">
          <GitBranch className="h-4 w-4 shrink-0" />
          <span className="text-xs capitalize">{plate.type}</span>
        </div>
        {plate.avg_rating > 0 && (
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            <Star className="h-3 w-3" />
            {plate.avg_rating.toFixed(1)}
          </div>
        )}
      </div>

      <div>
        <p className="truncate font-semibold text-foreground transition-colors group-hover:text-primary">
          {plate.name}
        </p>
        {plate.description && (
          <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">
            {plate.description}
          </p>
        )}
      </div>

      <div className="mt-auto flex items-center justify-between border-t border-border pt-2">
        <span className="text-xs capitalize text-muted-foreground">{plate.category}</span>
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
    <section className="bg-background py-16">
      <div className="container mx-auto px-4">
        <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Discover
            </p>
            <h2 className="mt-2 text-3xl font-bold tracking-tight text-foreground">Most Used Plates</h2>
            <p className="mt-1 text-sm text-muted-foreground">
              Templates developers repeatedly trust in production.
            </p>
          </div>
          <Link
            href="/explore"
            className="inline-flex items-center gap-1.5 border border-border px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted"
          >
            Explore all
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="h-40 animate-pulse border border-border bg-muted/30" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {plates.map((plate) => (
              <FeaturedPlateCard key={plate.id} plate={plate} />
            ))}
          </div>
        )}

        <div className="mt-6 border border-border bg-card p-5">
          <div className="mb-4 flex items-center gap-2 text-muted-foreground">
            <Sparkles className="h-4 w-4" />
            <p className="text-xs font-semibold uppercase tracking-widest">Contribute</p>
          </div>
          <h3 className="text-xl font-bold text-foreground">Publish your own plate</h3>
          <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
            Share your starter with the community and help other developers ship faster.
          </p>

          <ul className="mt-5 space-y-2 border-t border-border pt-4 text-sm text-muted-foreground">
            <li>Repository templates</li>
            <li>Clear ownership and metadata</li>
            <li>Discoverable in explore and search</li>
          </ul>

          <div className="mt-6 flex flex-wrap gap-2">
            <Link
              href="/submit"
              className={cn(buttonVariants({ size: "sm" }), "gap-1.5")}
            >
              Submit a plate
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}