"use client"

import { useState } from "react"
import { Building2, Plus, Loader2, Pencil, Check, X, Trash2 } from "lucide-react"
import { useCreateOrganization, useMyOrganizations, useRemoveOrganization, useUpdateOrganization } from "@/src/presentation/hooks/useOrganizations"
import { Button } from "@/components/ui/button"
import Image from "next/image"

export function OrganizationsManager() {
  const { data: organizations, isLoading } = useMyOrganizations()
  const createOrg = useCreateOrganization()
  const updateOrg = useUpdateOrganization()
  const removeOrg = useRemoveOrganization()

  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [logoUrl, setLogoUrl] = useState("")
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editName, setEditName] = useState("")
  const [editDescription, setEditDescription] = useState("")
  const [editLogoUrl, setEditLogoUrl] = useState("")
  const [removingId, setRemovingId] = useState<string | null>(null)
  const [confirmingOrg, setConfirmingOrg] = useState<{ id: string; name: string } | null>(null)
  const [confirmOrgNameInput, setConfirmOrgNameInput] = useState("")

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    try {
      await createOrg.mutateAsync({
        name: name.trim(),
        description: description.trim(),
        logo_url: logoUrl.trim() || undefined,
      })
      setName("")
      setDescription("")
      setLogoUrl("")
    } catch {
    }
  }

  function startEditing(id: string, currentName: string, currentDescription: string, currentLogoUrl?: string) {
    setEditingId(id)
    setEditName(currentName)
    setEditDescription(currentDescription)
    setEditLogoUrl(currentLogoUrl ?? "")
  }

  function cancelEditing() {
    setEditingId(null)
    setEditName("")
    setEditDescription("")
    setEditLogoUrl("")
  }

  async function handleUpdate(e: React.FormEvent) {
    e.preventDefault()
    if (!editingId) return

    try {
      await updateOrg.mutateAsync({
        id: editingId,
        input: {
          name: editName.trim(),
          description: editDescription.trim(),
          logo_url: editLogoUrl.trim() || "",
        },
      })
      cancelEditing()
    } catch {
    }
  }

  async function onConfirmRemoveOrganization() {
    if (!confirmingOrg) return
    setRemovingId(confirmingOrg.id)
    try {
      await removeOrg.mutateAsync(confirmingOrg.id)
      if (editingId === confirmingOrg.id) {
        cancelEditing()
      }
      setConfirmingOrg(null)
      setConfirmOrgNameInput("")
    } finally {
      setRemovingId(null)
    }
  }

  const errorMsg = createOrg.error instanceof Error ? createOrg.error.message : null
  const updateErrorMsg = updateOrg.error instanceof Error ? updateOrg.error.message : null
  const removeErrorMsg = removeOrg.error instanceof Error ? removeOrg.error.message : null
  const canRemoveOrganization = Boolean(confirmingOrg && confirmOrgNameInput === confirmingOrg.name)

  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Organizations
        </p>

        <form onSubmit={handleCreate} className="space-y-3 border border-border bg-card p-4">
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Organization name</label>
            <input
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="acme-platform"
              className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
              placeholder="What this organization builds"
              className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors resize-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Logo URL</label>
            <input
              type="url"
              value={logoUrl}
              onChange={(e) => setLogoUrl(e.target.value)}
              placeholder="https://example.com/logo.png"
              className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
            />
          </div>

          {errorMsg && (
            <p className="text-sm text-destructive">{errorMsg}</p>
          )}

          <Button type="submit" disabled={createOrg.isPending || !name.trim()} className="gap-2">
            {createOrg.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
            Create organization
          </Button>
        </form>
      </div>

      <div className="space-y-3">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Your organizations
        </p>

        {isLoading && (
          <div className="text-sm text-muted-foreground">Loading organizations...</div>
        )}

        {!isLoading && (!organizations || organizations.length === 0) && (
          <div className="border border-border bg-muted/20 p-4 text-sm text-muted-foreground">
            You have no organizations yet. Create one above, then use it when submitting plates.
          </div>
        )}

        {organizations?.map((org) => (
          <div key={org.id} className="border border-border bg-card p-4">
            {editingId === org.id ? (
              <form onSubmit={handleUpdate} className="space-y-3">
                <div className="space-y-1.5">
                  <label className="text-sm font-medium text-foreground">Organization name</label>
                  <input
                    required
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-sm font-medium text-foreground">Description</label>
                  <textarea
                    value={editDescription}
                    onChange={(e) => setEditDescription(e.target.value)}
                    rows={2}
                    className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors resize-none"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-sm font-medium text-foreground">Logo URL</label>
                  <input
                    type="url"
                    value={editLogoUrl}
                    onChange={(e) => setEditLogoUrl(e.target.value)}
                    placeholder="https://example.com/logo.png"
                    className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
                  />
                </div>

                {updateErrorMsg && (
                  <p className="text-sm text-destructive">{updateErrorMsg}</p>
                )}

                {removeErrorMsg && removingId === editingId && (
                  <p className="text-sm text-destructive">
                    {removeErrorMsg}
                    {removeErrorMsg.toLowerCase().includes("contains plates") && " Move plates out of this organization first."}
                  </p>
                )}

                <div className="flex items-center gap-2">
                  <Button type="submit" disabled={updateOrg.isPending || !editName.trim()} className="gap-2">
                    {updateOrg.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}
                    Save changes
                  </Button>
                  <Button
                    type="button"
                    variant="destructive"
                    onClick={() => {
                      removeOrg.reset()
                      setConfirmingOrg({ id: org.id, name: org.name })
                      setConfirmOrgNameInput("")
                    }}
                    disabled={updateOrg.isPending || removeOrg.isPending}
                    className="gap-2"
                  >
                    {(removeOrg.isPending && removingId === org.id) ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                    Delete
                  </Button>
                  <Button type="button" variant="outline" onClick={cancelEditing} disabled={updateOrg.isPending} className="gap-2">
                    <X className="h-4 w-4" />
                    Cancel
                  </Button>
                </div>
              </form>
            ) : (
              <div className="flex items-start gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden border border-border bg-muted">
                  {org.logo_url ? (
                    <Image
                      src={org.logo_url}
                      alt={`${org.name} logo`}
                      width={40}
                      height={40}
                      unoptimized
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <Building2 className="h-4 w-4 text-muted-foreground" />
                  )}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <p className="font-medium text-foreground">{org.name}</p>
                      <p className="mt-1 text-sm text-muted-foreground">{org.description || "No description"}</p>
                      {org.logo_url && (
                        <p className="mt-1 truncate text-xs text-muted-foreground">{org.logo_url}</p>
                      )}
                    </div>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="gap-1.5"
                      disabled={removeOrg.isPending}
                      onClick={() => startEditing(org.id, org.name, org.description || "", org.logo_url)}
                    >
                      <Pencil className="h-3.5 w-3.5" />
                      Edit
                    </Button>
                  </div>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {confirmingOrg && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-full max-w-sm rounded-none border border-border bg-background p-5">
            <div className="space-y-2">
              <h3 className="font-heading text-base font-medium">Delete organization</h3>
              <p className="text-sm text-muted-foreground">
                This will permanently delete &quot;{confirmingOrg.name}&quot;.
              </p>
              <p className="text-xs text-muted-foreground">
                Type the full organization name to confirm: <span className="font-medium text-foreground">{confirmingOrg.name}</span>
              </p>
              <input
                value={confirmOrgNameInput}
                onChange={(e) => setConfirmOrgNameInput(e.target.value)}
                placeholder="Enter full organization name"
                className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none transition-colors focus:border-ring"
                autoFocus
              />
            </div>

            {removeErrorMsg && (
              <p className="mt-3 text-xs text-destructive">
                {removeErrorMsg}
                {removeErrorMsg.toLowerCase().includes("contains plates") && " Move plates out of this organization first."}
              </p>
            )}

            <div className="mt-5 flex items-center justify-end gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={removeOrg.isPending}
                onClick={() => {
                  removeOrg.reset()
                  setConfirmingOrg(null)
                  setConfirmOrgNameInput("")
                }}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                size="sm"
                disabled={removeOrg.isPending || !canRemoveOrganization}
                onClick={onConfirmRemoveOrganization}
              >
                {removeOrg.isPending ? "Deleting..." : "Yes, delete"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
