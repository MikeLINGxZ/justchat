import { describe, it, expect } from 'vitest'
import { shouldShowTimestamp } from '../lib/utils'
import type { Message } from '../types'

const makeMsg = (daysAgo: number, id: number): Message => ({
  id,
  sessionId: 1,
  parentId: null,
  role: 'user',
  contentType: 'text',
  content: 'test',
  modelName: '',
  agentName: '',
  tokensIn: 0,
  tokensOut: 0,
  extra: '',
  createdAt: new Date(Date.now() - daysAgo * 86400000).toISOString(),
})

describe('shouldShowTimestamp', () => {
  it('always shows timestamp for first message', () => {
    const msgs = [makeMsg(0, 1)]
    expect(shouldShowTimestamp(msgs, 0)).toBe(true)
  })

  it('does not show timestamp for same-day messages', () => {
    const msgs = [makeMsg(0, 1), makeMsg(0, 2)]
    expect(shouldShowTimestamp(msgs, 1)).toBe(false)
  })

  it('shows timestamp when messages cross day boundary', () => {
    const msgs = [makeMsg(1, 1), makeMsg(0, 2)]
    expect(shouldShowTimestamp(msgs, 1)).toBe(true)
  })
})
