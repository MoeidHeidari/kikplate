"use client"

import Link from "next/link"
import { useStats } from "@/src/presentation/hooks/usePlates"
import { formatCount } from "@/src/presentation/utils/plateUtils"
import { Github, Users, Layers, Download, Package } from "lucide-react"

const SlackIcon = () => (
  <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor">
    <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z" />
  </svg>
)

interface StatRowProps {
  icon: React.ReactNode
  value: string
  label: string
}

function StatRow({ icon, value, label }: StatRowProps) {
  return (
    <div className="flex items-center justify-between py-4 border-b border-gray-100 last:border-0">
      <div className="flex items-center gap-3 text-gray-400">
        {icon}
        <span className="text-sm text-gray-700">{label}</span>
      </div>
      <span className="text-2xl font-bold tabular-nums text-gray-900">{value}</span>
    </div>
  )
}

export function StatsBanner() {
  const { data } = useStats()

  const leftStats = [
    {
      icon: <Package className="h-4 w-4" />,
      value: data ? formatCount(data.total_plates) : "—",
      label: "Templates",
    },
    {
      icon: <Users className="h-4 w-4" />,
      value: data ? formatCount(data.total_contributors) : "—",
      label: "Contributors",
    },
    {
      icon: <Layers className="h-4 w-4" />,
      value: data ? formatCount(data.total_categories) : "—",
      label: "Categories",
    },
  ]

  const rightStats = [
    {
      icon: <Download className="h-4 w-4" />,
      value: data ? formatCount(data.total_uses) : "—",
      label: "Total Uses",
    },
    {
      icon: <Github className="h-4 w-4" />,
      value: "—",
      label: "GitHub Stars",
    },
    {
      icon: <SlackIcon />,
      value: "—",
      label: "Slack Members",
    },
  ]

  return (
    <section className="bg-white border-t border-gray-100 py-20">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-16">

          <div className="flex flex-col justify-center gap-6">
            <div>
              <p className="text-xs uppercase tracking-widest text-gray-400 mb-3">
                By the numbers
              </p>
              <h2 className="text-3xl font-bold leading-tight text-gray-900">
                Built by the community,<br />for the community
              </h2>
            </div>
            <p className="text-gray-500 leading-relaxed max-w-md">
              KikPlate is an open source registry of production-ready project templates.
              Every template is submitted by a developer, reviewed by the community,
              and free to use forever.
            </p>
            <div className="flex items-center gap-4">
              <Link
                href="https://github.com/kickplate"
                target="_blank"
                className="text-sm font-medium text-gray-900 underline underline-offset-4 hover:text-gray-500 transition-colors"
              >
                View on GitHub →
              </Link>
              <Link
                href="#"
                className="text-sm font-medium text-gray-900 underline underline-offset-4 hover:text-gray-500 transition-colors"
              >
                Join Slack →
              </Link>
            </div>
          </div>

          <div className="flex flex-col gap-0 lg:pl-16 lg:border-l border-gray-100">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-12">
              <div>
                {leftStats.map((s) => (
                  <StatRow key={s.label} icon={s.icon} value={s.value} label={s.label} />
                ))}
              </div>
              <div>
                {rightStats.map((s) => (
                  <StatRow key={s.label} icon={s.icon} value={s.value} label={s.label} />
                ))}
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}