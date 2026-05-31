import { CheckCircle2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'

export function AddProviderStepDone() {
  const { t } = useTranslation()
  return (
    <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
      <CheckCircle2 size={56} className="text-primary" />
      <h2 className="text-xl font-semibold text-foreground">{t('settingsPage.addProvider.doneTitle')}</h2>
      <p className="max-w-xs text-sm text-muted-foreground">{t('settingsPage.addProvider.doneMessage')}</p>
    </div>
  )
}
