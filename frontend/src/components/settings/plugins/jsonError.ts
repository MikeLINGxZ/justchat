export type JsonError =
  | { valid: true }
  | { valid: false, errorLine: number, errorMessage: string }

export function getJsonError(text: string): JsonError {
  if (text.trim() === '') return { valid: false, errorLine: 1, errorMessage: 'empty' }

  try {
    JSON.parse(text)
    return { valid: true }
  } catch (err) {
    if (err instanceof SyntaxError) {
      const pos = extractErrorPos(text, err.message)
      const line = getLineNumberAtPos(text, pos)
      return { valid: false, errorLine: Math.max(line, 1), errorMessage: err.message }
    }

    return { valid: false, errorLine: 1, errorMessage: 'parse error' }
  }
}

function getLineNumberAtPos(text: string, pos: number): number {
  return text.slice(0, Math.min(Math.max(pos, 0), text.length)).split('\n').length
}

function getPosFromLineAndColumn(text: string, line: number, column: number): number {
  if (line <= 1) return Math.max(column - 1, 0)

  let currentLine = 1

  for (let i = 0; i < text.length; i++) {
    if (currentLine === line) {
      return i + Math.max(column - 1, 0)
    }
    if (text[i] === '\n') currentLine += 1
  }

  return text.length
}

function extractErrorPos(text: string, message: string): number {
  const positionMatch = message.match(/position\s+(\d+)/i)
  if (positionMatch) return parseInt(positionMatch[1], 10)

  const lineColumnMatch = message.match(/line\s+(\d+)\s+column\s+(\d+)/i)
  if (lineColumnMatch) {
    const [, line, column] = lineColumnMatch
    return getPosFromLineAndColumn(text, parseInt(line, 10), parseInt(column, 10))
  }

  return scanErrorPos(text)
}

function scanErrorPos(text: string): number {
  const scanner = new JsonScanner(text)
  return scanner.findErrorPos()
}

class JsonScanner {
  private index = 0

  constructor(private readonly text: string) {}

  findErrorPos(): number {
    this.skipWhitespace()
    if (this.index >= this.text.length) return this.text.length

    const valueError = this.parseValue()
    if (valueError !== null) return valueError

    this.skipWhitespace()
    return this.index < this.text.length ? this.index : this.text.length
  }

  private parseValue(): number | null {
    if (this.index >= this.text.length) return this.text.length

    const ch = this.text[this.index]
    if (ch === '{') return this.parseObject()
    if (ch === '[') return this.parseArray()
    if (ch === '"') return this.parseString()
    if (ch === '-' || this.isDigit(ch)) return this.parseNumber()
    if (this.text.startsWith('true', this.index)) {
      this.index += 4
      return null
    }
    if (this.text.startsWith('false', this.index)) {
      this.index += 5
      return null
    }
    if (this.text.startsWith('null', this.index)) {
      this.index += 4
      return null
    }

    return this.index
  }

  private parseObject(): number | null {
    this.index += 1
    this.skipWhitespace()

    if (this.peek() === '}') {
      this.index += 1
      return null
    }

    while (this.index < this.text.length) {
      if (this.peek() !== '"') return this.index

      const keyError = this.parseString()
      if (keyError !== null) return keyError

      this.skipWhitespace()
      if (this.peek() !== ':') return this.index
      this.index += 1

      this.skipWhitespace()
      const valueError = this.parseValue()
      if (valueError !== null) return valueError

      this.skipWhitespace()
      const ch = this.peek()
      if (ch === ',') {
        this.index += 1
        this.skipWhitespace()
        continue
      }
      if (ch === '}') {
        this.index += 1
        return null
      }
      return this.index
    }

    return this.text.length
  }

  private parseArray(): number | null {
    this.index += 1
    this.skipWhitespace()

    if (this.peek() === ']') {
      this.index += 1
      return null
    }

    while (this.index < this.text.length) {
      const valueError = this.parseValue()
      if (valueError !== null) return valueError

      this.skipWhitespace()
      const ch = this.peek()
      if (ch === ',') {
        this.index += 1
        this.skipWhitespace()
        continue
      }
      if (ch === ']') {
        this.index += 1
        return null
      }
      return this.index
    }

    return this.text.length
  }

  private parseString(): number | null {
    this.index += 1

    while (this.index < this.text.length) {
      const ch = this.text[this.index]

      if (ch === '"') {
        this.index += 1
        return null
      }

      if (ch === '\\') {
        this.index += 1
        if (this.index >= this.text.length) return this.text.length

        const escaped = this.text[this.index]
        if (escaped === 'u') {
          for (let i = 1; i <= 4; i++) {
            const hex = this.text[this.index + i]
            if (hex === undefined || !/[0-9a-fA-F]/.test(hex)) return this.index + i
          }
          this.index += 5
          continue
        }

        this.index += 1
        continue
      }

      if (ch === '\n' || ch === '\r') return this.index
      this.index += 1
    }

    return this.text.length
  }

  private parseNumber(): number | null {
    const start = this.index

    if (this.peek() === '-') this.index += 1
    if (!this.isDigit(this.peek())) return this.index

    if (this.peek() === '0') {
      this.index += 1
    } else {
      this.consumeDigits()
    }

    if (this.peek() === '.') {
      this.index += 1
      if (!this.isDigit(this.peek())) return this.index
      this.consumeDigits()
    }

    if (this.peek() === 'e' || this.peek() === 'E') {
      this.index += 1
      if (this.peek() === '+' || this.peek() === '-') this.index += 1
      if (!this.isDigit(this.peek())) return this.index
      this.consumeDigits()
    }

    return this.index === start ? start : null
  }

  private consumeDigits() {
    while (this.isDigit(this.peek())) {
      this.index += 1
    }
  }

  private skipWhitespace() {
    while (/\s/.test(this.peek())) {
      this.index += 1
    }
  }

  private peek(): string {
    return this.text[this.index] ?? ''
  }

  private isDigit(ch: string): boolean {
    return ch >= '0' && ch <= '9'
  }
}
