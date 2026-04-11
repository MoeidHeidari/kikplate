"use client"

import { useState } from "react"
import { Copy, Check, AlertCircle, MoreVertical, RotateCw, Trash2 } from "lucide-react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { Plate } from "@/src/domain/entities/Plate"
import { useVerifyRepository } from "@/src/presentation/hooks/usePlates"

interface Props {
  plate: Plate
  onRemove?: (plate: Plate) => void
  removing?: boolean
}

export function PendingVerification({ plate, onRemove, removing = false }: Props) {
  const [copied, setCopied] = useState(false)
  const verifyMutation = useVerifyRepository()

  if (plate.status !== "pending" || !plate.verification_token) {
    return null
  }

  async function handleCopy() {
    await navigator.clipboard.writeText(plate.verification_token!)
    setCopied(true)
    toast.success("Token copied to clipboard")
    setTimeout(() => setCopied(false), 2000)
  }

  async function handleRetryVerification() {
    try {
      await verifyMutation.mutateAsync(plate.id)
      toast.success("Plate verified and published!")
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Verification failed"
      toast.error(msg)
    }
  }

  const snippetYAML = `verification_token: ${plate.verification_token}`

  return (
    <div className="space-y-4 rounded-lg border border-amber-300/40 bg-amber-50/50 p-4 dark:border-amber-400/30 dark:bg-amber-950/20">
      <div className="flex items-start gap-3">
        <AlertCircle className="mt-0.5 h-5 w-5 text-amber-600 dark:text-amber-400 shrink-0" />
        <div className="min-w-0 flex-1 space-y-3">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <p className="font-semibold text-amber-900 dark:text-amber-100">Verification Pending</p>
              <p className="mt-1 text-sm text-amber-800 dark:text-amber-200">
                To publish this plate, add the verification token to your kikplate.yaml and push it to the repository.
              </p>
            </div>
            {onRemove ? (
              <DropdownMenu>
                <DropdownMenuTrigger className="inline-flex h-7 w-7 shrink-0 items-center justify-center rounded-none border border-transparent text-amber-800 outline-none transition-colors hover:border-amber-400/50 hover:bg-amber-100/60 hover:text-amber-950 focus-visible:border-amber-500 focus-visible:ring-2 focus-visible:ring-amber-500/40 dark:text-amber-200 dark:hover:bg-amber-950/40 dark:hover:text-amber-50">
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">Pending plate actions</span>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem
                    variant="destructive"
                    className="cursor-pointer gap-2"
                    disabled={removing}
                    onClick={() => onRemove(plate)}
                  >
                    <Trash2 className="h-4 w-4" />
                    {removing ? "Removing…" : "Remove plate"}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : null}
          </div>

          <div className="space-y-2">
            <p className="text-sm font-medium text-amber-900 dark:text-amber-100">Add this to your kikplate.yaml:</p>
            <div className="flex gap-2">
              <code className="flex-1 overflow-auto rounded bg-amber-100/50 px-3 py-2 font-mono text-xs text-foreground dark:bg-amber-950/40">
                {snippetYAML}
              </code>
              <Button
                size="sm"
                variant="outline"
                onClick={handleCopy}
                className="shrink-0"
              >
                {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <p className="text-sm text-amber-800 dark:text-amber-200">
              Once you&apos;ve pushed the update, click the button below to verify:
            </p>
            <div className="flex flex-wrap items-center gap-2">
              <Button
                onClick={handleRetryVerification}
                disabled={verifyMutation.isPending}
                className="gap-2"
                size="sm"
              >
                {verifyMutation.isPending && <RotateCw className="h-4 w-4 animate-spin" />}
                {verifyMutation.isPending ? "Verifying..." : "Retry Verification"}
              </Button>
            </div>
          </div>

          {verifyMutation.error && (
            <p className="text-sm text-red-600 dark:text-red-400">
              {verifyMutation.error instanceof Error ? verifyMutation.error.message : "Verification failed"}
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
