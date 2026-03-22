function repoToRaw(repoUrl: string, branch: string, file: string): string {
  const path = repoUrl
    .replace("https://github.com/", "")
    .replace("http://github.com/", "")
    .replace(/\/$/, "")
  return `https://raw.githubusercontent.com/${path}/${branch}/${file}`
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