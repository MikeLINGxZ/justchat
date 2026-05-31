import type { Attachment, AttachmentKind } from '@/types'

export const ATTACHMENT_MAX_COUNT = 10
export const ATTACHMENT_MAX_BYTES = 20 * 1024 * 1024

const EXT_MIME: Record<string, string> = {
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.webp': 'image/webp',
  '.gif': 'image/gif',
  '.pdf': 'application/pdf',
  '.txt': 'text/plain',
  '.log': 'text/plain',
  '.md': 'text/markdown',
  '.markdown': 'text/markdown',
  '.json': 'application/json',
  '.html': 'text/html',
  '.css': 'text/css',
  '.js': 'text/javascript',
  '.ts': 'application/typescript',
  '.tsx': 'application/typescript',
  '.go': 'text/x-go',
  '.py': 'text/x-python',
  '.rs': 'text/x-rust',
  '.yaml': 'application/yaml',
  '.yml': 'application/yaml',
  '.toml': 'application/toml',
}

function basename(path: string): string {
  // Supports both unix and windows path separators.
  const idx = Math.max(path.lastIndexOf('/'), path.lastIndexOf('\\'))
  return idx >= 0 ? path.slice(idx + 1) : path
}

function extOf(path: string): string {
  const name = basename(path)
  const dot = name.lastIndexOf('.')
  return dot >= 0 ? name.slice(dot).toLowerCase() : ''
}

export function inferAttachmentMeta(path: string): Attachment {
  const name = basename(path)
  const mime = EXT_MIME[extOf(path)] ?? 'application/octet-stream'
  const kind: AttachmentKind = mime.startsWith('image/') ? 'image' : 'file'
  return { path, name, mime, kind }
}

export function isImageMime(mime: string): boolean {
  return mime.startsWith('image/')
}

export function isPdfMime(mime: string): boolean {
  return mime === 'application/pdf'
}
