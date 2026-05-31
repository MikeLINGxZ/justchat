import { Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'

export function SettingsActionBar(props: {
  primaryLabel: string
  primaryDisabled?: boolean
  onPrimaryClick: () => void
  secondaryLabel?: string
  secondaryDisabled?: boolean
  onSecondaryClick?: () => void
  dangerLabel?: string
  onDangerClick?: () => void
  reloadLabel?: string
  onReloadClick?: () => void
  regenerateLabel?: string
  regenerateLoading?: boolean
  regenerateDisabled?: boolean
  onRegenerateClick?: () => void
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <div className="flex items-center gap-3">
        {props.regenerateLabel && props.onRegenerateClick && (
          <button
            type="button"
            disabled={props.regenerateDisabled}
            onClick={props.onRegenerateClick}
            className="flex items-center gap-1.5 rounded-xl border border-border px-4 py-2 text-sm text-foreground transition-colors hover:bg-accent disabled:cursor-not-allowed disabled:opacity-60"
          >
            {props.regenerateLoading && <Loader2 size={14} className="animate-spin" />}
            {props.regenerateLabel}
          </button>
        )}
        {props.dangerLabel && props.onDangerClick && (
          <button
            type="button"
            onClick={props.onDangerClick}
            className="rounded-xl border border-red-300 px-4 py-2 text-sm text-red-500 transition-colors hover:bg-red-500/10"
          >
            {props.dangerLabel}
          </button>
        )}
        {props.reloadLabel && props.onReloadClick && (
          <button
            type="button"
            onClick={props.onReloadClick}
            className="rounded-xl border border-border px-4 py-2 text-sm text-foreground transition-colors hover:bg-accent"
          >
            {props.reloadLabel}
          </button>
        )}
      </div>

      <div className="flex items-center gap-3">
        {props.secondaryLabel && props.onSecondaryClick && (
          <button
            type="button"
            disabled={props.secondaryDisabled}
            onClick={props.onSecondaryClick}
            className={cn(
              'rounded-xl border border-border px-4 py-2 text-sm transition-colors',
              props.secondaryDisabled
                ? 'cursor-not-allowed text-muted-foreground'
                : 'text-foreground hover:bg-accent'
            )}
          >
            {props.secondaryLabel}
          </button>
        )}
        <button
          type="button"
          disabled={props.primaryDisabled}
          onClick={props.onPrimaryClick}
          className={cn(
            'rounded-xl px-4 py-2 text-sm font-medium transition-colors',
            props.primaryDisabled
              ? 'cursor-not-allowed bg-muted text-muted-foreground'
              : 'bg-primary text-primary-foreground hover:opacity-90'
          )}
        >
          {props.primaryLabel}
        </button>
      </div>
    </div>
  )
}
