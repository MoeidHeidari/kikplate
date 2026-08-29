"use client"

import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Copy, Check, CheckCircle2, XCircle, Pencil, Trash2, Github, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import type { MeResult } from "@/src/domain/entities/User"
import { EditProfileModal } from "./EditProfileModal"
import { DeleteAccountModal } from "./DeleteAccountModal"
import { useLogout } from "@/src/presentation/hooks/useAuth"
import { CLIENT_API_BASE } from "@/src/lib/client-api"

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Button
      onClick={handleCopy}
      variant="ghost"
      size="sm"
      className="h-6 w-6 p-0 text-muted-foreground hover:text-foreground"
    >
      {copied ? <Check className="h-3.5 w-3.5 text-green-500" /> : <Copy className="h-3.5 w-3.5" />}
    </Button>
  )
}

interface Row {
  label: string
  value: React.ReactNode
  copyable?: string
}

function oauthOrTrustedProvider(provider: string): boolean {
  return provider !== "local"
}

export function ProfileDetails({ me }: { me: MeResult }) {
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [connectingGitHub, setConnectingGitHub] = useState(false)
  const router = useRouter()
  const searchParams = useSearchParams()
  const logout = useLogout()
  const githubAccountStatus = searchParams.get("github_account_status")

  function connectGitHub() {
    setConnectingGitHub(true)
    window.location.assign(`${CLIENT_API_BASE}/auth/github/connect`)
  }

  const rows: Row[] = [
    {
      label: "Account ID",
      value: (
        <span className="max-w-xs truncate font-mono text-xs text-muted-foreground">
          {me.account_id}
        </span>
      ),
      copyable: me.account_id,
    },
    me.username
      ? {
          label: "Username",
          value: <span className="text-sm">{me.username}</span>,
        }
      : null,
    me.email
      ? {
          label: "Email",
          value: <span className="text-sm">{me.email}</span>,
        }
      : null,
    {
      label: "Provider",
      value: <span className="text-sm capitalize">{me.provider}</span>,
    },
    {
      label: "GitHub access",
      value: (
        <span className="text-sm text-muted-foreground">
          {me.github_connected
            ? `Connected${me.github_installation_account ? ` to ${me.github_installation_account}` : ""}`
            : "Not connected"}
        </span>
      ),
    },
    me.role
      ? {
          label: "Role",
          value: <span className="text-sm capitalize">{me.role}</span>,
        }
      : null,
    oauthOrTrustedProvider(me.provider) || me.is_active !== undefined
      ? {
          label: "Email verified",
          value:
            oauthOrTrustedProvider(me.provider) || me.is_active === true ? (
              <span className="flex items-center gap-1 text-sm text-green-600">
                <CheckCircle2 className="h-3.5 w-3.5" /> Verified
              </span>
            ) : (
              <span className="flex items-center gap-1 text-sm text-destructive">
                <XCircle className="h-3.5 w-3.5" /> Not verified
              </span>
            ),
        }
      : null,
  ].filter(Boolean) as Row[]

  return (
    <>
      <div className="max-w-lg space-y-6">
        {githubAccountStatus && (
          <div className="rounded-xl border border-border bg-muted/20 px-4 py-3 text-sm text-muted-foreground">
            {githubAccountStatus === "connected" && "GitHub account access is now connected for personal private repositories."}
            {githubAccountStatus === "type-mismatch" && "The GitHub App was installed on an organization. Personal access requires installation on your personal GitHub account."}
            {githubAccountStatus === "error" && "GitHub account connection did not complete. Try again from this page."}
          </div>
        )}

        <div>
          <p className="mb-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Account details
          </p>
          <div className="divide-y divide-border rounded-xl border border-border overflow-hidden">
            {rows.map((row) => (
              <div
                key={row.label}
                className="flex flex-col gap-2 px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
              >
                <span className="shrink-0 text-xs font-medium uppercase tracking-wide text-muted-foreground sm:w-32">
                  {row.label}
                </span>
                <div className="flex min-w-0 items-center gap-2 self-start sm:self-auto">
                  {row.value}
                  {row.copyable && <CopyButton value={row.copyable} />}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <Button
            onClick={connectGitHub}
            variant="outline"
            className="gap-1.5"
            disabled={connectingGitHub}
          >
            {connectingGitHub ? <Loader2 className="h-4 w-4 animate-spin" /> : <Github className="h-4 w-4" />}
            {me.github_connected ? "Reconnect GitHub" : "Connect GitHub"}
          </Button>
          <Button
            onClick={() => setEditOpen(true)}
            variant="outline"
            className="gap-1.5"
          >
            <Pencil className="h-4 w-4" />
            Edit profile
          </Button>
          <Button
            onClick={() => setDeleteOpen(true)}
            variant="destructive"
            className="gap-1.5"
          >
            <Trash2 className="h-3 w-3" />
            Delete account
          </Button>
        </div>
      </div>

      {editOpen && (
        <EditProfileModal
          me={me}
          onClose={() => setEditOpen(false)}
          onSaved={() => setEditOpen(false)}
        />
      )}

      {deleteOpen && (
        <DeleteAccountModal
          username={me.username ?? me.account_id}
          onClose={() => setDeleteOpen(false)}
          onDeleted={() => {
            setDeleteOpen(false)
            logout()
            router.push("/")
            router.refresh()
          }}
        />
      )}
    </>
  )
}
