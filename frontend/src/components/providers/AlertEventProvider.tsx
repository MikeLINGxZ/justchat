import { useEffect } from 'react'
import { Events } from '@wailsio/runtime'
import { APP_ALERT_EVENT, normalizeAlertEventPayload } from '@/alert/event'
import { useAlertStore } from '@/alert/store'

type AlertEventProviderProps = {
  children: React.ReactNode
}

export function AlertEventProvider({ children }: AlertEventProviderProps) {
  useEffect(() => {
    const off = Events.On(APP_ALERT_EVENT, (event: { data: unknown }) => {
      const payload = normalizeAlertEventPayload(event.data)
      if (!payload) {
        return
      }

      useAlertStore.getState().pushAlert(payload)
    })

    return () => off()
  }, [])

  return <>{children}</>
}
