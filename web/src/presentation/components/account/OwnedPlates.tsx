"use client"

import Link from "next/link"
import { useState } from "react"
import { usePlates, useRemovePlate, useMovePlateOrganization } from "@/src/presentation/hooks/usePlates"
import { useMyOrganizations } from "@/src/presentation/hooks/useOrganizations"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { PendingVerification } from "@/src/presentation/components/plates/PendingVerification"
import { Button } from "@/components/ui/button"
import { GitBranch, Plus, Trash2 } from "lucide-react"

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
            {publishedPlates.map((plate) => (
            <div key={plate.id} className="rounded-xl border border-border bg-card p-4">
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <Link href={`/plates/${plate.slug}`} className="line-clamp-1 text-sm font-semibold hover:underline">
                    {plate.name}
                  </Link>
                  <div className="flex shrink-0 items-center gap-1.5">
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

              <div className="mt-4 flex items-center justify-between gap-2">
                <p className="text-[11px] uppercase tracking-wide text-muted-foreground">{plate.category}</p>
                <button
                  type="button"
                  onClick={() => {
                    setConfirmingPlate({ id: plate.id, name: plate.name })
                    setConfirmNameInput("")
                    setRemoveError(null)
                  }}
                  disabled={removePlate.isPending && removingId === plate.id}
                  className="inline-flex items-center gap-1 rounded-md border border-destructive/30 px-2 py-1 text-xs text-destructive transition-colors hover:bg-destructive/10 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  {removePlate.isPending && removingId === plate.id ? "Removing..." : "Remove"}
                </button>
              </div>

              <div className="mt-3 space-y-2 border-t border-border pt-3">
                <label className="text-[11px] uppercase tracking-wide text-muted-foreground">Organization</label>
                <div className="flex items-center gap-2">
                  <select
                    value={selectedOrgByPlate[plate.id] ?? (plate.organization_id ?? "")}
                    onChange={(e) => {
                      const value = e.target.value
                      setSelectedOrgByPlate((prev) => ({ ...prev, [plate.id]: value }))
                    }}
                    className="h-8 flex-1 border border-input bg-transparent px-2 text-xs outline-none transition-colors focus:border-ring"
                    disabled={orgsLoading || (movePlateOrganization.isPending && movingId === plate.id)}
                  >
                    <option value="">Personal (no organization)</option>
                    {organizations.map((org) => (
                      <option key={org.id} value={org.id}>
                        {org.name}
                      </option>
                    ))}
                  </select>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      const currentOrgId = plate.organization_id ?? ""
                      const selectedOrgId = selectedOrgByPlate[plate.id] ?? currentOrgId

                      if (selectedOrgId === currentOrgId) {
                        setMoveError("Select a different organization before moving this plate")
                        return
                      }

                      const targetOrgName = selectedOrgId === ""
                        ? "Personal (no organization)"
                        : (organizations.find((org) => org.id === selectedOrgId)?.name ?? "selected organization")

                      setMoveError(null)
                      setMoveConfirmNameInput("")
                      setConfirmingMove({
                        plateId: plate.id,
                        plateName: plate.name,
                        targetOrgId: selectedOrgId,
                        targetOrgName,
                      })
                    }}
                    disabled={orgsLoading || (movePlateOrganization.isPending && movingId === plate.id)}
                  >
                    {movePlateOrganization.isPending && movingId === plate.id ? "Moving..." : "Move"}
                  </Button>
                </div>
              </div>
            </div>
          ))}
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