import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ConversationList } from '@/components/sidebar/ConversationList'
import i18n from '@/i18n'
import { useAppStore } from '@/store/appStore'
import { useChatStore } from '@/store/chatStore'

vi.mock('@wailsio/runtime', () => ({
  Events: {
    On: vi.fn(),
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

describe('ConversationList', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
    useAppStore.setState({
      theme: 'auto',
      fontSize: 'md',
      language: 'en',
    })
    useChatStore.setState({
      conversations: [
        {
          id: 1,
          title: 'Test chat',
          kind: 'user' as const,
          createdAt: '2026-05-18T00:00:00Z',
          updatedAt: '2026-05-18T00:00:00Z',
          starred: false,
          status: 'idle',
        },
      ],
      currentConversationId: 1,
    })
  })

  it('uses the thin sidebar scrollbar style on the conversation scroll container', () => {
    render(<ConversationList />)

    const footer = screen.getByText('All chats loaded (1)')
    const scrollContainer = footer.parentElement

    expect(scrollContainer).toHaveClass('sidebar-scrollbar')
  })

  it('pins only task conversations in the task section', () => {
    useChatStore.setState({
      conversations: [
        {
          id: 1,
          title: 'Regular chat',
          kind: 'user' as const,
          createdAt: '2026-05-18T00:00:00Z',
          updatedAt: '2026-05-18T00:00:00Z',
          starred: false,
          status: 'idle',
        },
        {
          id: 2,
          title: 'Install task',
          kind: 'task' as const,
          createdAt: '2026-05-18T00:00:00Z',
          updatedAt: '2026-05-18T00:00:00Z',
          starred: false,
          status: 'loading',
        },
        {
          id: 3,
          title: 'Old background',
          kind: 'user' as const,
          createdAt: '2026-05-18T00:00:00Z',
          updatedAt: '2026-05-18T00:00:00Z',
          starred: false,
          status: 'loading',
        },
      ],
      currentConversationId: 1,
    })

    render(<ConversationList />)

    const taskHeading = screen.getByText('Tasks')
    const taskSection = taskHeading.parentElement

    expect(taskSection).toHaveTextContent('Install task')
    expect(taskSection).not.toHaveTextContent('Old background')
    expect(screen.getByText('Old background')).toBeInTheDocument()
  })
})
