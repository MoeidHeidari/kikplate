"use client"

import { useState } from "react"
import { useSearchParams } from "next/navigation"
import { Building2, Plus, Loader2, Pencil, Check, X, Trash2 } from "lucide-react"
import {
  useAcceptOrganizationInvitation,
  useCreateOrganization,
  useDeclineOrganizationInvitation,
  useInviteOrganizationMember,
  useOrganizationInvitations,
  useOrganizationMembers,
  useRemoveOrganizationMember,
  useRevokeOrganizationInvitation,
  useMyOrganizationInvitations,
  useMyOrganizations,
  useLeaveOrganization,
  useRemoveOrganization,
  useUpdateOrganization,
} from "@/src/presentation/hooks/useOrganizations"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import Image from "next/image"
import type { Organization } from "@/src/domain/entities/Organization"
import { CLIENT_API_BASE } from "@/src/lib/client-api"

function OrganizationMembersModal({
  org,
  open,
  onClose,
}: {
  org: Organization | null
  open: boolean
  onClose: () => void
}) {
  const canManage = open && !!org && (org.membership_role === "owner" || org.membership_role === "admin")
  const { data: members, isLoading: membersLoading } = useOrganizationMembers(org?.id ?? "", canManage)
  const { data: invitations, isLoading: invitationsLoading } = useOrganizationInvitations(org?.id ?? "", canManage)
  const inviteMember = useInviteOrganizationMember()
  const removeMember = useRemoveOrganizationMember()
  const revokeInvitation = useRevokeOrganizationInvitation()
  const [inviteEmail, setInviteEmail] = useState("")
  const [inviteRole, setInviteRole] = useState<"admin" | "member">("member")

  if (!open || !org || !canManage) {
    return null
  }

  async function onInvite() {
    if (!org) return
    const email = inviteEmail.trim()
    if (!email) return
    await inviteMember.mutateAsync({ organizationId: org.id, input: { email, role: inviteRole } })
    setInviteEmail("")
  }

  const inviteError = inviteMember.error instanceof Error ? inviteMember.error.message : null
  const removeError = removeMember.error instanceof Error ? removeMember.error.message : null
  const revokeError = revokeInvitation.error instanceof Error ? revokeInvitation.error.message : null

  return (
    <Dialog open={open} onOpenChange={(next) => !next && onClose()}>
      <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-4xl lg:max-w-5xl">
        <div className="flex items-center justify-between gap-3">
          <h3 className="font-heading text-base font-medium">Manage members: {org.name}</h3>
        </div>

        <div className="mt-4 grid gap-2 sm:grid-cols-[1fr_auto_auto]">
          <Input
            type="email"
            value={inviteEmail}
            onChange={(e) => setInviteEmail(e.target.value)}
            placeholder="Invite by email"
          />
          <Select value={inviteRole} onValueChange={(value) => setInviteRole(value as "admin" | "member")}>
            <SelectTrigger>
              <SelectValue placeholder="Role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="member">Member</SelectItem>
              <SelectItem value="admin">Admin</SelectItem>
            </SelectContent>
          </Select>
          <Button type="button" onClick={onInvite} disabled={inviteMember.isPending || !inviteEmail.trim()}>
            {inviteMember.isPending ? "Inviting..." : "Invite"}
          </Button>
        </div>
        {inviteError && <p className="mt-2 text-xs text-destructive">{inviteError}</p>}

        <div className="mt-5 border border-border bg-muted/20 p-3">
          <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-muted-foreground">Members</p>
          {membersLoading && <p className="text-xs text-muted-foreground">Loading members...</p>}
          {!membersLoading && (!members || members.length === 0) && (
            <p className="text-xs text-muted-foreground">No members yet.</p>
          )}
          <div className="space-y-2">
            {members?.map((member) => (
              <div key={member.id} className="flex items-center justify-between gap-3 text-sm">
                <div className="min-w-0">
                  <p className="truncate text-foreground">
                    {member.profile?.display_name || member.profile?.username || member.account_id}
                  </p>
                  <p className="text-xs uppercase tracking-wider text-muted-foreground">
                    {member.role} • {member.status}
                  </p>
                </div>
                {member.role !== "owner" && (
                  <Button
                    type="button"
                    size="sm"
                    variant="destructive"
                    disabled={removeMember.isPending}
                    onClick={() => removeMember.mutate({ organizationId: org.id, accountId: member.account_id })}
                  >
                    Remove
                  </Button>
                )}
              </div>
            ))}
          </div>
          {removeError && <p className="mt-2 text-xs text-destructive">{removeError}</p>}
        </div>

        <div className="mt-4 border border-border bg-muted/20 p-3">
          <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-muted-foreground">Invitations</p>
          {invitationsLoading && <p className="text-xs text-muted-foreground">Loading invitations...</p>}
          {!invitationsLoading && (!invitations || invitations.length === 0) && (
            <p className="text-xs text-muted-foreground">No invitations yet.</p>
          )}
          <div className="space-y-2">
            {invitations?.map((inv) => (
              <div key={inv.id} className="flex items-center justify-between gap-3 text-sm">
                <div className="min-w-0">
                  <p className="truncate text-foreground">{inv.email}</p>
                  <p className="text-xs uppercase tracking-wider text-muted-foreground">
                    {inv.role} • {inv.status}
                  </p>
                </div>
                {inv.status === "pending" && (
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    disabled={revokeInvitation.isPending}
                    onClick={() => revokeInvitation.mutate({ organizationId: org.id, invitationId: inv.id })}
                  >
                    Revoke
                  </Button>
                )}
              </div>
            ))}
          </div>
          {revokeError && <p className="mt-2 text-xs text-destructive">{revokeError}</p>}
        </div>

      </DialogContent>
    </Dialog>
  )
}

export function OrganizationsManager() {
  const searchParams = useSearchParams()
  const { data: config } = useConfig()
  const { data: organizations, isLoading } = useMyOrganizations()
  const createOrg = useCreateOrganization()
  const updateOrg = useUpdateOrganization()
  const removeOrg = useRemoveOrganization()
  const leaveOrg = useLeaveOrganization()
  const { data: invitations, isLoading: invitationsLoading } = useMyOrganizationInvitations()
  const acceptInvitation = useAcceptOrganizationInvitation()
  const declineInvitation = useDeclineOrganizationInvitation()
  
  const privateOrgEnabled = config?.features?.private_organizations_enabled ?? false

  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [logoUrl, setLogoUrl] = useState("")
  const [visibility, setVisibility] = useState<"public" | "private">("public")
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editName, setEditName] = useState("")
  const [editDescription, setEditDescription] = useState("")
  const [editLogoUrl, setEditLogoUrl] = useState("")
  const [editVisibility, setEditVisibility] = useState<"public" | "private">("public")
  const [removingId, setRemovingId] = useState<string | null>(null)
  const [confirmingOrg, setConfirmingOrg] = useState<{ id: string; name: string } | null>(null)
  const [confirmOrgNameInput, setConfirmOrgNameInput] = useState("")
  const [confirmingLeaveOrg, setConfirmingLeaveOrg] = useState<{ id: string; name: string } | null>(null)
  const [leavingId, setLeavingId] = useState<string | null>(null)
  const [managingMembersOrg, setManagingMembersOrg] = useState<Organization | null>(null)
  const [connectingOrgId, setConnectingOrgId] = useState<string | null>(null)

  const githubOrgStatus = searchParams.get("github_org_status")
  const githubOrgId = searchParams.get("github_org_id")

  function connectOrganizationGitHub(orgId: string) {
    setConnectingOrgId(orgId)
    window.location.assign(`${CLIENT_API_BASE}/organizations/${orgId}/github/connect`)
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    try {
      await createOrg.mutateAsync({
        name: name.trim(),
        description: description.trim(),
        logo_url: logoUrl.trim() || undefined,
        visibility,
      })
      setName("")
      setDescription("")
      setLogoUrl("")
      setVisibility("public")
    } catch {
    }
  }

  function startEditing(id: string, currentName: string, currentDescription: string, currentLogoUrl?: string, currentVisibility: "public" | "private" = "public") {
    setEditingId(id)
    setEditName(currentName)
    setEditDescription(currentDescription)
    setEditLogoUrl(currentLogoUrl ?? "")
    setEditVisibility(currentVisibility)
  }

  function cancelEditing() {
    setEditingId(null)
    setEditName("")
    setEditDescription("")
    setEditLogoUrl("")
    setEditVisibility("public")
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
          visibility: editVisibility,
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

  async function onConfirmLeaveOrganization() {
    if (!confirmingLeaveOrg) return
    setLeavingId(confirmingLeaveOrg.id)
    try {
      await leaveOrg.mutateAsync(confirmingLeaveOrg.id)
      if (editingId === confirmingLeaveOrg.id) {
        cancelEditing()
      }
      setConfirmingLeaveOrg(null)
    } finally {
      setLeavingId(null)
    }
  }

  const errorMsg = createOrg.error instanceof Error ? createOrg.error.message : null
  const updateErrorMsg = updateOrg.error instanceof Error ? updateOrg.error.message : null
  const removeErrorMsg = removeOrg.error instanceof Error ? removeOrg.error.message : null
  const leaveErrorMsg = leaveOrg.error instanceof Error ? leaveOrg.error.message : null
  const canRemoveOrganization = Boolean(confirmingOrg && confirmOrgNameInput === confirmingOrg.name)

  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Organizations
        </p>

        <form onSubmit={handleCreate} className="space-y-3 border border-border bg-card rounded-lg p-4">
          <div className="space-y-1.5">
            <Label>Organization name</Label>
            <Input
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="acme-platform"
            />
          </div>

          <div className="space-y-1.5">
            <Label>Description</Label>
            <Textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
              placeholder="What this organization builds"
              className="resize-none"
            />
          </div>

          <div className="space-y-1.5">
            <Label>Logo URL</Label>
            <Input
              type="url"
              value={logoUrl}
              onChange={(e) => setLogoUrl(e.target.value)}
              placeholder="https://example.com/logo.png"
            />
          </div>

          <div className="space-y-1.5">
            <Label>Visibility</Label>
            <Select value={visibility} onValueChange={(value) => setVisibility(value as "public" | "private")}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Visibility" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="public">Public</SelectItem>
                {privateOrgEnabled && <SelectItem value="private">Private</SelectItem>}
              </SelectContent>
            </Select>
            {!privateOrgEnabled && <p className="text-xs text-muted-foreground mt-1">Private organizations are disabled</p>}
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

        {githubOrgStatus && (
          <div className="border border-border bg-muted/20 p-3 text-sm text-muted-foreground rounded-lg">
            {githubOrgStatus === "connected" && "Organization GitHub access connected successfully."}
            {githubOrgStatus === "type-mismatch" && "The GitHub App was installed on a personal account. Organization access requires installation on a GitHub organization."}
            {githubOrgStatus === "error" && "Organization GitHub connection did not complete. Try again from the organization card."}
            {githubOrgId && <span className="block mt-1 text-xs">Organization ID: {githubOrgId}</span>}
          </div>
        )}

        {isLoading && (
          <div className="text-sm text-muted-foreground">Loading organizations...</div>
        )}

        {!isLoading && (!organizations || organizations.length === 0) && (
          <div className="border border-border bg-muted/20 p-4 text-sm text-muted-foreground">
            You have no organizations yet. Create one above, then use it when submitting plates.
          </div>
        )}

        {organizations?.map((org) => (
          <div key={org.id} className="border border-border bg-card rounded-lg p-4">
            {editingId === org.id ? (
              <form onSubmit={handleUpdate} className="space-y-3">
                <div className="space-y-1.5">
                  <Label>Organization name</Label>
                  <Input
                    required
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                  />
                </div>

                <div className="space-y-1.5">
                  <Label>Description</Label>
                  <Textarea
                    value={editDescription}
                    onChange={(e) => setEditDescription(e.target.value)}
                    rows={2}
                    className="resize-none"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label>Logo URL</Label>
                  <Input
                    type="url"
                    value={editLogoUrl}
                    onChange={(e) => setEditLogoUrl(e.target.value)}
                    placeholder="https://example.com/logo.png"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label>Visibility</Label>
                  <Select value={editVisibility} onValueChange={(value) => setEditVisibility(value as "public" | "private")}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Visibility" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="public">Public</SelectItem>
                      {privateOrgEnabled && <SelectItem value="private">Private</SelectItem>}
                    </SelectContent>
                  </Select>
                  {!privateOrgEnabled && <p className="text-xs text-muted-foreground mt-1">Private organizations are disabled</p>}
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
                      <p className="mt-1 text-xs uppercase tracking-wider text-muted-foreground">{org.visibility}</p>
                      <p className="mt-1 text-sm text-muted-foreground">{org.description || "No description"}</p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        GitHub access: {org.github_connected ? `connected${org.github_installation_account ? ` to ${org.github_installation_account}` : ""}` : "not connected"}
                      </p>
                      {org.logo_url && (
                        <p className="mt-1 truncate text-xs text-muted-foreground">{org.logo_url}</p>
                      )}
                    </div>
                    {(org.membership_role === "owner" || org.membership_role === "admin") && (
                      <div className="flex items-center gap-2">
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          className="gap-1.5"
                          disabled={removeOrg.isPending || leaveOrg.isPending || connectingOrgId === org.id}
                          onClick={() => connectOrganizationGitHub(org.id)}
                        >
                          {connectingOrgId === org.id ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : null}
                          {org.github_connected ? "Reconnect GitHub" : "Connect GitHub"}
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          className="gap-1.5"
                          disabled={removeOrg.isPending || leaveOrg.isPending}
                          onClick={() => setManagingMembersOrg(org)}
                        >
                          Manage members
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          className="gap-1.5"
                          disabled={removeOrg.isPending || leaveOrg.isPending}
                          onClick={() => startEditing(org.id, org.name, org.description || "", org.logo_url, org.visibility ?? "public")}
                        >
                          <Pencil className="h-3.5 w-3.5" />
                          Edit
                        </Button>
                      </div>
                    )}

                    {org.membership_role === "member" && (
                      <Button
                        type="button"
                        variant="destructive"
                        size="sm"
                        className="gap-1.5"
                        disabled={leaveOrg.isPending || removeOrg.isPending}
                        onClick={() => {
                          leaveOrg.reset()
                          setConfirmingLeaveOrg({ id: org.id, name: org.name })
                        }}
                      >
                        {leaveOrg.isPending && leavingId === org.id ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <X className="h-3.5 w-3.5" />}
                        Leave
                      </Button>
                    )}
                  </div>

                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      <div className="space-y-3">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Pending invitations</p>
        {invitationsLoading && <div className="text-sm text-muted-foreground">Loading invitations...</div>}
        {!invitationsLoading && (!invitations || invitations.length === 0) && (
          <div className="border border-border bg-muted/20 p-4 text-sm text-muted-foreground">No pending invitations.</div>
        )}
        {invitations?.map((inv) => (
          <div key={inv.id} className="border border-border bg-card rounded-lg p-4">
            <div className="flex items-center justify-between gap-3">
              <div>
                <p className="font-medium text-foreground">Organization invitation</p>
                <p className="text-sm text-muted-foreground">Role: {inv.role}</p>
                <p className="text-xs text-muted-foreground">Expires: {new Date(inv.expires_at).toLocaleDateString()}</p>
              </div>
              <div className="flex items-center gap-2">
                <Button size="sm" onClick={() => acceptInvitation.mutate(inv.id)} disabled={acceptInvitation.isPending}>Accept</Button>
                <Button size="sm" variant="outline" onClick={() => declineInvitation.mutate(inv.id)} disabled={declineInvitation.isPending}>Decline</Button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {confirmingOrg && (
        <Dialog
          open={!!confirmingOrg}
          onOpenChange={(next) => {
            if (!next) {
              removeOrg.reset()
              setConfirmingOrg(null)
              setConfirmOrgNameInput("")
            }
          }}
        >
          <DialogContent className="max-w-sm">
            <DialogHeader>
              <DialogTitle>Delete organization</DialogTitle>
              <DialogDescription>
                This will permanently delete &quot;{confirmingOrg.name}&quot;.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-2">
              <p className="text-xs text-muted-foreground">
                Type the full organization name to confirm: <span className="font-medium text-foreground">{confirmingOrg.name}</span>
              </p>
              <Input
                value={confirmOrgNameInput}
                onChange={(e) => setConfirmOrgNameInput(e.target.value)}
                placeholder="Enter full organization name"
                autoFocus
              />
            </div>

            {removeErrorMsg && (
              <p className="text-xs text-destructive">
                {removeErrorMsg}
                {removeErrorMsg.toLowerCase().includes("contains plates") && " Move plates out of this organization first."}
              </p>
            )}

            <DialogFooter>
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
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}

      {confirmingLeaveOrg && (
        <Dialog
          open={!!confirmingLeaveOrg}
          onOpenChange={(next) => {
            if (!next) {
              leaveOrg.reset()
              setConfirmingLeaveOrg(null)
            }
          }}
        >
          <DialogContent className="max-w-sm">
            <DialogHeader>
              <DialogTitle>Leave organization</DialogTitle>
              <DialogDescription>
                Are you sure you want to leave &quot;{confirmingLeaveOrg.name}&quot;?
              </DialogDescription>
            </DialogHeader>

            {leaveErrorMsg && (
              <p className="text-xs text-destructive">{leaveErrorMsg}</p>
            )}

            <DialogFooter>
              <Button
                variant="outline"
                size="sm"
                disabled={leaveOrg.isPending}
                onClick={() => {
                  leaveOrg.reset()
                  setConfirmingLeaveOrg(null)
                }}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                size="sm"
                disabled={leaveOrg.isPending}
                onClick={onConfirmLeaveOrganization}
              >
                {leaveOrg.isPending ? "Leaving..." : "Yes, leave"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}

      <OrganizationMembersModal
        org={managingMembersOrg}
        open={!!managingMembersOrg}
        onClose={() => setManagingMembersOrg(null)}
      />
    </div>
  )
}
