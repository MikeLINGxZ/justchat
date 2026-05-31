import { describe, expect, it } from 'vitest'
import { getJsonError } from '@/components/settings/plugins/jsonError'

describe('getJsonError', () => {
  it('reports the offending line for trailing invalid characters inside an object', () => {
    const text = [
      '{',
      '  "mcpServers": {',
      '    "browsermcp": {',
      '      "command": "npx",',
      '      "args": ["@browsermcp/mcp@latest"]',
      '    }',
      '  }1',
      '}',
      '',
    ].join('\n')

    expect(getJsonError(text)).toMatchObject({
      valid: false,
      errorLine: 7,
    })
  })

  it('parses line and column from engines that do not expose a position', () => {
    const text = '{\n  "a": 1,\n}\n'
    const error = new SyntaxError('JSON Parse error: Expected property name or \'}\' at line 3 column 1')

    const originalParse = JSON.parse
    JSON.parse = (() => {
      throw error
    }) as typeof JSON.parse

    try {
      expect(getJsonError(text)).toMatchObject({
        valid: false,
        errorLine: 3,
      })
    } finally {
      JSON.parse = originalParse
    }
  })

  it('falls back to the actual invalid token when the parser only returns a generic syntax error', () => {
    const text = '{\n  "a": 1\n}1\n'
    const error = new SyntaxError('JSON Parse error: Unexpected identifier')

    const originalParse = JSON.parse
    JSON.parse = (() => {
      throw error
    }) as typeof JSON.parse

    try {
      expect(getJsonError(text)).toMatchObject({
        valid: false,
        errorLine: 3,
      })
    } finally {
      JSON.parse = originalParse
    }
  })
})
