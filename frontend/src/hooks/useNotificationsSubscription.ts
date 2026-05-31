import { useEffect } from 'react'
import { Events } from '@wailsio/runtime'
import { Notification as NotificationBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification'
import { useNotificationsStore } from '@/store/notificationsStore'
import type { NotificationItem } from '@/types/notifications'

type WailsEvent<T> = {
  data: T | T[]
}

function eventData<T>(event: WailsEvent<T>): T | undefined {
  return Array.isArray(event.data) ? event.data[0] : event.data
}

export function useNotificationsSubscription() {
  const setItems = useNotificationsStore((state) => state.setItems)
  const prepend = useNotificationsStore((state) => state.prepend)
  const remove = useNotificationsStore((state) => state.remove)

  useEffect(() => {
    void NotificationBinding.ListNotifications({ include_resolved: false }).then((result) => {
      setItems((result?.items ?? []) as NotificationItem[])
    })

    const offCreate = Events.On('notification.created', (event: WailsEvent<NotificationItem>) => {
      const data = eventData(event)
      if (data) prepend(data)
    })
    const offResolve = Events.On('notification.resolved', (event: WailsEvent<{ id: number }>) => {
      const data = eventData(event)
      if (data) remove(data.id)
    })

    return () => {
      offCreate()
      offResolve()
    }
  }, [prepend, remove, setItems])
}
