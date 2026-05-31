import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { AttachmentChips } from '@/components/chat/AttachmentChips'
import type { Attachment } from '@/types'

const items: Attachment[] = [
  { path: '/a/foo.png', name: 'foo.png', mime: 'image/png', kind: 'image' },
  { path: '/a/spec.pdf', name: 'spec.pdf', mime: 'application/pdf', kind: 'file' },
]

function renderWithI18n(node: React.ReactNode) {
  return render(<I18nextProvider i18n={i18n}>{node}</I18nextProvider>)
}

describe('AttachmentChips', () => {
  it('renders one chip per attachment', () => {
    renderWithI18n(<AttachmentChips items={items} variant="message" />)
    expect(screen.getByText('foo.png')).toBeInTheDocument()
    expect(screen.getByText('spec.pdf')).toBeInTheDocument()
  })

  it('invokes onRemove when × clicked in input variant', async () => {
    const onRemove = vi.fn()
    const user = userEvent.setup()
    renderWithI18n(
      <AttachmentChips items={items} variant="input" onRemove={onRemove} />
    )
    const removes = screen.getAllByLabelText(/remove attachment|移除附件/i)
    expect(removes).toHaveLength(2)
    await user.click(removes[1])
    expect(onRemove).toHaveBeenCalledWith(1)
  })

  it('does not render remove buttons in message variant', () => {
    renderWithI18n(<AttachmentChips items={items} variant="message" />)
    expect(screen.queryByLabelText(/remove attachment|移除附件/i)).toBeNull()
  })

  it('invokes onOpen when a message attachment is clicked', async () => {
    const onOpen = vi.fn()
    const user = userEvent.setup()
    renderWithI18n(<AttachmentChips items={items} variant="message" onOpen={onOpen} />)

    await user.click(screen.getByRole('button', { name: 'spec.pdf' }))
    expect(onOpen).toHaveBeenCalledWith(items[1])
  })

  it('renders nothing when items empty', () => {
    const { container } = renderWithI18n(
      <AttachmentChips items={[]} variant="input" />
    )
    expect(container.firstChild).toBeNull()
  })
})
