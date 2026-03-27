"use client"

import { useState } from "react"
import { Copy, Check, CheckCircle2, XCircle, Pencil, Trash2 } from "lucide-react"
import type { MeResult } from "@/src/domain/entities/User"
import { EditProfileModal } from "./EditProfileModal"
import { DeleteAccountModal } from "./DeleteAccountModal"

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)
  async function handleCopy() {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <button onClick={handleCopy} className="text-muted-foreground transition-colors hover:text-foreground">
      {copied ? <Check className="h-3.5 w-3.5 text-green-500" /> : <Copy className="h-3.5 w-3.5" />}
    </button>
  )
}





interface Row {
  label: string
  value: React.ReactNode
  copyable?: string
}

export function ProfileDetails({ me }: { me: MeResult }) {
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)

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
    me.role
      ? {
          label: "Role",
          value: <span className="text-sm capitalize">{me.role}</span>,
        }
      : null,
    me.is_active !== undefined
      ? {
          label: "Email verified",
          value: me.is_active ? (
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
        <div>
          <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Account details
          </p>
          <div className="divide-y divide-border border border-border">
            {rows.map((row) => (
              <div key={row.label} className="flex items-center justify-between bg-card px-4 py-3">
                <span className="w-32 shrink-0 text-xs text-muted-foreground">{row.label}</span>
                <div className="flex min-w-0 items-center gap-2">
                  {row.value}
                  {row.copyable && <CopyButton value={row.copyable} />}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => setEditOpen(true)}
            className="flex items-center gap-1.5 h-8 px-3 text-xs border border-border text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
          >
            <Pencil className="h-3 w-3" />
            Edit profile
          </button>
          <button
            onClick={() => setDeleteOpen(true)}
            className="flex items-center gap-1.5 h-8 px-3 text-xs border border-destructive/40 text-destructive/70 hover:text-destructive hover:border-destructive hover:bg-destructive/5 transition-colors"
          >
            <Trash2 className="h-3 w-3" />
            Delete account
          </button>
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
          onDeleted={() => setDeleteOpen(false)}
        />
      )}
    </>
  )
}