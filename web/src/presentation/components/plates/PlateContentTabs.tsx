"use client"

import { useState, useEffect } from "react"

interface Props {
  readme: string | null
  license: string | null
}

type Tab = "readme" | "license"

export function PlateContentTabs({ readme, license }: Props) {
  const [active, setActive] = useState<Tab>("readme")
  const [MarkdownComponent, setMarkdownComponent] = useState<any>(null)
  const [plugins, setPlugins] = useState<any[]>([])

  useEffect(() => {
    Promise.all([
      import("react-markdown"),
      import("remark-gfm"),
      import("remark-breaks"),
    ]).then(([md, gfm, breaks]) => {
      setMarkdownComponent(() => md.default)
      setPlugins([gfm.default, breaks.default])
    })
  }, [])

  const tabs = [
    { id: "readme" as Tab, label: "README", available: Boolean(readme) },
    { id: "license" as Tab, label: "License", available: Boolean(license) },
  ]

  const content = active === "readme" ? readme : license

  return (
    <div>
      <div className="flex border-b border-border">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => tab.available && setActive(tab.id)}
            disabled={!tab.available}
            className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              active === tab.id
                ? "border-foreground text-foreground"
                : "border-transparent text-muted-foreground hover:text-foreground"
            } disabled:opacity-40 disabled:cursor-not-allowed`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div>
        {content ? (
          MarkdownComponent ? (
            <div className="prose prose-neutral dark:prose-invert max-w-none p-6
              prose-headings:font-semibold
              prose-headings:border-b prose-headings:border-border prose-headings:pb-2 prose-headings:mb-4
              prose-h1:text-2xl prose-h2:text-xl prose-h3:text-lg
              prose-p:text-sm prose-p:leading-relaxed prose-p:text-foreground
              prose-a:text-blue-500 prose-a:no-underline hover:prose-a:underline
              prose-code:bg-muted prose-code:px-1.5 prose-code:py-0.5 prose-code:text-xs prose-code:font-mono prose-code:border prose-code:border-border prose-code:rounded-none prose-code:before:content-none prose-code:after:content-none
              prose-pre:bg-muted prose-pre:border prose-pre:border-border prose-pre:rounded-none prose-pre:p-4 prose-pre:text-xs prose-pre:overflow-x-auto
              prose-blockquote:border-l-4 prose-blockquote:border-border prose-blockquote:text-muted-foreground prose-blockquote:not-italic
              prose-table:text-sm prose-th:text-left prose-th:font-semibold prose-th:border prose-th:border-border prose-th:px-3 prose-th:py-2 prose-th:bg-muted
              prose-td:border prose-td:border-border prose-td:px-3 prose-td:py-2
              prose-img:border prose-img:border-border
              prose-hr:border-border
              prose-li:text-sm prose-li:marker:text-muted-foreground
              prose-strong:text-foreground prose-strong:font-semibold
            ">
              <MarkdownComponent remarkPlugins={plugins}>
                {content}
              </MarkdownComponent>
            </div>
          ) : (
            <div className="p-6">
              <p className="text-xs text-muted-foreground">Loading…</p>
            </div>
          )
        ) : (
          <div className="p-6">
            <p className="text-sm text-muted-foreground text-center py-4">
              {active === "readme" ? "No README found" : "No LICENSE found"}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}