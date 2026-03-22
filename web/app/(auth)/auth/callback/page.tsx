"use client"

import { useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { AuthService } from "@/src/domain/services/AuthService"

export default function CallbackPage() {
  const router = useRouter()
  const params = useSearchParams()

  useEffect(() => {
    const token = params.get("token")
    const error = params.get("error")

    if (token) {
      AuthService.setToken(token)
      router.replace("/")
    } else if (error) {
      router.replace("/login?error=" + error)
    } else {
      router.replace("/login")
    }
  }, [])

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-sm text-muted-foreground">Completing sign in…</p>
    </div>
  )
}