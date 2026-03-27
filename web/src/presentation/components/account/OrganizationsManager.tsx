"use client"

import { useState } from "react"
import { Building2, Plus, Loader2, Pencil, Check, X } from "lucide-react"
import { useCreateOrganization, useMyOrganizations, useUpdateOrganization } from "@/src/presentation/hooks/useOrganizations"
import { Button } from "@/components/ui/button"

export function OrganizationsManager() {
  const { data: organizations, isLoading } = useMyOrganizations()
  const createOrg = useCreateOrganization()
  const updateOrg = useUpdateOrganization()

  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [logoUrl, setLogoUrl] = useState("")
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editName, setEditName] = useState("")
  const [editDescription, setEditDescription] = useState("")
  const [editLogoUrl, setEditLogoUrl] = useState("")

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

  const errorMsg = createOrg.error instanceof Error ? createOrg.error.message : null
  const updateErrorMsg = updateOrg.error instanceof Error ? updateOrg.error.message : null

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

                <div className="flex items-center gap-2">
                  <Button type="submit" disabled={updateOrg.isPending || !editName.trim()} className="gap-2">
                    {updateOrg.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}
                    Save changes
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
                    <img src={org.logo_url} alt={`${org.name} logo`} className="h-full w-full object-cover" />
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
    </div>
  )
}
