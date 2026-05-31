export type AlertKind = 'error' | 'info' | 'success' | 'warning'

export type AlertPlacement = 'toast' | 'banner'

export type AlertActionStyle = 'primary' | 'secondary' | 'danger'

export type AlertCloseReason =
  | 'user'
  | 'timeout'
  | 'programmatic'
  | 'replaced'
  | 'action'

export type AlertAriaPriority = 'assertive' | 'polite'

export type AlertSource = 'frontend' | 'backend'

export type AlertDedupeBehavior = 'count' | 'replace'

export type AlertAction = {
  id: string
  label: string
  style: AlertActionStyle
  closeOnClick: boolean
  href?: string
  onClick?: () => void
}

export type AlertInput = {
  kind: AlertKind
  placement: AlertPlacement
  title: string
  message: string
  detail?: string | null
  code?: string | null
  actions?: AlertAction[]
  dismissible?: boolean
  autoClose?: boolean
  durationMs?: number | null
  dedupeKey?: string
  dedupeBehavior?: AlertDedupeBehavior
  showCount?: boolean
  source?: AlertSource
}

export type AlertItem = {
  id: string
  kind: AlertKind
  placement: AlertPlacement
  title: string
  message: string
  detail: string | null
  code: string | null
  actions: AlertAction[]
  dismissible: boolean
  autoClose: boolean
  durationMs: number | null
  count: number
  dedupeKey: string
  dedupeBehavior: AlertDedupeBehavior
  showCount: boolean
  createdAt: number
  updatedAt: number
  ariaPriority: AlertAriaPriority
  source: AlertSource
}

export type ClosedAlertRecord = {
  id: string
  reason: AlertCloseReason
}
