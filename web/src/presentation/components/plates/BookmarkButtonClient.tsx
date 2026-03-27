"use client"

import { useState } from "react"
import { Heart } from "lucide-react"
import { useSetBookmark } from "@/src/presentation/hooks/usePlates"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { toast } from "sonner"

interface Props {
  plateId: string
  isBookmarked?: boolean
  prominent?: boolean
  className?: string
}

export function BookmarkButtonClient({ plateId, isBookmarked = false, prominent = false, className }: Props) {
  const [bookmarked, setBookmarked] = useState(isBookmarked)
  const mutation = useSetBookmark()

  async function handleBookmarkToggle() {
    try {
      await mutation.mutateAsync({ id: plateId, bookmarked: !bookmarked })
      setBookmarked(!bookmarked)
      toast.success(bookmarked ? "Removed from bookmarks" : "Added to bookmarks")
    } catch (error) {
      toast.error("Failed to update bookmark")
    }
  }

  return (
    <Button
      onClick={handleBookmarkToggle}
      disabled={mutation.isPending}
      variant={bookmarked ? "default" : "outline"}
      className={cn(
        prominent
          ? "h-11 w-full gap-2 border-primary/90 bg-primary text-primary-foreground shadow-lg hover:bg-primary/90"
          : "h-9 gap-2 border-border/80 bg-background px-3.5 text-foreground hover:border-primary/40 hover:bg-muted",
        className,
      )}
    >
      <Heart className={cn("h-3.5 w-3.5", bookmarked && "fill-current")} />
      <span className="text-sm font-semibold">{bookmarked ? "Bookmarked" : "Bookmark"}</span>
    </Button>
  )
}
