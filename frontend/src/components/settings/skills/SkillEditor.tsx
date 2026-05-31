import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { SkillFull } from '@/types/skills'

type Props = {
  initial: Partial<SkillFull>
  readOnly?: boolean
  saving?: boolean
  nameError?: string | null
  dangerLabel?: string
  onDangerClick?: () => void
  onSubmit: (draft: { name: string; description: string; body: string }) => Promise<void> | void
  onCancel: () => void
}

export function SkillEditor({
  initial,
  readOnly = false,
  saving: savingProp = false,
  nameError = null,
  dangerLabel,
  onDangerClick,
  onSubmit,
  onCancel,
}: Props) {
  const { t } = useTranslation()
  const [name, setName] = useState(initial.name ?? '')
  const [description, setDescription] = useState(initial.description ?? '')
  const [body, setBody] = useState(initial.body ?? '')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    setName(initial.name ?? '')
    setDescription(initial.description ?? '')
    setBody(initial.body ?? '')
  }, [initial.name, initial.description, initial.body])

  const handleSave = async () => {
    if (savingProp) return
    setSaving(true)
    try {
      await onSubmit({ name, description, body })
    } finally {
      setSaving(false)
    }
  }

  const busy = saving || savingProp

  return (
    <div className="flex h-full flex-col space-y-5 pb-4 pt-2">
      <div className="shrink-0 space-y-1.5">
        <label className="text-sm font-medium text-foreground">
          {t('settingsPage.skills.editor.name')}
        </label>
        <input
          value={name}
          disabled={readOnly}
          onChange={(e) => setName(e.target.value)}
          className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary disabled:cursor-not-allowed disabled:opacity-70"
        />
        {nameError ? <p className="text-xs text-destructive">{nameError}</p> : null}
      </div>
      <div className="shrink-0 space-y-1.5">
        <label className="text-sm font-medium text-foreground">
          {t('settingsPage.skills.editor.description')}
        </label>
        <textarea
          rows={2}
          value={description}
          disabled={readOnly}
          onChange={(e) => setDescription(e.target.value)}
          className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary disabled:cursor-not-allowed disabled:opacity-70"
        />
      </div>
      <div className="flex min-h-0 flex-1 flex-col space-y-1.5">
        <label className="text-sm font-medium text-foreground">
          {t('settingsPage.skills.editor.body')}
        </label>
        <textarea
          value={body}
          disabled={readOnly}
          onChange={(e) => setBody(e.target.value)}
          className="min-h-[260px] flex-1 rounded-xl border border-border bg-background p-3 font-mono text-xs outline-none focus:border-primary disabled:cursor-not-allowed disabled:opacity-70"
        />
      </div>
      {!readOnly && (
        <div className="flex items-center justify-between gap-3">
          <div>
            {dangerLabel && onDangerClick ? (
              <button
                type="button"
                onClick={onDangerClick}
                className="rounded-xl border border-red-300 px-4 py-2 text-sm text-red-500 transition-colors hover:bg-red-500/10"
              >
                {dangerLabel}
              </button>
            ) : null}
          </div>
          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={onCancel}
              className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
            >
              {t('settingsPage.skills.editor.cancel')}
            </button>
            <button
              type="button"
              onClick={handleSave}
              disabled={busy}
              className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground transition-opacity disabled:opacity-40"
            >
              {busy ? '...' : t('settingsPage.skills.editor.save')}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
