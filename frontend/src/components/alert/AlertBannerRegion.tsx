import { useAlertStore } from '@/alert/store'
import { AlertCard } from '@/components/alert/AlertCard'

export function AlertBannerRegion() {
  const banners = useAlertStore((state) => state.banners)

  return (
    <section
      aria-label="Banner alerts"
      className="pointer-events-none fixed inset-x-0 top-16 z-50 flex justify-center px-6"
    >
      {banners.map((alert) => (
        <div key={alert.id} className="pointer-events-auto w-full max-w-2xl">
          <AlertCard alert={alert} />
        </div>
      ))}
    </section>
  )
}
