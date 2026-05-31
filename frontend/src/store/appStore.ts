import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Theme, FontSize, Language } from '../types'
import i18n from '../i18n'

interface AppStore {
  theme: Theme
  fontSize: FontSize
  language: Language
  enabledToolIds: string[]
  setTheme: (theme: Theme) => void
  setFontSize: (size: FontSize) => void
  setLanguage: (lang: Language) => void
  setEnabledToolIds: (ids: string[]) => void
}

export const useAppStore = create<AppStore>()(
  persist(
    (set) => ({
      theme: 'auto',
      fontSize: 'xl',
      language: 'zh-CN',
      enabledToolIds: [] as string[],
      setTheme: (theme) => set({ theme }),
      setFontSize: (fontSize) => set({ fontSize }),
      setLanguage: (language) => {
        i18n.changeLanguage(language)
        set({ language })
      },
      setEnabledToolIds: (enabledToolIds) => set({ enabledToolIds }),
    }),
    {
      name: 'lemontea-app-settings',
      onRehydrateStorage: () => (state) => {
        if (state?.language) {
          i18n.changeLanguage(state.language)
        }
      },
    }
  )
)
