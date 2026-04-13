import { PlateGridClient } from "@/src/presentation/components/plates/PlateGridClient"

interface Props {
  searchParams: Promise<{ search?: string; tag?: string; category?: string; badge?: string; type?: string }>
}

export default async function ExplorePage({ searchParams }: Props) {
  const params = await searchParams

  return (
    <div className="container mx-auto px-4 py-8 sm:py-10">
      <header className="mb-8">
        <div className="w-full border border-border bg-card px-5 py-5 sm:px-6 sm:py-6">
          <h1 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">Explore</h1>
          <p className="mt-3 text-[15px] leading-relaxed text-muted-foreground">
            Search templates and refine by category, tags, or badges—filters apply instantly.
          </p>
        </div>
      </header>
      <PlateGridClient
        initialSearch={params.search ?? ""}
        initialTag={params.tag ?? ""}
        initialCategory={params.category ?? ""}
        initialBadge={params.badge ?? ""}
      />
    </div>
  )
}