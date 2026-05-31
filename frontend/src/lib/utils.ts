import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { Message, Conversation, DisplayMessage } from '../types'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const isMac =
  typeof navigator !== 'undefined' &&
  (/Mac/i.test(navigator.platform) || /Mac/i.test(navigator.userAgent))

function toDate(value: string | Date): Date {
  return value instanceof Date ? value : new Date(value)
}

export function shouldShowTimestamp(messages: Message[], index: number): boolean {
  if (index === 0) return true
  const prev = toDate(messages[index - 1].createdAt)
  const curr = toDate(messages[index].createdAt)
  const prevDay = new Date(prev.getFullYear(), prev.getMonth(), prev.getDate())
  const currDay = new Date(curr.getFullYear(), curr.getMonth(), curr.getDate())
  return prevDay.getTime() !== currDay.getTime()
}

export function formatTimestamp(dateValue: string | Date, language: string): string {
  const date = toDate(dateValue)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  const msgDay = new Date(date.getFullYear(), date.getMonth(), date.getDate())

  const timeStr = date.toLocaleTimeString(language === 'zh-CN' ? 'zh-CN' : 'en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })

  const isZh = language === 'zh-CN'

  if (msgDay.getTime() === today.getTime()) {
    return isZh ? `今天 ${timeStr}` : `Today ${timeStr}`
  }
  if (msgDay.getTime() === yesterday.getTime()) {
    return isZh ? `昨天 ${timeStr}` : `Yesterday ${timeStr}`
  }
  if (date.getFullYear() === now.getFullYear()) {
    return isZh
      ? `${date.getMonth() + 1}月${date.getDate()}日 ${timeStr}`
      : `${date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} ${timeStr}`
  }
  return isZh
    ? `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${timeStr}`
    : `${date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })} ${timeStr}`
}

const CONTENT_TYPE_PRIORITY: Record<string, number> = {
  thinking: 0,
  tool_call: 1,
  text: 2,
}

export function groupMessagesForDisplay(messages: Message[]): DisplayMessage[] {
  const toolResultMap = new Map<number, string>()
  const toolConfirmMap = new Map<number, { action: 'approve' | 'reject' | 'comment'; comment?: string }>()

  for (let i = 0; i < messages.length; i++) {
    const m = messages[i]
    if (m.contentType === 'tool_result') {
      for (let j = i - 1; j >= 0; j--) {
        if (messages[j].contentType === 'tool_call') {
          toolResultMap.set(messages[j].id, m.content)
          break
        }
      }
      continue
    }

    if (m.contentType === 'confirm_response') {
      try {
        const parsed = JSON.parse(m.extra || '{}') as { action?: 'approve' | 'reject' | 'comment'; comment?: string }
        if (parsed.action === 'approve' || parsed.action === 'reject' || parsed.action === 'comment') {
          for (let j = i - 1; j >= 0; j--) {
            if (messages[j].contentType === 'tool_call') {
              toolConfirmMap.set(messages[j].id, {
                action: parsed.action,
                comment: parsed.comment,
              })
              break
            }
          }
        }
      } catch {
        // Ignore malformed metadata.
      }
    }
  }

  const filtered = messages.filter(
    m => m.contentType !== 'tool_result' && m.contentType !== 'confirm_response'
  )

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
          const confirmMeta = toolConfirmMap.get(m.id)
          if (confirmMeta) {
            display.toolConfirmAction = confirmMeta.action
            display.toolConfirmComment = confirmMeta.comment
          }
        }
        result.push(display)
      }
    }

    if (!isEnd && isUserText) {
      result.push({ ...filtered[i] })
      roundStart = i + 1
    }
  }

  const grouped: DisplayMessage[] = []
  let toolGroup: DisplayMessage[] = []

  const flushToolGroup = () => {
    if (toolGroup.length === 0) return
    if (toolGroup.length === 1) {
      grouped.push(toolGroup[0])
    } else {
      const first = toolGroup[0]
      grouped.push({
        ...first,
        id: first.id,
        isToolGroup: true,
        groupedTools: toolGroup,
      })
    }
    toolGroup = []
  }

  for (const m of result) {
    if (m.contentType === 'tool_call') {
      toolGroup.push(m)
    } else {
      flushToolGroup()
      grouped.push(m)
    }
  }
  flushToolGroup()

  return grouped
}

export function groupConversationsByDate(
  conversations: Conversation[],
  language: string
): { label: string; items: Conversation[] }[] {
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  const sevenDaysAgo = new Date(today.getTime() - 7 * 86400000)

  const groups: Record<string, Conversation[]> = {
    today: [],
    yesterday: [],
    week: [],
    earlier: [],
  }

  for (const conv of conversations) {
    const updatedAt = toDate(conv.updatedAt)
    const d = new Date(
      updatedAt.getFullYear(),
      updatedAt.getMonth(),
      updatedAt.getDate()
    )
    if (d >= today) groups.today.push(conv)
    else if (d >= yesterday) groups.yesterday.push(conv)
    else if (d >= sevenDaysAgo) groups.week.push(conv)
    else groups.earlier.push(conv)
  }

  const labels = language === 'zh-CN'
    ? { today: '今天', yesterday: '昨天', week: '过去7天', earlier: '更久以前' }
    : { today: 'Today', yesterday: 'Yesterday', week: 'Last 7 days', earlier: 'Earlier' }

  return Object.entries(groups)
    .filter(([, items]) => items.length > 0)
    .map(([key, items]) => ({ label: labels[key as keyof typeof labels], items }))
}
