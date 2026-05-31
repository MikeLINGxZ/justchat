import { AlertBannerRegion } from '@/components/alert/AlertBannerRegion'
import { AlertToastStack } from '@/components/alert/AlertToastStack'

export function AlertViewport() {
  return (
    <>
      <AlertBannerRegion />
      <AlertToastStack />
    </>
  )
}
