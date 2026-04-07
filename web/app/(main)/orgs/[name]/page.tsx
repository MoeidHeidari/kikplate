import { notFound } from "next/navigation"
import Link from "next/link"
import Image from "next/image"
import { ArrowLeft, Calendar } from "lucide-react"
import { PlateGrid } from "@/src/presentation/components/plates/PlateGrid"
import { relativeTime } from "@/src/presentation/utils/plateUtils"
import type { Plate, PlateListResponse } from "@/src/domain/entities/Plate"
import type { Organization } from "@/src/domain/entities/Organization"
import { getServerApiBaseUrl } from "@/src/lib/api"

interface Props {
  params: Promise<{ name: string }>
}

export default async function PublicOrgPage({ params }: Props) {
  const { name } = await params
  const base = await getServerApiBaseUrl()

  const orgRes = await fetch(
    `${base}/organizations/by-name/${encodeURIComponent(name)}`,
    { cache: "no-store" },
  )
  if (!orgRes.ok) notFound()
  const org = (await orgRes.json()) as Organization

  const platesRes = await fetch(
    `${base}/plates?organization_id=${org.id}&limit=100`,
    { cache: "no-store" },
  )
  const platesData: PlateListResponse = platesRes.ok
    ? await platesRes.json()
    : { data: [], total: 0 }

  const plates: Plate[] = platesData.data ?? []

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border">
        <div className="container mx-auto px-4 py-6">
          <Link
            href="/explore"
            className="mb-5 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <ArrowLeft className="h-3.5 w-3.5" />
            Back to explore
          </Link>

          <div className="flex items-center gap-4">
            <div className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden border border-border bg-muted">
              {org.logo_url ? (
                <Image
                  src={org.logo_url}
                  alt={org.name}
                  width={64}
                  height={64}
                  unoptimized
                  className="h-full w-full object-cover"
                />
              ) : (
                <span className="text-lg font-bold uppercase text-muted-foreground">
                  {org.name.slice(0, 2)}
                </span>
              )}
            </div>
            <div>
              <h1 className="text-2xl font-black tracking-tight text-foreground">
                {org.name}
              </h1>
              {org.description && (
                <p className="mt-1 text-sm text-muted-foreground">{org.description}</p>
              )}
              <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                <span className="flex items-center gap-1.5">
                  <Calendar className="h-3 w-3" />
                  Created {relativeTime(org.created_at)}
                </span>
                {org.owner && (
                  <span>
                    Owner:{" "}
                    <Link
                      href={`/users/${encodeURIComponent(org.owner.username ?? "")}`}
                      className="font-medium text-foreground hover:underline"
                    >
                      {org.owner.username ?? org.owner.display_name ?? "Unknown"}
                    </Link>
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-10">
        <section>
          <h2 className="mb-4 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">
            Plates ({plates.length})
          </h2>
          <PlateGrid plates={plates} />
        </section>
      </div>
    </div>
  )
}
