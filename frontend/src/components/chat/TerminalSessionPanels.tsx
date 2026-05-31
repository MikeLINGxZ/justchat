import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Events } from '@wailsio/runtime'
import { Terminal as TerminalBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/terminal'
import type { Info } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal/models'
import { InteractiveTerminalBlock } from '@/components/chat/InteractiveTerminalBlock'

type TerminalView = Info & {
  output: string
  cursor: number
}

type OutputEvent = {
  terminal_id: string
  session_id: number
  cursor: number
  data: string
  status: string
  visible?: boolean
}

type ExitEvent = {
  terminal_id: string
  session_id: number
  status: string
  exit_code: number
  visible?: boolean
}

type Props = {
  sessionId: number | null
}

export function TerminalSessionPanels({ sessionId }: Props) {
  const { t } = useTranslation()
  const [terminals, setTerminals] = useState<Record<string, TerminalView>>({})

  useEffect(() => {
    if (!sessionId) {
      setTerminals({})
      return
    }

    let cancelled = false
    void TerminalBinding.ListTerminals({ session_id: sessionId }).then(async (result) => {
      if (cancelled) return
      const visible = (result?.items ?? []).filter((item) => item.visible)
      const entries = await Promise.all(visible.map(async (item) => {
        const output = await TerminalBinding.ReadTerminalOutput({
          terminal_id: item.id,
          cursor: 0,
        })
        const text = (output?.chunks ?? []).map((chunk) => chunk.data).join('')
        const cursor = (output?.chunks ?? []).reduce((max, chunk) => Math.max(max, chunk.cursor_end), 0)
        return [item.id, { ...item, output: text, cursor }] as const
      }))
      if (!cancelled) {
        setTerminals(Object.fromEntries(entries))
      }
    })

    return () => {
      cancelled = true
    }
  }, [sessionId])

  useEffect(() => {
    if (!sessionId) return

    const offOutput = Events.On('terminal.output', (event: { data: OutputEvent | OutputEvent[] }) => {
      const items = Array.isArray(event.data) ? event.data : [event.data]
      setTerminals((previous) => {
        let next = previous
        for (const item of items) {
          if (!item || item.session_id !== sessionId) continue
          const current = next[item.terminal_id]
          if (!current) {
            if (item.visible !== true) continue
            next = {
              ...next,
              [item.terminal_id]: {
                id: item.terminal_id,
                session_id: item.session_id,
                title: t('chat.interactiveTerminal'),
                command: '',
                args: '',
                cwd: '',
                status: item.status || 'active',
                visible: item.visible,
                pid: 0,
                current_cursor: item.cursor,
                output: item.data,
                cursor: item.cursor,
              },
            }
            continue
          }
          if (item.cursor <= current.cursor) continue
          next = {
            ...next,
            [item.terminal_id]: {
              ...current,
              output: current.output + item.data,
              cursor: item.cursor,
              status: item.status || current.status,
              visible: item.visible ?? current.visible,
            },
          }
        }
        return next
      })
    })

    const offExit = Events.On('terminal.exited', (event: { data: ExitEvent | ExitEvent[] }) => {
      const items = Array.isArray(event.data) ? event.data : [event.data]
      setTerminals((previous) => {
        let next = previous
        for (const item of items) {
          if (!item || item.session_id !== sessionId || !next[item.terminal_id]) continue
          next = {
            ...next,
            [item.terminal_id]: {
              ...next[item.terminal_id],
              status: item.status,
              exit_code: item.exit_code,
              visible: item.visible ?? next[item.terminal_id].visible,
            },
          }
        }
        return next
      })
    })

    return () => {
      offOutput()
      offExit()
    }
  }, [sessionId, t])

  const ordered = useMemo(
    () => Object.values(terminals).filter((item) => item.visible),
    [terminals]
  )

  if (ordered.length === 0) return null

  return (
    <div className="flex flex-col gap-2">
      {ordered.map((terminal) => (
        <InteractiveTerminalBlock
          key={terminal.id}
          output={terminal.output}
          active={terminal.status === 'active' || terminal.status === 'waiting'}
          onInput={(data) => {
            void TerminalBinding.WriteTerminalInput({
              terminal_id: terminal.id,
              data,
            })
          }}
          onResize={(rows, cols) => {
            void TerminalBinding.ResizeTerminal({
              terminal_id: terminal.id,
              rows,
              cols,
            })
          }}
        />
      ))}
    </div>
  )
}
