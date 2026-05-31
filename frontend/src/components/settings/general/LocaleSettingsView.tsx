import { useTranslation } from 'react-i18next'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsDirtyGuard } from '@/components/settings/common/SettingsDirtyGuard'
import { SettingsFieldRow } from '@/components/settings/common/SettingsFieldRow'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { SettingsSelect } from '@/components/settings/common/SettingsSelect'
import type { Language } from '@/types'
import type { SettingsOption } from '@/types/settings'

export function LocaleSettingsView(props: {
  locale: string
  language: Language
  regions: SettingsOption[]
  languages: SettingsOption[]
  dirty: boolean
  onLocaleChange: (value: string) => void
  onLanguageChange: (value: Language) => void
  onApply: () => void
}) {
  const { t } = useTranslation()
  const regionOptions = props.regions.map((item) => ({ value: item.id, label: item.name, icon: item.icon }))
  const languageOptions = props.languages.map((item) => ({ value: item.id, label: item.name }))

  return (
    <SettingsDirtyGuard dirty={props.dirty}>
      <SettingsContentLayout
        header={
          <SettingsPanelHeader
            title={t('settingsPage.general.locale.title')}
            description={t('settingsPage.general.locale.description')}
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
            label={t('settingsPage.general.locale.regionLabel')}
            description={t('settingsPage.general.locale.regionDescription')}
          >
            <SettingsSelect
              ariaLabel={t('settingsPage.general.locale.regionLabel')}
              value={props.locale}
              options={[...regionOptions]}
              onChange={props.onLocaleChange}
            />
          </SettingsFieldRow>

          <SettingsFieldRow
            label={t('settingsPage.general.locale.languageLabel')}
            description={t('settingsPage.general.locale.languageDescription')}
          >
            <SettingsSelect
              ariaLabel={t('settingsPage.general.locale.languageLabel')}
              value={props.language}
              options={[...languageOptions]}
              onChange={(value) => props.onLanguageChange(value as Language)}
            />
          </SettingsFieldRow>
        </div>
      </SettingsContentLayout>
    </SettingsDirtyGuard>
  )
}
