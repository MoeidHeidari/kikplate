const DEFAULT_REPO = 'kikplate/kikplate'
const DEFAULT_REF = 'main'

export function docsGithubRepo(): string {
  return process.env.DOCS_GITHUB_REPO?.trim() || DEFAULT_REPO
}

export function docsGithubRef(): string {
  return process.env.DOCS_GITHUB_REF?.trim() || DEFAULT_REF
}

export function docsRawBaseUrl(): string {
  const [owner, repo] = parseRepo(docsGithubRepo())
  const ref = docsGithubRef()
  return `https://raw.githubusercontent.com/${owner}/${repo}/${ref}/docs`
}

function parseRepo(spec: string): [string, string] {
  const parts = spec.split('/').filter(Boolean)
  if (parts.length !== 2) {
    throw new Error(`DOCS_GITHUB_REPO must be "owner/repo", got: ${spec}`)
  }
  return [parts[0], parts[1]]
}

export function docsContentsApiUrl(): string {
  const [owner, repo] = parseRepo(docsGithubRepo())
  const ref = docsGithubRef()
  return `https://api.github.com/repos/${owner}/${repo}/contents/docs?ref=${encodeURIComponent(ref)}`
}

export function githubApiHeaders(): HeadersInit {
  const headers: Record<string, string> = {
    Accept: 'application/vnd.github+json',
    'X-GitHub-Api-Version': '2022-11-28',
  }
  const token = process.env.GITHUB_TOKEN?.trim()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return headers
}

export function rewriteDocAssets(markdown: string): string {
  const base = docsRawBaseUrl()

  let out = markdown.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (full, alt, src: string) => {
      const s = src.trim()
      if (/^https?:\/\//i.test(s) || s.startsWith('data:')) {
        return full
      }
      const path = s.replace(/^\.\//, '').replace(/^\//, '')
      return `![${alt}](${base}/${path})`
    },
  )

  out = out.replace(
    /<img\b([^>]*)\bsrc=(["'])([^"']+)\2([^>]*)>/gi,
    (full, before, quote, src: string, after) => {
      const s = String(src).trim()
      if (/^https?:\/\//i.test(s) || s.startsWith('data:')) {
        return full
      }
      const path = s.replace(/^\//, '').replace(/^\.\//, '')
      return `<img${before}src=${quote}${base}/${path}${quote}${after}>`
    },
  )

  return out
}

export const DOC_SLUG_ORDER: string[] = [
  'getting-started',
  'how-it-works',
  'architecture',
  'database',
  'configuration',
  'cli',
  'kubernetes',
  'helm',
  'contributing',
]

interface GitHubContentFile {
  name: string
  type: string
}

export function sortDocSlugs(slugs: { slug: string; name: string; file: string }[]) {
  const order = DOC_SLUG_ORDER
  return [...slugs].sort((a, b) => {
    const ia = order.indexOf(a.slug)
    const ib = order.indexOf(b.slug)
    const va = ia === -1 ? 999 : ia
    const vb = ib === -1 ? 999 : ib
    if (va !== vb) return va - vb
    return a.slug.localeCompare(b.slug)
  })
}

export async function fetchDocsIndex(): Promise<{ slug: string; name: string; file: string }[]> {
  const res = await fetch(docsContentsApiUrl(), {
    headers: githubApiHeaders(),
    next: { revalidate: 300 },
  })

  if (!res.ok) {
    const text = await res.text().catch(() => '')
    throw new Error(`GitHub API ${res.status}: ${text.slice(0, 200)}`)
  }

  const data = (await res.json()) as GitHubContentFile[] | { message?: string }
  if (!Array.isArray(data)) {
    throw new Error(
      typeof data === 'object' && data && 'message' in data
        ? String(data.message)
        : 'Unexpected GitHub API response',
    )
  }

  const mdFiles = data
    .filter((e) => e.type === 'file' && e.name.endsWith('.md') && e.name !== 'README.md')
    .map((e) => {
      const slug = e.name.replace(/\.md$/, '')
      const name = slug.replace(/-/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())
      return { slug, name, file: e.name }
    })

  return sortDocSlugs(mdFiles)
}

export async function fetchDocMarkdown(slug: string): Promise<string> {
  const safe = slug.replace(/[^a-z0-9-]/g, '')
  if (!safe) {
    throw new Error('Invalid doc slug')
  }

  const url = `${docsRawBaseUrl()}/${safe}.md`
  const res = await fetch(url, {
    headers: { Accept: 'text/plain' },
    next: { revalidate: 300 },
  })

  if (res.status === 404) {
    throw new Error('Not found')
  }
  if (!res.ok) {
    throw new Error(`Raw fetch ${res.status}`)
  }

  const text = await res.text()
  if (text.includes('404: Not Found') && text.length < 80) {
    throw new Error('Not found')
  }

  return rewriteDocAssets(text)
}
