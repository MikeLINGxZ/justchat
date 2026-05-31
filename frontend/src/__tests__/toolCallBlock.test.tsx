import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { ToolCallBlock } from '@/components/chat/ToolCallBlock'
import { ToolCallsGroup } from '@/components/chat/ToolCallsGroup'
import type { DisplayMessage } from '@/types'

const toolMessage: DisplayMessage = {
  id: 1,
  sessionId: 1,
  parentId: null,
  role: 'assistant',
  contentType: 'tool_call',
  content: JSON.stringify({ name: 'shell', args: { command: 'pwd' } }),
  modelName: '',
  agentName: '',
  tokensIn: 0,
  tokensOut: 0,
  extra: 'Run command',
  createdAt: '2026-05-15T00:00:00Z',
  toolResult: '/tmp/project',
}

describe('ToolCallBlock', () => {
  it('uses compact thought styling with a lined detail panel', async () => {
    const user = userEvent.setup()
    render(
      <ToolCallBlock
        toolName="shell"
        purpose="Run command"
        args='{"command":"pwd"}'
        result="/tmp/project"
        status="completed"
      />
    )

    const toggle = screen.getByRole('button', { name: /shell/i })
    expect(toggle).toHaveClass('rounded-md')
    expect(toggle).not.toHaveClass('border')
    expect(screen.getByText('shell').nextElementSibling).toHaveClass('lucide-chevron-down')
    expect(toggle.querySelector('.rotate-0')).toBeInTheDocument()

    await user.click(toggle)

    expect(toggle.querySelector('.rotate-180')).toBeInTheDocument()
    expect(screen.getByText('Result').closest('.border-l')).toBeInTheDocument()
  })
})

describe('ToolCallsGroup', () => {
  it('uses the same compact lined style for grouped tools', async () => {
    const user = userEvent.setup()
    render(<ToolCallsGroup tools={[toolMessage]} />)

    const toggle = screen.getByRole('button', { name: /工具|tools used/i })
    expect(toggle).toHaveClass('rounded-md')
    expect(toggle).not.toHaveClass('border')
    expect(screen.getByText(/工具|tools used/i).nextElementSibling).toHaveClass('lucide-chevron-down')
    expect(toggle.querySelector('.rotate-0')).toBeInTheDocument()
    expect(toggle).not.toHaveTextContent('shell')

    await user.click(toggle)

    expect(toggle.querySelector('.rotate-180')).toBeInTheDocument()
    expect(screen.getAllByText('shell').some((node) => Boolean(node.closest('.border-l')))).toBe(true)
  })
})
