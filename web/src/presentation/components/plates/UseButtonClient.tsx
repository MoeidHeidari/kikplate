"use client"

import { useState } from "react"
import { Download } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { UseModal } from "./UseModal"

interface Props {
  plateId: string
  slug: string
  repoUrl?: string
  prominent?: boolean
  className?: string
}

export function UseButtonClient({ plateId, slug, repoUrl, prominent = false, className }: Props) {
  const [openModal, setOpenModal] = useState(false)

  return (
    <>
      <Button
        onClick={() => setOpenModal(true)}
        variant="default"
        className={cn(
          prominent
            ? "h-11 w-full gap-2"
            : "h-9 gap-2 border-border/80 bg-background text-foreground hover:bg-muted",
          className,
        )}
      >
        <Download className="h-3.5 w-3.5" />
        <span className="text-sm font-semibold">Use plate</span>
      </Button>
      <UseModal open={openModal} onClose={() => setOpenModal(false)} repoUrl={repoUrl} slug={slug} />
    </>
  )
}
