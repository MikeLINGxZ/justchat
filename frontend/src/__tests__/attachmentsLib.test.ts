import { describe, expect, it } from 'vitest'
import {
  inferAttachmentMeta,
  ATTACHMENT_MAX_COUNT,
  ATTACHMENT_MAX_BYTES,
  isImageMime,
  isPdfMime,
} from '@/lib/attachments'

describe('inferAttachmentMeta', () => {
  it('classifies png as image', () => {
    const meta = inferAttachmentMeta('/foo/bar/baz.PNG')
    expect(meta.name).toBe('baz.PNG')
    expect(meta.mime).toBe('image/png')
    expect(meta.kind).toBe('image')
  })

  it('classifies pdf as file with application/pdf', () => {
    const meta = inferAttachmentMeta('/x/spec.pdf')
    expect(meta.kind).toBe('file')
    expect(meta.mime).toBe('application/pdf')
  })

  it('classifies unknown extension as octet-stream file', () => {
    const meta = inferAttachmentMeta('/x/data.xyz')
    expect(meta.kind).toBe('file')
    expect(meta.mime).toBe('application/octet-stream')
  })

  it('handles windows-style backslash path', () => {
    const meta = inferAttachmentMeta('C:\\Users\\me\\notes.md')
    expect(meta.name).toBe('notes.md')
    expect(meta.kind).toBe('file')
  })
})

describe('mime helpers', () => {
  it('isImageMime true for image/* false otherwise', () => {
    expect(isImageMime('image/png')).toBe(true)
    expect(isImageMime('application/pdf')).toBe(false)
  })
  it('isPdfMime checks pdf', () => {
    expect(isPdfMime('application/pdf')).toBe(true)
    expect(isPdfMime('text/plain')).toBe(false)
  })
})

describe('limits', () => {
  it('exports sane defaults', () => {
    expect(ATTACHMENT_MAX_COUNT).toBe(10)
    expect(ATTACHMENT_MAX_BYTES).toBe(20 * 1024 * 1024)
  })
})
