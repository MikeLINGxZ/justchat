# Base Settings Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a standalone settings window with reusable settings UI primitives, persisted general settings, provider management, custom data-directory migration, and file-based logging.

**Architecture:** Add backend foundations first: config/data-dir resolution, logger, and provider/settings services. Then split the frontend into `MainLayout` and `SettingsApp`, build a reusable settings-shell/component framework, and wire each settings section to typed backend APIs through a dedicated settings store.

**Tech Stack:** Go, Wails v3 services/bindings, Gorm + SQLite, React 18, Zustand, i18next, Vitest, Testing Library, Tailwind CSS.

---

### Task 1: Data Directory, Config, and Logger Foundations

**Files:**
- Modify: `backend/models/data_models/config.go`
- Modify: `backend/pkg/dir/dir.go`
- Create: `backend/pkg/logger/logger.go`
- Create: `backend/pkg/logger/logger_test.go`
- Create: `backend/pkg/dir/dir_test.go`
- Test: `backend/pkg/logger/logger_test.go`
- Test: `backend/pkg/dir/dir_test.go`

- [ ] **Step 1: Write the failing data-dir and logger tests**

```go
package dir

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetDataDirUsesLocatorFileBeforeDefaultDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("LEMONTEA_DATA_DIR", "")

	targetDir := filepath.Join(tempHome, "custom-data")
	metaDir := filepath.Join(tempHome, ".lemontea")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content, err := json.Marshal(map[string]string{"data_dir": targetDir})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(metaDir, "data_dir.json"), content, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := GetDataDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != targetDir {
		t.Fatalf("expected custom data dir %q, got %q", targetDir, got)
	}
}
```

```go
package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerWritesOnlyMessagesAtOrAboveConfiguredLevel(t *testing.T) {
	tempDir := t.TempDir()
	clock := func() time.Time {
		return time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	}

	lg, err := New(Options{
		DataDir:  tempDir,
		Level:    "warn",
		Now:      clock,
		FileMode: 0o644,
		DirMode:  0o755,
	})
	if err != nil {
		t.Fatal(err)
	}

	lg.Info("skip this")
	lg.Warn("keep this")
	lg.Error("keep this too")

	bytes, err := os.ReadFile(filepath.Join(tempDir, "logs", "2026-05-10.log"))
	if err != nil {
		t.Fatal(err)
	}

	logText := string(bytes)
	if strings.Contains(logText, "skip this") {
		t.Fatalf("info message should not have been written: %s", logText)
	}
	if !strings.Contains(logText, "keep this") || !strings.Contains(logText, "keep this too") {
		t.Fatalf("expected warn and error messages in log file: %s", logText)
	}
}
```

- [ ] **Step 2: Run the backend tests to verify they fail**

Run: `go test ./backend/pkg/dir ./backend/pkg/logger`

Expected: FAIL because the locator-file resolution and logger package do not exist yet.

- [ ] **Step 3: Implement config fields, locator-file resolution, and logger package**

```go
// backend/models/data_models/config.go
package data_models

// Config data structure of config file (config.json)
type Config struct {
	Locale            string `json:"locale"`
	Language          string `json:"language"`
	FontSize          string `json:"font_size"`
	DataDir           string `json:"data_dir"`
	LogLevel          string `json:"log_level"`
	DefaultProviderID uint   `json:"default_provider_id"`
}
```

```go
// backend/pkg/dir/dir.go
package dir

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type dataDirLocator struct {
	DataDir string `json:"data_dir"`
}

func GetDefaultBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	baseDir := filepath.Join(homeDir, defaultBaseDirName)
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}
	return baseDir, nil
}

func GetLocatorFilePath() (string, error) {
	baseDir, err := GetDefaultBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "data_dir.json"), nil
}

func ReadLocatorDataDir() (string, error) {
	locatorPath, err := GetLocatorFilePath()
	if err != nil {
		return "", err
	}
	bytes, err := os.ReadFile(locatorPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var locator dataDirLocator
	if err := json.Unmarshal(bytes, &locator); err != nil {
		return "", err
	}
	return locator.DataDir, nil
}
```

```go
// backend/pkg/logger/logger.go
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	dataDir  string
	level    int
	now      func() time.Time
	fileMode os.FileMode
	dirMode  os.FileMode
	mu       sync.Mutex
}

type Options struct {
	DataDir  string
	Level    string
	Now      func() time.Time
	FileMode os.FileMode
	DirMode  os.FileMode
}

func New(options Options) (*Logger, error) {
	level, err := parseLevel(options.Level)
	if err != nil {
		return nil, err
	}
	if options.Now == nil {
		options.Now = time.Now
	}
	if options.FileMode == 0 {
		options.FileMode = 0o644
	}
	if options.DirMode == 0 {
		options.DirMode = 0o755
	}
	return &Logger{
		dataDir:  options.DataDir,
		level:    level,
		now:      options.Now,
		fileMode: options.FileMode,
		dirMode:  options.DirMode,
	}, nil
}

func (l *Logger) Info(message string) { _ = l.write("INFO", infoLevel, message) }
func (l *Logger) Warn(message string) { _ = l.write("WARN", warnLevel, message) }
func (l *Logger) Error(message string) { _ = l.write("ERROR", errorLevel, message) }
```

- [ ] **Step 4: Run the backend tests to verify they pass**

Run: `go test ./backend/pkg/dir ./backend/pkg/logger`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/models/data_models/config.go backend/pkg/dir/dir.go backend/pkg/dir/dir_test.go backend/pkg/logger/logger.go backend/pkg/logger/logger_test.go
git commit -m "feat: add config and logger foundations"
```

### Task 2: Database Bootstrap and Provider Backend Services

**Files:**
- Create: `backend/pkg/db/db.go`
- Create: `backend/models/view_model/provider.go`
- Create: `backend/service/provider.go`
- Create: `backend/service/provider_test.go`
- Modify: `backend/models/data_models/provider.go`
- Modify: `main.go`
- Test: `backend/service/provider_test.go`

- [ ] **Step 1: Write the failing provider service tests**

```go
package service

import (
	"testing"

	"changeme/backend/models/data_models"
)

func TestProviderServiceSetDefaultProviderClearsPreviousDefault(t *testing.T) {
	svc := newTestProviderService(t)

	first := data_models.Provider{ProviderName: "DeepSeek", ProviderType: data_models.ProviderTypeDeepseek, BaseUrl: "https://api.deepseek.com", Enable: true}
	second := data_models.Provider{ProviderName: "OpenRouter", ProviderType: data_models.ProviderTypeOpenrouter, BaseUrl: "https://openrouter.ai/api/v1", Enable: true}

	createdFirst, err := svc.CreateProvider(first)
	if err != nil {
		t.Fatal(err)
	}
	createdSecond, err := svc.CreateProvider(second)
	if err != nil {
		t.Fatal(err)
	}

	if err := svc.SetDefaultProvider(createdFirst.ID); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetDefaultProvider(createdSecond.ID); err != nil {
		t.Fatal(err)
	}

	items, err := svc.ListProviders()
	if err != nil {
		t.Fatal(err)
	}

	var firstDefault, secondDefault bool
	for _, item := range items {
		if item.ID == createdFirst.ID {
			firstDefault = item.IsDefault
		}
		if item.ID == createdSecond.ID {
			secondDefault = item.IsDefault
		}
	}

	if firstDefault {
		t.Fatal("first provider should no longer be default")
	}
	if !secondDefault {
		t.Fatal("second provider should be default")
	}
}
```

- [ ] **Step 2: Run the provider backend tests to verify they fail**

Run: `go test ./backend/service -run Provider`

Expected: FAIL because the DB bootstrap and provider service do not exist yet.

- [ ] **Step 3: Implement database bootstrap, provider view models, and provider service**

```go
// backend/pkg/db/db.go
package db

import (
	"path/filepath"

	"changeme/backend/models/data_models"
	"changeme/backend/pkg/dir"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, dir.DataBaseFileName)
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&data_models.Provider{}, &data_models.ProviderDefaultModel{}); err != nil {
		return nil, err
	}
	return gormDB, nil
}
```

```go
// backend/models/view_model/provider.go
package view_model

type ProviderItem struct {
	ID           uint   `json:"id"`
	ProviderName string `json:"provider_name"`
	ProviderType string `json:"provider_type"`
	BaseURL      string `json:"base_url"`
	Enabled      bool   `json:"enabled"`
	IsDefault    bool   `json:"is_default"`
	ModelCount   int    `json:"model_count"`
	Icon         string `json:"icon"`
}
```

```go
// backend/service/provider.go
package service

import (
	"changeme/backend/models/data_models"
	"changeme/backend/models/view_model"

	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

func NewProvider(db *gorm.DB) *Provider {
	return &Provider{db: db}
}

func (p *Provider) CreateProvider(input data_models.Provider) (*data_models.Provider, error) {
	if err := p.db.Create(&input).Error; err != nil {
		return nil, err
	}
	return &input, nil
}
```

- [ ] **Step 4: Run the provider backend tests to verify they pass**

Run: `go test ./backend/service -run Provider`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/db/db.go backend/models/view_model/provider.go backend/service/provider.go backend/service/provider_test.go backend/models/data_models/provider.go main.go go.mod go.sum
git commit -m "feat: add provider backend services"
```

### Task 3: Settings Backend Service, Config Persistence, and Migration Flow

**Files:**
- Create: `backend/models/view_model/settings.go`
- Create: `backend/service/settings.go`
- Create: `backend/service/settings_test.go`
- Modify: `backend/pkg/dir/dir.go`
- Modify: `main.go`
- Test: `backend/service/settings_test.go`

- [ ] **Step 1: Write the failing settings service tests**

```go
package service

import (
	"os"
	"path/filepath"
	"testing"

	"changeme/backend/models/data_models"
)

func TestSettingsServiceApplyFileSettingsMigratesDataAndWritesLocator(t *testing.T) {
	svc := newTestSettingsService(t)
	sourceDir := t.TempDir()
	targetDir := filepath.Join(t.TempDir(), "custom-target")

	if err := os.WriteFile(filepath.Join(sourceDir, "config.json"), []byte(`{"language":"zh-CN"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "data.db"), []byte("sqlite"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := svc.ApplyFileSettings(data_models.Config{DataDir: targetDir}, sourceDir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "config.json")); err != nil {
		t.Fatalf("expected migrated config file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetDir, "logs")); err != nil {
		t.Fatalf("expected logs directory: %v", err)
	}
}
```

- [ ] **Step 2: Run the settings backend tests to verify they fail**

Run: `go test ./backend/service -run Settings`

Expected: FAIL because config persistence and migration helpers do not exist yet.

- [ ] **Step 3: Implement settings view models, config persistence, and migration helpers**

```go
// backend/models/view_model/settings.go
package view_model

type SettingsBootstrap struct {
	Locale            string         `json:"locale"`
	Language          string         `json:"language"`
	FontSize          string         `json:"font_size"`
	DataDir           string         `json:"data_dir"`
	LogLevel          string         `json:"log_level"`
	DefaultProviderID uint           `json:"default_provider_id"`
	Version           string         `json:"version"`
	Providers         []ProviderItem `json:"providers"`
}
```

```go
// backend/service/settings.go
package service

import (
	"encoding/json"
	"os"
	"path/filepath"

	"changeme/backend/models/data_models"
	"changeme/backend/models/view_model"
	"changeme/backend/pkg/dir"
)

type Settings struct {
	providerService *Provider
}

func NewSettings(providerService *Provider) *Settings {
	return &Settings{providerService: providerService}
}

func (s *Settings) LoadBootstrap() (*view_model.SettingsBootstrap, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, err
	}
	providers, err := s.providerService.ListProviders()
	if err != nil {
		return nil, err
	}
	return &view_model.SettingsBootstrap{
		Locale:            config.Locale,
		Language:          config.Language,
		FontSize:          config.FontSize,
		DataDir:           config.DataDir,
		LogLevel:          config.LogLevel,
		DefaultProviderID: config.DefaultProviderID,
		Version:           "v0.0.1-dev",
		Providers:         providers,
	}, nil
}
```

- [ ] **Step 4: Run the settings backend tests to verify they pass**

Run: `go test ./backend/service -run Settings`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/models/view_model/settings.go backend/service/settings.go backend/service/settings_test.go backend/pkg/dir/dir.go main.go
git commit -m "feat: add settings backend service"
```

### Task 4: Frontend Settings Entry and Reusable Settings Framework

**Files:**
- Modify: `frontend/src/App.tsx`
- Create: `frontend/src/types/settings.ts`
- Create: `frontend/src/store/settingsStore.ts`
- Create: `frontend/src/components/settings/SettingsApp.tsx`
- Create: `frontend/src/components/settings/SettingsShell.tsx`
- Create: `frontend/src/components/settings/SettingsPrimaryMenu.tsx`
- Create: `frontend/src/components/settings/common/SettingsSectionLayout.tsx`
- Create: `frontend/src/components/settings/common/SettingsSubmenuList.tsx`
- Create: `frontend/src/components/settings/common/SettingsPanelHeader.tsx`
- Create: `frontend/src/components/settings/common/SettingsActionBar.tsx`
- Create: `frontend/src/components/settings/common/SettingsFieldRow.tsx`
- Create: `frontend/src/components/settings/common/SettingsDirtyGuard.tsx`
- Create: `frontend/src/hooks/useSettingsBootstrap.ts`
- Create: `frontend/src/__tests__/settingsStore.test.ts`
- Test: `frontend/src/__tests__/settingsStore.test.ts`

- [ ] **Step 1: Write the failing frontend settings-store tests**

```ts
import { beforeEach, describe, expect, it } from 'vitest'
import { useSettingsStore } from '@/store/settingsStore'

beforeEach(() => {
  useSettingsStore.setState(useSettingsStore.getInitialState())
})

describe('settingsStore', () => {
  it('marks display settings dirty when draft font size differs from applied font size', () => {
    useSettingsStore.getState().hydrate({
      fontSize: 'md',
      language: 'zh-CN',
      locale: 'zh-CN',
      dataDir: '/tmp/a',
      logLevel: 'info',
      providers: [],
      version: 'v0.0.1-dev',
      defaultProviderID: 0,
    })

    useSettingsStore.getState().setDisplayDraft({ fontSize: 'xl' })

    expect(useSettingsStore.getState().displayDirty).toBe(true)
  })
})
```

- [ ] **Step 2: Run the frontend settings tests to verify they fail**

Run: `cd frontend && npx vitest run src/__tests__/settingsStore.test.ts`

Expected: FAIL because the store and settings components do not exist yet.

- [ ] **Step 3: Implement the settings entry, typed store, and reusable UI framework**

```tsx
// frontend/src/App.tsx
import { ThemeProvider } from '@/components/providers/ThemeProvider'
import { FontSizeProvider } from '@/components/providers/FontSizeProvider'
import { MainLayout } from '@/components/layout/MainLayout'
import { SettingsApp } from '@/components/settings/SettingsApp'

function App() {
  const params = new URLSearchParams(window.location.search)
  const entry = params.get('entry')

  return (
    <ThemeProvider>
      <FontSizeProvider>
        {entry === 'settings' ? <SettingsApp /> : <MainLayout />}
      </FontSizeProvider>
    </ThemeProvider>
  )
}
```

```ts
// frontend/src/store/settingsStore.ts
import { create } from 'zustand'
import type { SettingsBootstrap, SettingsPrimaryTab } from '@/types/settings'

type SettingsState = {
  activeTab: SettingsPrimaryTab
  fontSize: string
  fontSizeDraft: string
  displayDirty: boolean
  hydrate: (payload: SettingsBootstrap) => void
  setDisplayDraft: (draft: { fontSize: string }) => void
  getInitialState: () => SettingsState
}
```

- [ ] **Step 4: Run the frontend settings tests to verify they pass**

Run: `cd frontend && npx vitest run src/__tests__/settingsStore.test.ts`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.tsx frontend/src/types/settings.ts frontend/src/store/settingsStore.ts frontend/src/components/settings frontend/src/hooks/useSettingsBootstrap.ts frontend/src/__tests__/settingsStore.test.ts
git commit -m "feat: add reusable settings framework"
```

### Task 5: General Settings Views and Global Runtime Synchronization

**Files:**
- Create: `frontend/src/components/settings/general/GeneralSettingsPanel.tsx`
- Create: `frontend/src/components/settings/general/DisplaySettingsView.tsx`
- Create: `frontend/src/components/settings/general/LocaleSettingsView.tsx`
- Create: `frontend/src/components/settings/general/FileSettingsView.tsx`
- Modify: `frontend/src/store/appStore.ts`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`
- Create: `frontend/src/__tests__/generalSettings.test.tsx`
- Test: `frontend/src/__tests__/generalSettings.test.tsx`

- [ ] **Step 1: Write the failing general-settings UI tests**

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { DisplaySettingsView } from '@/components/settings/general/DisplaySettingsView'

describe('DisplaySettingsView', () => {
  it('enables apply when font size draft changes', async () => {
    const user = userEvent.setup()
    const onApply = vi.fn()

    render(
      <DisplaySettingsView
        value="md"
        draft="md"
        onDraftChange={() => {}}
        onApply={onApply}
      />
    )

    await user.click(screen.getByRole('button', { name: 'Large' }))

    expect(screen.getByRole('button', { name: 'Apply' })).toBeEnabled()
  })
})
```

- [ ] **Step 2: Run the general-settings UI tests to verify they fail**

Run: `cd frontend && npx vitest run src/__tests__/generalSettings.test.tsx`

Expected: FAIL because the general settings views do not exist yet.

- [ ] **Step 3: Implement the general settings panels and runtime sync**

```tsx
// frontend/src/components/settings/general/DisplaySettingsView.tsx
export function DisplaySettingsView(props: {
  value: FontSize
  draft: FontSize
  onDraftChange: (next: FontSize) => void
  onApply: () => void
  onReset: () => void
}) {
  const dirty = props.value !== props.draft

  return (
    <div className="flex h-full flex-col gap-6">
      <SettingsPanelHeader
        title="Display settings"
        description="Adjust reading size and preview the result before applying it."
      />
      <div className="grid grid-cols-5 gap-3">
        {FONT_SIZE_OPTIONS.map(option => (
          <button key={option.value} onClick={() => props.onDraftChange(option.value)}>
            {option.label}
          </button>
        ))}
      </div>
      <SettingsActionBar
        primaryLabel="Apply"
        secondaryLabel="Reset"
        primaryDisabled={!dirty}
        onPrimaryClick={props.onApply}
        onSecondaryClick={props.onReset}
      />
    </div>
  )
}
```

- [ ] **Step 4: Run the general-settings UI tests to verify they pass**

Run: `cd frontend && npx vitest run src/__tests__/generalSettings.test.tsx`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/settings/general frontend/src/store/appStore.ts frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts frontend/src/__tests__/generalSettings.test.tsx
git commit -m "feat: add general settings views"
```

### Task 6: Provider Settings UI and About View

**Files:**
- Create: `frontend/src/components/settings/providers/ProviderSettingsPanel.tsx`
- Create: `frontend/src/components/settings/providers/ProviderList.tsx`
- Create: `frontend/src/components/settings/providers/ProviderListItem.tsx`
- Create: `frontend/src/components/settings/providers/ProviderDetailView.tsx`
- Create: `frontend/src/components/settings/providers/AddProviderDialog.tsx`
- Create: `frontend/src/components/settings/about/AboutSettingsView.tsx`
- Create: `frontend/src/__tests__/providerSettings.test.tsx`
- Test: `frontend/src/__tests__/providerSettings.test.tsx`

- [ ] **Step 1: Write the failing provider-settings UI tests**

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { ProviderList } from '@/components/settings/providers/ProviderList'

describe('ProviderList', () => {
  it('shows the default badge for the current default provider', async () => {
    render(
      <ProviderList
        items={[
          { id: 1, provider_name: 'DeepSeek', provider_type: 'deepseek', base_url: 'https://api.deepseek.com', enabled: true, is_default: true, model_count: 0, icon: 'deepseek' },
        ]}
        selectedId={1}
        onSelect={vi.fn()}
        onCreate={vi.fn()}
      />
    )

    expect(screen.getByText('Default')).toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run the provider-settings UI tests to verify they fail**

Run: `cd frontend && npx vitest run src/__tests__/providerSettings.test.tsx`

Expected: FAIL because the provider settings views do not exist yet.

- [ ] **Step 3: Implement the provider settings and about views**

```tsx
// frontend/src/components/settings/providers/ProviderListItem.tsx
export function ProviderListItem({ item, selected, onSelect, onToggle }: ProviderListItemProps) {
  return (
    <button
      type="button"
      onClick={() => onSelect(item.id)}
      className={cn('flex w-full items-start gap-3 rounded-xl border p-3 text-left', selected && 'border-primary bg-accent/40')}
    >
      <ProviderIcon type={item.provider_type} />
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className="truncate font-medium">{item.provider_name}</span>
          {item.is_default && <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">Default</span>}
        </div>
        <p className="truncate text-xs text-muted-foreground">{item.base_url}</p>
      </div>
      <Switch checked={item.enabled} onCheckedChange={onToggle} />
    </button>
  )
}
```

- [ ] **Step 4: Run the provider-settings UI tests to verify they pass**

Run: `cd frontend && npx vitest run src/__tests__/providerSettings.test.tsx`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/settings/providers frontend/src/components/settings/about/AboutSettingsView.tsx frontend/src/__tests__/providerSettings.test.tsx
git commit -m "feat: add provider and about settings views"
```

### Task 7: End-to-End Wiring, Bindings Refresh, and Verification

**Files:**
- Modify: `frontend/src/components/sidebar/SidebarFooter.tsx`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`
- Modify: `frontend/src/components/settings/general/GeneralSettingsPanel.tsx`
- Modify: `frontend/src/components/settings/general/DisplaySettingsView.tsx`
- Modify: `frontend/src/components/settings/general/LocaleSettingsView.tsx`
- Modify: `frontend/src/components/settings/general/FileSettingsView.tsx`
- Modify: `frontend/src/components/settings/providers/ProviderSettingsPanel.tsx`
- Modify: `frontend/src/components/settings/providers/ProviderList.tsx`
- Modify: `frontend/src/components/settings/providers/ProviderListItem.tsx`
- Modify: `frontend/src/components/settings/providers/ProviderDetailView.tsx`
- Modify: `frontend/src/components/settings/providers/AddProviderDialog.tsx`
- Modify: `frontend/bindings/changeme/backend/service/*.ts`
- Test: `frontend/src/__tests__/appStore.test.ts`
- Test: `frontend/src/__tests__/settingsStore.test.ts`
- Test: `frontend/src/__tests__/generalSettings.test.tsx`
- Test: `frontend/src/__tests__/providerSettings.test.tsx`

- [ ] **Step 1: Write the final integration checks before wiring**

```tsx
// Extend the existing appStore tests to ensure runtime state updates after settings apply.
it('setFontSize updates fontSize', () => {
  useAppStore.getState().setFontSize('xl')
  expect(useAppStore.getState().fontSize).toBe('xl')
})
```

```go
// Add one final service integration assertion after wiring both services into main.go.
func TestSettingsBootstrapIncludesProviders(t *testing.T) {
	providerService := newTestProviderService(t)
	settingsService := NewSettings(providerService)

	created, err := providerService.CreateProvider(data_models.Provider{
		ProviderName: "DeepSeek",
		ProviderType: data_models.ProviderTypeDeepseek,
		BaseUrl:      "https://api.deepseek.com",
		Enable:       true,
	})
	if err != nil {
		t.Fatal(err)
	}

	payload, err := settingsService.LoadBootstrap()
	if err != nil {
		t.Fatal(err)
	}

	if len(payload.Providers) != 1 {
		t.Fatalf("expected one provider in bootstrap payload, got %d", len(payload.Providers))
	}
	if payload.Providers[0].ID != created.ID {
		t.Fatalf("expected provider id %d in bootstrap payload, got %d", created.ID, payload.Providers[0].ID)
	}
}
```

- [ ] **Step 2: Run the full verification suite and confirm any failing gaps**

Run: `go test ./...`

Run: `cd frontend && npx vitest run src/__tests__/appStore.test.ts src/__tests__/settingsStore.test.ts src/__tests__/generalSettings.test.tsx src/__tests__/providerSettings.test.tsx`

Expected: any remaining failures point to missing bindings or store wiring.

- [ ] **Step 3: Wire the live settings flow and regenerate bindings**

```tsx
// frontend/src/components/sidebar/SidebarFooter.tsx
import { OpenSettings } from '@/../bindings/changeme/backend/service/window'

<button
  className="w-full flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent text-left"
  onClick={() => {
    void OpenSettings()
    close()
  }}
>
  <Settings size={14} />
  {t('settings.settings')}
</button>
```

```go
// main.go
dbConn, err := db.Open()
if err != nil {
	log.Fatal(err)
}

providerService := service.NewProvider(dbConn)

app := application.New(application.Options{
	Services: []application.Service{
		application.NewService(service.NewSettings(providerService)),
		application.NewService(providerService),
		application.NewService(&service.Process{}),
		application.NewService(&service.Window{}),
	},
})
```

- [ ] **Step 4: Run the full verification suite again**

Run: `go test ./...`

Run: `cd frontend && npx vitest run src/__tests__/appStore.test.ts src/__tests__/settingsStore.test.ts src/__tests__/generalSettings.test.tsx src/__tests__/providerSettings.test.tsx`

Run: `cd frontend && npm run build`

Expected: all tests pass and the frontend build completes successfully.

- [ ] **Step 5: Commit**

```bash
git add main.go backend/service frontend/src/components/sidebar/SidebarFooter.tsx frontend/bindings/changeme/backend/service
git commit -m "feat: wire full settings experience"
```
