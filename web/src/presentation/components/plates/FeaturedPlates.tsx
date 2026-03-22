"use client"

import Link from "next/link"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { GitBranch, FileText, Star, Download } from "lucide-react"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import type { Plate } from "@/src/domain/entities/Plate"

function FeaturedPlateCard({ plate }: { plate: Plate }) {
  return (
    <Link
      href={`/plates/${plate.slug}`}
      className="group flex flex-col gap-3 border border-gray-100 bg-white p-5 hover:border-gray-300 hover:shadow-sm transition-all"
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2 text-gray-400">
          {plate.type === "file"
            ? <FileText className="h-4 w-4 shrink-0" />
            : <GitBranch className="h-4 w-4 shrink-0" />
          }
          <span className="text-xs capitalize text-gray-400">{plate.type}</span>
        </div>
        {plate.avg_rating > 0 && (
          <div className="flex items-center gap-1 text-xs text-gray-400">
            <Star className="h-3 w-3" />
            {plate.avg_rating.toFixed(1)}
          </div>
        )}
      </div>

      <div>
        <p className="font-semibold text-gray-900 group-hover:text-gray-600 transition-colors truncate">
          {plate.name}
        </p>
        {plate.description && (
          <p className="text-xs text-gray-400 mt-1 line-clamp-2">
            {plate.description}
          </p>
        )}
      </div>

      <div className="flex items-center justify-between mt-auto pt-2 border-t border-gray-100">
        <span className="text-xs text-gray-400 capitalize">{plate.category}</span>
        <div className="flex items-center gap-1 text-xs text-gray-400">
          <Download className="h-3 w-3" />
          {formatCount(plate.use_count)}
        </div>
      </div>
    </Link>
  )
}

export function FeaturedPlates() {
  const { data, isLoading } = usePlates({ limit: 8 })
  const plates = data?.data ?? []

  return (
    <section className="bg-white border-t border-gray-100 py-16">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Most used</h2>
            <p className="text-sm text-gray-400 mt-1">
              Templates developers reach for most
            </p>
          </div>
          <Link
            href="/explore"
            className="text-sm text-gray-400 hover:text-gray-900 transition-colors"
          >
            View all →
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="h-40 border border-gray-100 bg-gray-50 animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {plates.map((plate) => (
              <FeaturedPlateCard key={plate.id} plate={plate} />
            ))}
          </div>
        )}
      </div>
    </section>
  )
}