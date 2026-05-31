import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

type Props = {
  open: boolean
  onClose: () => void
  onSubmit: (npmPackage: string, name: string) => void
}

export function CliNpmInstallModal({ open, onClose, onSubmit }: Props) {
  const { t } = useTranslation()
  const [npmPkg, setNpmPkg] = useState('')
  const [name, setName] = useState('')

  useEffect(() => {
    if (!open) {
      setNpmPkg('')
      setName('')
    }
  }, [open])

  if (!open) return null

  const handleSubmit = () => {
    if (!npmPkg.trim()) return
    onSubmit(npmPkg.trim(), name.trim())
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button type="button" className="absolute inset-0 bg-black/35" onClick={onClose} />
      <div role="dialog" className="relative z-10 w-full max-w-md rounded-3xl border bg-background p-5 shadow-2xl">
        <h3 className="mb-4 text-base font-semibold">{t('settingsPage.plugins.cliFromNpm')}</h3>
        <input
          value={npmPkg}
          onChange={(e) => setNpmPkg(e.target.value)}
          placeholder={t('settingsPage.plugins.cliNpmPlaceholder')}
          className="mb-3 w-full rounded-lg border bg-background px-3 py-2 text-sm"
          autoFocus
          onKeyDown={(e) => { if (e.key === 'Enter') handleSubmit() }}
        />
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder={t('settingsPage.plugins.cliNamePlaceholder')}
          className="mb-4 w-full rounded-lg border bg-background px-3 py-2 text-sm"
          onKeyDown={(e) => { if (e.key === 'Enter') handleSubmit() }}
        />
        <div className="flex justify-end gap-2">
          <button type="button" onClick={onClose} className="rounded-xl border px-4 py-2 text-sm">
            {t('settingsPage.plugins.cliCancel')}
          </button>
          <button
            type="button"
            onClick={handleSubmit}
            disabled={!npmPkg.trim()}
            className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
          >
            {t('settingsPage.plugins.cliInstall')}
          </button>
        </div>
      </div>
    </div>
  )
}
