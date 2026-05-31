import { useEffect } from 'react'

const INPUT_ASSISTANCE_SELECTOR = 'input, textarea, [contenteditable="true"]'

function disableAssistanceForElement(element: Element): void {
  element.setAttribute('autocomplete', 'off')
  element.setAttribute('autocorrect', 'off')
  element.setAttribute('autocapitalize', 'off')
  element.setAttribute('spellcheck', 'false')
}

function disableInputAssistance(root: ParentNode): void {
  root.querySelectorAll(INPUT_ASSISTANCE_SELECTOR).forEach(disableAssistanceForElement)
}

// useDisableInputAssistance disables Apple/WebKit typing suggestions globally.
export function useDisableInputAssistance(): void {
  useEffect(() => {
    disableInputAssistance(document)
    document.documentElement.setAttribute('spellcheck', 'false')

    const observer = new MutationObserver((mutations: MutationRecord[]) => {
      for (const mutation of mutations) {
        for (const node of mutation.addedNodes) {
          if (!(node instanceof Element)) continue
          if (node.matches(INPUT_ASSISTANCE_SELECTOR)) {
            disableAssistanceForElement(node)
          }
          disableInputAssistance(node)
        }
      }
    })

    observer.observe(document.body, { childList: true, subtree: true })
    return () => observer.disconnect()
  }, [])
}
