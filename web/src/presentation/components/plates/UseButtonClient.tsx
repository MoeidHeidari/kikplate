"use client"

import { useState } from "react"
import { Download } from "lucide-react"
import { UseModal } from "./UseModal"

interface Props {
  plateId: string
  slug: string
  repoUrl?: string
}

export function UseButtonClient({ plateId, slug, repoUrl }: Props) {
  const [open, setOpen] = useState(false)

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="w-full flex items-center justify-center gap-2 h-10 bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
      >
        <Download className="h-4 w-4" />
        Use this plate
      </button>

      <UseModal
        open={open}
        onClose={() => setOpen(false)}
        repoUrl={repoUrl}
        slug={slug}
      />
    </>
  )
}