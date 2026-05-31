export type InteractiveTerminalState = {
  visible: boolean
  output: string
  active: boolean
}

export function cleanTerminalOutput(output: string): string {
  const extracted = extractTerminalOutput(output, 0)
  return stripProgressEnvelope(extracted)
    .replace(/\x1b\[[0-9;?]*[ -/]*[@-~]/g, '')
    .replace(/^\[(stdout|stderr)\]\n?/gim, '')
    .replace(/^Tool confirmation:.*\n*/gim, '')
    .replace(/^\s*\n+/, '')
    .trimEnd()
}

function extractTerminalOutput(output: string, depth: number): string {
  if (depth > 4) return output

  const trimmed = output.trim()
  const jsonPayload = findJSONObject(trimmed)
  if (!jsonPayload) return output

  try {
    const parsed = JSON.parse(jsonPayload) as {
      terminal_output?: unknown
      output?: unknown
      qr?: unknown
      stdout?: unknown
      stderr?: unknown
      result?: unknown
    }
    if (typeof parsed.terminal_output === 'string') return extractTerminalOutput(parsed.terminal_output, depth + 1)
    if (typeof parsed.output === 'string') return extractTerminalOutput(parsed.output, depth + 1)
    if (typeof parsed.qr === 'string') return extractTerminalOutput(parsed.qr, depth + 1)
    if (typeof parsed.result === 'string') return extractTerminalOutput(parsed.result, depth + 1)

    const stdout = typeof parsed.stdout === 'string' ? parsed.stdout : ''
    const stderr = typeof parsed.stderr === 'string' ? parsed.stderr : ''
    if (stdout || stderr) return extractTerminalOutput([stdout, stderr].filter(Boolean).join('\n'), depth + 1)
  } catch {
    // Fall back to the raw tool output.
  }
  return output
}

function findJSONObject(text: string): string | null {
  const start = text.indexOf('{')
  const end = text.lastIndexOf('}')
  if (start < 0 || end <= start) return null
  return text.slice(start, end + 1)
}

function stripProgressEnvelope(output: string): string {
  let text = output.replace(/\r\n/g, '\n')

  if (text.includes('\\n') && /[█▀▄]|OUTPUT_SO_FAR|TIMEOUT_REACHED/.test(text)) {
    text = text.replace(/\\n/g, '\n')
  }

  const outputSoFar = text.match(/OUTPUT_SO_FAR:\s*([\s\S]*)/i)
  if (outputSoFar?.[1]) {
    text = outputSoFar[1]
  }

  text = text.replace(/^TIMEOUT_REACHED\s*\n?/gim, '')
  return text
}

export function getInteractiveTerminalState(output: string): InteractiveTerminalState | null {
  const payload = parseToolJSON(output)
  if (!payload || payload.interactive_terminal !== true) return null

  const status = String(payload.terminal_status ?? payload.status ?? 'active').toLowerCase()
  if (status === 'done' || status === 'hidden' || status === 'hide' || status === 'complete' || payload.visible === false) {
    return { visible: false, output: '', active: false }
  }

  const terminalOutput = stringValue(payload.terminal_output)
    ?? stringValue(payload.output)
    ?? stringValue(payload.qr)
    ?? [stringValue(payload.stdout), stringValue(payload.stderr)].filter(Boolean).join('\n')

  return {
    visible: true,
    output: cleanTerminalOutput(terminalOutput),
    active: status === 'active' || status === 'waiting',
  }
}

export function isInteractiveTerminalControl(output: string): boolean {
  return getInteractiveTerminalState(output) !== null
}

export function isInteractiveTerminalOutput(output: string): boolean {
  return getInteractiveTerminalState(output)?.visible === true
}

function parseToolJSON(output: string): Record<string, unknown> | null {
  if (!output) return null
  const jsonPayload = findJSONObject(output.trim())
  if (!jsonPayload) return null
  try {
    const parsed = JSON.parse(jsonPayload)
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed)
      ? parsed as Record<string, unknown>
      : null
  } catch {
    return null
  }
}

function stringValue(value: unknown): string | undefined {
  return typeof value === 'string' ? value : undefined
}
