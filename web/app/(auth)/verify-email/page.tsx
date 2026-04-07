"use client"

import Link from "next/link"
import { Suspense, useEffect, useMemo, useRef, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { authRepository } from "@/src/data/repositories/AuthRepository"
import { ApiError } from "@/src/data/repositories/httpClient"
import { AuthService } from "@/src/domain/services/AuthService"

function VerifyEmailContent() {
  const router = useRouter()
  const params = useSearchParams()
  const hasAttemptedRef = useRef(false)
  const redirectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const [status, setStatus] = useState<"idle" | "success" | "error">("idle")
  const [errorMessage, setErrorMessage] = useState("")

  const token = useMemo(() => params.get("token")?.trim() ?? "", [params])

  useEffect(() => {
    let cancelled = false

    async function runVerification() {
      if (hasAttemptedRef.current) {
        return
      }
      hasAttemptedRef.current = true

      if (!token) {
        setStatus("error")
        setErrorMessage("Verification token is missing. Please use the link from your email.")
        return
      }

      try {
        const result = await authRepository.verifyEmail(token)
        if (cancelled) return
        AuthService.setToken(result.token)
        setStatus("success")

        redirectTimerRef.current = setTimeout(() => {
          if (!cancelled) router.replace("/")
        }, 1200)
      } catch (err) {
        if (cancelled) return
        setStatus("error")

        if (err instanceof ApiError && err.isUnprocessable()) {
          setErrorMessage("This verification link is invalid or expired.")
          return
        }

        setErrorMessage(err instanceof Error ? err.message : "Email verification failed.")
      }
    }

    runVerification()

    return () => {
      cancelled = true
      if (redirectTimerRef.current) {
        clearTimeout(redirectTimerRef.current)
      }
    }
  }, [router, token])

  return (
    <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm border border-border bg-card p-6 text-center">
        {status === "idle" && <h1 className="text-2xl font-bold">Verifying your email...</h1>}

        {status === "success" && (
          <>
            <h1 className="text-2xl font-bold">Email verified</h1>
            <p className="mt-2 text-sm text-muted-foreground">Your account is now active. Redirecting...</p>
          </>
        )}

        {status === "error" && (
          <>
            <h1 className="text-2xl font-bold">Verification failed</h1>
            <p className="mt-2 text-sm text-muted-foreground">{errorMessage}</p>
            <div className="mt-4">
              <Link href="/login" className="text-sm text-foreground underline underline-offset-4 hover:text-primary">
                Go to login
              </Link>
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
          <div className="w-full max-w-sm border border-border bg-card p-6 text-center">
            <h1 className="text-2xl font-bold">Verifying your email...</h1>
          </div>
        </div>
      }
    >
      <VerifyEmailContent />
    </Suspense>
  )
}
