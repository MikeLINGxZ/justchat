import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { ChatInput } from '@/components/chat/ChatInput'
import { useChatStore } from '@/store/chatStore'

vi.mock('@wailsio/runtime', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wailsio/runtime')>()
  return {
    ...actual,
    Events: {
      On: vi.fn(() => vi.fn()),
      Types: {
        Mac: {
          WindowFileDraggingEntered: 'mac:WindowFileDraggingEntered',
          WindowFileDraggingExited: 'mac:WindowFileDraggingExited',
          WindowFileDraggingPerformed: 'mac:WindowFileDraggingPerformed',
        },
      },
    },
  }
})

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider', () => ({
  Provider: {
    ProviderAndModelList: vi.fn().mockResolvedValue({
      provider_models: [
        {
          provider: {
            id: 1,
            provider_name: 'P',
            enabled: true,
            is_default: true,
            base_url: 'http://x',
            api_key: 'k',
            provider_type: 'aliyun',
          },
          models: [{ id: 1, model: 'm', alias: 'm', is_default: true, enable: true }],
        },
      ],
    }),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin', () => ({
  Plugin: {
    ListAvailableTools: vi.fn().mockResolvedValue({
      tools: [
        { id: 'web_fetch', name: 'web_fetch', description: 'Fetch a URL and return its content', category: 'builtin' },
        { id: 'mcp:filesystem:1.0.0', name: 'filesystem', description: 'Filesystem tools', category: 'mcp' },
      ],
    }),
  },
}))

const openSettings = vi.fn()
vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window', () => ({
  Window: {
    OpenSettings: (...args: unknown[]) => openSettings(...args),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  File: {
    SelectFile: vi.fn(),
    SaveTempFile: vi.fn(),
  },
}))

function renderInput() {
  return render(<I18nextProvider i18n={i18n}><ChatInput /></I18nextProvider>)
}

describe('Chat input tools', () => {
  beforeEach(async () => {
    openSettings.mockReset()
    useChatStore.setState({
      conversations: [{
        id: 1, title: 't', kind: 'user' as const, createdAt: '', updatedAt: '',
        starred: false, status: 'idle',
      }],
      currentConversationId: 1,
    })
    await i18n.changeLanguage('en')
  })

  it('renders aggregated mcp entries and opens plugin settings', async () => {
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByRole('button', { name: /tools/i }))
    expect(await screen.findByText('filesystem')).toBeInTheDocument()
    expect(screen.queryByText('read_file')).not.toBeInTheDocument()
    expect(screen.queryByText('Web Search')).not.toBeInTheDocument()
    expect(screen.queryByText('Search the web')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /manage tools & plugins/i }))
    await waitFor(() => expect(openSettings).toHaveBeenCalledWith({ tab: 'plugins' }))
  })
})
