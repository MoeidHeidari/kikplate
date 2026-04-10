import { NextResponse } from "next/server"

export const revalidate = 600

function parseRepo(spec: string): { owner: string; repo: string } | null {
  const parts = spec.split("/").filter(Boolean)
  if (parts.length !== 2) return null
  return { owner: parts[0], repo: parts[1] }
}

async function fetchGitHubStars(): Promise<number | null> {
  const spec = process.env.STATS_GITHUB_REPO?.trim() || "kikplate/kikplate"
  const parsed = parseRepo(spec)
  if (!parsed) return null

  const headers: Record<string, string> = {
    Accept: "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
  }
  const pat = process.env.GITHUB_STATS_TOKEN?.trim()
  if (pat && pat !== "unset") {
    headers.Authorization = `Bearer ${pat}`
  }

  const res = await fetch(`https://api.github.com/repos/${parsed.owner}/${parsed.repo}`, {
    headers,
    next: { revalidate: 600 },
  })

  if (!res.ok) return null
  const data = (await res.json()) as { stargazers_count?: number }
  return typeof data.stargazers_count === "number" ? data.stargazers_count : null
}

async function fetchSlackHumanMemberCount(): Promise<number | null> {
  const token = process.env.SLACK_BOT_TOKEN?.trim()
  if (!token || token === "unset" || token === "changeme") {
    return null
  }

  let total = 0
  let cursor: string | undefined

  for (;;) {
    const body = new URLSearchParams()
    body.set("limit", "200")
    if (cursor) body.set("cursor", cursor)

    const res = await fetch("https://slack.com/api/users.list", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body,
      next: { revalidate: 600 },
    })

    const data = (await res.json()) as {
      ok?: boolean
      members?: Array<{ id?: string; deleted?: boolean; is_bot?: boolean }>
      response_metadata?: { next_cursor?: string }
    }

    if (!data.ok) return null

    for (const m of data.members ?? []) {
      if (!m.deleted && !m.is_bot && m.id && m.id !== "USLACKBOT") {
        total++
      }
    }

    const next = data.response_metadata?.next_cursor
    if (!next) break
    cursor = next
  }

  return total
}

export async function GET() {
  try {
    const [githubStars, slackMembers] = await Promise.all([
      fetchGitHubStars(),
      fetchSlackHumanMemberCount(),
    ])
    return NextResponse.json({ githubStars, slackMembers })
  } catch (e) {
    console.error("external stats:", e)
    return NextResponse.json({ githubStars: null, slackMembers: null }, { status: 200 })
  }
}
