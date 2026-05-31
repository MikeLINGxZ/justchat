import { type ReactNode, useContext } from 'react'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { cn } from '@/lib/utils'

export function SettingsSubmenuList<T extends string | number>(props: {
  title: string
  items: { key: T; label: string; icon?: ReactNode }[]
  value: T
  onChange: (value: T) => void
  action?: ReactNode
}) {
  const { isH5, onCollapseSubmenu } = useContext(SettingsMenuContext)

  return (
    <div className="flex h-full flex-col">
      <div className="mb-4 flex items-center justify-between gap-3 px-2">
        <h2 className="text-lg font-semibold text-foreground">{props.title}</h2>
        {props.action && <div className="flex items-center gap-1">{props.action}</div>}
      </div>

      <div className="space-y-1">
        {props.items.map((item) => (
          <button
            key={String(item.key)}
            type="button"
            onClick={() => { props.onChange(item.key); onCollapseSubmenu() }}
            className={cn(
              'flex w-full items-center gap-3 rounded-2xl px-3 py-2.5 text-left transition-colors',
              !isH5 && item.key === props.value
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-accent hover:text-foreground'
            )}
          >
            {item.icon && <span className="shrink-0">{item.icon}</span>}
            <span className="text-sm font-medium">{item.label}</span>
          </button>
        ))}
      </div>
    </div>
  )
}
