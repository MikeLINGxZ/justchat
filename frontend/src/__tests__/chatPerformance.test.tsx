import { Profiler } from 'react'
import { act, render } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ConversationList } from '@/components/sidebar/ConversationList'
import { useAppStore } from '@/store/appStore'
import { useChatStore } from '@/store/chatStore'

const eventHandlers = new Map<string, (event: { data: unknown }) => void>()

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
    CreateSession: vi.fn(),
    DeleteSession: vi.fn(),
    RenameSession: vi.fn(),
    ToggleStarSession: vi.fn(),
    ListSessions: vi.fn(),
    LoadSessionMessages: vi.fn(),
    MarkSessionRead: vi.fn(),
    SendMessage: vi.fn(),
    StopGeneration: vi.fn(),
    RespondToConfirm: vi.fn(),
    GenerateTitle: vi.fn(),
  },
}))

describe('chat rendering performance', () => {
  beforeEach(() => {
    eventHandlers.clear()
    useAppStore.setState({
      theme: 'auto',
      fontSize: 'xl',
      language: 'zh-CN',
    })
    useChatStore.setState({
      conversations: [
        {
          id: 1,
          title: 'Streaming session',
          kind: 'user' as const,
          createdAt: '2026-05-15T00:00:00Z',
          updatedAt: '2026-05-15T00:00:00Z',
          starred: false,
          status: 'loading',
        },
      ],
      messages: {
        1: [
          {
            id: 11,
            sessionId: 1,
            parentId: null,
            role: 'user',
            contentType: 'text',
            content: 'hello',
            modelName: '',
            agentName: '',
            tokensIn: 0,
            tokensOut: 0,
            extra: '',
            createdAt: '2026-05-15T00:00:00Z',
          },
        ],
      },
      currentConversationId: 1,
      streamingMessages: {
        1: {
          content: '',
          thinking: '',
        },
      },
      pendingConfirms: {},
      sessionsLoading: false,
      hasMoreSessions: true,
      sessionsCursor: 0,
    })
  })

  it('does not re-render the conversation list for stream chunks', () => {
    let commitCount = 0

    useChatStore.getState().initEventListeners()

    render(
      <Profiler
        id="conversation-list"
        onRender={() => {
          commitCount += 1
        }}
      >
        <ConversationList />
      </Profiler>
    )

    commitCount = 0

    act(() => {
      eventHandlers.get('agent:stream:chunk')?.({
        data: {
          sessionId: 1,
          delta: 'hello',
          contentType: 'text',
        },
      })
    })

    expect(commitCount).toBe(0)
  })
})
