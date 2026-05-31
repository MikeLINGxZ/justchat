import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { AddProviderStepForm } from '@/components/settings/providers/AddProviderStepForm'
import { ProviderDetailView } from '@/components/settings/providers/ProviderDetailView'
import { ProviderList } from '@/components/settings/providers/ProviderList'

const deleteProviderMock = vi.fn().mockResolvedValue(undefined)
const editProviderMock = vi.fn().mockResolvedValue(undefined)

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider', () => ({
  Provider: {
    DeleteProvider: (...args: unknown[]) => deleteProviderMock(...args),
    EditProvider: (...args: unknown[]) => editProviderMock(...args),
    SetDefault: vi.fn().mockResolvedValue(undefined),
    DeleteModel: vi.fn().mockResolvedValue(undefined),
    RequestProviderModelList: vi.fn().mockResolvedValue({ models: [] }),
    CreateProvider: vi.fn().mockResolvedValue(undefined),
  },
}))

describe('ProviderList', () => {
  it('shows the default badge for the current default provider', () => {
    render(
      <ProviderList
        items={[
          {
            id: 1,
            provider_name: 'DeepSeek',
            provider_type: 'deepseek',
            base_url: 'https://api.deepseek.com',
            api_key: '',
            enabled: true,
            is_default: true,
            model_count: 0,
            icon: 'deepseek',
            models: [],
          },
        ]}
        selectedId={1}
        onSelect={vi.fn()}
        onCreate={vi.fn()}
        onToggleEnable={vi.fn()}
        onSetDefault={vi.fn()}
        onDelete={vi.fn()}
      />
    )

    expect(screen.getByText('默认')).toBeInTheDocument()
  })

  it('shows a confirmation dialog before deleting from the list', async () => {
    const user = userEvent.setup()
    const onDelete = vi.fn()

    render(
      <ProviderList
        items={[
          {
            id: 1,
            provider_name: 'A very very long provider name',
            provider_type: 'deepseek',
            base_url: 'https://example.com/really/long/provider/url',
            api_key: '',
            enabled: true,
            is_default: false,
            model_count: 0,
            icon: 'deepseek',
            models: [],
          },
        ]}
        selectedId={1}
        onSelect={vi.fn()}
        onCreate={vi.fn()}
        onToggleEnable={vi.fn()}
        onSetDefault={vi.fn()}
        onDelete={onDelete}
      />
    )

    expect(screen.getByText('A very very long provider name')).toHaveAttribute('title', 'A very very long provider name')
    expect(screen.getByText('https://example.com/really/long/provider/url')).toHaveAttribute('title', 'https://example.com/really/long/provider/url')

    await user.click(screen.getByRole('button', { name: '更多' }))
    await user.click(screen.getByRole('button', { name: '删除' }))

    expect(onDelete).not.toHaveBeenCalled()
    const dialog = screen.getByRole('dialog')
    expect(dialog).toBeInTheDocument()

    await user.click(within(dialog).getByRole('button', { name: '删除' }))

    expect(onDelete).toHaveBeenCalledWith(1)
  })
})

describe('ProviderDetailView', () => {
  it('shows a confirmation dialog before deleting from the detail page', async () => {
    const user = userEvent.setup()
    deleteProviderMock.mockClear()

    render(
      <ProviderDetailView
        provider={{
          id: 1,
          provider_name: 'DeepSeek',
          provider_type: 'deepseek',
          base_url: 'https://api.deepseek.com',
          api_key: '',
          enabled: true,
          is_default: false,
          model_count: 0,
          icon: 'deepseek',
          models: [],
        }}
        onUpdated={vi.fn()}
        onDeleted={vi.fn()}
      />
    )

    await user.click(screen.getByRole('button', { name: '删除' }))

    const dialog = screen.getByRole('dialog')
    expect(dialog).toBeInTheDocument()
    expect(deleteProviderMock).not.toHaveBeenCalled()

    await user.click(within(dialog).getByRole('button', { name: '删除' }))

    await waitFor(() => {
      expect(deleteProviderMock).toHaveBeenCalledWith({ provider_id: 1 })
    })
  })

  it('keeps a newly added custom model visible immediately', async () => {
    const user = userEvent.setup()
    editProviderMock.mockClear()

    render(
      <ProviderDetailView
        provider={{
          id: 1,
          provider_name: 'DeepSeek',
          provider_type: 'deepseek',
          base_url: 'https://api.deepseek.com',
          api_key: '',
          enabled: true,
          is_default: false,
          model_count: 1,
          icon: 'deepseek',
          models: [
            {
              id: 2,
              provider_id: 1,
              model: 'deepseek-chat',
              owned_by: '',
              object: 'model',
              enable: true,
              alias: null,
              is_custom: false,
              is_default: true,
            },
          ],
        }}
        onUpdated={vi.fn()}
        onDeleted={vi.fn()}
      />
    )

    await user.click(screen.getByRole('button', { name: '添加自定义模型' }))
    await user.type(screen.getByPlaceholderText('输入模型名称'), 'my-custom-model')
    await user.keyboard('{Enter}')

    expect(screen.getByText('my-custom-model')).toBeInTheDocument()
    expect(screen.getByText('自定义')).toBeInTheDocument()
  })

  it('syncs the model list when the parent passes refreshed provider data', () => {
    const { rerender } = render(
      <ProviderDetailView
        provider={{
          id: 1,
          provider_name: 'DeepSeek',
          provider_type: 'deepseek',
          base_url: 'https://api.deepseek.com',
          api_key: '',
          enabled: true,
          is_default: false,
          model_count: 1,
          icon: 'deepseek',
          models: [
            {
              id: 2,
              provider_id: 1,
              model: 'deepseek-chat',
              owned_by: '',
              object: 'model',
              enable: true,
              alias: null,
              is_custom: false,
              is_default: true,
            },
          ],
        }}
        onUpdated={vi.fn()}
        onDeleted={vi.fn()}
      />
    )

    expect(screen.getByText('deepseek-chat')).toBeInTheDocument()

    rerender(
      <ProviderDetailView
        provider={{
          id: 1,
          provider_name: 'DeepSeek',
          provider_type: 'deepseek',
          base_url: 'https://api.deepseek.com',
          api_key: '',
          enabled: true,
          is_default: false,
          model_count: 1,
          icon: 'deepseek',
          models: [
            {
              id: 3,
              provider_id: 1,
              model: 'deepseek-reasoner',
              owned_by: '',
              object: 'model',
              enable: true,
              alias: null,
              is_custom: true,
              is_default: true,
            },
          ],
        }}
        onUpdated={vi.fn()}
        onDeleted={vi.fn()}
      />
    )

    expect(screen.queryByText('deepseek-chat')).not.toBeInTheDocument()
    expect(screen.getByText('deepseek-reasoner')).toBeInTheDocument()
  })
})

describe('AddProviderStepForm', () => {
  it('keeps the custom tag visible for custom models', async () => {
    const user = userEvent.setup()

    render(
      <AddProviderStepForm
        provider={{
          type: 'deepseek',
          icon: 'deepseek',
          name: 'DeepSeek',
          description: '',
          base_url: 'https://api.deepseek.com',
        }}
        onDone={vi.fn()}
      />
    )

    await user.click(screen.getByRole('button', { name: '添加自定义模型' }))
    await user.type(screen.getByPlaceholderText('输入模型名称'), 'custom-model')
    await user.keyboard('{Enter}')

    expect(screen.getByText('custom-model')).toBeInTheDocument()
    expect(screen.getByText('自定义')).toBeInTheDocument()
  })
})
