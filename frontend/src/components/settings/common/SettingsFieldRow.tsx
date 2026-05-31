import type { ReactNode } from 'react'

export function SettingsFieldRow(props: {
  label: string
  description?: string
  children: ReactNode
}) {
  return (
    <div className="space-y-3 rounded-2xl bg-card/30 p-1">
      <div className="space-y-1">
        <p className="text-sm font-medium text-foreground">{props.label}</p>
        {props.description && <p className="text-sm text-muted-foreground">{props.description}</p>}
      </div>
      <div className="min-w-0">{props.children}</div>
    </div>
  )
}
