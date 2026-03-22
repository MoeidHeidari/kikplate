import Link from "next/link"

export function Footer() {
  return (
    <footer className="border-t border-border bg-background">
      <div className="container mx-auto px-4 py-12">

        {/* top — brand + columns */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-8 mb-12">

          {/* brand */}
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

          {/* product */}
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Product
            </p>
            <ul className="space-y-2">
              {[
                { label: "Explore", href: "/explore" },
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

          {/* community */}
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Community
            </p>
            <ul className="space-y-2">
              {[
                { label: "GitHub", href: "https://github.com/kickplate" },
                { label: "Slack", href: "#" },
                { label: "Twitter / X", href: "#" },
              ].map((link) => (
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

          {/* resources */}
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Resources
            </p>
            <ul className="space-y-2">
              {[
                { label: "Documentation", href: "#" },
                { label: "CLI", href: "#" },
                { label: "Changelog", href: "#" },
                { label: "Contributing", href: "https://github.com/kickplate/kickplate/blob/main/CONTRIBUTING.md" },
              ].map((link) => (
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

        </div>

        {/* bottom — copyright + legal */}
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