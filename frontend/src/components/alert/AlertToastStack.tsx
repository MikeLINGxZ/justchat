import { useAlertStore } from '@/alert/store'
import { AlertCard } from '@/components/alert/AlertCard'

export function AlertToastStack() {
  const toasts = useAlertStore((state) => state.toasts)

  return (
    <section
      aria-label="Toast alerts"
      className="pointer-events-none fixed bottom-4 right-4 z-50 flex flex-col items-end gap-3"
    >
      {toasts.map((alert) => (
        <div key={alert.id} className="pointer-events-auto w-auto max-w-[min(26rem,calc(100vw-2rem))]">
          <AlertCard alert={alert} />
        </div>
      ))}
    </section>
  )
}
