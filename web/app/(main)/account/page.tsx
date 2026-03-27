"use client"

import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { useMe, useLogout } from "@/src/presentation/hooks/useAuth"
import { usePlates } from "@/src/presentation/hooks/usePlates"
import { useMounted } from "@/src/presentation/hooks/useMounted"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { AccountHeader } from "@/src/presentation/components/account/AccountHeader"
import { AccountStats } from "@/src/presentation/components/account/AccountStats"
import { AccountTabs, type AccountTab } from "@/src/presentation/components/account/AccountTabs"
import { ProfileDetails } from "@/src/presentation/components/account/ProfileDetails"
import { OwnedPlates } from "@/src/presentation/components/account/OwnedPlates"
import { BookmarkedPlates } from "@/src/presentation/components/account/BookmarkedPlates"
import { OrganizationsManager } from "@/src/presentation/components/account/OrganizationsManager"

export default function AccountPage() {
  const mounted = useMounted()
  const { data: me, isLoading } = useMe()
  const logout = useLogout()
  const router = useRouter()
  const searchParams = useSearchParams()

  const resolveTab = (value: string | null): AccountTab => {
    if (value === "plates" || value === "bookmarked" || value === "organizations" || value === "profile") {
      return value
    }
    return "profile"
  }

  const [tab, setTab] = useState<AccountTab>(() => resolveTab(searchParams.get("tab")))

  const { data: ownedData } = usePlates(
    me ? { owner_id: me.account_id, limit: 1000 } : {}
  )

  const totalStars = ownedData?.data?.reduce((sum, plate) => {
    return sum + (plate.avg_rating ? Math.round(plate.avg_rating * 10) / 10 : 0)
  }, 0) ?? 0

  function handleLogout() {
    logout()
    router.push("/")
    router.refresh()
  }

  function handleTabChange(nextTab: AccountTab) {
    setTab(nextTab)
    if (nextTab === "profile") {
      router.replace("/account")
      return
    }
    router.replace(`/account?tab=${nextTab}`)
  }

  if (!mounted || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <LoadingSpinner />
      </div>
    )
  }

  if (!me) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground mb-4">
            You need to sign in to view your account.
          </p>
          <Link href="/login" className="text-sm font-medium underline underline-offset-4">
            Sign in →
          </Link>
        </div>
      </div>
    )
  }

  const stats = [
    { label: "Plates", value: ownedData?.total ?? "—" },
    { label: "Stars",  value: totalStars > 0 ? totalStars.toFixed(1) : "—" },
  ]

  return (
    <div className="min-h-screen bg-background">
      <div className="border-b border-border bg-muted/10">
        <div className="container mx-auto px-4 py-8 space-y-6">
          <AccountHeader me={me} onLogout={handleLogout} />
          <AccountStats stats={stats} />
        </div>
      </div>

      <div className="sticky top-[80px] z-10 border-b border-border bg-background">
        <div className="container mx-auto px-4">
          <AccountTabs active={tab} onChange={handleTabChange} />
        </div>
      </div>

      <div className="container mx-auto px-4 py-10">
        {tab === "profile" && <ProfileDetails me={me} />}
        {tab === "plates"  && <OwnedPlates accountId={me.account_id} />}
        {tab === "bookmarked"    && <BookmarkedPlates accountId={me.account_id} />}
        {tab === "organizations" && <OrganizationsManager />}
      </div>
    </div>
  )
}