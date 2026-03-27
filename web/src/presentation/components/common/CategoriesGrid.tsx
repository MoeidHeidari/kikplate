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
    <section className="bg-muted/10 py-16">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-12">
          <div className="border border-border bg-card p-5 sm:p-6 lg:col-span-4">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Taxonomy
            </p>
            <h2 className="mt-2 text-3xl font-bold tracking-tight text-foreground">Browse by category</h2>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Find templates by domain, from backend services to mobile apps and infrastructure.
            </p>
            <p className="mt-6 border-t border-border pt-4 text-xs text-muted-foreground">8 core domains</p>
          </div>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:col-span-8">
            {CATEGORIES.map(({ slug, label, icon: Icon, description }, i) => (
              <button
                key={slug}
                onClick={() => router.push(`/explore?category=${slug}`)}
                className="group flex items-start gap-4 border border-border bg-card p-5 text-left transition-colors hover:border-foreground/20 hover:bg-background"
              >
                <div className="mt-0.5 flex flex-col items-center gap-2 text-muted-foreground">
                  <Icon className="h-5 w-5 transition-colors group-hover:text-foreground" />
                  <span className="text-[10px] tabular-nums">{String(i + 1).padStart(2, "0")}</span>
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-foreground">{label}</p>
                  <p className="mt-1 text-xs text-muted-foreground">{description}</p>
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}