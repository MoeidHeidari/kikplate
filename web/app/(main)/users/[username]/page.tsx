import { notFound } from "next/navigation"
import Link from "next/link"
import Image from "next/image"
import { Calendar, ArrowLeft, Building2 } from "lucide-react"
import { PlateGrid } from "@/src/presentation/components/plates/PlateGrid"
import { relativeTime } from "@/src/presentation/utils/plateUtils"
import type { Plate, PlateListResponse } from "@/src/domain/entities/Plate"
import type { Organization } from "@/src/domain/entities/Organization"
import { getServerApiBaseUrl } from "@/src/lib/api"

interface UserProfile {
  username: string
  display_name?: string
  avatar_url?: string
  account_id: string
  created_at: string
  organizations: Organization[]
}

interface Props {
  params: Promise<{ username: string }>
}

export default async function PublicUserPage({ params }: Props) {
  const { username } = await params
  const base = await getServerApiBaseUrl()

  const userRes = await fetch(`${base}/users/${encodeURIComponent(username)}`, {
    cache: "no-store",
  })
  if (!userRes.ok) notFound()
  const profile = (await userRes.json()) as UserProfile

  const platesRes = await fetch(
    `${base}/plates?owner_id=${profile.account_id}&limit=100`,
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
              {profile.avatar_url ? (
                <Image
                  src={profile.avatar_url}
                  alt={profile.username}
                  width={64}
                  height={64}
                  unoptimized
                  className="h-full w-full object-cover"
                />
              ) : (
                <span className="text-lg font-bold uppercase text-muted-foreground">
                  {(profile.display_name ?? profile.username).slice(0, 2)}
                </span>
              )}
            </div>
            <div>
              <h1 className="text-2xl font-black tracking-tight text-foreground">
                {profile.display_name ?? profile.username}
              </h1>
              <p className="text-sm text-muted-foreground">@{profile.username}</p>
              <div className="mt-1 flex items-center gap-1.5 text-xs text-muted-foreground">
                <Calendar className="h-3 w-3" />
                Joined {relativeTime(profile.created_at)}
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-10">
        {profile.organizations && profile.organizations.length > 0 && (
          <section className="mb-10">
            <h2 className="mb-4 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">
              Organizations
            </h2>
            <div className="flex flex-wrap gap-3">
              {profile.organizations.map((org) => (
                <Link
                  key={org.id}
                  href={`/orgs/${encodeURIComponent(org.name)}`}
                  className="flex items-center gap-2.5 border border-border bg-card px-4 py-3 transition-colors hover:border-foreground/20"
                >
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden border border-border bg-muted">
                    {org.logo_url ? (
                      <Image
                        src={org.logo_url}
                        alt={org.name}
                        width={32}
                        height={32}
                        unoptimized
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <Building2 className="h-4 w-4 text-muted-foreground" />
                    )}
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-foreground">{org.name}</p>
                    {org.description && (
                      <p className="text-xs text-muted-foreground line-clamp-1">
                        {org.description}
                      </p>
                    )}
                  </div>
                </Link>
              ))}
            </div>
          </section>
        )}

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
