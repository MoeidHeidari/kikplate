"use client"

import { useEffect, useState } from "react"

type HeaderTab = "readme" | "license" | "files" | "content"

interface Props {
  isRepository: boolean
  hasReadme?: boolean
  hasLicense?: boolean
  hasTree?: boolean
}

export function PlateHeaderTabs({ isRepository, hasReadme = false, hasLicense = false, hasTree = false }: Props) {
  const [active, setActive] = useState<HeaderTab>(isRepository ? "readme" : "content")

  useEffect(() => {
    if (!isRepository) {
      setActive("content")
      return
    }

    const syncFromHash = () => {
      if (window.location.hash === "#files" && hasTree) {
        setActive("files")
        return
      }
      if (window.location.hash === "#license" && hasLicense) {
        setActive("license")
        return
      }
      if (window.location.hash === "#readme" && hasReadme) {
        setActive("readme")
        return
      }
      setActive(hasReadme ? "readme" : hasLicense ? "license" : "files")
    }

    syncFromHash()
    window.addEventListener("hashchange", syncFromHash)
    return () => window.removeEventListener("hashchange", syncFromHash)
  }, [isRepository, hasReadme, hasLicense, hasTree])

  const setHashWithoutScroll = (tab: "readme" | "license" | "files") => {
    setActive(tab)
    const current = window.location.href.split("#")[0]
    window.history.replaceState(window.history.state, "", `${current}#${tab}`)
    window.dispatchEvent(new HashChangeEvent("hashchange"))
  }

  return (
    <nav className="flex h-12 items-end gap-1 overflow-x-auto text-sm">
      {isRepository ? (
        <>
          <button
            type="button"
            onClick={() => setHashWithoutScroll("readme")}
            disabled={!hasReadme}
            className={`inline-flex h-10 items-center border-b-2 px-3 font-semibold transition-colors ${
              active === "readme"
                ? "border-foreground text-foreground"
                : "border-transparent text-muted-foreground hover:border-foreground/30 hover:text-foreground"
            } disabled:cursor-not-allowed disabled:opacity-40`}
          >
            README
          </button>
          <button
            type="button"
            onClick={() => setHashWithoutScroll("license")}
            disabled={!hasLicense}
            className={`inline-flex h-10 items-center border-b-2 px-3 font-semibold transition-colors ${
              active === "license"
                ? "border-foreground text-foreground"
                : "border-transparent text-muted-foreground hover:border-foreground/30 hover:text-foreground"
            } disabled:cursor-not-allowed disabled:opacity-40`}
          >
            License
          </button>
          <button
            type="button"
            onClick={() => setHashWithoutScroll("files")}
            disabled={!hasTree}
            className={`inline-flex h-10 items-center border-b-2 px-3 font-semibold transition-colors ${
              active === "files"
                ? "border-foreground text-foreground"
                : "border-transparent text-muted-foreground hover:border-foreground/30 hover:text-foreground"
            } disabled:cursor-not-allowed disabled:opacity-40`}
          >
            Files
          </button>
        </>
      ) : (
        <span className="inline-flex h-10 items-center border-b-2 border-foreground px-3 font-semibold text-foreground">
          Content
        </span>
      )}
    </nav>
  )
}
