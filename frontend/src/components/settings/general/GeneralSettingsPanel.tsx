import { Settings } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/settings'
import { File } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'
import { DisplaySettingsView } from '@/components/settings/general/DisplaySettingsView'
import { FileSettingsView } from '@/components/settings/general/FileSettingsView'
import { LocaleSettingsView } from '@/components/settings/general/LocaleSettingsView'
import { useAppStore } from '@/store/appStore'
import { useSettingsStore } from '@/store/settingsStore'

export function GeneralSettingsPanel() {
  const generalTab = useSettingsStore((state) => state.generalTab)
  const fontSize = useSettingsStore((state) => state.fontSize)
  const dataDir = useSettingsStore((state) => state.dataDir)
  const displayDraft = useSettingsStore((state) => state.displayDraft)
  const localeDraft = useSettingsStore((state) => state.localeDraft)
  const fileDraft = useSettingsStore((state) => state.fileDraft)
  const localeDirty = useSettingsStore((state) => state.localeDirty)
  const fileDirty = useSettingsStore((state) => state.fileDirty)
  const languages = useSettingsStore((state) => state.languages)
  const regions = useSettingsStore((state) => state.regions)
  const setDisplayDraft = useSettingsStore((state) => state.setDisplayDraft)
  const setLocaleDraft = useSettingsStore((state) => state.setLocaleDraft)
  const setFileDraft = useSettingsStore((state) => state.setFileDraft)
  const applyDisplaySettings = useSettingsStore((state) => state.applyDisplaySettings)
  const applyLocaleSettings = useSettingsStore((state) => state.applyLocaleSettings)
  const setFontSize = useAppStore((state) => state.setFontSize)
  const setLanguage = useAppStore((state) => state.setLanguage)

  if (generalTab === 'locale') {
    return (
      <LocaleSettingsView
        locale={localeDraft.locale}
        language={localeDraft.language}
        languages={languages}
        regions={regions}
        dirty={localeDirty}
        onLocaleChange={(value) => setLocaleDraft({ locale: value })}
        onLanguageChange={(value) => setLocaleDraft({ language: value })}
        onApply={async () => {
          await Settings.SaveLocaleSettings({ locale: localeDraft.locale, language: localeDraft.language })
          applyLocaleSettings({ locale: localeDraft.locale, language: localeDraft.language })
          setLanguage(localeDraft.language)
        }}
      />
    )
  }

  if (generalTab === 'file') {
    return (
      <FileSettingsView
        dataDir={fileDraft.dataDir}
        currentDataDir={dataDir}
        dirty={fileDirty}
        onSelectFolder={async () => {
          const result = await File.SelectFolder({ folder_path: fileDraft.dataDir })
          if (result?.folder_path) {
            setFileDraft({ dataDir: result.folder_path })
          }
        }}
        onApply={() => undefined}
      />
    )
  }

  return (
    <DisplaySettingsView
      value={fontSize}
      draft={displayDraft.fontSize}
      onDraftChange={(value) => setDisplayDraft({ fontSize: value })}
      onApply={async () => {
        await Settings.SaveDisplaySettings({ font_size: displayDraft.fontSize })
        applyDisplaySettings(displayDraft.fontSize)
        setFontSize(displayDraft.fontSize)
      }}
      onReset={() => setDisplayDraft({ fontSize })}
    />
  )
}
