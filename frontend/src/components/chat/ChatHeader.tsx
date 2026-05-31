import { useState } from 'react'
import { PanelLeftOpen, Pencil, Check, X, RefreshCw } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Provider as ProviderBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider'
import type { ProviderWrapper } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models'
import { Type as ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider/models'
import { useChatStore } from '@/store/chatStore'
import { useSettingsStore } from '@/store/settingsStore'
import { cn } from '@/lib/utils'
import { useShallow } from 'zustand/react/shallow'
import type { Message } from '@/types'
import type { ProviderItem } from '@/types/settings'
import { NotificationBell } from '@/components/notifications/NotificationBell'

const EMPTY_MESSAGES: Message[] = []

interface ChatHeaderProps {
  sidebarCollapsed: boolean
  onExpandSidebar: () => void
  showWindowControlSpace: boolean
  animateWindowControlSpace: boolean
}

function pickProviderModel(provider: ProviderItem, preferredModelName?: string) {
  if (preferredModelName) {
    const matchedModel = provider.models.find(
      (model) => model.enable && model.model === preferredModelName
    )
    if (matchedModel) return matchedModel
  }

  return (
    provider.models.find((model) => model.enable && model.is_default) ??
    provider.models.find((model) => model.enable) ??
    null
  )
}

function mapWrapperToProviderItems(wrappers: ProviderWrapper[]): ProviderItem[] {
  return wrappers.map((wrapper) => ({
    id: wrapper.providers.id,
    provider_name: wrapper.providers.provider_name,
    provider_type: wrapper.providers.provider_type,
    base_url: wrapper.providers.base_url,
    api_key: wrapper.providers.api_key,
    enabled: wrapper.providers.enabled,
    is_default: wrapper.providers.is_default,
    model_count: wrapper.models.length,
    icon: wrapper.providers.icon,
    models: wrapper.models.map((model) => ({
      id: model.id,
      provider_id: model.provider_id,
      model: model.model,
      owned_by: model.owned_by,
      object: model.object,
      enable: model.enable,
      alias: model.alias,
      is_custom: model.is_custom,
      is_default: model.is_default,
    })),
  }))
}

export function ChatHeader({
  sidebarCollapsed,
  onExpandSidebar,
  showWindowControlSpace,
  animateWindowControlSpace,
}: ChatHeaderProps) {
  const { t } = useTranslation()
  const {
    currentConversationId,
    renameConversation,
    generateTitleDraft,
    title,
    hasMessages,
    sessionRunConfig,
    messages,
  } = useChatStore(useShallow((state) => {
    const currentConversation = state.conversations.find(
      conversation => conversation.id === state.currentConversationId
    )
    const conversationMessages = state.currentConversationId
      ? (state.messages[state.currentConversationId] ?? EMPTY_MESSAGES)
      : EMPTY_MESSAGES

    return {
      currentConversationId: state.currentConversationId,
      renameConversation: state.renameConversation,
      generateTitleDraft: state.generateTitleDraft,
      title: currentConversation?.title,
      hasMessages: conversationMessages.length > 0,
      sessionRunConfig: state.currentConversationId
        ? state.sessionRunConfigs[state.currentConversationId]
        : undefined,
      messages: conversationMessages,
    }
  }))
  const providers = useSettingsStore((state) => state.providers)
  const [editing, setEditing] = useState(false)
  const [editValue, setEditValue] = useState('')
  const [isGeneratingTitle, setIsGeneratingTitle] = useState(false)

  const displayTitle = title ?? t('chat.newChat')

  const startEdit = () => {
    setEditValue(displayTitle)
    setEditing(true)
  }

  const saveEdit = () => {
    if (editValue.trim() && currentConversationId) {
      void renameConversation(currentConversationId, editValue.trim())
    }
    setEditing(false)
  }

  const resolveTitleConfigFromProviders = (providerList: ProviderItem[]) => {
    const latestModelName = [...messages]
      .reverse()
      .find((message) => message.modelName)?.modelName

    const enabledProviders = providerList.filter((provider) => provider.enabled)
    if (enabledProviders.length === 0) return null

    const matchingProvider = latestModelName
      ? enabledProviders.find((provider) =>
          provider.models.some((model) => model.enable && model.model === latestModelName)
        )
      : null

    const selectedProvider =
      matchingProvider ??
      enabledProviders.find((provider) => provider.is_default) ??
      enabledProviders[0]

    if (!selectedProvider) return null

    const selectedModel = pickProviderModel(selectedProvider, latestModelName)
    if (!selectedModel) return null

    return {
      baseUrl: selectedProvider.base_url,
      apiKey: selectedProvider.api_key,
      modelName: selectedModel.model,
      providerType: selectedProvider.provider_type as ProviderType,
    }
  }

  const fallbackTitleConfig = sessionRunConfig ?? resolveTitleConfigFromProviders(providers)

  const regenerateTitle = async () => {
    if (!currentConversationId || isGeneratingTitle) return

    setIsGeneratingTitle(true)
    try {
      let titleConfig = fallbackTitleConfig
      if (!titleConfig) {
        const providersResult = await ProviderBinding.ListProviders({})
        titleConfig = providersResult?.providers
          ? resolveTitleConfigFromProviders(mapWrapperToProviderItems(providersResult.providers))
          : null
      }

      if (!titleConfig) return

      const nextTitle = await generateTitleDraft(currentConversationId, titleConfig)
      if (nextTitle) {
        setEditValue(nextTitle)
      }
    } finally {
      setIsGeneratingTitle(false)
    }
  }

  return (
    <div className={cn(
      'flex items-center gap-2 px-4 py-3 border-b border-border h-14 shrink-0',
      animateWindowControlSpace && 'transition-[padding-left] ease-[cubic-bezier(0.22,1,0.36,1)] will-change-[padding-left]',
      sidebarCollapsed && showWindowControlSpace ? cn('pl-20', animateWindowControlSpace && 'duration-150') : animateWindowControlSpace && 'duration-300'
    )}>
      {sidebarCollapsed && (
        <button
          type="button"
          aria-label="展开侧边栏"
          onClick={onExpandSidebar}
          className="shrink-0 rounded-lg p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        >
          <PanelLeftOpen size={16} />
        </button>
      )}
      {editing ? (
        <>
          <input
            className="flex-1 bg-transparent border-b border-primary outline-none text-lg font-medium text-foreground"
            value={editValue}
            onChange={e => setEditValue(e.target.value)}
            onKeyDown={e => {
              if (e.key === 'Enter') saveEdit()
              if (e.key === 'Escape' && !isGeneratingTitle) setEditing(false)
            }}
            autoFocus
          />
          <button onClick={saveEdit} className="p-1 rounded hover:bg-accent text-muted-foreground hover:text-foreground">
            <Check size={14} />
          </button>
          <button
            onClick={() => void regenerateTitle()}
            disabled={isGeneratingTitle}
            aria-label={t('chat.regenerateTitle')}
            title={t('chat.regenerateTitle')}
            className="p-1 rounded hover:bg-accent text-muted-foreground hover:text-foreground disabled:cursor-not-allowed disabled:opacity-50"
          >
            <RefreshCw size={14} className={cn(isGeneratingTitle && 'animate-spin')} />
          </button>
          <button
            onClick={() => setEditing(false)}
            disabled={isGeneratingTitle}
            className="p-1 rounded hover:bg-accent text-muted-foreground hover:text-foreground disabled:cursor-not-allowed disabled:opacity-50"
          >
            <X size={14} />
          </button>
        </>
      ) : (
        <div className="flex max-w-full min-w-1 items-center gap-3">
          <h1 className="text-lg font-medium text-foreground truncate">{displayTitle}</h1>
          {hasMessages && (
            <button
              onClick={startEdit}
              className="shrink-0 p-1 rounded hover:bg-accent text-muted-foreground hover:text-foreground"
            >
              <Pencil size={14} />
            </button>
          )}
        </div>
      )}
      <div className="ml-auto flex shrink-0 items-center">
        <NotificationBell />
      </div>
    </div>
  )
}
