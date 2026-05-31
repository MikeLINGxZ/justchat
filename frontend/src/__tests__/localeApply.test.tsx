import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { GeneralSettingsPanel } from '@/components/settings/general/GeneralSettingsPanel'
import { getSettingsInitialState, useSettingsStore } from '@/store/settingsStore'
import { useAppStore } from '@/store/appStore'

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/settings', () => ({
  Settings: {
    SaveLocaleSettings: vi.fn().mockResolvedValue(undefined),
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
  useSettingsStore.setState({
    generalTab: 'locale',
    locale: 'zh-CN',
    language: 'zh-CN',
    localeDraft: {
      locale: 'zh-CN',
      language: 'zh-CN',
    },
    languages: [
      { id: 'zh-CN', name: '简体中文' },
      { id: 'en', name: 'English' },
    ],
    regions: [
      { id: 'zh-CN', name: '中国', icon: '🇨🇳' },
      { id: 'en-US', name: '美国', icon: '🇺🇸' },
    ],
    localeDirty: false,
  })
})

describe('GeneralSettingsPanel locale apply', () => {
  it('applies language and locale changes when the user clicks apply', async () => {
    const user = userEvent.setup()

    render(<GeneralSettingsPanel />)

    await user.click(screen.getByRole('combobox', { name: '语言' }))
    await user.click(screen.getByRole('option', { name: 'English' }))
    await user.click(screen.getByRole('combobox', { name: '地区' }))
    await user.click(screen.getByRole('option', { name: '美国' }))
    await user.click(screen.getByRole('button', { name: '应用' }))

    await waitFor(() => {
      expect(useAppStore.getState().language).toBe('en')
      expect(useSettingsStore.getState().language).toBe('en')
      expect(useSettingsStore.getState().locale).toBe('en-US')
      expect(useSettingsStore.getState().localeDirty).toBe(false)
    })
  })
})
