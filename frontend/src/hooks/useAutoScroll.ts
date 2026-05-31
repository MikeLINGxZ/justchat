import { useEffect, useLayoutEffect, useRef, useState, useCallback } from 'react'

function scrollElementToBottom(el: HTMLDivElement, smooth = false) {
  const top = el.scrollHeight
  if (typeof el.scrollTo === 'function') {
    el.scrollTo({ top, behavior: smooth ? 'smooth' : 'auto' })
    return
  }
  el.scrollTop = top
}

export function useAutoScroll(deps: unknown[], resetKey?: unknown) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [isAtBottom, setIsAtBottom] = useState(true)
  const autoScrollRef = useRef(true)

  const scrollToBottom = useCallback((smooth = false) => {
    const el = containerRef.current
    if (!el) return
    scrollElementToBottom(el, smooth)
    autoScrollRef.current = true
    setIsAtBottom(true)
  }, [])

  useEffect(() => {
    const el = containerRef.current
    if (!el) return
    const handler = () => {
      const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 50
      setIsAtBottom(atBottom)
      autoScrollRef.current = atBottom
    }
    el.addEventListener('scroll', handler, { passive: true })
    return () => el.removeEventListener('scroll', handler)
  }, [])

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useLayoutEffect(() => {
    if (resetKey === undefined) return
    autoScrollRef.current = true
    const el = containerRef.current
    if (!el) return
    scrollElementToBottom(el)
    setIsAtBottom(true)
  }, [resetKey])

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useLayoutEffect(() => {
    if (autoScrollRef.current) {
      const el = containerRef.current
      if (!el) return
      scrollElementToBottom(el)
    }
  }, deps)

  return { containerRef, isAtBottom, scrollToBottom, autoScrollRef }
}
