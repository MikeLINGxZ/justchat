import { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { Check, ChevronDown, Plus, RotateCcw, Search } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Memory as MemoryBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory'
import { Window } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window'
import { CreateMemoryInput, ForgetMemoryInput, ListMemoriesInput, RestoreMemoryInput, UpdateMemoryInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory/memory_dto/models'
import type { MemoryItem } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory/memory_dto/models'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { SettingsSelect } from '@/components/settings/common/SettingsSelect'
import { cn } from '@/lib/utils'

type MemoryDraft = {
  id: number | null
  summary: string
  content: string
  type: string
  target: string
  importance: number
  confidence: number
  pinned: boolean
}

type MemoryController = {
  items: MemoryItem[]
  selected: MemoryItem | null
  selectedId: number | null
  draft: MemoryDraft
  draftDirty: boolean
  query: string
  target: string
  type: string
  includeForgotten: boolean
  loading: boolean
  setSelectedId: (id: number | null) => void
  setDraft: (draft: MemoryDraft) => void
  resetDraft: () => void
  setQuery: (value: string) => void
  setTarget: (value: string) => void
  setType: (value: string) => void
  setIncludeForgotten: (value: boolean) => void
  loadMemories: () => Promise<void>
  saveDraft: () => Promise<void>
  forgetSelected: () => Promise<void>
  restoreSelected: () => Promise<void>
  startNew: () => void
}

const EMPTY_DRAFT: MemoryDraft = {
  id: null,
  summary: '',
  content: '',
  type: 'information',
  target: 'user',
  importance: 70,
  confidence: 85,
  pinned: false,
}

function memoryDraftsEqual(left: MemoryDraft, right: MemoryDraft): boolean {
  return left.id === right.id &&
    left.summary === right.summary &&
    left.content === right.content &&
    left.type === right.type &&
    left.target === right.target &&
    left.importance === right.importance &&
    left.confidence === right.confidence &&
    left.pinned === right.pinned
}

// useMemoryController owns memory list/detail state for the settings split view.
export function useMemoryController(active: boolean): MemoryController {
  const [items, setItems] = useState<MemoryItem[]>([])
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [query, setQuery] = useState('')
  const [target, setTarget] = useState('')
  const [type, setType] = useState('')
  const [includeForgotten, setIncludeForgotten] = useState(false)
  const [loading, setLoading] = useState(false)
  const [draft, setDraft] = useState<MemoryDraft>(EMPTY_DRAFT)
  const [initialDraft, setInitialDraft] = useState<MemoryDraft>(EMPTY_DRAFT)

  const selected = useMemo(
    () => items.find((item) => item.id === selectedId) ?? null,
    [items, selectedId],
  )

  const draftDirty = useMemo(
    () => !memoryDraftsEqual(draft, initialDraft),
    [draft, initialDraft],
  )

  const loadMemories = useCallback(async () => {
    setLoading(true)
    try {
      const listResult = await MemoryBinding.ListMemories(new ListMemoriesInput({
        query,
        target,
        type,
        include_forgotten: includeForgotten,
        offset: 0,
        limit: 100,
      }))
      const nextItems = listResult?.items ?? []
      setItems(nextItems)
      setSelectedId((current) => {
        if (current !== null && nextItems.some((item) => item.id === current)) return current
        return nextItems[0]?.id ?? null
      })
    } finally {
      setLoading(false)
    }
  }, [includeForgotten, query, target, type])

  useEffect(() => {
    if (active) {
      void loadMemories()
    }
  }, [active, loadMemories])

  useEffect(() => {
    if (!active) return
    const handleFocus = () => { void loadMemories() }
    window.addEventListener('focus', handleFocus)
    return () => window.removeEventListener('focus', handleFocus)
  }, [active, loadMemories])

  useEffect(() => {
    if (!selected) {
      setDraft(EMPTY_DRAFT)
      setInitialDraft(EMPTY_DRAFT)
      return
    }
    const nextDraft: MemoryDraft = {
      id: selected.id,
      summary: selected.summary,
      content: selected.content,
      type: selected.type,
      target: selected.target,
      importance: selected.importance,
      confidence: selected.confidence,
      pinned: selected.pinned,
    }
    setDraft(nextDraft)
    setInitialDraft(nextDraft)
  }, [selected])

  const startNew = useCallback(() => {
    setSelectedId(null)
    setDraft(EMPTY_DRAFT)
    setInitialDraft(EMPTY_DRAFT)
  }, [])

  const resetDraft = useCallback(() => {
    setDraft(initialDraft)
  }, [initialDraft])

  const saveDraft = useCallback(async () => {
    if (draft.id === null) {
      await MemoryBinding.CreateMemory(new CreateMemoryInput({
        summary: draft.summary,
        content: draft.content,
        type: draft.type,
        target: draft.target,
        source: 'manual',
        importance: draft.importance,
        confidence: draft.confidence,
        pinned: draft.pinned,
      }))
    } else {
      await MemoryBinding.UpdateMemory(new UpdateMemoryInput({
        id: draft.id,
        summary: draft.summary,
        content: draft.content,
        type: draft.type,
        target: draft.target,
        source: selected?.source ?? 'manual',
        importance: draft.importance,
        confidence: draft.confidence,
        pinned: draft.pinned,
      }))
    }
    await loadMemories()
  }, [draft, loadMemories, selected?.source])

  const forgetSelected = useCallback(async () => {
    if (!selected) return
    await MemoryBinding.ForgetMemory(new ForgetMemoryInput({ id: selected.id }))
    await loadMemories()
  }, [loadMemories, selected])

  const restoreSelected = useCallback(async () => {
    if (!selected) return
    await MemoryBinding.RestoreMemory(new RestoreMemoryInput({ id: selected.id }))
    await loadMemories()
  }, [loadMemories, selected])

  return {
    items,
    selected,
    selectedId,
    draft,
    draftDirty,
    query,
    target,
    type,
    includeForgotten,
    loading,
    setSelectedId,
    setDraft,
    resetDraft,
    setQuery,
    setTarget,
    setType,
    setIncludeForgotten,
    loadMemories,
    saveDraft,
    forgetSelected,
    restoreSelected,
    startNew,
  }
}

// MemoryList renders the left-hand memory navigation list.
export function MemoryList(props: { controller: MemoryController }) {
  const { t } = useTranslation()
  const { isH5, onCollapseSubmenu } = useContext(SettingsMenuContext)
  const controller = props.controller

  return (
    <div className="flex h-full flex-col gap-3 p-3">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-foreground">{t('settingsPage.memory.title')}</h3>
        <button
          type="button"
          aria-label={t('settingsPage.memory.new')}
          onClick={() => {
            void Window.OpenAddMemory({})
            onCollapseSubmenu()
          }}
          className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <Plus size={16} />
        </button>
      </div>

      <label className="flex h-8 items-center gap-2 rounded-lg border border-border/80 bg-card/80 px-2.5 text-xs text-foreground transition-colors focus-within:border-primary">
        <Search size={14} className="text-muted-foreground" />
        <input
          value={controller.query}
          onChange={(event) => controller.setQuery(event.target.value)}
          placeholder={t('settingsPage.memory.search')}
          className="min-w-0 flex-1 appearance-none bg-transparent text-xs outline-none placeholder:text-muted-foreground"
        />
      </label>

      <div className="grid grid-cols-2 gap-2">
        <CompactSelect
          ariaLabel={t('settingsPage.memory.target')}
          value={controller.target}
          onChange={controller.setTarget}
          options={[
            { value: '', label: t('settingsPage.memory.allTargets') },
            { value: 'user', label: t('settingsPage.memory.targetUser') },
            { value: 'memory', label: t('settingsPage.memory.targetMemory') },
          ]}
        />
        <CompactSelect
          ariaLabel={t('settingsPage.memory.type')}
          value={controller.type}
          onChange={controller.setType}
          options={[
            { value: '', label: t('settingsPage.memory.allTypes') },
            { value: 'fact', label: t('settingsPage.memory.typeFact') },
            { value: 'information', label: t('settingsPage.memory.typeInformation') },
            { value: 'event', label: t('settingsPage.memory.typeEvent') },
          ]}
        />
      </div>

      <button
        type="button"
        role="switch"
        aria-checked={controller.includeForgotten}
        onClick={() => controller.setIncludeForgotten(!controller.includeForgotten)}
        className="flex items-center justify-between rounded-xl px-2 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      >
        <span>{t('settingsPage.memory.showForgotten')}</span>
        <span className={cn(
          'relative inline-flex h-[18px] w-[31px] shrink-0 items-center rounded-full transition-colors duration-200',
          controller.includeForgotten ? 'bg-primary' : 'bg-muted',
        )}>
          <span className={cn(
            'inline-block h-[14px] w-[14px] rounded-full bg-white shadow-sm transition-transform duration-200',
            controller.includeForgotten ? 'translate-x-[15px]' : 'translate-x-[2px]',
          )} />
        </span>
      </button>

      <div className="min-h-0 flex-1 overflow-y-auto">
        {controller.items.length === 0 && !controller.loading ? (
          <p className="px-2 py-6 text-center text-sm text-muted-foreground">{t('settingsPage.memory.empty')}</p>
        ) : (
          <ul className="flex flex-col gap-1">
            {controller.items.map((item) => (
              <li key={item.id}>
                <button
                  type="button"
                  onClick={() => {
                    controller.setSelectedId(item.id)
                    onCollapseSubmenu()
                  }}
                  className={cn(
                    'group flex w-full items-start rounded-xl px-3 py-2.5 text-left transition-colors',
                    !isH5 && item.id === controller.selectedId
                      ? 'bg-primary/10 text-primary'
                      : 'bg-background hover:bg-accent',
                  )}
                >
                  <span className="min-w-0 flex-1">
                    <span className="flex min-w-0 items-center gap-2">
                      <span className="truncate text-sm font-medium" title={item.summary}>
                        {item.summary || t('settingsPage.memory.untitled')}
                      </span>
                      {item.is_forgotten ? (
                        <span className="shrink-0 rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground">
                          {t('settingsPage.memory.forgotten')}
                        </span>
                      ) : null}
                    </span>
                    <span className="mt-1 block line-clamp-2 text-xs leading-5 text-muted-foreground">
                      {item.content}
                    </span>
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

// MemorySettingsView renders the right-hand editor for the selected memory.
export function MemorySettingsView(props: { controller: MemoryController }) {
  const { t } = useTranslation()
  const { isH5 } = useContext(SettingsMenuContext)
  const { controller } = props
  const { draft, selected } = controller

  return (
    <SettingsContentLayout
      noContentScroll
      header={
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border/50 py-5">
          <SettingsPanelHeader
            title={draft.id === null ? t('settingsPage.memory.new') : (selected?.summary || t('settingsPage.memory.untitled'))}
          />
          <button
            type="button"
            onClick={() => { void controller.loadMemories() }}
            className="inline-flex h-9 items-center gap-2 rounded-lg px-3 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            <RotateCcw size={15} />
            {t('settingsPage.memory.refresh')}
          </button>
        </div>
      }
      footprint={
        <SettingsActionBar
          primaryLabel={t('settingsPage.memory.save')}
          primaryDisabled={!controller.draftDirty || (!draft.summary.trim() && !draft.content.trim())}
          onPrimaryClick={() => { void controller.saveDraft() }}
          secondaryLabel={controller.draftDirty ? t('settingsPage.providers.cancel') : (selected?.is_forgotten ? t('settingsPage.memory.restore') : undefined)}
          onSecondaryClick={controller.draftDirty ? controller.resetDraft : (selected?.is_forgotten ? () => { void controller.restoreSelected() } : undefined)}
          dangerLabel={selected && !selected.is_forgotten ? t('settingsPage.memory.forget') : undefined}
          onDangerClick={selected && !selected.is_forgotten ? () => { void controller.forgetSelected() } : undefined}
        />
      }
    >
      <div className={cn(
        'flex min-h-0 flex-1 flex-col overflow-y-auto pb-4 pt-6',
        isH5 ? '-mr-5 pr-5' : '-mr-10 pr-10',
      )}>
        <div className="max-w-3xl space-y-5">
          <Field label={t('settingsPage.memory.summary')}>
            <TextInput value={draft.summary} onChange={(value) => controller.setDraft({ ...draft, summary: value })} />
          </Field>

          <Field label={t('settingsPage.memory.content')}>
            <textarea
              value={draft.content}
              onChange={(event) => controller.setDraft({ ...draft, content: event.target.value })}
              rows={8}
              className="w-full resize-none appearance-none rounded-xl border border-border/80 bg-card/80 px-3.5 py-3 text-sm leading-6 text-foreground outline-none transition-colors placeholder:text-muted-foreground focus:border-primary"
            />
          </Field>

          <div className="grid gap-4 sm:grid-cols-2">
            <Field label={t('settingsPage.memory.target')}>
              <SettingsSelect
                ariaLabel={t('settingsPage.memory.target')}
                value={draft.target}
                onChange={(value) => controller.setDraft({ ...draft, target: value })}
                options={[
                  { value: 'user', label: t('settingsPage.memory.targetUser') },
                  { value: 'memory', label: t('settingsPage.memory.targetMemory') },
                ]}
              />
            </Field>
            <Field label={t('settingsPage.memory.type')}>
              <SettingsSelect
                ariaLabel={t('settingsPage.memory.type')}
                value={draft.type}
                onChange={(value) => controller.setDraft({ ...draft, type: value })}
                options={[
                  { value: 'fact', label: t('settingsPage.memory.typeFact') },
                  { value: 'information', label: t('settingsPage.memory.typeInformation') },
                  { value: 'event', label: t('settingsPage.memory.typeEvent') },
                ]}
              />
            </Field>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <Field label={t('settingsPage.memory.importance')}>
              <TextInput
                type="number"
                value={String(draft.importance)}
                onChange={(value) => controller.setDraft({ ...draft, importance: Number(value) })}
              />
            </Field>
            <Field label={t('settingsPage.memory.confidence')}>
              <TextInput
                type="number"
                value={String(draft.confidence)}
                onChange={(value) => controller.setDraft({ ...draft, confidence: Number(value) })}
              />
            </Field>
          </div>

          <button
            type="button"
            aria-pressed={draft.pinned}
            onClick={() => controller.setDraft({ ...draft, pinned: !draft.pinned })}
            className="inline-flex items-center gap-2 rounded-xl px-2 py-1.5 text-sm text-foreground transition-colors hover:bg-accent"
          >
            <span className={cn(
              'flex h-5 w-5 items-center justify-center rounded-md border transition-colors',
              draft.pinned ? 'border-primary bg-primary text-primary-foreground' : 'border-border bg-card',
            )}>
              {draft.pinned ? <Check size={13} /> : null}
            </span>
            {t('settingsPage.memory.pinned')}
          </button>
        </div>
      </div>
    </SettingsContentLayout>
  )
}

function CompactSelect(props: {
  ariaLabel: string
  value: string
  onChange: (value: string) => void
  options: Array<{ value: string; label: string }>
}) {
  const [open, setOpen] = useState(false)
  const selected = props.options.find((option) => option.value === props.value) ?? props.options[0]

  return (
    <div className="relative">
      <button
        type="button"
        aria-label={props.ariaLabel}
        aria-expanded={open}
        onClick={() => setOpen((value) => !value)}
        className="flex h-8 w-full items-center justify-between gap-1.5 rounded-lg border border-border/80 bg-card/80 px-2.5 text-left text-xs text-foreground transition-colors hover:bg-accent focus:border-primary focus:outline-none"
      >
        <span className="min-w-0 flex-1 truncate">{selected?.label}</span>
        <ChevronDown size={13} className={cn('shrink-0 text-muted-foreground transition-transform', open ? 'rotate-180' : '')} />
      </button>
      {open ? (
        <div className="absolute left-0 right-0 top-full z-20 mt-1 overflow-hidden rounded-lg border border-border bg-popover py-1 shadow-lg">
          {props.options.map((option) => (
            <button
              key={option.value}
              type="button"
              onClick={() => {
                props.onChange(option.value)
                setOpen(false)
              }}
              className={cn(
                'flex h-7 w-full items-center justify-between gap-2 px-2.5 text-left text-xs transition-colors hover:bg-accent',
                option.value === props.value ? 'text-primary' : 'text-foreground',
              )}
            >
              <span className="min-w-0 flex-1 truncate">{option.label}</span>
              {option.value === props.value ? <Check size={12} /> : null}
            </button>
          ))}
        </div>
      ) : null}
    </div>
  )
}

function Field(props: { label: string; children: ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-sm font-medium text-foreground">{props.label}</span>
      {props.children}
    </label>
  )
}

function TextInput(props: {
  value: string
  onChange: (value: string) => void
  type?: 'text' | 'number'
}) {
  return (
    <input
      type="text"
      inputMode={props.type === 'number' ? 'numeric' : undefined}
      value={props.value}
      onChange={(event) => props.onChange(event.target.value)}
      className="h-11 w-full appearance-none rounded-xl border border-border/80 bg-card/80 px-3.5 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground focus:border-primary"
    />
  )
}
