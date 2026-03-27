"use client"

import Link from "next/link"
import { useRouter, usePathname } from "next/navigation"
import { useState, useEffect } from "react"
import { useMe, useLogout } from "@/src/presentation/hooks/useAuth"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { NavbarSearch } from "./NavbarSearch"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { LogOut, User, UserPlus, Sun, Moon, Github } from "lucide-react"
import { useTheme } from "next-themes"

const SlackIcon = () => (
  <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor">
    <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z" />
  </svg>
)

export function Navbar() {
  const { data: me } = useMe()
  const { data: appConfig } = useConfig()
  const logout = useLogout()
  const router = useRouter()
  const { theme, setTheme } = useTheme()
  const pathname = usePathname()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (mounted && me && !me.username && !window.location.pathname.startsWith("/set-username")) {
      router.push("/set-username")
    }
  }, [mounted, me])

  const initials = me?.username
    ? me.username.slice(0, 2).toUpperCase()
    : me?.display_name
    ? me.display_name.slice(0, 2).toUpperCase()
    : "?"

  function handleLogout() {
    logout()
    router.push("/")
    router.refresh()
  }

  function toggleTheme() {
    setTheme(theme === "dark" ? "light" : "dark")
  }

  if (!mounted) {
    return (
      <nav className="dark sticky top-0 z-50 bg-background">
        <div className="container mx-auto flex h-14 items-center justify-between px-4">
          <div className="h-7 w-32 rounded-sm bg-mutedanimate-pulse" />
          <div className="h-8 w-8 rounded-sm bg-mutedanimate-pulse" />
        </div>
      </nav>
    )
  }

  return (
    <nav className="dark sticky top-0 z-50 bg-background">
      <div className="container mx-auto flex h-20 items-center px-4 gap-4">
        <div className="flex items-center gap-2 cursor-pointer shrink-0" onClick={() => router.push("/")}>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={appConfig?.logo ?? "/kikplate-logo-on-dark.svg"}
            alt="logo"
            width={40}
            height={14}
          />
          <span className="text-xl tracking-tight text-white hidden sm:inline">
            Kik<span className="font-bold">Plate</span>
          </span>
        </div>

        {pathname !== "/" && !pathname?.startsWith("/explore") && (
          <div className="hidden md:flex">
            <NavbarSearch />
          </div>
        )}

        <div className="flex items-center gap-4 shrink-0 ml-auto">
          <div className="hidden sm:flex items-center gap-4 text-sm text-white/60">
            <Link href="/explore" className="hover:text-white transition-colors">
              Explore
            </Link>
            {me && (
              <Link href="/submit" className="hover:text-white transition-colors">
                Submit
              </Link>
            )}
          </div>

          <div className="hidden md:flex items-center gap-3 text-white/40 px-3">
            <Link href="#" className="hover:text-white transition-colors" title="Docs">
              docs
            </Link>
            <Link href="/stats" className="hover:text-white transition-colors" title="Stats">
              stats
            </Link>
            <Link
              href="https://github.com/kickplate"
              target="_blank"
              className="hover:text-white transition-colors"
              title="GitHub"
            >
              <Github className="h-4 w-4" />
            </Link>
            <Link href="#" className="hover:text-white transition-colors" title="Slack">
              <SlackIcon />
            </Link>
          </div>

          <div className="hidden sm:block h-4 w-px bg-white/10" />

          {me ? (
            <DropdownMenu>
              <DropdownMenuTrigger className="outline-none focus-visible:ring-2 focus-visible:ring-white/30 focus-visible:ring-offset-2 ring-offset-background">
                <Avatar className="h-8 w-8 cursor-pointer rounded-sm">
                  {me.avatar_url && <AvatarImage src={me.avatar_url} alt={me.username ?? me.display_name ?? "avatar"} />}
                  <AvatarFallback className="text-xs bg-white/20 text-white rounded-sm font-semibold">
                    {initials}
                  </AvatarFallback>
                </Avatar>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-52 rounded-sm">
                <div className="px-2 py-2 border-b border-border">
                  <p className="text-sm font-semibold">{me.username ?? me.display_name ?? "User"}</p>
                  <p className="text-xs text-muted-foreground truncate mt-0.5">{me.email ?? me.provider}</p>
                </div>
                <div className="py-1">
                  <DropdownMenuItem
                    className="cursor-pointer rounded-none gap-2 text-sm"
                    onClick={() => router.push("/account")}
                  >
                    <User className="h-4 w-4" />
                    Account
                  </DropdownMenuItem>
                </div>
                <div className="border-t border-border py-1">
                  <DropdownMenuItem
                    onClick={toggleTheme}
                    className="cursor-pointer rounded-none gap-2 text-sm"
                  >
                    {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                    {theme === "dark" ? "Light mode" : "Dark mode"}
                  </DropdownMenuItem>
                </div>
                <div className="border-t border-border py-1">
                  <DropdownMenuItem
                    onClick={handleLogout}
                    className="cursor-pointer rounded-none gap-2 text-sm text-destructive focus:text-destructive"
                  >
                    <LogOut className="h-4 w-4" />
                    Sign out
                  </DropdownMenuItem>
                </div>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <DropdownMenu>
              <DropdownMenuTrigger className="outline-none focus-visible:ring-2 focus-visible:ring-white/30 focus-visible:ring-offset-2 ring-offset-[#1a1f2e]">
                <Avatar className="h-8 w-8 cursor-pointer rounded-sm">
                  <AvatarFallback className="text-xs bg-muted text-white/60 rounded-sm">
                    ?
                  </AvatarFallback>
                </Avatar>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-52 rounded-sm">
                <div className="px-2 py-2 border-b border-border">
                  <p className="text-sm text-muted-foreground">Not signed in</p>
                </div>
                <div className="py-1">
                  <DropdownMenuItem
                    className="cursor-pointer rounded-none gap-2 text-sm"
                    onClick={() => router.push("/login")}
                  >
                    <User className="h-4 w-4" />
                    Sign in
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="cursor-pointer rounded-none gap-2 text-sm"
                    onClick={() => router.push("/register")}
                  >
                    <UserPlus className="h-4 w-4" />
                    Sign up
                  </DropdownMenuItem>
                </div>
                <div className="border-t border-border py-1">
                  <DropdownMenuItem
                    onClick={toggleTheme}
                    className="cursor-pointer rounded-none gap-2 text-sm"
                  >
                    {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                    {theme === "dark" ? "Light mode" : "Dark mode"}
                  </DropdownMenuItem>
                </div>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>

      </div>
    </nav>
  )
}