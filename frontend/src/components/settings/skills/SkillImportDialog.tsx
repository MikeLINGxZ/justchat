import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'
import { ImportSkillInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'
import { SelectFileInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto/models'

type Props = {
  open: boolean
  onClose: () => void
  onImported: () => void
}

export function SkillImportDialog({ open, onClose, onImported }: Props) {
  const { t } = useTranslation()
  const [filePath, setFilePath] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  if (!open) return null

  const handleSelectFile = async () => {
    const result = await FileBinding.SelectFile(new SelectFileInput())
    if (result?.file_path) {
      setFilePath(result.file_path)
      setError(null)
    }
  }

  const handleSubmit = async () => {
    setSubmitting(true)
    setError(null)
    try {
      await SkillsBinding.ImportSkill(new ImportSkillInput({ path: filePath }))
      setFilePath('')
      onImported()
      onClose()
    } catch (err: unknown) {
      setError((err as Error).message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-[600px] max-w-[90vw] rounded-md bg-card p-4 shadow-xl">
        <h2 className="mb-3 text-sm font-semibold">{t('settingsPage.skills.importDialog.title')}</h2>
        <p className="mb-2 text-xs text-muted-foreground">{t('settingsPage.skills.importDialog.hint')}</p>
        <div className="flex items-center gap-2">
          <input
            className="flex-1 rounded-md border px-2 py-1.5 text-xs"
            readOnly
            value={filePath}
            placeholder={t('settingsPage.skills.importDialog.placeholder')}
          />
          <button
            type="button"
            className="shrink-0 rounded-md border px-3 py-1.5 text-sm hover:bg-muted"
            onClick={() => { void handleSelectFile() }}
          >
            {t('settingsPage.skills.import')}
          </button>
        </div>
        {error && <p className="mt-2 text-xs text-red-600">{error}</p>}
        <div className="mt-3 flex justify-end gap-2">
          <button type="button" className="rounded-md border px-3 py-1 text-sm" onClick={onClose}>
            {t('settingsPage.skills.editor.cancel')}
          </button>
          <button
            type="button"
            className="rounded-md bg-primary px-3 py-1 text-sm text-primary-foreground disabled:opacity-50"
            onClick={() => { void handleSubmit() }}
            disabled={submitting || !filePath}
          >
            {t('settingsPage.skills.importDialog.submit')}
          </button>
        </div>
      </div>
    </div>
  )
}
