import { fireEvent, render, screen } from '@testing-library/react'
import { useContext } from 'react'
import { beforeEach, describe, expect, it } from 'vitest'
import i18n from '@/i18n'
import { SettingsPrimaryMenu } from '@/components/settings/SettingsPrimaryMenu'
import { SettingsMenuContext, SettingsShell } from '@/components/settings/SettingsShell'
import { getSettingsInitialState, useSettingsStore } from '@/store/settingsStore'

function OpenMenuConsumer() {
  const { onOpenMenu } = useContext(SettingsMenuContext)
  if (!onOpenMenu) return null
  return (
    <button type="button" aria-label="Open settings menu" onClick={onOpenMenu}>
      open
    </button>
  )
}

beforeEach(() => {
  useSettingsStore.setState(getSettingsInitialState())
  void i18n.changeLanguage('zh-CN')
})

describe('SettingsShell', () => {
  it('collapses the primary settings menu on narrow windows and allows reopening it', () => {
    window.innerWidth = 760

    render(
      <SettingsShell>
        <OpenMenuConsumer />
      </SettingsShell>
    )

    fireEvent(window, new Event('resize'))

    expect(screen.queryByText('Lemontea')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Open settings menu' })).toBeInTheDocument()
  })

  it('shows the primary menu as a full-screen h5 view after reopening on narrow windows', () => {
    window.innerWidth = 760

    render(
      <SettingsShell>
        <OpenMenuConsumer />
      </SettingsShell>
    )

    fireEvent(window, new Event('resize'))
    fireEvent.click(screen.getByRole('button', { name: 'Open settings menu' }))

    expect(screen.getByText('Lemontea')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Open settings menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Close settings menu' })).not.toBeInTheDocument()
  })

  it('renders translated settings labels when the app language changes', async () => {
    await i18n.changeLanguage('en')

    render(
      <SettingsPrimaryMenu
        activeTab="general"
        onChange={() => undefined}
      />
    )

    expect(screen.getByText('General')).toBeInTheDocument()
    expect(screen.getByText('Providers')).toBeInTheDocument()
    expect(screen.getByText('About')).toBeInTheDocument()
  })
})
