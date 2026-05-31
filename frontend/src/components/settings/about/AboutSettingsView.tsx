import { useTranslation } from 'react-i18next'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'

export function AboutSettingsView(props: { version: string }) {
  const { t } = useTranslation()

  return (
    <SettingsContentLayout
      header={
        <SettingsPanelHeader
          title="Lemontea"
          description={t('settingsPage.about.description')}
          aside={<span className="rounded-full bg-primary/10 px-3 py-1 text-sm text-primary">{props.version}</span>}
        />
      }
    />
  )
}
