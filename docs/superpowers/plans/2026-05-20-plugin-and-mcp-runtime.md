# Plugin And MCP Runtime Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a settings-driven MCP management flow that imports local MCP servers, runs enabled stdio servers, exposes discovered tools in the chat input, and lets the existing `trpc-agent-go` runtime execute them.

**Architecture:** Add a backend `plugin` service plus a focused `backend/pkg/mcp` runtime manager that owns MCP import, config persistence, server lifecycle, and discovered tool metadata. Extend the existing agent tool registry from static metadata to dynamic tool definitions, then wire the settings page and chat input to real backend data so enabled MCP tools and actual runnable tools stay consistent.

**Tech Stack:** Go, Wails v3 services/bindings, `trpc-agent-go`, React 18, Zustand, Vitest, Testing Library, Tailwind CSS.

---

### Task 1: Backend MCP Runtime Foundations

**Files:**
- Create: `backend/pkg/mcp/manifest.go`
- Create: `backend/pkg/mcp/config.go`
- Create: `backend/pkg/mcp/importer.go`
- Create: `backend/pkg/mcp/manager.go`
- Create: `backend/pkg/mcp/manager_test.go`
- Modify: `backend/models/data_models/config.go`
- Test: `backend/pkg/mcp/manager_test.go`

- [ ] **Step 1: Write the failing MCP import and config tests**

```go
func TestImportMCPDirectoryCopiesIntoVersionedDataDir(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Name:        "filesystem",
		Version:     "1.2.3",
		Command:     "node",
		Args:        []string{"server.js"},
		ConfigFile:  "mcp.json",
		Description: "Filesystem tools",
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{
		SourceDir: sourceDir,
		Enabled:   true,
	})
	if err != nil {
		t.Fatalf("import mcp: %v", err)
	}
	if !strings.Contains(item.RootDir, filepath.Join("mcp", "filesystem", "1.2.3")) {
		t.Fatalf("expected versioned mcp path, got %q", item.RootDir)
	}
}

func TestImportMCPDirectoryFallsBackToDefaultVersionDir(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Name:       "clock",
		Command:    "node",
		Args:       []string{"server.js"},
		ConfigFile: "mcp.json",
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{SourceDir: sourceDir})
	if err != nil {
		t.Fatalf("import mcp: %v", err)
	}
	if !strings.Contains(item.RootDir, filepath.Join("mcp", "clock", "default")) {
		t.Fatalf("expected default version dir, got %q", item.RootDir)
	}
}
```

- [ ] **Step 2: Run the MCP tests to verify they fail**

Run: `go test ./backend/pkg/mcp`

Expected: FAIL because the package and import/runtime manager do not exist yet.

- [ ] **Step 3: Implement persisted MCP config models and import helpers**

```go
// backend/models/data_models/config.go
type ExtensionTool struct {
	ToolID          string `json:"tool_id"`
	ServerID        string `json:"server_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	RequiresConfirm bool   `json:"requires_confirm"`
}

type ExtensionItem struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Author         string          `json:"author"`
	Version        string          `json:"version"`
	Kind           string          `json:"kind"`
	Enabled        bool            `json:"enabled"`
	RootDir        string          `json:"root_dir"`
	SourceDir      string          `json:"source_dir"`
	ConfigFilePath string          `json:"config_file_path"`
	Tools          []ExtensionTool `json:"tools"`
}

type Config struct {
	Locale            string          `json:"locale"`
	Language          string          `json:"language"`
	FontSize          string          `json:"font_size"`
	DataDir           string          `json:"data_dir"`
	LogLevel          string          `json:"log_level"`
	DefaultProviderID uint            `json:"default_provider_id"`
	Extensions        []ExtensionItem `json:"extensions"`
}
```

- [ ] **Step 4: Implement `backend/pkg/mcp` import parsing and copying logic**

```go
type ImportInput struct {
	SourceDir string
	Enabled   bool
}

func (m *Manager) ImportMCP(ctx context.Context, input ImportInput) (data_models.ExtensionItem, error) {
	manifest, err := loadManifest(input.SourceDir)
	if err != nil {
		return data_models.ExtensionItem{}, err
	}
	version := manifest.Version
	if strings.TrimSpace(version) == "" {
		version = "default"
	}
	targetDir := filepath.Join(m.dataDir, "mcp", manifest.Name, version)
	if err := copyDirectory(input.SourceDir, targetDir); err != nil {
		return data_models.ExtensionItem{}, err
	}
	return data_models.ExtensionItem{
		ID:             buildExtensionID("mcp", manifest.Name, version),
		Name:           manifest.Name,
		Description:    manifest.Description,
		Author:         manifest.Author,
		Version:        manifest.Version,
		Kind:           "mcp",
		Enabled:        input.Enabled,
		RootDir:        targetDir,
		SourceDir:      input.SourceDir,
		ConfigFilePath: filepath.Join(targetDir, manifest.ConfigFile),
		Tools:          []data_models.ExtensionTool{},
	}, nil
}
```

- [ ] **Step 5: Run the MCP tests to verify they pass**

Run: `go test ./backend/pkg/mcp`

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/models/data_models/config.go backend/pkg/mcp/manifest.go backend/pkg/mcp/config.go backend/pkg/mcp/importer.go backend/pkg/mcp/manager.go backend/pkg/mcp/manager_test.go
git commit -m "feat: add mcp runtime foundations"
```

### Task 2: Plugin Service and Backend APIs

**Files:**
- Create: `backend/service/plugin/plugin.go`
- Create: `backend/service/plugin/plugin_implement.go`
- Create: `backend/service/plugin/plugin_internal.go`
- Create: `backend/service/plugin/plugin_test.go`
- Create: `backend/service/plugin/plugin_dto/list_extensions.go`
- Create: `backend/service/plugin/plugin_dto/import_extension.go`
- Create: `backend/service/plugin/plugin_dto/get_extension_detail.go`
- Create: `backend/service/plugin/plugin_dto/toggle_extension.go`
- Create: `backend/service/plugin/plugin_dto/reload_extension.go`
- Create: `backend/service/plugin/plugin_dto/delete_extension.go`
- Create: `backend/service/plugin/plugin_dto/save_extension_config.go`
- Modify: `backend/service/settings/settings.go`
- Test: `backend/service/plugin/plugin_test.go`

- [ ] **Step 1: Write the failing plugin service tests**

```go
func TestListExtensionsReturnsImportedItems(t *testing.T) {
	service := newTestPluginService(t)
	_, err := service.ImportExtension(context.Background(), plugin_dto.ImportExtensionInput{
		Kind: "mcp",
		Path: writeTestMCPBundleDir(t),
	})
	if err != nil {
		t.Fatalf("import extension: %v", err)
	}

	out, err := service.ListExtensions(context.Background(), plugin_dto.ListExtensionsInput{})
	if err != nil {
		t.Fatalf("list extensions: %v", err)
	}
	if len(out.Extensions) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(out.Extensions))
	}
}
```

- [ ] **Step 2: Run the plugin service tests to verify they fail**

Run: `go test ./backend/service/plugin`

Expected: FAIL because the plugin service does not exist yet.

- [ ] **Step 3: Implement typed DTOs and plugin service methods**

```go
func (s *Plugin) ListExtensions(ctx context.Context, input plugin_dto.ListExtensionsInput) (*plugin_dto.ListExtensionsOutput, error)
func (s *Plugin) ImportExtension(ctx context.Context, input plugin_dto.ImportExtensionInput) (*plugin_dto.ImportExtensionOutput, error)
func (s *Plugin) GetExtensionDetail(ctx context.Context, input plugin_dto.GetExtensionDetailInput) (*plugin_dto.GetExtensionDetailOutput, error)
func (s *Plugin) ToggleExtension(ctx context.Context, input plugin_dto.ToggleExtensionInput) (*plugin_dto.ToggleExtensionOutput, error)
func (s *Plugin) ReloadExtension(ctx context.Context, input plugin_dto.ReloadExtensionInput) (*plugin_dto.ReloadExtensionOutput, error)
func (s *Plugin) DeleteExtension(ctx context.Context, input plugin_dto.DeleteExtensionInput) (*plugin_dto.DeleteExtensionOutput, error)
func (s *Plugin) SaveExtensionConfig(ctx context.Context, input plugin_dto.SaveExtensionConfigInput) (*plugin_dto.SaveExtensionConfigOutput, error)
```

- [ ] **Step 4: Extend settings bootstrap with extension data**

Run code change in `backend/service/settings/settings.go` so `LoadBootstrap` includes saved extensions or a companion `plugin.ListExtensions` call can be used by the settings frontend.

- [ ] **Step 5: Run the plugin service tests to verify they pass**

Run: `go test ./backend/service/plugin`

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/service/plugin backend/service/settings/settings.go
git commit -m "feat: add plugin service apis"
```

### Task 3: Dynamic Agent Tool Registry and Runner Cache Fix

**Files:**
- Modify: `backend/pkg/agent/tools/registry.go`
- Modify: `backend/pkg/agent/tools/registry_test.go`
- Modify: `backend/pkg/agent/manager.go`
- Modify: `backend/pkg/agent/manager_test.go`
- Create: `backend/pkg/agent/tools/mcp_tool.go`
- Test: `backend/pkg/agent/tools/registry_test.go`
- Test: `backend/pkg/agent/manager_test.go`

- [ ] **Step 1: Write the failing registry and runner cache tests**

```go
func TestRegistryBuildsDynamicToolFromFactory(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterDefinition(ToolDefinition{
		Meta: ToolMeta{Name: "mcp:filesystem:read_file", Category: CategoryUser},
		Factory: func() toolpkg.Tool { return stubTool{name: "mcp:filesystem:read_file"} },
	})

	tools := registry.BuildTools([]string{"mcp:filesystem:read_file"})
	if len(tools) != 1 {
		t.Fatalf("expected dynamic tool to be built, got %d", len(tools))
	}
}

func TestGetOrCreateRunnerIncludesEnabledToolsInCacheKey(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	first, _ := manager.GetOrCreateRunner("https://api.example.com/v1", "key", "gpt-test", pkgProvider.OpenAiCompatibility, []string{"web_search"})
	second, _ := manager.GetOrCreateRunner("https://api.example.com/v1", "key", "gpt-test", pkgProvider.OpenAiCompatibility, []string{"code_exec"})
	if first == second {
		t.Fatal("expected different runner when enabled tools differ")
	}
}
```

- [ ] **Step 2: Run the agent tests to verify they fail**

Run: `go test ./backend/pkg/agent/...`

Expected: FAIL because the registry only stores metadata and the runner cache ignores tool selections.

- [ ] **Step 3: Refactor the tool registry to store definitions with factories**

```go
type ToolDefinition struct {
	Meta    ToolMeta
	Factory func() toolpkg.Tool
}

func (r *Registry) RegisterDefinition(def ToolDefinition)
func (r *Registry) BuildTools(enabled []string) []toolpkg.Tool
func (r *Registry) Remove(name string)
```

- [ ] **Step 4: Update `backend/pkg/agent/manager.go` to use registry factories and tool-aware cache keys**

```go
func runnerCacheKey(baseURL string, modelName string, providerType pkgProvider.Type, enabledUserTools []string) string {
	names := append([]string(nil), enabledUserTools...)
	slices.Sort(names)
	return fmt.Sprintf("%s|%s|%s|%s", baseURL, modelName, providerType, strings.Join(names, ","))
}
```

- [ ] **Step 5: Run the agent tests to verify they pass**

Run: `go test ./backend/pkg/agent/...`

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/pkg/agent/tools/registry.go backend/pkg/agent/tools/registry_test.go backend/pkg/agent/tools/mcp_tool.go backend/pkg/agent/manager.go backend/pkg/agent/manager_test.go
git commit -m "feat: support dynamic mcp tools in agent runtime"
```

### Task 4: Settings Page Plugin UI

**Files:**
- Modify: `frontend/src/types/settings.ts`
- Modify: `frontend/src/store/settingsStore.ts`
- Modify: `frontend/src/components/settings/SettingsPrimaryMenu.tsx`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`
- Create: `frontend/src/components/settings/plugins/PluginToolList.tsx`
- Create: `frontend/src/components/settings/plugins/PluginToolListItem.tsx`
- Create: `frontend/src/components/settings/plugins/PluginToolDetailView.tsx`
- Create: `frontend/src/components/settings/plugins/PluginToolConfigEditor.tsx`
- Create: `frontend/src/__tests__/pluginSettings.test.tsx`
- Test: `frontend/src/__tests__/pluginSettings.test.tsx`

- [ ] **Step 1: Write the failing settings plugin UI tests**

```tsx
it('renders plugins as a primary settings tab', () => {
  useSettingsStore.setState({
    ...getSettingsInitialState(),
    activeTab: 'plugins',
    extensions: [mockExtension()],
    selectedExtensionId: 'ext-1',
  })

  render(<SettingsApp />)

  expect(screen.getByText('Plugins & Tools')).toBeInTheDocument()
  expect(screen.getByText('filesystem')).toBeInTheDocument()
})
```

- [ ] **Step 2: Run the frontend settings test to verify it fails**

Run: `npm test -- --run frontend/src/__tests__/pluginSettings.test.tsx`

Expected: FAIL because the plugins tab and components do not exist yet.

- [ ] **Step 3: Implement typed settings state and plugin settings components**

```ts
export type SettingsPrimaryTab = 'general' | 'providers' | 'plugins' | 'about'

export type ExtensionToolItem = {
  tool_id: string
  server_id: string
  name: string
  description: string
  enabled: boolean
  requires_confirm: boolean
}

export type ExtensionItem = {
  id: string
  name: string
  description: string
  author: string
  version: string
  kind: 'mcp' | 'plugin'
  enabled: boolean
  runtime_status: 'stopped' | 'starting' | 'running' | 'error'
  runtime_message: string
  config_file_path: string
  tool_count: number
  tools: ExtensionToolItem[]
}
```

- [ ] **Step 4: Wire the new settings tab into `SettingsApp`**

Add plugin list/detail rendering that mirrors the providers split view and calls the new plugin service bindings for add, toggle, reload, delete, restore, and apply actions.

- [ ] **Step 5: Run the frontend settings test to verify it passes**

Run: `npm test -- --run frontend/src/__tests__/pluginSettings.test.tsx`

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add frontend/src/types/settings.ts frontend/src/store/settingsStore.ts frontend/src/components/settings/SettingsPrimaryMenu.tsx frontend/src/components/settings/SettingsApp.tsx frontend/src/components/settings/plugins frontend/src/__tests__/pluginSettings.test.tsx
git commit -m "feat: add plugin settings ui"
```

### Task 5: Chat Input Tool Integration

**Files:**
- Modify: `frontend/src/components/chat/ChatInput.tsx`
- Create: `frontend/src/__tests__/chatInputTools.test.tsx`
- Modify: `frontend/src/hooks/useSettingsBootstrap.ts`
- Test: `frontend/src/__tests__/chatInputTools.test.tsx`

- [ ] **Step 1: Write the failing chat input tool tests**

```tsx
it('renders enabled mcp tools and opens plugin settings', async () => {
  render(<ChatInput />)

  await user.click(screen.getByRole('button', { name: /tools/i }))

  expect(screen.getByText('read_file')).toBeInTheDocument()
  await user.click(screen.getByRole('button', { name: /manage tools and plugins/i }))
  expect(Window.OpenSettings).toHaveBeenCalledWith({ tab: 'plugins' })
})
```

- [ ] **Step 2: Run the chat input tool test to verify it fails**

Run: `npm test -- --run frontend/src/__tests__/chatInputTools.test.tsx`

Expected: FAIL because the tool list still uses `mockTools` and has no plugin-management footer.

- [ ] **Step 3: Replace `mockTools` with backend-provided tool data**

```ts
type ChatToolOption = {
  id: string
  name: string
  description: string
  category: 'builtin' | 'mcp'
  enabled: boolean
}
```

Load available tools from the plugin service, group them by category, keep the checked state in component state, and send selected tool ids through `enabledUserTools`.

- [ ] **Step 4: Add the fixed footer action**

Add a button to the bottom of the tool popover:

```tsx
<button onClick={() => { void Window.OpenSettings({ tab: 'plugins' }) }}>
  {t('input.manageToolsPlugins')}
</button>
```

- [ ] **Step 5: Run the chat input tool test to verify it passes**

Run: `npm test -- --run frontend/src/__tests__/chatInputTools.test.tsx`

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/chat/ChatInput.tsx frontend/src/hooks/useSettingsBootstrap.ts frontend/src/__tests__/chatInputTools.test.tsx
git commit -m "feat: connect chat input to mcp tools"
```

### Task 6: End-to-End Verification

**Files:**
- Modify: any files needed from previous tasks

- [ ] **Step 1: Run backend tests**

Run: `go test ./backend/pkg/mcp ./backend/service/plugin ./backend/pkg/agent/... ./backend/service/agent/...`

Expected: PASS

- [ ] **Step 2: Run frontend tests**

Run: `npm test -- --run frontend/src/__tests__/pluginSettings.test.tsx frontend/src/__tests__/chatInputTools.test.tsx frontend/src/__tests__/settingsApp.test.tsx`

Expected: PASS

- [ ] **Step 3: Run the production frontend build**

Run: `npm run build --prefix frontend`

Expected: PASS

- [ ] **Step 4: Commit final verification fixes**

```bash
git add .
git commit -m "test: verify plugin and mcp runtime flow"
```

---

## Self-Review

- Spec coverage: backend import/copy path, settings page management, runtime lifecycle, dynamic tool registration, runner cache invalidation, and chat input integration are all mapped to tasks above.
- Placeholder scan: no `TODO` / `TBD` placeholders remain; each task includes test-first verification and concrete files.
- Type consistency: `ExtensionItem`, `ExtensionToolItem`, dynamic `tool_id`, and `plugins` settings tab use the same names across backend and frontend tasks.

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-05-20-plugin-and-mcp-runtime.md`. The user has already chosen implementation now, so proceed with inline execution in this session using `superpowers:executing-plans`.
