"use client"

import Link from "next/link"
import {
  GitBranch, FileText, Download, Star,
  Shield, Tag, CheckCircle2
} from "lucide-react"
import type { Plate } from "@/src/domain/entities/Plate"
import { formatCount, tierColour } from "@/src/presentation/utils/plateUtils"

export function PlateCard({ plate }: { plate: Plate }) {
  return (
    <Link
      href={`/plates/${plate.slug}`}
      className="group flex flex-col gap-0 border border-border bg-card hover:border-foreground/20 hover:shadow-sm transition-all overflow-hidden"
    >
      {/* header */}
      <div className="flex items-start justify-between gap-2 px-4 pt-4 pb-3">
        <div className="flex items-center gap-2 min-w-0">
          <div className="shrink-0 text-muted-foreground">
            {plate.type === "file"
              ? <FileText className="h-4 w-4" />
              : <GitBranch className="h-4 w-4" />
            }
          </div>
          <p className="font-semibold text-sm text-foreground truncate group-hover:text-primary transition-colors">
            {plate.name}
          </p>
        </div>
        {plate.is_verified && (
          <CheckCircle2 className="h-3.5 w-3.5 text-green-500 shrink-0 mt-0.5" />
        )}
      </div>

      {/* description */}
      <div className="px-4 pb-3 min-h-[2.5rem]">
        {plate.description ? (
          <p className="text-xs text-muted-foreground line-clamp-2 leading-relaxed">
            {plate.description}
          </p>
        ) : (
          <p className="text-xs text-muted-foreground/40 italic">No description</p>
        )}
      </div>

      {/* tags */}
      {plate.tags && plate.tags.length > 0 && (
        <div className="px-4 pb-3 flex flex-wrap gap-1">
          {plate.tags.slice(0, 4).map((t) => (
            <span
              key={t.id}
              className="inline-flex items-center gap-1 px-1.5 py-0.5 text-[10px] border border-border text-muted-foreground"
            >
              <Tag className="h-2.5 w-2.5" />
              {t.tag}
            </span>
          ))}
          {plate.tags.length > 4 && (
            <span className="text-[10px] text-muted-foreground/50 py-0.5">
              +{plate.tags.length - 4}
            </span>
          )}
        </div>
      )}

      {/* badges */}
      {plate.badges && plate.badges.length > 0 && (
        <div className="px-4 pb-3 flex flex-wrap gap-1">
          {plate.badges.slice(0, 3).map((pb) => pb.badge && (
            <span
              key={pb.id}
              className={`inline-flex items-center gap-1 px-1.5 py-0.5 text-[10px] border font-medium ${tierColour(pb.badge.tier)}`}
            >
              <Shield className="h-2.5 w-2.5" />
              {pb.badge.name}
            </span>
          ))}
        </div>
      )}

      {/* footer */}
      <div className="mt-auto flex items-center justify-between px-4 py-3 border-t border-border bg-muted/20">
        <span className="text-xs text-muted-foreground capitalize">{plate.category}</span>
        <div className="flex items-center gap-3">
          {plate.avg_rating > 0 && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Star className="h-3 w-3" />
              {plate.avg_rating.toFixed(1)}
            </div>
          )}
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            <Download className="h-3 w-3" />
            {formatCount(plate.use_count)}
          </div>
        </div>
      </div>
    </Link>
  )
}