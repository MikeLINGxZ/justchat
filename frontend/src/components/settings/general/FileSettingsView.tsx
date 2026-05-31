import { FolderOpen } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsDirtyGuard } from '@/components/settings/common/SettingsDirtyGuard'
import { SettingsFieldRow } from '@/components/settings/common/SettingsFieldRow'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'

export function FileSettingsView(props: {
  dataDir: string
  currentDataDir: string
  dirty: boolean
  onSelectFolder: () => void
  onApply: () => void
}) {
  const { t } = useTranslation()

  return (
    <SettingsDirtyGuard dirty={props.dirty}>
      <SettingsContentLayout
        header={
          <SettingsPanelHeader
            title={t('settingsPage.general.file.title')}
            description={t('settingsPage.general.file.description')}
          />
        }
        footprint={
          <SettingsActionBar
            primaryLabel={t('settingsPage.actions.apply')}
            primaryDisabled={!props.dirty}
            onPrimaryClick={props.onApply}
          />
        }
      >
        <div className="space-y-4 pt-6">
          <SettingsFieldRow
            label={t('settingsPage.general.file.dataDirLabel')}
            description={t('settingsPage.general.file.dataDirDescription')}
          >
            <div className="flex items-center gap-3">
              <input
                value={props.dataDir}
                readOnly
                className="w-full rounded-xl border border-border bg-background px-3 py-2"
                placeholder={props.currentDataDir}
              />
              <button
                type="button"
                aria-label="Select folder"
                onClick={props.onSelectFolder}
                className="inline-flex h-10 shrink-0 items-center justify-center rounded-xl border border-border px-3 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
              >
                <FolderOpen size={16} />
              </button>
            </div>
          </SettingsFieldRow>
        </div>
      </SettingsContentLayout>
    </SettingsDirtyGuard>
  )
}
