import { act, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { ChatInput } from '@/components/chat/ChatInput'
import { useChatStore } from '@/store/chatStore'
import { ATTACHMENT_MAX_COUNT } from '@/lib/attachments'

const eventListeners = new Map<string, Array<(event?: unknown) => void>>()

vi.mock('@wailsio/runtime', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wailsio/runtime')>()
  return {
    ...actual,
    Events: {
      On: vi.fn((eventName: string, callback: (event?: unknown) => void) => {
        const listeners = eventListeners.get(eventName) ?? []
        listeners.push(callback)
        eventListeners.set(eventName, listeners)
        return vi.fn(() => {
          const nextListeners = (eventListeners.get(eventName) ?? []).filter(listener => listener !== callback)
          eventListeners.set(eventName, nextListeners)
        })
      }),
      Types: {
        Mac: {
          WindowFileDraggingEntered: 'mac:WindowFileDraggingEntered',
          WindowFileDraggingExited: 'mac:WindowFileDraggingExited',
          WindowFileDraggingPerformed: 'mac:WindowFileDraggingPerformed',
        },
        Common: {
          WindowFilesDropped: 'common:WindowFilesDropped',
        },
      },
    },
  }
})

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  File: {
    SelectFile: vi.fn(),
    SaveTempFile: vi.fn(),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider', () => ({
  Provider: {
    ProviderAndModelList: vi.fn().mockResolvedValue({
      provider_models: [
        {
          provider: {
            id: 1, provider_name: 'P', enabled: true, is_default: true,
            base_url: 'http://x', api_key: 'k', provider_type: 'aliyun',
          },
          models: [{ id: 1, model: 'm', alias: 'm', is_default: true, enable: true }],
        },
      ],
    }),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin', () => ({
  Plugin: {
    ListAvailableTools: vi.fn().mockResolvedValue({ tools: [] }),
  },
}))

import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'

beforeEach(() => {
  eventListeners.clear()
  useChatStore.setState({
    conversations: [{
      id: 1, title: 't', kind: 'user' as const, createdAt: '', updatedAt: '',
      starred: false, status: 'idle',
    }],
    currentConversationId: 1,
  })
})

function renderInput() {
  return render(<I18nextProvider i18n={i18n}><ChatInput /></I18nextProvider>)
}

describe('ChatInput attachments', () => {
  it('adds a chip when SelectFile returns a path', async () => {
    ;(FileBinding.SelectFile as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce({ file_path: '/x/foo.png' })
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    await waitFor(() => expect(screen.getByText('foo.png')).toBeInTheDocument())
  })

  it('does nothing when SelectFile returns empty', async () => {
    ;(FileBinding.SelectFile as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce(null)
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    // No chip appears, so foo.png doesn't exist
    expect(screen.queryByText('foo.png')).toBeNull()
  })

  it('removes a chip via × button', async () => {
    ;(FileBinding.SelectFile as unknown as ReturnType<typeof vi.fn>).mockResolvedValueOnce({ file_path: '/x/foo.png' })
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    await screen.findByText('foo.png')

    await user.click(screen.getByLabelText(/remove attachment|移除附件/i))
    expect(screen.queryByText('foo.png')).toBeNull()
  })
})

describe('ChatInput clipboard paste', () => {
  beforeEach(() => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockReset()
  })

  it('calls SaveTempFile and adds chip when pasting an image', async () => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      file_path: '/tmp/lemontea/123-screenshot.png',
    })
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const file = new File(['fake-image-data'], 'screenshot.png', { type: 'image/png' })

    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: {
        items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }],
      },
    })
    container.dispatchEvent(pasteEvent)

    await waitFor(() => expect(screen.getByText('screenshot.png')).toBeInTheDocument())
    expect(FileBinding.SaveTempFile).toHaveBeenCalledOnce()
  })

  it('does not call SaveTempFile when pasting text only', async () => {
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'string', type: 'text/plain' }] },
    })
    container.dispatchEvent(pasteEvent)

    await waitFor(() => expect(FileBinding.SaveTempFile).not.toHaveBeenCalled())
  })

  it('respects ATTACHMENT_MAX_COUNT when pasting', async () => {
    // Fill attachments to max via SelectFile first
    const mockSelect = FileBinding.SelectFile as ReturnType<typeof vi.fn>
    for (let i = 0; i < ATTACHMENT_MAX_COUNT; i++) {
      mockSelect.mockResolvedValueOnce({ file_path: `/x/file${i}.png` })
    }
    const user = userEvent.setup()
    renderInput()
    for (let i = 0; i < ATTACHMENT_MAX_COUNT; i++) {
      await user.click(screen.getByLabelText(/attach file|附加文件/i))
    }
    await waitFor(() => expect(screen.getAllByText(/file\d\.png/).length).toBe(ATTACHMENT_MAX_COUNT))

    // Now try to paste one more — should be ignored
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      file_path: '/tmp/lemontea/extra.png',
    })
    const container = document.querySelector('.chat-input-area')!
    const file = new File(['data'], 'extra.png', { type: 'image/png' })
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }] },
    })
    container.dispatchEvent(pasteEvent)

    await waitFor(() => expect(FileBinding.SaveTempFile).not.toHaveBeenCalled())
  })

  it('skips file silently when SaveTempFile rejects', async () => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockRejectedValueOnce(new Error('disk full'))
    renderInput()

    const container = document.querySelector('.chat-input-area')!
    const file = new File(['data'], 'broken.png', { type: 'image/png' })
    const pasteEvent = new ClipboardEvent('paste', { bubbles: true, cancelable: true })
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: { items: [{ kind: 'file', type: 'image/png', getAsFile: () => file }] },
    })
    container.dispatchEvent(pasteEvent)

    await waitFor(() => expect(screen.queryByText('broken.png')).toBeNull())
  })
})

describe('ChatInput drag-and-drop', () => {
  beforeEach(() => {
    ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockReset()
  })

  it('adds chip when a file is dropped on the input area', async () => {
    renderInput()

    const listeners = eventListeners.get('files-dropped') ?? []
    expect(listeners).toHaveLength(1)

    act(() => {
      listeners[0]({
        data: {
          files: ['/tmp/lemontea/123-doc.pdf'],
        },
      })
    })

    await waitFor(() => expect(screen.getByText('123-doc.pdf')).toBeInTheDocument())
    expect(FileBinding.SaveTempFile).not.toHaveBeenCalled()
  })

  it('shows drop overlay on dragover with files', async () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', {
      value: { types: ['Files'] },
    })
    inputArea.dispatchEvent(dragoverEvent)

    await waitFor(() => expect(screen.getByText(/drop files to attach|释放以添加文件/i)).toBeInTheDocument())
  })

  it('hides drop overlay on dragleave when leaving the container', async () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!

    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', { value: { types: ['Files'] } })
    inputArea.dispatchEvent(dragoverEvent)
    await waitFor(() => expect(screen.getByText(/drop files to attach|释放以添加文件/i)).toBeInTheDocument())

    const dragleaveEvent = new DragEvent('dragleave', { bubbles: true, cancelable: true, relatedTarget: document.body })
    inputArea.dispatchEvent(dragleaveEvent)
    await waitFor(() => expect(screen.queryByText(/drop files to attach|释放以添加文件/i)).toBeNull())
  })

  it('caps dropped files at ATTACHMENT_MAX_COUNT', async () => {
    for (let i = 0; i < ATTACHMENT_MAX_COUNT; i++) {
      ;(FileBinding.SaveTempFile as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
        file_path: `/tmp/lemontea/img${i}.png`,
      })
    }
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const files = Array.from({ length: 15 }, (_, i) => new File(['d'], `img${i}.png`, { type: 'image/png' }))
    const dropEvent = new DragEvent('drop', { bubbles: true, cancelable: true })
    Object.defineProperty(dropEvent, 'dataTransfer', {
      value: { files, types: ['Files'] },
    })
    inputArea.dispatchEvent(dropEvent)

    await waitFor(() => expect(screen.getAllByText(/img\d+\.png/).length).toBe(ATTACHMENT_MAX_COUNT))
    expect(FileBinding.SaveTempFile).toHaveBeenCalledTimes(ATTACHMENT_MAX_COUNT)
  })

  it('does not show overlay when dragging non-file content', () => {
    renderInput()

    const inputArea = document.querySelector('.chat-input-area')!
    const dragoverEvent = new DragEvent('dragover', { bubbles: true, cancelable: true })
    Object.defineProperty(dragoverEvent, 'dataTransfer', {
      value: { types: ['text/plain'] },
    })
    inputArea.dispatchEvent(dragoverEvent)

    expect(screen.queryByText(/drop files to attach|释放以添加文件/i)).toBeNull()
  })
})
