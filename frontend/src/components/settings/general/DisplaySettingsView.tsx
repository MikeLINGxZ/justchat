import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsDirtyGuard } from '@/components/settings/common/SettingsDirtyGuard'
import { SettingsFieldRow } from '@/components/settings/common/SettingsFieldRow'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { cn } from '@/lib/utils'
import type { FontSize } from '@/types'

export function DisplaySettingsView(props: {
  value: FontSize
  draft: FontSize
  onDraftChange: (next: FontSize) => void
  onApply: () => void
  onReset: () => void
}) {
  const { t } = useTranslation()
  const [previewFontSize, setPreviewFontSize] = useState<FontSize>(props.draft)
  const options: { value: FontSize; label: string; preview: string }[] = [
    { value: 'xs', label: t('settingsPage.fontSizes.xs'), preview: t('settingsPage.fontPreview.xs') },
    { value: 'sm', label: t('settingsPage.fontSizes.sm'), preview: t('settingsPage.fontPreview.sm') },
    { value: 'md', label: t('settingsPage.fontSizes.md'), preview: t('settingsPage.fontPreview.md') },
    { value: 'lg', label: t('settingsPage.fontSizes.lg'), preview: t('settingsPage.fontPreview.lg') },
    { value: 'xl', label: t('settingsPage.fontSizes.xl'), preview: t('settingsPage.fontPreview.xl') },
  ]

  useEffect(() => {
    setPreviewFontSize(props.draft)
  }, [props.draft])

  const dirty = props.value !== previewFontSize

  return (
    <SettingsDirtyGuard dirty={dirty}>
      <SettingsContentLayout
        header={
          <SettingsPanelHeader
            title={t('settingsPage.general.display.title')}
            description={t('settingsPage.general.display.description')}
          />
        }
        footprint={
          <SettingsActionBar
            primaryLabel={t('settingsPage.actions.apply')}
            secondaryLabel={t('settingsPage.actions.reset')}
            primaryDisabled={!dirty}
            onPrimaryClick={props.onApply}
            onSecondaryClick={props.onReset}
          />
        }
      >
        <div className="space-y-4 pt-6">
          <SettingsFieldRow
            label={t('settingsPage.general.display.textSizeLabel')}
            description={t('settingsPage.general.display.textSizeDescription')}
          >
            <div className="grid gap-2 sm:grid-cols-2 xl:grid-cols-5">
              {options.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  aria-label={option.label}
                  onClick={() => {
                    setPreviewFontSize(option.value)
                    props.onDraftChange(option.value)
                  }}
                  className={cn(
                    'flex min-h-20 flex-col justify-between rounded-xl border px-2.5 py-2.5 text-left transition-colors',
                    previewFontSize === option.value
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-border bg-background hover:bg-accent'
                  )}
                >
                  <span className="block text-[11px] font-semibold tracking-wide">{option.label}</span>
                  <span className="mt-1.5 block text-lg leading-none">{option.preview}</span>
                </button>
              ))}
            </div>
          </SettingsFieldRow>

          <SettingsFieldRow
            label={t('settingsPage.general.display.previewLabel')}
            description={t('settingsPage.general.display.previewDescription')}
          >
            <div className="rounded-2xl border border-border/70 bg-card/60 p-5">
              <p className={cn(
                'text-foreground',
                previewFontSize === 'xs' && 'text-xs',
                previewFontSize === 'sm' && 'text-sm',
                previewFontSize === 'md' && 'text-base',
                previewFontSize === 'lg' && 'text-lg',
                previewFontSize === 'xl' && 'text-xl'
              )}
              >
                {t('settingsPage.general.display.previewBody')}
              </p>
            </div>
          </SettingsFieldRow>
        </div>
      </SettingsContentLayout>
    </SettingsDirtyGuard>
  )
}
