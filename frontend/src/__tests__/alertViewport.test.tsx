import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it } from 'vitest'
import { AlertViewport } from '@/components/alert/AlertViewport'
import { useAlertStore } from '@/alert/store'

describe('AlertViewport', () => {
  beforeEach(() => {
    useAlertStore.getState().reset()
  })

  it('renders toast and banner regions from the shared store', () => {
    useAlertStore.getState().pushAlert({
      kind: 'success',
      placement: 'toast',
      title: 'Saved',
      message: 'Profile saved',
    })
    useAlertStore.getState().pushAlert({
      kind: 'warning',
      placement: 'banner',
      title: 'Attention',
      message: 'Banner warning',
    })

    render(<AlertViewport />)

    expect(screen.getByLabelText('Banner alerts')).toHaveClass('top-16')
    expect(screen.getByText('Profile saved')).toBeInTheDocument()
    expect(screen.getByText('Banner warning')).toBeInTheDocument()
  })

  it('shows duplicate count and exposes error alert semantics', () => {
    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'banner',
      title: 'Failed',
      message: 'Request failed',
    })
    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'banner',
      title: 'Failed',
      message: 'Request failed',
    })

    render(<AlertViewport />)

    const countBadge = screen.getByText('x2')
    const message = screen.getByText('Request failed')

    expect(countBadge).toBeInTheDocument()
    expect(countBadge.parentElement).toContainElement(message)
    expect(countBadge.compareDocumentPosition(message) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(screen.getByRole('alert')).toBeInTheDocument()
  })

  it('supports expanding details and dismissing from the keyboard', async () => {
    const user = userEvent.setup()

    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'banner',
      title: 'Failed',
      message: 'Failed to fetch model list',
      detail: 'stack trace',
    })

    render(<AlertViewport />)

    expect(screen.queryByText('stack trace')).not.toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /show details|展开详情/i }))
    expect(screen.getByText('stack trace')).toBeInTheDocument()

    const dismissButton = screen.getByRole('button', { name: /dismiss alert|关闭提醒/i })
    dismissButton.focus()
    await user.keyboard('{Enter}')
    expect(screen.queryByText('Failed to fetch model list')).not.toBeInTheDocument()
  })

  it('keeps error banner actions compact and aligned on the right', () => {
    useAlertStore.getState().pushAlert({
      kind: 'error',
      placement: 'banner',
      title: '',
      message: 'Failed to fetch model list',
      detail: 'stack trace',
    })

    render(<AlertViewport />)

    const alert = screen.getByRole('alert')
    expect(alert).toHaveClass('px-4', 'py-3')

    const showDetailsButton = screen.getByRole('button', { name: /show details|展开详情/i })
    const dismissButton = screen.getByRole('button', { name: /dismiss alert|关闭提醒/i })
    const message = screen.getByText('Failed to fetch model list')
    const messageRow = message.parentElement

    expect(showDetailsButton).toHaveClass('text-xs')
    expect(dismissButton).toHaveClass('text-lg')
    expect(message).not.toHaveClass('flex-1')
    expect(messageRow).toHaveClass('items-center')
  })

  it('renders at most two action buttons', () => {
    useAlertStore.getState().pushAlert({
      kind: 'warning',
      placement: 'banner',
      title: 'Action needed',
      message: 'Choose',
      actions: [
        { id: 'one', label: 'One', style: 'primary', closeOnClick: false },
        { id: 'two', label: 'Two', style: 'secondary', closeOnClick: false },
        { id: 'three', label: 'Three', style: 'danger', closeOnClick: false },
      ],
    })

    render(<AlertViewport />)

    expect(screen.getByRole('button', { name: 'One' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Two' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Three' })).not.toBeInTheDocument()
  })
})
