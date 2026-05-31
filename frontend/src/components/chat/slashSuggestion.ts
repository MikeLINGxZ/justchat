import { ReactRenderer } from '@tiptap/react'
import tippy, { type Instance as TippyInstance } from 'tippy.js'
import type { SuggestionProps, SuggestionKeyDownProps } from '@tiptap/suggestion'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { SlashSuggestionPopup, type SlashSuggestionItem } from './SlashSuggestionPopup'

export function buildSlashSuggestion() {
  return {
    char: '/',
    startOfLine: false,
    items: async ({ query }: { query: string }) => {
      try {
        const list = await SkillsBinding.ListSkills({})
        const enabled = (list?.skills ?? []).filter((s) => !s.disabled)
        return enabled
          .filter((s) => s.name.toLowerCase().startsWith(query.toLowerCase()))
          .slice(0, 10)
          .map((s) => ({ name: s.name, description: s.description })) as SlashSuggestionItem[]
      } catch {
        return []
      }
    },
    render: () => {
      let component: ReactRenderer
      let popup: TippyInstance[] = []
      return {
        onStart: (props: SuggestionProps) => {
          component = new ReactRenderer(SlashSuggestionPopup, { props, editor: undefined as never })
          popup = tippy('body', {
            getReferenceClientRect: props.clientRect as () => DOMRect,
            appendTo: () => document.body,
            content: component.element,
            showOnCreate: true,
            interactive: true,
            trigger: 'manual',
            placement: 'top-start',
          })
        },
        onUpdate: (props: SuggestionProps) => {
          component.updateProps(props)
          popup[0]?.setProps({ getReferenceClientRect: props.clientRect as () => DOMRect })
        },
        onKeyDown: (props: SuggestionKeyDownProps) => (component.ref as { onKeyDown?: (p: unknown) => boolean })?.onKeyDown?.(props) ?? false,
        onExit: () => {
          popup[0]?.destroy()
          component.destroy()
        },
      }
    },
  }
}
