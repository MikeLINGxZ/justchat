import { useEffect } from 'react'
import { useAppStore } from '@/store/appStore'
import type { FontSize } from '@/types'

const FONT_SIZE_MAP: Record<FontSize, string> = {
  xs: '13px',
  sm: '14px',
  md: '16px',
  lg: '18px',
  xl: '20px',
}

export function FontSizeProvider({ children }: { children: React.ReactNode }) {
  const fontSize = useAppStore(s => s.fontSize)

  useEffect(() => {
    document.documentElement.style.fontSize = FONT_SIZE_MAP[fontSize]
  }, [fontSize])

  return <>{children}</>
}
