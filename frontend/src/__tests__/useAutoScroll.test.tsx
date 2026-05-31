import { act, fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useAutoScroll } from '@/hooks/useAutoScroll'

type AutoScrollState = ReturnType<typeof useAutoScroll>

function Harness({
  deps,
  resetKey,
  onHook,
}: {
  deps: unknown[]
  resetKey?: unknown
  onHook: (state: AutoScrollState) => void
}) {
  const state = useAutoScroll(deps, resetKey)
  onHook(state)
  return <div ref={state.containerRef} data-testid="scroll-container" />
}

function setup(initialDeps: unknown[] = ['chunk-1'], resetKey?: unknown) {
  let hookState: AutoScrollState | undefined
  const prototypeScrollTo = vi.fn()
  Object.defineProperty(HTMLElement.prototype, 'scrollTo', {
    configurable: true,
    value: prototypeScrollTo,
  })

  const rendered = render(
    <Harness
      deps={initialDeps}
      resetKey={resetKey}
      onHook={(state) => {
        hookState = state
      }}
    />
  )
  const element = screen.getByTestId('scroll-container')
  const scrollTo = vi.fn(({ top }: ScrollToOptions) => {
    element.scrollTop = Number(top)
  })

  Object.defineProperties(element, {
    clientHeight: { configurable: true, value: 100 },
    scrollHeight: { configurable: true, value: 500 },
    scrollTo: { configurable: true, value: scrollTo },
  })
  scrollTo.mockClear()

  return {
    ...rendered,
    element,
    scrollTo,
    getHook: () => {
      if (!hookState) throw new Error('hook state not captured')
      return hookState
    },
  }
}

describe('useAutoScroll', () => {
  it('keeps following new output while auto-scroll is enabled', () => {
    const { rerender, scrollTo } = setup()

    rerender(<Harness deps={['chunk-2']} onHook={() => undefined} />)

    expect(scrollTo).toHaveBeenCalledWith({ top: 500, behavior: 'auto' })
  })

  it('pauses after the user scrolls away and resumes when requested', () => {
    let hookState: AutoScrollState | undefined
    const { element, rerender, scrollTo } = setup()

    act(() => {
      element.scrollTop = 100
      fireEvent.scroll(element)
    })

    scrollTo.mockClear()
    rerender(<Harness deps={['chunk-2']} onHook={(state) => { hookState = state }} />)
    expect(scrollTo).not.toHaveBeenCalled()

    act(() => {
      hookState?.scrollToBottom(true)
    })
    expect(scrollTo).toHaveBeenCalledWith({ top: 500, behavior: 'smooth' })

    scrollTo.mockClear()
    rerender(<Harness deps={['chunk-3']} onHook={() => undefined} />)
    expect(scrollTo).toHaveBeenCalledWith({ top: 500, behavior: 'auto' })
  })

  it('can be reset when a new send starts', () => {
    const { element, rerender, scrollTo } = setup(['chunk-1'], 1)

    act(() => {
      element.scrollTop = 100
      fireEvent.scroll(element)
    })

    scrollTo.mockClear()
    rerender(<Harness deps={['chunk-2']} resetKey={2} onHook={() => undefined} />)

    expect(scrollTo).toHaveBeenCalledWith({ top: 500, behavior: 'auto' })
  })
})
