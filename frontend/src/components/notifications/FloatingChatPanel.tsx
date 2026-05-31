import { useEffect } from 'react'
import { X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Notification as NotificationBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification'
import { useChatStore } from '@/store/chatStore'
import { ChatInput } from '@/components/chat/ChatInput'
import { ChatMessages } from '@/components/chat/ChatMessages'
import type { NotificationItem } from '@/types/notifications'

type FloatingChatPanelProps = {
  notification: NotificationItem
  onClose: () => void
}

export function FloatingChatPanel({ notification, onClose }: FloatingChatPanelProps) {
  const { t } = useTranslation()
  const loadMessages = useChatStore((state) => state.loadMessages)

  useEffect(() => {
    void loadMessages(notification.session_id)
  }, [loadMessages, notification.session_id])

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [onClose])

  const handleResolve = async () => {
    await NotificationBinding.ResolveNotification({ id: notification.id })
    onClose()
  }

  const handleReject = async () => {
    await NotificationBinding.RejectNotification({ id: notification.id })
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/35 px-4 py-6">
      <div className="flex h-[720px] max-h-full w-[920px] max-w-full flex-col overflow-hidden rounded-2xl border border-border bg-background shadow-2xl">
        <div className="flex items-start justify-between gap-4 border-b border-border px-4 py-3">
          <div className="min-w-0">
            <h2 className="truncate text-base font-semibold text-foreground">{notification.title}</h2>
            <p className="mt-1 text-sm text-muted-foreground">{notification.message}</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg p-2 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            <X size={16} />
          </button>
        </div>

        <ChatMessages sessionId={notification.session_id} />

        <div className="border-t border-border bg-background/95">
          <ChatInput sessionId={notification.session_id} />
          {notification.kind === 'needs_attention' && (
            <div className="flex items-center justify-end gap-2 px-4 pb-4">
              <button
                type="button"
                onClick={() => { void handleReject() }}
                className="rounded-lg border border-border px-3 py-2 text-sm text-foreground transition-colors hover:bg-accent"
              >
                {t('common.reject')}
              </button>
              <button
                type="button"
                onClick={() => { void handleResolve() }}
                className="rounded-lg bg-primary px-3 py-2 text-sm text-primary-foreground transition-colors hover:bg-primary/90"
              >
                {t('common.confirm')}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
