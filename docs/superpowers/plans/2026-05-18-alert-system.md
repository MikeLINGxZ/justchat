# Alert System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a unified alert system that supports toast and banner placement, frontend store-based triggering, and backend Wails-event triggering.

**Architecture:** Add a dedicated frontend alert module with one normalized store, one Wails event adapter, and two render regions. Mount the alert provider and viewport at the app root, define a shared backend event ID in `pkg/event_id`, and drive the implementation with Vitest store, adapter, and UI coverage.

**Tech Stack:** React 18, TypeScript, Zustand, Vitest, Testing Library, Wails runtime events, Tailwind CSS

---

### Task 1: Define alert types and store behavior with failing store tests

**Files:**
- Create: `frontend/src/__tests__/alertStore.test.ts`
- Create: `frontend/src/alert/types.ts`
- Create: `frontend/src/alert/store.ts`

- [ ] **Step 1: Write the failing test**

```ts
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd frontend && npx vitest run src/__tests__/alertStore.test.ts`
Expected: FAIL because `@/alert/store` does not exist yet

- [ ] **Step 3: Write minimal implementation**

Add explicit alert types and a focused Zustand alert store with:

- normalized alert shape
- placement split into `toasts` and `banners`
- default duration rules
- dedupe support
- timer registration and reset
- close reason capture
- placement limit enforcement

- [ ] **Step 4: Run test to verify it passes**

Run: `cd frontend && npx vitest run src/__tests__/alertStore.test.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/__tests__/alertStore.test.ts frontend/src/alert/types.ts frontend/src/alert/store.ts
git commit -m "feat(alert): add alert store"
```

### Task 2: Add backend event ID and frontend event adapter with failing adapter tests

**Files:**
- Create: `frontend/src/__tests__/alertEventProvider.test.tsx`
- Create: `frontend/src/alert/event.ts`
- Create: `frontend/src/components/providers/AlertEventProvider.tsx`
- Modify: `backend/pkg/event_id/*.go`

- [ ] **Step 1: Write the failing test**

```tsx
import { render } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { AlertEventProvider } from '@/components/providers/AlertEventProvider'
import { useAlertStore } from '@/alert/store'
import { APP_ALERT_EVENT } from '@/alert/event'

const handlers = new Map<string, (event: { data: unknown }) => void>()

vi.mock('@wailsio/runtime', () => ({
  Events: {
    On: vi.fn((name: string, handler: (event: { data: unknown }) => void) => {
      handlers.set(name, handler)
      return () => handlers.delete(name)
    }),
  },
}))

describe('AlertEventProvider', () => {
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
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd frontend && npx vitest run src/__tests__/alertEventProvider.test.tsx`
Expected: FAIL because provider and event adapter do not exist yet

- [ ] **Step 3: Write minimal implementation**

Add:

- one shared frontend event constant and payload normalizer
- one provider that subscribes to the Wails event and pushes alerts into the store
- one backend `pkg/event_id` constant for `app:alert`

- [ ] **Step 4: Run test to verify it passes**

Run: `cd frontend && npx vitest run src/__tests__/alertEventProvider.test.tsx`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/__tests__/alertEventProvider.test.tsx frontend/src/alert/event.ts frontend/src/components/providers/AlertEventProvider.tsx backend/pkg/event_id
git commit -m "feat(alert): add backend alert event bridge"
```

### Task 3: Build alert UI with failing render tests

**Files:**
- Create: `frontend/src/__tests__/alertViewport.test.tsx`
- Create: `frontend/src/components/alert/AlertViewport.tsx`
- Create: `frontend/src/components/alert/AlertToastStack.tsx`
- Create: `frontend/src/components/alert/AlertBannerRegion.tsx`
- Create: `frontend/src/components/alert/AlertCard.tsx`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`
- Modify: `frontend/src/index.css`

- [ ] **Step 1: Write the failing test**

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it } from 'vitest'
import { AlertViewport } from '@/components/alert/AlertViewport'
import { useAlertStore } from '@/alert/store'

describe('AlertViewport', () => {
  beforeEach(() => {
    useAlertStore.getState().reset()
  })

  it('renders toast and banner regions from the shared store', () => {
    useAlertStore.getState().pushAlert({
      kind: 'success',
      placement: 'toast',
      title: 'Saved',
      message: 'Profile saved',
    })
    useAlertStore.getState().pushAlert({
      kind: 'warning',
      placement: 'banner',
      title: 'Attention',
      message: 'Banner warning',
    })

    render(<AlertViewport />)

    expect(screen.getByText('Profile saved')).toBeInTheDocument()
    expect(screen.getByText('Banner warning')).toBeInTheDocument()
  })

  it('shows duplicate count and exposes error alert semantics', () => {
    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'toast',
      title: 'Failed',
      message: 'Request failed',
    })
    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'toast',
      title: 'Failed',
      message: 'Request failed',
    })

    render(<AlertViewport />)

    expect(screen.getByText('x2')).toBeInTheDocument()
    expect(screen.getByRole('alert')).toBeInTheDocument()
  })

  it('supports expanding details and dismissing from the keyboard', async () => {
    const user = userEvent.setup()

    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'toast',
      title: 'Failed',
      message: 'Request failed',
      detail: 'stack trace',
    })

    render(<AlertViewport />)

    await user.click(screen.getByRole('button', { name: /show details/i }))
    expect(screen.getByText('stack trace')).toBeInTheDocument()

    await user.tab()
    await user.keyboard('{Enter}')
    expect(screen.queryByText('Request failed')).not.toBeInTheDocument()
  })

  it('renders at most two action buttons', () => {
    useAlertStore.getState().pushAlert({
      kind: 'warning',
      placement: 'banner',
      title: 'Action needed',
      message: 'Choose',
      actions: [
        { id: 'one', label: 'One', style: 'primary', closeOnClick: false },
        { id: 'two', label: 'Two', style: 'secondary', closeOnClick: false },
        { id: 'three', label: 'Three', style: 'danger', closeOnClick: false },
      ],
    })

    render(<AlertViewport />)

    expect(screen.getByRole('button', { name: 'One' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Two' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Three' })).not.toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd frontend && npx vitest run src/__tests__/alertViewport.test.tsx`
Expected: FAIL because viewport components do not exist yet

- [ ] **Step 3: Write minimal implementation**

Build shared alert UI that:

- renders banners and toasts from the store
- uses translated labels
- exposes close button, detail toggle, and duplicate count
- applies theme-friendly styles in `index.css`
- limits rendered actions to two buttons

- [ ] **Step 4: Run test to verify it passes**

Run: `cd frontend && npx vitest run src/__tests__/alertViewport.test.tsx`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/__tests__/alertViewport.test.tsx frontend/src/components/alert frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts frontend/src/index.css
git commit -m "feat(alert): add alert viewport and UI"
```

### Task 4: Mount alert system in app root and verify integrated coverage

**Files:**
- Modify: `frontend/src/App.tsx`

- [ ] **Step 1: Write the failing test**

Add one integration assertion to `frontend/src/__tests__/alertEventProvider.test.tsx` or `frontend/src/__tests__/alertViewport.test.tsx` that renders `<App />` and confirms backend event delivery reaches the mounted viewport.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd frontend && npx vitest run src/__tests__/alertEventProvider.test.tsx src/__tests__/alertViewport.test.tsx`
Expected: FAIL because `App.tsx` does not mount the alert provider and viewport yet

- [ ] **Step 3: Write minimal implementation**

Mount:

- `AlertEventProvider`
- `AlertViewport`

inside `frontend/src/App.tsx` alongside the existing root providers.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd frontend && npx vitest run src/__tests__/alertEventProvider.test.tsx src/__tests__/alertViewport.test.tsx`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.tsx frontend/src/__tests__/alertEventProvider.test.tsx frontend/src/__tests__/alertViewport.test.tsx
git commit -m "feat(alert): mount global alert system"
```

### Task 5: Run full verification for the alert feature slice

**Files:**
- No additional files required unless failures demand fixes

- [ ] **Step 1: Run focused alert tests**

Run: `cd frontend && npx vitest run src/__tests__/alertStore.test.ts src/__tests__/alertEventProvider.test.tsx src/__tests__/alertViewport.test.tsx`
Expected: PASS

- [ ] **Step 2: Run the broader frontend test suite**

Run: `cd frontend && npx vitest run`
Expected: PASS

- [ ] **Step 3: Run frontend production build**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 4: Commit any final fixes if verification required them**

```bash
git add <fixed-files>
git commit -m "fix(alert): address verification issues"
```
If no fixes are needed, skip this commit.
