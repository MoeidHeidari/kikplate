"use client"

import { useEffect } from "react"

const TOKEN_KEY = "kp_token"

export function AuthCookieSync() {
  useEffect(() => {
    const token = localStorage.getItem(TOKEN_KEY)
    if (token) {
      document.cookie = `${TOKEN_KEY}=${encodeURIComponent(token)}; Path=/; Max-Age=2592000; SameSite=Lax`
      return
    }
    document.cookie = `${TOKEN_KEY}=; Path=/; Max-Age=0; SameSite=Lax`
  }, [])

  return null
}
