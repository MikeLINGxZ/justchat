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
      'tool_call', // grouped (file_read + code_exec)
      'text',      // assistant reply
    ])
    expect(result[2].isToolGroup).toBe(true)
    expect(result[2].groupedTools).toHaveLength(2)
    expect(result[2].groupedTools![0].toolResult).toBe('result1')
    expect(result[2].groupedTools![1].toolResult).toBe('result2')
  })

  it('hides confirm_response from display output', () => {
    const messages = [
      msg(1, 'user', 'text', 'read file'),
      msg(2, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(3, 'user', 'confirm_response', 'approved', '{"action":"approve","comment":""}'),
      msg(4, 'tool', 'tool_result', 'contents', '', 'file_read'),
      msg(5, 'assistant', 'text', 'here it is'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result.map(m => m.contentType)).toEqual([
      'text',      // user
      'tool_call', // with result merged
      'text',      // assistant reply
    ])
    expect(result[1].toolResult).toBe('contents')
    expect(result[1].toolConfirmAction).toBe('approve')
    expect(result.find(m => m.contentType === 'confirm_response')).toBeUndefined()
  })

  it('attaches structured reject action to the related tool call', () => {
    const messages = [
      msg(1, 'user', 'text', 'read file'),
      msg(2, 'assistant', 'tool_call', '{"name":"file_read","args":{}}'),
      msg(3, 'user', 'confirm_response', 'rejected', '{"action":"reject","comment":""}'),
      msg(4, 'tool', 'tool_result', 'Tool confirmation: rejected.\n\n{"error":"tool execution rejected by user"}', '', 'file_read'),
    ]
    const result = groupMessagesForDisplay(messages)
    expect(result).toHaveLength(2)
    expect(result[1].contentType).toBe('tool_call')
    expect(result[1].toolConfirmAction).toBe('reject')
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
