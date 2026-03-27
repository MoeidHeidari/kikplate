import { LoginForm } from "@/src/presentation/components/auth/LoginForm"

export default function LoginPage() {
  return (
    <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm border border-border bg-card p-6">
        <h1 className="mb-6 text-2xl font-bold">Sign in</h1>
        <LoginForm />
      </div>
    </div>
  )
}
