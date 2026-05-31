import { render } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AlertEventProvider } from '@/components/providers/AlertEventProvider'
import { useAlertStore } from '@/alert/store'
import { APP_ALERT_EVENT } from '@/alert/event'
import App from '@/App'

const handlers = new Map<string, (event: { data: unknown }) => void>()

vi.mock('@/components/layout/MainLayout', () => ({
  MainLayout: () => <div>Main Layout</div>,
}))

vi.mock('@/components/settings/SettingsApp', () => ({
  SettingsApp: () => <div>Settings App</div>,
}))

vi.mock('@/components/settings/providers/AddProviderApp', () => ({
  AddProviderApp: () => <div>Add Provider App</div>,
}))

vi.mock('@wailsio/runtime', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wailsio/runtime')>()
  return {
    ...actual,
    Events: {
      ...actual.Events,
      On: vi.fn((name: string, handler: (event: { data: unknown }) => void) => {
        handlers.set(name, handler)
        return () => handlers.delete(name)
      }),
    },
  }
})

describe('AlertEventProvider', () => {
  beforeEach(() => {
    handlers.clear()
    useAlertStore.getState().reset()
  })

  it('maps backend Wails payload into a toast alert', () => {
    render(
      <AlertEventProvider>
        <div>child</div>
      </AlertEventProvider>,
    )

    handlers.get(APP_ALERT_EVENT)?.({
      data: {
        kind: 'success',
        placement: 'toast',
        title: 'Saved',
        message: 'Saved from backend',
      },
    })

    expect(useAlertStore.getState().toasts).toHaveLength(1)
    expect(useAlertStore.getState().toasts[0]).toMatchObject({
      kind: 'success',
      title: 'Saved',
      message: 'Saved from backend',
      source: 'backend',
    })
  })

  it('ignores malformed payloads without crashing', () => {
    render(
      <AlertEventProvider>
        <div>child</div>
      </AlertEventProvider>,
    )

    handlers.get(APP_ALERT_EVENT)?.({
      data: {
        placement: 'toast',
      },
    })

    expect(useAlertStore.getState().toasts).toHaveLength(0)
    expect(useAlertStore.getState().banners).toHaveLength(0)
  })

  it('wires backend events into the mounted app viewport', () => {
    render(<App />)

    handlers.get(APP_ALERT_EVENT)?.({
      data: {
        kind: 'success',
        placement: 'toast',
        title: 'Saved',
        message: 'Mounted in app',
      },
    })

    expect(useAlertStore.getState().toasts).toHaveLength(1)
  })
})
