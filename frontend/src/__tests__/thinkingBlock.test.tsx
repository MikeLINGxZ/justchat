import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { ThinkingBlock } from '@/components/chat/ThinkingBlock'

describe('ThinkingBlock', () => {
  it('shows active thinking label while streaming', () => {
    render(<ThinkingBlock content="reasoning" active={true} />)
    expect(screen.getByText('思考中')).toBeInTheDocument()
  })

  it('shows completed thinking label after streaming finishes', () => {
    render(<ThinkingBlock content="reasoning" active={false} />)
    expect(screen.getByText('思考')).toBeInTheDocument()
  })

  it('uses compact thought styling and reveals content when expanded', async () => {
    const user = userEvent.setup()
    render(<ThinkingBlock content="reasoning details" active={false} />)

    const toggle = screen.getByRole('button', { name: /思考/ })
    expect(toggle).toHaveClass('rounded-md')
    expect(toggle).toHaveClass('w-full')
    expect(toggle).not.toHaveClass('border')
    expect(toggle.querySelector('svg')).toBeInTheDocument()
    expect(screen.getByText('思考').nextElementSibling).toHaveClass('lucide-chevron-down')
    expect(toggle.querySelector('.rotate-0')).toBeInTheDocument()
    expect(screen.queryByText('reasoning details')).not.toBeInTheDocument()

    await user.click(toggle)

    expect(toggle.querySelector('.rotate-180')).toBeInTheDocument()
    expect(screen.getByText('reasoning details')).toHaveClass('border-l')
  })
})
