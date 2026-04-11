"use client"

import Link from "next/link"
import { useState } from "react"
import { usePlates, useRemovePlate, useMovePlateOrganization } from "@/src/presentation/hooks/usePlates"
import { useMyOrganizations } from "@/src/presentation/hooks/useOrganizations"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { PendingVerification } from "@/src/presentation/components/plates/PendingVerification"
import { PlateGetBadgeModal } from "@/src/presentation/components/account/PlateGetBadgeModal"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { Plate } from "@/src/domain/entities/Plate"
import {
  Building2,
  Code2,
  Copy,
  GitBranch,
  MoreVertical,
  Plus,
  Trash2,
} from "lucide-react"
import { toast } from "sonner"

export function OwnedPlates({ accountId }: { accountId: string }) {
  const { data, isLoading, isError } = usePlates({ owner_id: accountId, limit: 48 })
  const { data: organizations = [], isLoading: orgsLoading } = useMyOrganizations()
  const removePlate = useRemovePlate()
  const movePlateOrganization = useMovePlateOrganization()
  const [removingId, setRemovingId] = useState<string | null>(null)
  const [movingId, setMovingId] = useState<string | null>(null)
  const [selectedOrgByPlate, setSelectedOrgByPlate] = useState<Record<string, string>>({})
  const [confirmingMove, setConfirmingMove] = useState<{
    plateId: string
    plateName: string
    targetOrgId: string
    targetOrgName: string
  } | null>(null)
  const [moveError, setMoveError] = useState<string | null>(null)
  const [confirmingPlate, setConfirmingPlate] = useState<{ id: string; name: string } | null>(null)
  const [removeError, setRemoveError] = useState<string | null>(null)
  const [confirmNameInput, setConfirmNameInput] = useState("")
  const [moveConfirmNameInput, setMoveConfirmNameInput] = useState("")
  const [getBadgePlate, setGetBadgePlate] = useState<Plate | null>(null)
  const [orgDialogPlate, setOrgDialogPlate] = useState<Plate | null>(null)

  const canRemove = Boolean(confirmingPlate && confirmNameInput === confirmingPlate.name)
  const canMove = Boolean(confirmingMove && moveConfirmNameInput === confirmingMove.plateName)

  const onConfirmRemove = async () => {
    if (!confirmingPlate) return

    try {
      setRemoveError(null)
      setRemovingId(confirmingPlate.id)
      await removePlate.mutateAsync(confirmingPlate.id)
      setConfirmingPlate(null)
    } catch (error) {
      setRemoveError(error instanceof Error ? error.message : "Failed to remove plate")
    } finally {
      setRemovingId(null)
    }
  }

  const onConfirmMoveOrganization = async () => {
    if (!confirmingMove) return

    const organizationId = confirmingMove.targetOrgId !== "" ? confirmingMove.targetOrgId : undefined

    try {
      setMoveError(null)
      setMovingId(confirmingMove.plateId)
      const updated = await movePlateOrganization.mutateAsync({ id: confirmingMove.plateId, organizationId })
      setSelectedOrgByPlate((prev) => ({
        ...prev,
        [confirmingMove.plateId]: updated.organization_id ?? "",
      }))
      setConfirmingMove(null)
      setMoveConfirmNameInput("")
      setOrgDialogPlate((p) => (p?.id === confirmingMove.plateId ? null : p))
    } catch (error) {
      setMoveError(error instanceof Error ? error.message : "Failed to move plate")
    } finally {
      setMovingId(null)
    }
  }

  if (isLoading) return <LoadingSpinner />

  if (isError) {
    return <p className="text-sm text-destructive">Failed to load your plates.</p>
  }

  const plates = data?.data ?? []
  const pendingPlates = plates.filter((p) => p.status === "pending")
  const publishedPlates = plates.filter((p) => p.status !== "pending")

  const visibilityBadgeClass = (visibility: string) => {
    switch (visibility) {
      case "private":
        return "border-amber-500/40 bg-amber-500/10 text-amber-700 dark:text-amber-300"
      case "unlisted":
        return "border-blue-500/40 bg-blue-500/10 text-blue-700 dark:text-blue-300"
      default:
        return "border-emerald-500/40 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300"
    }
  }

  return (
    <div className="space-y-6">
      {pendingPlates.length > 0 && (
        <div className="space-y-3">
          <p className="text-xs font-medium uppercase tracking-widest text-muted-foreground">
            Pending Verification ({pendingPlates.length})
          </p>
          <div className="space-y-3">
            {pendingPlates.map((plate) => (
              <PendingVerification
                key={plate.id}
                plate={plate}
                removing={removePlate.isPending && removingId === plate.id}
                onRemove={(targetPlate) => {
                  setConfirmingPlate({ id: targetPlate.id, name: targetPlate.name })
                  setConfirmNameInput("")
                  setRemoveError(null)
                }}
              />
            ))}
          </div>
        </div>
      )}

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <p className="text-xs text-muted-foreground">
            {publishedPlates.length} published plate{publishedPlates.length !== 1 ? "s" : ""}
          </p>
          <Link
            href="/submit"
            className="flex items-center gap-1.5 border border-border px-3 py-1.5 text-xs font-medium text-foreground transition-colors hover:bg-muted"
          >
            <Plus className="h-3.5 w-3.5" />
            New plate
          </Link>
        </div>

        {publishedPlates.length === 0 ? (
          <div className="flex flex-col items-center gap-3 py-20 text-muted-foreground">
            <GitBranch className="h-8 w-8 opacity-30" />
            <p className="text-sm">You haven&apos;t published any plates yet.</p>
            <Link
              href="/submit"
              className="text-sm text-foreground underline underline-offset-4"
            >
              Submit your first plate →
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {publishedPlates.map((plate) => {
              const organizationLabel =
                plate.organization?.name ??
                (plate.organization_id
                  ? organizations.find((o) => o.id === plate.organization_id)?.name
                  : undefined)

              return (
            <div key={plate.id} className="group relative rounded-xl border border-border bg-card p-4">
              <Link
                href={`/plates/${encodeURIComponent(plate.slug)}`}
                className="absolute inset-0 z-0 rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 ring-offset-background"
                aria-label={`Open plate: ${plate.name}`}
              />
              <div className="relative z-10 space-y-2 pointer-events-none">
                <div className="flex items-start gap-2">
                  <span className="min-w-0 flex-1 line-clamp-2 text-sm font-semibold leading-snug group-hover:underline">
                    {plate.name}
                  </span>
                  <div className="flex shrink-0 flex-wrap items-center justify-end gap-1">
                    <span className={`inline-flex items-center border px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide ${visibilityBadgeClass(plate.visibility)}`}>
                      {plate.visibility}
                    </span>
                    {plate.sync_status === "failed" && (
                      <span className="inline-flex items-center border border-red-500/40 bg-red-500/10 px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide text-red-700 dark:text-red-300">
                        Sync Failed
                      </span>
                    )}
                    {!plate.is_verified && (
                      <span className="inline-flex items-center border border-amber-500/40 bg-amber-500/10 px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide text-amber-700 dark:text-amber-300">
                        Unverified
                      </span>
                    )}
                    <div className="pointer-events-auto">
                    <DropdownMenu>
                      <DropdownMenuTrigger className="inline-flex h-7 w-7 shrink-0 items-center justify-center rounded-none border border-transparent text-muted-foreground outline-none transition-colors hover:border-border hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50">
                        <MoreVertical className="h-4 w-4" />
                        <span className="sr-only">Plate actions</span>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          className="cursor-pointer gap-2"
                          onClick={() => setGetBadgePlate(plate)}
                        >
                          <Code2 className="h-4 w-4" />
                          Get badge
                        </DropdownMenuItem>
                        {plate.verification_token ? (
                          <DropdownMenuItem
                            className="cursor-pointer gap-2"
                            onClick={async () => {
                              try {
                                await navigator.clipboard.writeText(plate.verification_token!)
                                toast.success("Verification token copied")
                              } catch {
                                toast.error("Could not copy")
                              }
                            }}
                          >
                            <Copy className="h-4 w-4" />
                            Copy verification token
                          </DropdownMenuItem>
                        ) : null}
                        <DropdownMenuItem
                          className="cursor-pointer gap-2"
                          onClick={() => setOrgDialogPlate(plate)}
                        >
                          <Building2 className="h-4 w-4" />
                          Organization…
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          variant="destructive"
                          className="cursor-pointer gap-2"
                          disabled={removePlate.isPending && removingId === plate.id}
                          onClick={() => {
                            setConfirmingPlate({ id: plate.id, name: plate.name })
                            setConfirmNameInput("")
                            setRemoveError(null)
                          }}
                        >
                          <Trash2 className="h-4 w-4" />
                          {removePlate.isPending && removingId === plate.id ? "Removing…" : "Remove"}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                    </div>
                  </div>
                </div>
                {plate.description ? (
                  <p className="line-clamp-2 text-xs text-muted-foreground">{plate.description}</p>
                ) : (
                  <p className="text-xs text-muted-foreground">No description</p>
                )}
                {(plate.last_synced_at || plate.next_sync_at) && (
                  <div className="space-y-1 text-[11px]">
                    {plate.last_synced_at && (
                      <p className="text-muted-foreground">
                        Last synced: <span className="font-medium text-foreground">{new Date(plate.last_synced_at).toLocaleString()}</span>
                      </p>
                    )}
                    {plate.next_sync_at && (
                      <p className="text-muted-foreground">
                        Next sync: <span className="font-medium text-foreground">{new Date(plate.next_sync_at).toLocaleString()}</span>
                      </p>
                    )}
                  </div>
                )}
                {(plate.sync_status === "failed" || plate.sync_status === "unverified") && plate.sync_error && (
                  <p className="text-xs text-red-600 dark:text-red-400 whitespace-pre-wrap break-words">{plate.sync_error}</p>
                )}
              </div>

              <p className="relative z-10 mt-4 text-[11px] uppercase tracking-wide text-muted-foreground pointer-events-none">
                {plate.category}
              </p>
              {organizationLabel ? (
                <p className="relative z-10 mt-1.5 flex items-center gap-1.5 text-xs text-muted-foreground pointer-events-none">
                  <Building2 className="h-3.5 w-3.5 shrink-0 opacity-70" aria-hidden />
                  <span className="min-w-0 truncate font-medium text-foreground/90">{organizationLabel}</span>
                </p>
              ) : null}
            </div>
              )
            })}
        </div>
      )}
      </div>

      {moveError && <p className="text-xs text-destructive">{moveError}</p>}

      {confirmingMove && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-full max-w-sm rounded-none border border-border bg-background p-5">
            <div className="space-y-2">
              <h3 className="font-heading text-base font-medium">Move plate</h3>
              <p className="text-sm text-muted-foreground">
                This will move &quot;{confirmingMove.plateName}&quot; to {confirmingMove.targetOrgName}.
              </p>
              <p className="text-xs text-muted-foreground">
                Type the full plate name to confirm: <span className="font-medium text-foreground">{confirmingMove.plateName}</span>
              </p>
              <input
                value={moveConfirmNameInput}
                onChange={(e) => setMoveConfirmNameInput(e.target.value)}
                placeholder="Enter full plate name"
                className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none transition-colors focus:border-ring"
                autoFocus
              />
            </div>

            {moveError && <p className="mt-3 text-xs text-destructive">{moveError}</p>}

            <div className="mt-5 flex items-center justify-end gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={movePlateOrganization.isPending}
                onClick={() => {
                  setConfirmingMove(null)
                  setMoveConfirmNameInput("")
                  setMoveError(null)
                }}
              >
                Cancel
              </Button>
              <Button
                variant="default"
                size="sm"
                disabled={movePlateOrganization.isPending || !canMove}
                onClick={onConfirmMoveOrganization}
              >
                {movePlateOrganization.isPending ? "Moving..." : "Yes, move"}
              </Button>
            </div>
          </div>
        </div>
      )}

      <Dialog open={Boolean(orgDialogPlate)} onOpenChange={(open) => { if (!open) setOrgDialogPlate(null) }}>
        <DialogContent className="sm:max-w-md">
          {orgDialogPlate ? (
            <>
              <DialogHeader>
                <DialogTitle>Organization</DialogTitle>
                <DialogDescription>
                  Move &quot;{orgDialogPlate.name}&quot; between your personal namespace and an organization you manage.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-2 py-1">
                <label htmlFor={`org-select-${orgDialogPlate.id}`} className="text-xs text-muted-foreground">
                  Organization
                </label>
                <select
                  id={`org-select-${orgDialogPlate.id}`}
                  value={selectedOrgByPlate[orgDialogPlate.id] ?? (orgDialogPlate.organization_id ?? "")}
                  onChange={(e) => {
                    const value = e.target.value
                    setSelectedOrgByPlate((prev) => ({ ...prev, [orgDialogPlate.id]: value }))
                  }}
                  className="h-9 w-full border border-input bg-transparent px-2 text-sm outline-none transition-colors focus:border-ring"
                  disabled={orgsLoading || (movePlateOrganization.isPending && movingId === orgDialogPlate.id)}
                >
                  <option value="">Personal (no organization)</option>
                  {organizations.map((org) => (
                    <option key={org.id} value={org.id}>
                      {org.name}
                    </option>
                  ))}
                </select>
              </div>
              <DialogFooter className="gap-2 sm:justify-end">
                <Button type="button" variant="outline" size="sm" onClick={() => setOrgDialogPlate(null)}>
                  Cancel
                </Button>
                <Button
                  type="button"
                  size="sm"
                  disabled={orgsLoading || (movePlateOrganization.isPending && movingId === orgDialogPlate.id)}
                  onClick={() => {
                    const p = orgDialogPlate
                    const currentOrgId = p.organization_id ?? ""
                    const selectedOrgId = selectedOrgByPlate[p.id] ?? currentOrgId

                    if (selectedOrgId === currentOrgId) {
                      setMoveError("Select a different organization before moving this plate")
                      return
                    }

                    const targetOrgName = selectedOrgId === ""
                      ? "Personal (no organization)"
                      : (organizations.find((org) => org.id === selectedOrgId)?.name ?? "selected organization")

                    setMoveError(null)
                    setMoveConfirmNameInput("")
                    setOrgDialogPlate(null)
                    setConfirmingMove({
                      plateId: p.id,
                      plateName: p.name,
                      targetOrgId: selectedOrgId,
                      targetOrgName,
                    })
                  }}
                >
                  {movePlateOrganization.isPending && movingId === orgDialogPlate.id ? "Moving…" : "Move"}
                </Button>
              </DialogFooter>
            </>
          ) : null}
        </DialogContent>
      </Dialog>

      <PlateGetBadgeModal plate={getBadgePlate} onClose={() => setGetBadgePlate(null)} />

      {confirmingPlate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-full max-w-sm rounded-none border border-border bg-background p-5">
            <div className="space-y-2">
              <h3 className="font-heading text-base font-medium">Remove plate</h3>
              <p className="text-sm text-muted-foreground">
                This will permanently remove &quot;{confirmingPlate.name}&quot;. This action cannot be undone.
              </p>
              <p className="text-xs text-muted-foreground">
                Type the full plate name to confirm: <span className="font-medium text-foreground">{confirmingPlate.name}</span>
              </p>
              <input
                value={confirmNameInput}
                onChange={(e) => setConfirmNameInput(e.target.value)}
                placeholder="Enter full plate name"
                className="h-9 w-full border border-input bg-transparent px-3 text-sm outline-none transition-colors focus:border-ring"
                autoFocus
              />
            </div>

            {removeError && <p className="mt-3 text-xs text-destructive">{removeError}</p>}

            <div className="mt-5 flex items-center justify-end gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={removePlate.isPending}
                onClick={() => {
                  setConfirmingPlate(null)
                  setConfirmNameInput("")
                  setRemoveError(null)
                }}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                size="sm"
                disabled={removePlate.isPending || !canRemove}
                onClick={onConfirmRemove}
              >
                {removePlate.isPending ? "Removing..." : "Yes, remove"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}