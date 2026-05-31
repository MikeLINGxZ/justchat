import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

type Tab = 'smart' | 'npm'

type Props = {
  open: boolean
  onClose: () => void
  onNpmSubmit: (npmPackage: string, name: string) => void
  onSmartSubmit: (content: string) => void
}

export function CliInstallModal({ open, onClose, onNpmSubmit, onSmartSubmit }: Props) {
  const { t } = useTranslation()
  const [tab, setTab] = useState<Tab>('smart')
  const [content, setContent] = useState('')
  const [npmPkg, setNpmPkg] = useState('')
  const [name, setName] = useState('')

  useEffect(() => {
    if (!open) {
      setContent('')
      setNpmPkg('')
      setName('')
      setTab('smart')
    }
  }, [open])

  if (!open) return null

  const handleSmartSubmit = () => {
    if (!content.trim()) return
    onSmartSubmit(content.trim())
    onClose()
  }

  const handleNpmSubmit = () => {
    if (!npmPkg.trim()) return
    onNpmSubmit(npmPkg.trim(), name.trim())
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button type="button" className="absolute inset-0 bg-black/35" onClick={onClose} />
      <div role="dialog" className="relative z-10 w-full max-w-lg rounded-3xl border bg-background p-5 shadow-2xl">
        <h3 className="mb-4 text-base font-semibold">{t('settingsPage.plugins.cliInstallTitle')}</h3>

        {/* Tab switcher */}
        <div className="mb-4 flex rounded-lg bg-muted p-0.5">
          <button
            type="button"
            onClick={() => setTab('smart')}
            className={cn(
              'flex-1 rounded-md py-1.5 text-sm font-medium transition-colors',
              tab === 'smart'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            {t('settingsPage.plugins.cliSmartInstall')}
          </button>
          <button
            type="button"
            onClick={() => setTab('npm')}
            className={cn(
              'flex-1 rounded-md py-1.5 text-sm font-medium transition-colors',
              tab === 'npm'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            {t('settingsPage.plugins.cliFromNpm')}
          </button>
        </div>

        {tab === 'smart' && (
          <>
            <p className="mb-3 text-xs text-muted-foreground">{t('settingsPage.plugins.cliFromDocsHint')}</p>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder={t('settingsPage.plugins.cliDocsPlaceholder')}
              rows={7}
              className="mb-4 w-full resize-none rounded-lg border bg-background px-3 py-2 text-sm"
              autoFocus
            />
            <div className="flex justify-end gap-2">
              <button type="button" onClick={onClose} className="rounded-xl border px-4 py-2 text-sm">
                {t('settingsPage.plugins.cliCancel')}
              </button>
              <button
                type="button"
                onClick={handleSmartSubmit}
                disabled={!content.trim()}
                className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
              >
                {t('settingsPage.plugins.cliInstall')}
              </button>
            </div>
          </>
        )}

        {tab === 'npm' && (
          <>
            <input
              value={npmPkg}
              onChange={(e) => setNpmPkg(e.target.value)}
              placeholder={t('settingsPage.plugins.cliNpmPlaceholder')}
              className="mb-3 w-full rounded-lg border bg-background px-3 py-2 text-sm"
              autoFocus
              onKeyDown={(e) => { if (e.key === 'Enter') handleNpmSubmit() }}
            />
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settingsPage.plugins.cliNamePlaceholder')}
              className="mb-4 w-full rounded-lg border bg-background px-3 py-2 text-sm"
              onKeyDown={(e) => { if (e.key === 'Enter') handleNpmSubmit() }}
            />
            <div className="flex justify-end gap-2">
              <button type="button" onClick={onClose} className="rounded-xl border px-4 py-2 text-sm">
                {t('settingsPage.plugins.cliCancel')}
              </button>
              <button
                type="button"
                onClick={handleNpmSubmit}
                disabled={!npmPkg.trim()}
                className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
              >
                {t('settingsPage.plugins.cliInstall')}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
