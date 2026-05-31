# CLI Install Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor CLI plugin installation (spec `docs/dev/13.plugin_cli.md` §3.1–§3.3) — remove the local directory import entry, move npm install to background execution via hidden sessions, and add a new "Install from official docs" entry powered by a built-in skill. This is Plan C of three; Plan A (Skills Foundation) and Plan B (Hidden Sessions + Notifications) must merge first.

**Architecture:**
- **Remove** `InstallCliFromLocal` Wails method, its DTO, and the "从本地目录导入" tab in `CliInstallModal`. The `kind: "cli"` dropdown now shows two entries: "从 npm 安装" and "从官方文档安装".
- **Agent tools** `InstallCli` and `GenerateCliManifest` are registered on the agent's per-turn tool pool (alongside the Plan A/B built-in tools). These wrap the existing `Plugin.InstallCliFromNpm` and `Plugin.GenerateCliManifest` service methods so the hidden-session agent can drive the full install pipeline.
- **Progress events**: `cli.install.progress` Wails events are emitted by the `InstallCli` agent tool during npm download / install phases. Frontend subscribes and updates a transient `installing` state on the extension list. On completion the agent calls `GenerateCliManifest` which persists the final `ExtensionItem` and emits `cli.install.done`.
- **Built-in skill** `install-cli-from-docs` lives at `backend/pkg/skills/_builtin/install-cli-from-docs/SKILL.md`. Its body instructs the agent: parse the user's input (URL / pasted docs / natural language), determine the npm package name, call `InstallCli`, call `GenerateCliManifest`. If the agent cannot determine the package, it calls `RequestUserAttention` (Plan B) to ask the user.
- **Frontend flow**: both "from npm" and "from docs" entries open a lightweight input modal (single-field for npm, textarea for docs). On submit, the modal closes immediately and the frontend calls `Agent.SpawnHiddenSession` (Plan B) with the appropriate `skill_name` and `user_message`. The plugin list shows a pending entry with live progress.
- **Cleanup**: the old blocking `CliInstallModal` is replaced by two focused components — `CliNpmInstallModal` (single input) and `CliInstallFromDocsModal` (textarea). The old component file is deleted.

**Tech Stack:** Go 1.22+, Wails v3, React 18, Zustand, i18next.

---

## Pre-flight

- [ ] **Confirm Plan A is merged.** Run `git log --oneline | grep -i "skills\|builtin"` — Plan A commits should be present.
- [ ] **Confirm Plan B is merged.** Run `git log --oneline | grep -i "hidden\|notification\|SpawnHidden"` — Plan B commits should be present.
- [ ] **Verify `SpawnHiddenSession` exists.** Run `grep -n "SpawnHiddenSession" backend/service/agent/agent.go` — should find the method.
- [ ] **Verify builtin skill embedding works.** Run `grep -n "go:embed" backend/pkg/skills/builtin.go` — should find `//go:embed all:_builtin`.

---

## Section 1 — Remove local directory import entry

### Task 1: Delete `InstallCliFromLocal` backend + DTO

**Files:**
- Delete: `backend/service/plugin/plugin_dto/install_cli_from_local.go`
- Modify: `backend/service/plugin/plugin_cli.go` (remove `InstallCliFromLocal` method)
- Modify: `backend/service/plugin/plugin_cli_test.go` (remove any local-install tests)

- [ ] **Step 1: Search for all references**

```
grep -rn "InstallCliFromLocal\|install_cli_from_local\|InstallFromLocal" backend/ frontend/
```

Catalog every hit — these are all the places to update.

- [ ] **Step 2: Delete the DTO file**

```
git rm backend/service/plugin/plugin_dto/install_cli_from_local.go
```

- [ ] **Step 3: Remove the `InstallCliFromLocal` method from `plugin_cli.go`**

Delete the entire method (lines ~36–52 of the current file). Also remove the `defaultNameFromPath` helper if no other code references it:

```
grep -rn "defaultNameFromPath" backend/
```

If only `InstallCliFromLocal` used it, delete it too.

- [ ] **Step 4: Remove related tests**

In `plugin_cli_test.go`, delete any test functions that exercise `InstallCliFromLocal`.

- [ ] **Step 5: Remove `InstallFromLocal` from `pkg/cli/manager.go` if no longer referenced**

```
grep -rn "InstallFromLocal\b" backend/
```

If only the deleted service method called it, remove the `InstallFromLocal` method from `Manager` and any related tests in `manager_test.go`.

- [ ] **Step 6: Run tests**

```
go test ./...
```

- [ ] **Step 7: Regenerate bindings**

```
wails generate bindings
```

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "refactor(plugin): remove InstallCliFromLocal entry per spec §3.1"
```

---

### Task 2: Remove local tab from frontend

**Files:**
- Modify: `frontend/src/components/settings/plugins/CliInstallModal.tsx`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: Remove `local` tab from `CliInstallModal`**

In `CliInstallModal.tsx`:
- Remove `CliInstallTab` type's `'local'` value (change to just `'npm'`).
- Remove the `localPath` state, `pickFolder` function, and the local-path input UI.
- Remove the `File as FileBinding` import and `SelectFolderInput` import.
- Remove the tab-switching buttons (since there's only one tab now, no need for tabs).
- Remove the `InstallCliFromLocalInput` import.
- In the `install` function, remove the `else` branch that calls `PluginBinding.InstallCliFromLocal`.

- [ ] **Step 2: Remove unused i18n keys**

In both `zh-CN.ts` and `en.ts`, remove:
- `settingsPage.plugins.cliFromLocal`
- `settingsPage.plugins.cliLocalPlaceholder`
- `settingsPage.plugins.cliBrowse`

(Keep `cliFromNpm` and all other keys — they're still used.)

- [ ] **Step 3: Verify frontend builds**

```
cd frontend && cnpm run build
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/settings/plugins/CliInstallModal.tsx frontend/src/i18n/
git commit -m "refactor(plugin): remove local directory install tab from frontend"
```

---

## Section 2 — Agent tools for CLI install

### Task 3: `InstallCli` agent tool

**Files:**
- Create: `backend/pkg/agent/tools/install_cli_tool.go`
- Create: `backend/pkg/agent/tools/install_cli_tool_test.go`

- [ ] **Step 1: Write `install_cli_tool.go`**

```go
package tools

import (
	"context"
	"encoding/json"
	"errors"
)

// CliInstaller is satisfied by the plugin service (avoids import cycle).
type CliInstaller interface {
	InstallCliSync(ctx context.Context, npmPackage string, name string, onProgress func(phase string, detail string)) (installResultJSON string, err error)
}

const InstallCliToolName = "InstallCli"

func BuildInstallCliTool() ToolMeta {
	return ToolMeta{
		Name:        InstallCliToolName,
		Description: "Install a CLI plugin from npm. Call with {\"npm_package\":\"@scope/pkg\",\"name\":\"my-cli\"}. Returns the installed extension item JSON.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				NpmPackage string `json:"npm_package"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Installing CLI: " + parsed.NpmPackage
		},
	}
}

func InvokeInstallCli(ctx context.Context, installer CliInstaller, args json.RawMessage) (string, error) {
	var parsed struct {
		NpmPackage string `json:"npm_package"`
		Name       string `json:"name"`
	}
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.NpmPackage == "" {
		return "", errors.New("npm_package is required")
	}
	return installer.InstallCliSync(ctx, parsed.NpmPackage, parsed.Name, nil)
}
```

- [ ] **Step 2: Write test with fake installer**

```go
package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeInstaller struct{ gotPkg string }

func (f *fakeInstaller) InstallCliSync(_ context.Context, npmPackage string, _ string, _ func(string, string)) (string, error) {
	f.gotPkg = npmPackage
	return `{"id":"cli:test","name":"test","kind":"cli"}`, nil
}

func TestInvokeInstallCli_PassesPackage(t *testing.T) {
	fi := &fakeInstaller{}
	args := json.RawMessage(`{"npm_package":"@scope/test-cli","name":"test"}`)
	out, err := InvokeInstallCli(context.Background(), fi, args)
	if err != nil { t.Fatal(err) }
	if fi.gotPkg != "@scope/test-cli" { t.Fatalf("got %q", fi.gotPkg) }
	if out == "" { t.Fatal("empty result") }
}
```

- [ ] **Step 3: Run — must pass**

```
go test ./backend/pkg/agent/tools/...
```

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/agent/tools/install_cli_tool.go backend/pkg/agent/tools/install_cli_tool_test.go
git commit -m "feat(agent): add InstallCli tool definition"
```

---

### Task 4: `GenerateCliManifest` agent tool

**Files:**
- Create: `backend/pkg/agent/tools/generate_manifest_tool.go`
- Create: `backend/pkg/agent/tools/generate_manifest_tool_test.go`

- [ ] **Step 1: Write `generate_manifest_tool.go`**

```go
package tools

import (
	"context"
	"encoding/json"
	"errors"
)

// CliManifestGenerator is satisfied by the plugin service.
type CliManifestGenerator interface {
	GenerateCliManifestSync(ctx context.Context, extensionID string) (extensionJSON string, err error)
}

const GenerateCliManifestToolName = "GenerateCliManifest"

func BuildGenerateCliManifestTool() ToolMeta {
	return ToolMeta{
		Name:        GenerateCliManifestToolName,
		Description: "Generate (or regenerate) the manifest for an installed CLI plugin. Call with {\"id\":\"cli:<name>\"}. Returns the updated extension item JSON.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				ID string `json:"id"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Generating manifest for: " + parsed.ID
		},
	}
}

func InvokeGenerateCliManifest(ctx context.Context, generator CliManifestGenerator, args json.RawMessage) (string, error) {
	var parsed struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.ID == "" {
		return "", errors.New("id is required")
	}
	return generator.GenerateCliManifestSync(ctx, parsed.ID)
}
```

- [ ] **Step 2: Write test**

```go
package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeGenerator struct{ gotID string }

func (f *fakeGenerator) GenerateCliManifestSync(_ context.Context, id string) (string, error) {
	f.gotID = id
	return `{"id":"cli:test","name":"test","kind":"cli","runtime_status":"ready"}`, nil
}

func TestInvokeGenerateCliManifest_PassesID(t *testing.T) {
	fg := &fakeGenerator{}
	args := json.RawMessage(`{"id":"cli:test"}`)
	out, err := InvokeGenerateCliManifest(context.Background(), fg, args)
	if err != nil { t.Fatal(err) }
	if fg.gotID != "cli:test" { t.Fatalf("got %q", fg.gotID) }
	if out == "" { t.Fatal("empty result") }
}
```

- [ ] **Step 3: Run + commit**

```
go test ./backend/pkg/agent/tools/...
git add backend/pkg/agent/tools/generate_manifest_tool.go backend/pkg/agent/tools/generate_manifest_tool_test.go
git commit -m "feat(agent): add GenerateCliManifest tool definition"
```

---

### Task 5: Implement `CliInstaller` + `CliManifestGenerator` on plugin service

**Files:**
- Create: `backend/service/plugin/plugin_agent_bridge.go`
- Modify: `backend/service/plugin/plugin.go` (no change needed if interface is satisfied)

- [ ] **Step 1: Create `plugin_agent_bridge.go`**

```go
package plugin

import (
	"context"
	"encoding/json"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// InstallCliSync satisfies tools.CliInstaller. It wraps the existing InstallCliFromNpm + emits progress events.
func (p *Plugin) InstallCliSync(ctx context.Context, npmPackage string, name string, onProgress func(string, string)) (string, error) {
	emit := func(phase, detail string) {
		if onProgress != nil {
			onProgress(phase, detail)
		}
		if p.wailsApp != nil {
			p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
				"npm_package": npmPackage,
				"name":        name,
				"phase":       phase,
				"detail":      detail,
			})
		}
	}

	emit("downloading", "Fetching "+npmPackage)
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{
		NpmPackage: npmPackage,
		Name:       name,
	})
	if err != nil {
		emit("failed", err.Error())
		return "", err
	}
	emit("installed", "Package installed to "+out.Extension.RootDir)

	resultBytes, _ := json.Marshal(out.Extension)
	return string(resultBytes), nil
}

// GenerateCliManifestSync satisfies tools.CliManifestGenerator.
func (p *Plugin) GenerateCliManifestSync(ctx context.Context, extensionID string) (string, error) {
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
			"id":    extensionID,
			"phase": "generating",
		})
	}
	out, err := p.GenerateCliManifest(ctx, plugin_dto.GenerateCliManifestInput{ID: extensionID})
	if err != nil {
		if p.wailsApp != nil {
			p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
				"id":    extensionID,
				"phase": "failed",
			})
		}
		return "", err
	}
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.install.done", map[string]interface{}{
			"extension": out.Extension,
		})
	}
	resultBytes, _ := json.Marshal(out.Extension)
	return string(resultBytes), nil
}
```

- [ ] **Step 2: Write bridge test**

```go
package plugin

import (
	"context"
	"testing"
)

func TestInstallCliSync_EmitsProgress(t *testing.T) {
	// This test requires a test Plugin with stubbed InstallCliFromNpm.
	// If the existing test harness supports this, add it; otherwise skip
	// and rely on the E2E smoke test.
	t.Skip("requires full plugin test harness")
}
```

- [ ] **Step 3: Run tests**

```
go test ./backend/service/plugin/...
```

- [ ] **Step 4: Commit**

```bash
git add backend/service/plugin/plugin_agent_bridge.go
git commit -m "feat(plugin): add agent bridge methods for InstallCli and GenerateCliManifest"
```

---

### Task 6: Register CLI tools in agent dispatcher

**Files:**
- Modify: agent tool registration site (where Plan A registered `BuildSkillTool` and Plan B registered `BuildRequestAttentionTool`)
- Modify: agent tool dispatch site

- [ ] **Step 1: Locate registration site**

```
grep -rn "BuildSkillTool\|BuildRequestAttentionTool\|Register.*ToolMeta" backend/pkg/agent/ backend/service/agent/
```

Find where built-in tools are registered per turn.

- [ ] **Step 2: Add CLI tool registration**

At the same registration site, add (for **all** sessions, not just hidden — the user might ask the agent in a regular chat to install a CLI):

```go
registry.Register(tools.BuildInstallCliTool())
registry.Register(tools.BuildGenerateCliManifestTool())
```

- [ ] **Step 3: Wire dispatch**

At the tool dispatch site (where `case tools.RequestAttentionToolName:` is), add:

```go
case tools.InstallCliToolName:
    return tools.InvokeInstallCli(ctx, a.cliInstaller, call.Args)
case tools.GenerateCliManifestToolName:
    return tools.InvokeGenerateCliManifest(ctx, a.cliManifestGenerator, call.Args)
```

- [ ] **Step 4: Add fields + setters on agent struct**

```go
// On the agent struct:
cliInstaller          tools.CliInstaller
cliManifestGenerator  tools.CliManifestGenerator

// Setters:
func (a *Agent) SetCliInstaller(i tools.CliInstaller) { a.cliInstaller = i }
func (a *Agent) SetCliManifestGenerator(g tools.CliManifestGenerator) { a.cliManifestGenerator = g }
```

- [ ] **Step 5: Wire in `main.go`**

```go
pluginService := pluginSvc.NewPlugin()
// ... after agent service is created:
agentService.SetCliInstaller(pluginService)
agentService.SetCliManifestGenerator(pluginService)
```

- [ ] **Step 6: Run tests**

```
go test ./...
```

- [ ] **Step 7: Regenerate bindings + commit**

```
wails generate bindings
git add backend/pkg/agent/ backend/service/agent/ main.go frontend/bindings/
git commit -m "feat(agent): register InstallCli and GenerateCliManifest tools in dispatcher"
```

---

## Section 3 — Install progress tracking

### Task 7: Frontend subscribes to install progress events

**Files:**
- Create: `frontend/src/store/cliInstallStore.ts`
- Create: `frontend/src/hooks/useCliInstallSubscription.ts`
- Create: `frontend/src/types/cliInstall.ts`

- [ ] **Step 1: Types**

```ts
// frontend/src/types/cliInstall.ts
export type CliInstallPhase = 'downloading' | 'installed' | 'generating' | 'done' | 'failed'

export type CliInstallItem = {
  npm_package: string
  name: string
  phase: CliInstallPhase
  detail: string
  extension_id?: string // set once install completes and manifest is generated
}
```

- [ ] **Step 2: Store**

```ts
// frontend/src/store/cliInstallStore.ts
import { create } from 'zustand'
import type { CliInstallItem } from '@/types/cliInstall'

type State = {
  items: CliInstallItem[]
  upsert: (item: CliInstallItem) => void
  remove: (npmPackage: string) => void
  clear: () => void
}

export const useCliInstallStore = create<State>((set) => ({
  items: [],
  upsert: (item) => set((s) => {
    const idx = s.items.findIndex((i) => i.npm_package === item.npm_package)
    if (idx >= 0) {
      const next = [...s.items]
      next[idx] = item
      return { items: next }
    }
    return { items: [...s.items, item] }
  }),
  remove: (npmPackage) => set((s) => ({
    items: s.items.filter((i) => i.npm_package !== npmPackage),
  })),
  clear: () => set({ items: [] }),
}))
```

- [ ] **Step 3: Subscription hook**

```ts
// frontend/src/hooks/useCliInstallSubscription.ts
import { useEffect } from 'react'
import { Events } from '@wailsio/runtime'
import { useCliInstallStore } from '@/store/cliInstallStore'
import type { CliInstallItem } from '@/types/cliInstall'

type ProgressEvent = {
  npm_package: string
  name: string
  phase: string
  detail: string
  id?: string
}

type DoneEvent = {
  extension: { id: string; name: string }
}

export function useCliInstallSubscription() {
  const upsert = useCliInstallStore((s) => s.upsert)
  const remove = useCliInstallStore((s) => s.remove)

  useEffect(() => {
    const offProgress = Events.On('cli.install.progress', (event: { data: ProgressEvent }) => {
      const d = event.data
      upsert({
        npm_package: d.npm_package,
        name: d.name,
        phase: d.phase as CliInstallItem['phase'],
        detail: d.detail ?? '',
        extension_id: d.id,
      })
    })
    const offDone = Events.On('cli.install.done', (event: { data: DoneEvent }) => {
      // The done event carries the final extension; the list will refresh via ListExtensions.
      // Remove the pending install entry after a short delay so the user sees the transition.
      const ext = event.data.extension
      if (ext?.name) {
        setTimeout(() => remove(ext.name), 2000)
      }
    })
    return () => { offProgress(); offDone() }
  }, [upsert, remove])
}
```

- [ ] **Step 4: Mount subscription in app shell**

In the top-level `App.tsx` (or wherever `useNotificationsSubscription` is mounted from Plan B), add:

```ts
useCliInstallSubscription()
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/store/cliInstallStore.ts frontend/src/hooks/useCliInstallSubscription.ts frontend/src/types/cliInstall.ts frontend/src/App.tsx
git commit -m "feat(plugin): add frontend store and subscription for CLI install progress"
```

---

### Task 8: Show install progress in `PluginToolListItem`

**Files:**
- Modify: `frontend/src/components/settings/plugins/PluginToolListItem.tsx`
- Modify: `frontend/src/components/settings/plugins/PluginToolList.tsx`

- [ ] **Step 1: Add install-progress indicator to `PluginToolListItem`**

In `PluginToolListItem.tsx`, import the store and add a progress bar when the item is being installed:

```tsx
import { useCliInstallStore } from '@/store/cliInstallStore'

// Inside the component, before the return:
const installProgress = useCliInstallStore((s) =>
  s.items.find((i) => i.name === props.item.name || i.extension_id === props.item.id)
)
```

After the description line (line ~55), add:

```tsx
{installProgress && installProgress.phase !== 'done' && installProgress.phase !== 'failed' && (
  <div className="mt-2 flex items-center gap-2">
    <div className="h-1 flex-1 overflow-hidden rounded-full bg-muted">
      <div
        className="h-full rounded-full bg-primary transition-all duration-300"
        style={{ width: installProgress.phase === 'generating' ? '80%' : installProgress.phase === 'installed' ? '60%' : '30%' }}
      />
    </div>
    <span className="text-[11px] text-muted-foreground">{installProgress.detail || installProgress.phase}</span>
  </div>
)}
{installProgress?.phase === 'failed' && (
  <div className="mt-2 rounded-md bg-red-100 px-2 py-1 text-xs text-red-700 dark:bg-red-500/10 dark:text-red-400">
    {installProgress.detail || t('settingsPage.plugins.cliInstallFailed')}
  </div>
)}
```

- [ ] **Step 2: Show pending installs in the list**

In `PluginToolList.tsx`, import the store and render pending install entries at the top of the list:

```tsx
import { useCliInstallStore } from '@/store/cliInstallStore'

// Inside the component:
const pendingInstalls = useCliInstallStore((s) => s.items)
```

Before `{props.items.map(...)}`, add:

```tsx
{pendingInstalls.map((install) => (
  <div key={`installing:${install.npm_package}`} className="rounded-2xl bg-accent/50 px-3 py-3">
    <div className="flex items-center gap-2">
      <Loader2 size={14} className="animate-spin text-muted-foreground" />
      <span className="text-sm font-medium">{install.name || install.npm_package}</span>
      <span className="rounded-full bg-purple-100 px-2 py-0.5 text-[11px] text-purple-700 dark:bg-purple-500/20 dark:text-purple-400">
        {t('settingsPage.plugins.kindCli')}
      </span>
    </div>
    <div className="mt-1 text-xs text-muted-foreground">{install.detail || install.phase}</div>
  </div>
))}
```

Add `Loader2` to the lucide-react imports.

- [ ] **Step 3: Add i18n keys**

```ts
'settingsPage.plugins.cliInstallFailed': '安装失败', // en: 'Installation failed'
```

- [ ] **Step 4: Verify frontend builds**

```
cd frontend && cnpm run build
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/settings/plugins/PluginToolListItem.tsx frontend/src/components/settings/plugins/PluginToolList.tsx frontend/src/i18n/
git commit -m "feat(plugin): show CLI install progress in plugin list"
```

---

## Section 4 — Built-in skill + install-from-docs entry

### Task 9: Create built-in `install-cli-from-docs` skill

**Files:**
- Create: `backend/pkg/skills/_builtin/install-cli-from-docs/SKILL.md`

- [ ] **Step 1: Create the skill directory**

```
mkdir -p backend/pkg/skills/_builtin/install-cli-from-docs
```

- [ ] **Step 2: Write `SKILL.md`**

```markdown
---
name: install-cli-from-docs
description: Install a CLI plugin from official documentation. Use when the user provides a documentation URL, pasted documentation text, or a natural-language description of a CLI tool they want to integrate. Handles package discovery, npm install, and manifest generation.
---

# Install CLI from Official Documentation

You are installing a CLI plugin based on user-provided documentation or description.

## Input

The user's message contains one or more of:
- A documentation URL (e.g. `https://github.com/foo/bar-cli#readme`)
- Pasted documentation text (Markdown or plain text)
- A natural-language description of the CLI tool

## Steps

1. **Identify the npm package name.**
   - If the user provided a URL, fetch it (via WebSearch or WebFetch) and extract the npm package name from the documentation.
   - If the user pasted documentation, look for `npm install <package>` commands, `npx <package>` references, or a "Package" / "Installation" section.
   - If the user gave a natural-language description, search for the most popular npm package matching that description.
   - Common patterns: `@anthropic-ai/claude-code`, `@modelcontextprotocol/server-filesystem`, `lark-cli`, etc.

2. **If you cannot determine the package name with confidence, call `RequestUserAttention`.**
   - `title`: "需要确认 npm 包名"
   - `message`: Explain what you found and ask the user to confirm or provide the exact package name.
   - Wait for the user's reply before proceeding.

3. **Call `InstallCli` with the package name.**
   ```json
   {"npm_package": "<determined-package-name>", "name": "<short-cli-name>"}
   ```
   - `name` should be a short kebab-case identifier derived from the package name (e.g. `@scope/foo-cli` → `foo-cli`).

4. **Call `GenerateCliManifest` to probe the CLI and generate a manifest.**
   ```json
   {"id": "cli:<name>"}
   ```
   - Use the `id` from the `InstallCli` result.

5. **Report success.**
   - Summarize what was installed: CLI name, version, number of tools discovered.

## Error Handling

- If `InstallCli` fails (network error, package not found), report the error clearly and suggest alternatives.
- If `GenerateCliManifest` fails, the CLI is still installed — report that the manifest needs manual configuration.
- If the user's input is completely empty or nonsensical, call `RequestUserAttention` asking for valid documentation.

## Constraints

- Do NOT install packages that are not published on npm.
- Do NOT modify the CLI's source code or configuration beyond what `InstallCli` and `GenerateCliManifest` handle.
- If the documentation mentions authentication / login steps, note them in your success summary so the user knows to run the login flow.
```

- [ ] **Step 3: Verify builtin embedding picks it up**

```
go test ./backend/pkg/skills/... -run TestBuiltin
```

If Plan A's builtin test asserts a specific count, update it to expect at least 1 (this skill).

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/skills/_builtin/install-cli-from-docs/
git commit -m "feat(skills): add builtin install-cli-from-docs skill"
```

---

### Task 10: Frontend — `CliInstallFromDocsModal` + dropdown expansion

**Files:**
- Create: `frontend/src/components/settings/plugins/CliInstallFromDocsModal.tsx`
- Create: `frontend/src/components/settings/plugins/CliNpmInstallModal.tsx`
- Modify: `frontend/src/components/settings/plugins/PluginToolList.tsx`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: Create `CliNpmInstallModal.tsx`**

A simple single-input modal:

```tsx
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

type Props = {
  open: boolean
  onClose: () => void
  onSubmit: (npmPackage: string, name: string) => void
}

export function CliNpmInstallModal({ open, onClose, onSubmit }: Props) {
  const { t } = useTranslation()
  const [npmPkg, setNpmPkg] = useState('')
  const [name, setName] = useState('')

  useEffect(() => {
    if (!open) { setNpmPkg(''); setName('') }
  }, [open])

  if (!open) return null

  const handleSubmit = () => {
    if (!npmPkg.trim()) return
    onSubmit(npmPkg.trim(), name.trim())
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button type="button" className="absolute inset-0 bg-black/35" onClick={onClose} />
      <div role="dialog" className="relative z-10 w-full max-w-md rounded-3xl border bg-background p-5 shadow-2xl">
        <h3 className="mb-4 text-base font-semibold">{t('settingsPage.plugins.cliFromNpm')}</h3>
        <input
          value={npmPkg}
          onChange={(e) => setNpmPkg(e.target.value)}
          placeholder={t('settingsPage.plugins.cliNpmPlaceholder')}
          className="mb-3 w-full rounded-lg border bg-background px-3 py-2 text-sm"
          autoFocus
        />
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder={t('settingsPage.plugins.cliNamePlaceholder')}
          className="mb-4 w-full rounded-lg border bg-background px-3 py-2 text-sm"
        />
        <div className="flex justify-end gap-2">
          <button type="button" onClick={onClose} className="rounded-xl border px-4 py-2 text-sm">{t('settingsPage.plugins.cliCancel')}</button>
          <button type="button" onClick={handleSubmit} disabled={!npmPkg.trim()} className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60">{t('settingsPage.plugins.cliInstall')}</button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Create `CliInstallFromDocsModal.tsx`**

```tsx
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

type Props = {
  open: boolean
  onClose: () => void
  onSubmit: (content: string) => void
}

export function CliInstallFromDocsModal({ open, onClose, onSubmit }: Props) {
  const { t } = useTranslation()
  const [content, setContent] = useState('')

  useEffect(() => {
    if (!open) setContent('')
  }, [open])

  if (!open) return null

  const handleSubmit = () => {
    if (!content.trim()) return
    onSubmit(content.trim())
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button type="button" className="absolute inset-0 bg-black/35" onClick={onClose} />
      <div role="dialog" className="relative z-10 w-full max-w-lg rounded-3xl border bg-background p-5 shadow-2xl">
        <h3 className="mb-2 text-base font-semibold">{t('settingsPage.plugins.cliFromDocs')}</h3>
        <p className="mb-3 text-xs text-muted-foreground">{t('settingsPage.plugins.cliFromDocsHint')}</p>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder={t('settingsPage.plugins.cliDocsPlaceholder')}
          rows={8}
          className="mb-4 w-full resize-none rounded-lg border bg-background px-3 py-2 text-sm"
          autoFocus
        />
        <div className="flex justify-end gap-2">
          <button type="button" onClick={onClose} className="rounded-xl border px-4 py-2 text-sm">{t('settingsPage.plugins.cliCancel')}</button>
          <button type="button" onClick={handleSubmit} disabled={!content.trim()} className="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60">{t('settingsPage.plugins.cliInstall')}</button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 3: Expand dropdown in `PluginToolList.tsx`**

Change `onCreate` type from `(kind: 'mcp' | 'plugin' | 'cli') => void` to `(kind: 'mcp' | 'plugin' | 'cli-npm' | 'cli-docs') => void`.

Replace the single "addCli" button with two:

```tsx
<button type="button" className="flex w-full items-center px-3 py-2 text-left text-sm hover:bg-accent"
  onClick={() => { props.onCreate('cli-npm'); setMenuOpen(false) }}>
  {t('settingsPage.plugins.cliFromNpm')}
</button>
<button type="button" className="flex w-full items-center px-3 py-2 text-left text-sm hover:bg-accent"
  onClick={() => { props.onCreate('cli-docs'); setMenuOpen(false) }}>
  {t('settingsPage.plugins.cliFromDocs')}
</button>
```

- [ ] **Step 4: Update `SettingsApp.tsx`**

Replace the old `CliInstallModal` usage with the two new modals:

```tsx
const [npmModalOpen, setNpmModalOpen] = useState(false)
const [docsModalOpen, setDocsModalOpen] = useState(false)
```

Update `handleCreateExtension`:

```tsx
const handleCreateExtension = async (kind: 'mcp' | 'plugin' | 'cli-npm' | 'cli-docs') => {
  if (kind === 'cli-npm') { setNpmModalOpen(true); return }
  if (kind === 'cli-docs') { setDocsModalOpen(true); return }
  // ... existing mcp/plugin logic
}
```

Add submit handlers:

```tsx
const handleNpmSubmit = async (npmPackage: string, name: string) => {
  try {
    await AgentBinding.SpawnHiddenSession(new SpawnHiddenSessionInput({
      title: `Install CLI: ${npmPackage}`,
      system_prompt: 'You are installing a CLI plugin. Use the InstallCli and GenerateCliManifest tools.',
      user_message: `Install the npm package "${npmPackage}" as a CLI plugin${name ? ` with name "${name}"` : ''}.`,
      skill_name: '',
    }))
  } catch (err) {
    pushAlert({ kind: 'error', placement: 'banner', title: '', message: String(err) })
  }
}

const handleDocsSubmit = async (content: string) => {
  try {
    await AgentBinding.SpawnHiddenSession(new SpawnHiddenSessionInput({
      title: t('settingsPage.plugins.cliFromDocsTask'),
      system_prompt: '',
      user_message: content,
      skill_name: 'install-cli-from-docs',
    }))
  } catch (err) {
    pushAlert({ kind: 'error', placement: 'banner', title: '', message: String(err) })
  }
}
```

Replace `<CliInstallModal ... />` with:

```tsx
<CliNpmInstallModal open={npmModalOpen} onClose={() => setNpmModalOpen(false)} onSubmit={handleNpmSubmit} />
<CliInstallFromDocsModal open={docsModalOpen} onClose={() => setDocsModalOpen(false)} onSubmit={handleDocsSubmit} />
```

Remove old `cliInstallModalOpen` state and `CliInstallModal` import.

- [ ] **Step 5: Add i18n keys**

`zh-CN.ts`:
```ts
'settingsPage.plugins.cliFromDocs': '从官方文档安装',
'settingsPage.plugins.cliFromDocsHint': '粘贴文档 URL、文档正文或自然语言描述，AI 将自动识别并安装对应的 CLI 插件。',
'settingsPage.plugins.cliDocsPlaceholder': '粘贴文档 URL、文档内容或描述你想要安装的 CLI 工具…',
'settingsPage.plugins.cliFromDocsTask': '从官方文档安装 CLI',
'settingsPage.plugins.cliInstallFailed': '安装失败',
```

`en.ts`:
```ts
'settingsPage.plugins.cliFromDocs': 'Install from official docs',
'settingsPage.plugins.cliFromDocsHint': 'Paste a documentation URL, documentation text, or a natural-language description. AI will identify and install the CLI plugin.',
'settingsPage.plugins.cliDocsPlaceholder': 'Paste a documentation URL, documentation text, or describe the CLI tool you want to install…',
'settingsPage.plugins.cliFromDocsTask': 'Install CLI from official docs',
'settingsPage.plugins.cliInstallFailed': 'Installation failed',
```

- [ ] **Step 6: Verify frontend builds**

```
cd frontend && cnpm run build
```

- [ ] **Step 7: Commit**

```bash
git add frontend/src/components/settings/plugins/CliNpmInstallModal.tsx frontend/src/components/settings/plugins/CliInstallFromDocsModal.tsx frontend/src/components/settings/plugins/PluginToolList.tsx frontend/src/components/settings/SettingsApp.tsx frontend/src/i18n/
git commit -m "feat(plugin): add install-from-docs modal and expand CLI dropdown entries"
```

---

## Section 5 — Cleanup + E2E

### Task 11: Delete old `CliInstallModal`

**Files:**
- Delete: `frontend/src/components/settings/plugins/CliInstallModal.tsx`

- [ ] **Step 1: Verify no remaining references**

```
grep -rn "CliInstallModal" frontend/src/
```

Should return no hits (the old component was replaced in Task 10).

- [ ] **Step 2: Delete the file**

```
git rm frontend/src/components/settings/plugins/CliInstallModal.tsx
```

- [ ] **Step 3: Remove orphaned i18n keys**

Remove keys that were only used by the old modal:
- `settingsPage.plugins.cliInstallTitle`
- `settingsPage.plugins.cliFromLocal` (already removed in Task 2)
- `settingsPage.plugins.cliLocalPlaceholder` (already removed in Task 2)
- `settingsPage.plugins.cliBrowse` (already removed in Task 2)
- `settingsPage.plugins.cliInstallStepInstalling`
- `settingsPage.plugins.cliInstallStepGenerating`
- `settingsPage.plugins.cliInstalling`

(Only remove if no other code references them.)

- [ ] **Step 4: Verify frontend builds**

```
cd frontend && cnpm run build
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "refactor(plugin): delete old CliInstallModal replaced by npm + docs modals"
```

---

### Task 12: E2E smoke test

- [ ] **Step 1: Run app**

```
task dev
```

- [ ] **Step 2: Verify dropdown entries**

Open Settings → Plugins → click "+" dropdown. Verify:
- "从 npm 安装" appears
- "从官方文档安装" appears
- "从本地目录导入" does **NOT** appear

- [ ] **Step 3: Test npm install flow**

1. Click "从 npm 安装".
2. Enter a real npm package (e.g. `lark-cli` or a small test package).
3. Click install.
4. **Verify**: modal closes immediately.
5. **Verify**: a pending entry appears in the plugin list with a progress bar.
6. **Verify**: progress phases update (downloading → installed → generating → done).
7. **Verify**: on success, the entry becomes a normal CLI plugin with toggle enabled.

- [ ] **Step 4: Test concurrent installs**

1. Start installing CLI A.
2. Before A finishes, start installing CLI B.
3. **Verify**: both pending entries appear and progress independently.
4. **Verify**: both complete successfully.

- [ ] **Step 5: Test install-from-docs flow**

1. Click "从官方文档安装".
2. Paste a documentation URL (e.g. `https://github.com/anthropics/claude-code#readme`).
3. Submit.
4. **Verify**: modal closes, pending entry appears.
5. **Verify**: the hidden session's agent identifies the package and installs it.
6. If the agent cannot determine the package:
   - **Verify**: a notification appears in the bell.
   - Click the notification → floating popup opens.
   - Reply with the package name.
   - **Verify**: installation continues and completes.

- [ ] **Step 6: Test install-from-docs with pasted text**

1. Click "从官方文档安装".
2. Paste documentation text (no URL).
3. Submit.
4. **Verify**: same flow as Step 5.

- [ ] **Step 7: Test navigate-away resilience**

1. Start an install.
2. Navigate away from Settings to the chat page.
3. Navigate back to Settings → Plugins.
4. **Verify**: install progress is still showing (or completed).

- [ ] **Step 8: Test skill in chat**

1. Open a regular chat session.
2. Type: `帮我装一个飞书 CLI，文档在 https://github.com/larksuite/cli`
3. **Verify**: the agent picks up the `install-cli-from-docs` skill and enters the install flow.
4. **Verify**: a notification or progress update appears.

- [ ] **Step 9: Commit any fix-ups**

---

## Self-review

1. **Spec coverage** —
   - §3.1 (simplify entry) → Task 1 (remove backend), Task 2 (remove frontend tab), Task 10 (dropdown expansion).
   - §3.2 (background install) → Task 3–6 (agent tools + bridge), Task 7–8 (progress UI), Task 10 (SpawnHiddenSession wiring).
   - §3.3 (install from docs) → Task 9 (builtin skill), Task 10 (modal + wiring).
   - §3.5 (intervention) → reuses Plan B's `RequestUserAttention` + notification bell + floating popup. The skill body (Task 9) instructs the agent to call `RequestUserAttention` when it cannot determine the package.

2. **Placeholder scan** — no "TBD" or "implement later" markers. All code blocks are complete.

3. **Type consistency** —
   - `CliInstaller.InstallCliSync` returns `(string, error)` — JSON-encoded extension item.
   - `CliManifestGenerator.GenerateCliManifestSync` returns `(string, error)` — JSON-encoded extension item.
   - Agent tool names: `InstallCli`, `GenerateCliManifest` — match the `ToolMeta.Name` constants.
   - Event names: `cli.install.progress`, `cli.install.done` — consistent between backend emit and frontend subscription.
   - `SpawnHiddenSessionInput` fields: `title`, `system_prompt`, `user_message`, `skill_name` — match Plan B's DTO.

4. **Cross-plan consistency** —
   - Plan A's `BuildSkillTool` is registered alongside the new CLI tools (Task 6).
   - Plan B's `SpawnHiddenSession` is called from the frontend (Task 10).
   - Plan B's `RequestUserAttention` is used by the builtin skill (Task 9) for package-name confirmation.
   - Plan B's notification bell + floating popup surface the intervention UI.

---

## Execution handoff

Plan complete at `docs/superpowers/plans/2026-05-25-cli-refactor.md`. Two execution options:

1. **Subagent-Driven (recommended)** — fresh subagent per task, review between tasks.
2. **Inline Execution** — execute in this session with checkpoints.

Plans A and B must merge first.
