"use client"

import { GitBranch, Heart, User, Building2 } from "lucide-react"

export type AccountTab = "profile" | "plates" | "bookmarked" | "organizations"

interface Props {
  active: AccountTab
  onChange: (tab: AccountTab) => void
}

const TABS: { id: AccountTab; label: string; icon: React.ReactNode }[] = [
  { id: "profile", label: "Profile",           icon: <User      className="h-3.5 w-3.5" /> },
  { id: "plates",  label: "My Plates",        icon: <GitBranch className="h-3.5 w-3.5" /> },
  { id: "bookmarked",    label: "Bookmarked",  icon: <Heart     className="h-3.5 w-3.5" /> },
  { id: "organizations", label: "Organizations", icon: <Building2 className="h-3.5 w-3.5" /> },
]

export function AccountTabs({ active, onChange }: Props) {
  return (
    <div className="flex gap-0">
      {TABS.map((t) => (
        <button
          key={t.id}
          onClick={() => onChange(t.id)}
          className={`flex items-center gap-1.5 border-b-2 px-4 py-3 text-sm font-medium transition-colors ${
            active === t.id
              ? "border-foreground text-foreground"
              : "border-transparent text-muted-foreground hover:text-foreground"
          }`}
        >
          {t.icon}
          {t.label}
        </button>
      ))}
    </div>
  )
}