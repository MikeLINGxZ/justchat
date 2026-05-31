import { useState } from 'react'
import { ChevronDown, Wrench } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { ToolCallBlock } from './ToolCallBlock'
import { InteractiveTerminalBlock } from '@/components/chat/InteractiveTerminalBlock'
import type { DisplayMessage } from '@/types'
import {
  getInteractiveTerminalState,
  isInteractiveTerminalControl,
  type InteractiveTerminalState,
} from '@/lib/terminalOutput'

type Props = {
  tools: DisplayMessage[]
}

export function ToolCallsGroup({ tools }: Props) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const terminalState = latestInteractiveTerminalState(tools)
  const collapsedTools = tools.filter((tool) => !isInteractiveTerminalControl(tool.toolResult ?? ''))

  return (
    <div className="my-0 flex flex-col gap-0">
      {terminalState?.visible && (
        <InteractiveTerminalBlock
          output={terminalState.output}
          active={terminalState.active}
        />
      )}

      {collapsedTools.length > 0 && (
        <div>
          <button
            onClick={() => setExpanded((v) => !v)}
            className="flex w-full min-w-0 items-center gap-1.5 rounded-md py-0.5 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            <Wrench size={14} className="shrink-0" />
            <span className="shrink-0 text-foreground">
              {t('toolCall.toolsUsed', { count: collapsedTools.length })}
            </span>
            <ChevronDown
              size={14}
              className={expanded ? 'shrink-0 rotate-180 transition-transform' : 'shrink-0 rotate-0 transition-transform'}
            />
            <span className="min-w-0 flex-1" />
          </button>
          {expanded && (
            <div className="ml-1.5 mt-1 space-y-1 border-l border-border pl-3">
              {collapsedTools.map((tool) => (
                <ToolCallBlock
                  key={tool.id}
                  toolName={(() => {
                    try {
                      const parsed = JSON.parse(tool.content) as { name?: string }
                      return parsed.name ?? ''
                    } catch {
                      return ''
                    }
                  })()}
                  purpose={tool.extra}
                  args={(() => {
                    try {
                      const parsed = JSON.parse(tool.content) as { args?: unknown }
                      return JSON.stringify(parsed.args ?? {}, null, 2)
                    } catch {
                      return tool.content
                    }
                  })()}
                  result={tool.toolResult ?? ''}
                  userComment={tool.toolConfirmComment}
                  status={
                    tool.toolConfirmAction === 'comment'
                      ? 'user_commented'
                      : tool.toolConfirmAction === 'reject'
                        ? 'user_rejected'
                        : 'completed'
                  }
                />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function latestInteractiveTerminalState(tools: DisplayMessage[]): InteractiveTerminalState | null {
  let latest: InteractiveTerminalState | null = null
  for (const tool of tools) {
    const state = getInteractiveTerminalState(tool.toolResult ?? '')
    if (state) latest = state
  }
  return latest
}
