export type Theme = 'auto' | 'light' | 'dark'
export type FontSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'
export type Language = 'zh-CN' | 'en'
export type ConversationTab = 'chats' | 'favorites'
export type ConversationStatus =
  | 'idle'
  | 'loading'
  | 'done-unread'
  | 'error-unread'
  | 'waiting-unread'

export type MessageContentType =
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'thinking'
  | 'confirm_request'
  | 'confirm_response'
  | 'error'

export type MessageRole = 'user' | 'assistant' | 'tool' | 'system'

export type AttachmentKind = 'image' | 'file'

export type Attachment = {
  path: string
  name: string
  mime: string
  kind: AttachmentKind
}

export type Conversation = {
  id: number
  title: string
  kind: 'user' | 'task'
  tags?: string[]
  createdAt: string
  updatedAt: string
  starred: boolean
  status: ConversationStatus
}

export type Message = {
  id: number
  sessionId: number
  parentId: number | null
  role: MessageRole
  contentType: MessageContentType
  content: string
  modelName: string
  agentName: string
  tokensIn: number
  tokensOut: number
  extra: string
  attachments?: Attachment[]
  createdAt: string
}

export type DisplayMessage = Message & {
  toolResult?: string
  toolConfirmAction?: 'approve' | 'reject' | 'comment'
  toolConfirmComment?: string
  isToolGroup?: boolean
  groupedTools?: DisplayMessage[]
}

export type StreamChunkEvent = {
  sessionId: number
  seq?: number
  delta: string
  content?: string
  contentType: 'text' | 'thinking'
}

export type ToolCallEvent = {
  sessionId: number
  toolName: string
  args: string
  purpose: string
}

export type ToolResultEvent = {
  sessionId: number
  toolName: string
  result: string
}

export type ConfirmRequestEvent = {
  sessionId: number
  requestId: string
  toolName: string
  args: string
  purpose: string
}

export type StreamDoneEvent = {
  sessionId: number
  usage: { input: number; output: number }
}

export type StreamErrorEvent = {
  sessionId: number
  error: string
}

export type SessionStatusEvent = {
  sessionId: number
  status: ConversationStatus
}

export type Tool = {
  id: string
  name: string
  description: string
  enabled: boolean
}

export type ModelProvider = {
  id: string
  name: string
  models: Model[]
}

export type Model = {
  id: string
  name: string
  providerId: string
}
