"use client"

import { useState } from "react"
import { X, Loader2, TriangleAlert } from "lucide-react"

interface Props {
  username: string
  onClose: () => void
  onDeleted: () => void
}

export function DeleteAccountModal({ username, onClose, onDeleted }: Props) {
  const [confirm, setConfirm] = useState("")
  const [deleting, setDeleting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const canDelete = confirm === username

  async function handleDelete() {
    if (!canDelete) return
    setDeleting(true)
    setError(null)
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000))
      onDeleted()
    } catch (error) {
      setError(error instanceof Error ? error.message : "Failed to delete account")
      setDeleting(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" onClick={onClose}>
      <div className="relative w-full max-w-md border border-border bg-background shadow-none" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <h2 className="text-sm font-semibold text-destructive">Delete account</h2>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground transition-colors">
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="px-5 py-5 space-y-4">
          <div className="flex items-start gap-3 p-3 border border-destructive/30 bg-destructive/5">
            <TriangleAlert className="h-4 w-4 text-destructive shrink-0 mt-0.5" />
            <div className="text-xs text-muted-foreground space-y-1">
              <p className="font-medium text-foreground">This action is permanent and cannot be undone.</p>
              <p>All your plates, reviews, and data will be deleted. Other users' forks and uses of your plates will remain.</p>
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs text-muted-foreground">
              Type <span className="font-mono font-medium text-foreground">{username}</span> to confirm
            </label>
            <input
              className="h-9 w-full border border-input bg-transparent px-3 text-sm font-mono outline-none focus:border-destructive transition-colors"
              value={confirm}
              onChange={e => setConfirm(e.target.value)}
              placeholder={username}
              autoComplete="off"
            />
          </div>

          {error && <p className="text-xs text-destructive">{error}</p>}
        </div>

        <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-border">
          <button
            onClick={onClose}
            className="h-9 border border-border px-4 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={!canDelete || deleting}
            className="flex h-9 items-center gap-2 border border-destructive bg-destructive px-4 text-sm text-white transition-colors hover:bg-destructive/90 disabled:opacity-40"
          >
            {deleting && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
            Delete my account
          </button>
        </div>
      </div>
    </div>
  )
}