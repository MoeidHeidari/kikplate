"use client"

import { useRouter } from "next/navigation"
import {
  Server, Globe, Layers, Smartphone,
  Terminal, Package, Wrench, MoreHorizontal
} from "lucide-react"

const CATEGORIES = [
  { slug: "backend",   label: "Backend",    icon: Server,         description: "APIs, services, databases" },
  { slug: "frontend",  label: "Frontend",   icon: Globe,          description: "Web UIs, SPAs, SSR apps" },
  { slug: "fullstack", label: "Full Stack",  icon: Layers,         description: "End-to-end project starters" },
  { slug: "mobile",    label: "Mobile",     icon: Smartphone,     description: "iOS, Android, React Native" },
  { slug: "cli",       label: "CLI",        icon: Terminal,       description: "Command line tools" },
  { slug: "devops",    label: "DevOps",     icon: Wrench,         description: "Docker, CI/CD, infra" },
  { slug: "library",   label: "Library",    icon: Package,        description: "Reusable packages and SDKs" },
  { slug: "other",     label: "Other",      icon: MoreHorizontal, description: "Everything else" },
]

export function CategoriesGrid() {
  const router = useRouter()

  return (
    <section className="container mx-auto px-4 py-16 border-t border-border">
      <div className="mb-8">
        <h2 className="text-2xl font-bold">Browse by category</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Find templates for your specific use case
        </p>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        {CATEGORIES.map(({ slug, label, icon: Icon, description }) => (
          <button
            key={slug}
            onClick={() => router.push(`/explore?category=${slug}`)}
            className="group flex flex-col gap-3 border border-border bg-card p-5 text-left hover:border-foreground/20 hover:bg-muted/50 transition-all"
          >
            <Icon className="h-5 w-5 text-muted-foreground group-hover:text-foreground transition-colors" />
            <div>
              <p className="text-sm font-semibold text-foreground">{label}</p>
              <p className="text-xs text-muted-foreground mt-0.5">{description}</p>
            </div>
          </button>
        ))}
      </div>
    </section>
  )
}