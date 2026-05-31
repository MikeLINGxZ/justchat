import type { ReactNode } from 'react'
import { useContext } from 'react'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { cn } from '@/lib/utils'

export function SettingsContentLayout(props: {
  header: ReactNode
  footprint?: ReactNode
  children?: ReactNode
  noContentScroll?: boolean
}) {
  const { isH5 } = useContext(SettingsMenuContext)
  const horizontalPadding = isH5 ? 'px-5' : 'px-10'

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className={cn('shrink-0', horizontalPadding)}>{props.header}</div>
      {props.children && (
        <div className={cn(
          props.noContentScroll ? 'flex min-h-0 flex-1 flex-col overflow-hidden' : 'min-h-0 flex-1 overflow-y-auto',
          horizontalPadding,
        )}>
          {props.children}
        </div>
      )}
      {!props.children && <div className="min-h-0 flex-1" />}
      {props.footprint && (
        <div className={cn('shrink-0 pb-6 pt-5', horizontalPadding)}>{props.footprint}</div>
      )}
    </div>
  )
}
