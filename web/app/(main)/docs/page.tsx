'use client'

import { Suspense, useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import { MarkdownRenderer } from '@/src/presentation/components/markdown/MarkdownRenderer'
import { BookOpen, Menu, X, ChevronLeft, ChevronRight } from 'lucide-react'

interface DocItem {
  slug: string
  name: string
}

function DocsPageContent() {
  const searchParams = useSearchParams()
  const [docs, setDocs] = useState<DocItem[]>([])
  const [currentDoc, setCurrentDoc] = useState(searchParams.get('doc') || 'getting-started')
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  // Load docs list
  useEffect(() => {
    const fetchList = async () => {
      try {
        const res = await fetch('/api/docs')
        if (res.ok) {
          const data = await res.json()
          setDocs(data)
        }
      } catch {
        console.error('Failed to load docs')
      }
    }
    fetchList()
  }, [])

  // Load current doc
  useEffect(() => {
    const fetchDoc = async () => {
      setLoading(true)
      setMobileOpen(false)
      try {
        const res = await fetch(`/api/docs?doc=${currentDoc}`)
        if (res.ok) {
          const data = await res.json()
          setContent(data.content)
        }
      } catch {
        console.error('Failed to load doc')
      } finally {
        setLoading(false)
      }
    }
    fetchDoc()
  }, [currentDoc])

  const currentDocName = docs.find((d) => d.slug === currentDoc)?.name
  const currentIdx = docs.findIndex((d) => d.slug === currentDoc)
  const prevDoc = currentIdx > 0 ? docs[currentIdx - 1] : null
  const nextDoc = currentIdx < docs.length - 1 ? docs[currentIdx + 1] : null

  return (
    <div className="min-h-screen bg-background mt-20">
      <div className="mx-auto max-w-7xl">
        {/* Mobile header */}
        <div className="lg:hidden border-b border-border bg-background/50 sticky top-20 z-40">
          <div className="flex items-center justify-between px-4 py-4">
            <div className="flex-1">
              <h1 className="text-lg font-semibold text-foreground truncate">{currentDocName || 'docs'}</h1>
            </div>
            <button
              onClick={() => setMobileOpen(!mobileOpen)}
              className="ml-4 p-2 hover:bg-muted rounded transition-colors"
              aria-label="Toggle menu"
            >
              {mobileOpen ? <X size={20} /> : <Menu size={20} />}
            </button>
          </div>

          {/* Mobile sidebar */}
          {mobileOpen && (
            <nav className="border-t border-border bg-background/50">
              <div className="px-4 py-3 space-y-1">
                {docs.map((doc) => (
                  <button
                    key={doc.slug}
                    onClick={() => setCurrentDoc(doc.slug)}
                    className={`w-full text-left px-3 py-2 rounded text-sm transition-colors ${
                      currentDoc === doc.slug
                        ? 'bg-primary text-primary-foreground'
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                    }`}
                  >
                    {doc.name}
                  </button>
                ))}
              </div>
            </nav>
          )}
        </div>

        {/* Main layout */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8 px-4 py-12 lg:py-16">
          {/* Sidebar */}
          <aside className="hidden lg:block">
            <nav className="sticky top-24 space-y-1 pr-4">
              {docs.map((doc) => (
                <button
                  key={doc.slug}
                  onClick={() => setCurrentDoc(doc.slug)}
                  className={`w-full text-left px-4 py-2.5 rounded text-sm font-medium transition-all ${
                    currentDoc === doc.slug
                      ? 'bg-primary text-primary-foreground'
                      : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                  }`}
                >
                  {doc.name}
                </button>
              ))}
            </nav>
          </aside>

          {/* Content */}
          <main className="lg:col-span-3">
            <div>
              {loading ? (
                <div className="space-y-6 animate-pulse p-6">
                  <div className="h-12 bg-muted rounded w-1/2" />
                  <div className="space-y-3">
                    <div className="h-4 bg-muted rounded" />
                    <div className="h-4 bg-muted rounded w-5/6" />
                  </div>
                </div>
              ) : content ? (
                <MarkdownRenderer content={content} />
              ) : (
              <div className="text-center py-16 text-muted-foreground">
                <BookOpen className="mx-auto mb-4 opacity-50" size={40} />
                <p>No documentation found</p>
              </div>
              )}
            </div>

            {/* Navigation */}
            {docs.length > 0 && !loading && (
              <div className="mt-16 pt-8 border-t border-border flex items-center justify-between gap-4">
                {prevDoc ? (
                  <button
                    onClick={() => setCurrentDoc(prevDoc.slug)}
                    className="flex items-center gap-3 px-4 py-2 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors group"
                  >
                    <ChevronLeft size={18} className="group-hover:-translate-x-0.5 transition-transform" />
                    <div className="text-left">
                      <div className="text-xs text-muted-foreground">Previous</div>
                      <div className="text-sm font-medium">{prevDoc.name}</div>
                    </div>
                  </button>
                ) : (
                  <div />
                )}

                {nextDoc ? (
                  <button
                    onClick={() => setCurrentDoc(nextDoc.slug)}
                    className="flex items-center gap-3 px-4 py-2 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors group ml-auto"
                  >
                    <div className="text-right">
                      <div className="text-xs text-muted-foreground">Next</div>
                      <div className="text-sm font-medium">{nextDoc.name}</div>
                    </div>
                    <ChevronRight size={18} className="group-hover:translate-x-0.5 transition-transform" />
                  </button>
                ) : (
                  <div />
                )}
              </div>
            )}
          </main>
        </div>
      </div>
    </div>
  )
}

export default function DocsPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-background mt-20">
          <div className="mx-auto max-w-7xl px-4 py-12 lg:py-16">
            <div className="space-y-6 animate-pulse p-6">
              <div className="h-12 bg-muted rounded w-1/2" />
              <div className="space-y-3">
                <div className="h-4 bg-muted rounded" />
                <div className="h-4 bg-muted rounded w-5/6" />
              </div>
            </div>
          </div>
        </div>
      }
    >
      <DocsPageContent />
    </Suspense>
  )
}
