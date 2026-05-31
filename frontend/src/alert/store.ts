import { create } from 'zustand'
import { createStore } from 'zustand/vanilla'
import type { StoreApi } from 'zustand'
import type {
  AlertAriaPriority,
  AlertCloseReason,
  AlertInput,
  AlertItem,
  ClosedAlertRecord,
} from '@/alert/types'

export const ALERT_DEFAULT_DURATION_MS = 5000
export const ALERT_MAX_DURATION_MS = 15000
export const ALERT_TOAST_LIMIT = 3
export const ALERT_BANNER_LIMIT = 1

type AlertState = {
  toasts: AlertItem[]
  banners: AlertItem[]
  lastClosed: ClosedAlertRecord | null
  pushAlert: (input: AlertInput) => string
  closeAlert: (id: string, reason: AlertCloseReason) => void
  reset: () => void
}

type TimerRegistry = Map<string, ReturnType<typeof setTimeout>>

let nextAlertId = 0

function buildDedupeKey(input: AlertInput): string {
  if (input.dedupeKey && input.dedupeKey.trim().length > 0) {
    return input.dedupeKey
  }

  return [
    input.kind,
    input.placement,
    input.title,
    input.message,
    input.detail ?? '',
  ].join('::')
}

function resolveAriaPriority(kind: AlertInput['kind']): AlertAriaPriority {
  return kind === 'error' ? 'assertive' : 'polite'
}

function resolveAutoClose(input: AlertInput): boolean {
  if (typeof input.autoClose === 'boolean') {
    return input.autoClose
  }

  const actionCount = input.actions?.length ?? 0
  if (actionCount > 0) {
    return false
  }

  return input.kind !== 'error'
}

function resolveDurationMs(input: AlertInput, autoClose: boolean): number | null {
  if (!autoClose) {
    return null
  }

  const candidate = input.durationMs ?? ALERT_DEFAULT_DURATION_MS
  if (candidate <= 0) {
    return ALERT_DEFAULT_DURATION_MS
  }

  return Math.min(candidate, ALERT_MAX_DURATION_MS)
}

function normalizeAlert(input: AlertInput): AlertItem {
  const now = Date.now()
  const autoClose = resolveAutoClose(input)
  return {
    id: `alert-${now}-${nextAlertId += 1}`,
    kind: input.kind,
    placement: input.placement,
    title: input.title,
    message: input.message,
    detail: input.detail ?? null,
    code: input.code ?? null,
    actions: input.actions ?? [],
    dismissible: input.dismissible ?? true,
    autoClose,
    durationMs: resolveDurationMs(input, autoClose),
    count: 1,
    dedupeKey: buildDedupeKey(input),
    dedupeBehavior: input.dedupeBehavior ?? 'count',
    showCount: input.showCount ?? true,
    createdAt: now,
    updatedAt: now,
    ariaPriority: resolveAriaPriority(input.kind),
    source: input.source ?? 'frontend',
  }
}

function clearAlertTimer(timers: TimerRegistry, id: string): void {
  const timer = timers.get(id)
  if (!timer) {
    return
  }

  clearTimeout(timer)
  timers.delete(id)
}

function priorityRank(kind: AlertItem['kind']): number {
  if (kind === 'error') {
    return 4
  }
  if (kind === 'warning') {
    return 3
  }
  if (kind === 'success') {
    return 2
  }
  return 1
}

function createAlertState(timers: TimerRegistry) {
  return (set: StoreApi<AlertState>['setState'], get: StoreApi<AlertState>['getState']): AlertState => {
    const scheduleTimer = (item: AlertItem): void => {
      clearAlertTimer(timers, item.id)
      if (!item.autoClose || item.durationMs === null) {
        return
      }

      const timer = setTimeout(() => {
        get().closeAlert(item.id, 'timeout')
      }, item.durationMs)
      timers.set(item.id, timer)
    }

    const applyToastLimit = (items: AlertItem[]): AlertItem[] => {
      if (items.length <= ALERT_TOAST_LIMIT) {
        return items
      }

      const removed = items.slice(0, items.length - ALERT_TOAST_LIMIT)
      for (const item of removed) {
        clearAlertTimer(timers, item.id)
      }

      return items.slice(items.length - ALERT_TOAST_LIMIT)
    }

    return {
      toasts: [],
      banners: [],
      lastClosed: null,
      pushAlert: (input) => {
        const existingList = input.placement === 'toast' ? get().toasts : get().banners
        const dedupeKey = buildDedupeKey(input)
        const existing = existingList.find((item) => item.dedupeKey === dedupeKey)

        if (existing) {
          const now = Date.now()
          const autoClose = resolveAutoClose(input)
          const shouldReplace = (input.dedupeBehavior ?? existing.dedupeBehavior) === 'replace'
          const updated: AlertItem = shouldReplace
            ? {
                ...existing,
                kind: input.kind,
                title: input.title,
                message: input.message,
                detail: input.detail ?? null,
                code: input.code ?? null,
                actions: input.actions ?? [],
                dismissible: input.dismissible ?? true,
                autoClose,
                durationMs: resolveDurationMs(input, autoClose),
                count: 1,
                updatedAt: now,
                ariaPriority: resolveAriaPriority(input.kind),
                source: input.source ?? existing.source,
                dedupeBehavior: input.dedupeBehavior ?? existing.dedupeBehavior,
                showCount: input.showCount ?? existing.showCount,
              }
            : {
                ...existing,
                count: existing.count + 1,
                updatedAt: now,
              }
          const nextList = existingList.map((item) => item.id === existing.id ? updated : item)
          if (input.placement === 'toast') {
            set({ toasts: nextList })
          } else {
            set({ banners: nextList })
          }
          scheduleTimer(updated)
          return existing.id
        }

        const item = normalizeAlert(input)
        if (item.placement === 'toast') {
          const nextToasts = applyToastLimit([...get().toasts, item])
          set({ toasts: nextToasts })
        } else {
          const currentBanner = get().banners[0]
          if (!currentBanner || priorityRank(item.kind) >= priorityRank(currentBanner.kind)) {
            if (currentBanner) {
              clearAlertTimer(timers, currentBanner.id)
            }
            set({
              banners: [item],
              lastClosed: currentBanner ? { id: currentBanner.id, reason: 'replaced' } : get().lastClosed,
            })
          }
        }

        scheduleTimer(item)
        return item.id
      },
      closeAlert: (id, reason) => {
        clearAlertTimer(timers, id)
        set((state) => ({
          toasts: state.toasts.filter((item) => item.id !== id),
          banners: state.banners.filter((item) => item.id !== id),
          lastClosed: { id, reason },
        }))
      },
      reset: () => {
        for (const id of Array.from(timers.keys())) {
          clearAlertTimer(timers, id)
        }
        set({
          toasts: [],
          banners: [],
          lastClosed: null,
        })
      },
    }
  }
}

export function createAlertStore(): StoreApi<AlertState> {
  const timers: TimerRegistry = new Map<string, ReturnType<typeof setTimeout>>()
  return createStore<AlertState>(createAlertState(timers))
}

const alertStore = createAlertStore()

export const useAlertStore = create<AlertState>()((set, get) =>
  createAlertState(new Map<string, ReturnType<typeof setTimeout>>())(set, get)
)

export function resetAlertStore(): void {
  useAlertStore.getState().reset()
}

export { alertStore }
export type { AlertState }
