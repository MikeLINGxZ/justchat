import { render, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { InteractiveTerminalBlock } from '@/components/chat/InteractiveTerminalBlock'

const writeSpy = vi.fn()

vi.mock('xterm', () => ({
  Terminal: vi.fn().mockImplementation(function TerminalMock() {
    return {
      open: vi.fn(),
      loadAddon: vi.fn(),
      write: writeSpy,
      focus: vi.fn(),
      onData: vi.fn(() => ({ dispose: vi.fn() })),
      reset: vi.fn(),
      dispose: vi.fn(),
    }
  }),
}))

vi.mock('xterm-addon-fit', () => ({
  FitAddon: vi.fn().mockImplementation(function FitAddonMock() {
    return {
      fit: vi.fn(),
      cols: 80,
      rows: 24,
    }
  }),
}))

vi.mock('xterm-addon-web-links', () => ({
  WebLinksAddon: vi.fn().mockImplementation(function WebLinksAddonMock() {
    return {}
  }),
}))

describe('InteractiveTerminalBlock', () => {
  it('streams ANSI output into xterm instead of rendering raw terminal controls', async () => {
    render(<InteractiveTerminalBlock output={'\x1b[32mDone\x1b[0m'} active />)

    await waitFor(() => {
      expect(writeSpy).toHaveBeenCalledWith('\x1b[32mDone\x1b[0m')
    })
  })
})
