import { describe, expect, it } from 'vitest'
import { cleanTerminalOutput, getInteractiveTerminalState, isInteractiveTerminalOutput } from '@/lib/terminalOutput'

describe('terminalOutput', () => {
  it('unwraps timed out command JSON and keeps only QR terminal output', () => {
    const wrapped = JSON.stringify({
      interactive_terminal: true,
      terminal_status: 'active',
      terminal_output: 'TIMEOUT_REACHED\nOUTPUT_SO_FAR: ████████\n██    ██\n████████\n',
    })

    expect(cleanTerminalOutput(wrapped)).toBe('████████\n██    ██\n████████')
    expect(isInteractiveTerminalOutput(wrapped)).toBe(true)
    expect(getInteractiveTerminalState(wrapped)).toEqual({
      visible: true,
      output: '████████\n██    ██\n████████',
      active: true,
    })
  })

  it('extracts QR output from QRCode tool JSON', () => {
    const wrapped = JSON.stringify({
      text: 'https://example.com',
      qr: '████████\n██    ██\n████████\n',
      interactive_terminal: true,
      terminal_status: 'active',
      terminal_output: '████████\n██    ██\n████████\n',
    })

    expect(cleanTerminalOutput(wrapped)).toBe('████████\n██    ██\n████████')
    expect(isInteractiveTerminalOutput(wrapped)).toBe(true)
  })

  it('does not promote terminal-looking output without an explicit AI decision', () => {
    const raw = '[stderr]\n████████\n██    ██\n████████\nscan to login'

    expect(isInteractiveTerminalOutput(raw)).toBe(false)
    expect(getInteractiveTerminalState(raw)).toBeNull()
  })

  it('treats done as a hidden terminal state', () => {
    const wrapped = JSON.stringify({
      interactive_terminal: true,
      terminal_status: 'done',
    })

    expect(isInteractiveTerminalOutput(wrapped)).toBe(false)
    expect(getInteractiveTerminalState(wrapped)).toEqual({
      visible: false,
      output: '',
      active: false,
    })
  })
})
