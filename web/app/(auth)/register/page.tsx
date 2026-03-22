import { RegisterForm } from "@/src/presentation/components/auth/RegisterForm"

export default function RegisterPage() {
  return (
    <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm">
        <h1 className="mb-6 text-2xl font-bold">Create account</h1>
        <RegisterForm />
      </div>
    </div>
  )
}
