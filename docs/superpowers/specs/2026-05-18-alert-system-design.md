# Alert System Design

## Background

`docs/dev/09.alert.md` defines a unified in-app alert capability covering errors, informational prompts, actionable prompts, deduplication, accessibility, and internationalization. The existing codebase already has shared theme, font size, i18n, Zustand stores, and Wails event infrastructure, but it does not yet have a centralized alert system that can be called consistently from both frontend code and backend event emitters.

This design defines a single alert center for the frontend and a single Wails event protocol for backend-to-frontend alert delivery. It follows the constraints from `docs/dev/00.rules.md`, especially:

- Reuse the existing project structure and coding style.
- Keep theme, font size, and i18n support consistent across the capability.
- Use explicit TypeScript types.
- Reuse Wails events by standardizing event IDs in `pkg/event_id`.

## Goals

- Build one unified alert center in the frontend.
- Support two presentation modes in the first version:
  - bottom-right stacked toast alerts
  - top banner alerts
- Allow frontend-local code to trigger alerts through the same store API.
- Allow backend code to trigger alerts by emitting a Wails event defined in `pkg/event_id`.
- Support the behaviors required by the alert document:
  - clear severity styling
  - optional title and message
  - expandable detail for errors
  - repeat deduplication with count aggregation
  - configurable auto close with max duration limit
  - actionable alerts with up to two primary actions
  - manual close and tracked close reasons
  - accessibility and keyboard support
  - i18n-ready copy

## Non-Goals

- No new backend Wails service API for alerts in the first version.
- No generic analytics pipeline implementation in this change. The design preserves event hooks and close reasons so analytics can be added later.
- No direct execution of destructive actions inside the alert itself. Alerts may guide the user into a confirmation flow, but not replace it.
- No mobile-specific gesture implementation in this first desktop-focused version.

## Existing Context

The implementation should align with the current frontend architecture:

- `frontend/src/store` uses Zustand for state.
- `frontend/src/components/providers` hosts app-wide providers.
- `frontend/src/types` contains shared explicit types.
- `frontend/src/i18n` already supports `zh-CN` and `en`.
- `frontend/src/App.tsx` mounts global providers near the root.
- Wails events are already consumed from `@wailsio/runtime`, especially in chat flows.

This makes a provider + store + typed event adapter pattern the best fit for alerts.

## Proposed Architecture

### Overview

The system will use one normalized alert model and one centralized frontend store. All alert sources flow through that store:

1. Frontend code calls a typed alert API.
2. Backend code emits one standardized Wails event.
3. The frontend event listener converts the incoming payload to the normalized alert model.
4. The alert store handles deduplication, lifecycle, queueing, close reasons, and placement routing.
5. The UI renders the same alert data into either a toast stack or a banner region.

This keeps behavior consistent across all entry points and avoids duplicating deduplication, timing, and accessibility logic.

### Frontend Module Layout

The first version should introduce a dedicated alert module under `frontend/src`, following the existing structure style:

- `frontend/src/alert/types.ts`
  - normalized alert types
  - event payload types
  - action and close reason types
- `frontend/src/alert/store.ts`
  - centralized Zustand alert store
- `frontend/src/alert/event.ts`
  - mapping from Wails event payload to normalized frontend data
- `frontend/src/components/providers/AlertEventProvider.tsx`
  - subscribes to the Wails event and forwards payloads into the store
- `frontend/src/components/alert/AlertViewport.tsx`
  - top-level viewport mounted once in app root
- `frontend/src/components/alert/AlertToastStack.tsx`
  - bottom-right stack renderer
- `frontend/src/components/alert/AlertBannerRegion.tsx`
  - top banner renderer
- `frontend/src/components/alert/AlertCard.tsx`
  - shared visual body used by both placements

This structure keeps the module cohesive while matching the provider and component organization that already exists in the repository.

### Backend Event Definition

The backend will not gain a new alert service in this version. Instead, a single standardized Wails event ID will be added to `pkg/event_id`.

Recommended event name:

- `AppAlert`

Recommended emitted string value:

- `app:alert`

Backend code that needs to display an alert will emit the Wails event with the standardized payload. This keeps the backend thin and matches the requested design of “backend emit, frontend listen, frontend display”.

## Data Model

### Normalized Frontend Alert Item

The frontend store will normalize all alerts into a single `AlertItem` type:

- `id: string`
- `kind: 'error' | 'info' | 'success' | 'warning'`
- `placement: 'toast' | 'banner'`
- `title: string`
- `message: string`
- `detail: string | null`
- `code: string | null`
- `actions: AlertAction[]`
- `dismissible: boolean`
- `autoClose: boolean`
- `durationMs: number | null`
- `count: number`
- `dedupeKey: string`
- `createdAt: number`
- `updatedAt: number`
- `ariaPriority: 'assertive' | 'polite'`
- `source: 'frontend' | 'backend'`

Supporting action type:

- `id: string`
- `label: string`
- `style: 'primary' | 'secondary' | 'danger'`
- `closeOnClick: boolean`
- `href?: string`

Supporting close reason type:

- `user`
- `timeout`
- `programmatic`
- `replaced`
- `action`

### Backend Event Payload

The backend event payload should remain declarative and UI-agnostic. Recommended fields:

- `kind`
- `placement`
- `title`
- `message`
- `detail`
- `code`
- `dismissible`
- `auto_close`
- `duration_ms`
- `dedupe_key`
- `actions`

The frontend is responsible for filling in missing defaults and deriving store-managed fields such as `id`, `count`, timestamps, and ARIA priority.

## Behavior Rules

### Severity and Styling

- `error` uses strong warning/destructive styling and the highest ARIA priority.
- `warning` uses caution styling.
- `success` uses positive styling.
- `info` uses neutral styling.

Both light and dark themes must be supported through the existing CSS variable system. Visual sizing must respect the existing font-size provider automatically by relying on `rem`-based styling and the app root font size.

### Placements

Two placements are supported in version one:

- `toast`
  - rendered as a bottom-right stacked list
- `banner`
  - rendered near the top of the page as a single active banner

The same normalized alert model is used for both placements. Placement controls only routing and layout, not behavior semantics.

### Default Lifetime Rules

- `error`
  - defaults to no auto close
- `info`
  - defaults to auto close after `5000ms`
- `success`
  - defaults to auto close after `5000ms`
- `warning`
  - defaults to auto close after `5000ms`

Additional rules:

- If the alert has actions, default `autoClose` becomes `false` unless explicitly overridden.
- Any configured duration must be clamped to a maximum of `15000ms`.
- When a duplicate alert is merged, its visible lifetime resets from the latest occurrence.

### Deduplication and Aggregation

Deduplication happens inside the centralized store.

Preferred key behavior:

1. Use explicit `dedupeKey` when provided.
2. Otherwise derive a stable composite key from:
   - `kind`
   - `placement`
   - `title`
   - `message`
   - `detail`

When a duplicate is detected:

- do not create a second rendered alert
- increment `count`
- update `updatedAt`
- reset auto close timing if the alert is auto-closing

Rendered duplicate count will be shown as `xN`.

### Visible Limits and Queue Strategy

Limits should be handled per placement instead of globally:

- toast limit: `3`
- banner limit: `1`

Overflow strategy:

- `toast`
  - keep the highest-priority and newest relevant items visible
  - remove older lower-priority items when the limit is exceeded
- `banner`
  - only one banner is visible at a time
  - a higher-priority incoming banner replaces the current banner
  - otherwise later banners wait in queue order

Priority order:

- `error > warning > success > info`

This satisfies the requirement that simultaneous alerts stay bounded and follow a clear handling policy.

### Actions

Alerts can expose at most two primary actions.

First-version supported action patterns:

- local callback actions wired by frontend code
- navigation or external-link actions via `href`

Rules:

- Actions must use i18n labels.
- Clicking an action may optionally close the alert based on `closeOnClick`.
- Dangerous or irreversible flows must not execute directly in the alert. The alert may instead open an existing confirmation dialog or route into a separate confirmation step.

### Detail Expansion

Error alerts may display optional details:

- collapsed by default
- expandable by the user
- may include error code and suggested next step

The UI should also support a copy-details action when details exist. The copied text should contain title, message, code, and detail in a readable format.

### Manual Close and Programmatic Removal

Any dismissible alert may be manually closed by the user. Alerts may also be removed by timeout or programmatic actions. All removals must carry a normalized close reason:

- `user`
- `timeout`
- `programmatic`
- `replaced`
- `action`

This is required so later analytics or debugging can distinguish how alerts disappear.

## Accessibility

Accessibility is a first-class behavior requirement.

Rules for the first version:

- `error` alerts should use assertive live region behavior.
- non-error alerts should use polite live region behavior.
- interactive controls must be keyboard reachable.
- dismiss button must be keyboard operable.
- action buttons must have clear accessible names.
- details toggle must expose expanded state.
- banner and toast regions should not trap keyboard focus.

Recommended role mapping:

- `error`
  - `role="alert"`
- `info`, `success`, `warning`
  - `role="status"` where appropriate, or live region semantics equivalent to polite announcement

## Internationalization

All frontend-owned strings must go through the existing i18n pipeline in `zh-CN` and `en`.

This includes:

- default action labels
- dismiss label
- copy-details label
- show-details and hide-details labels
- repeated count label formatting if needed
- timeout-related helper copy if displayed

Backend-provided content may already be localized before emission, but frontend fallback labels must always come from i18n.

## Root Integration

`App.tsx` should mount the alert capability near the existing global providers:

- mount `AlertEventProvider`
- mount `AlertViewport`

This ensures:

- alerts work in the main layout
- alerts work in settings-related entries
- the event listener is established once per window

## Frontend API

The frontend should expose a typed API on top of the store for local usage, such as:

- `pushAlert(payload)`
- `closeAlert(id, reason)`
- `closeAllAlerts(reason)`

Optional convenience wrappers can be added if they remain thin:

- `showErrorAlert(...)`
- `showSuccessAlert(...)`
- `showInfoAlert(...)`
- `showWarningAlert(...)`

These wrappers must still funnel into the same normalized store path.

## Testing Strategy

Implementation must follow TDD. Tests should be written before production code.

### Store Tests

Create focused tests for the alert store covering:

- duplicate aggregation increments count
- duplicate aggregation resets timeout
- error default does not auto close
- info and success defaults auto close at `5000ms`
- duration is clamped at `15000ms`
- toast limit enforces max three visible items
- banner limit enforces one visible item
- close reasons are preserved on removal

### Event Adapter Tests

Test the event bridge in isolation:

- backend payload maps correctly to normalized alert data
- missing optional fields receive correct defaults
- malformed or partial payloads do not crash the app

### UI Tests

Render-based tests should cover:

- toast stack renders bottom-right items
- banner region renders top banner item
- duplicate count renders as `xN`
- error alerts expose higher-priority live semantics
- detail expansion works
- copy-details action is shown only when detail exists
- no more than two actions are rendered
- manual dismiss works from keyboard interaction

## Risks and Mitigations

### Risk: duplicated logic across placements

Mitigation:

- keep all lifecycle and deduplication logic in one store
- keep card body shared between toast and banner views

### Risk: backend payload drift

Mitigation:

- define one canonical event payload type in frontend code
- define one canonical event ID in `pkg/event_id`
- document the payload format near the listener and tests

### Risk: timer-heavy behavior becoming flaky in tests

Mitigation:

- use fake timers for store lifecycle tests
- isolate timing logic in small store helpers where possible

### Risk: i18n omissions

Mitigation:

- add all fallback UI labels together in both locale files
- include UI assertions in tests using translated labels where practical

## Implementation Summary

The implementation should introduce one centralized alert system that:

- receives alerts from frontend code and backend Wails events
- normalizes them into one typed store
- deduplicates repeated alerts with count aggregation
- supports timeout and manual lifecycle handling
- renders both toast and banner views from the same source of truth
- respects theme, font size, i18n, and accessibility constraints

This design keeps the first version complete enough to satisfy the alert requirements while avoiding unnecessary backend surface area.
