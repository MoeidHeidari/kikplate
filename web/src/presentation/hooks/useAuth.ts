"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { authRepository } from "@/src/data/repositories/AuthRepository"
import { LoginUseCase }       from "@/src/domain/usecases/LoginUseCase"
import { RegisterUseCase }    from "@/src/domain/usecases/RegisterUseCase"
import { GetMeUseCase }       from "@/src/domain/usecases/GetMeUseCase"
import { VerifyEmailUseCase } from "@/src/domain/usecases/VerifyEmailUseCase"
import { AuthService }        from "@/src/domain/services/AuthService"
import type { LoginInput, RegisterInput } from "@/src/domain/entities/User"

const loginUseCase       = new LoginUseCase(authRepository)
const registerUseCase    = new RegisterUseCase(authRepository)
const getMeUseCase       = new GetMeUseCase(authRepository)
const verifyEmailUseCase = new VerifyEmailUseCase(authRepository)

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: () => getMeUseCase.execute(),
    enabled: typeof window !== "undefined" && AuthService.isAuthenticated(),
    retry: false,
    staleTime: 5 * 60_000,
  })
}

export function useLogin() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: LoginInput) => loginUseCase.execute(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useRegister() {
  return useMutation({
    mutationFn: (input: RegisterInput) => registerUseCase.execute(input),
  })
}

export function useVerifyEmail() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (token: string) => verifyEmailUseCase.execute(token),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useLogout() {
  const qc = useQueryClient()
  return () => {
    AuthService.clearToken()
    qc.clear()
  }
}

export function useProviders() {
  return useQuery({
    queryKey: ["auth-providers"],
    queryFn: () => authRepository.providers(),
    staleTime: Infinity,
  })
}
