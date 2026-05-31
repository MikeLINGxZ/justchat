# Clipboard Paste & Drag-and-Drop Attachments Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Allow users to paste images/files from clipboard (Cmd+V) and drag files onto the chat input box to add them as attachments.

**Architecture:** A new Go `SaveTempFile` method accepts base64-encoded binary data, writes it to `<os.TempDir()>/lemontea/`, and returns the disk path. The frontend reads clipboard/dropped File objects via `Blob.arrayBuffer()`, encodes to base64, calls `SaveTempFile` via Wails IPC, then uses the returned path with the existing `inferAttachmentMeta` pipeline. Drag-and-drop and paste handlers both live on `div.chat-input-area` as React event handlers.

**Tech Stack:** Go (Wails v3), React 18, Tiptap, Vitest, @testing-library/react

---

## File Map

| Action | File |
|---|---|
| Create | `backend/service/file/file_dto/save_temp_file.go` |
| Modify | `backend/pkg/ierror/error.go` |
| Modify | `backend/pkg/i18n/resources_en_us.go` |
| Modify | `backend/pkg/i18n/resources_en_zh_cn.go` |
| Modify | `backend/service/file/file.go` |
| Modify | `backend/service/file/file_implement.go` |
| Create | `backend/service/file/file_temp_test.go` |
| Auto-regen | `frontend/bindings/.../file_dto/models.ts` |
| Auto-regen | `frontend/bindings/.../file.ts` |
| Modify | `frontend/src/i18n/locales/zh-CN.ts` |
| Modify | `frontend/src/i18n/locales/en.ts` |
| Modify | `frontend/src/components/chat/ChatInput.tsx` |
| Modify | `frontend/src/__tests__/chatInputAttachments.test.tsx` |

---

## Task 1: Backend DTO + error code + i18n strings

**Files:**
- Create: `backend/service/file/file_dto/save_temp_file.go`
- Modify: `backend/pkg/ierror/error.go`
- Modify: `backend/pkg/i18n/resources_en_us.go`
- Modify: `backend/pkg/i18n/resources_en_zh_cn.go`

- [ ] **Step 1: Create the DTO file**

Create `backend/service/file/file_dto/save_temp_file.go`:

```go
package file_dto

type SaveTempFileInput struct {
	Name string `json:"name"` // original filename, e.g. "screenshot.png"
	Data string `json:"data"` // base64-encoded file content
	Mime string `json:"mime"` // e.g. "image/png"
}

type SaveTempFileOutput struct {
	FilePath string `json:"file_path"`
}
```

- [ ] **Step 2: Add error code**

In `backend/pkg/ierror/error.go`, add after `ErrFileOpen`:

```go
ErrFileSaveTempFile errorCode = "ierror.file.save_temp_file"
```

- [ ] **Step 3: Add i18n string (English)**

In `backend/pkg/i18n/resources_en_us.go`, add after `"ierror.file.open"`:

```go
"ierror.file.save_temp_file": "Failed to save temporary file",
```

- [ ] **Step 4: Add i18n string (Chinese)**

In `backend/pkg/i18n/resources_en_zh_cn.go`, add after `"ierror.file.open"`:

```go
"ierror.file.save_temp_file": "保存临时文件失败",
```

- [ ] **Step 5: Verify Go compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 6: Commit**

```bash
git add backend/service/file/file_dto/save_temp_file.go \
        backend/pkg/ierror/error.go \
        backend/pkg/i18n/resources_en_us.go \
        backend/pkg/i18n/resources_en_zh_cn.go
git commit -m "feat(file): add SaveTempFile DTO and error code"
```

---

## Task 2: Backend SaveTempFile implementation + startup cleanup (TDD)

**Files:**
- Create: `backend/service/file/file_temp_test.go`
- Modify: `backend/service/file/file.go`
- Modify: `backend/service/file/file_implement.go`

- [ ] **Step 1: Write failing tests**

Create `backend/service/file/file_temp_test.go`:

```go
package file

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto"
)

func TestSaveTempFile_Valid(t *testing.T) {
	f := &File{}
	payload := []byte("hello world")
	encoded := base64.StdEncoding.EncodeToString(payload)

	out, err := f.SaveTempFile(context.Background(), file_dto.SaveTempFileInput{
		Name: "test.txt",
		Data: encoded,
		Mime: "text/plain",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || out.FilePath == "" {
		t.Fatal("expected non-empty file path")
	}
	t.Cleanup(func() { os.Remove(out.FilePath) })

	got, err := os.ReadFile(out.FilePath)
	if err != nil {
		t.Fatalf("read temp file: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(got))
	}
}

func TestSaveTempFile_InvalidBase64(t *testing.T) {
	f := &File{}
	_, err := f.SaveTempFile(context.Background(), file_dto.SaveTempFileInput{
		Name: "test.txt",
		Data: "!!!not-valid-base64!!!",
		Mime: "text/plain",
	})
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestCleanTempDir_RemovesFiles(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "lemontea")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for _, name := range []string{"a.png", "b.txt"} {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	cleanTempDir()

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			t.Errorf("expected file %s to be removed", e.Name())
		}
	}
}
```

- [ ] **Step 2: Run tests — confirm they fail**

```bash
go test ./backend/service/file/... -v -run "TestSaveTempFile|TestCleanTempDir"
```

Expected: FAIL with "undefined: File.SaveTempFile" or similar.

- [ ] **Step 3: Implement SaveTempFile**

In `backend/service/file/file.go`, add the following (after existing imports, add `encoding/base64`, `fmt`, `os`, `path/filepath`, `time`):

```go
import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto"

	"github.com/wailsapp/wails/v3/pkg/application"
)
```

Add method to `file.go`:

```go
// SaveTempFile decodes base64 data and writes it to the app temp directory.
// Returns the absolute file path on success.
func (f *File) SaveTempFile(ctx context.Context, input file_dto.SaveTempFileInput) (*file_dto.SaveTempFileOutput, error) {
	data, err := base64.StdEncoding.DecodeString(input.Data)
	if err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	tmpDir := filepath.Join(os.TempDir(), "lemontea")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	safeName := filepath.Base(input.Name)
	if safeName == "" || safeName == "." {
		safeName = "file"
	}
	name := fmt.Sprintf("%d-%s", time.Now().UnixNano(), safeName)
	path := filepath.Join(tmpDir, name)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	return &file_dto.SaveTempFileOutput{FilePath: path}, nil
}
```

- [ ] **Step 4: Add cleanTempDir + call from ServiceStartup**

Replace the content of `backend/service/file/file_implement.go`:

```go
package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (f *File) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	f.wailsApp = application.Get()
	cleanTempDir()
	return nil
}

// cleanTempDir removes all regular files from the app temp directory.
// Called at startup to prevent unbounded disk growth across sessions.
func cleanTempDir() {
	dir := filepath.Join(os.TempDir(), "lemontea")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}
```

- [ ] **Step 5: Run tests — confirm they pass**

```bash
go test ./backend/service/file/... -v -run "TestSaveTempFile|TestCleanTempDir"
```

Expected: PASS for all three tests.

- [ ] **Step 6: Verify full build still compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 7: Commit**

```bash
git add backend/service/file/file.go \
        backend/service/file/file_implement.go \
        backend/service/file/file_temp_test.go
git commit -m "feat(file): implement SaveTempFile with temp dir cleanup on startup"
```

---

## Task 3: Regenerate Wails TypeScript bindings

**Files:**
- Auto-regen: `frontend/bindings/.../file_dto/models.ts`
- Auto-regen: `frontend/bindings/.../file.ts`

- [ ] **Step 1: Run binding generation**

From project root:

```bash
task common:generate:bindings
```

Expected: command completes without errors.

- [ ] **Step 2: Verify new types appear**

```bash
grep -n "SaveTempFile" frontend/bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file.ts
grep -n "SaveTempFileInput" frontend/bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto/models.ts
```

Expected: both greps find at least one match.

- [ ] **Step 3: Commit**

```bash
git add frontend/bindings/
git commit -m "chore: regenerate Wails bindings for SaveTempFile"
```

---

## Task 4: Frontend i18n key

**Files:**
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: Add key to zh-CN.ts**

In `frontend/src/i18n/locales/zh-CN.ts`, inside the `input` object, add after `attachLimitSize`:

```typescript
dropFiles: '释放以添加文件',
```

- [ ] **Step 2: Add key to en.ts**

In `frontend/src/i18n/locales/en.ts`, inside the `input` object, add after `attachLimitSize`:

```typescript
dropFiles: 'Drop files to attach',
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts
git commit -m "feat(i18n): add input.dropFiles key for drag-and-drop overlay"
```

---

## Task 5: Frontend clipboard paste (TDD)

**Files:**
- Modify: `frontend/src/__tests__/chatInputAttachments.test.tsx`
- Modify: `frontend/src/components/chat/ChatInput.tsx`

- [ ] **Step 1: Update test mock to include SaveTempFile**

In `frontend/src/__tests__/chatInputAttachments.test.tsx`, update the existing `vi.mock` block for the file binding:

```typescript
vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  File: {
    SelectFile: vi.fn(),
    SaveTempFile: vi.fn(),
  },
}))
```

Also update the import at the bottom of the file to include `SaveTempFile` reference (the binding import line already exists — just confirm it's present):

```typescript
import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'
```

- [ ] **Step 2: Write failing paste tests**

Add the following `describe` block to `frontend/src/__tests__/chatInputAttachments.test.tsx`:

```typescript
describe('ChatInput clipboard paste', () => {
  beforeEach(() => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockReset()
  })

  it('calls SaveTempFile and adds chip when pasting an image', async () => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      file_path: '/tmp/lemontea/123-screenshot.png',
    })
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const file = new File(['fake-image-data'], 'screenshot.png', { type: 'image/png' })

    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: {
        items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }],
      },
    })
    container.dispatchEvent(pasteEvent)

    await waitFor(() => expect(screen.getByText('screenshot.png')).toBeInTheDocument())
    expect(FileBinding.SaveTempFile).toHaveBeenCalledOnce()
  })

  it('does not call SaveTempFile when pasting text only', async () => {
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'string', type: 'text/plain' }] },
    })
    container.dispatchEvent(pasteEvent)

    await new Promise(r => setTimeout(r, 50))
    expect(FileBinding.SaveTempFile).not.toHaveBeenCalled()
  })

  it('respects ATTACHMENT_MAX_COUNT when pasting', async () => {
    // Fill attachments to max via SelectFile first
    const mockSelect = FileBinding.SelectFile as ReturnType<typeof vi.fn>
    for (let i = 0; i < 10; i++) {
      mockSelect.mockResolvedValueOnce({ file_path: `/x/file${i}.png` })
    }
    const user = userEvent.setup()
    renderInput()
    for (let i = 0; i < 10; i++) {
      await user.click(screen.getByLabelText(/attach file|附加文件/i))
    }
    await waitFor(() => expect(screen.getAllByText(/file\d\.png/).length).toBe(10))

    // Now try to paste one more — should be ignored
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      file_path: '/tmp/lemontea/extra.png',
    })
    const container = document.querySelector('.chat-input-area')!
    const file = new File(['data'], 'extra.png', { type: 'image/png' })
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }] },
    })
    container.dispatchEvent(pasteEvent)

    await new Promise(r => setTimeout(r, 50))
    expect(FileBinding.SaveTempFile).not.toHaveBeenCalled()
  })

  it('skips file silently when SaveTempFile rejects', async () => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockRejectedValueOnce(new Error('disk full'))
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const file = new File(['data'], 'broken.png', { type: 'image/png' })
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }] },
    })
    container.dispatchEvent(pasteEvent)

    await new Promise(r => setTimeout(r, 100))
    expect(screen.queryByText('broken.png')).toBeNull()
  })
})
```

- [ ] **Step 3: Run tests — confirm they fail**

```bash
cd frontend && npx vitest run src/__tests__/chatInputAttachments.test.tsx
```

Expected: the 4 new paste tests fail (SaveTempFile not called / chip doesn't appear).

- [ ] **Step 4: Add blobToBase64 helper and paste handler to ChatInput.tsx**

Update the React import at the top of `frontend/src/components/chat/ChatInput.tsx` to include `useRef`:

```typescript
import { useState, useEffect, useRef } from 'react'
```

Add `SaveTempFileInput` to the existing file_dto import line:

```typescript
import { SelectFileInput, SaveTempFileInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto'
```

Add this module-level utility function just before `const SingleLineCodeBlock = Extension.create(...)`:

```typescript
async function blobToBase64(blob: Blob): Promise<string> {
  const buf = await blob.arrayBuffer()
  const bytes = new Uint8Array(buf)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) binary += String.fromCharCode(bytes[i])
  return btoa(binary)
}
```

Inside the `ChatInput` component function, after the `attachments` state declaration, add a ref that tracks current attachments (needed to read inside async callbacks without stale closure):

```typescript
const attachmentsRef = useRef<Attachment[]>([])
attachmentsRef.current = attachments
```

Then add the paste handler (add after `handleStop`):

```typescript
const handlePaste = async (event: React.ClipboardEvent<HTMLDivElement>) => {
  const items = Array.from(event.clipboardData.items)
  const fileItems = items.filter(item => item.kind === 'file')
  if (fileItems.length === 0) return

  event.preventDefault()

  const remaining = ATTACHMENT_MAX_COUNT - attachmentsRef.current.length
  const toProcess = fileItems.slice(0, remaining).map(item => {
    const blob = item.getAsFile()
    if (!blob) return null
    const name = blob.name || `pasted-${Date.now()}.${item.type.split('/')[1] ?? 'bin'}`
    return { blob, name, mime: item.type || 'application/octet-stream' }
  }).filter((x): x is NonNullable<typeof x> => x !== null)

  const results = await Promise.allSettled(
    toProcess.map(async ({ blob, name, mime }) => {
      const data = await blobToBase64(blob)
      const result = await FileBinding.SaveTempFile(new SaveTempFileInput({ name, data, mime }))
      if (!result?.file_path) throw new Error('no path')
      return inferAttachmentMeta(result.file_path)
    })
  )

  const newAtts = results
    .filter((r): r is PromiseFulfilledResult<Attachment> => r.status === 'fulfilled')
    .map(r => r.value)

  if (newAtts.length > 0) {
    setAttachments(prev => [...prev, ...newAtts])
  }
}
```

Wire the handler to the outer div in the JSX. Change the return statement's outermost div:

```typescript
return (
  <div
    className="chat-input-area shrink-0 pb-4 pt-2"
    onPaste={handlePaste}
  >
```

- [ ] **Step 5: Run tests — confirm paste tests pass**

```bash
cd frontend && npx vitest run src/__tests__/chatInputAttachments.test.tsx
```

Expected: all tests pass including the 4 new paste tests.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/chat/ChatInput.tsx \
        frontend/src/__tests__/chatInputAttachments.test.tsx
git commit -m "feat(chat): support pasting images and files from clipboard into input"
```

Expected: commit succeeds. Running `npx vitest run` from `frontend/` should show 7 passing tests (3 original + 4 new paste).
```

---

## Task 6: Frontend drag-and-drop (TDD)

**Files:**
- Modify: `frontend/src/__tests__/chatInputAttachments.test.tsx`
- Modify: `frontend/src/components/chat/ChatInput.tsx`

- [ ] **Step 1: Write failing drag-and-drop tests**

Add the following `describe` block to `frontend/src/__tests__/chatInputAttachments.test.tsx`:

```typescript
describe('ChatInput drag-and-drop', () => {
  beforeEach(() => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockReset()
  })

  it('adds chip when a file is dropped on the input area', async () => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      file_path: '/tmp/lemontea/123-doc.pdf',
    })
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const file = new File(['fake-pdf'], 'doc.pdf', { type: 'application/pdf' })
    const dropEvent = new DragEvent('drop', { bubbles: true, cancelable: true })
    Object.defineProperty(dropEvent, 'dataTransfer', {
      value: { files: [file], types: ['Files'] },
    })
    inputArea.dispatchEvent(dropEvent)

    await waitFor(() => expect(screen.getByText('doc.pdf')).toBeInTheDocument())
    expect(FileBinding.SaveTempFile).toHaveBeenCalledOnce()
  })

  it('shows drop overlay on dragover with files', () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', {
      value: { types: ['Files'] },
    })
    inputArea.dispatchEvent(dragoverEvent)

    expect(screen.getByText(/drop files to attach|释放以添加文件/i)).toBeInTheDocument()
  })

  it('hides drop overlay on dragleave when leaving the container', () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!

    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', { value: { types: ['Files'] } })
    inputArea.dispatchEvent(dragoverEvent)
    expect(screen.getByText(/drop files to attach|释放以添加文件/i)).toBeInTheDocument()

    const dragleaveEvent = new DragEvent('dragleave', { bubbles: true, cancelable: true, relatedTarget: document.body })
    inputArea.dispatchEvent(dragleaveEvent)
    expect(screen.queryByText(/drop files to attach|释放以添加文件/i)).toBeNull()
  })

  it('caps dropped files at ATTACHMENT_MAX_COUNT', async () => {
    // Return paths whose basename matches the dropped file name so chip text is predictable
    for (let i = 0; i < 10; i++) {
      ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
        file_path: `/tmp/lemontea/img${i}.png`,
      })
    }
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const files = Array.from({ length: 15 }, (_, i) => new File(['d'], `img${i}.png`, { type: 'image/png' }))
    const dropEvent = new DragEvent('drop', { bubbles: true, cancelable: true })
    Object.defineProperty(dropEvent, 'dataTransfer', {
      value: { files, types: ['Files'] },
    })
    inputArea.dispatchEvent(dropEvent)

    await waitFor(() => expect(screen.getAllByText(/img\d+\.png/).length).toBe(10))
    expect(FileBinding.SaveTempFile).toHaveBeenCalledTimes(10)
  })

  it('does not show overlay when dragging non-file content', () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', {
      value: { types: ['text/plain'] },
    })
    inputArea.dispatchEvent(dragoverEvent)

    expect(screen.queryByText(/drop files to attach|释放以添加文件/i)).toBeNull()
  })
})
```

- [ ] **Step 2: Run tests — confirm they fail**

```bash
cd frontend && npx vitest run src/__tests__/chatInputAttachments.test.tsx
```

Expected: the 5 new drag-and-drop tests fail.

- [ ] **Step 3: Add isDraggingOver state and drag handlers to ChatInput.tsx**

Inside the `ChatInput` component function, after the `attachmentsRef` line, add:

```typescript
const [isDraggingOver, setIsDraggingOver] = useState(false)
```

After `handlePaste`, add the drag handlers:

```typescript
const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
  if (!event.dataTransfer.types.includes('Files')) return
  event.preventDefault()
  setIsDraggingOver(true)
}

const handleDragLeave = (event: React.DragEvent<HTMLDivElement>) => {
  if ((event.currentTarget as HTMLDivElement).contains(event.relatedTarget as Node)) return
  setIsDraggingOver(false)
}

const handleDrop = async (event: React.DragEvent<HTMLDivElement>) => {
  event.preventDefault()
  setIsDraggingOver(false)

  const dropped = Array.from(event.dataTransfer.files)
  const remaining = ATTACHMENT_MAX_COUNT - attachments.length
  const toProcess = dropped.slice(0, remaining)

  const results = await Promise.allSettled(
    toProcess.map(async (file) => {
      const data = await blobToBase64(file)
      const result = await FileBinding.SaveTempFile(
        new SaveTempFileInput({ name: file.name, data, mime: file.type || 'application/octet-stream' })
      )
      if (!result?.file_path) throw new Error('no path')
      return inferAttachmentMeta(result.file_path)
    })
  )

  const newAtts = results
    .filter((r): r is PromiseFulfilledResult<Attachment> => r.status === 'fulfilled')
    .map(r => r.value)

  if (newAtts.length > 0) {
    setAttachments(prev => [...prev, ...newAtts])
  }
}
```

Wire the drag handlers to the outer div (the one already has `onPaste`):

```typescript
return (
  <div
    className="chat-input-area shrink-0 pb-4 pt-2"
    onPaste={handlePaste}
    onDragOver={handleDragOver}
    onDragLeave={handleDragLeave}
    onDrop={handleDrop}
  >
```

- [ ] **Step 4: Add visual drop overlay to the inner rounded container**

Find the inner `rounded-2xl` div in the JSX. Change it to:

```typescript
<div className={cn(
  "rounded-2xl border border-border bg-background shadow-sm focus-within:border-primary/40 transition-colors relative",
  isDraggingOver && "border-primary"
)}>
  {isDraggingOver && (
    <div className="absolute inset-0 z-10 rounded-2xl bg-primary/5 flex items-center justify-center pointer-events-none">
      <span className="text-sm text-primary font-medium">{t('input.dropFiles')}</span>
    </div>
  )}
  {/* rest of content unchanged */}
```

- [ ] **Step 5: Run tests — confirm all pass**

```bash
cd frontend && npx vitest run src/__tests__/chatInputAttachments.test.tsx
```

Expected: all 12 tests pass (3 original + 4 paste + 5 drag-and-drop).

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/chat/ChatInput.tsx \
        frontend/src/__tests__/chatInputAttachments.test.tsx
git commit -m "feat(chat): support drag-and-drop files onto input with visual drop overlay"
```
