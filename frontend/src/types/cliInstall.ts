export type CliInstallPhase =
  | 'pending'
  | 'downloading'
  | 'installed'
  | 'generating'
  | 'initializing'
  | 'waiting_auth'
  | 'verifying'
  | 'done'
  | 'failed'

export type CliInstallItem = {
  npm_package: string
  name: string
  phase: CliInstallPhase
  detail: string
  extension_id?: string
  session_id?: number
  action_url?: string
  action_label?: string
  expires_at?: string
}
