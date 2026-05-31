# Main Agent UX Improvements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix four UX issues in the main agent chat: delay conversation creation until first message, fix tool confirmation rendering, reorder messages for visual grouping (tool_result merged into tool_call), and remove scroll animation on conversation switch.

**Architecture:** All changes are frontend-only (TypeScript/React). No backend modifications. Changes touch the chat store, message display components, and auto-scroll hook. A new `groupMessagesForDisplay` utility handles message reordering and tool_result merging.

**Tech Stack:** React 18, Zustand, TypeScript, Vitest, Tailwind CSS

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `frontend/src/components/sidebar/ConversationList.tsx` | Modify | Remove `createConversation` call from "New Chat" handler |
| `frontend/src/components/chat/ChatMessages.tsx` | Modify | Expand `isStreaming` to include `waiting-unread`; remove `scroll-smooth` class; apply message grouping |
| `frontend/src/components/chat/ChatInput.tsx` | Modify | Expand `isStreaming` to include `waiting-unread` |
| `frontend/src/components/chat/MessageItem.tsx` | Modify | Accept `toolResult` prop for tool_call messages; remove standalone `tool_result` branch |
| `frontend/src/lib/utils.ts` | Modify | Add `groupMessagesForDisplay` function |
| `frontend/src/types/index.ts` | Modify | Add `DisplayMessage` type |
| `frontend/src/__tests__/groupMessages.test.ts` | Create | Tests for `groupMessagesForDisplay` |

---

### Task 1: Delay Conversation Creation Until First Message

**Files:**
- Modify: `frontend/src/components/sidebar/ConversationList.tsx:33-40`

- [ ] **Step 1: Replace `handleNewChat` to skip backend call**

In `frontend/src/components/sidebar/ConversationList.tsx`, replace the `handleNewChat` function:

```typescript
// OLD (lines 33-40):
  const handleNewChat = async () => {
    const id = await createConversation()
    if (id) {
      startTransition(() => {
        setCurrentConversation(id)
      })
    }
  }

// NEW:
  const handleNewChat = () => {
    startTransition(() => {
      setCurrentConversation(null)
    })
  }
```

Also remove `createConversation` from the `useChatStore` destructuring on line 15-21 since it is no longer used here:

```typescript
// OLD:
  const { conversations, currentConversationId, setCurrentConversation, createConversation } =
    useChatStore(useShallow((state) => ({
      conversations: state.conversations,
      currentConversationId: state.currentConversationId,
      setCurrentConversation: state.setCurrentConversation,
      createConversation: state.createConversation,
    })))

// NEW:
  const { conversations, currentConversationId, setCurrentConversation } =
    useChatStore(useShallow((state) => ({
      conversations: state.conversations,
      currentConversationId: state.currentConversationId,
      setCurrentConversation: state.setCurrentConversation,
    })))
```

Remove the `startTransition` import if it is no longer used elsewhere — but it IS still used in the `onClick` of `ConversationItem` (line 99), so keep the import.

- [ ] **Step 2: Verify build passes**

Run:
```bash
cd frontend && npx tsc --noEmit
```
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/sidebar/ConversationList.tsx
git commit -m "feat(chat): delay conversation creation until first message sent"
```

---

### Task 2: Fix Tool Confirmation — Expand `isStreaming` Check

**Files:**
- Modify: `frontend/src/components/chat/ChatMessages.tsx:100-118`
- Modify: `frontend/src/components/chat/ChatInput.tsx:120-133`

- [ ] **Step 1: Fix `isStreaming` in `ChatMessages.tsx`**

In `frontend/src/components/chat/ChatMessages.tsx`, change the `isStreaming` calculation inside `useChatStore(useShallow(...))` (around line 113):

```typescript
// OLD:
      isStreaming: currentConversation?.status === 'loading',

// NEW:
      isStreaming: currentConversation?.status === 'loading' || currentConversation?.status === 'waiting-unread',
```

- [ ] **Step 2: Fix `isStreaming` in `ChatInput.tsx`**

In `frontend/src/components/chat/ChatInput.tsx`, change the `isStreaming` calculation inside `useChatStore(useShallow(...))` (around line 127):

```typescript
// OLD:
      isStreaming: currentConversation?.status === 'loading',

// NEW:
      isStreaming: currentConversation?.status === 'loading' || currentConversation?.status === 'waiting-unread',
```

- [ ] **Step 3: Verify build passes**

Run:
```bash
cd frontend && npx tsc --noEmit
```
Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/chat/ChatMessages.tsx frontend/src/components/chat/ChatInput.tsx
git commit -m "fix(chat): show tool confirmation card when session status is waiting-unread"
```

---

### Task 3: Add `DisplayMessage` Type and `groupMessagesForDisplay` Utility

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/lib/utils.ts`
- Create: `frontend/src/__tests__/groupMessages.test.ts`

- [ ] **Step 1: Add `DisplayMessage` type**

In `frontend/src/types/index.ts`, add after the `Message` type (after line 44):

```typescript
export type DisplayMessage = Message & {
  toolResult?: string
}
```

- [ ] **Step 2: Write the failing test**

Create `frontend/src/__tests__/groupMessages.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { groupMessagesForDisplay } from '../lib/utils'
import type { Message } from '../types'

const msg = (
  id: number,
  role: Message['role'],
  contentType: Message['contentType'],
  content = '',
  extra = '',
  agentName = '',
): Message => ({
  id,
  sessionId: 1,
  parentId: null,
  role,
  contentType,
  content,
  modelName: '',
  agentName,
  tokensIn: 0,
  tokensOut: 0,
  extra,
  createdAt: new Date().toISOString(),
})

describe('groupMessagesForDisplay', () => {
  it('returns empty array for empty input', () => {
    expect(groupMessagesForDisplay([])).toEqual([])
  })

  it('passes through user messages unchanged', () => {
    const messages = [msg(1, 'user', 'text', 'hello')]
    const result = groupMessagesForDisplay(messages)
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(1)
  })

  it('sorts assistant round: thinking before tool_call before text', () => {
    const messages = [
      msg(1, 'user', 'text', 'hi'),
      msg(2, 'assistant', 'text', 'reply'),
      msg(3, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(4, 'assistant', 'thinking', 'let me think'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result.map(m => m.contentType)).toEqual([
      'text',      // user
      'thinking',  // sorted first in assistant round
      'tool_call', // sorted second
      'text',      // sorted last
    ])
  })

  it('merges tool_result into preceding tool_call', () => {
    const messages = [
      msg(1, 'user', 'text', 'read a file'),
      msg(2, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(3, 'tool', 'tool_result', 'file contents here', '', 'file_read'),
      msg(4, 'assistant', 'text', 'here is the content'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result).toHaveLength(3) // user, tool_call(+result), text
    expect(result[1].contentType).toBe('tool_call')
    expect(result[1].toolResult).toBe('file contents here')
    expect(result.find(m => m.contentType === 'tool_result')).toBeUndefined()
  })

  it('handles multiple tool calls with results', () => {
    const messages = [
      msg(1, 'user', 'text', 'do stuff'),
      msg(2, 'assistant', 'thinking', 'planning'),
      msg(3, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(4, 'tool', 'tool_result', 'result1', '', 'file_read'),
      msg(5, 'assistant', 'tool_call', '{"name":"code_exec","args":{}}'),
      msg(6, 'tool', 'tool_result', 'result2', '', 'code_exec'),
      msg(7, 'assistant', 'text', 'done'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result.map(m => m.contentType)).toEqual([
      'text',      // user
      'thinking',
      'tool_call', // file_read with result1
      'tool_call', // code_exec with result2
      'text',      // assistant reply
    ])
    expect(result[2].toolResult).toBe('result1')
    expect(result[3].toolResult).toBe('result2')
  })

  it('handles confirm_response in correct position', () => {
    const messages = [
      msg(1, 'user', 'text', 'read file'),
      msg(2, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(3, 'user', 'confirm_response', 'approved'),
      msg(4, 'tool', 'tool_result', 'contents', '', 'file_read'),
      msg(5, 'assistant', 'text', 'here it is'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result.map(m => m.contentType)).toEqual([
      'text',             // user
      'tool_call',        // with result merged
      'confirm_response', // kept after tool
      'text',             // assistant reply
    ])
    expect(result[1].toolResult).toBe('contents')
  })

  it('preserves messages across multiple user rounds', () => {
    const messages = [
      msg(1, 'user', 'text', 'first'),
      msg(2, 'assistant', 'text', 'reply1'),
      msg(3, 'user', 'text', 'second'),
      msg(4, 'assistant', 'thinking', 'hmm'),
      msg(5, 'assistant', 'text', 'reply2'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result.map(m => m.contentType)).toEqual([
      'text',     // user 1
      'text',     // assistant reply1
      'text',     // user 2
      'thinking', // sorted first in round 2
      'text',     // assistant reply2
    ])
  })
})
```

- [ ] **Step 3: Run test to verify it fails**

Run:
```bash
cd frontend && npx vitest run src/__tests__/groupMessages.test.ts
```
Expected: FAIL — `groupMessagesForDisplay` is not exported from `../lib/utils`.

- [ ] **Step 4: Implement `groupMessagesForDisplay`**

In `frontend/src/lib/utils.ts`, add the import and function at the end of the file:

```typescript
import type { Message, Conversation, DisplayMessage } from '../types'
```

(Update the existing import on line 3 to add `DisplayMessage`.)

Then add the function:

```typescript
const CONTENT_TYPE_PRIORITY: Record<string, number> = {
  thinking: 0,
  tool_call: 1,
  confirm_response: 2,
  text: 3,
}

export function groupMessagesForDisplay(messages: Message[]): DisplayMessage[] {
  const toolResultMap = new Map<number, string>()

  for (let i = 0; i < messages.length; i++) {
    const m = messages[i]
    if (m.contentType !== 'tool_result') continue

    for (let j = i - 1; j >= 0; j--) {
      if (messages[j].contentType === 'tool_call') {
        toolResultMap.set(messages[j].id, m.content)
        break
      }
    }
  }

  const filtered = messages.filter(m => m.contentType !== 'tool_result')

  const result: DisplayMessage[] = []
  let roundStart = 0

  for (let i = 0; i <= filtered.length; i++) {
    const isEnd = i === filtered.length
    const isUserText = !isEnd && filtered[i].role === 'user' && filtered[i].contentType === 'text'

    if ((isEnd || isUserText) && i > roundStart) {
      const round = filtered.slice(roundStart, i)
      const sorted = [...round].sort((a, b) => {
        const pa = CONTENT_TYPE_PRIORITY[a.contentType] ?? 3
        const pb = CONTENT_TYPE_PRIORITY[b.contentType] ?? 3
        return pa - pb
      })

      for (const m of sorted) {
        const display: DisplayMessage = { ...m }
        if (m.contentType === 'tool_call') {
          display.toolResult = toolResultMap.get(m.id)
        }
        result.push(display)
      }
    }

    if (!isEnd && isUserText) {
      result.push({ ...filtered[i] })
      roundStart = i + 1
    }
  }

  return result
}
```

- [ ] **Step 5: Run test to verify it passes**

Run:
```bash
cd frontend && npx vitest run src/__tests__/groupMessages.test.ts
```
Expected: All 7 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/types/index.ts frontend/src/lib/utils.ts frontend/src/__tests__/groupMessages.test.ts
git commit -m "feat(chat): add groupMessagesForDisplay utility for visual message ordering"
```

---

### Task 4: Wire Up Message Grouping in Chat Display

**Files:**
- Modify: `frontend/src/components/chat/MessageItem.tsx`
- Modify: `frontend/src/components/chat/ChatMessages.tsx`

- [ ] **Step 1: Update `MessageItem` to accept `toolResult` and remove `tool_result` branch**

In `frontend/src/components/chat/MessageItem.tsx`:

First update the import:

```typescript
// OLD (line 10):
import type { Message } from '@/types'

// NEW:
import type { DisplayMessage } from '@/types'
```

Update the `Props` interface:

```typescript
// OLD:
interface Props {
  message: Message
  // ... rest unchanged
}

// NEW:
interface Props {
  message: DisplayMessage
  // ... rest unchanged
}
```

In the `tool_call` handler (around lines 62-82), pass the `toolResult` from the message to `ToolCallBlock`:

```typescript
// OLD (lines 62-82):
  if (contentType === 'tool_call') {
    let toolName = ''
    let args = content
    try {
      const parsed = JSON.parse(content)
      toolName = parsed.name ?? ''
      args = JSON.stringify(parsed.args ?? parsed, null, 2)
    } catch {
      // Fall back to raw content.
    }

    return (
      <ToolCallBlock
        toolName={toolName}
        purpose={extra}
        args={args}
        result=""
        status="completed"
      />
    )
  }

// NEW:
  if (contentType === 'tool_call') {
    let toolName = ''
    let args = content
    try {
      const parsed = JSON.parse(content)
      toolName = parsed.name ?? ''
      args = JSON.stringify(parsed.args ?? parsed, null, 2)
    } catch {
      // Fall back to raw content.
    }

    return (
      <ToolCallBlock
        toolName={toolName}
        purpose={extra}
        args={args}
        result={message.toolResult ?? ''}
        status="completed"
      />
    )
  }
```

Remove the standalone `tool_result` branch entirely (lines 84-93):

```typescript
// DELETE these lines:
  if (contentType === 'tool_result') {
    return (
      <ToolCallBlock
        toolName={message.agentName}
        purpose={extra}
        args=""
        result={content}
        status="completed"
      />
    )
  }
```

- [ ] **Step 2: Apply `groupMessagesForDisplay` in `ChatMessages.tsx`**

In `frontend/src/components/chat/ChatMessages.tsx`, add the import:

```typescript
// Add to imports:
import { groupMessagesForDisplay } from '@/lib/utils'
```

Also update the `Message` type import:

```typescript
// OLD:
import type { Message } from '@/types'

// NEW:
import type { DisplayMessage } from '@/types'
```

Update `HistoricalMessages` to accept `DisplayMessage[]`:

```typescript
// OLD (lines 19-20):
interface MessageGroupProps {
  language: string
  messages: Message[]
}

// NEW:
interface MessageGroupProps {
  language: string
  messages: DisplayMessage[]
}
```

In the `ChatMessages` function body, apply grouping to `historicalMessages`. Change the `useMemo` that computes `historicalMessages` (around lines 134-137):

```typescript
// OLD:
  const historicalMessages = useMemo(
    () => (streamingMessage ? currentMessages.slice(0, -1) : currentMessages),
    [currentMessages, streamingMessage]
  )

// NEW:
  const historicalMessages = useMemo(
    () => groupMessagesForDisplay(streamingMessage ? currentMessages.slice(0, -1) : currentMessages),
    [currentMessages, streamingMessage]
  )
```

Also need to keep `shouldShowTimestamp` and `formatTimestamp` import working with `Message` for timestamps — these still work because `DisplayMessage` extends `Message`.

- [ ] **Step 3: Verify build passes**

Run:
```bash
cd frontend && npx tsc --noEmit
```
Expected: No errors.

- [ ] **Step 4: Run all tests**

Run:
```bash
cd frontend && npx vitest run
```
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/chat/MessageItem.tsx frontend/src/components/chat/ChatMessages.tsx
git commit -m "feat(chat): wire up message grouping with tool_result merged into tool_call"
```

---

### Task 5: Remove `scroll-smooth` CSS Class

**Files:**
- Modify: `frontend/src/components/chat/ChatMessages.tsx:160`

- [ ] **Step 1: Remove `scroll-smooth` from container className**

In `frontend/src/components/chat/ChatMessages.tsx`, update the scroll container div:

```typescript
// OLD (around line 160):
        className="chat-scroll-area flex-1 overflow-y-auto py-4 scroll-smooth"

// NEW:
        className="chat-scroll-area flex-1 overflow-y-auto py-4"
```

- [ ] **Step 2: Verify build passes**

Run:
```bash
cd frontend && npx tsc --noEmit
```
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/chat/ChatMessages.tsx
git commit -m "fix(chat): remove scroll-smooth to eliminate animation on conversation switch"
```

---

### Task 6: Final Verification

- [ ] **Step 1: Run all tests**

Run:
```bash
cd frontend && npx vitest run
```
Expected: All tests pass.

- [ ] **Step 2: Build check**

Run:
```bash
cd frontend && npx tsc --noEmit
```
Expected: No errors.

- [ ] **Step 3: Manual testing checklist**

Start the app and verify:

1. Click "New Chat" -> sidebar has no new item -> send message -> conversation appears with title
2. Click "New Chat" multiple times -> no empty conversations
3. Trigger tool call -> confirmation card appears -> approve/reject works
4. Historical messages: thinking first, then tool calls (collapsed, result inside), then text
5. Switch conversations -> instant jump to bottom, no animation
6. Click "scroll to bottom" button -> smooth animation preserved
