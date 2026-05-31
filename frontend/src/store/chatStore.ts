import { create } from 'zustand'
import { Events } from '@wailsio/runtime'
import { Agent as AgentBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent'
import type { Type as ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider/models'
import type {
  Attachment,
  ConfirmRequestEvent,
  Conversation,
  ConversationStatus,
  Message,
  SessionStatusEvent,
  StreamChunkEvent,
  StreamDoneEvent,
  StreamErrorEvent,
  ToolCallEvent,
  ToolResultEvent,
} from '../types'

type StreamingMessageState = { content: string; thinking: string }
type PendingStreamingPatch = {
  contentDelta: string
  thinkingDelta: string
  contentSnapshot?: string
  thinkingSnapshot?: string
}
type PendingTitleGeneration = {
  baseUrl: string
  apiKey: string
  modelName: string
  providerType: ProviderType
}
type TitleGenerationConfig = PendingTitleGeneration

interface ChatStore {
  conversations: Conversation[]
  messages: Record<number, Message[]>
  sessionStatuses: Record<number, ConversationStatus | undefined>
  currentConversationId: number | null
  streamingMessages: Record<number, StreamingMessageState>
  pendingTitleGenerations: Record<number, PendingTitleGeneration | undefined>
  sessionRunConfigs: Record<number, PendingTitleGeneration | undefined>
  pendingConfirms: Record<number, {
    requestId: string
    toolName: string
    args: string
    purpose: string
  } | null>
  sessionsLoading: boolean
  hasMoreSessions: boolean
  sessionsCursor: number
  loadSessions: (reset?: boolean) => Promise<void>
  loadMessages: (sessionId: number) => Promise<void>
  createConversation: (params?: { title?: string; tags?: string[] }) => Promise<number>
  deleteConversation: (id: number) => Promise<void>
  renameConversation: (id: number, title: string) => Promise<void>
  toggleStar: (id: number, starred: boolean) => Promise<void>
  markAsRead: (id: number) => Promise<void>
  setCurrentConversation: (id: number | null) => void
  setConversationStatus: (id: number, status: ConversationStatus) => void
  sendMessage: (params: {
    sessionId: number
    content: string
    baseUrl: string
    apiKey: string
    modelName: string
    providerType: ProviderType
    enabledUserTools: string[]
    attachments: Attachment[]
    systemPrompt?: string
    skillName?: string
  }) => Promise<void>
  stopGeneration: (sessionId: number) => Promise<void>
  respondToConfirm: (
    sessionId: number,
    approved: boolean,
    message: string,
    action?: 'approve' | 'reject' | 'comment'
  ) => Promise<void>
  generateTitle: (
    sessionId: number,
    baseUrl: string,
    apiKey: string,
    modelName: string,
    providerType: ProviderType
  ) => Promise<void>
  generateTitleDraft: (sessionId: number, fallbackConfig?: TitleGenerationConfig) => Promise<string | null>
  initEventListeners: () => () => void
}

function eventItems<T>(event: { data: T | T[] }): T[] {
  return Array.isArray(event.data) ? event.data : [event.data]
}

const EMPTY_STREAMING_MESSAGE: StreamingMessageState = { content: '', thinking: '' }
const EMPTY_STREAMING_PATCH: PendingStreamingPatch = { contentDelta: '', thinkingDelta: '' }

function scheduleNextFrame(callback: FrameRequestCallback): number {
  if (typeof globalThis.requestAnimationFrame === 'function') {
    return globalThis.requestAnimationFrame(callback)
  }

  return globalThis.setTimeout(() => callback(Date.now()), 16) as unknown as number
}

function cancelScheduledFrame(handle: number) {
  if (typeof globalThis.cancelAnimationFrame === 'function') {
    globalThis.cancelAnimationFrame(handle)
    return
  }

  globalThis.clearTimeout(handle)
}

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

function isStreamingStatus(status: ConversationStatus | undefined): boolean {
  return status === 'loading' || status === 'waiting-unread'
}

export const useChatStore = create<ChatStore>()((set, get) => {
  const pendingStreamingBuffers = new Map<number, PendingStreamingPatch>()
  const scheduledStreamingFlushes = new Map<number, number>()
  const latestStreamingChunkSeq = new Map<string, number>()

  const clearStreamingBuffer = (sessionId: number) => {
    const handle = scheduledStreamingFlushes.get(sessionId)
    if (handle !== undefined) {
      cancelScheduledFrame(handle)
      scheduledStreamingFlushes.delete(sessionId)
    }
    pendingStreamingBuffers.delete(sessionId)
  }

  const clearStreamingSeqs = (sessionId: number) => {
    latestStreamingChunkSeq.delete(`${sessionId}:content`)
    latestStreamingChunkSeq.delete(`${sessionId}:thinking`)
  }

  const flushStreamingBuffer = (sessionId: number) => {
    const pending = pendingStreamingBuffers.get(sessionId)
    clearStreamingBuffer(sessionId)

    if (
      !pending ||
      (!pending.contentDelta &&
        !pending.thinkingDelta &&
        pending.contentSnapshot === undefined &&
        pending.thinkingSnapshot === undefined)
    ) return

    set((state) => {
      const previous = state.streamingMessages[sessionId] ?? EMPTY_STREAMING_MESSAGE

      return {
        streamingMessages: {
          ...state.streamingMessages,
          [sessionId]: {
            content: pending.contentSnapshot ?? previous.content + pending.contentDelta,
            thinking: pending.thinkingSnapshot ?? previous.thinking + pending.thinkingDelta,
          },
        },
      }
    })
  }

  const scheduleStreamingFlush = (sessionId: number) => {
    if (scheduledStreamingFlushes.has(sessionId)) return

    const handle = scheduleNextFrame(() => {
      flushStreamingBuffer(sessionId)
    })

    scheduledStreamingFlushes.set(sessionId, handle)
  }

  const appendStreamingDelta = (data: StreamChunkEvent) => {
    const key = data.contentType === 'thinking' ? 'thinking' : 'content'
    const seqKey = `${data.sessionId}:${key}`
    if (typeof data.seq === 'number') {
      const previousSeq = latestStreamingChunkSeq.get(seqKey) ?? 0
      if (data.seq <= previousSeq) return
      latestStreamingChunkSeq.set(seqKey, data.seq)
    }

    const previous = pendingStreamingBuffers.get(data.sessionId) ?? EMPTY_STREAMING_PATCH
    if (typeof data.content === 'string') {
      pendingStreamingBuffers.set(data.sessionId, {
        ...previous,
        [`${key}Snapshot`]: data.content,
        [`${key}Delta`]: '',
      })
    } else {
      const deltaKey = `${key}Delta` as 'contentDelta' | 'thinkingDelta'
      pendingStreamingBuffers.set(data.sessionId, {
        ...previous,
        [deltaKey]: previous[deltaKey] + data.delta,
      })
    }

    scheduleStreamingFlush(data.sessionId)
  }

  const ensureStreamingPlaceholder = (sessionId: number) => {
    set((state) => {
      const conversation = state.conversations.find((item) => item.id === sessionId)
      const status = state.sessionStatuses[sessionId] ?? conversation?.status
      if (!isStreamingStatus(status)) {
        return state
      }

      const sessionMessages = state.messages[sessionId]
      if (!sessionMessages) {
        return state
      }

      const lastMessage = sessionMessages[sessionMessages.length - 1]
      if (lastMessage?.role === 'assistant' && lastMessage?.contentType !== 'tool_call') {
        return state
      }

      return {
        messages: {
          ...state.messages,
          [sessionId]: [
            ...sessionMessages,
            {
              id: -(sessionId * 1000) - sessionMessages.length - 1,
              sessionId,
              parentId: lastMessage?.role === 'user' ? lastMessage.id : null,
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
      }
    })
  }

  return ({
  conversations: [],
  messages: {},
  sessionStatuses: {},
  currentConversationId: null,
  streamingMessages: {},
  pendingTitleGenerations: {},
  sessionRunConfigs: {},
  pendingConfirms: {},
  sessionsLoading: false,
  hasMoreSessions: true,
  sessionsCursor: 0,

  loadSessions: async (reset = false) => {
    if (get().sessionsLoading) return
    set({ sessionsLoading: true })

    try {
      const cursor = reset ? 0 : get().sessionsCursor
      const result = await AgentBinding.ListSessions({
        cursor,
        limit: 20,
        starred_only: false,
        include_hidden: true,
      })
      if (!result) return

      const nextConversations: Conversation[] = (result.sessions ?? []).map((session) => ({
        id: session.id,
        title: session.title,
        kind: ((session.tags ?? []).includes('task') || session.kind === 'task' ? 'task' : 'user') as Conversation['kind'],
        tags: session.tags ?? (session.kind === 'task' ? ['task'] : []),
        createdAt: session.created,
        updatedAt: session.updated,
        starred: session.starred,
        status: session.status as ConversationStatus,
      }))

      set((state) => ({
        conversations: reset ? nextConversations : [...state.conversations, ...nextConversations],
        sessionStatuses: {
          ...state.sessionStatuses,
          ...Object.fromEntries(nextConversations.map((conversation) => [conversation.id, conversation.status])),
        },
        hasMoreSessions: result.has_more,
        sessionsCursor: result.next_cursor,
      }))
    } finally {
      set({ sessionsLoading: false })
    }
  },

  loadMessages: async (sessionId) => {
    const result = await AgentBinding.LoadSessionMessages({
      session_id: sessionId,
      offset: 0,
      limit: 100,
    })
    if (!result) return

    const messages: Message[] = (result.messages ?? []).map((message) => {
      let attachments: Attachment[] | undefined
      const raw = message.attachments
      if (raw && raw.trim() !== '') {
        try {
          const parsed = JSON.parse(raw) as Attachment[]
          if (Array.isArray(parsed) && parsed.length > 0) attachments = parsed
        } catch {
          // ignore malformed payload
        }
      }
      return {
        id: message.id,
        sessionId: message.session_id,
        parentId: message.parent_id ?? null,
        role: message.role as Message['role'],
        contentType: message.content_type as Message['contentType'],
        content: message.content,
        modelName: message.model_name,
        agentName: message.agent_name,
        tokensIn: message.tokens_in,
        tokensOut: message.tokens_out,
        extra: message.extra,
        attachments,
        createdAt: message.created,
      }
    })

    // For task sessions still running, append a temp assistant placeholder
    // so the streaming message container is valid when the user opens the session
    const sessionStatus =
      get().sessionStatuses[sessionId] ??
      get().conversations.find((c) => c.id === sessionId)?.status
    const lastMsg = messages[messages.length - 1]
    if (
      (sessionStatus === 'loading' || sessionStatus === 'waiting-unread') &&
      (!lastMsg || lastMsg.role !== 'assistant' || lastMsg.contentType === 'tool_call')
    ) {
      messages.push({
        id: -(sessionId * 1000),
        sessionId,
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
      })
    }

    set((state) => ({
      messages: { ...state.messages, [sessionId]: messages },
    }))
  },

  createConversation: async (params) => {
    const result = await AgentBinding.CreateSession({ title: params?.title ?? '', tags: params?.tags ?? [] } as any)
    if (!result) return 0

    const tags = result.tags ?? params?.tags ?? []
    const conversation: Conversation = {
      id: result.session_id,
      title: result.title,
      kind: tags.includes('task') ? 'task' : 'user',
      tags,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      starred: false,
      status: 'idle',
    }

    set((state) => ({
      conversations: [conversation, ...state.conversations],
      messages: { ...state.messages, [result.session_id]: [] },
      pendingTitleGenerations: {
        ...state.pendingTitleGenerations,
        [result.session_id]: undefined,
      },
      sessionRunConfigs: {
        ...state.sessionRunConfigs,
        [result.session_id]: undefined,
      },
      sessionStatuses: {
        ...state.sessionStatuses,
        [result.session_id]: 'idle',
      },
    }))
    return result.session_id
  },

  deleteConversation: async (id) => {
    await AgentBinding.DeleteSession({ session_id: id })
    clearStreamingBuffer(id)
    clearStreamingSeqs(id)
    set((state) => {
      const nextMessages = { ...state.messages }
      const nextPendingTitleGenerations = { ...state.pendingTitleGenerations }
      const nextSessionRunConfigs = { ...state.sessionRunConfigs }
      const nextPendingConfirms = { ...state.pendingConfirms }
      const nextStreamingMessages = { ...state.streamingMessages }
      const nextSessionStatuses = { ...state.sessionStatuses }

      delete nextMessages[id]
      delete nextPendingTitleGenerations[id]
      delete nextSessionRunConfigs[id]
      delete nextPendingConfirms[id]
      delete nextStreamingMessages[id]
      delete nextSessionStatuses[id]

      return {
        conversations: state.conversations.filter((conversation) => conversation.id !== id),
        currentConversationId: state.currentConversationId === id ? null : state.currentConversationId,
        messages: nextMessages,
        pendingTitleGenerations: nextPendingTitleGenerations,
        sessionRunConfigs: nextSessionRunConfigs,
        pendingConfirms: nextPendingConfirms,
        streamingMessages: nextStreamingMessages,
        sessionStatuses: nextSessionStatuses,
      }
    })
  },

  renameConversation: async (id, title) => {
    await AgentBinding.RenameSession({ session_id: id, title })
    set((state) => ({
      conversations: state.conversations.map((conversation) =>
        conversation.id === id
          ? { ...conversation, title, updatedAt: new Date().toISOString() }
          : conversation
      ),
    }))
  },

  toggleStar: async (id, starred) => {
    await AgentBinding.ToggleStarSession({ session_id: id, starred })
    set((state) => ({
      conversations: state.conversations.map((conversation) =>
        conversation.id === id ? { ...conversation, starred } : conversation
      ),
    }))
  },

  markAsRead: async (id) => {
    const conversation = get().conversations.find((item) => item.id === id)
    if (!conversation || conversation.status === 'idle' || conversation.status === 'loading') return

    set((state) => ({
      conversations: state.conversations.map((item) =>
        item.id === id ? { ...item, status: 'idle' } : item
      ),
    }))

    await AgentBinding.MarkSessionRead({ session_id: id })
  },

  setCurrentConversation: (id) => {
    if (id) void get().markAsRead(id)
    if (id && get().messages[id] === undefined) {
      void get().loadMessages(id)
    }
    set({ currentConversationId: id })
  },

  setConversationStatus: (id, status) =>
    set((state) => ({
      conversations: state.conversations.map((conversation) =>
        conversation.id === id ? { ...conversation, status } : conversation
      ),
    })),

  sendMessage: async (params) => {
    const now = new Date().toISOString()
    const tempUserId = -Date.now()
    const tempAssistantId = tempUserId - 1
    const isFirstMessage = (get().messages[params.sessionId]?.length ?? 0) === 0
    clearStreamingBuffer(params.sessionId)
    clearStreamingSeqs(params.sessionId)

    set((state) => ({
      conversations: state.conversations.map((conversation) =>
        conversation.id === params.sessionId
          ? {
              ...conversation,
              status: 'loading',
              updatedAt: now,
            }
          : conversation
      ),
      messages: {
        ...state.messages,
        [params.sessionId]: [
          ...(state.messages[params.sessionId] ?? []),
          {
            id: tempUserId,
            sessionId: params.sessionId,
            parentId: null,
            role: 'user',
            contentType: 'text',
            content: params.content,
            modelName: '',
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            attachments: params.attachments.length > 0 ? params.attachments : undefined,
            createdAt: now,
          },
          {
            id: tempAssistantId,
            sessionId: params.sessionId,
            parentId: tempUserId,
            role: 'assistant',
            contentType: 'text',
            content: '',
            modelName: params.modelName,
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            createdAt: now,
          },
        ],
      },
      pendingConfirms: {
        ...state.pendingConfirms,
        [params.sessionId]: null,
      },
      sessionStatuses: {
        ...state.sessionStatuses,
        [params.sessionId]: 'loading',
      },
      sessionRunConfigs: {
        ...state.sessionRunConfigs,
        [params.sessionId]: {
          baseUrl: params.baseUrl,
          apiKey: params.apiKey,
          modelName: params.modelName,
          providerType: params.providerType,
        },
      },
      streamingMessages: {
        ...state.streamingMessages,
        [params.sessionId]: EMPTY_STREAMING_MESSAGE,
      },
      pendingTitleGenerations: {
        ...state.pendingTitleGenerations,
        [params.sessionId]: isFirstMessage
          ? {
              baseUrl: params.baseUrl,
              apiKey: params.apiKey,
              modelName: params.modelName,
              providerType: params.providerType,
            }
          : state.pendingTitleGenerations[params.sessionId],
      },
    }))

    try {
      await AgentBinding.SendMessage({
        session_id: params.sessionId,
        content: params.content,
        system_prompt: params.systemPrompt ?? '',
        skill_name: params.skillName ?? '',
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
      } as any)
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
          sessionStatuses: {
            ...state.sessionStatuses,
            [params.sessionId]: 'idle',
          },
          pendingTitleGenerations: {
            ...state.pendingTitleGenerations,
            [params.sessionId]: undefined,
          },
        }
      })
    }
  },

  stopGeneration: async (sessionId) => {
    await AgentBinding.StopGeneration({ session_id: sessionId })
  },

  respondToConfirm: async (sessionId, approved, message, action = approved ? 'approve' : 'reject') => {
    await AgentBinding.RespondToConfirm({
      session_id: sessionId,
      approved,
      message,
      action,
    })
    set((state) => ({
      pendingConfirms: { ...state.pendingConfirms, [sessionId]: null },
    }))
  },

  generateTitle: async (sessionId, baseUrl, apiKey, modelName, providerType) => {
    const result = await AgentBinding.GenerateTitle({
      session_id: sessionId,
      base_url: baseUrl,
      api_key: apiKey,
      model_name: modelName,
      provider_type: providerType,
    })
    if (!result?.title) return

    set((state) => ({
      conversations: state.conversations.map((conversation) =>
        conversation.id === sessionId ? { ...conversation, title: result.title } : conversation
      ),
    }))
  },

  generateTitleDraft: async (sessionId, fallbackConfig) => {
    const runConfig = get().sessionRunConfigs[sessionId] ?? fallbackConfig
    if (!runConfig) return null

    const result = await AgentBinding.GenerateTitle({
      session_id: sessionId,
      base_url: runConfig.baseUrl,
      api_key: runConfig.apiKey,
      model_name: runConfig.modelName,
      provider_type: runConfig.providerType,
    })

    return result?.title?.trim() || null
  },

  initEventListeners: () => {
    if ((globalThis as { __lemonteaChatListenersInitialized?: boolean }).__lemonteaChatListenersInitialized) {
      return () => {}
    }
    ;(globalThis as { __lemonteaChatListenersInitialized?: boolean }).__lemonteaChatListenersInitialized = true

    const offs: Array<() => void> = []

    offs.push(Events.On('agent:stream:chunk', (event: { data: StreamChunkEvent | StreamChunkEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        ensureStreamingPlaceholder(data.sessionId)
        appendStreamingDelta(data)
      }
    }))

    offs.push(Events.On('agent:stream:confirm_request', (event: { data: ConfirmRequestEvent | ConfirmRequestEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        set((state) => ({
          pendingConfirms: {
            ...state.pendingConfirms,
            [data.sessionId]: {
              requestId: data.requestId,
              toolName: data.toolName,
              args: data.args,
              purpose: data.purpose,
            },
          },
        }))
      }
    }))

    offs.push(Events.On('agent:stream:tool_call', (event: { data: ToolCallEvent | ToolCallEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        let parsedArgs: Record<string, unknown>
        try {
          parsedArgs = JSON.parse(data.args) as Record<string, unknown>
        } catch {
          parsedArgs = {}
        }
        const payload = JSON.stringify({ name: data.toolName, args: parsedArgs })
        const tempId = -Date.now() - (Math.random() * 1000 | 0)

        set((state) => {
          const sessionMessages = state.messages[data.sessionId] ?? []
          const insertIndex = sessionMessages.length > 0 ? sessionMessages.length - 1 : 0

          const toolCallMessage: Message = {
            id: tempId,
            sessionId: data.sessionId,
            parentId: null,
            role: 'assistant',
            contentType: 'tool_call',
            content: payload,
            modelName: '',
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: data.purpose,
            createdAt: new Date().toISOString(),
          }

          const newMessages = [...sessionMessages]
          newMessages.splice(insertIndex, 0, toolCallMessage)

          return {
            messages: {
              ...state.messages,
              [data.sessionId]: newMessages,
            },
          }
        })
      }
    }))

    offs.push(Events.On('agent:stream:tool_result', (event: { data: ToolResultEvent | ToolResultEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        set((state) => {
          const sessionMessages = state.messages[data.sessionId] ?? []
          const insertIndex = sessionMessages.length > 0 ? sessionMessages.length - 1 : 0
          const latestToolResultIndex = findLatestStreamToolResultIndex(sessionMessages, data.toolName)

          if (latestToolResultIndex >= 0) {
            const nextMessages = [...sessionMessages]
            nextMessages[latestToolResultIndex] = {
              ...nextMessages[latestToolResultIndex],
              content: data.result,
              createdAt: new Date().toISOString(),
            }
            return {
              messages: {
                ...state.messages,
                [data.sessionId]: nextMessages,
              },
            }
          }

          const tempId = -Date.now() - (Math.random() * 1000 | 0)

          const toolResultMessage: Message = {
            id: tempId,
            sessionId: data.sessionId,
            parentId: null,
            role: 'tool',
            contentType: 'tool_result',
            content: data.result,
            modelName: '',
            agentName: data.toolName,
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            createdAt: new Date().toISOString(),
          }

          const newMessages = [...sessionMessages]
          newMessages.splice(insertIndex, 0, toolResultMessage)

          return {
            messages: {
              ...state.messages,
              [data.sessionId]: newMessages,
            },
          }
        })
      }
    }))

    offs.push(Events.On('agent:stream:done', (event: { data: StreamDoneEvent | StreamDoneEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        flushStreamingBuffer(data.sessionId)
        clearStreamingSeqs(data.sessionId)

        set((state) => {
          const nextStreaming = { ...state.streamingMessages }
          delete nextStreaming[data.sessionId]
          return {
            streamingMessages: nextStreaming,
            sessionStatuses: {
              ...state.sessionStatuses,
              [data.sessionId]: 'idle',
            },
          }
        })

        void get().loadMessages(data.sessionId)
      }
    }))

    offs.push(Events.On('agent:stream:error', (event: { data: StreamErrorEvent | StreamErrorEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        flushStreamingBuffer(data.sessionId)
        clearStreamingSeqs(data.sessionId)

        set((state) => {
          const nextStreaming = { ...state.streamingMessages }
          delete nextStreaming[data.sessionId]
          return {
            streamingMessages: nextStreaming,
            sessionStatuses: {
              ...state.sessionStatuses,
              [data.sessionId]: 'error-unread',
            },
          }
        })

        void get().loadMessages(data.sessionId)
      }
    }))

    offs.push(Events.On('agent:session:spawned', (event: { data: { sessionId: number; title: string; kind: string; tags?: string[]; userMessage?: string } | { sessionId: number; title: string; kind: string; tags?: string[]; userMessage?: string }[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue
        if (get().conversations.some((c) => c.id === data.sessionId)) continue

        const now = new Date().toISOString()
        const tempUserId = -Date.now()
        const tempAssistantId = tempUserId - 1
        const tags = data.tags ?? (data.kind === 'task' ? ['task'] : [])
        const newConversation: Conversation = {
          id: data.sessionId,
          title: data.title,
          kind: tags.includes('task') ? 'task' : 'user',
          tags,
          createdAt: now,
          updatedAt: now,
          starred: false,
          status: 'loading',
        }

        set((state) => ({
          conversations: [newConversation, ...state.conversations],
          messages: data.userMessage
            ? {
                ...state.messages,
                [data.sessionId]: [
                  {
                    id: tempUserId,
                    sessionId: data.sessionId,
                    parentId: null,
                    role: 'user',
                    contentType: 'text',
                    content: data.userMessage,
                    modelName: '',
                    agentName: '',
                    tokensIn: 0,
                    tokensOut: 0,
                    extra: '',
                    createdAt: now,
                  },
                  {
                    id: tempAssistantId,
                    sessionId: data.sessionId,
                    parentId: tempUserId,
                    role: 'assistant',
                    contentType: 'text',
                    content: '',
                    modelName: '',
                    agentName: '',
                    tokensIn: 0,
                    tokensOut: 0,
                    extra: '',
                    createdAt: now,
                  },
                ],
              }
            : state.messages,
          sessionStatuses: {
            ...state.sessionStatuses,
            [data.sessionId]: 'loading',
          },
        }))

        void get().loadSessions(true)
      }
    }))

    offs.push(Events.On('agent:session:status', (event: { data: SessionStatusEvent | SessionStatusEvent[] }) => {
      for (const data of eventItems(event)) {
        if (!data) continue

        // Unknown session (e.g. a task spawned from another window) — reload list
        const known = get().conversations.some((c) => c.id === data.sessionId)
        if (!known) {
          void get().loadSessions(true)
          continue
        }

        const shouldFinalizeCurrentConversation =
          data.status === 'done-unread' || data.status === 'error-unread'

        const pendingTitleGeneration = get().pendingTitleGenerations[data.sessionId]
        if (shouldFinalizeCurrentConversation && pendingTitleGeneration) {
          set((state) => ({
            pendingTitleGenerations: {
              ...state.pendingTitleGenerations,
              [data.sessionId]: undefined,
            },
          }))
          void get().generateTitle(
            data.sessionId,
            pendingTitleGeneration.baseUrl,
            pendingTitleGeneration.apiKey,
            pendingTitleGeneration.modelName,
            pendingTitleGeneration.providerType
          )
        }

        if (data.sessionId === get().currentConversationId && shouldFinalizeCurrentConversation) {
          set((state) => ({
            sessionStatuses: {
              ...state.sessionStatuses,
              [data.sessionId]: 'idle',
            },
            conversations: state.conversations.map((conversation) =>
              conversation.id === data.sessionId ? { ...conversation, status: 'idle' } : conversation
            ),
          }))
          void AgentBinding.MarkSessionRead({ session_id: data.sessionId })
          continue
        }

        set((state) => ({
          sessionStatuses: {
            ...state.sessionStatuses,
            [data.sessionId]: data.status,
          },
          conversations: state.conversations.map((conversation) =>
            conversation.id === data.sessionId ? { ...conversation, status: data.status } : conversation
          ),
        }))

        ensureStreamingPlaceholder(data.sessionId)
      }
    }))

    return () => {
      offs.forEach((off) => off())
      ;(globalThis as { __lemonteaChatListenersInitialized?: boolean }).__lemonteaChatListenersInitialized = false
    }
  },
  })
})

function findLatestStreamToolResultIndex(messages: Message[], toolName: string): number {
  for (let i = messages.length - 1; i >= 0; i--) {
    const message = messages[i]
    if (message.contentType === 'tool_result' && message.agentName === toolName) {
      return i
    }
    if (message.contentType === 'tool_call') {
      break
    }
  }
  return -1
}
