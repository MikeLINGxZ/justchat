# Error Message Display Design

**Date:** 2026-05-17  
**Status:** Approved  
**Scope:** Backend IError fix + frontend error bubble + stream error blank fix

---

## Summary

Two problems to solve:

1. **Blank assistant bubble after LLM stream error** — backend saves error to DB but frontend never reloads messages after `agent:stream:error`.
2. **No error display** — neither streaming LLM errors nor Wails call errors (e.g., attachment validation) are shown to users.

---

## Architecture

### Error Data Structure

**Backend `IError`** (after fix, in `backend/pkg/ierror/ierrors.go`):

```go
type IError struct {
    Detail string `json:"detail"` // raw err.Error() string
    Msg    string `json:"msg"`    // localized user-facing message
}
```

Old `Err error` field replaced by `Detail string` — `error` interface does not serialize its message via `encoding/json`. `Unwrap()` method is removed (nothing to unwrap with a string detail).

**Error message content stored in DB** (new `"error"` content type):

```json
{ "msg": "大模型响应出错", "detail": "<raw LLM provider error string>" }
```

Stored as `ContentType: "error"` (previously `"text"`) so frontend can distinguish error messages.

---

## Backend Changes

### 1. `backend/pkg/ierror/ierrors.go`

Change `IError` struct:

```go
type IError struct {
    Detail string `json:"detail"`
    Msg    string `json:"msg"`
}
```

Remove `Unwrap()` method (no longer wraps an error).

Update `Error()` constructor:

```go
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

Remove `"errors"` import if `Unwrap()` is removed and `errors.As` is still used (keep it).

### 2. `backend/pkg/ierror/error.go`

Add new error code constant in the Agent section:

```go
ErrAgentStreamError errorCode = "ierror.agent.stream_error"
```

### 3. `backend/pkg/i18n/resources_en_us.go`

```go
"ierror.agent.stream_error": "Model response error",
```

### 4. `backend/pkg/i18n/resources_en_zh_cn.go`

```go
"ierror.agent.stream_error": "大模型响应出错",
```

### 5. `backend/pkg/agent/chat_handler.go`

Add `ierror` import:

```go
"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
```

Change stream error message storage (in the `evt.Error != nil` block):

```go
// Build structured error content
errContent, _ := json.Marshal(map[string]string{
    "msg":    ierror.ErrAgentStreamError.Msg(),
    "detail": evt.Error.Message,
})
_, _ = stor.CreateMessage(data_models.Message{
    SessionID:   sessionID,
    Role:        "assistant",
    ContentType: "error",        // was "text"
    Content:     string(errContent),
    ModelName:   modelName,
})
```

Remove the `fmt.Sprintf("Error: %s", ...)` call. `json` is already imported.

---

## Frontend Changes

### 1. `frontend/src/types/index.ts`

```typescript
export type MessageContentType =
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'thinking'
  | 'confirm_request'
  | 'confirm_response'
  | 'error'                // NEW
```

### 2. `frontend/src/store/chatStore.ts`

**Helper function** (add near top, before `useChatStore`):

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

**Stream error handler** — add `loadMessages` call:

```typescript
offs.push(Events.On('agent:stream:error', (event) => {
    const data = eventData(event)
    if (!data) return

    flushStreamingBuffer(data.sessionId)
    set((state) => {
        const nextStreaming = { ...state.streamingMessages }
        delete nextStreaming[data.sessionId]
        return { streamingMessages: nextStreaming }
    })
    void get().loadMessages(data.sessionId)  // NEW — fixes blank bubble
}))
```

**sendMessage** — wrap `AgentBinding.SendMessage` call in try/catch. On error, replace the temp assistant message (which has `id === tempAssistantId`) with an error message, and reset status to `'idle'`:

```typescript
try {
    await AgentBinding.SendMessage({ ... })
} catch (err: unknown) {
    const parsed = parseIError(err)
    set((state) => ({
        messages: {
            ...state.messages,
            [params.sessionId]: (state.messages[params.sessionId] ?? []).map((m) =>
                m.id === tempAssistantId
                    ? {
                          ...m,
                          contentType: 'error' as const,
                          content: JSON.stringify(parsed),
                      }
                    : m
            ),
        },
        conversations: state.conversations.map((c) =>
            c.id === params.sessionId ? { ...c, status: 'idle' } : c
        ),
        streamingMessages: (() => {
            const next = { ...state.streamingMessages }
            delete next[params.sessionId]
            return next
        })(),
    }))
}
```

### 3. New `frontend/src/components/chat/ErrorMessageBubble.tsx`

```typescript
import { useState } from 'react'
import { AlertCircle, ChevronDown, ChevronUp } from 'lucide-react'
import { useTranslation } from 'react-i18next'

interface Props {
    content: string  // JSON string: { msg: string, detail: string }
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
    } catch { /* fallback: display raw string */ }

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

### 4. `frontend/src/components/chat/MessageItem.tsx`

Add import at top:

```typescript
import { ErrorMessageBubble } from './ErrorMessageBubble'
```

Add branch before `displayContent` block (after the `contentType === 'thinking'` check):

```typescript
if (contentType === 'error') {
    return <ErrorMessageBubble content={content} />
}
```

### 5. `frontend/src/i18n/locales/en.ts`

Add to `chat` section:

```typescript
errorDetail: 'Show error details',
hideErrorDetail: 'Hide error details',
```

### 6. `frontend/src/i18n/locales/zh-CN.ts`

Add to `chat` section:

```typescript
errorDetail: '展开错误详情',
hideErrorDetail: '收起错误详情',
```

---

## Data Flow

```
LLM Provider Error
    → chat_handler emits agent:stream:error
    → chat_handler saves { role: assistant, contentType: "error", content: JSON } to DB
    → frontend stream:error handler:
        1. clears streaming state
        2. calls loadMessages() → fetches error message from DB
        3. MessageItem reads contentType === "error" → renders ErrorMessageBubble

Wails call error (e.g. attachment validation)
    → AgentBinding.SendMessage() throws
    → chatStore.sendMessage catch block
    → replaces tempAssistantId message with { contentType: "error", content: JSON }
    → MessageItem reads contentType === "error" → renders ErrorMessageBubble
```

---

## Out of Scope

- Other Wails binding errors (delete session, rename, etc.) — not handled this iteration
- Toast notification component — not needed; all errors use in-chat bubble
- Migrating existing "Error: ..." text messages in DB — legacy records remain as plain text
