function repoToRaw(repoUrl: string, branch: string, file: string): string {
  const path = repoToPath(repoUrl)
  return `https://raw.githubusercontent.com/${path}/${branch}/${file}`
}

function repoToPath(repoUrl: string): string {
  return repoUrl
    .replace("https://github.com/", "")
    .replace("http://github.com/", "")
    .replace(/\.git$/, "")
    .replace(/\/$/, "")
}

export interface RepoTreeEntry {
  path: string
  type: "blob" | "tree"
}

export async function fetchRepoFile(
  repoUrl: string,
  branch: string,
  file: string
): Promise<string | null> {
  try {
    const res = await fetch(repoToRaw(repoUrl, branch, file), {
      next: { revalidate: 3600 },
    })
    if (!res.ok) return null
    return res.text()
  } catch {
    return null
  }
}

export async function fetchRepoTree(
  repoUrl: string,
  branch: string
): Promise<RepoTreeEntry[] | null> {
  try {
    const repoPath = repoToPath(repoUrl)
    const res = await fetch(
      `https://api.github.com/repos/${repoPath}/git/trees/${encodeURIComponent(branch)}?recursive=1`,
      { next: { revalidate: 3600 } }
    )
    if (!res.ok) return null

    const data = await res.json() as { tree?: Array<{ path?: string; type?: string }> }
    if (!Array.isArray(data.tree)) return null

    return data.tree
      .filter((entry): entry is { path: string; type: "blob" | "tree" } => {
        return Boolean(entry.path) && (entry.type === "blob" || entry.type === "tree")
      })
      .slice(0, 1200)
  } catch {
    return null
  }
}