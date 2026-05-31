import { Plus } from 'lucide-react'
import { useContext } from 'react'
import { useTranslation } from 'react-i18next'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { ProviderListItem } from '@/components/settings/providers/ProviderListItem'
import type { ProviderItem } from '@/types/settings'

export function ProviderList(props: {
  items: ProviderItem[]
  selectedId: number | null
  onSelect: (id: number) => void
  onCreate: () => void
  onToggleEnable: (item: ProviderItem) => void
  onSetDefault: (id: number) => void
  onDelete: (id: number) => void
}) {
  const { t } = useTranslation()
  const { isH5, onCollapseSubmenu } = useContext(SettingsMenuContext)

  return (
    <div className="flex h-full flex-col">
      <div className="mb-4 flex items-center justify-between gap-3 px-2">
        <h2 className="text-lg font-semibold text-foreground">{t('settingsPage.primary.providers')}</h2>
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={props.onCreate}
            className="inline-flex h-8 items-center gap-1.5 rounded-lg px-2 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
            aria-label={t('settingsPage.providers.add')}
          >
            <Plus size={14} />
            <span>{t('settingsPage.providers.add')}</span>
          </button>
        </div>
      </div>

      <div className="space-y-2">
        {props.items.map((item) => (
          <ProviderListItem
            key={item.id}
            item={item}
            selected={!isH5 && item.id === props.selectedId}
            onSelect={(id) => { props.onSelect(id); onCollapseSubmenu() }}
            onToggleEnable={props.onToggleEnable}
            onSetDefault={props.onSetDefault}
            onDelete={props.onDelete}
          />
        ))}
      </div>
    </div>
  )
}
