import { ChevronDown, Check } from 'lucide-react'
import { useEffect, useId, useRef, useState } from 'react'
import { cn } from '@/lib/utils'

type SettingsSelectOption<T extends string> = {
  value: T
  label: string
  icon?: string
}

export function SettingsSelect<T extends string>(props: {
  ariaLabel: string
  value: T
  options: SettingsSelectOption<T>[]
  onChange: (value: T) => void
}) {
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement | null>(null)
  const listboxId = useId()
  const selectedOption = props.options.find((option) => option.value === props.value) ?? props.options[0]

  useEffect(() => {
    if (!open) {
      return undefined
    }

    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setOpen(false)
      }
    }

    document.addEventListener('mousedown', handlePointerDown)
    document.addEventListener('keydown', handleEscape)
    return () => {
      document.removeEventListener('mousedown', handlePointerDown)
      document.removeEventListener('keydown', handleEscape)
    }
  }, [open])

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        role="combobox"
        aria-label={props.ariaLabel}
        aria-expanded={open}
        aria-controls={listboxId}
        onClick={() => setOpen((current) => !current)}
        className={cn(
          'flex h-11 w-full items-center justify-between rounded-xl border border-border/80 bg-card/80 px-3.5 text-left',
          'text-sm text-foreground transition-colors hover:bg-accent/60 focus:outline-none focus:ring-2 focus:ring-primary/20'
        )}
      >
        <span className="inline-flex items-center gap-2 truncate">
          {selectedOption?.icon && <span className="shrink-0" aria-hidden="true">{selectedOption.icon}</span>}
          <span className="truncate">{selectedOption?.label}</span>
        </span>
        <ChevronDown size={16} className={cn('shrink-0 text-muted-foreground transition-transform', open && 'rotate-180')} />
      </button>

      {open && (
        <div
          id={listboxId}
          role="listbox"
          aria-label={props.ariaLabel}
          className="absolute left-0 top-[calc(100%+0.5rem)] z-30 w-full rounded-2xl border border-border/80 bg-popover p-1.5 shadow-xl max-h-72 overflow-y-auto overscroll-contain"
        >
          {props.options.map((option) => {
            const selected = option.value === props.value
            return (
              <button
                key={option.value}
                type="button"
                role="option"
                aria-selected={selected}
                onClick={() => {
                  props.onChange(option.value)
                  setOpen(false)
                }}
                className={cn(
                  'flex w-full items-center gap-2 rounded-xl px-3 py-2.5 text-left text-sm transition-colors',
                  selected
                    ? 'bg-primary/10 font-medium text-primary'
                    : 'text-foreground hover:bg-accent'
                )}
              >
                <span className="inline-flex items-center gap-2 min-w-0 flex-1 truncate">
                  {option.icon && <span className="shrink-0" aria-hidden="true">{option.icon}</span>}
                  <span className="truncate">{option.label}</span>
                </span>
                {selected && <Check size={14} className="shrink-0" />}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
