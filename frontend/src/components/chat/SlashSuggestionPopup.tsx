import { forwardRef, useImperativeHandle, useState, useEffect } from 'react'

export type SlashSuggestionItem = {
  name: string
  description: string
}

type Props = {
  items: SlashSuggestionItem[]
  command: (item: SlashSuggestionItem) => void
}

export const SlashSuggestionPopup = forwardRef<unknown, Props>(({ items, command }, ref) => {
  const [index, setIndex] = useState(0)

  useEffect(() => {
    setIndex((i) => Math.min(i, Math.max(items.length - 1, 0)))
  }, [items])

  useImperativeHandle(ref, () => ({
    onKeyDown: ({ event }: { event: KeyboardEvent }) => {
      if (event.key === 'ArrowDown') {
        setIndex((i) => (i + 1) % items.length)
        return true
      }
      if (event.key === 'ArrowUp') {
        setIndex((i) => (i + items.length - 1) % items.length)
        return true
      }
      if (event.key === 'Enter') {
        const target = items[index]
        if (target) command(target)
        return true
      }
      return false
    },
  }), [items, index, command])

  if (items.length === 0) return null

  return (
    <div className="max-h-60 w-72 overflow-y-auto rounded-md border bg-popover p-1 shadow-lg">
      {items.map((it, i) => (
        <button
          key={it.name}
          type="button"
          onClick={() => command(it)}
          className={`flex w-full flex-col items-start gap-0.5 rounded-md px-2 py-1 text-left text-xs ${
            i === index ? 'bg-accent text-accent-foreground' : 'hover:bg-muted'
          }`}
        >
          <span className="font-medium">/{it.name}</span>
          <span className="line-clamp-2 opacity-70">{it.description}</span>
        </button>
      ))}
    </div>
  )
})
