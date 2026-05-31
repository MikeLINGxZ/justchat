# Main Agent UX Improvements Design

## Overview

Four targeted optimizations to improve the main agent chat experience. All changes follow existing code patterns and conventions.

---

## 1. Delay Conversation Creation Until First Message

### Problem

Clicking "New Chat" immediately calls `CreateSession` on the backend and adds an empty entry to the sidebar. If the user doesn't send a message, these empty sessions accumulate.

### Solution: Frontend Draft Mode

When the user clicks "New Chat", set `currentConversationId = null` without calling the backend. The sidebar shows no new entry. Only when the first message is sent does `createConversation` + `sendMessage` execute.

### Files Changed

| File | Change |
|------|--------|
| `frontend/src/components/sidebar/ConversationList.tsx` | `handleNewChat` no longer calls `createConversation()`. Instead calls `setCurrentConversation(null)`. |
| `frontend/src/store/chatStore.ts` | No change needed. `createConversation` is already called from `ChatInput.handleSend` when `!sessionId`. |
| `frontend/src/components/chat/ChatInput.tsx` | Already handles `!sessionId` case. No change needed. |
| `frontend/src/components/chat/ChatMessages.tsx` | Already shows `WelcomeScreen` when `currentMessages.length === 0`. No change needed. |

### Edge Cases

- Repeated "New Chat" clicks: Each click resets to `null`, no side effects.
- Switch to existing conversation and back to new: Works normally.

---

## 2. Fix Tool Confirmation Flow

### Problem

Tool confirmation is fully implemented in backend code (`RequiresConfirm: true`, `confirm_request` events, `WaitForConfirm`) and frontend code (`ToolConfirmCard`, `pendingConfirms`), but the confirmation card never renders at runtime.

### Root Cause

When the backend emits `agent:stream:confirm_request`, it also changes the session status to `waiting-unread` via `updateSessionStatus`. The frontend determines `isStreaming` by checking `currentConversation?.status === 'loading'`. Once status becomes `waiting-unread`, `isStreaming` is `false`, `StreamingMessageRow` is not rendered, and the `pendingConfirm` prop is never passed to `MessageItem`.

### Solution

Expand the `isStreaming` condition to include `waiting-unread`:

```typescript
// In ChatMessages.tsx and ChatInput.tsx
isStreaming: currentConversation?.status === 'loading' || currentConversation?.status === 'waiting-unread'
```

Optionally introduce a helper for clarity:

```typescript
const isActiveSession = (status: ConversationStatus) =>
  status === 'loading' || status === 'waiting-unread'
```

### Files Changed

| File | Change |
|------|--------|
| `frontend/src/components/chat/ChatMessages.tsx` | Expand `isStreaming` check to include `waiting-unread`. |
| `frontend/src/components/chat/ChatInput.tsx` | Same expansion for the send/stop button state. |

### Verification

After fix: trigger a tool call (e.g., file_read) -> backend emits confirm_request -> status becomes `waiting-unread` -> `isStreaming` still true -> `StreamingMessageRow` renders -> `ToolConfirmCard` shows -> user approves/rejects -> status returns to `loading` -> flow continues.

---

## 3. Tool Call Visual Grouping and Sorting

### Problem

Messages are stored and rendered as flat independent records. Tool calls, tool results, thinking, and text content appear in storage order. The desired display order is: thinking -> tool calls (with results embedded) -> text reply.

### Solution: Frontend Display Grouping

1. Add a `groupMessagesForDisplay` function that processes the flat message list into display groups.
2. Messages are split into "rounds" separated by `role=user` messages.
3. Within each round, messages are sorted by `contentType` priority:
   - `thinking` = 0
   - `tool_call` = 1
   - `confirm_response` = 2
   - `text` = 3
4. `tool_result` messages are merged into their corresponding `tool_call` (matched by adjacent position or tool name), not rendered separately.

### Data Structure

```typescript
type DisplayMessage = Message & {
  toolResult?: string  // populated for tool_call messages, merged from the matching tool_result
}
```

### Files Changed

| File | Change |
|------|--------|
| `frontend/src/lib/utils.ts` | Add `groupMessagesForDisplay(messages: Message[]): DisplayMessage[]` function. |
| `frontend/src/components/chat/ChatMessages.tsx` | Apply `groupMessagesForDisplay` to `historicalMessages` before rendering. |
| `frontend/src/components/chat/MessageItem.tsx` | Accept `DisplayMessage` type. For `tool_call`, pass `toolResult` to `ToolCallBlock`. Remove standalone `tool_result` rendering branch. |
| `frontend/src/components/chat/ToolCallBlock.tsx` | Already accepts `result` prop. No change needed. |

### Sorting Logic

```
Input:  [user, thinking, tool_call, tool_result, tool_call, tool_result, text]
Output: [user, thinking, tool_call(+result), tool_call(+result), text]
```

Backend storage order is not changed.

---

## 4. Eliminate Scroll Animation on Conversation Switch

### Problem

Switching conversations shows a visible scroll-to-bottom animation despite `scrollToBottom(false)` being called. The CSS class `scroll-smooth` on the scroll container overrides the JS `behavior: 'auto'` parameter.

### Root Cause

`ChatMessages.tsx` line:
```tsx
<div ref={containerRef} className="... scroll-smooth">
```

CSS `scroll-behavior: smooth` applies globally to the container, overriding `scrollTo({ behavior: 'auto' })`.

### Solution

Remove the static `scroll-smooth` CSS class from the scroll container. Smooth scrolling is already correctly handled by the JS `scrollTo({ behavior: 'smooth' })` parameter when explicitly requested (e.g., user clicking "scroll to bottom" button).

### Files Changed

| File | Change |
|------|--------|
| `frontend/src/components/chat/ChatMessages.tsx` | Remove `scroll-smooth` from the container's className. |

### Behavior After Fix

| Scenario | Behavior |
|----------|----------|
| Switch conversation | Instant jump to bottom (no animation) |
| Streaming auto-scroll | Instant follow (no animation, high-frequency updates) |
| User clicks "Back to bottom" button | Smooth scroll (`scrollToBottom(true)`) |

---

## Testing Checklist

- [ ] Click "New Chat" -> sidebar has no new item -> type message and send -> conversation appears in sidebar with auto-generated title
- [ ] Click "New Chat" multiple times -> no empty conversations created
- [ ] Trigger tool call (file_read, code_exec) -> confirmation card appears -> approve -> tool executes -> result shows
- [ ] Trigger tool call -> reject with message -> tool skipped
- [ ] Historical messages display: thinking block first, then tool calls (collapsed, with results inside), then text reply
- [ ] Switch between conversations -> no scroll animation, instant bottom position
- [ ] During streaming -> auto-scroll follows content instantly
- [ ] Click "scroll to bottom" button -> smooth scroll animation
- [ ] Light/dark theme: all changes look correct in both themes
