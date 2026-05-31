import type { ReactNode } from 'react'

export function SettingsDirtyGuard(props: {
  dirty: boolean
  children: ReactNode
}) {
  return (
    <div data-dirty={props.dirty ? 'true' : 'false'} className="flex min-h-0 flex-1 flex-col">
      {props.children}
    </div>
  )
}
