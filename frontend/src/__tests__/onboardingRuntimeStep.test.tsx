import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { OnboardingStepRuntime } from '@/components/onboarding/OnboardingStepRuntime'

const {
  getStatusMock,
  markDownloadLaterMock,
  downloadNodeMock,
  cancelDownloadMock,
  enterHomeMock,
  onMock,
} = vi.hoisted(() => ({
  getStatusMock: vi.fn(),
  markDownloadLaterMock: vi.fn(),
  downloadNodeMock: vi.fn(),
  cancelDownloadMock: vi.fn(),
  enterHomeMock: vi.fn(),
  onMock: vi.fn(() => vi.fn()),
}))

vi.mock('@wailsio/runtime', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wailsio/runtime')>()
  return {
    ...actual,
    Events: {
      ...actual.Events,
      On: onMock,
    },
  }
})

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/runtime/index', () => ({
  Runtime: {
    GetStatus: (...args: unknown[]) => getStatusMock(...args),
    MarkDownloadLater: (...args: unknown[]) => markDownloadLaterMock(...args),
    DownloadNode: (...args: unknown[]) => downloadNodeMock(...args),
    CancelDownload: (...args: unknown[]) => cancelDownloadMock(...args),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/onboarding/index', () => ({
  Onboarding: {
    EnterHome: (...args: unknown[]) => enterHomeMock(...args),
  },
}))

describe('OnboardingStepRuntime', () => {
  beforeEach(() => {
    getStatusMock.mockReset()
    markDownloadLaterMock.mockReset()
    downloadNodeMock.mockReset()
    cancelDownloadMock.mockReset()
    enterHomeMock.mockReset()
    onMock.mockClear()

    getStatusMock.mockResolvedValue({
      state: 'missing',
      version: '',
      install_dir: '',
      node_path: '',
      npm_path: '',
      error_msg: '',
    })
    markDownloadLaterMock.mockResolvedValue(undefined)
    enterHomeMock.mockResolvedValue(undefined)
  })

  it('lets the user defer download and enter the app', async () => {
    const user = userEvent.setup()
    const onBack = vi.fn()

    render(<OnboardingStepRuntime onBack={onBack} />)

    await waitFor(() => {
      expect(screen.getByText('尚未安装')).toBeInTheDocument()
    })
    expect(screen.getByText('插件运行时')).toBeInTheDocument()
    expect(screen.getByText('安装后可以启用更多 CLI 与 MCP 插件能力。可以现在下载，也可以稍后在设置里完成。')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: '稍后下载' }))

    expect(markDownloadLaterMock).toHaveBeenCalledWith({})
    expect(enterHomeMock).toHaveBeenCalledWith({})
  })

  it('supports returning to the previous step when not downloading', async () => {
    const user = userEvent.setup()
    const onBack = vi.fn()

    render(<OnboardingStepRuntime onBack={onBack} />)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: '返回' })).toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: '返回' }))

    expect(onBack).toHaveBeenCalledTimes(1)
  })
})
