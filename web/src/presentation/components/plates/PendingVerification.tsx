"use client"

import { useState } from "react"
import { Copy, Check, AlertCircle, RotateCw, Trash2 } from "lucide-react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
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
        <div className="flex-1 space-y-3">
          <div>
            <p className="font-semibold text-amber-900 dark:text-amber-100">Verification Pending</p>
            <p className="mt-1 text-sm text-amber-800 dark:text-amber-200">
              To publish this plate, add the verification token to your kickplate.yaml and push it to the repository.
            </p>
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
              Once you've pushed the update, click the button below to verify:
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
              {onRemove && (
                <Button
                  type="button"
                  variant="destructive"
                  size="sm"
                  onClick={() => onRemove(plate)}
                  disabled={removing}
                  className="gap-2"
                >
                  <Trash2 className="h-4 w-4" />
                  {removing ? "Removing..." : "Remove Plate"}
                </Button>
              )}
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
