import type { Plate } from "@/src/domain/entities/Plate"
import { PlateCard } from "./PlateCard"

export function PlateGrid({ plates }: { plates: Plate[] }) {
  if (plates.length === 0) {
    return <p className="text-sm text-muted-foreground py-8 text-center">No plates found.</p>
  }
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {plates.map((p) => <PlateCard key={p.id} plate={p} />)}
    </div>
  )
}
