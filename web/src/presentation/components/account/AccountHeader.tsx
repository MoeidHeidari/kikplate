"use client"

import Image from "next/image"
import { LogOut } from "lucide-react"
import type { MeResult } from "@/src/domain/entities/User"

interface Props {
  me: MeResult
  onLogout: () => void
}

export function AccountHeader({ me, onLogout }: Props) {
  const displayName = me.username ?? me.display_name ?? "User"
  const initials = displayName.slice(0, 2).toUpperCase()

  return (
    <div className="flex items-start justify-between gap-6">
      <div className="flex items-center gap-5">
        <div className="relative flex h-16 w-16 shrink-0 items-center justify-center border border-border bg-card text-xl font-bold text-foreground overflow-hidden">
          {me.avatar_url ? (
            <Image
              src={me.avatar_url}
              alt={displayName}
              width={64}
              height={64}
              unoptimized
              className="h-full w-full object-cover"
            />
          ) : (
            initials
          )}
        </div>

        <div>
          <h1 className="text-2xl font-bold text-foreground">{displayName}</h1>
          <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            {me.email && <span>{me.email}</span>}
            {me.email && <span>·</span>}
            <span className="capitalize">{me.provider}</span>
            {me.role && (
              <>
                <span>·</span>
                <span className="capitalize">{me.role}</span>
              </>
            )}
          </div>
        </div>
      </div>

      <button
        onClick={onLogout}
        className="flex shrink-0 items-center gap-1.5 border border-border px-3 py-2 text-xs text-muted-foreground transition-colors hover:border-destructive/50 hover:text-destructive"
      >
        <LogOut className="h-3.5 w-3.5" />
        Sign out
      </button>
    </div>
  )
}