export type ParsedIError = {
  msg: string
  detail: string
}

type UnknownRecord = Record<string, unknown>

function isRecord(value: unknown): value is UnknownRecord {
  return typeof value === 'object' && value !== null
}

function normalizeString(value: unknown): string | null {
  return typeof value === 'string' && value.length > 0 ? value : null
}

function parseStructuredIError(value: unknown): ParsedIError | null {
  if (!isRecord(value)) {
    return null
  }

  const msg = normalizeString(value.msg)
  const detail = normalizeString(value.detail)
  if (!msg && !detail) {
    return null
  }

  return {
    msg: msg ?? detail ?? 'Unknown error',
    detail: detail ?? msg ?? 'Unknown error',
  }
}

function parseSerializedIError(value: string): ParsedIError | null {
  try {
    const parsed = JSON.parse(value) as unknown
    return parseUnknownIError(parsed)
  } catch {
    return null
  }
}

export function parseUnknownIError(value: unknown): ParsedIError {
  const structured = parseStructuredIError(value)
  if (structured) {
    return structured
  }

  if (typeof value === 'string') {
    const parsed = parseSerializedIError(value)
    if (parsed) {
      return parsed
    }

    return { msg: value, detail: value }
  }

  if (value instanceof Error) {
    const parsedFromMessage = parseSerializedIError(value.message)
    if (parsedFromMessage) {
      return parsedFromMessage
    }

    if ('cause' in value) {
      return parseUnknownIError((value as Error & { cause?: unknown }).cause ?? value.message)
    }

    return { msg: value.message, detail: value.message }
  }

  if (isRecord(value)) {
    const messageValue = value.message
    if (typeof messageValue === 'string') {
      const parsed = parseSerializedIError(messageValue)
      if (parsed) {
        return parsed
      }

      if (isRecord(value.cause)) {
        return parseUnknownIError(value.cause)
      }

      return { msg: messageValue, detail: messageValue }
    }

    if (messageValue !== undefined) {
      return parseUnknownIError(messageValue)
    }

    if (value.cause !== undefined) {
      return parseUnknownIError(value.cause)
    }
  }

  const fallback = String(value)
  return { msg: fallback, detail: fallback }
}
