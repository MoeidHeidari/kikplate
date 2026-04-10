"use client"

import Link from "next/link"
import { useConfig } from "@/src/presentation/hooks/useConfig"
import { getSocialLink } from "@/src/lib/socialLinks"

const COMMUNITY = [
  { label: "GitHub", type: "github" as const },
  { label: "Slack", type: "slack" as const },
  { label: "X", type: "x" as const },
  { label: "LinkedIn", type: "linkedin" as const },
]

export function Footer() {
  const { data: appConfig } = useConfig()

  const communityLinks = COMMUNITY.map(({ label, type }) => ({
    label,
    href: getSocialLink(appConfig?.social_media, type),
  })).filter((l): l is { label: string; href: string } => Boolean(l.href))

  return (
    <footer className="border-t border-border bg-background">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-8 mb-12">
          <div className="col-span-2 sm:col-span-1 space-y-3">
            <p className="text-sm font-semibold text-foreground">
              Kik<span className="font-bold">Plate</span>
            </p>
            <p className="text-xs text-muted-foreground leading-relaxed">
              The biggest library of production-ready project templates and boilerplates.
            </p>
            <Link
              href="https://www.gnu.org/licenses/agpl-3.0.html"
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-muted-foreground underline underline-offset-4 hover:text-foreground transition-colors"
            >
              AGPL-3.0 License
            </Link>
          </div>

          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Product
            </p>
            <ul className="space-y-2">
              {[
                { label: "Explore", href: "/explore" },
                { label: "Stats", href: "/stats" },
                { label: "Submit a plate", href: "/submit" },
                { label: "Account", href: "/account" },
              ].map((link) => (
                <li key={link.label}>
                  <Link
                    href={link.href}
                    className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Community
            </p>
            <ul className="space-y-2">
              {communityLinks.map((link) => (
                <li key={link.label}>
                  <Link
                    href={link.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Resources
            </p>
            <ul className="space-y-2">
              <li>
                <Link
                  href="/docs"
                  className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  Documentation
                </Link>
              </li>
              <li>
                <Link
                  href="/docs?doc=cli"
                  className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  CLI
                </Link>
              </li>
              <li>
                <Link
                  href="https://github.com/kikplate/kikplate/releases"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  Changelog
                </Link>
              </li>
              <li>
                <Link
                  href="https://github.com/kikplate/kikplate/blob/main/CONTRIBUTING.md"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  Contributing
                </Link>
              </li>
            </ul>
          </div>
        </div>
        <div className="border-t border-border pt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-xs text-muted-foreground">
            © {new Date().getFullYear()} KikPlate. Open source and free forever.
          </p>
          <div className="flex items-center gap-4">
            {[
              { label: "Privacy", href: "#" },
              { label: "Terms", href: "#" },
            ].map((link) => (
              <Link
                key={link.label}
                href={link.href}
                className="text-xs text-muted-foreground hover:text-foreground transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>
      </div>
    </footer>
  )
}
