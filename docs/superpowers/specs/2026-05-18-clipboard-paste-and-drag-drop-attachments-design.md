# Clipboard Paste & Drag-and-Drop Attachments

**Date:** 2026-05-18  
**Branch:** chat-file  
**Status:** Approved

## Overview

Add two new ways to attach files to the chat input:

1. **Clipboard paste** — pressing Cmd/Ctrl+V while the editor is focused pastes image or file attachments from the clipboard.
2. **Drag-and-drop** — dragging files from the OS onto the chat input area adds them as attachments.

Both features reuse the existing `Attachment` model and `inferAttachmentMeta` pipeline. The core enabler is a new backend `SaveTempFile` API that bridges browser binary data (Blob/base64) to on-disk file paths that the existing attachment system requires.

## Architecture

The backend's `BuildUserMessage` reads attachments via `os.ReadFile(a.Path)`. Because the browser's File API and Clipboard API provide binary data but not native file system paths, both new input methods must convert binary data to a temp file on disk and obtain a path before creating an `Attachment` record.

```
Browser (Paste/Drop)
  └─ Read file content as base64 (FileReader)
       └─ IPC: File.SaveTempFile(name, base64, mime) → file_path
            └─ inferAttachmentMeta(file_path) → Attachment
                 └─ existing send flow (unchanged)
```

## Backend: `SaveTempFile` API

**Location:** `backend/service/file/file.go`

```go
func (f *File) SaveTempFile(ctx context.Context, input file_dto.SaveTempFileInput) (*file_dto.SaveTempFileOutput, error)
```

**DTOs** (`backend/service/file/file_dto/`):

```go
type SaveTempFileInput struct {
    Name string `json:"name"` // original filename, e.g. "screenshot.png"
    Data string `json:"data"` // base64-encoded file content
    Mime string `json:"mime"` // e.g. "image/png"
}

type SaveTempFileOutput struct {
    FilePath string `json:"file_path"`
}
```

**Temp directory:** `<os.TempDir()>/lemontea/`

**File naming:** `<unix-timestamp-ns>-<original-name>` to avoid collisions.

**Startup cleanup:** `ServiceStartup` deletes all regular files directly inside the temp directory (non-recursive) before the app begins accepting requests. This prevents unbounded disk growth across sessions.

**Error cases:**
- Invalid base64 → return error (frontend skips that file)
- Disk write failure → return error (frontend skips that file)

## Frontend: Clipboard Paste

**Location:** `frontend/src/components/chat/ChatInput.tsx`

Add `onPaste` handler to the editor container `div`. The handler runs before Tiptap's own paste processing.

**Logic:**
1. Inspect `event.clipboardData.items` for items with `kind === 'file'`.
2. If no file items exist, return without calling `preventDefault()` — Tiptap handles text paste normally.
3. If file items exist, call `event.preventDefault()` to suppress Tiptap's default.
4. For each file item (up to remaining capacity toward `ATTACHMENT_MAX_COUNT`):
   - Call `blob.arrayBuffer()` → convert to base64.
   - Call `File.SaveTempFile({ name, data, mime })`.
   - On success: call `inferAttachmentMeta(path)` and append to `attachments` state.
   - On failure: skip silently (no toast — consistent with existing attach error handling).
5. Files beyond the remaining capacity are silently ignored.

**Text paste is unaffected** — the guard exits early if no file items are found.

## Frontend: Drag-and-Drop

**Location:** `frontend/src/components/chat/ChatInput.tsx`

**Drop zone:** the `div.chat-input-area` wrapper (the outermost div of `ChatInput`).

**New state:** `isDraggingOver: boolean` (local `useState`).

**Event handlers on the wrapper div:**

| Event | Action |
|---|---|
| `onDragOver` | If `dataTransfer.types` includes `'Files'`, call `preventDefault()` and set `isDraggingOver(true)` |
| `onDragLeave` | If `event.relatedTarget` is outside the container, set `isDraggingOver(false)` |
| `onDrop` | `preventDefault()`, set `isDraggingOver(false)`, process dropped files |

**Drop processing** (same as paste, steps 4–5 above):
- Read each `File` from `event.dataTransfer.files` as base64.
- Call `File.SaveTempFile` and `inferAttachmentMeta`.
- Append to `attachments` state, capped at `ATTACHMENT_MAX_COUNT`.

**Visual feedback when `isDraggingOver` is true:**
- The `rounded-2xl` input box border changes to `primary` color.
- A semi-transparent overlay covers the input box area with centered text: "释放以添加文件" (i18n key: `input.dropFiles`).
- The overlay uses `pointer-events: none` so it does not interfere with drop event bubbling.

## i18n

Add one key to both locale files:

| Key | zh-CN | en |
|---|---|---|
| `input.dropFiles` | 释放以添加文件 | Drop files to attach |

## Error Handling

Both paste and drop handle `SaveTempFile` failures by silently skipping the failed file. This is consistent with how the existing `handleAttach` (file picker) handles failures — no error toast is shown. Remaining files in the same batch are still processed.

## Testing

**`frontend/src/__tests__/chatInputAttachments.test.tsx`** (extend existing file):

- Paste with image clipboard item → `SaveTempFile` called → attachment added
- Paste with plain text only → `SaveTempFile` not called → editor text updated normally
- Paste when at `ATTACHMENT_MAX_COUNT` → no new attachments added
- Paste with 3 items when 1 slot remains → only 1 attachment added
- `SaveTempFile` rejects → attachment skipped, no crash
- Drop single file → attachment added
- Drop multiple files → all added (capped at limit)
- Drop non-file drag data (e.g. text/plain drag) → `isDraggingOver` not set, no attachments

**`backend/service/file/` (new test file `file_temp_test.go`):**

- Valid base64 input → file written to temp dir → path returned
- Invalid base64 input → error returned
- `ServiceStartup` → pre-existing files in temp dir are removed
