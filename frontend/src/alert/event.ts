import type { AlertAction, AlertInput, AlertKind, AlertPlacement } from '@/alert/types'

export const APP_ALERT_EVENT = 'app:alert'

export type AlertEventPayload = {
  kind?: AlertKind
  placement?: AlertPlacement
  title?: string
  message?: string
  detail?: string | null
  code?: string | null
  dismissible?: boolean
  auto_close?: boolean
  duration_ms?: number | null
  dedupe_key?: string
  actions?: AlertAction[]
}

function isAlertKind(value: unknown): value is AlertKind {
  return value === 'error' || value === 'info' || value === 'success' || value === 'warning'
}

function isAlertPlacement(value: unknown): value is AlertPlacement {
  return value === 'toast' || value === 'banner'
}

function isAlertAction(value: unknown): value is AlertAction {
  if (typeof value !== 'object' || value === null) {
    return false
  }

  const candidate = value as Record<string, unknown>
  const style = candidate.style
  return typeof candidate.id === 'string'
    && typeof candidate.label === 'string'
    && (style === 'primary' || style === 'secondary' || style === 'danger')
    && typeof candidate.closeOnClick === 'boolean'
    && (candidate.href === undefined || typeof candidate.href === 'string')
}

export function normalizeAlertEventPayload(payload: unknown): AlertInput | null {
  if (typeof payload !== 'object' || payload === null) {
    return null
  }

  const candidate = payload as Record<string, unknown>
  if (!isAlertKind(candidate.kind) || !isAlertPlacement(candidate.placement)) {
    return null
  }
  if (typeof candidate.title !== 'string' || typeof candidate.message !== 'string') {
    return null
  }

  const actions = Array.isArray(candidate.actions)
    ? candidate.actions.filter((item) => isAlertAction(item))
    : []

  return {
    kind: candidate.kind,
    placement: candidate.placement,
    title: candidate.title,
    message: candidate.message,
    detail: typeof candidate.detail === 'string' ? candidate.detail : null,
    code: typeof candidate.code === 'string' ? candidate.code : null,
    dismissible: typeof candidate.dismissible === 'boolean' ? candidate.dismissible : undefined,
    autoClose: typeof candidate.auto_close === 'boolean' ? candidate.auto_close : undefined,
    durationMs: typeof candidate.duration_ms === 'number' ? candidate.duration_ms : undefined,
    dedupeKey: typeof candidate.dedupe_key === 'string' ? candidate.dedupe_key : undefined,
    actions,
    source: 'backend',
  }
}
