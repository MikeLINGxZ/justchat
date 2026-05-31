import { cn } from '@/lib/utils'
import { Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { SkillItem, SkillSource } from '@/types/skills'

const sourceKeySuffix: Record<SkillSource, string> = {
  builtin: 'Builtin',
  user: 'User',
  ai: 'AI',
}

type Props = {
  item: SkillItem
  selected: boolean
  onSelect: (name: string) => void
  onToggle: (item: SkillItem) => void
  onDelete: (item: SkillItem) => void
}

export function SkillsListItem({ item, selected, onSelect, onToggle, onDelete }: Props) {
  const { t } = useTranslation()
  const badge = t(`settingsPage.skills.source${sourceKeySuffix[item.source]}`)
  const canDelete = item.source !== 'builtin'

  return (
    <div
      className={cn(
        'group relative flex w-full items-start gap-3 rounded-2xl px-3 py-3 transition-colors',
        selected ? 'bg-primary/10 text-primary' : 'bg-background hover:bg-accent'
      )}
    >
      <div
        role="button"
        tabIndex={0}
        onClick={() => onSelect(item.name)}
        onKeyDown={(event) => {
          if (event.key === 'Enter' || event.key === ' ') {
            onSelect(item.name)
          }
        }}
        className="flex min-w-0 flex-1 cursor-pointer items-start"
      >
        <div className="min-w-0 flex-1">
          <div className="flex min-w-0 items-center gap-2">
            <span className="truncate text-sm font-medium text-foreground" title={item.name}>
              {item.name}
            </span>
            <span className="shrink-0 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
              {badge}
            </span>
          </div>
          <p className="mt-1 line-clamp-2 text-xs leading-5 text-muted-foreground">
            {item.description}
          </p>
        </div>
      </div>

      <div className="flex shrink-0 flex-col items-center gap-2">
        <button
          type="button"
          role="switch"
          aria-checked={!item.disabled}
          aria-label={t('settingsPage.skills.enabled')}
          onClick={(event) => {
            event.stopPropagation()
            onToggle(item)
          }}
          className={cn(
            'relative mt-0.5 inline-flex h-[18px] w-[31px] items-center rounded-full transition-colors duration-200',
            item.disabled ? 'bg-muted' : 'bg-primary'
          )}
        >
          <span
            className={cn(
              'inline-block h-[14px] w-[14px] rounded-full bg-white shadow-sm transition-transform duration-200',
              item.disabled ? 'translate-x-[2px]' : 'translate-x-[15px]'
            )}
          />
        </button>
        {canDelete ? (
          <button
            type="button"
            aria-label={t('settingsPage.skills.delete')}
            title={t('settingsPage.skills.delete')}
            onClick={(event) => {
              event.stopPropagation()
              onDelete(item)
            }}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
          >
            <Trash2 size={14} />
          </button>
        ) : null}
      </div>
    </div>
  )
}
