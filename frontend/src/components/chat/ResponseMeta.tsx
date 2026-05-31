import { ArrowDownToLine, ArrowUpFromLine } from 'lucide-react'

interface Props {
  modelName?: string
  tokensIn?: number
  tokensOut?: number
}

export function ResponseMeta({ modelName = '', tokensIn = 0, tokensOut = 0 }: Props) {
  if (!modelName && tokensIn <= 0 && tokensOut <= 0) {
    return null
  }

  return (
    <div className="mt-0.5 flex w-full items-center gap-3 text-xs text-muted-foreground">
      {modelName && <span className="opacity-70">{modelName}</span>}
      {tokensIn > 0 && (
        <span className="flex items-center gap-1">
          <ArrowUpFromLine size={10} />
          {tokensIn}
        </span>
      )}
      {tokensOut > 0 && (
        <span className="flex items-center gap-1">
          <ArrowDownToLine size={10} />
          {tokensOut}
        </span>
      )}
    </div>
  )
}
