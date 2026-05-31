import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { ConfirmDialog } from '@/components/settings/common/ConfirmDialog'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { CliToolsSection } from '@/components/settings/plugins/CliToolsSection'
import { ExtensionDetailInfo } from '@/components/settings/plugins/ExtensionDetailInfo'
import { getJsonError } from '@/components/settings/plugins/jsonError'
import type { ExtensionToolItem } from '@/types/settings'
import type { ExtensionItem } from '@/types/settings'

export function PluginToolDetailView(props: {
  extension: ExtensionItem | null
  configText: string
  loading: boolean
  onReload: () => void
  onSave: (configText: string) => void
  onLoginCli?: () => void
  onResetData?: () => void
  onRegenerateCliManifest?: () => void
}) {
  const { t } = useTranslation()
  const [draft, setDraft] = useState(props.configText)
  const [confirmResetOpen, setConfirmResetOpen] = useState(false)
  const [confirmRegenerateOpen, setConfirmRegenerateOpen] = useState(false)

  useEffect(() => {
    setDraft(props.configText)
  }, [props.configText, props.extension?.id])

  const dirty = useMemo(() => draft !== props.configText, [draft, props.configText])
  const jsonError = useMemo(() => getJsonError(draft), [draft])
  const jsonValid = jsonError.valid
  const errorLine = jsonError.valid ? null : jsonError.errorLine
  const cliTools = useMemo(() => deriveCliTools(props.extension, draft, props.configText), [draft, props.configText, props.extension])
  const cliManifest = useMemo(() => deriveCliManifest(props.extension, draft, props.configText), [draft, props.configText, props.extension])

  const kindBadge = props.extension ? (
    <span className={cn(
      'rounded-full px-2.5 py-0.5 text-xs font-medium',
      props.extension.kind === 'mcp'
        ? 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-400'
        : props.extension.kind === 'cli'
          ? 'bg-purple-100 text-purple-700 dark:bg-purple-500/20 dark:text-purple-400'
          : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-400'
    )}>
      {props.extension.kind === 'mcp'
        ? t('settingsPage.plugins.kindMcp')
        : props.extension.kind === 'cli'
          ? t('settingsPage.plugins.kindCli')
          : t('settingsPage.plugins.kindPlugin')}
    </span>
  ) : undefined

  if (!props.extension) {
    return (
      <SettingsContentLayout
        header={(
          <SettingsPanelHeader
            title={t('settingsPage.plugins.title')}
            description={t('settingsPage.plugins.emptyDescription')}
          />
        )}
      />
    )
  }

  return (
    <>
      <SettingsContentLayout
        header={(
          <SettingsPanelHeader
            title={props.extension.name}
            description={props.extension.description || props.extension.author || props.extension.version}
            badge={kindBadge}
          />
        )}
        footprint={(props.extension.kind === 'mcp' || props.extension.kind === 'cli') && props.extension.config_file_path ? (
          <SettingsActionBar
            primaryLabel={t('settingsPage.actions.apply')}
            secondaryLabel={t('settingsPage.actions.reset')}
            primaryDisabled={!dirty || props.loading || !jsonValid}
            secondaryDisabled={!dirty || props.loading}
            onPrimaryClick={() => props.onSave(draft)}
            onSecondaryClick={() => setDraft(props.configText)}
            reloadLabel={t('settingsPage.plugins.reload')}
            onReloadClick={props.onReload}
            regenerateLabel={props.extension.kind === 'cli' ? (props.loading ? t('settingsPage.plugins.cliRegenerating') : t('settingsPage.plugins.cliRegenerate')) : undefined}
            regenerateLoading={props.loading}
            regenerateDisabled={props.loading || !props.onRegenerateCliManifest}
            onRegenerateClick={props.extension.kind === 'cli' ? () => setConfirmRegenerateOpen(true) : undefined}
          />
        ) : undefined}
      >
        <ExtensionDetailInfo
          extension={props.extension}
          draft={draft}
          jsonValid={jsonValid}
          errorLine={errorLine}
          onDraftChange={setDraft}
        />
        {props.extension.kind === 'cli' && cliManifest ? (
          <section className="mt-6 space-y-3 border-t border-border pt-4">
            <div className="flex items-center justify-between gap-4">
              <div className="space-y-1">
                <h3 className="text-sm font-semibold text-foreground">
                  {t('settingsPage.plugins.cliLoginTitle')}
                </h3>
                <p className="text-xs text-muted-foreground">
                  {t('settingsPage.plugins.cliLoginHint')}
                </p>
              </div>
              <button
                type="button"
                disabled={
                  props.loading
                  || !props.onLoginCli
                  || cliManifest.login_command.length === 0
                  || dirty
                  || !jsonValid
                }
                onClick={() => props.onLoginCli?.()}
                className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {t('settingsPage.plugins.cliLoginAction')}
              </button>
            </div>
            <label className="flex flex-col gap-1 text-sm text-foreground">
              <span className="text-xs text-muted-foreground">
                {t('settingsPage.plugins.cliLoginCommandLabel')}
              </span>
              <input
                type="text"
                placeholder={t('settingsPage.plugins.cliLoginCommandPlaceholder')}
                value={cliManifest.login_command.join(' ')}
                onChange={(event) => {
                  setDraft(updateCliManifestDraftLoginCommand(
                    draft,
                    parseLoginCommandInput(event.target.value),
                  ))
                }}
                className="rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm focus:border-primary focus:outline-none disabled:opacity-60"
                disabled={props.loading || !jsonValid}
              />
              {cliManifest.login_command.length === 0 ? (
                <span className="text-xs text-muted-foreground">
                  {t('settingsPage.plugins.cliLoginCommandEmptyHint')}
                </span>
              ) : dirty ? (
                <span className="text-xs text-amber-600 dark:text-amber-400">
                  {t('settingsPage.plugins.cliLoginCommandDirtyHint')}
                </span>
              ) : null}
            </label>
            <label className="flex items-center gap-2 text-sm text-foreground">
              <input
                type="checkbox"
                checked={cliManifest.isolation === 'shared'}
                onChange={(event) => {
                  setDraft(updateCliManifestDraftIsolation(
                    draft,
                    event.target.checked ? 'shared' : 'isolated',
                  ))
                }}
                disabled={props.loading || !jsonValid}
              />
              <span>{t('settingsPage.plugins.cliUseSharedIsolation')}</span>
            </label>
          </section>
        ) : null}
        {props.extension.kind === 'cli' ? (
          <CliToolsSection
            tools={cliTools}
            busy={props.loading}
            jsonValid={jsonValid}
            onToggleEnabled={(toolName, enabled) => {
              setDraft(updateCliToolDraft(draft, toolName, (tool) => ({
                ...tool,
                enabled,
              })))
            }}
            onToggleRequiresConfirm={(toolName, requiresConfirm) => {
              setDraft(updateCliToolDraft(draft, toolName, (tool) => ({
                ...tool,
                requires_confirm: requiresConfirm,
              })))
            }}
          />
        ) : null}
        {props.extension.kind === 'cli' && props.onResetData ? (
          <div className="mt-6 flex gap-3 border-t border-border pt-4">
            <button
              type="button"
              onClick={() => setConfirmResetOpen(true)}
              className="rounded-lg border border-red-300 bg-red-50 px-4 py-2 text-sm text-red-700 hover:bg-red-100 dark:border-red-700 dark:bg-red-500/10 dark:text-red-400"
            >
              {t('settingsPage.plugins.cliResetData')}
            </button>
          </div>
        ) : null}
      </SettingsContentLayout>
      <ConfirmDialog
        open={confirmResetOpen}
        title={t('settingsPage.plugins.cliResetDataConfirmTitle')}
        description={props.extension.name}
        confirmLabel={t('settingsPage.plugins.cliResetData')}
        cancelLabel={t('settingsPage.providers.cancel')}
        confirmTone="danger"
        onConfirm={() => {
          if (props.onResetData) {
            props.onResetData()
          }
          setConfirmResetOpen(false)
        }}
        onCancel={() => setConfirmResetOpen(false)}
      />
      <ConfirmDialog
        open={confirmRegenerateOpen}
        title={t('settingsPage.plugins.cliRegenerateConfirmTitle')}
        description={t('settingsPage.plugins.cliRegenerateConfirmDescription', {
          name: props.extension.name,
        })}
        confirmLabel={t('settingsPage.plugins.cliRegenerate')}
        cancelLabel={t('settingsPage.providers.cancel')}
        busy={props.loading}
        onConfirm={() => {
          if (props.onRegenerateCliManifest) {
            props.onRegenerateCliManifest()
          }
          setConfirmRegenerateOpen(false)
        }}
        onCancel={() => setConfirmRegenerateOpen(false)}
      />
    </>
  )
}

type DraftManifestTool = {
  name: string
  description: string
  enabled: boolean
  requires_confirm: boolean
}

type CliIsolationMode = 'isolated' | 'shared'

type DraftCliManifest = {
  login_command: string[]
  isolation: CliIsolationMode
}

function deriveCliTools(extension: ExtensionItem | null, draft: string, fallbackConfigText: string): ExtensionToolItem[] {
  if (!extension || extension.kind !== 'cli') {
    return []
  }
  const parsedDraft = parseCliToolItems(draft)
  if (parsedDraft) {
    return parsedDraft
  }
  const parsedFallback = parseCliToolItems(fallbackConfigText)
  if (parsedFallback) {
    return parsedFallback
  }
  return extension.tools
}

function parseCliToolItems(configText: string): ExtensionToolItem[] | null {
  try {
    const parsed = JSON.parse(configText) as {
      tools?: Array<{
        name?: string
        description?: string
        enabled?: boolean
        requires_confirm?: boolean
      }>
    }
    if (!Array.isArray(parsed.tools)) {
      return []
    }
    return parsed.tools.map((tool) => ({
      tool_id: tool.name ?? '',
      server_id: '',
      name: tool.name ?? '',
      description: tool.description ?? '',
      enabled: Boolean(tool.enabled),
      requires_confirm: Boolean(tool.requires_confirm),
    }))
  } catch {
    return null
  }
}

function updateCliToolDraft(
  draft: string,
  toolName: string,
  updater: (tool: DraftManifestTool) => DraftManifestTool,
): string {
  try {
    const parsed = JSON.parse(draft) as { tools?: DraftManifestTool[] }
    if (!Array.isArray(parsed.tools)) {
      return draft
    }
    parsed.tools = parsed.tools.map((tool) => (tool.name === toolName ? updater(tool) : tool))
    return `${JSON.stringify(parsed, null, 2)}\n`
  } catch {
    return draft
  }
}

function deriveCliManifest(extension: ExtensionItem | null, draft: string, fallbackConfigText: string): DraftCliManifest | null {
  if (!extension || extension.kind !== 'cli') {
    return null
  }

  const parsedDraft = parseCliManifestDraft(draft)
  if (parsedDraft) {
    return parsedDraft
  }

  return parseCliManifestDraft(fallbackConfigText)
}

function parseCliManifestDraft(configText: string): DraftCliManifest | null {
  try {
    const parsed = JSON.parse(configText) as {
      login_command?: unknown
      isolation?: unknown
    }
    const loginCommand: string[] = Array.isArray(parsed.login_command)
      ? parsed.login_command.filter((part: unknown): part is string => typeof part === 'string')
      : []
    const isolation: CliIsolationMode = parsed.isolation === 'shared' ? 'shared' : 'isolated'

    return {
      login_command: loginCommand,
      isolation,
    }
  } catch {
    return null
  }
}

function updateCliManifestDraftIsolation(draft: string, isolation: CliIsolationMode): string {
  try {
    const parsed = JSON.parse(draft) as Record<string, unknown>
    parsed.isolation = isolation
    return `${JSON.stringify(parsed, null, 2)}\n`
  } catch {
    return draft
  }
}

function updateCliManifestDraftLoginCommand(draft: string, loginCommand: string[]): string {
  try {
    const parsed = JSON.parse(draft) as Record<string, unknown>
    parsed.login_command = loginCommand
    return `${JSON.stringify(parsed, null, 2)}\n`
  } catch {
    return draft
  }
}

function parseLoginCommandInput(value: string): string[] {
  return value.trim().split(/\s+/).filter((token) => token.length > 0)
}
