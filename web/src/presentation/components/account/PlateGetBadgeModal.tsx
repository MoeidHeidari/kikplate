"use client"

import { useState } from "react"
import type { Plate } from "@/src/domain/entities/Plate"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Check, Copy } from "lucide-react"
import { toast } from "sonner"

function buildBadgeSnippets(plate: Plate) {
  const origin = typeof window !== "undefined" ? window.location.origin : ""
  const platePath = `/plates/${encodeURIComponent(plate.slug)}`
  const plateUrl = `${origin}${platePath}`
  // Use static/v1 query params so slugs with hyphens are not parsed as extra /badge/ path segments (which yields shields "404 badge not found").
  const qs = new URLSearchParams({
    label: "KikPlate",
    message: plate.slug,
    color: "0366d6",
    style: "flat-square",
  })
  const shieldUrl = `https://img.shields.io/static/v1?${qs.toString()}`
  const markdown = `[![View on KikPlate](${shieldUrl})](${plateUrl})`
  const asciidoc = `image:${shieldUrl}[View on KikPlate,link="${plateUrl}"]`
  return { plateUrl, shieldUrl, markdown, asciidoc }
}

function CopyRow({
  label,
  value,
  copyId,
  activeCopy,
  onCopy,
}: {
  label: string
  value: string
  copyId: string
  activeCopy: string | null
  onCopy: (id: string, text: string) => void
}) {
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground">{label}</span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="h-7 gap-1 px-2 text-xs"
          onClick={() => onCopy(copyId, value)}
        >
          {activeCopy === copyId ? (
            <Check className="h-3 w-3 text-green-600" />
          ) : (
            <Copy className="h-3 w-3" />
          )}
          Copy
        </Button>
      </div>
      <pre className="max-h-36 overflow-auto whitespace-pre-wrap break-all border border-border bg-muted/30 p-3 font-mono text-[11px] leading-relaxed text-foreground">
        {value}
      </pre>
    </div>
  )
}

export function PlateGetBadgeModal({
  plate,
  onClose,
}: {
  plate: Plate | null
  onClose: () => void
}) {
  const [activeCopy, setActiveCopy] = useState<string | null>(null)

  async function handleCopy(id: string, text: string) {
    try {
      await navigator.clipboard.writeText(text)
      setActiveCopy(id)
      toast.success("Copied to clipboard")
      setTimeout(() => setActiveCopy(null), 2000)
    } catch {
      toast.error("Could not copy")
    }
  }

  const snippets = plate ? buildBadgeSnippets(plate) : null

  return (
    <Dialog open={Boolean(plate)} onOpenChange={(open) => { if (!open) onClose() }}>
      <DialogContent className="max-h-[min(90vh,640px)] overflow-y-auto sm:max-w-lg">
        {snippets && plate ? (
          <>
            <DialogHeader>
              <DialogTitle>Get badge</DialogTitle>
              <DialogDescription>
                Preview how the badge looks, then copy Markdown or AsciiDoc for your repository.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-2">
              <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Preview</p>
              <div className="flex flex-wrap items-center gap-3 rounded-md border border-border bg-muted/20 px-4 py-4">
                <a
                  href={snippets.plateUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center overflow-hidden rounded-md"
                >
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img src={snippets.shieldUrl} alt="View on KikPlate" className="block max-w-full" />
                </a>
              </div>
            </div>

            <div className="space-y-4 border-t border-border pt-4">
              <CopyRow
                label="Markdown"
                value={snippets.markdown}
                copyId="md"
                activeCopy={activeCopy}
                onCopy={handleCopy}
              />
              <CopyRow
                label="AsciiDoc"
                value={snippets.asciidoc}
                copyId="adoc"
                activeCopy={activeCopy}
                onCopy={handleCopy}
              />
            </div>
          </>
        ) : null}
      </DialogContent>
    </Dialog>
  )
}
