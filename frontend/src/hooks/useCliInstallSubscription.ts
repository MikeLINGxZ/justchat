import { useEffect } from 'react'
import { Events } from '@wailsio/runtime'
import { useCliInstallStore } from '@/store/cliInstallStore'
import type { CliInstallItem } from '@/types/cliInstall'

type ProgressEvent = {
  session_id?: number
  npm_package: string
  name: string
  phase: string
  detail: string
  id?: string
  extension_id?: string
  action_url?: string
  action_label?: string
  expires_at?: string
}

type DoneEvent = {
  extension: { id: string; name: string }
}

export function useCliInstallSubscription() {
  const upsert = useCliInstallStore((s) => s.upsert)
  const remove = useCliInstallStore((s) => s.remove)

  useEffect(() => {
    const offProgress = Events.On('cli.install.progress', (event: { data: ProgressEvent }) => {
      const d = event.data
      upsert({
        npm_package: d.npm_package ?? d.id ?? '',
        name: d.name ?? '',
        phase: d.phase as CliInstallItem['phase'],
        detail: d.detail ?? '',
        extension_id: d.extension_id ?? d.id,
        session_id: d.session_id,
        action_url: d.action_url,
        action_label: d.action_label,
        expires_at: d.expires_at,
      })
      if (d.phase === 'done') {
        setTimeout(() => remove(d.session_id ?? d.extension_id ?? d.npm_package), 2000)
      }
    })
    const offDone = Events.On('cli.install.done', (event: { data: DoneEvent }) => {
      const ext = event.data.extension
      if (ext?.name) {
        setTimeout(() => remove(ext.id ?? ext.name), 2000)
      }
    })
    return () => {
      offProgress()
      offDone()
    }
  }, [upsert, remove])
}
