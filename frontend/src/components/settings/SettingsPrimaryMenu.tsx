import { Blocks, Brain, Info, Package2, SlidersHorizontal, Sparkles } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import type { SettingsPrimaryTab } from '@/types/settings'

export function SettingsPrimaryMenu(props: {
  activeTab: SettingsPrimaryTab
  onChange: (tab: SettingsPrimaryTab) => void
  hideHeader?: boolean
  isH5?: boolean
}) {
  const { t } = useTranslation()
  const menuItems: { key: SettingsPrimaryTab; label: string; icon: typeof SlidersHorizontal }[] = [
    { key: 'general', label: t('settingsPage.primary.general'), icon: SlidersHorizontal },
    { key: 'providers', label: t('settingsPage.primary.providers'), icon: Package2 },
    { key: 'plugins', label: t('settingsPage.primary.plugins'), icon: Blocks },
    { key: 'skills', label: t('settingsPage.primary.skills'), icon: Sparkles },
    { key: 'memory', label: t('settingsPage.primary.memory'), icon: Brain },
    { key: 'about', label: t('settingsPage.primary.about'), icon: Info },
  ]

  return (
    <aside className="flex h-full w-64 shrink-0 flex-col border-r border-border/50 bg-card/50 px-3 pb-3 pt-[var(--settings-top-padding)] max-md:w-full">
      {!props.hideHeader && (
        <div className="px-2.5 pb-4 pt-2">
          <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">Lemontea</p>
          <h1 className="mt-2 text-2xl font-semibold text-foreground">{t('settings.settings')}</h1>
        </div>
      )}

      <nav className="space-y-1">
        {menuItems.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            type="button"
            onClick={() => props.onChange(key)}
            className={cn(
              'flex w-full items-center gap-3 rounded-2xl px-3 py-2.5 text-left text-sm transition-colors',
              !props.isH5 && props.activeTab === key
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-accent hover:text-foreground'
            )}
          >
            <Icon size={17} />
            <span className="font-medium">{label}</span>
          </button>
        ))}
      </nav>
    </aside>
  )
}
