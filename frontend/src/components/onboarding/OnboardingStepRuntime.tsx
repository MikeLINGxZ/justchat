import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Events } from '@wailsio/runtime'
import { Onboarding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/onboarding/index'
import { Runtime } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/runtime/index'
import { cn } from '@/lib/utils'

type Phase = 'download' | 'extract' | 'verify'

type RuntimeViewState =
  | { kind: 'idle' }
  | { kind: 'downloading'; phase: Phase; received: number; total: number; percent: number }
  | { kind: 'ready'; version: string; installDir: string }
  | { kind: 'failed'; message: string }

type ProgressEventData = {
  phase: Phase
  received: number
  total: number
  percent: number
}

type Props = {
  onBack: () => void
}

export function OnboardingStepRuntime(props: Props) {
  const { t } = useTranslation()
  const [state, setState] = useState<RuntimeViewState>({ kind: 'idle' })

  useEffect(() => {
    let cancelled = false

    Runtime.GetStatus({})
      .then((res) => {
        if (cancelled || !res) {
          return
        }
        if (res.state === 'ready') {
          setState({ kind: 'ready', version: res.version, installDir: res.install_dir })
          return
        }
        if (res.state === 'downloading') {
          setState({ kind: 'downloading', phase: 'download', received: 0, total: 0, percent: 0 })
          return
        }
        if (res.state === 'failed') {
          setState({ kind: 'failed', message: res.error_msg || '' })
          return
        }
        setState({ kind: 'idle' })
      })
      .catch(() => undefined)

    const off = Events.On('runtime.node.progress', (evt: { data: ProgressEventData }) => {
      const progress = evt.data
      setState((prev) => {
        if (prev.kind === 'ready' || prev.kind === 'failed') {
          return prev
        }
        return {
          kind: 'downloading',
          phase: progress.phase,
          received: progress.received,
          total: progress.total,
          percent: progress.percent,
        }
      })

      if (progress.phase === 'verify' && progress.percent === 100) {
        Runtime.GetStatus({})
          .then((res) => {
            if (!res) {
              return
            }
            if (res.state === 'ready') {
              setState({ kind: 'ready', version: res.version, installDir: res.install_dir })
            } else if (res.state === 'failed') {
              setState({ kind: 'failed', message: res.error_msg || '' })
            }
          })
          .catch(() => undefined)
      }
    })

    return () => {
      cancelled = true
      off()
    }
  }, [])

  const startDownload = () => {
    setState({ kind: 'downloading', phase: 'download', received: 0, total: 0, percent: 0 })
    void Runtime.DownloadNode({})
  }

  const downloadLater = async () => {
    await Runtime.MarkDownloadLater({})
    await Onboarding.EnterHome({})
  }

  const cancelDownload = () => {
    void Runtime.CancelDownload({})
    setState({ kind: 'idle' })
  }

  const enterHome = () => {
    void Onboarding.EnterHome({})
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-4">
      <div className="space-y-1">
        <h2 className="text-lg font-semibold">{t('onboarding.runtime.title')}</h2>
        <p className="text-xs text-muted-foreground">{t('onboarding.runtime.intro')}</p>
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto rounded-2xl border border-border bg-card p-6">
        {state.kind === 'idle' && (
          <div className="flex flex-col items-center justify-center gap-3 py-8 text-center">
            <div className="text-sm text-muted-foreground">{t('onboarding.runtime.stateMissing')}</div>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={downloadLater}
                className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
              >
                {t('onboarding.actions.downloadLater')}
              </button>
              <button
                type="button"
                onClick={startDownload}
                className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground"
              >
                {t('onboarding.actions.downloadNow')}
              </button>
            </div>
          </div>
        )}

        {state.kind === 'downloading' && (
          <div className="space-y-3">
            <div className="text-sm font-medium">{t(`onboarding.runtime.phase${capitalize(state.phase)}`)}</div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
              <div
                className={cn('h-full bg-primary transition-all')}
                style={{ width: `${state.percent}%` }}
              />
            </div>
            <div className="text-xs text-muted-foreground">
              {formatBytes(state.received)} / {state.total > 0 ? formatBytes(state.total) : '…'} ({state.percent}%)
            </div>
            <div className="pt-2">
              <button
                type="button"
                onClick={cancelDownload}
                className="rounded-xl border border-border px-4 py-1.5 text-xs font-medium text-foreground transition-colors hover:bg-accent"
              >
                {t('onboarding.actions.cancel')}
              </button>
            </div>
          </div>
        )}

        {state.kind === 'ready' && (
          <div className="flex flex-col items-center justify-center gap-3 py-8 text-center">
            <div className="text-sm font-semibold text-primary">{t('onboarding.runtime.stateReady')}</div>
            <div className="text-xs text-muted-foreground">{state.version}</div>
            <div className="text-xs text-muted-foreground/70">{state.installDir}</div>
          </div>
        )}

        {state.kind === 'failed' && (
          <div className="flex flex-col items-center justify-center gap-3 py-8 text-center">
            <div className="text-sm font-semibold text-destructive">{t('onboarding.runtime.stateFailed')}</div>
            <div className="text-xs text-muted-foreground">{state.message}</div>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={downloadLater}
                className="rounded-xl border border-border px-4 py-1.5 text-xs font-medium text-foreground hover:bg-accent"
              >
                {t('onboarding.actions.skip')}
              </button>
              <button
                type="button"
                onClick={startDownload}
                className="rounded-xl bg-primary px-4 py-1.5 text-xs font-medium text-primary-foreground"
              >
                {t('onboarding.actions.retry')}
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="flex shrink-0 items-center justify-between border-t border-border pt-4">
        <button
          type="button"
          disabled={state.kind === 'downloading'}
          onClick={props.onBack}
          className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent disabled:opacity-40"
        >
          {t('onboarding.actions.back')}
        </button>
        <button
          type="button"
          disabled={state.kind === 'downloading'}
          onClick={enterHome}
          className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground transition-opacity disabled:opacity-40"
        >
          {t('onboarding.actions.enterHome')}
        </button>
      </div>
    </div>
  )
}

function capitalize(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1)
}

function formatBytes(value: number): string {
  if (value <= 0) {
    return '0 B'
  }
  const units = ['B', 'KB', 'MB', 'GB']
  let unitIndex = 0
  let current = value
  while (current >= 1024 && unitIndex < units.length - 1) {
    current /= 1024
    unitIndex += 1
  }
  return `${current.toFixed(1)} ${units[unitIndex]}`
}
