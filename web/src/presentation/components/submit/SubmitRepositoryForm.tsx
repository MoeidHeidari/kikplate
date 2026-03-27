"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useSubmitRepository } from "@/src/presentation/hooks/usePlates"
import { useMyOrganizations } from "@/src/presentation/hooks/useOrganizations"
import { useMe } from "@/src/presentation/hooks/useAuth"
import { Loader2, AlertCircle, Info } from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"

export function SubmitRepositoryForm() {
  const router = useRouter()
  const submit = useSubmitRepository()
  const { data: me } = useMe()
  const { data: organizations } = useMyOrganizations()
  const [repoUrl, setRepoUrl] = useState("")
  const [branch, setBranch] = useState("main")
  const [organizationId, setOrganizationId] = useState("")
  const selectedOrganization = organizations?.find((org) => org.id === organizationId)
  const isPersonalSubmission = !organizationId
  const ownerHint = selectedOrganization?.name ?? me?.username ?? "your-username"

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await submit.mutateAsync({
        repo_url: repoUrl,
        branch,
        organization_id: organizationId || undefined,
      })
      toast.success("Plate submitted. Complete verification from your account.")
      router.replace("/account?tab=plates")
      router.refresh()
    } catch {
    }
  }

  const errorMsg = submit.error instanceof Error ? submit.error.message : null

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="border border-border bg-muted/30 px-4 py-3 space-y-2">
        <div className="flex items-center gap-1.5">
          <Info className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
          <p className="text-xs font-semibold text-foreground">Before you submit</p>
        </div>
        <ul className="text-xs text-muted-foreground space-y-1 pl-5">
          <li>· The repository must be public</li>
          <li>
            · It must contain a{" "}
            <code className="font-mono bg-muted px-1 py-0.5">kickplate.yaml</code>{" "}
            at the root
          </li>
          <li>
            · The <code className="font-mono bg-muted px-1 py-0.5">owner</code> field in that
            file must match: <span className="font-mono text-foreground font-medium">{ownerHint}</span>{" "}
            ({isPersonalSubmission ? "personal submit uses your username" : "org submit uses organization name"})
          </li>
        </ul>
      </div>

      <div className="space-y-4">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Repository details
        </p>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">
            GitHub URL <span className="text-destructive">*</span>
          </label>
          <input
            required
            type="url"
            placeholder="https://github.com/username/repo"
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">Branch</label>
          <input
            type="text"
            placeholder="main"
            value={branch}
            onChange={(e) => setBranch(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">Organization</label>
          <select
            value={organizationId}
            onChange={(e) => setOrganizationId(e.target.value)}
            className="w-full border border-input bg-background px-3 py-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 transition-colors"
          >
            <option value="">Personal (no organization)</option>
            {organizations?.map((org) => (
              <option key={org.id} value={org.id}>{org.name}</option>
            ))}
          </select>
          {!organizations || organizations.length === 0 ? (
            <p className="text-xs text-muted-foreground">
              No organizations found. You can still submit personally using your username as owner.
            </p>
          ) : (
            <p className="text-xs text-muted-foreground">
              Choose Personal to submit under your account, or pick an organization.
            </p>
          )}
        </div>
      </div>

      {errorMsg && (
        <div className="border border-destructive/40 bg-destructive/5 px-4 py-3 space-y-2">
          <div className="flex items-start gap-2 text-sm text-destructive">
            <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
            <span className="font-medium">{errorMsg}</span>
          </div>
          {errorMsg.includes("owner field") && (
            <p className="text-xs text-muted-foreground pl-6">
              Open your <code className="font-mono bg-muted px-1 py-0.5">kickplate.yaml</code> and
              set <code className="font-mono bg-muted px-1 py-0.5">owner: {ownerHint}</code> then push and try again.
            </p>
          )}
          {errorMsg.includes("username") && (
            <p className="text-xs text-muted-foreground pl-6">
              Repository plates require a local account with a username. Your account was
              created via SSO and does not have one. Register a new account with email and
              password.
            </p>
          )}
          {(errorMsg.includes("not found") || errorMsg.includes("fetch")) && (
            <p className="text-xs text-muted-foreground pl-6">
              Make sure the repository is public, the URL is correct, and
              the <code className="font-mono bg-muted px-1 py-0.5">kickplate.yaml</code> exists
              on the <code className="font-mono bg-muted px-1 py-0.5">{branch}</code> branch.
            </p>
          )}
          {errorMsg.includes("conflict") && (
            <p className="text-xs text-muted-foreground pl-6">
              A plate with this name already exists. Rename your plate in{" "}
              <code className="font-mono bg-muted px-1 py-0.5">kickplate.yaml</code> and try again.
            </p>
          )}
        </div>
      )}

      <div className="border-t border-border pt-5">
        <Button
          type="submit"
          disabled={submit.isPending || !repoUrl}
          className="gap-2"
        >
          {submit.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
          Submit repository plate
        </Button>
      </div>
    </form>
  )
}