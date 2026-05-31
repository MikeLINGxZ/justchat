import { createContext, type ReactNode, useEffect, useRef, useState } from 'react'
import { SettingsPrimaryMenu } from '@/components/settings/SettingsPrimaryMenu'
import { isH5Width } from '@/lib/responsive'
import { useSettingsStore } from '@/store/settingsStore'

export const SettingsMenuContext = createContext<{
  isH5: boolean
  onOpenMenu?: () => void
  submenuCollapsed: boolean
  onCollapseSubmenu: () => void
  onOpenSubmenu?: () => void
}>({
  isH5: false,
  submenuCollapsed: false,
  onCollapseSubmenu: () => {},
})

export function SettingsShell(props: { children: ReactNode }) {
  const activeTab = useSettingsStore((state) => state.activeTab)
  const setActiveTab = useSettingsStore((state) => state.setActiveTab)
  const [menuCollapsed, setMenuCollapsed] = useState(false)
  const [submenuCollapsed, setSubmenuCollapsed] = useState(false)
  const [isH5, setIsH5] = useState(false)
  const autoCollapsedRef = useRef(false)
  const wasNarrowRef = useRef(false)

  useEffect(() => {
    const handleResize = () => {
      const narrow = isH5Width(window.innerWidth)
      setIsH5(narrow)

      if (narrow && !wasNarrowRef.current) {
        autoCollapsedRef.current = true
        setMenuCollapsed(true)
        setSubmenuCollapsed(true)
      }
      if (!narrow && wasNarrowRef.current && autoCollapsedRef.current) {
        autoCollapsedRef.current = false
        setMenuCollapsed(false)
        setSubmenuCollapsed(false)
      }

      wasNarrowRef.current = narrow
    }

    handleResize()
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  const handleChangeTab = (tab: typeof activeTab) => {
    setActiveTab(tab)
    if (isH5) {
      setMenuCollapsed(true)
      setSubmenuCollapsed(false)
    }
  }

  const onOpenMenu = (menuCollapsed && isH5)
    ? () => { autoCollapsedRef.current = false; setMenuCollapsed(false) }
    : undefined

  const onOpenSubmenu = (submenuCollapsed && isH5)
    ? () => setSubmenuCollapsed(false)
    : undefined

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground [--settings-top-padding:2.5rem]">
      {!menuCollapsed && !isH5 && <SettingsPrimaryMenu activeTab={activeTab} onChange={handleChangeTab} />}
      {!menuCollapsed && isH5 ? (
        <main className="relative flex h-full w-full flex-col bg-background">
          <SettingsPrimaryMenu activeTab={activeTab} onChange={handleChangeTab} isH5 />
        </main>
      ) : (
        <main className="relative flex min-h-0 min-w-0 flex-1 flex-col">
          <SettingsMenuContext.Provider value={{
            isH5,
            onOpenMenu,
            submenuCollapsed,
            onCollapseSubmenu: () => setSubmenuCollapsed(true),
            onOpenSubmenu,
          }}>
            <div className="flex min-h-0 min-w-0 flex-1 flex-col">
              {props.children}
            </div>
          </SettingsMenuContext.Provider>
        </main>
      )}
    </div>
  )
}
