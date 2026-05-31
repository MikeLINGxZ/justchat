# Settings Content Layout Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix the settings content area so header is pinned at top, content scrolls, and footprint stays fixed at bottom.

**Architecture:** The flex height chain in `SettingsShell` is broken — the wrapper div around `props.children` is not a flex container, so `SettingsContentLayout`'s `flex-1` root has no effect and height is driven by content instead of being constrained. Adding `flex flex-col` to that wrapper div closes the gap.

**Tech Stack:** React, Tailwind CSS, Vitest + jsdom

---

### Task 1: Fix wrapper div in SettingsShell

**Files:**
- Modify: `frontend/src/components/settings/SettingsShell.tsx:85`

- [ ] **Step 1: Apply the fix**

In `frontend/src/components/settings/SettingsShell.tsx`, find line 85:

```tsx
<div className="min-h-0 min-w-0 flex-1">
```

Change it to:

```tsx
<div className="flex min-h-0 min-w-0 flex-1 flex-col">
```

Full context for orientation (lines 83–88):

```tsx
        </div>
        <div className="flex min-h-0 min-w-0 flex-1 flex-col">
          {props.children}
        </div>
      </main>
      )}
```

- [ ] **Step 2: Run existing tests**

```bash
cd frontend && npx vitest run src/__tests__/settingsShell.test.tsx
```

Expected: all 3 tests pass (menu collapse/expand behavior is unaffected by this change).

- [ ] **Step 3: Visual verification**

Start the app and open Settings. Verify:
1. Header (title + description) stays pinned at the top when scrolling.
2. Content area scrolls when there is enough content to overflow.
3. Footprint (action buttons: 应用 / 重置 / 删除 etc.) stays pinned at the bottom and is never pushed off screen.

Test with DisplaySettings, LocaleSettings, FileSettings, and ProviderDetail views.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/settings/SettingsShell.tsx
git commit -m "fix: restore flex height chain in SettingsShell so content area scrolls and footprint stays fixed"
```
