import { useRef, useState } from 'react'
import { RefreshCw, Plus, Check, X, Search } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Onboarding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/onboarding/index'
import { Provider as ProviderBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider'
import { CreateProviderInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models'
import { Model as ViewModel } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model/models'
import { Type as ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider/models'
import { useAlertStore } from '@/alert/store'
import { parseUnknownIError } from '@/lib/ierror'
import { cn } from '@/lib/utils'
import type { SupportedProvider } from '@/types/settings'

type FetchedModel = { model: string; owned_by: string; object: string; source: 'fetched' | 'custom' }

export function AddProviderStepForm(props: {
  provider: SupportedProvider
  onDone: () => void
  mode?: 'settings' | 'onboarding'
}) {
  const { t } = useTranslation()
  const pushAlert = useAlertStore((state) => state.pushAlert)
  const [enable, setEnable] = useState(true)
  const [name, setName] = useState(props.provider.name)
  const [apiKey, setApiKey] = useState('')
  const [baseUrl, setBaseUrl] = useState(props.provider.base_url)
  const [defaultModel, setDefaultModel] = useState('')
  const [models, setModels] = useState<FetchedModel[]>([])
  const [refreshing, setRefreshing] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [addingCustom, setAddingCustom] = useState(false)
  const [customModelInput, setCustomModelInput] = useState('')
  const customInputRef = useRef<HTMLInputElement>(null)
  const customModelsRef = useRef<FetchedModel[]>([])
  const searchInputRef = useRef<HTMLInputElement>(null)
  const [searchOpen, setSearchOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  const handleRefreshModels = async () => {
    if (!baseUrl || refreshing) return
    setRefreshing(true)
    try {
      const result = await ProviderBinding.RequestProviderModelList({ base_url: baseUrl, api_key: apiKey })
      if (result?.models) {
        const fetched = result.models.map((m) => ({ model: m.model, owned_by: m.owned_by, object: m.object, source: 'fetched' as const }))
        const existing = new Map<string, FetchedModel>()
        for (const m of customModelsRef.current) existing.set(m.model, m)
        for (const m of fetched) {
          if (!existing.has(m.model)) existing.set(m.model, m)
        }
        const merged = Array.from(existing.values())
        setModels(merged)
        if (merged.length > 0 && !defaultModel) setDefaultModel(merged[0].model)
      }
    } catch (err) {
      const parsed = parseUnknownIError(err)
      pushAlert({
        kind: 'error',
        placement: 'banner',
        title: '',
        message: parsed.msg || t('settingsPage.addProvider.errors.fetchModels'),
        detail: parsed.detail !== parsed.msg ? parsed.detail : null,
      })
    } finally {
      setRefreshing(false)
    }
  }

  const handleConfirmAddCustom = (keepOpen = false) => {
    const trimmed = customModelInput.trim()
    if (!trimmed || models.some((m) => m.model === trimmed)) {
      if (keepOpen) {
        setCustomModelInput('')
        requestAnimationFrame(() => customInputRef.current?.focus())
      } else {
        setAddingCustom(false)
        setCustomModelInput('')
      }
      return
    }
    const custom: FetchedModel = { model: trimmed, owned_by: '', object: 'model', source: 'custom' }
    customModelsRef.current = [...customModelsRef.current, custom]
    setModels((prev) => [...prev, custom])
    if (!defaultModel) setDefaultModel(trimmed)
    setCustomModelInput('')
    if (keepOpen) {
      requestAnimationFrame(() => customInputRef.current?.focus())
    } else {
      setAddingCustom(false)
    }
  }

  const handleStartAddCustom = () => {
    if (addingCustom) {
      if (customModelInput.trim() === '') return
      handleConfirmAddCustom(true)
      return
    }
    setAddingCustom(true)
    setCustomModelInput('')
    requestAnimationFrame(() => customInputRef.current?.focus())
  }

  const handleCancelAddCustom = () => {
    setAddingCustom(false)
    setCustomModelInput('')
  }

  const handleDeleteCustomModel = (modelName: string) => {
    customModelsRef.current = customModelsRef.current.filter((m) => m.model !== modelName)
    setModels((prev) => {
      const next = prev.filter((m) => !(m.source === 'custom' && m.model === modelName))
      if (defaultModel === modelName && next.length > 0) setDefaultModel(next[0].model)
      else if (defaultModel === modelName) setDefaultModel('')
      return next
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (submitting) return
    setSubmitting(true)

    try {
      const allModels: ViewModel[] = models.map((m) => new ViewModel({
        id: 0,
        provider_id: 0,
        model: m.model,
        owned_by: m.owned_by,
        object: m.object,
        enable: true,
        alias: null,
        is_custom: m.source === 'custom',
        is_default: false,
      }))

      const input = new CreateProviderInput({
        provider_name: name,
        provider_type: props.provider.type as ProviderType,
        base_url: baseUrl,
        api_key: apiKey,
        enable,
        default_model: defaultModel || null,
        models: allModels,
      })

      if (props.mode === 'onboarding') {
        await Onboarding.SaveProviderAndMarkInitialized({
          provider_name: input.provider_name,
          provider_type: input.provider_type,
          base_url: input.base_url,
          api_key: input.api_key,
          enable: input.enable,
          default_model: input.default_model,
          models: input.models,
        })
      } else {
        await ProviderBinding.CreateProvider(input)
      }
      props.onDone()
    } catch (err) {
      const parsed = parseUnknownIError(err)
      pushAlert({
        kind: 'error',
        placement: 'banner',
        title: '',
        message: parsed.msg || t('settingsPage.addProvider.errors.createProvider'),
        detail: parsed.detail !== parsed.msg ? parsed.detail : null,
      })
    } finally {
      setSubmitting(false)
    }
  }

  const filteredModels = searchOpen && searchQuery
    ? models.filter(m => m.model.toLowerCase().includes(searchQuery.toLowerCase()))
    : models

  const openSearch = () => {
    setSearchOpen(true)
    requestAnimationFrame(() => searchInputRef.current?.focus())
  }

  const closeSearch = () => {
    setSearchOpen(false)
    setSearchQuery('')
  }

  return (
    <form id="add-provider-form" onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col space-y-5 pb-4 pt-2">
      {/* Enable toggle */}
      <div className="flex shrink-0 flex-col gap-2">
        <span className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.enable')}</span>
        <button
          type="button"
          role="switch"
          aria-checked={enable}
          onClick={() => setEnable((v) => !v)}
          className={cn(
            'relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors',
            enable ? 'bg-primary' : 'bg-muted'
          )}
        >
          <span className={cn(
            'inline-block h-4 w-4 rounded-full bg-white shadow transition-transform',
            enable ? 'translate-x-6' : 'translate-x-1'
          )} />
        </button>
      </div>

      {/* Provider name */}
      <div className="shrink-0 space-y-1.5">
        <label className="text-sm font-medium text-foreground">
          {t('settingsPage.addProvider.form.name')} <span className="text-destructive">*</span>
        </label>
        <input
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
        />
      </div>

      {/* Base URL */}
      <div className="shrink-0 space-y-1.5">
        <label className="text-sm font-medium text-foreground">
          {t('settingsPage.addProvider.form.baseUrl')} <span className="text-destructive">*</span>
        </label>
        <input
          required
          value={baseUrl}
          onChange={(e) => setBaseUrl(e.target.value)}
          className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
        />
      </div>

      {/* API key */}
      <div className="shrink-0 space-y-1.5">
        <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.apiKey')}</label>
        <input
          type="password"
          value={apiKey}
          onChange={(e) => setApiKey(e.target.value)}
          className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
        />
      </div>

      {/* Model list */}
      <div className="flex min-h-0 flex-1 flex-col space-y-1.5">
        <div className="flex shrink-0 items-center justify-between">
          <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.modelList')}</label>
          <div className="flex items-center gap-1">
            {/* Search box — expands from right to left */}
            <div className={cn(
              'flex items-center overflow-hidden rounded-lg border border-border bg-background px-2 transition-all duration-200',
              searchOpen ? 'w-36 opacity-100' : 'w-0 border-0 opacity-0'
            )}>
              <input
                ref={searchInputRef}
                tabIndex={searchOpen ? 0 : -1}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder={t('settingsPage.providers.searchModelPlaceholder')}
                className="min-w-0 flex-1 bg-transparent py-1 text-xs outline-none"
              />
              <button
                type="button"
                onClick={closeSearch}
                className="ml-1 shrink-0 text-muted-foreground hover:text-foreground"
              >
                <X size={11} />
              </button>
            </div>
            {!searchOpen && (
              <button
                type="button"
                onClick={openSearch}
                title={t('settingsPage.providers.searchModelPlaceholder')}
                className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent"
              >
                <Search size={14} />
              </button>
            )}
            <button
              type="button"
              onClick={handleRefreshModels}
              disabled={refreshing}
              title={t('settingsPage.addProvider.form.refreshModels')}
              aria-label={t('settingsPage.addProvider.form.refreshModels')}
              className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent disabled:opacity-40"
            >
              <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
            </button>
            <button
              type="button"
              onClick={handleStartAddCustom}
              title={t('settingsPage.addProvider.form.addCustomModel')}
              aria-label={t('settingsPage.addProvider.form.addCustomModel')}
              className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent"
            >
              <Plus size={14} />
            </button>
          </div>
        </div>

        <div className="scrollbar-track-transparent min-h-0 flex-1 overflow-y-auto rounded-xl border border-border">
          {addingCustom && (
            <div className="flex items-center gap-1.5 border-b border-border px-3 py-2">
              <input
                ref={customInputRef}
                value={customModelInput}
                onChange={(e) => setCustomModelInput(e.target.value)}
                onKeyDown={(e) => { if (e.key === 'Enter') handleConfirmAddCustom(); else if (e.key === 'Escape') handleCancelAddCustom() }}
                placeholder={t('settingsPage.addProvider.form.customModelPlaceholder')}
                className="min-w-0 flex-1 bg-transparent text-sm outline-none"
              />
              <button
                type="button"
                onClick={() => handleConfirmAddCustom()}
                disabled={customModelInput.trim() === ''}
                className="flex h-6 w-6 items-center justify-center rounded text-primary hover:bg-accent disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent"
              >
                <Check size={14} />
              </button>
              <button type="button" onClick={handleCancelAddCustom} className="flex h-6 w-6 items-center justify-center rounded text-muted-foreground hover:bg-accent">
                <X size={14} />
              </button>
            </div>
          )}
          {models.length === 0 && !addingCustom && (
            <div className="flex items-center justify-center gap-1.5 px-3 py-4 text-xs text-muted-foreground">
              {refreshing ? (
                <RefreshCw size={12} className="animate-spin" />
              ) : (
                <>
                  {t('settingsPage.addProvider.form.noModels')}
                  <span
                    className="cursor-pointer underline hover:text-foreground"
                    onClick={handleRefreshModels}
                  >
                    {t('settingsPage.addProvider.form.refreshModels')}
                  </span>
                </>
              )}
            </div>
          )}
          {models.length > 0 && filteredModels.length === 0 && searchQuery && (
            <div className="flex items-center justify-center px-3 py-4 text-xs text-muted-foreground">
              {t('settingsPage.providers.noSearchResults')}
            </div>
          )}
          {filteredModels.map((m, idx) => (
            <div
              key={`${m.source}-${m.model}-${idx}`}
              className="group flex items-center gap-2 px-3 py-2 hover:bg-accent"
            >
              <span className="min-w-0 flex-1 truncate text-sm">{m.model}</span>
              <span className="flex shrink-0 items-center gap-1.5">
                {defaultModel !== m.model && (
                  <button
                    type="button"
                    onClick={() => setDefaultModel(m.model)}
                    className="hidden rounded px-1.5 py-0.5 text-[10px] text-muted-foreground hover:bg-accent hover:text-foreground group-hover:inline-block"
                  >
                    {t('settingsPage.addProvider.form.setDefault')}
                  </button>
                )}
                {m.source === 'custom' && (
                  <button
                    type="button"
                    onClick={() => handleDeleteCustomModel(m.model)}
                    className="hidden rounded px-1.5 py-0.5 text-[10px] text-destructive/70 hover:bg-accent hover:text-destructive group-hover:inline-block"
                  >
                    {t('settingsPage.providers.delete')}
                  </button>
                )}
                {defaultModel === m.model && (
                  <span className="rounded bg-primary/10 px-1.5 py-0.5 text-[10px] text-primary">
                    {t('settingsPage.addProvider.form.defaultTag')}
                  </span>
                )}
                {m.source === 'custom' && (
                  <span className="rounded border border-border px-1.5 py-0.5 text-[10px] text-muted-foreground">
                    {t('settingsPage.addProvider.form.customTag')}
                  </span>
                )}
              </span>
            </div>
          ))}
        </div>
      </div>
    </form>
  )
}
