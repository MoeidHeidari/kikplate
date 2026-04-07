'use client'

import { useState, useEffect, ComponentPropsWithoutRef, Children, isValidElement, ReactNode } from 'react'

interface MarkdownRendererProps {
  content: string
}

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [MarkdownComponent, setMarkdownComponent] = useState<React.ComponentType<any> | null>(null)
  const [remarkPlugins, setRemarkPlugins] = useState<unknown[]>([])
  const [rehypePlugins, setRehypePlugins] = useState<unknown[]>([])

  useEffect(() => {
    Promise.all([
      import('react-markdown'),
      import('remark-gfm'),
      import('remark-breaks'),
      import('rehype-raw'),
    ]).then(([md, gfm, breaks, rehypeRaw]) => {
      setMarkdownComponent(() => md.default)
      setRemarkPlugins([gfm.default, breaks.default])
      setRehypePlugins([rehypeRaw.default])
    })
  }, [])

  const markdownComponents = {
    p: ({ children, ...props }: ComponentPropsWithoutRef<'p'> & { children?: ReactNode }) => {
      const nodes = Children.toArray(children)
      const badgeLikeNodes = nodes.filter((node) => {
        if (typeof node === 'string') return node.trim().length === 0
        if (!isValidElement<{ children?: ReactNode }>(node)) return false
        if (node.type !== 'a') return false
        const anchorChildren = Children.toArray(node.props.children)
        return anchorChildren.some((child) => {
          if (!isValidElement<{ src?: string }>(child)) return false
          return child.type === 'img' && typeof child.props.src === 'string' && child.props.src.includes('shields.io')
        })
      })

      const isBadgeOnlyParagraph = nodes.length > 0 && badgeLikeNodes.length === nodes.length
      if (isBadgeOnlyParagraph) {
        return (
          <span {...props} className="mx-1 inline-block align-middle">
            {children}
          </span>
        )
      }

      return (
        <p {...props} className="text-sm leading-relaxed text-foreground">
          {children}
        </p>
      )
    },
    pre: ({ children, ...props }: ComponentPropsWithoutRef<'pre'>) => (
      <pre
        {...props}
        className="my-4 overflow-x-auto border !border-border !bg-muted dark:!bg-card !p-6 text-xs leading-relaxed !text-foreground"
      >
        {children}
      </pre>
    ),
    code: ({ inline, className, children, ...props }: ComponentPropsWithoutRef<'code'> & { inline?: boolean }) => {
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
          className={className ? `${className} !text-foreground font-mono` : '!text-foreground font-mono'}
        >
          {children}
        </code>
      )
    },
  }

  if (!MarkdownComponent) {
    return (
      <div className="space-y-6 animate-pulse p-6">
        <div className="h-12 bg-muted rounded w-1/2" />
        <div className="space-y-3">
          <div className="h-4 bg-muted rounded" />
          <div className="h-4 bg-muted rounded w-5/6" />
        </div>
      </div>
    )
  }

  return (
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
      prose-img:rounded-md
      prose-hr:border-border
      prose-li:text-sm prose-li:marker:text-muted-foreground
      prose-strong:text-foreground prose-strong:font-semibold
      [&_div[align='center']_a]:inline-flex
      [&_div[align='center']_img]:align-middle
    ">
      <MarkdownComponent
        remarkPlugins={remarkPlugins}
        rehypePlugins={rehypePlugins}
        components={markdownComponents}
      >
        {content}
      </MarkdownComponent>
    </div>
  )
}

