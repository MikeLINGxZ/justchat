import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  ALERT_BANNER_LIMIT,
  ALERT_TOAST_LIMIT,
  createAlertStore,
} from '@/alert/store'

describe('alert store', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  it('aggregates duplicate toast alerts and resets timeout', () => {
    const store = createAlertStore()
    const firstId = store.getState().pushAlert({
      kind: 'success',
      placement: 'toast',
      title: 'Saved',
      message: 'Profile saved',
    })

    vi.advanceTimersByTime(4000)

    store.getState().pushAlert({
      kind: 'success',
      placement: 'toast',
      title: 'Saved',
      message: 'Profile saved',
    })

    expect(store.getState().toasts).toHaveLength(1)
    expect(store.getState().toasts[0].id).toBe(firstId)
    expect(store.getState().toasts[0].count).toBe(2)

    vi.advanceTimersByTime(4999)
    expect(store.getState().toasts).toHaveLength(1)

    vi.advanceTimersByTime(1)
    expect(store.getState().toasts).toHaveLength(0)
  })

  it('replaces duplicate toast content without increasing count when dedupe behavior is replace', () => {
    const store = createAlertStore()
    const firstId = store.getState().pushAlert({
      kind: 'info',
      placement: 'toast',
      title: 'Plugin runtime',
      message: 'Download later',
      dedupeKey: 'runtime-toast',
      dedupeBehavior: 'replace',
      showCount: false,
      autoClose: false,
    })

    const secondId = store.getState().pushAlert({
      kind: 'info',
      placement: 'toast',
      title: 'Plugin runtime',
      message: 'Downloading 42%',
      detail: 'Download · 4.2 MB / 10.0 MB',
      dedupeKey: 'runtime-toast',
      dedupeBehavior: 'replace',
      showCount: false,
      autoClose: false,
    })

    expect(secondId).toBe(firstId)
    expect(store.getState().toasts).toHaveLength(1)
    expect(store.getState().toasts[0].message).toBe('Downloading 42%')
    expect(store.getState().toasts[0].detail).toBe('Download · 4.2 MB / 10.0 MB')
    expect(store.getState().toasts[0].count).toBe(1)
    expect(store.getState().toasts[0].showCount).toBe(false)
  })

  it('keeps error alerts open by default', () => {
    const store = createAlertStore()

    store.getState().pushAlert({
      kind: 'error',
      placement: 'toast',
      title: 'Failed',
      message: 'Request failed',
    })

    vi.advanceTimersByTime(20000)
    expect(store.getState().toasts).toHaveLength(1)
  })

  it('clamps duration to fifteen seconds', () => {
    const store = createAlertStore()

    store.getState().pushAlert({
      kind: 'info',
      placement: 'toast',
      title: 'Syncing',
      message: 'Sync in progress',
      durationMs: 30000,
    })

    vi.advanceTimersByTime(14999)
    expect(store.getState().toasts).toHaveLength(1)

    vi.advanceTimersByTime(1)
    expect(store.getState().toasts).toHaveLength(0)
  })

  it('limits visible toasts to the configured maximum', () => {
    const store = createAlertStore()

    for (const title of ['One', 'Two', 'Three', 'Four']) {
      store.getState().pushAlert({
        kind: 'info',
        placement: 'toast',
        title,
        message: title,
      })
    }

    expect(ALERT_TOAST_LIMIT).toBe(3)
    expect(store.getState().toasts).toHaveLength(3)
    expect(store.getState().toasts.map((item) => item.title)).toEqual(['Two', 'Three', 'Four'])
  })

  it('replaces lower-priority banner with higher-priority banner', () => {
    const store = createAlertStore()

    store.getState().pushAlert({
      kind: 'info',
      placement: 'banner',
      title: 'Heads up',
      message: 'Information',
    })
    store.getState().pushAlert({
      kind: 'error',
      placement: 'banner',
      title: 'Critical',
      message: 'Danger',
    })

    expect(ALERT_BANNER_LIMIT).toBe(1)
    expect(store.getState().banners).toHaveLength(1)
    expect(store.getState().banners[0].title).toBe('Critical')
  })

  it('records close reason when dismissed programmatically', () => {
    const store = createAlertStore()
    const id = store.getState().pushAlert({
      kind: 'warning',
      placement: 'toast',
      title: 'Warning',
      message: 'Check input',
    })

    store.getState().closeAlert(id, 'programmatic')

    expect(store.getState().toasts).toHaveLength(0)
    expect(store.getState().lastClosed?.reason).toBe('programmatic')
  })
})
