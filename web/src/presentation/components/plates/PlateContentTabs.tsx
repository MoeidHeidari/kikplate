"use client"

import { useState, useEffect, useMemo } from "react"
import type { ComponentPropsWithoutRef } from "react"
import { ChevronDown, ChevronRight, File, Folder } from "lucide-react"
import type { RepoTreeEntry } from "@/src/data/repositories/githubClient"

interface Props {
  readme: string | null
  license: string | null
  tree: RepoTreeEntry[] | null
}

type Tab = "readme" | "license" | "files"

interface FileTreeNode {
  name: string
  path: string
  type: "blob" | "tree"
  children: FileTreeNode[]
}

interface BuilderNode {
  name: string
  path: string
  type: "blob" | "tree"
  children: Map<string, BuilderNode>
}

function buildFileTree(entries: RepoTreeEntry[]): FileTreeNode[] {
  const root: BuilderNode = {
    name: "",
    path: "",
    type: "tree",
    children: new Map(),
  }

  const sorted = [...entries].sort((a, b) => a.path.localeCompare(b.path))

  for (const entry of sorted) {
    const parts = entry.path.split("/").filter(Boolean)
    if (parts.length === 0) continue

    let current = root
    let currentPath = ""

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      currentPath = currentPath ? `${currentPath}/${part}` : part
      const isLeaf = i === parts.length - 1

      const existing = current.children.get(part)
      if (existing) {
        if (!isLeaf) existing.type = "tree"
        if (isLeaf) existing.type = entry.type
        current = existing
        continue
      }

      const created: BuilderNode = {
        name: part,
        path: currentPath,
        type: isLeaf ? entry.type : "tree",
        children: new Map(),
      }

      current.children.set(part, created)
      current = created
    }
  }

  const toArray = (node: BuilderNode): FileTreeNode[] => {
    const items = Array.from(node.children.values())
      .sort((a, b) => {
        if (a.type !== b.type) return a.type === "tree" ? -1 : 1
        return a.name.localeCompare(b.name)
      })
      .map((child) => ({
        name: child.name,
        path: child.path,
        type: child.type,
        children: toArray(child),
      }))

    return items
  }

  return toArray(root)
}

function collectFolderPaths(nodes: FileTreeNode[]): string[] {
  const folderPaths: string[] = []

  const walk = (items: FileTreeNode[]) => {
    for (const item of items) {
      if (item.type === "tree") {
        folderPaths.push(item.path)
        if (item.children.length > 0) {
          walk(item.children)
        }
      }
    }
  }

  walk(nodes)
  return folderPaths
}

function TreeView({
  nodes,
  expandedPaths,
  onToggle,
  depth = 0,
}: {
  nodes: FileTreeNode[]
  expandedPaths: Set<string>
  onToggle: (path: string) => void
  depth?: number
}) {
  return (
    <ul className="space-y-1">
      {nodes.map((node) => (
        <li key={node.path}>
          <div
            className="flex items-center gap-2 py-0.5 text-sm"
            style={{ paddingLeft: `${depth * 16}px` }}
          >
            {node.type === "tree" ? (
              <button
                type="button"
                onClick={() => onToggle(node.path)}
                className="inline-flex items-center gap-1 text-left"
                aria-label={`${expandedPaths.has(node.path) ? "Collapse" : "Expand"} ${node.name}`}
              >
                {expandedPaths.has(node.path) ? (
                  <ChevronDown className="h-3.5 w-3.5 text-muted-foreground" />
                ) : (
                  <ChevronRight className="h-3.5 w-3.5 text-muted-foreground" />
                )}
                <Folder className="h-3.5 w-3.5 text-muted-foreground" />
                <span className="font-mono text-xs text-foreground">{node.name}</span>
              </button>
            ) : (
              <>
                <span className="inline-block w-3.5" />
                <File className="h-3.5 w-3.5 text-muted-foreground" />
                <span className="font-mono text-xs text-foreground">{node.name}</span>
              </>
            )}
          </div>
          {node.type === "tree" && node.children.length > 0 && expandedPaths.has(node.path) ? (
            <TreeView
              nodes={node.children}
              expandedPaths={expandedPaths}
              onToggle={onToggle}
              depth={depth + 1}
            />
          ) : null}
        </li>
      ))}
    </ul>
  )
}

export function PlateContentTabs({ readme, license, tree }: Props) {
  const [active, setActive] = useState<Tab>("readme")
  const [MarkdownComponent, setMarkdownComponent] = useState<any>(null)
  const [plugins, setPlugins] = useState<any[]>([])
  const treeNodes = useMemo(() => (tree ? buildFileTree(tree) : []), [tree])
  const folderPaths = useMemo(() => collectFolderPaths(treeNodes), [treeNodes])
  const [expandedPaths, setExpandedPaths] = useState<Set<string>>(new Set())

  useEffect(() => {
    const topLevelFolders = treeNodes
      .filter((node) => node.type === "tree")
      .map((node) => node.path)
    setExpandedPaths(new Set(topLevelFolders))
  }, [treeNodes])

  const toggleFolder = (path: string) => {
    setExpandedPaths((prev) => {
      const next = new Set(prev)
      if (next.has(path)) {
        next.delete(path)
      } else {
        next.add(path)
      }
      return next
    })
  }

  const expandAll = () => setExpandedPaths(new Set(folderPaths))
  const collapseAll = () => setExpandedPaths(new Set())

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

  useEffect(() => {
    const syncFromHash = () => {
      if (window.location.hash === "#files" && tree && tree.length > 0) {
        setActive("files")
        return
      }
      if (window.location.hash === "#license" && license) {
        setActive("license")
        return
      }
      if (window.location.hash === "#readme" && readme) {
        setActive("readme")
        return
      }

      if (readme) {
        setActive("readme")
        return
      }
      if (license) {
        setActive("license")
        return
      }
      if (tree && tree.length > 0) {
        setActive("files")
      }
    }

    syncFromHash()
    window.addEventListener("hashchange", syncFromHash)
    return () => window.removeEventListener("hashchange", syncFromHash)
  }, [readme, license, tree])

  const content = active === "readme" ? readme : active === "license" ? license : null

  const markdownComponents = {
    pre: ({ children, ...props }: ComponentPropsWithoutRef<"pre">) => (
      <pre
        {...props}
        className="my-4 overflow-x-auto border !border-border !bg-muted dark:!bg-card !p-6 text-xs leading-relaxed !text-foreground"
      >
        {children}
      </pre>
    ),
    code: ({ inline, className, children, ...props }: ComponentPropsWithoutRef<"code"> & { inline?: boolean }) => {
      if (inline) {
        return (
          <code
            {...props}
            className="rounded-none border border-border bg-muted/80 px-2 py-1 font-mono text-xs text-foreground dark:bg-card"
          >
            {children}
          </code>
        )
      }

      return (
        <code
          {...props}
          className={className ? `${className} !text-foreground font-mono` : "!text-foreground font-mono"}
        >
          {children}
        </code>
      )
    },
  }

  return (
    <div>
      <div>
        {active === "files" ? (
          treeNodes.length > 0 ? (
            <div className="border border-border bg-card p-6">
              <div className="mb-4 flex items-center justify-end gap-2 border-b border-border pb-3">
                <button
                  type="button"
                  onClick={expandAll}
                  className="h-8 border border-border px-3 text-xs font-medium text-foreground transition-colors hover:bg-muted"
                >
                  Expand all
                </button>
                <button
                  type="button"
                  onClick={collapseAll}
                  className="h-8 border border-border px-3 text-xs font-medium text-foreground transition-colors hover:bg-muted"
                >
                  Collapse all
                </button>
              </div>
              <TreeView
                nodes={treeNodes}
                expandedPaths={expandedPaths}
                onToggle={toggleFolder}
              />
            </div>
          ) : (
            <div className="p-6">
              <p className="py-4 text-center text-sm text-muted-foreground">No file tree found</p>
            </div>
          )
        ) : content ? (
          MarkdownComponent ? (
            <div className="prose prose-neutral dark:prose-invert max-w-none p-6
              prose-headings:font-semibold
              prose-headings:border-b prose-headings:border-border prose-headings:pb-2 prose-headings:mb-4
              prose-h1:text-2xl prose-h2:text-xl prose-h3:text-lg
              prose-p:text-sm prose-p:leading-relaxed prose-p:text-foreground
              prose-a:text-blue-500 prose-a:no-underline hover:prose-a:underline
              prose-code:before:content-none prose-code:after:content-none
              prose-pre:bg-transparent prose-pre:border-0 prose-pre:p-0 prose-pre:shadow-none
              prose-blockquote:border-l-4 prose-blockquote:border-border prose-blockquote:text-muted-foreground prose-blockquote:not-italic
              prose-table:text-sm prose-th:text-left prose-th:font-semibold prose-th:border prose-th:border-border prose-th:px-3 prose-th:py-2 prose-th:bg-muted
              prose-td:border prose-td:border-border prose-td:px-3 prose-td:py-2
              prose-img:border prose-img:border-border
              prose-hr:border-border
              prose-li:text-sm prose-li:marker:text-muted-foreground
              prose-strong:text-foreground prose-strong:font-semibold
            ">
              <MarkdownComponent remarkPlugins={plugins} components={markdownComponents}>
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