import { useState } from 'react'
import { AlertCircle, ChevronDown, ChevronUp } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { ResponseMeta } from './ResponseMeta'

interface Props {
  content: string
  modelName?: string
  tokensIn?: number
  tokensOut?: number
}

export function ErrorMessageBubble({ content, modelName, tokensIn, tokensOut }: Props) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)

  let msg = content
  let detail = ''
  try {
    const parsed = JSON.parse(content) as { msg?: string; detail?: string }
    msg = parsed.msg ?? content
    detail = parsed.detail ?? ''
  } catch {
    // fallback: display raw string
  }

  const localizedStreamErrors = new Set([
    '大模型响应出错',
    'Model response error',
  ])

  if (localizedStreamErrors.has(msg)) {
    msg = t('chat.modelResponseError')
  }

  return (
    <div className="flex max-w-[70%] flex-col gap-1">
      <div className="flex w-full items-center justify-start gap-1 text-left text-sm text-destructive">
        <span className="select-text">{msg}</span>
        {detail && (
          <button
            onClick={() => setExpanded(!expanded)}
            className="flex-shrink-0 flex items-center gap-0.5 text-destructive/70 hover:text-destructive transition-colors"
            title={expanded ? t('chat.hideErrorDetail') : t('chat.errorDetail')}
          >
            <AlertCircle size={14} />
            {expanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
          </button>
        )}
      </div>
      {expanded && detail && (
        <div className="w-full rounded-lg border border-border bg-muted/50 px-3 py-2 text-left font-mono text-xs text-muted-foreground whitespace-pre-wrap break-all select-text">
          {detail}
        </div>
      )}
      <ResponseMeta modelName={modelName} tokensIn={tokensIn} tokensOut={tokensOut} />
    </div>
  )
}
