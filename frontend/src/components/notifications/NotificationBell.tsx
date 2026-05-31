import { useMemo, useState } from 'react'
import { Bell } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useNotificationsStore } from '@/store/notificationsStore'
import type { NotificationItem } from '@/types/notifications'
import { FloatingChatPanel } from './FloatingChatPanel'

export function NotificationBell() {
  const { t } = useTranslation()
  const items = useNotificationsStore((state) => state.items)
  const [open, setOpen] = useState(false)
  const [activeNotification, setActiveNotification] = useState<NotificationItem | null>(null)

  const unresolvedCount = items.length
  const sortedItems = useMemo(
    () => [...items].sort((left, right) => right.created_at.localeCompare(left.created_at)),
    [items],
  )

  return (
    <>
      <div className="relative">
        <button
          type="button"
          onClick={() => setOpen((value) => !value)}
          className="relative rounded-lg p-2 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          aria-label={t('notifications.open')}
        >
          <Bell size={16} />
          {unresolvedCount > 0 && (
            <span className="absolute -right-1 -top-1 inline-flex min-w-[1.1rem] items-center justify-center rounded-full bg-destructive px-1 text-[10px] font-medium leading-4 text-destructive-foreground">
              {unresolvedCount}
            </span>
          )}
        </button>

        {open && (
          <>
            <button type="button" className="fixed inset-0 z-10 cursor-default" onClick={() => setOpen(false)} />
            <div className="absolute right-0 z-20 mt-2 w-80 overflow-hidden rounded-xl border border-border bg-popover shadow-xl">
              <div className="border-b border-border px-3 py-2 text-sm font-medium text-foreground">
                {t('notifications.title')}
              </div>
              {sortedItems.length === 0 ? (
                <div className="px-4 py-6 text-center text-sm text-muted-foreground">
                  {t('notifications.empty')}
                </div>
              ) : (
                <div className="max-h-96 overflow-y-auto p-2">
                  {sortedItems.map((item) => (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => {
                        setActiveNotification(item)
                        setOpen(false)
                      }}
                      className="flex w-full flex-col items-start gap-1 rounded-lg px-3 py-2 text-left transition-colors hover:bg-accent"
                    >
                      <span className="text-sm font-medium text-foreground">{item.title}</span>
                      <span className="line-clamp-2 text-xs text-muted-foreground">{item.message}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </>
        )}
      </div>

      {activeNotification && (
        <FloatingChatPanel
          notification={activeNotification}
          onClose={() => setActiveNotification(null)}
        />
      )}
    </>
  )
}
