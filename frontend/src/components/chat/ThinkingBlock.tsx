import { useState } from 'react'
import { Brain, ChevronDown } from 'lucide-react'
import { useTranslation } from 'react-i18next'

interface Props {
  content: string
  defaultExpanded?: boolean
  active?: boolean
}

export function ThinkingBlock({ content, defaultExpanded = false, active = false }: Props) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(defaultExpanded)

  return (
    <div className="mb-0">
      <button
        onClick={() => setExpanded(v => !v)}
        className="flex w-full min-w-0 items-center gap-1.5 rounded-md py-0.5 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
      >
        <Brain size={14} className="shrink-0" />
        <span>{t(active ? 'chat.thinking' : 'chat.thought')}</span>
        <ChevronDown
          size={14}
          className={expanded ? 'shrink-0 rotate-180 transition-transform' : 'shrink-0 rotate-0 transition-transform'}
        />
        <span className="min-w-0 flex-1" />
      </button>
      {expanded && (
        <div className="ml-1.5 mt-1 select-text whitespace-pre-wrap border-l border-border pl-3 text-sm leading-relaxed text-foreground">
          {content}
        </div>
      )}
    </div>
  )
}
