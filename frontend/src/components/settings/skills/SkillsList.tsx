import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, Upload } from 'lucide-react'
import { Window } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window'
import { ConfirmDialog } from '@/components/settings/common/ConfirmDialog'
import { useSettingsStore } from '@/store/settingsStore'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { DeleteSkillInput, ListSkillsInput, ToggleSkillInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'
import { SkillsListItem } from './SkillsListItem'
import { SkillImportDialog } from './SkillImportDialog'
import type { SkillItem } from '@/types/skills'

export function SkillsList() {
  const { t } = useTranslation()
  const skills = useSettingsStore((s) => s.skills)
  const setSkills = useSettingsStore((s) => s.setSkills)
  const selectedSkillName = useSettingsStore((s) => s.selectedSkillName)
  const setSelectedSkillName = useSettingsStore((s) => s.setSelectedSkillName)
  const updateSkill = useSettingsStore((s) => s.updateSkill)
  const removeSkill = useSettingsStore((s) => s.removeSkill)
  const [importOpen, setImportOpen] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<SkillItem | null>(null)
  const [deleting, setDeleting] = useState(false)

  const refresh = useCallback(async () => {
    const result = await SkillsBinding.ListSkills(new ListSkillsInput())
    setSkills((result?.skills ?? []) as unknown as SkillItem[])
  }, [setSkills])

  useEffect(() => {
    void refresh()
  }, [refresh])

  useEffect(() => {
    const handleFocus = () => {
      void refresh()
    }
    window.addEventListener('focus', handleFocus)
    return () => window.removeEventListener('focus', handleFocus)
  }, [refresh])

  const handleToggle = useCallback(async (item: SkillItem) => {
    const nextDisabled = !item.disabled
    const optimistic: SkillItem = { ...item, disabled: nextDisabled }
    updateSkill(optimistic)
    try {
      const result = await SkillsBinding.ToggleSkill(new ToggleSkillInput({ name: item.name, disabled: nextDisabled }))
      if (result?.skill) {
        updateSkill(result.skill as unknown as SkillItem)
      }
    } catch {
      updateSkill(item)
    }
  }, [updateSkill])

  const handleDelete = useCallback(async () => {
    if (!deleteTarget) return
    if (deleting) return
    setDeleting(true)
    try {
      await SkillsBinding.DeleteSkill(new DeleteSkillInput({ name: deleteTarget.name }))
      removeSkill(deleteTarget.name)
      if (selectedSkillName === deleteTarget.name) {
        setSelectedSkillName(null)
      }
      setDeleteTarget(null)
    } finally {
      setDeleting(false)
    }
  }, [deleteTarget, deleting, removeSkill, selectedSkillName, setSelectedSkillName])

  const handleStartNew = () => { void Window.OpenAddSkill({}) }

  return (
    <div className="flex h-full flex-col gap-2 p-3">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold">{t('settingsPage.skills.title')}</h3>
        <div className="flex gap-1">
          <button
            type="button"
            aria-label={t('settingsPage.skills.addNew')}
            onClick={handleStartNew}
            className="rounded-md p-1 hover:bg-muted"
          >
            <Plus size={16} />
          </button>
          <button
            type="button"
            aria-label={t('settingsPage.skills.import')}
            onClick={() => setImportOpen(true)}
            className="rounded-md p-1 hover:bg-muted"
          >
            <Upload size={16} />
          </button>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto">
        {skills.length === 0 ? (
          <p className="px-2 py-6 text-center text-sm opacity-60">{t('settingsPage.skills.empty')}</p>
        ) : (
          <ul className="flex flex-col gap-1">
            {skills.map((s) => (
              <li key={s.name}>
                <SkillsListItem
                  item={s}
                  selected={selectedSkillName === s.name}
                  onSelect={setSelectedSkillName}
                  onToggle={handleToggle}
                  onDelete={setDeleteTarget}
                />
              </li>
            ))}
          </ul>
        )}
      </div>
      <SkillImportDialog open={importOpen} onClose={() => setImportOpen(false)} onImported={() => { void refresh() }} />
      <ConfirmDialog
        open={deleteTarget !== null}
        title={t('settingsPage.skills.deleteConfirm')}
        description={deleteTarget?.name}
        confirmLabel={t('settingsPage.skills.delete')}
        cancelLabel={t('settingsPage.skills.editor.cancel')}
        confirmTone="danger"
        busy={deleting}
        onConfirm={() => { void handleDelete() }}
        onCancel={() => {
          if (!deleting) setDeleteTarget(null)
        }}
      />
    </div>
  )
}
