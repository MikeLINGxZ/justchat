# Error Message Display Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Display backend errors as expandable in-chat bubbles, and fix the blank-bubble bug when LLM streaming fails.

**Architecture:** Add a new `"error"` message content type that carries `{msg, detail}` JSON. Backend stores stream errors with this type; frontend renders them as a red bubble with a collapsible detail panel. Non-streaming Wails call errors are caught in `sendMessage` and inserted as local-only error messages with the same type.

**Tech Stack:** Go (backend), React + TypeScript + Zustand + Tailwind + Lucide (frontend), Vitest + Testing Library (frontend tests), Go `testing` package (backend tests)

---

## File Map

| File | Action | Purpose |
|------|--------|---------|
| `backend/pkg/ierror/ierrors.go` | Modify | Replace `Err error` with `Detail string`; remove `Unwrap()` |
| `backend/pkg/ierror/ierrors_test.go` | Create | Tests for IError struct behaviour |
| `backend/pkg/ierror/error.go` | Modify | Add `ErrAgentStreamError` constant |
| `backend/pkg/i18n/resources_en_us.go` | Modify | Add stream error English string |
| `backend/pkg/i18n/resources_en_zh_cn.go` | Modify | Add stream error Chinese string |
| `backend/pkg/agent/chat_handler.go` | Modify | Store stream error as `ContentType: "error"` with JSON content |
| `frontend/src/types/index.ts` | Modify | Add `"error"` to `MessageContentType` |
| `frontend/src/i18n/locales/en.ts` | Modify | Add `errorDetail` / `hideErrorDetail` keys |
| `frontend/src/i18n/locales/zh-CN.ts` | Modify | Add same keys in Chinese |
| `frontend/src/components/chat/ErrorMessageBubble.tsx` | Create | Error bubble UI component |
| `frontend/src/__tests__/errorMessageBubble.test.tsx` | Create | Component tests |
| `frontend/src/components/chat/MessageItem.tsx` | Modify | Add `"error"` branch + import |
| `frontend/src/__tests__/messageItem.test.tsx` | Modify | Add test for error content type |
| `frontend/src/store/chatStore.ts` | Modify | Add `parseIError`, fix stream:error handler, add sendMessage try/catch |
| `frontend/src/__tests__/chatStore.test.ts` | Modify | Tests for error handling in store |

---

## Task 1: Fix IError Struct

**Files:**
- Modify: `backend/pkg/ierror/ierrors.go`
- Create: `backend/pkg/ierror/ierrors_test.go`

- [ ] **Step 1: Write failing tests**

Create `backend/pkg/ierror/ierrors_test.go`:

```go
package ierror

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError_NilReturnsNil(t *testing.T) {
	result := Error(ErrAgentSendMessage, nil)
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestError_CreatesIErrorWithDetail(t *testing.T) {
	underlying := errors.New("connection refused")
	result := Error(ErrAgentSendMessage, underlying)

	var iErr *IError
	if !errors.As(result, &iErr) {
		t.Fatalf("expected *IError, got %T", result)
	}
	if iErr.Detail != "connection refused" {
		t.Errorf("Detail = %q, want %q", iErr.Detail, "connection refused")
	}
	if iErr.Msg == "" {
		t.Error("Msg should not be empty")
	}
}

func TestError_PassthroughExistingIError(t *testing.T) {
	original := &IError{Detail: "original detail", Msg: "original msg"}
	result := Error(ErrAgentSendMessage, original)

	var iErr *IError
	if !errors.As(result, &iErr) {
		t.Fatalf("expected *IError")
	}
	if iErr.Msg != "original msg" {
		t.Errorf("expected passthrough of original IError, got Msg=%q", iErr.Msg)
	}
}

func TestIError_JSONSerializesDetailAndMsg(t *testing.T) {
	e := &IError{Detail: "raw error text", Msg: "user message"}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out["detail"] != "raw error text" {
		t.Errorf("detail = %q, want %q", out["detail"], "raw error text")
	}
	if out["msg"] != "user message" {
		t.Errorf("msg = %q, want %q", out["msg"], "user message")
	}
}

func TestIError_IsMatchesByMsg(t *testing.T) {
	a := &IError{Detail: "d1", Msg: "same msg"}
	b := &IError{Detail: "d2", Msg: "same msg"}
	c := &IError{Detail: "d3", Msg: "different msg"}
	if !errors.Is(a, b) {
		t.Error("expected a.Is(b) == true for same Msg")
	}
	if errors.Is(a, c) {
		t.Error("expected a.Is(c) == false for different Msg")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/backend && go test ./pkg/ierror/...
```

Expected: FAIL (tests compile but some may fail because `IError` still has `Err error` field, not `Detail string`)

- [ ] **Step 3: Update `backend/pkg/ierror/ierrors.go`**

Replace the entire file:

```go
package ierror

import (
	"encoding/json"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

// IError is the standard error type returned by all Wails-bound service methods.
// Detail carries the raw error string; Msg is the localized user-facing message.
type IError struct {
	Detail string `json:"detail"`
	Msg    string `json:"msg"`
}

// Error returns the JSON representation of the IError.
func (e *IError) Error() string {
	errBytes, err := json.Marshal(e)
	if err != nil {
		return i18n.TCurrent("ierror.unknown_error", nil)
	}
	return string(errBytes)
}

// Is enables errors.Is to compare IError values by Msg.
func (e *IError) Is(target error) bool {
	if t, ok := target.(*IError); ok {
		return e.Msg == t.Msg
	}
	return false
}

// Error wraps err in an IError using the given error code for the localized Msg.
// Returns nil when err is nil. Passes through an existing *IError unchanged.
func Error(code errorCode, err error) error {
	if err == nil {
		return nil
	}

	var iErr *IError
	if errors.As(err, &iErr) {
		return iErr
	}

	return &IError{
		Detail: err.Error(),
		Msg:    code.Msg(),
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/backend && go test ./pkg/ierror/...
```

Expected: `ok gitlab.linhf.cn/.../pkg/ierror`

- [ ] **Step 5: Verify full backend still compiles**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/backend && go build ./...
```

Expected: no output (clean build)

- [ ] **Step 6: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/ierror/ierrors.go backend/pkg/ierror/ierrors_test.go && git commit -m "fix(ierror): replace Err error with Detail string for correct JSON serialization"
```

---

## Task 2: Add ErrAgentStreamError + i18n

**Files:**
- Modify: `backend/pkg/ierror/error.go`
- Modify: `backend/pkg/i18n/resources_en_us.go`
- Modify: `backend/pkg/i18n/resources_en_zh_cn.go`

- [ ] **Step 1: Add constant to `backend/pkg/ierror/error.go`**

In the Agent error codes `const` block, add after `ErrAgentAttachmentSize`:

```go
ErrAgentStreamError errorCode = "ierror.agent.stream_error"
```

The const block should look like:

```go
const (
	// Agent error codes
	ErrAgentSendMessage        errorCode = "ierror.agent.send_message"
	ErrAgentNoConfirm          errorCode = "ierror.agent.no_confirm"
	ErrAgentCreateSession      errorCode = "ierror.agent.create_session"
	ErrAgentListSessions       errorCode = "ierror.agent.list_sessions"
	ErrAgentListMessages       errorCode = "ierror.agent.list_messages"
	ErrAgentCountMessages      errorCode = "ierror.agent.count_messages"
	ErrAgentMarkRead           errorCode = "ierror.agent.mark_read"
	ErrAgentRename             errorCode = "ierror.agent.rename"
	ErrAgentDelete             errorCode = "ierror.agent.delete"
	ErrAgentToggleStar         errorCode = "ierror.agent.toggle_star"
	ErrAgentTooManyAttachments errorCode = "ierror.agent.too_many_attachments"
	ErrAgentAttachmentPath     errorCode = "ierror.agent.attachment_path_empty"
	ErrAgentAttachmentNotFound errorCode = "ierror.agent.attachment_not_found"
	ErrAgentAttachmentSize     errorCode = "ierror.agent.attachment_too_large"
	ErrAgentStreamError        errorCode = "ierror.agent.stream_error"

	// File error codes
	ErrFileSelectFolder errorCode = "ierror.file.select_folder"
	ErrFileSelectFile   errorCode = "ierror.file.select_file"

	// Settings error codes
	ErrSettingsLoadConfig   errorCode = "ierror.settings.load_config"
	ErrSettingsSaveConfig   errorCode = "ierror.settings.save_config"
	ErrSettingsCreateDir    errorCode = "ierror.settings.create_dir"
	ErrSettingsCopyFile     errorCode = "ierror.settings.copy_file"
	ErrSettingsReadConfig   errorCode = "ierror.settings.read_config"
	ErrSettingsParseConfig  errorCode = "ierror.settings.parse_config"
	ErrSettingsWriteLocator errorCode = "ierror.settings.write_locator"
	ErrSettingsTargetDir    errorCode = "ierror.settings.target_dir_required"

	// Provider error codes
	ErrProviderCreate        errorCode = "ierror.provider.create"
	ErrProviderCreateModels  errorCode = "ierror.provider.create_models"
	ErrProviderListProviders errorCode = "ierror.provider.list_providers"
	ErrProviderListModels    errorCode = "ierror.provider.list_models"
	ErrProviderDelete        errorCode = "ierror.provider.delete"
	ErrProviderUpdate        errorCode = "ierror.provider.update"
	ErrProviderDeleteModel   errorCode = "ierror.provider.delete_model"
	ErrProviderInvalidModel  errorCode = "ierror.provider.invalid_model_id"
	ErrProviderSetDefault    errorCode = "ierror.provider.set_default"
	ErrProviderFetchModels   errorCode = "ierror.provider.fetch_model_list"
)
```

- [ ] **Step 2: Add English string to `backend/pkg/i18n/resources_en_us.go`**

In the `// Agent` section, add after `"ierror.agent.attachment_too_large"`:

```go
"ierror.agent.stream_error": "Model response error",
```

- [ ] **Step 3: Add Chinese string to `backend/pkg/i18n/resources_en_zh_cn.go`**

In the `// Agent` section, add after `"ierror.agent.attachment_too_large"`:

```go
"ierror.agent.stream_error": "大模型响应出错",
```

- [ ] **Step 4: Compile check**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/backend && go build ./...
```

Expected: no output

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/ierror/error.go backend/pkg/i18n/resources_en_us.go backend/pkg/i18n/resources_en_zh_cn.go && git commit -m "feat(ierror): add ErrAgentStreamError code with i18n strings"
```

---

## Task 3: Update chat_handler.go Stream Error Storage

**Files:**
- Modify: `backend/pkg/agent/chat_handler.go`

- [ ] **Step 1: Add `ierror` import**

In `backend/pkg/agent/chat_handler.go`, add the import path to the import block:

```go
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	agentpkg "trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)
```

- [ ] **Step 2: Replace the `evt.Error != nil` block**

Find the current block (around line 249):

```go
if evt.Error != nil {
    ch.emitEvent(event_id.AgentStreamError, map[string]any{
        "sessionId": sessionID,
        "error":     evt.Error.Message,
    })
    ch.updateSessionStatus(sessionID, "error-unread")
    _, _ = stor.CreateMessage(data_models.Message{
        SessionID:   sessionID,
        Role:        "assistant",
        ContentType: "text",
        Content:     fmt.Sprintf("Error: %s", evt.Error.Message),
        ModelName:   modelName,
    })
    return
}
```

Replace with:

```go
if evt.Error != nil {
    ch.emitEvent(event_id.AgentStreamError, map[string]any{
        "sessionId": sessionID,
        "error":     evt.Error.Message,
    })
    ch.updateSessionStatus(sessionID, "error-unread")
    errContent, _ := json.Marshal(map[string]string{
        "msg":    ierror.ErrAgentStreamError.Msg(),
        "detail": evt.Error.Message,
    })
    _, _ = stor.CreateMessage(data_models.Message{
        SessionID:   sessionID,
        Role:        "assistant",
        ContentType: "error",
        Content:     string(errContent),
        ModelName:   modelName,
    })
    return
}
```

- [ ] **Step 3: Remove unused `fmt` import if needed**

Check if `fmt` is still used elsewhere in the file. If the only usage was `fmt.Sprintf("Error: %s", ...)`, remove `"fmt"` from the imports. (Scan the file for other `fmt.` usages to confirm — if any remain, keep it.)

- [ ] **Step 4: Compile check**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/backend && go build ./...
```

Expected: no output

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/chat_handler.go && git commit -m "feat(agent): store stream errors as structured error content type"
```

---

## Task 4: Frontend Types + i18n

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/i18n/locales/en.ts`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`

- [ ] **Step 1: Add `"error"` to `MessageContentType` in `frontend/src/types/index.ts`**

Change:

```typescript
export type MessageContentType =
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'thinking'
  | 'confirm_request'
  | 'confirm_response'
```

To:

```typescript
export type MessageContentType =
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'thinking'
  | 'confirm_request'
  | 'confirm_response'
  | 'error'
```

- [ ] **Step 2: Add i18n keys to `frontend/src/i18n/locales/en.ts`**

In the `chat` section, add after `attachmentRemove`:

```typescript
errorDetail: 'Show error details',
hideErrorDetail: 'Hide error details',
```

- [ ] **Step 3: Add i18n keys to `frontend/src/i18n/locales/zh-CN.ts`**

In the `chat` section, add after `attachmentRemove`:

```typescript
errorDetail: '展开错误详情',
hideErrorDetail: '收起错误详情',
```

- [ ] **Step 4: Verify TypeScript compiles**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx tsc --noEmit
```

Expected: no errors

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add frontend/src/types/index.ts frontend/src/i18n/locales/en.ts frontend/src/i18n/locales/zh-CN.ts && git commit -m "feat(frontend): add error message content type and i18n keys"
```

---

## Task 5: Create ErrorMessageBubble Component

**Files:**
- Create: `frontend/src/components/chat/ErrorMessageBubble.tsx`
- Create: `frontend/src/__tests__/errorMessageBubble.test.tsx`

- [ ] **Step 1: Write failing tests**

Create `frontend/src/__tests__/errorMessageBubble.test.tsx`:

```typescript
import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { ErrorMessageBubble } from '@/components/chat/ErrorMessageBubble'

describe('ErrorMessageBubble', () => {
  it('renders the user-facing msg from JSON content', () => {
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.getByText('Model response error')).toBeInTheDocument()
  })

  it('does not show detail text by default', () => {
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.queryByText('rate limit exceeded')).not.toBeInTheDocument()
  })

  it('expands detail when the icon button is clicked', () => {
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    const button = screen.getByRole('button')
    fireEvent.click(button)
    expect(screen.getByText('rate limit exceeded')).toBeInTheDocument()
  })

  it('collapses detail on second click', () => {
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    const button = screen.getByRole('button')
    fireEvent.click(button)
    fireEvent.click(button)
    expect(screen.queryByText('rate limit exceeded')).not.toBeInTheDocument()
  })

  it('renders raw string content when JSON parse fails', () => {
    render(<ErrorMessageBubble content="plain error text" />)
    expect(screen.getByText('plain error text')).toBeInTheDocument()
  })

  it('shows no expand button when detail is empty', () => {
    const content = JSON.stringify({ msg: 'Something failed', detail: '' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/errorMessageBubble.test.tsx
```

Expected: FAIL with "Cannot find module '@/components/chat/ErrorMessageBubble'"

- [ ] **Step 3: Create `frontend/src/components/chat/ErrorMessageBubble.tsx`**

```typescript
import { useState } from 'react'
import { AlertCircle, ChevronDown, ChevronUp } from 'lucide-react'
import { useTranslation } from 'react-i18next'

interface Props {
  content: string
}

export function ErrorMessageBubble({ content }: Props) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)

  let msg = content
  let detail = ''
  try {
    const parsed = JSON.parse(content) as { msg?: string; detail?: string }
    msg = parsed.msg ?? content
    detail = parsed.detail ?? ''
  } catch {
    // fallback: display raw string
  }

  return (
    <div className="flex flex-col gap-1 max-w-[70%]">
      <div className="flex items-center gap-2 rounded-2xl rounded-bl-sm bg-destructive/10 border border-destructive/20 text-destructive px-4 py-2.5 text-sm">
        <span className="flex-1 select-text">{msg}</span>
        {detail && (
          <button
            onClick={() => setExpanded(!expanded)}
            className="flex-shrink-0 flex items-center gap-0.5 text-destructive/70 hover:text-destructive transition-colors"
            title={expanded ? t('chat.hideErrorDetail') : t('chat.errorDetail')}
          >
            <AlertCircle size={14} />
            {expanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
          </button>
        )}
      </div>
      {expanded && detail && (
        <div className="rounded-lg bg-muted/50 border border-border px-3 py-2 text-xs text-muted-foreground font-mono whitespace-pre-wrap break-all select-text">
          {detail}
        </div>
      )}
    </div>
  )
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/errorMessageBubble.test.tsx
```

Expected: all 6 tests PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add frontend/src/components/chat/ErrorMessageBubble.tsx frontend/src/__tests__/errorMessageBubble.test.tsx && git commit -m "feat(chat): add ErrorMessageBubble component with expand/collapse detail"
```

---

## Task 6: Update MessageItem

**Files:**
- Modify: `frontend/src/components/chat/MessageItem.tsx`
- Modify: `frontend/src/__tests__/messageItem.test.tsx`

- [ ] **Step 1: Add failing test to `frontend/src/__tests__/messageItem.test.tsx`**

Append this test block to the existing `describe('MessageItem', ...)` block:

```typescript
it('renders ErrorMessageBubble for error content type', () => {
  const errorMessage: Message = {
    ...baseMessage,
    id: 99,
    role: 'assistant',
    contentType: 'error',
    content: JSON.stringify({ msg: 'Model response error', detail: 'timeout' }),
  }
  render(<MessageItem message={errorMessage} />)
  expect(screen.getByText('Model response error')).toBeInTheDocument()
})
```

- [ ] **Step 2: Run to confirm the test fails**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/messageItem.test.tsx
```

Expected: FAIL on the new test ("Model response error" not found)

- [ ] **Step 3: Update `frontend/src/components/chat/MessageItem.tsx`**

Add import after the existing imports (after the `AttachmentChips` import):

```typescript
import { ErrorMessageBubble } from './ErrorMessageBubble'
```

Add the `"error"` branch. Insert it after the `contentType === 'thinking'` block (around line 91) and before the `displayContent` section:

```typescript
if (contentType === 'error') {
  return <ErrorMessageBubble content={content} />
}
```

The relevant section of MessageItem after the change:

```typescript
if (contentType === 'thinking') {
  return <ThinkingBlock content={content} defaultExpanded={false} active={false} />
}

if (contentType === 'error') {
  return <ErrorMessageBubble content={content} />
}

const displayContent = isStreaming && streamingContent !== undefined ? streamingContent : content
```

- [ ] **Step 4: Run tests**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/messageItem.test.tsx
```

Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add frontend/src/components/chat/MessageItem.tsx frontend/src/__tests__/messageItem.test.tsx && git commit -m "feat(chat): render error content type as ErrorMessageBubble in MessageItem"
```

---

## Task 7: Fix chatStore Error Handling

**Files:**
- Modify: `frontend/src/store/chatStore.ts`
- Modify: `frontend/src/__tests__/chatStore.test.ts`

- [ ] **Step 1: Write failing tests**

Open `frontend/src/__tests__/chatStore.test.ts`. Read the file to find the existing `beforeEach` reset block. Then add two new `describe` blocks at the end of the file:

```typescript
describe('stream error handling', () => {
  it('calls loadMessages after stream:error event', async () => {
    const off = useChatStore.getState().initEventListeners()

    useChatStore.setState({
      conversations: [{ id: 1, title: 'chat', createdAt: '', updatedAt: '', starred: false, status: 'loading' }],
      messages: { 1: [] },
      streamingMessages: { 1: { content: 'partial', thinking: '' } },
    })

    const handler = eventHandlers.get('agent:stream:error')
    expect(handler).toBeDefined()
    handler!({ data: { sessionId: 1, error: 'rate limit' } })

    await Promise.resolve()

    expect(AgentBinding.LoadSessionMessages).toHaveBeenCalledWith(
      expect.objectContaining({ session_id: 1 })
    )
    expect(useChatStore.getState().streamingMessages[1]).toBeUndefined()

    off()
  })
})

describe('sendMessage error handling', () => {
  it('replaces temp assistant message with error bubble when SendMessage throws', async () => {
    vi.mocked(AgentBinding.SendMessage).mockRejectedValueOnce(
      new Error(JSON.stringify({ msg: 'Failed to send message', detail: 'api key invalid' }))
    )

    useChatStore.setState({
      conversations: [{ id: 1, title: 'chat', createdAt: '', updatedAt: '', starred: false, status: 'idle' }],
      messages: { 1: [] },
      streamingMessages: {},
    })

    await useChatStore.getState().sendMessage({
      sessionId: 1,
      content: 'hello',
      baseUrl: 'http://localhost',
      apiKey: 'test',
      modelName: 'gpt-4',
      providerType: 'openai_compatibility' as ProviderType,
      enabledUserTools: [],
      attachments: [],
    })

    const messages = useChatStore.getState().messages[1] ?? []
    const errorMsg = messages.find((m) => m.contentType === 'error')
    expect(errorMsg).toBeDefined()
    const parsed = JSON.parse(errorMsg!.content) as { msg: string; detail: string }
    expect(parsed.msg).toBe('Failed to send message')
    expect(parsed.detail).toBe('api key invalid')

    const status = useChatStore.getState().conversations.find((c) => c.id === 1)?.status
    expect(status).toBe('idle')
  })
})
```

- [ ] **Step 2: Run to confirm tests fail**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/chatStore.test.ts
```

Expected: new tests FAIL (stream:error test fails because `LoadSessionMessages` is not called; sendMessage test fails because no error message is inserted)

- [ ] **Step 3: Add `parseIError` helper and update `chatStore.ts`**

Open `frontend/src/store/chatStore.ts`.

**3a.** Add `parseIError` function just before the `useChatStore` declaration (after the `cancelScheduledFrame` function):

```typescript
function parseIError(err: unknown): { msg: string; detail: string } {
  try {
    const message = err instanceof Error ? err.message : String(err)
    const parsed = JSON.parse(message) as { msg?: string; detail?: string }
    return {
      msg: parsed.msg ?? message,
      detail: parsed.detail ?? message,
    }
  } catch {
    return { msg: String(err), detail: String(err) }
  }
}
```

**3b.** In the `agent:stream:error` event handler, add `void get().loadMessages(data.sessionId)` after the `set(...)` call:

Change:

```typescript
offs.push(Events.On('agent:stream:error', (event: { data: StreamErrorEvent | StreamErrorEvent[] }) => {
  const data = eventData(event)
  if (!data) return

  flushStreamingBuffer(data.sessionId)

  set((state) => {
    const nextStreaming = { ...state.streamingMessages }
    delete nextStreaming[data.sessionId]
    return { streamingMessages: nextStreaming }
  })
}))
```

To:

```typescript
offs.push(Events.On('agent:stream:error', (event: { data: StreamErrorEvent | StreamErrorEvent[] }) => {
  const data = eventData(event)
  if (!data) return

  flushStreamingBuffer(data.sessionId)

  set((state) => {
    const nextStreaming = { ...state.streamingMessages }
    delete nextStreaming[data.sessionId]
    return { streamingMessages: nextStreaming }
  })

  void get().loadMessages(data.sessionId)
}))
```

**3c.** In `sendMessage`, wrap `await AgentBinding.SendMessage(...)` in try/catch. Replace:

```typescript
await AgentBinding.SendMessage({
  session_id: params.sessionId,
  content: params.content,
  base_url: params.baseUrl,
  api_key: params.apiKey,
  model_name: params.modelName,
  provider_type: params.providerType,
  enabled_user_tools: params.enabledUserTools,
  attachments: params.attachments.map(att => ({
    path: att.path,
    name: att.name,
    mime: att.mime,
  })),
})
```

With:

```typescript
try {
  await AgentBinding.SendMessage({
    session_id: params.sessionId,
    content: params.content,
    base_url: params.baseUrl,
    api_key: params.apiKey,
    model_name: params.modelName,
    provider_type: params.providerType,
    enabled_user_tools: params.enabledUserTools,
    attachments: params.attachments.map(att => ({
      path: att.path,
      name: att.name,
      mime: att.mime,
    })),
  })
} catch (err: unknown) {
  const parsed = parseIError(err)
  set((state) => {
    const nextStreaming = { ...state.streamingMessages }
    delete nextStreaming[params.sessionId]
    return {
      messages: {
        ...state.messages,
        [params.sessionId]: (state.messages[params.sessionId] ?? []).map((m) =>
          m.id === tempAssistantId
            ? { ...m, contentType: 'error' as const, content: JSON.stringify(parsed) }
            : m
        ),
      },
      conversations: state.conversations.map((c) =>
        c.id === params.sessionId ? { ...c, status: 'idle' } : c
      ),
      streamingMessages: nextStreaming,
      pendingTitleGenerations: {
        ...state.pendingTitleGenerations,
        [params.sessionId]: undefined,
      },
    }
  })
}
```

- [ ] **Step 4: Run tests**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run src/__tests__/chatStore.test.ts
```

Expected: all tests PASS

- [ ] **Step 5: Run full frontend test suite**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx vitest run
```

Expected: all tests PASS

- [ ] **Step 6: TypeScript check**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx tsc --noEmit
```

Expected: no errors

- [ ] **Step 7: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add frontend/src/store/chatStore.ts frontend/src/__tests__/chatStore.test.ts && git commit -m "fix(chat): reload messages after stream error and catch sendMessage failures"
```

---

## Self-Review

**Spec coverage check:**

| Spec requirement | Covered by |
|---|---|
| `IError.Err error` → `Detail string` | Task 1 |
| `ErrAgentStreamError` constant | Task 2 |
| Stream error i18n keys (en + zh-CN) | Task 2 |
| `chat_handler.go` stores `ContentType: "error"` with JSON | Task 3 |
| Frontend `"error"` added to `MessageContentType` | Task 4 |
| Frontend i18n keys `errorDetail` / `hideErrorDetail` | Task 4 |
| `ErrorMessageBubble` component (msg + expand icon + detail) | Task 5 |
| `MessageItem` handles `contentType === "error"` | Task 6 |
| `chatStore` stream:error handler calls `loadMessages` | Task 7 |
| `chatStore` sendMessage try/catch with error bubble insertion | Task 7 |
| `parseIError` helper | Task 7 |

All spec requirements are covered. No gaps.

**Type consistency check:**
- `MessageContentType` gains `'error'` in Task 4; used in Task 5 (`ErrorMessageBubble`), Task 6 (`MessageItem`), Task 7 (`chatStore` catch block with `'error' as const`)
- `IError` fields `Detail` + `Msg` defined in Task 1; consumed by Task 3 (`chat_handler.go`) and Task 7 (`parseIError` reads `msg`/`detail` keys)
- `parseIError` defined in Task 7 step 3a; used in Task 7 step 3c — same file, consistent

**Placeholder scan:** No TBDs, TODOs, or vague steps. All code blocks are complete.
