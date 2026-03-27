interface Stat {
  label: string
  value: number | string
}

interface Props {
  stats: Stat[]
}

export function AccountStats({ stats }: Props) {
  return (
    <div className="flex w-fit border border-border bg-card">
      {stats.map((s, i) => (
        <div
          key={s.label}
          className={`flex flex-col items-center gap-0.5 px-6 py-3 ${
            i < stats.length - 1 ? "border-r border-border" : ""
          }`}
        >
          <span className="text-xl font-bold tabular-nums text-foreground">
            {s.value}
          </span>
          <span className="text-xs text-muted-foreground">{s.label}</span>
        </div>
      ))}
    </div>
  )
}