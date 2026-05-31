import { Loader2 } from 'lucide-react'
import type { ReactNode } from 'react'

type ConfirmDialogProps = {
  open: boolean
  title: string
  description?: ReactNode
  confirmLabel: string
  cancelLabel: string
  confirmTone?: 'default' | 'danger'
  busy?: boolean
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmDialog(props: ConfirmDialogProps) {
  if (!props.open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        aria-label={props.cancelLabel}
        className="absolute inset-0 bg-black/35"
        onClick={props.onCancel}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
        className="relative z-10 w-full max-w-sm rounded-3xl border border-border bg-background p-5 shadow-2xl"
      >
        <div className="space-y-2">
          <h3 id="confirm-dialog-title" className="text-base font-semibold text-foreground">
            {props.title}
          </h3>
          {props.description && (
            <div className="text-sm leading-6 text-muted-foreground">
              {props.description}
            </div>
          )}
        </div>

        <div className="mt-5 flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={props.onCancel}
            className="rounded-xl border border-border px-4 py-2 text-sm text-foreground transition-colors hover:bg-accent"
          >
            {props.cancelLabel}
          </button>
          <button
            type="button"
            disabled={props.busy}
            onClick={props.onConfirm}
            className={
              props.confirmTone === 'danger'
                ? 'flex items-center gap-1.5 rounded-xl bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60'
                : 'flex items-center gap-1.5 rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60'
            }
          >
            {props.busy && <Loader2 size={14} className="animate-spin" />}
            {props.confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}
