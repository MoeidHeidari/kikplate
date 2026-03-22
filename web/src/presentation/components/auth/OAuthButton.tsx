"use client"

import { Github, Chrome, Gitlab } from "lucide-react"
import { authRepository } from "@/src/data/repositories/AuthRepository"

const providerConfig: Record<string, { label: string; icon: React.ReactNode; bg: string }> = {
  github: {
    label: "GitHub",
    icon: <Github className="h-4 w-4" />,
    bg: "bg-[#24292e] hover:bg-[#2f363d] text-white border-[#24292e]",
  },
  google: {
    label: "Google",
    icon: <Chrome className="h-4 w-4" />,
    bg: "bg-white hover:bg-gray-50 text-gray-900 border-gray-300 dark:bg-white dark:text-gray-900",
  },
  gitlab: {
    label: "GitLab",
    icon: <Gitlab className="h-4 w-4" />,
    bg: "bg-[#fc6d26] hover:bg-[#e24329] text-white border-[#fc6d26]",
  },
}

interface Props {
  provider: string
}

export function OAuthButton({ provider }: Props) {
  const config = providerConfig[provider] ?? {
    label: provider.charAt(0).toUpperCase() + provider.slice(1),
    icon: null,
    bg: "bg-muted hover:bg-muted/80 text-foreground border-border",
  }

  function handleClick() {
    window.location.href = authRepository.oauthRedirectURL(provider)
  }

  return (
    <button
      onClick={handleClick}
      className={`flex h-10 w-full items-center justify-center gap-2 border text-sm font-medium transition-colors ${config.bg}`}
    >
      {config.icon}
      Continue with {config.label}
    </button>
  )
}