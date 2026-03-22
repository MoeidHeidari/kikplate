"use client"

import { useState } from "react"
import { Copy, Check, X, Terminal, GitBranch } from "lucide-react"
import { toast } from "sonner"

interface Props {
  open: boolean
  onClose: () => void
  repoUrl?: string
  slug: string
}

function CopyField({ label, icon, value }: { label: string; icon: React.ReactNode; value: string }) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    toast.success("Copied to clipboard")
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2 text-xs text-muted-foreground">
        {icon}
        <span>{label}</span>
      </div>
      <div className="flex items-center gap-0 border border-border">
        <code className="flex-1 px-3 py-2.5 text-xs font-mono text-foreground bg-muted/30 truncate">
          {value}
        </code>
        <button
          onClick={handleCopy}
          className="px-3 py-2.5 border-l border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
        >
          {copied
            ? <Check className="h-3.5 w-3.5 text-green-500" />
            : <Copy className="h-3.5 w-3.5" />
          }
        </button>
      </div>
    </div>
  )
}

export function UseModal({ open, onClose, repoUrl, slug }: Props) {
  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* modal */}
      <div className="relative w-full max-w-md bg-background border border-border shadow-xl">
        {/* header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <h2 className="text-sm font-semibold">Use this plate</h2>
          <button
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground transition-colors"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* content */}
        <div className="px-5 py-5 space-y-5">
          <p className="text-xs text-muted-foreground">
            Choose how you want to use this template:
          </p>

          {repoUrl && (
            <CopyField
              label="Clone with Git"
              icon={<GitBranch className="h-3.5 w-3.5" />}
              value={`git clone ${repoUrl}`}
            />
          )}

          <CopyField
            label="Scaffold with KikPlate CLI"
            icon={<Terminal className="h-3.5 w-3.5" />}
            value={`kickplate scaf ${slug}`}
          />

          <div className="border-t border-border pt-4">
            <p className="text-xs text-muted-foreground">
              Don't have the CLI?{" "}
              <span className="text-foreground underline underline-offset-4 cursor-pointer">
                Install kickplate
              </span>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}