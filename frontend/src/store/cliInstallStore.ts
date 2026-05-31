import { create } from 'zustand'
import type { CliInstallItem } from '@/types/cliInstall'

type State = {
  items: CliInstallItem[]
  upsert: (item: CliInstallItem) => void
  remove: (id: string | number) => void
  clear: () => void
}

export const useCliInstallStore = create<State>((set) => ({
  items: [],
  upsert: (item) =>
    set((s) => {
      const idx = s.items.findIndex((i) =>
        (item.session_id != null && i.session_id === item.session_id) ||
        (!!item.extension_id && i.extension_id === item.extension_id) ||
        (!!item.npm_package && i.npm_package === item.npm_package),
      )
      if (idx >= 0) {
        const next = [...s.items]
        next[idx] = { ...next[idx], ...item }
        return { items: next }
      }
      return { items: [...s.items, item] }
    }),
  remove: (id) =>
    set((s) => ({
      items: s.items.filter((i) =>
        typeof id === 'number'
          ? i.session_id !== id
          : i.npm_package !== id && i.extension_id !== id,
      ),
    })),
  clear: () => set({ items: [] }),
}))
