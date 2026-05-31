import { useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Events } from '@wailsio/runtime'
import { Plugin as PluginBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin'
import { ResizeLoginCliInput, SendLoginStdinInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto/models'
import 'xterm/css/xterm.css'

type TerminalCtor = {
  new (options: {
    convertEol: boolean
    fontFamily: string
    fontSize: number
    theme: {
      background: string
      foreground: string
    }
  }): {
    open: (element: HTMLElement) => void
    loadAddon: (addon: unknown) => void
    write: (data: string) => void
    focus: () => void
    onData: (listener: (data: string) => void) => { dispose: () => void }
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

type LoginOutputEvent = {
  data: {
    id: string
    data: string
  }
}

type LoginDoneEvent = {
  data: {
    id: string
    exit_code: number
    error: string
  }
}

type DoneState = {
  exitCode: number
  error: string
}

export type CliLoginDialogProps = {
  open: boolean
  extensionId: string
  title: string
  onStart: () => Promise<void>
  onClose: () => void
  onCancel: () => void
  onDone?: (exitCode: number) => void
  onStartError?: (error: unknown) => void
}

// CliLoginDialog renders a PTY-backed terminal modal for interactive CLI login flows.
export function CliLoginDialog(props: CliLoginDialogProps) {
  const { t } = useTranslation()
  const terminalHostRef = useRef<HTMLDivElement | null>(null)
  const onStartRef = useRef(props.onStart)
  const onCloseRef = useRef(props.onClose)
  const onDoneRef = useRef(props.onDone)
  const onStartErrorRef = useRef(props.onStartError)
  const [doneState, setDoneState] = useState<DoneState | null>(null)

  useEffect(() => {
    onStartRef.current = props.onStart
    onCloseRef.current = props.onClose
    onDoneRef.current = props.onDone
    onStartErrorRef.current = props.onStartError
  }, [props.onClose, props.onDone, props.onStart, props.onStartError])

  useEffect(() => {
    if (!props.open || !terminalHostRef.current) {
      return
    }

    let disposed = false
    let terminalDispose: (() => void) | null = null
    let doneTimer: number | null = null
    let terminalWriter: ((data: string) => void) | null = null
    const pendingOutput: string[] = []

    const offOutput = Events.On('cli.login.output', (event: LoginOutputEvent) => {
      if (event.data.id !== props.extensionId) {
        return
      }
      if (terminalWriter) {
        terminalWriter(event.data.data)
        return
      }
      pendingOutput.push(event.data.data)
    })

    const offDone = Events.On('cli.login.done', (event: LoginDoneEvent) => {
      if (event.data.id !== props.extensionId) {
        return
      }
      setDoneState({
        exitCode: event.data.exit_code,
        error: event.data.error,
      })
      onDoneRef.current?.(event.data.exit_code)
      doneTimer = window.setTimeout(() => {
        onCloseRef.current()
      }, 3000)
    })

    const setupTerminal = async (): Promise<void> => {
      const [{ Terminal }, { FitAddon }, { WebLinksAddon }] = await Promise.all([
        import('xterm'),
        import('xterm-addon-fit'),
        import('xterm-addon-web-links'),
      ])

      if (disposed || !terminalHostRef.current) {
        return
      }

      const terminal = new (Terminal as unknown as TerminalCtor)({
        convertEol: true,
        fontFamily: 'Menlo, Consolas, monospace',
        fontSize: 13,
        theme: {
          background: '#1b2636',
          foreground: '#f8fafc',
        },
      })
      const fitAddon = new (FitAddon as unknown as FitAddonCtor)()
      const webLinksAddon = new (WebLinksAddon as unknown as WebLinksAddonCtor)()

      terminal.open(terminalHostRef.current)
      terminal.loadAddon(fitAddon)
      terminal.loadAddon(webLinksAddon)
      fitAddon.fit()
      terminal.focus()
      terminalWriter = (data: string) => terminal.write(data)

      for (const chunk of pendingOutput.splice(0)) {
        terminal.write(chunk)
      }

      const emitResize = (): void => {
        void Promise.resolve(
          PluginBinding.ResizeLoginCli(new ResizeLoginCliInput({
            id: props.extensionId,
            rows: Math.max(fitAddon.rows, 1),
            cols: Math.max(fitAddon.cols, 1),
          })),
        ).catch(() => undefined)
      }

      emitResize()

      const dataSubscription = terminal.onData((data: string) => {
        void Promise.resolve(
          PluginBinding.SendLoginStdin(new SendLoginStdinInput({
            id: props.extensionId,
            data,
          })),
        ).catch(() => undefined)
      })

      let resizeObserver: ResizeObserver | null = null
      if (typeof ResizeObserver !== 'undefined' && terminalHostRef.current) {
        resizeObserver = new ResizeObserver(() => {
          fitAddon.fit()
          emitResize()
        })
        resizeObserver.observe(terminalHostRef.current)
      }

      await onStartRef.current().catch((error: unknown) => {
        onStartErrorRef.current?.(error)
        onCloseRef.current()
      })

      terminalDispose = () => {
        if (doneTimer !== null) {
          window.clearTimeout(doneTimer)
        }
        terminalWriter = null
        dataSubscription.dispose()
        resizeObserver?.disconnect()
        terminal.dispose()
      }
    }

    void setupTerminal()

    return () => {
      disposed = true
      offOutput()
      offDone()
      terminalDispose?.()
    }
  }, [props.extensionId, props.open])

  useEffect(() => {
    if (!props.open) {
      setDoneState(null)
    }
  }, [props.open])

  if (!props.open) {
    return null
  }

  const doneSuccess = doneState?.exitCode === 0

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        aria-label={t('settingsPage.providers.cancel')}
        className="absolute inset-0 bg-black/45"
        onClick={props.onClose}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="cli-login-dialog-title"
        className="relative z-10 flex h-[min(80vh,720px)] w-full max-w-4xl flex-col overflow-hidden rounded-3xl border border-border bg-background shadow-2xl"
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h3 id="cli-login-dialog-title" className="text-base font-semibold text-foreground">
            {props.title}
          </h3>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={props.onCancel}
              className="rounded-xl border border-border px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-accent"
            >
              {t('settingsPage.plugins.cliLoginCancel')}
            </button>
            <button
              type="button"
              onClick={props.onClose}
              className="rounded-xl border border-border px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-accent"
            >
              {t('settingsPage.plugins.cliLoginClose')}
            </button>
          </div>
        </div>
        {doneState ? (
          <div className={doneSuccess
            ? 'border-b border-emerald-200 bg-emerald-50 px-5 py-3 text-sm text-emerald-700 dark:border-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-300'
            : 'border-b border-red-200 bg-red-50 px-5 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-950/40 dark:text-red-300'}
          >
            {doneSuccess
              ? t('settingsPage.plugins.cliLoginDoneSuccess')
              : t('settingsPage.plugins.cliLoginDoneFailure', {
                code: doneState.exitCode,
                detail: doneState.error || t('settingsPage.plugins.cliLoginDoneFailureUnknown'),
              })}
          </div>
        ) : null}
        <div className="min-h-0 flex-1 bg-slate-950 px-3 py-3">
          <div ref={terminalHostRef} className="h-full w-full overflow-hidden rounded-2xl" />
        </div>
      </div>
    </div>
  )
}
