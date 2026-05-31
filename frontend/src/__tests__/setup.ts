import '@testing-library/jest-dom'
import '@/i18n'

// ClipboardEvent is not implemented in jsdom – provide a minimal polyfill so
// paste-related tests can construct and dispatch clipboard events.

if (typeof ClipboardEvent === 'undefined') {
  class ClipboardEventPolyfill extends Event {
    clipboardData: DataTransfer | null
    constructor(type: string, init?: ClipboardEventInit) {
      super(type, init)
      this.clipboardData = (init as { clipboardData?: DataTransfer | null } | undefined)?.clipboardData ?? null
    }
  }
  globalThis.ClipboardEvent = ClipboardEventPolyfill
}

// document.elementFromPoint is not implemented in jsdom – stub it so that
// @wailsio/runtime drag listeners do not crash during tests.
if (typeof document.elementFromPoint !== 'function') {
  document.elementFromPoint = () => null
}

if (typeof window.matchMedia !== 'function') {
  window.matchMedia = ((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: () => undefined,
    removeEventListener: () => undefined,
    addListener: () => undefined,
    removeListener: () => undefined,
    dispatchEvent: () => false,
  })) as typeof window.matchMedia
}

// DragEvent is not implemented in jsdom – provide a minimal polyfill so
// drag-and-drop related tests can construct and dispatch drag events.

if (typeof DragEvent === 'undefined') {
  class DragEventPolyfill extends Event {
    dataTransfer: DataTransfer | null
    relatedTarget: EventTarget | null
    constructor(type: string, init?: DragEventInit) {
      super(type, init)
      this.dataTransfer = (init as { dataTransfer?: DataTransfer | null } | undefined)?.dataTransfer ?? null
      this.relatedTarget = (init as { relatedTarget?: EventTarget | null } | undefined)?.relatedTarget ?? null
    }
  }
  globalThis.DragEvent = DragEventPolyfill as unknown as typeof DragEvent
}
