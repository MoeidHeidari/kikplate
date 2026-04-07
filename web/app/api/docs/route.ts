import { readFile, readdir } from 'fs/promises'
import { join } from 'path'
import { NextRequest, NextResponse } from 'next/server'

const docsDir = join(process.cwd(), '..', 'docs')

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const doc = searchParams.get('doc')

    if (!doc) {
      // Return list of available docs
      const files = await readdir(docsDir)
      const mdFiles = files
        .filter((file) => file.endsWith('.md') && file !== 'README.md')
        .map((file) => ({
          slug: file.replace('.md', ''),
          name: file.replace('.md', '').replace(/-/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase()),
          file,
        }))
        .sort((a, b) => {
          const order = [
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
          return order.indexOf(a.slug) - order.indexOf(b.slug)
        })

      return NextResponse.json(mdFiles)
    }

    // Return specific doc content
    const safeName = doc.replace(/[^a-z0-9-]/g, '')
    const filePath = join(docsDir, `${safeName}.md`)
    
    // Security: ensure file is within docs directory
    if (!filePath.startsWith(docsDir)) {
      return NextResponse.json({ error: 'Invalid doc' }, { status: 400 })
    }

    const content = await readFile(filePath, 'utf-8')
    return NextResponse.json({ content })
  } catch (error) {
    console.error('Error reading docs:', error)
    return NextResponse.json(
      { error: 'Failed to load documentation' },
      { status: 500 }
    )
  }
}
