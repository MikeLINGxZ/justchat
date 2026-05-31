import { memo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import 'highlight.js/styles/github-dark-dimmed.min.css'
import { ThinkingBlock } from './ThinkingBlock'
import { ToolCallBlock } from './ToolCallBlock'
import { ToolCallsGroup } from './ToolCallsGroup'
import { ToolConfirmCard } from './ToolConfirmCard'
import { AttachmentChips } from './AttachmentChips'
import { ErrorMessageBubble } from './ErrorMessageBubble'
import { ResponseMeta } from './ResponseMeta'
import { InteractiveTerminalBlock } from '@/components/chat/InteractiveTerminalBlock'
import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'
import type { DisplayMessage } from '@/types'
import { cn } from '@/lib/utils'
import {
  getInteractiveTerminalState,
  isInteractiveTerminalControl,
  type InteractiveTerminalState,
} from '@/lib/terminalOutput'

interface Props {
  message: DisplayMessage
  isStreaming?: boolean
  streamingContent?: string
  streamingThinking?: string
  toolMessages?: DisplayMessage[]
  pendingConfirm?: {
    requestId: string
    toolName: string
    args: string
    purpose: string
  } | null
  onConfirm?: (message: string) => void
  onReject?: (message: string) => void
  onSubmitComment?: (message: string) => void
}

export const MessageItem = memo(function MessageItem({
  message,
  isStreaming,
  streamingContent,
  streamingThinking,
  toolMessages,
  pendingConfirm,
  onConfirm,
  onReject,
  onSubmitComment,
}: Props) {
  const { role, contentType, content, modelName, tokensIn, tokensOut, extra } = message
  const openAttachment = (path: string) => {
    void FileBinding.OpenFile({ path })
  }

  if (role === 'user' && contentType === 'text') {
    const attachments = message.attachments ?? []

    return (
      <div className="flex justify-end">
        <div className="max-w-[70%] select-text rounded-2xl rounded-br-sm bg-primary text-primary-foreground px-4 py-2.5 text-sm dark:bg-muted dark:text-foreground">
          {content && <div className="whitespace-pre-wrap break-words">{content}</div>}
          {attachments.length > 0 && (
            <AttachmentChips
              items={attachments}
              variant="message"
              onOpen={(attachment) => openAttachment(attachment.path)}
            />
          )}
        </div>
      </div>
    )
  }

  if (role === 'user' && contentType === 'confirm_response') {
    return null
  }

  if (message.isToolGroup && message.groupedTools) {
    return <ToolCallsGroup tools={message.groupedTools} />
  }

  if (contentType === 'tool_call') {
    let toolName = ''
    let args = content
    let status: 'completed' | 'failed' | 'executing' | 'user_commented' | 'user_rejected' = 'completed'
    try {
      const parsed = JSON.parse(content)
      toolName = parsed.name ?? ''
      args = JSON.stringify(parsed.args ?? parsed, null, 2)
    } catch {
      // Fall back to raw content.
    }

    if (message.toolConfirmAction === 'comment') {
      status = 'user_commented'
    } else if (message.toolConfirmAction === 'reject') {
      status = 'user_rejected'
    }

    const terminalState = message.toolResult ? getInteractiveTerminalState(message.toolResult) : null
    if (terminalState) {
      if (!terminalState.visible) return null
      return (
        <InteractiveTerminalBlock
          output={terminalState.output}
          active={terminalState.active || Boolean(isStreaming)}
        />
      )
    }

    return (
      <ToolCallBlock
        toolName={toolName}
        purpose={extra}
        args={args}
        result={message.toolResult ?? ''}
        userComment={message.toolConfirmComment}
        status={status}
      />
    )
  }

  if (contentType === 'thinking') {
    return <ThinkingBlock content={content} defaultExpanded={false} active={false} />
  }

  if (contentType === 'error') {
    return (
      <ErrorMessageBubble
        content={content}
        modelName={modelName}
        tokensIn={tokensIn}
        tokensOut={tokensOut}
      />
    )
  }

  const displayContent = isStreaming && streamingContent !== undefined ? streamingContent : content
  const displayThinking = isStreaming && streamingThinking !== undefined ? streamingThinking : ''
  const shouldRenderMarkdown = Boolean(displayContent)
  const flattenedToolMessages = flattenToolMessages(toolMessages ?? [])
  const terminalState = isStreaming ? latestInteractiveTerminalState(flattenedToolMessages) : null
  const terminalControlIds = new Set(
    flattenedToolMessages
      .filter((toolMsg) => isInteractiveTerminalControl(toolMsg.toolResult ?? ''))
      .map((toolMsg) => toolMsg.id)
  )
  const nonInteractiveToolMessages = flattenedToolMessages.filter((toolMsg) => !terminalControlIds.has(toolMsg.id))

  return (
    <div className="flex flex-col gap-0.5">
      {displayThinking && (
        <ThinkingBlock
          content={displayThinking}
          defaultExpanded={Boolean(isStreaming && !displayContent)}
          active={Boolean(isStreaming)}
        />
      )}

      {terminalState?.visible && (
        <InteractiveTerminalBlock
          output={terminalState.output}
          active={terminalState.active || Boolean(isStreaming)}
        />
      )}

      {nonInteractiveToolMessages.length > 0 && (
        <div className="flex flex-col gap-0.5">
          {nonInteractiveToolMessages.length === 1 ? (
            <MessageItem message={nonInteractiveToolMessages[0]} />
          ) : (
            <ToolCallsGroup tools={nonInteractiveToolMessages} />
          )}
        </div>
      )}

      {displayContent && (
        <div className="llm-markdown select-text text-sm text-foreground leading-relaxed">
          {shouldRenderMarkdown ? (
            <ReactMarkdown
              remarkPlugins={[remarkGfm]}
              rehypePlugins={[rehypeHighlight]}
              components={{
                p: ({ className, children }) => (
                  <p className={cn('llm-markdown-paragraph', className)}>{children}</p>
                ),
                h1: ({ className, children }) => (
                  <h1 className={cn('llm-markdown-heading llm-markdown-h1', className)}>{children}</h1>
                ),
                h2: ({ className, children }) => (
                  <h2 className={cn('llm-markdown-heading llm-markdown-h2', className)}>{children}</h2>
                ),
                h3: ({ className, children }) => (
                  <h3 className={cn('llm-markdown-heading llm-markdown-h3', className)}>{children}</h3>
                ),
                h4: ({ className, children }) => (
                  <h4 className={cn('llm-markdown-heading llm-markdown-h4', className)}>{children}</h4>
                ),
                ul: ({ className, children }) => (
                  <ul className={cn('llm-markdown-list llm-markdown-list-disc', className)}>{children}</ul>
                ),
                ol: ({ className, children }) => (
                  <ol className={cn('llm-markdown-list llm-markdown-list-decimal', className)}>{children}</ol>
                ),
                li: ({ className, children }) => (
                  <li className={cn('llm-markdown-list-item', className)}>{children}</li>
                ),
                blockquote: ({ className, children }) => (
                  <blockquote className={cn('llm-markdown-blockquote', className)}>{children}</blockquote>
                ),
                hr: ({ className }) => <hr className={cn('llm-markdown-divider', className)} />,
                a: ({ className, children, ...props }) => (
                  <a
                    className={cn('llm-markdown-link', className)}
                    target="_blank"
                    rel="noreferrer"
                    {...props}
                  >
                    {children}
                  </a>
                ),
                table: ({ className, children }) => (
                  <div className="llm-markdown-table">
                    <table className={cn(className)}>{children}</table>
                  </div>
                ),
                thead: ({ className, children }) => (
                  <thead className={cn('llm-markdown-thead', className)}>{children}</thead>
                ),
                th: ({ className, children }) => (
                  <th className={cn('llm-markdown-th', className)}>{children}</th>
                ),
                td: ({ className, children }) => (
                  <td className={cn('llm-markdown-td', className)}>{children}</td>
                ),
                pre: ({ children }) => (
                  <pre className="llm-markdown-pre">{children}</pre>
                ),
                code: ({ children, className }) => {
                  const isBlock = className?.includes('language-')
                  return isBlock ? (
                    <code className={className}>{children}</code>
                  ) : (
                    <code className="llm-inline-code">{children}</code>
                  )
                },
              }}
            >
              {displayContent}
            </ReactMarkdown>
          ) : (
            <div className="whitespace-pre-wrap break-words">
              {displayContent}
            </div>
          )}
        </div>
      )}

      {pendingConfirm && onConfirm && onReject && onSubmitComment && (
        <ToolConfirmCard
          toolName={pendingConfirm.toolName}
          purpose={pendingConfirm.purpose}
          args={pendingConfirm.args}
          onConfirm={onConfirm}
          onReject={onReject}
          onSubmitComment={onSubmitComment}
          resolved={false}
          approved={false}
        />
      )}

      {!isStreaming && (
        <ResponseMeta modelName={modelName} tokensIn={tokensIn} tokensOut={tokensOut} />
      )}
    </div>
  )
})

MessageItem.displayName = 'MessageItem'

function flattenToolMessages(toolMessages: DisplayMessage[]): DisplayMessage[] {
  return toolMessages.flatMap((message) => message.isToolGroup && message.groupedTools ? message.groupedTools : [message])
}

function latestInteractiveTerminalState(toolMessages: DisplayMessage[]): InteractiveTerminalState | null {
  let latest: InteractiveTerminalState | null = null
  for (const toolMessage of toolMessages) {
    const state = getInteractiveTerminalState(toolMessage.toolResult ?? '')
    if (state) latest = state
  }
  return latest
}
