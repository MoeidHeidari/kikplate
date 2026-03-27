import { PlateGridClient } from "@/src/presentation/components/plates/PlateGridClient"

interface Props {
  searchParams: Promise<{ search?: string; tag?: string; category?: string; type?: string }>
}

export default async function ExplorePage({ searchParams }: Props) {
  const params = await searchParams

  return (
    <div className="container mx-auto px-4 py-10">
      <div className="mb-6 border border-border bg-muted/20 px-5 py-4">
        <h1 className="text-2xl font-bold">Explore</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Browse and discover production-ready templates
        </p>
      </div>
      <PlateGridClient
        initialSearch={params.search ?? ""}
        initialTag={params.tag ?? ""}
        initialCategory={params.category ?? ""}
      />
    </div>
  )
}