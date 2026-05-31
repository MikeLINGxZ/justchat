import { useCallback, useContext, useEffect, useState } from 'react'
import { LockKeyhole } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { ConfirmDialog } from '@/components/settings/common/ConfirmDialog'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { toSkillName } from '@/lib/skills'
import { useSettingsStore } from '@/store/settingsStore'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { DeleteSkillInput, GetSkillInput, UpdateSkillInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'
import { SkillEditor } from './SkillEditor'
import { cn } from '@/lib/utils'
import type { SkillFull, SkillItem } from '@/types/skills'

export function SkillsDetailPane() {
  const { t } = useTranslation()
  const { isH5 } = useContext(SettingsMenuContext)
  const selectedSkillName = useSettingsStore((s) => s.selectedSkillName)
  const updateSkill = useSettingsStore((s) => s.updateSkill)
  const removeSkill = useSettingsStore((s) => s.removeSkill)
  const setSelectedSkillName = useSettingsStore((s) => s.setSelectedSkillName)
  const [fullSkill, setFullSkill] = useState<SkillFull | null>(null)
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [nameError, setNameError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    if (!selectedSkillName) {
      setFullSkill(null)
      return
    }
    const result = await SkillsBinding.GetSkill(new GetSkillInput({ name: selectedSkillName }))
    if (result?.skill) {
      setFullSkill(result.skill as unknown as SkillFull)
      setNameError(null)
    }
  }, [selectedSkillName])

  useEffect(() => { void refresh() }, [refresh])

  const handleUpdate = useCallback(async (draft: { name: string; description: string; body: string }) => {
    if (!fullSkill) return
    const nextName = toSkillName(draft.name)
    if (!nextName) {
      setNameError(t('settingsPage.skills.editor.nameInvalid'))
      return
    }
    setNameError(null)
    const result = await SkillsBinding.UpdateSkill(new UpdateSkillInput({
      name: fullSkill.name,
      new_name: nextName,
      description: draft.description,
      body: draft.body,
    }))
    if (result?.skill) {
      const nextSkill = result.skill as unknown as SkillFull
      if (nextSkill.name !== fullSkill.name) {
        removeSkill(fullSkill.name)
      }
      updateSkill(nextSkill as SkillItem)
      setFullSkill(nextSkill)
      setSelectedSkillName(nextSkill.name)
    }
  }, [fullSkill, removeSkill, setSelectedSkillName, t, updateSkill])

  const handleDelete = useCallback(async () => {
    if (!fullSkill) return
    if (deleting) return
    setDeleting(true)
    try {
      await SkillsBinding.DeleteSkill(new DeleteSkillInput({ name: fullSkill.name }))
      removeSkill(fullSkill.name)
      setSelectedSkillName(null)
      setConfirmDeleteOpen(false)
    } finally {
      setDeleting(false)
    }
  }, [deleting, fullSkill, removeSkill, setSelectedSkillName])

  if (!fullSkill) {
    return (
      <div className={cn('flex h-full items-center justify-center text-sm opacity-60', isH5 ? 'p-5' : 'p-10')}>
        {t('settingsPage.skills.selectHint')}
      </div>
    )
  }

  const readOnly = fullSkill.source === 'builtin'
  const sourceLabel = t(`settingsPage.skills.source${fullSkill.source === 'builtin' ? 'Builtin' : fullSkill.source === 'user' ? 'User' : 'AI'}`)

  return (
    <div className="flex h-full flex-col">
      <div className={cn('pb-4 pt-4', isH5 ? 'px-5' : 'px-10')}>
        <SettingsPanelHeader
          title={fullSkill.name}
          badge={
            <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs font-medium text-primary">
              {sourceLabel}
            </span>
          }
        />
      </div>
      {readOnly && (
        <div className={cn(
          'mb-4 flex shrink-0 items-start gap-2 rounded-xl bg-amber-500/10 px-3 py-2.5 text-sm text-amber-700 dark:text-amber-400',
          isH5 ? 'mx-5' : 'mx-10',
        )}>
          <LockKeyhole size={14} className="mt-0.5 shrink-0" />
          <span>
            <span className="font-medium">{t('settingsPage.skills.builtinLockedTitle')}</span>
            {' '}
            {t('settingsPage.skills.builtinLockedMsg')}
          </span>
        </div>
      )}
      <div className={cn('min-h-0 flex-1', isH5 ? 'px-5' : 'px-10')}>
        <SkillEditor
          initial={fullSkill}
          readOnly={readOnly}
          nameError={nameError}
          dangerLabel={t('settingsPage.skills.delete')}
          onDangerClick={() => setConfirmDeleteOpen(true)}
          onSubmit={handleUpdate}
          onCancel={() => setSelectedSkillName(null)}
        />
      </div>
      <ConfirmDialog
        open={confirmDeleteOpen}
        title={t('settingsPage.skills.deleteConfirm')}
        description={fullSkill.name}
        confirmLabel={t('settingsPage.skills.delete')}
        cancelLabel={t('settingsPage.skills.editor.cancel')}
        confirmTone="danger"
        busy={deleting}
        onConfirm={() => { void handleDelete() }}
        onCancel={() => {
          if (!deleting) setConfirmDeleteOpen(false)
        }}
      />
    </div>
  )
}
