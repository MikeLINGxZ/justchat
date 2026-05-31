import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { GeneralSettingsPanel } from '@/components/settings/general/GeneralSettingsPanel'
import { DisplaySettingsView } from '@/components/settings/general/DisplaySettingsView'
import { useAppStore } from '@/store/appStore'
import { getSettingsInitialState, useSettingsStore } from '@/store/settingsStore'

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/settings', () => ({
  Settings: {
    SaveDisplaySettings: vi.fn().mockResolvedValue(undefined),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  SelectFolder: vi.fn().mockResolvedValue(''),
}))

beforeEach(() => {
  useSettingsStore.setState(getSettingsInitialState())
  useAppStore.setState({
    theme: 'auto',
    fontSize: 'md',
    language: 'zh-CN',
  })
})

describe('DisplaySettingsView', () => {
  it('enables apply when font size draft changes', async () => {
    const user = userEvent.setup()
    const onApply = vi.fn()

    render(
      <DisplaySettingsView
        value="md"
        draft="md"
        onDraftChange={() => undefined}
        onApply={onApply}
        onReset={() => undefined}
      />
    )

    await user.click(screen.getByRole('button', { name: '大' }))

    expect(screen.getByRole('button', { name: '应用' })).toBeEnabled()
  })

  it('applies the font size draft when the user clicks apply', async () => {
    const user = userEvent.setup()

    useSettingsStore.setState({
      generalTab: 'display',
      fontSize: 'md',
      displayDraft: { fontSize: 'md' },
      displayDirty: false,
    })

    render(<GeneralSettingsPanel />)

    await user.click(screen.getByRole('button', { name: '大' }))
    await user.click(screen.getByRole('button', { name: '应用' }))

    await waitFor(() => {
      expect(useAppStore.getState().fontSize).toBe('lg')
      expect(useSettingsStore.getState().fontSize).toBe('lg')
      expect(useSettingsStore.getState().displayDirty).toBe(false)
    })
  })
})
