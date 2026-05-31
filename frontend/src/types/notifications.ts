export type NotificationKind = 'needs_attention' | 'info'

export type NotificationItem = {
  id: number
  session_id: number
  kind: NotificationKind
  title: string
  message: string
  created_at: string
  resolved_at: string | null
}
