import { vi } from 'vitest'
import { describe, it, expect, beforeEach } from 'vitest'
import { useChatStore } from '../store/chatStore'
import { Agent as AgentBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent'
import { Type as ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider/models'

const eventHandlers = new Map<string, (event: { data: unknown }) => void>()
const rafCallbacks = new Map<number, FrameRequestCallback>()
let rafId = 0

vi.mock('@wailsio/runtime', () => ({
  Events: {
    On: vi.fn((name: string, handler: (event: { data: unknown }) => void) => {
      eventHandlers.set(name, handler)
      return () => eventHandlers.delete(name)
    }),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent', () => ({
  Agent: {
    CreateSession: vi.fn().mockResolvedValue({ session_id: 1, title: 'New Chat' }),
    DeleteSession: vi.fn().mockResolvedValue(undefined),
    RenameSession: vi.fn().mockResolvedValue(undefined),
    ToggleStarSession: vi.fn().mockResolvedValue(undefined),
    ListSessions: vi.fn().mockResolvedValue({ sessions: [], has_more: false, next_cursor: 0 }),
    LoadSessionMessages: vi.fn().mockResolvedValue({ messages: [], total: 0, has_more: false }),
    MarkSessionRead: vi.fn().mockResolvedValue(undefined),
    SendMessage: vi.fn().mockResolvedValue(undefined),
    StopGeneration: vi.fn().mockResolvedValue(undefined),
    RespondToConfirm: vi.fn().mockResolvedValue(undefined),
    GenerateTitle: vi.fn().mockResolvedValue({ title: 'New Chat' }),
  },
}))

beforeEach(() => {
  eventHandlers.clear()
  rafCallbacks.clear()
  rafId = 0
  ;(globalThis as { __lemonteaChatListenersInitialized?: boolean }).__lemonteaChatListenersInitialized = false
  vi.stubGlobal('requestAnimationFrame', vi.fn((callback: FrameRequestCallback) => {
    rafId += 1
    rafCallbacks.set(rafId, callback)
    return rafId
  }))
  vi.stubGlobal('cancelAnimationFrame', vi.fn((id: number) => {
    rafCallbacks.delete(id)
  }))
  useChatStore.setState({
    conversations: [],
    messages: {},
    currentConversationId: null,
    streamingMessages: {},
    pendingTitleGenerations: {},
    sessionRunConfigs: {},
    pendingConfirms: {},
    sessionsLoading: false,
    hasMoreSessions: true,
    sessionsCursor: 0,
  })
})

describe('chatStore', () => {
  it('createConversation adds a conversation and returns its id', async () => {
    const id = await useChatStore.getState().createConversation()
    const convs = useChatStore.getState().conversations
    expect(convs).toHaveLength(1)
    expect(convs[0].id).toBe(id)
    expect(convs[0].title).toBe('New Chat')
    expect(convs[0].status).toBe('idle')
  })

  it('deleteConversation removes the conversation', async () => {
    const id = await useChatStore.getState().createConversation()
    await useChatStore.getState().deleteConversation(id)
    expect(useChatStore.getState().conversations).toHaveLength(0)
  })

  it('renameConversation updates title', async () => {
    const id = await useChatStore.getState().createConversation()
    await useChatStore.getState().renameConversation(id, 'My Chat')
    const conv = useChatStore.getState().conversations.find(c => c.id === id)
    expect(conv?.title).toBe('My Chat')
  })

  it('generateTitleDraft returns a suggested title without mutating the conversation title', async () => {
    const id = await useChatStore.getState().createConversation()
    const generateTitleMock = vi.mocked(AgentBinding.GenerateTitle)
    generateTitleMock.mockResolvedValueOnce({ title: 'Suggested Title' })

    await useChatStore.getState().sendMessage({
      sessionId: id,
      content: 'hello world',
      baseUrl: 'https://example.com',
      apiKey: 'test-key',
      modelName: 'gpt-test',
      providerType: ProviderType.OpenAiCompatibility,
      enabledUserTools: [],
      attachments: [],
    })

    const draftTitle = await useChatStore.getState().generateTitleDraft(id)

    expect(draftTitle).toBe('Suggested Title')
    expect(generateTitleMock).toHaveBeenLastCalledWith({
      session_id: id,
      base_url: 'https://example.com',
      api_key: 'test-key',
      model_name: 'gpt-test',
      provider_type: ProviderType.OpenAiCompatibility,
    })
    expect(useChatStore.getState().conversations.find(c => c.id === id)?.title).toBe('New Chat')
  })

  it('toggleStar updates starred field', async () => {
    const id = await useChatStore.getState().createConversation()
    await useChatStore.getState().toggleStar(id, true)
    expect(useChatStore.getState().conversations[0].starred).toBe(true)
    await useChatStore.getState().toggleStar(id, false)
    expect(useChatStore.getState().conversations[0].starred).toBe(false)
  })

  it('markAsRead resets unread status to idle and persists it', async () => {
    const id = await useChatStore.getState().createConversation()
    const markSessionReadMock = vi.mocked(AgentBinding.MarkSessionRead)

    useChatStore.setState(s => ({
      conversations: s.conversations.map(c => c.id === id ? { ...c, status: 'done-unread' as const } : c)
    }))
    markSessionReadMock.mockClear()

    await useChatStore.getState().markAsRead(id)

    const conv = useChatStore.getState().conversations.find(c => c.id === id)
    expect(conv?.status).toBe('idle')
    expect(markSessionReadMock).toHaveBeenCalledWith({ session_id: id })
  })

  it('createConversation initializes an empty message list', async () => {
    const id = await useChatStore.getState().createConversation()
    expect(useChatStore.getState().messages[id]).toEqual([])
  })

  it('setCurrentConversation does not reload messages that are already initialized', async () => {
    const id = await useChatStore.getState().createConversation()
    const loadMessagesMock = vi.mocked(AgentBinding.LoadSessionMessages)

    loadMessagesMock.mockClear()
    useChatStore.getState().setCurrentConversation(id)

    expect(loadMessagesMock).not.toHaveBeenCalled()
  })

  it('sendMessage adds optimistic user and assistant messages immediately', async () => {
    const id = await useChatStore.getState().createConversation()
    const generateTitleMock = vi.mocked(AgentBinding.GenerateTitle)
    generateTitleMock.mockClear()

    await useChatStore.getState().sendMessage({
      sessionId: id,
      content: 'hello world',
      baseUrl: 'https://example.com',
      apiKey: 'test-key',
      modelName: 'gpt-test',
      providerType: ProviderType.OpenAiCompatibility,
      enabledUserTools: [],
      attachments: [],
    })

    const sessionMessages = useChatStore.getState().messages[id]
    expect(sessionMessages).toHaveLength(2)
    expect(sessionMessages[0]).toMatchObject({
      sessionId: id,
      role: 'user',
      contentType: 'text',
      content: 'hello world',
    })
    expect(sessionMessages[1]).toMatchObject({
      sessionId: id,
      role: 'assistant',
      contentType: 'text',
      content: '',
      modelName: 'gpt-test',
    })
    expect(useChatStore.getState().streamingMessages[id]).toEqual({
      content: '',
      thinking: '',
    })
    expect(generateTitleMock).not.toHaveBeenCalled()
  })

  it('generates a title only after the first response reaches a terminal state', async () => {
    const id = await useChatStore.getState().createConversation()
    const generateTitleMock = vi.mocked(AgentBinding.GenerateTitle)
    generateTitleMock.mockClear()

    await useChatStore.getState().sendMessage({
      sessionId: id,
      content: 'hello world',
      baseUrl: 'https://example.com',
      apiKey: 'test-key',
      modelName: 'gpt-test',
      providerType: ProviderType.OpenAiCompatibility,
      enabledUserTools: [],
      attachments: [],
    })

    useChatStore.getState().initEventListeners()

    expect(generateTitleMock).not.toHaveBeenCalled()

    eventHandlers.get('agent:session:status')?.({
      data: {
        sessionId: id,
        status: 'waiting-unread',
      },
    })

    expect(generateTitleMock).not.toHaveBeenCalled()

    eventHandlers.get('agent:session:status')?.({
      data: {
        sessionId: id,
        status: 'done-unread',
      },
    })

    expect(generateTitleMock).toHaveBeenCalledTimes(1)
    expect(generateTitleMock).toHaveBeenCalledWith({
      session_id: id,
      base_url: 'https://example.com',
      api_key: 'test-key',
      model_name: 'gpt-test',
      provider_type: ProviderType.OpenAiCompatibility,
    })
  })

  it('initEventListeners appends stream chunks from Wails event data', () => {
    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 1,
        delta: 'hello',
        contentType: 'text',
      },
    })
    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 1,
        delta: ' world',
        contentType: 'text',
      },
    })
    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 1,
        delta: 'thinking',
        contentType: 'thinking',
      },
    })

    expect(useChatStore.getState().streamingMessages[1]).toBeUndefined()

    rafCallbacks.forEach((callback) => callback(0))
    rafCallbacks.clear()

    expect(useChatStore.getState().streamingMessages[1]).toEqual({
      content: 'hello world',
      thinking: 'thinking',
    })
  })

  it('allows event listeners to be reinitialized after cleanup', () => {
    const off = useChatStore.getState().initEventListeners()
    expect(eventHandlers.has('agent:stream:chunk')).toBe(true)

    off()
    expect(eventHandlers.has('agent:stream:chunk')).toBe(false)

    useChatStore.getState().initEventListeners()
    expect(eventHandlers.has('agent:stream:chunk')).toBe(true)
  })

  it('batches multiple stream chunks into one visible update per frame', () => {
    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 2,
        delta: 'hello',
        contentType: 'text',
      },
    })
    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 2,
        delta: ' world',
        contentType: 'text',
      },
    })

    expect(useChatStore.getState().streamingMessages[2]).toBeUndefined()
    expect(vi.mocked(requestAnimationFrame)).toHaveBeenCalledTimes(1)

    rafCallbacks.forEach((callback) => callback(0))
    rafCallbacks.clear()

    expect(useChatStore.getState().streamingMessages[2]).toEqual({
      content: 'hello world',
      thinking: '',
    })
  })

  it('preserves every chunk when Wails delivers batched stream events', () => {
    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:stream:chunk')?.({
      data: [
        {
          sessionId: 8,
          delta: 'A',
          contentType: 'text',
        },
        {
          sessionId: 8,
          delta: 'B',
          contentType: 'text',
        },
        {
          sessionId: 8,
          delta: 'C',
          contentType: 'text',
        },
      ],
    })

    rafCallbacks.forEach((callback) => callback(0))
    rafCallbacks.clear()

    expect(useChatStore.getState().streamingMessages[8]).toEqual({
      content: 'ABC',
      thinking: '',
    })
  })

  it('ignores stale chunk snapshots when stream events arrive out of order', () => {
    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 9,
        seq: 2,
        delta: 'B',
        content: 'AB',
        contentType: 'thinking',
      },
    })
    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 9,
        seq: 1,
        delta: 'A',
        content: 'A',
        contentType: 'thinking',
      },
    })

    rafCallbacks.forEach((callback) => callback(0))
    rafCallbacks.clear()

    expect(useChatStore.getState().streamingMessages[9]).toEqual({
      content: '',
      thinking: 'AB',
    })
  })

  it('flushes pending chunks before stream completion cleanup', async () => {
    const loadMessagesMock = vi.mocked(AgentBinding.LoadSessionMessages)
    loadMessagesMock.mockClear()

    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:stream:chunk')?.({
      data: {
        sessionId: 3,
        delta: 'final chunk',
        contentType: 'text',
      },
    })

    expect(useChatStore.getState().streamingMessages[3]).toBeUndefined()

    eventHandlers.get('agent:stream:done')?.({
      data: {
        sessionId: 3,
        usage: { input: 1, output: 1 },
      },
    })

    expect(useChatStore.getState().streamingMessages[3]).toBeUndefined()
    expect(loadMessagesMock).toHaveBeenCalledWith({
      session_id: 3,
      offset: 0,
      limit: 100,
    })
  })

  it('updates the latest streaming tool result instead of appending duplicates', () => {
    useChatStore.getState().initEventListeners()
    useChatStore.setState({
      messages: {
        5: [
          {
            id: 11,
            sessionId: 5,
            parentId: null,
            role: 'assistant',
            contentType: 'tool_call',
            content: '{"name":"shell","args":{"command":"echo hi"}}',
            modelName: '',
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: 'run shell',
            createdAt: new Date().toISOString(),
          },
          {
            id: -21,
            sessionId: 5,
            parentId: null,
            role: 'tool',
            contentType: 'tool_result',
            content: 'partial',
            modelName: '',
            agentName: 'shell',
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            createdAt: new Date().toISOString(),
          },
          {
            id: -22,
            sessionId: 5,
            parentId: null,
            role: 'assistant',
            contentType: 'text',
            content: '',
            modelName: '',
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            createdAt: new Date().toISOString(),
          },
        ],
      },
    })

    eventHandlers.get('agent:stream:tool_result')?.({
      data: {
        sessionId: 5,
        toolName: 'shell',
        result: 'partial\nmore',
      },
    })

    const sessionMessages = useChatStore.getState().messages[5]
    const toolResults = sessionMessages.filter((message) => message.contentType === 'tool_result')
    expect(toolResults).toHaveLength(1)
    expect(toolResults[0].content).toBe('partial\nmore')
  })

  it('keeps the active conversation read when done-unread status arrives', async () => {
    const id = await useChatStore.getState().createConversation()
    const markSessionReadMock = vi.mocked(AgentBinding.MarkSessionRead)

    useChatStore.setState({
      currentConversationId: id,
      conversations: useChatStore.getState().conversations.map(c =>
        c.id === id ? { ...c, status: 'loading' as const } : c
      ),
    })

    useChatStore.getState().initEventListeners()
    markSessionReadMock.mockClear()

    eventHandlers.get('agent:session:status')?.({
      data: {
        sessionId: id,
        status: 'done-unread',
      },
    })

    const conv = useChatStore.getState().conversations.find(c => c.id === id)
    expect(conv?.status).toBe('idle')
    expect(markSessionReadMock).toHaveBeenCalledWith({ session_id: id })
  })

  it('keeps waiting-unread for the active conversation so confirm UI can stay visible', async () => {
    const id = await useChatStore.getState().createConversation()
    const markSessionReadMock = vi.mocked(AgentBinding.MarkSessionRead)
    const generateTitleMock = vi.mocked(AgentBinding.GenerateTitle)

    useChatStore.setState({
      currentConversationId: id,
      conversations: useChatStore.getState().conversations.map(c =>
        c.id === id ? { ...c, status: 'loading' as const } : c
      ),
      pendingTitleGenerations: {
        [id]: {
          baseUrl: 'https://example.com',
          apiKey: 'test-key',
          modelName: 'gpt-test',
          providerType: ProviderType.OpenAiCompatibility,
        },
      },
    })

    useChatStore.getState().initEventListeners()
    markSessionReadMock.mockClear()
    generateTitleMock.mockClear()

    eventHandlers.get('agent:session:status')?.({
      data: {
        sessionId: id,
        status: 'waiting-unread',
      },
    })

    const conv = useChatStore.getState().conversations.find(c => c.id === id)
    expect(conv?.status).toBe('waiting-unread')
    expect(markSessionReadMock).not.toHaveBeenCalled()
    expect(generateTitleMock).not.toHaveBeenCalled()
  })

  it('refreshes the conversation list after a task session is spawned', () => {
    const listSessionsMock = vi.mocked(AgentBinding.ListSessions)
    listSessionsMock.mockClear()

    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:session:spawned')?.({
      data: {
        sessionId: 42,
        title: 'Install CLI from official docs',
        kind: 'task',
        userMessage: 'install feishu cli',
      },
    })

    expect(useChatStore.getState().conversations[0]).toMatchObject({
      id: 42,
      title: 'Install CLI from official docs',
      kind: 'task',
      status: 'loading',
    })
    expect(useChatStore.getState().messages[42]?.[0]).toMatchObject({
      role: 'user',
      contentType: 'text',
      content: 'install feishu cli',
    })
    expect(useChatStore.getState().messages[42]?.[1]).toMatchObject({
      role: 'assistant',
      contentType: 'text',
      content: '',
    })
    expect(listSessionsMock).toHaveBeenCalledWith({
      cursor: 0,
      limit: 20,
      starred_only: false,
      include_hidden: true,
    })
  })

  it('adds an assistant placeholder when an open session enters loading', async () => {
    const id = await useChatStore.getState().createConversation()

    useChatStore.setState({
      currentConversationId: id,
      conversations: useChatStore.getState().conversations.map(c =>
        c.id === id ? { ...c, status: 'idle' as const } : c
      ),
      sessionStatuses: {
        [id]: 'idle',
      },
      messages: {
        [id]: [{
          id: 100,
          sessionId: id,
          parentId: null,
          role: 'user',
          contentType: 'text',
          content: 'install this cli',
          modelName: '',
          agentName: '',
          tokensIn: 0,
          tokensOut: 0,
          extra: '',
          createdAt: '2026-05-26T00:00:00Z',
        }],
      },
    })

    useChatStore.getState().initEventListeners()

    eventHandlers.get('agent:session:status')?.({
      data: {
        sessionId: id,
        status: 'loading',
      },
    })

    const messages = useChatStore.getState().messages[id] ?? []
    expect(messages).toHaveLength(2)
    expect(messages[1]).toMatchObject({
      role: 'assistant',
      contentType: 'text',
      content: '',
    })
  })
})

describe('stream error handling', () => {
  it('calls loadMessages after stream:error event', () => {
    const loadMessagesMock = vi.mocked(AgentBinding.LoadSessionMessages)
    loadMessagesMock.mockClear()

    const off = useChatStore.getState().initEventListeners()

    useChatStore.setState({
      conversations: [{ id: 1, title: 'chat', kind: 'user' as const, createdAt: '', updatedAt: '', starred: false, status: 'loading' }],
      messages: { 1: [] },
      streamingMessages: { 1: { content: 'partial', thinking: '' } },
    })

    const handler = eventHandlers.get('agent:stream:error')
    expect(handler).toBeDefined()
    handler!({ data: { sessionId: 1, error: 'rate limit' } })

    expect(loadMessagesMock).toHaveBeenCalledWith(
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
      conversations: [{ id: 1, title: 'chat', kind: 'user' as const, createdAt: '', updatedAt: '', starred: false, status: 'idle' }],
      messages: { 1: [] },
      streamingMessages: {},
    })

    await useChatStore.getState().sendMessage({
      sessionId: 1,
      content: 'hello',
      baseUrl: 'http://localhost',
      apiKey: 'test',
      modelName: 'gpt-4',
      providerType: ProviderType.OpenAiCompatibility,
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
