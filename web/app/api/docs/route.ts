import { NextRequest, NextResponse } from 'next/server'
import { fetchDocMarkdown, fetchDocsIndex } from '@/lib/githubDocs'

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const doc = searchParams.get('doc')

    if (!doc) {
      const mdFiles = await fetchDocsIndex()
      return NextResponse.json(mdFiles)
    }

    const content = await fetchDocMarkdown(doc)
    return NextResponse.json({ content })
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error'
    if (message === 'Not found' || message === 'Invalid doc slug') {
      return NextResponse.json({ error: 'Documentation not found' }, { status: 404 })
    }
    console.error('Error loading docs from GitHub:', error)
    return NextResponse.json(
      { error: 'Failed to load documentation' },
      { status: 500 },
    )
  }
}
