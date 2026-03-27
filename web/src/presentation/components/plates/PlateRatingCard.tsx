"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Star, CheckCircle2 } from "lucide-react"
import { toast } from "sonner"
import { ApiError } from "@/src/data/repositories/httpClient"
import { AuthService } from "@/src/domain/services/AuthService"
import { useMe } from "@/src/presentation/hooks/useAuth"
import { useRatePlate } from "@/src/presentation/hooks/usePlates"
import { useMounted } from "@/src/presentation/hooks/useMounted"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface Props {
  plateId: string
  plateSlug: string
  plateOwnerId: string
  avgRating: number
  userRating?: number
}

export function PlateRatingCard({ plateId, plateSlug, plateOwnerId, avgRating, userRating: initialUserRating }: Props) {
  const router = useRouter()
  const mounted = useMounted()
  const { data: me } = useMe()
  const isOwnPlate = mounted && me?.account_id === plateOwnerId
  const isAuthenticated = mounted && AuthService.isAuthenticated()
  const [rating, setRating] = useState(0)
  const [userRating, setUserRating] = useState<number | null>(initialUserRating ?? null)
  const [communityRating, setCommunityRating] = useState(avgRating)
  const mutation = useRatePlate()

  useEffect(() => {
    setUserRating(initialUserRating ?? null)
  }, [initialUserRating])

  useEffect(() => {
    setCommunityRating(avgRating)
  }, [avgRating])

  async function submitRating() {
    if (!rating) return

    try {
      console.log("Submitting rating:", { plateId, rating })
      await mutation.mutateAsync({ id: plateId, slug: plateSlug, rating })
      setUserRating(rating)
      setRating(0)
      toast.success("Thanks for rating this plate!")
      router.refresh()
    } catch (error) {
      console.error("Rating submission error:", error)
      if (error instanceof ApiError) {
        if (error.isConflict()) {
          setUserRating(rating)
          toast.info("You've already rated this plate")
          return
        }
        if (error.status === 403) {
          toast.error("You cannot rate your own plate")
          return
        }
        if (error.status === 400) {
          toast.error(error.message || "Invalid rating")
          return
        }
      }
      toast.error("Could not submit your rating")
    }
  }

  return (
    <div className="border border-border bg-card p-5">
      <p className="mb-3 text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">Your Rating</p>

      {isOwnPlate && (
        <div className="mb-4 rounded-lg bg-blue-500/5 p-3 border border-blue-500/20">
          <p className="text-xs text-blue-700 dark:text-blue-400">You cannot rate your own plate</p>
        </div>
      )}

      {!isOwnPlate && userRating ? (
        <div className="mb-4 rounded-lg bg-amber-500/5 p-3 border border-amber-500/20">
          <div className="flex items-center gap-2 mb-2">
            <CheckCircle2 className="h-4 w-4 text-amber-600 dark:text-amber-500" />
            <p className="text-xs font-medium text-amber-700 dark:text-amber-400">You rated this {userRating} star{userRating > 1 ? 's' : ''}</p>
          </div>
          <button
            onClick={() => setUserRating(null)}
            className="text-xs text-amber-600 dark:text-amber-400 hover:underline"
          >
            Change rating
          </button>
        </div>
      ) : !isOwnPlate ? (
        <>
          <div className="mb-3 flex items-center gap-2">
            {Array.from({ length: 5 }).map((_, index) => {
              const value = index + 1
              const active = value <= rating

              return (
                <button
                  key={value}
                  type="button"
                  onClick={() => setRating(value)}
                  disabled={!isAuthenticated || mutation.isPending}
                  className={cn(
                    "inline-flex h-8 w-8 items-center justify-center border transition-colors",
                    active
                      ? "border-amber-500/50 bg-amber-500/10 text-amber-500"
                      : "border-border bg-background text-muted-foreground hover:text-foreground",
                    (!isAuthenticated || mutation.isPending) && "cursor-not-allowed opacity-60",
                  )}
                  aria-label={`Rate ${value} star${value > 1 ? "s" : ""}`}
                >
                  <Star className={cn("h-4 w-4", active && "fill-current")} />
                </button>
              )
            })}
          </div>

          {isAuthenticated && (
            <Button
              onClick={submitRating}
              disabled={!rating || mutation.isPending}
              size="sm"
              className="w-full mb-3"
            >
              {mutation.isPending ? "Submitting..." : "Submit rating"}
            </Button>
          )}
        </>
      ) : null}

      <p className="text-xs text-muted-foreground">
        Community rating: {communityRating > 0 ? communityRating.toFixed(1) : "Not rated yet"}
      </p>

      {isAuthenticated && !isOwnPlate && !userRating && me?.username && (
        <p className="mt-3 text-xs text-muted-foreground">Signed in as <span className="text-foreground">@{me.username}</span></p>
      )}

      {!isAuthenticated && (
        <p className="mt-3 text-xs text-muted-foreground">
          <Link href="/login" className="text-foreground underline underline-offset-4">Log in</Link> to rate this plate.
        </p>
      )}
    </div>
  )
}