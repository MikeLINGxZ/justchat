import { create } from 'zustand'
import type { NotificationItem } from '@/types/notifications'

type NotificationsStore = {
  items: NotificationItem[]
  setItems: (items: NotificationItem[]) => void
  prepend: (item: NotificationItem) => void
  remove: (id: number) => void
}

export const useNotificationsStore = create<NotificationsStore>((set) => ({
  items: [],
  setItems: (items) => set({ items }),
  prepend: (item) => set((state) => ({
    items: [item, ...state.items.filter((current) => current.id !== item.id)],
  })),
  remove: (id) => set((state) => ({
    items: state.items.filter((item) => item.id !== id),
  })),
}))
