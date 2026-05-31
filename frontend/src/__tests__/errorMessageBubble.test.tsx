import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { ErrorMessageBubble } from '@/components/chat/ErrorMessageBubble'
import i18n from '@/i18n'

describe('ErrorMessageBubble', () => {
  it('renders the user-facing msg from JSON content', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.getByText('Model response error')).toBeInTheDocument()
  })

  it('left aligns the red error text row', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)

    const textRow = screen.getByText('Model response error').closest('div')
    expect(textRow).toHaveClass('w-full', 'justify-start', 'text-left')
    expect(textRow).not.toHaveClass('justify-center', 'text-center')
  })

  it('renders model metadata below the error text using the normal response footer style', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(
      <ErrorMessageBubble content={content} modelName="qwen-plus" tokensIn={128} tokensOut={32} />
    )

    const modelInfo = screen.getByText('qwen-plus').closest('div')
    expect(modelInfo).toHaveClass('mt-0.5', 'flex', 'w-full', 'items-center', 'gap-3', 'text-xs', 'text-muted-foreground')
    expect(screen.getByText('128')).toBeInTheDocument()
    expect(screen.getByText('32')).toBeInTheDocument()
  })

  it('does not show detail text by default', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.queryByText('rate limit exceeded')).not.toBeInTheDocument()
  })

  it('expands detail when the icon button is clicked', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(
      <ErrorMessageBubble content={content} modelName="qwen-plus" tokensIn={128} tokensOut={32} />
    )
    const button = screen.getByRole('button')
    fireEvent.click(button)
    const detail = screen.getByText('rate limit exceeded')
    expect(detail).toBeInTheDocument()

    const textRow = screen.getByText('Model response error').closest('div')
    const modelInfo = screen.getByText('qwen-plus').closest('div')
    const parent = textRow?.parentElement
    expect(parent?.children[0]).toBe(textRow)
    expect(parent?.children[1]).toBe(detail.closest('div'))
    expect(parent?.children[2]).toBe(modelInfo)
  })

  it('collapses detail on second click', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Model response error', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    const button = screen.getByRole('button')
    fireEvent.click(button)
    fireEvent.click(button)
    expect(screen.queryByText('rate limit exceeded')).not.toBeInTheDocument()
  })

  it('renders raw string content when JSON parse fails', () => {
    render(<ErrorMessageBubble content="plain error text" />)
    expect(screen.getByText('plain error text')).toBeInTheDocument()
  })

  it('localizes known stream error titles to the current language', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: '大模型响应出错', detail: 'rate limit exceeded' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.getByText('Model response error')).toBeInTheDocument()
  })

  it('shows no expand button when detail is empty', async () => {
    await i18n.changeLanguage('en')
    const content = JSON.stringify({ msg: 'Something failed', detail: '' })
    render(<ErrorMessageBubble content={content} />)
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })
})
