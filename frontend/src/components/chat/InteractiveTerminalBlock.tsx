import { useEffect, useRef } from 'react'
import { ScanLine, Terminal } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import 'xterm/css/xterm.css'

type TerminalCtor = {
  new (options: {
    convertEol: boolean
    cursorBlink: boolean
    fontFamily: string
    fontSize: number
    scrollback: number
    theme: {
      background: string
      foreground: string
      cursor: string
      selectionBackground: string
    }
  }): {
    open: (element: HTMLElement) => void
    loadAddon: (addon: unknown) => void
    write: (data: string) => void
    focus: () => void
    onData: (listener: (data: string) => void) => { dispose: () => void }
    reset: () => void
    dispose: () => void
  }
}

type FitAddonCtor = {
  new (): {
    fit: () => void
    readonly cols: number
    readonly rows: number
  }
}

type WebLinksAddonCtor = {
  new (): unknown
}

type Props = {
  output: string
  active?: boolean
  onInput?: (data: string) => void
  onResize?: (rows: number, cols: number) => void
}

export function InteractiveTerminalBlock({ output, active = false, onInput, onResize }: Props) {
  const { t } = useTranslation()
  const terminalHostRef = useRef<HTMLDivElement | null>(null)
  const terminalWriterRef = useRef<((data: string) => void) | null>(null)
  const lastOutputRef = useRef('')
  const pendingOutputRef = useRef(output)
  const activeRef = useRef(active)
  const onInputRef = useRef(onInput)
  const onResizeRef = useRef(onResize)

  useEffect(() => {
    activeRef.current = active
    onInputRef.current = onInput
    onResizeRef.current = onResize
  }, [active, onInput, onResize])

  useEffect(() => {
    const writer = terminalWriterRef.current
    if (!writer) {
      pendingOutputRef.current = output
      return
    }

    if (output.startsWith(lastOutputRef.current)) {
      const nextChunk = output.slice(lastOutputRef.current.length)
      if (nextChunk) {
        writer(nextChunk)
      }
    } else {
      writer('\x1bc')
      if (output) {
        writer(output)
      }
    }
    lastOutputRef.current = output
  }, [output])

  useEffect(() => {
    if (!terminalHostRef.current) return

    let disposed = false
    let cleanup: (() => void) | null = null

    const setupTerminal = async (): Promise<void> => {
      const [{ Terminal: XTerm }, { FitAddon }, { WebLinksAddon }] = await Promise.all([
        import('xterm'),
        import('xterm-addon-fit'),
        import('xterm-addon-web-links'),
      ])

      if (disposed || !terminalHostRef.current) return

      const terminal = new (XTerm as unknown as TerminalCtor)({
        convertEol: true,
        cursorBlink: true,
        fontFamily: 'Menlo, Consolas, monospace',
        fontSize: 13,
        scrollback: 5000,
        theme: {
          background: '#fffbeb',
          foreground: '#0f172a',
          cursor: '#92400e',
          selectionBackground: '#fbbf24',
        },
      })
      const fitAddon = new (FitAddon as unknown as FitAddonCtor)()
      const webLinksAddon = new (WebLinksAddon as unknown as WebLinksAddonCtor)()

      terminal.open(terminalHostRef.current)
      terminal.loadAddon(fitAddon)
      terminal.loadAddon(webLinksAddon)
      fitAddon.fit()
      terminalWriterRef.current = (data: string) => terminal.write(data)

      if (pendingOutputRef.current) {
        terminal.write(pendingOutputRef.current)
        lastOutputRef.current = pendingOutputRef.current
      }

      const emitResize = (): void => {
        onResizeRef.current?.(Math.max(fitAddon.rows, 1), Math.max(fitAddon.cols, 1))
      }

      emitResize()

      const dataSubscription = terminal.onData((data: string) => {
        if (activeRef.current) {
          onInputRef.current?.(data)
        }
      })

      let resizeObserver: ResizeObserver | null = null
      if (typeof ResizeObserver !== 'undefined' && terminalHostRef.current) {
        resizeObserver = new ResizeObserver(() => {
          fitAddon.fit()
          emitResize()
        })
        resizeObserver.observe(terminalHostRef.current)
      }

      if (activeRef.current) {
        terminal.focus()
      }

      cleanup = () => {
        terminalWriterRef.current = null
        dataSubscription.dispose()
        resizeObserver?.disconnect()
        terminal.dispose()
      }
    }

    void setupTerminal()

    return () => {
      disposed = true
      cleanup?.()
    }
  }, [])

  return (
    <div className="my-1 overflow-hidden rounded-lg border border-amber-300 bg-amber-50/80 shadow-sm dark:border-amber-500/50 dark:bg-amber-950/20">
      <div className="flex items-center gap-2 border-b border-amber-200 px-3 py-2 text-xs text-amber-900 dark:border-amber-500/40 dark:text-amber-100">
        <ScanLine size={14} />
        <span className="font-medium">
          {t(active ? 'chat.interactiveTerminalActive' : 'chat.interactiveTerminal')}
        </span>
        <Terminal size={12} className="ml-auto opacity-70" />
      </div>
      <span className="sr-only">{stripAnsiForAccessibleText(output)}</span>
      <div
        ref={terminalHostRef}
        role="textbox"
        aria-label={t('chat.interactiveTerminalInput')}
        aria-readonly={!active}
        className="h-[280px] max-h-[420px] w-full overflow-hidden px-3 py-2 [&_.xterm]:h-full [&_.xterm-viewport]:bg-transparent [&_.xterm-screen]:bg-transparent"
      />
    </div>
  )
}

function stripAnsiForAccessibleText(output: string): string {
  return output.replace(/\x1b\[[0-9;?]*[ -/]*[@-~]/g, '')
}
