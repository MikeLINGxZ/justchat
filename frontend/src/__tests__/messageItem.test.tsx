import type { ComponentProps, ReactNode } from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { MessageItem } from '@/components/chat/MessageItem'
import type { Message } from '@/types'
import i18n from '@/i18n'

type MarkdownProps = ComponentProps<typeof MessageItem> extends { message: infer _M }
  ? {
      children?: ReactNode
      components?: Record<string, (...args: any[]) => ReactNode>
    }
  : {
      children?: ReactNode
      components?: Record<string, (...args: any[]) => ReactNode>
    }

const markdownRenderSpy = vi.fn((props: MarkdownProps) => (
  <div data-testid="markdown-renderer">{props.children}</div>
))

vi.mock('react-markdown', () => ({
  default: (props: MarkdownProps) => markdownRenderSpy(props),
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  File: {
    OpenFile: vi.fn(),
  },
}))

import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'

const baseMessage: Message = {
  id: 1,
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
}

const assistantMessage: Message = {
  ...baseMessage,
  id: 2,
  role: 'assistant',
  content: '```ts\nconst value = 1\n```',
}

describe('MessageItem', () => {
  it('keeps distinct user bubble colors for light and dark modes', () => {
    render(<MessageItem message={baseMessage} />)

    expect(screen.getByText('hello').closest('.bg-primary')).toHaveClass(
      'bg-primary',
      'text-primary-foreground',
      'dark:bg-muted',
      'dark:text-foreground'
    )
  })

  it('renders read-only attachment chips inside the user bubble', () => {
    render(
      <MessageItem
        message={{
          ...baseMessage,
          attachments: [
            { path: '/tmp/photo.png', name: 'photo.png', mime: 'image/png', kind: 'image' },
            { path: '/tmp/spec.pdf', name: 'spec.pdf', mime: 'application/pdf', kind: 'file' },
          ],
        }}
      />
    )

    expect(screen.getByText('photo.png')).toBeInTheDocument()
    expect(screen.getByText('spec.pdf')).toBeInTheDocument()
    expect(screen.queryByLabelText(/remove attachment|移除附件/i)).toBeNull()
  })

  it('opens the attachment file when a message chip is clicked', async () => {
    const user = userEvent.setup()
    render(
      <MessageItem
        message={{
          ...baseMessage,
          attachments: [
            { path: '/tmp/spec.pdf', name: 'spec.pdf', mime: 'application/pdf', kind: 'file' },
          ],
        }}
      />
    )

    await user.click(screen.getByRole('button', { name: 'spec.pdf' }))
    expect(FileBinding.OpenFile).toHaveBeenCalledWith({ path: '/tmp/spec.pdf' })
  })

  it('keeps markdown rendering while the assistant message is still streaming', () => {
    markdownRenderSpy.mockClear()

    render(
      <MessageItem
        message={assistantMessage}
        isStreaming={true}
        streamingContent={assistantMessage.content}
      />
    )

    expect(markdownRenderSpy).toHaveBeenCalledOnce()
    expect(screen.getByTestId('markdown-renderer')).toHaveTextContent('const value = 1')
  })

  it('wraps assistant markdown with dedicated typography classes', () => {
    markdownRenderSpy.mockClear()

    render(<MessageItem message={assistantMessage} />)

    expect(markdownRenderSpy).toHaveBeenCalledOnce()
    expect(screen.getByTestId('markdown-renderer').closest('.llm-markdown')).toBeInTheDocument()
  })

  it('provides a scrollable table renderer for markdown tables', () => {
    markdownRenderSpy.mockClear()

    render(<MessageItem message={{ ...assistantMessage, content: '| A | B |\\n| - | - |\\n| 1 | 2 |' }} />)

    const markdownProps = markdownRenderSpy.mock.calls[0]?.[0] as MarkdownProps
    const tableRenderer = markdownProps.components?.table

    expect(tableRenderer).toBeTypeOf('function')

    const rendered = render(
      <>{tableRenderer?.({ children: <tbody><tr><td>1</td></tr></tbody> })}</>
    )

    expect(rendered.container.querySelector('.llm-markdown-table')).toBeInTheDocument()
    expect(rendered.container.querySelector('.llm-markdown-table table')).toBeInTheDocument()
  })

  it('distinguishes inline code styling from fenced code blocks', () => {
    markdownRenderSpy.mockClear()

    render(<MessageItem message={assistantMessage} />)

    const markdownProps = markdownRenderSpy.mock.calls[0]?.[0] as MarkdownProps
    const codeRenderer = markdownProps.components?.code

    expect(codeRenderer).toBeTypeOf('function')

    const inline = render(<>{codeRenderer?.({ children: 'npm install' })}</>)
    expect(inline.container.querySelector('code')).toHaveClass('llm-inline-code')

    const block = render(
      <>{codeRenderer?.({ children: 'const value = 1', className: 'language-ts' })}</>
    )
    expect(block.container.querySelector('code')).toHaveClass('language-ts')
    expect(block.container.querySelector('code')).not.toHaveClass('llm-inline-code')
  })

  it('renders ErrorMessageBubble for error content type', async () => {
    await i18n.changeLanguage('en')
    const errorMessage = {
      ...baseMessage,
      id: 99,
      role: 'assistant' as const,
      contentType: 'error' as const,
      content: JSON.stringify({ msg: 'Model response error', detail: 'timeout' }),
      modelName: 'qwen-plus',
      tokensIn: 128,
      tokensOut: 32,
    }
    render(<MessageItem message={errorMessage} />)
    expect(screen.getByText('Model response error')).toBeInTheDocument()
    expect(screen.getByText('qwen-plus')).toBeInTheDocument()
    expect(screen.getByText('128')).toBeInTheDocument()
    expect(screen.getByText('32')).toBeInTheDocument()
  })

  it('promotes explicitly active interactive terminal output outside the tool details', async () => {
    await i18n.changeLanguage('en')
    const toolMessage = {
      ...assistantMessage,
      id: 100,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'RunCliCommand', args: { argv: ['auth', 'login'] } }),
      extra: 'Running CLI command',
      toolResult: JSON.stringify({
        interactive_terminal: true,
        terminal_status: 'active',
        terminal_output: '[stderr]\n████████████\n████ █ █ ████\n████████████\nscan to login',
      }),
    }

    render(
      <MessageItem
        message={{ ...assistantMessage, id: 101, content: '' }}
        isStreaming={true}
        toolMessages={[toolMessage]}
      />
    )

    expect(screen.getByText('Waiting for scan or input')).toBeInTheDocument()
    expect(screen.getByText(/scan to login/)).toBeInTheDocument()
    expect(screen.queryByText('RunCliCommand')).toBeNull()
  })

  it('promotes explicitly active interactive terminal output from collapsed tool groups', async () => {
    await i18n.changeLanguage('en')
    const qrTool = {
      ...assistantMessage,
      id: 110,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'RunCliCommand', args: { argv: ['auth', 'login'] } }),
      extra: 'Running CLI command',
      toolResult: JSON.stringify({
        interactive_terminal: true,
        terminal_status: 'active',
        terminal_output: '[stderr]\n████████████\n████ █ █ ████\n████████████\nscan to login',
      }),
    }
    const otherTool = {
      ...assistantMessage,
      id: 111,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'file_read', args: { path: '/tmp/a' } }),
      extra: 'Read file',
      toolResult: 'ok',
    }

    render(
      <MessageItem
        message={{
          ...assistantMessage,
          id: 112,
          isToolGroup: true,
          groupedTools: [qrTool, otherTool],
        }}
      />
    )

    expect(screen.getByText('Waiting for scan or input')).toBeInTheDocument()
    expect(screen.getByText(/scan to login/)).toBeInTheDocument()
    expect(screen.getByText('1 tools used')).toBeInTheDocument()
  })

  it('hides interactive terminal output after a done control message', async () => {
    await i18n.changeLanguage('en')
    const activeTool = {
      ...assistantMessage,
      id: 120,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'InteractiveTerminal', args: { status: 'active' } }),
      extra: 'Show terminal',
      toolResult: JSON.stringify({
        interactive_terminal: true,
        terminal_status: 'active',
        terminal_output: 'scan to login',
      }),
    }
    const doneTool = {
      ...assistantMessage,
      id: 121,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'InteractiveTerminal', args: { status: 'done' } }),
      extra: 'Hide terminal',
      toolResult: JSON.stringify({
        interactive_terminal: true,
        terminal_status: 'done',
      }),
    }
    const otherTool = {
      ...assistantMessage,
      id: 122,
      contentType: 'tool_call' as const,
      content: JSON.stringify({ name: 'file_read', args: { path: '/tmp/a' } }),
      extra: 'Read file',
      toolResult: 'ok',
    }

    render(
      <MessageItem
        message={{
          ...assistantMessage,
          id: 123,
          isToolGroup: true,
          groupedTools: [activeTool, doneTool, otherTool],
        }}
      />
    )

    expect(screen.queryByText('Interactive terminal output')).toBeNull()
    expect(screen.queryByText(/scan to login/)).toBeNull()
    expect(screen.getByText('1 tools used')).toBeInTheDocument()
  })
})
