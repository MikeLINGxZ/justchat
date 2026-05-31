import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AddProviderApp } from '@/components/settings/providers/AddProviderApp'

const createProviderMock = vi.fn()
const supportedProviderListMock = vi.fn()

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider', () => ({
  Provider: {
    CreateProvider: (...args: unknown[]) => createProviderMock(...args),
    RequestProviderModelList: vi.fn().mockResolvedValue({ models: [] }),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/config', () => ({
  Config: {
    SupportedProviderList: (...args: unknown[]) => supportedProviderListMock(...args),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window', () => ({
  Window: {
    CloseAddProvider: vi.fn().mockResolvedValue(undefined),
  },
}))

describe('AddProviderApp', () => {
  beforeEach(() => {
    createProviderMock.mockReset()
    supportedProviderListMock.mockReset()
    supportedProviderListMock.mockResolvedValue({
      supported_providers: [
        {
          type: 'deepseek',
          icon: '',
          name: 'DeepSeek',
          description: 'Provider description',
          base_url: 'https://api.deepseek.com',
        },
      ],
    })
  })

  it('shows an error alert when creating a provider fails', async () => {
    const user = userEvent.setup()
    createProviderMock.mockRejectedValueOnce({
      message: JSON.stringify({
        detail: '网络连接失败',
        msg: '创建提供商失败',
      }),
      cause: {
        detail: '网络连接失败',
        msg: '创建提供商失败',
      },
      kind: 'RuntimeError',
    })

    render(<AddProviderApp />)

    await user.click(await screen.findByRole('button', { name: 'Select DeepSeek' }))
    await user.click(screen.getByRole('button', { name: '下一步' }))
    await user.click(screen.getByRole('button', { name: '添加' }))

    await waitFor(() => {
      expect(screen.getByLabelText('Banner alerts')).toBeInTheDocument()
      expect(screen.getByRole('alert')).toHaveTextContent('创建提供商失败')
      expect(screen.queryByText('网络连接失败')).not.toBeInTheDocument()
      expect(screen.queryByText(/RuntimeError/)).not.toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: /show details|展开详情/i }))

    expect(screen.getByText('网络连接失败')).toBeInTheDocument()
  })
})
