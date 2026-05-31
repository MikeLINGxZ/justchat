# Memory Core MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first usable long-term memory loop for Lemontea: store memories, manage them from settings, inject them into chat turns, and asynchronously encode new memories.

**Architecture:** Use the app's existing SQLite/GORM storage as the source of truth. Render core and retrieval memory per chat turn into hidden context instead of mutating the cached runner instruction. Use a bounded background encoder with a structured JSON LLM response and backend validation.

**Tech Stack:** Go, GORM SQLite, Wails service bindings, React, Zustand, i18next.

---

### Task 1: Backend Memory Storage

**Files:**
- Create: `backend/models/data_models/memory.go`
- Create: `backend/storage/memory.go`
- Test: `backend/storage/memory_test.go`
- Modify: `backend/storage/storage.go`

- [ ] Add memory and embedding data models with fields from `docs/dev/17.memory_core.md`, plus `importance`, `confidence`, and `pinned`.
- [ ] Write storage tests covering create, update, soft delete, restore, filters, stats, core rendering order, and LIKE search.
- [ ] Run `go test ./backend/storage` and verify the tests fail before implementation.
- [ ] Implement storage methods and AutoMigrate registration.
- [ ] Run `go test ./backend/storage` and verify the tests pass.

### Task 2: Backend Memory Service

**Files:**
- Create/expand: `backend/service/memory/*.go`
- Create: `backend/service/memory/memory_dto/*.go`
- Test: `backend/service/memory/memory_test.go`
- Modify: `main.go`

- [ ] Add Wails-style methods: list, get, create, update, forget, restore, stats.
- [ ] Wrap public errors with `ierror`.
- [ ] Register the service in `main.go`.
- [ ] Run `go test ./backend/service/memory`.

### Task 3: Chat Runtime Integration

**Files:**
- Modify: `backend/pkg/agent/chat_handler.go`
- Modify: `backend/pkg/agent/manager.go`
- Test: `backend/pkg/agent/*memory*_test.go`

- [ ] Add a memory provider interface to avoid coupling chat runtime to Wails service.
- [ ] Before each run, render core memory and search retrieval memory, then add a fenced memory block to hidden run content.
- [ ] After stream completion, enqueue bounded async encoding using the latest user message and assistant response.
- [ ] Skip encoding for attachment-heavy turns and failures.
- [ ] Run `go test ./backend/pkg/agent`.

### Task 4: Frontend Management UI

**Files:**
- Modify: `frontend/src/types/settings.ts`
- Modify: `frontend/src/store/settingsStore.ts`
- Modify: `frontend/src/components/settings/SettingsPrimaryMenu.tsx`
- Create: `frontend/src/components/settings/memory/MemorySettingsView.tsx`
- Modify: `frontend/src/i18n/locales/en.ts`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`

- [ ] Add a memory primary settings tab.
- [ ] Add list, filter, edit, forget, restore, and create controls.
- [ ] Add i18n strings for English and Simplified Chinese.
- [ ] Run `npm test -- --run`.

### Task 5: Verification

- [ ] Run `go test ./backend/...`.
- [ ] Run `npm test -- --run`.
- [ ] Run `go build ./...`.
- [ ] Report any remaining gaps, especially Wails binding generation if unavailable in the current shell.
