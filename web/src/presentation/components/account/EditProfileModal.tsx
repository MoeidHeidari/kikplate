"use client"

import { useState } from "react"
import { X, Loader2 } from "lucide-react"
import type { MeResult } from "@/src/domain/entities/User"
import { Button } from "@/components/ui/button"
import { useSetUsername, useUpdateProfile } from "@/src/presentation/hooks/useAuth"

interface Props {
  me: MeResult
  onClose: () => void
  onSaved: () => void
}

export function EditProfileModal({ me, onClose, onSaved }: Props) {
  const [username, setUsername] = useState(me.username ?? "")
  const [displayName, setDisplayName] = useState(me.display_name ?? "")
  const [avatarUrl, setAvatarUrl] = useState(me.avatar_url ?? "")
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const setUsernameMutation = useSetUsername()
  const updateProfileMutation = useUpdateProfile()

  async function handleSave() {
    setSaving(true)
    setError(null)
    try {
      if (me.provider === "local" && me.username && username.trim() && username.trim() !== me.username) {
        await setUsernameMutation.mutateAsync(username.trim())
      }

      await updateProfileMutation.mutateAsync({
        display_name: displayName.trim() || undefined,
        avatar_url: avatarUrl.trim() || undefined,
      })

      onSaved()
    } catch (e) {
      const message = e instanceof Error ? e.message : "Failed to save"
      setError(message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" onClick={onClose}>
      <div className="relative w-full max-w-md border border-border bg-background shadow-none" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <h2 className="text-sm font-semibold">Edit profile</h2>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground transition-colors">
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="px-5 py-5 space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs text-muted-foreground">Username</label>
            <input
              className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none focus:border-ring transition-colors"
              value={username}
              onChange={e => setUsername(e.target.value)}
              placeholder="your_username"
              disabled={me.provider !== "local"}
            />
            {me.provider !== "local" && (
              <p className="text-xs text-muted-foreground/60">
                Username is managed by your {me.provider} account.
              </p>
            )}
          </div>

          <div className="space-y-1.5">
            <label className="text-xs text-muted-foreground">Display name</label>
            <input
              className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none focus:border-ring transition-colors"
              value={displayName}
              onChange={e => setDisplayName(e.target.value)}
              placeholder="Your Name"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs text-muted-foreground">Avatar URL</label>
            <input
              className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none focus:border-ring transition-colors"
              value={avatarUrl}
              onChange={e => setAvatarUrl(e.target.value)}
              placeholder="https://..."
            />
          </div>

          {avatarUrl && (
            <div className="flex items-center gap-3">
              <img
                src={avatarUrl}
                alt="Preview"
                className="h-10 w-10 object-cover border border-border"
                onError={e => (e.currentTarget.style.display = "none")}
              />
              <span className="text-xs text-muted-foreground">Avatar preview</span>
            </div>
          )}

          {error && <p className="text-xs text-destructive">{error}</p>}
        </div>

        <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-border">
          <button
            onClick={onClose}
            className="h-9 border border-border px-4 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            Cancel
          </button>
          <Button
            onClick={handleSave}
            disabled={saving}
            size="sm"
            className="gap-2"
          >
            {saving && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
            Save changes
          </Button>
        </div>
      </div>
    </div>
  )
}