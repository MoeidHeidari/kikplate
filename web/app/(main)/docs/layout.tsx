import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Documentation - Kikplate',
  description: 'Complete documentation for Kikplate including guides, architecture, configuration, and API reference.',
}

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return <>{children}</>
}
