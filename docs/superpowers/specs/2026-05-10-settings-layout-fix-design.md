# Settings Content Layout Fix

**Date:** 2026-05-10

## Problem

The settings content area does not scroll when content overflows, and the footprint (action bar) gets pushed off screen. Root cause: the wrapper `<div>` in `SettingsShell` is not a flex container, so `SettingsContentLayout`'s `flex-1` root has no effect — height is driven by content instead of being constrained by the parent.

## Layout Structure (desired)

```
h-screen (root)
  └─ main (flex-col, flex-1)
       └─ div wrapper (flex-col, flex-1)   ← fix here
            └─ SettingsContentLayout (flex-col, flex-1)
                 ├─ header  (shrink-0)     — fixed top, title + description
                 ├─ content (flex-1, overflow-y-auto) — scrollable
                 └─ footprint (shrink-0)   — fixed bottom, action buttons
```

- No border separator between content and footprint.
- Header style unchanged.
- All views already route action buttons through the `footprint` prop.

## Change

**File:** `frontend/src/components/settings/SettingsShell.tsx`

Add `flex flex-col` to the wrapper div around `{props.children}`:

```tsx
// before
<div className="min-h-0 min-w-0 flex-1">

// after
<div className="flex min-h-0 min-w-0 flex-1 flex-col">
```

No other files need changes.
