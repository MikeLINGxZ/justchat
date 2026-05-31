import { useTranslation } from 'react-i18next'
import { FileText, Image as ImageIcon, Paperclip, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { isImageMime, isPdfMime } from '@/lib/attachments'
import type { Attachment } from '@/types'

interface Props {
  items: Attachment[]
  variant: 'input' | 'message'
  onRemove?: (index: number) => void
  onOpen?: (item: Attachment) => void
}

function chipIcon(att: Attachment) {
  if (isImageMime(att.mime)) return <ImageIcon size={12} />
  if (isPdfMime(att.mime)) return <FileText size={12} />
  return <Paperclip size={12} />
}

export function AttachmentChips({ items, variant, onRemove, onOpen }: Props) {
  const { t } = useTranslation()

  if (items.length === 0) return null

  return (
    <div className="flex flex-wrap gap-1.5 px-3 py-2">
      {items.map((att, index) => (
        variant === 'message' ? (
          <button
            key={`${att.path}-${index}`}
            type="button"
            onClick={() => onOpen?.(att)}
            className={cn(
              'inline-flex items-center gap-1 max-w-[18rem] rounded-md border border-border bg-muted/60 px-2 py-1 text-xs text-foreground',
              'bg-background/60 hover:bg-background/80 transition-colors'
            )}
            title={att.name}
          >
            {chipIcon(att)}
            <span className="truncate">{att.name}</span>
          </button>
        ) : (
          <span
            key={`${att.path}-${index}`}
            className="inline-flex items-center gap-1 max-w-[18rem] rounded-md border border-border bg-muted/60 px-2 py-1 text-xs text-foreground"
          >
            {chipIcon(att)}
            <span className="truncate" title={att.name}>{att.name}</span>
            {onRemove && (
            <button
              type="button"
              aria-label={t('chat.attachmentRemove')}
              onClick={() => onRemove(index)}
              className="ml-0.5 rounded text-muted-foreground hover:text-foreground hover:bg-accent p-0.5"
            >
              <X size={10} />
            </button>
            )}
          </span>
        )
      ))}
    </div>
  )
}
