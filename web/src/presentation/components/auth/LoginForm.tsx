"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { useLogin, useProviders } from "@/src/presentation/hooks/useAuth"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { OAuthButton } from "./OAuthButton"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export function LoginForm() {
  const router = useRouter()
  const login = useLogin()
  const { data: providersData } = useProviders()
  const providers = providersData?.providers ?? []
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await login.mutateAsync({ email, password })
      toast.success("Welcome back!")
      router.push("/")
      router.refresh()
    } catch (err: any) {
      toast.error(err?.message ?? "Login failed")
    }
  }

  return (
    <div className="space-y-6">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoComplete="email"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            placeholder="••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
          />
        </div>
        <Button type="submit" className="w-full" disabled={login.isPending}>
          {login.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Sign in
        </Button>
      </form>

      {providers.length > 0 && (
        <div className="space-y-3">
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center text-xs">
              <span className="bg-background px-2 text-muted-foreground">or continue with</span>
            </div>
          </div>
          <div className="space-y-2">
            {providers.map((p) => (
              <OAuthButton key={p} provider={p} />
            ))}
          </div>
        </div>
      )}

      <p className="text-center text-sm text-muted-foreground">
        Don't have an account?{" "}
        <Link href="/register" className="text-foreground underline underline-offset-4 hover:text-primary">
          Sign up
        </Link>
      </p>
    </div>
  )
}