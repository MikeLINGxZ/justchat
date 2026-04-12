# Progressive Skill Injection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Change skill injection from all-at-once system prompt injection to progressive on-demand loading via a `load_skill` tool, and add `when` field to skill metadata.

**Architecture:** Create a `load_skill` builtin tool that LLM can call to fetch skill content on demand. Inject only skill summaries (name + description + when) into agent system prompts. Update skill metadata to include a `when` (trigger condition) field. Add folder import support to the skill settings UI.

**Tech Stack:** Go (backend), React + Ant Design (frontend), Wails v3 (desktop bridge), Eino ADK (LLM framework)

---

### Task 1: Extend Skill Metadata with `when` Field

**Files:**
- Modify: `backend/pkg/skills/skills.go:14-19`
- Modify: `backend/models/view_models/skill.go:1-13`
- Modify: `backend/service/skill_settings.go:28-34,50-57,128-142`

- [ ] **Step 1: Add `when` field to `SkillMeta` struct**

In `backend/pkg/skills/skills.go`, add `When` field to `SkillMeta`:

```go
type SkillMeta struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	When        string   `yaml:"when" json:"when"`
	Version     string   `yaml:"version" json:"version"`
	Tags        []string `yaml:"tags" json:"tags"`
}
```

- [ ] **Step 2: Add `when` field to view models**

In `backend/models/view_models/skill.go`:

```go
type SkillSummary struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	When        string   `json:"when"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}
```

- [ ] **Step 3: Update service layer to pass `when` field through**

In `backend/service/skill_settings.go`, update `ListSkills()` result mapping (~line 28) to include `When: m.When`. Update `GetSkill()` result (~line 51) to include `When: skill.When`. Update `viewModelToSkill()` (~line 134) to include `When: strings.TrimSpace(input.When)`.

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/skills/skills.go backend/models/view_models/skill.go backend/service/skill_settings.go
git commit -m "feat(skills): add 'when' trigger condition field to skill metadata"
```

---

### Task 2: Add Skill Summary Resolver for Progressive Injection

**Files:**
- Modify: `backend/pkg/skills/resolver.go`

- [ ] **Step 1: Add `ResolveSkillSummaries` function**

Replace the content of `backend/pkg/skills/resolver.go` with both the existing function and a new summary function:

```go
package skills

import "strings"

// ResolveSkillContents loads the specified skills and concatenates their content
// into a structured prompt section for injection into agent instructions.
func ResolveSkillContents(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}

	var sections []string
	for _, name := range skillNames {
		skill, err := LoadSkill(name)
		if err != nil {
			continue
		}

		content := strings.TrimSpace(skill.Content)
		if content == "" {
			continue
		}

		sections = append(sections, "### "+skill.Name+"\n"+content)
	}

	if len(sections) == 0 {
		return ""
	}

	return "## Skills\n\n" + strings.Join(sections, "\n\n")
}

// ResolveSkillSummaries builds a prompt section listing available skills
// with their name, description, and trigger condition, instructing the LLM
// to call load_skill when needed.
func ResolveSkillSummaries(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}

	var lines []string
	for _, name := range skillNames {
		skill, err := LoadSkill(name)
		if err != nil {
			continue
		}

		desc := strings.TrimSpace(skill.Description)
		when := strings.TrimSpace(skill.When)
		if desc == "" {
			continue
		}

		line := "- **" + skill.Name + "**: " + desc
		if when != "" {
			line += " (when: " + when + ")"
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return ""
	}

	return "## Available Skills\n\n" +
		"You have the following skills available. When a user's request matches a skill's trigger condition, " +
		"call the `load_skill` tool with the skill name to load its full content before responding.\n\n" +
		strings.Join(lines, "\n")
}

// ResolveAllSkillSummaries loads all skills from disk and builds a summary prompt section.
func ResolveAllSkillSummaries() string {
	metas, err := ListSkills()
	if err != nil || len(metas) == 0 {
		return ""
	}

	var names []string
	for _, m := range metas {
		names = append(names, m.Name)
	}
	return ResolveSkillSummaries(names)
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/pkg/skills/resolver.go
git commit -m "feat(skills): add skill summary resolver for progressive injection"
```

---

### Task 3: Create `load_skill` Builtin Tool

**Files:**
- Create: `backend/pkg/llm_provider/tools/skill_tool.go`
- Modify: `backend/pkg/llm_provider/tools/common.go:20-30` (register the tool)
- Modify: `backend/pkg/i18n/resources_zh_cn.go` (add i18n keys)
- Modify: `backend/pkg/i18n/resources_en_us.go` (add i18n keys)

- [ ] **Step 1: Add i18n keys for the skill tool**

In `backend/pkg/i18n/resources_zh_cn.go`, add after the `"tool.shell.description"` line:

```go
"tool.load_skill.name":        "加载技能",
"tool.load_skill.description": "按名称加载技能的完整内容。当用户请求匹配某个可用技能的触发条件时，调用此工具获取完整指令后再回复。",
```

In `backend/pkg/i18n/resources_en_us.go`, add after the `"tool.shell.description"` line:

```go
"tool.load_skill.name":        "Load Skill",
"tool.load_skill.description": "Load the full content of a skill by name. Call this tool when the user's request matches an available skill's trigger condition, then follow the loaded instructions.",
```

- [ ] **Step 2: Create the skill tool implementation**

Create `backend/pkg/llm_provider/tools/skill_tool.go`:

```go
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

type LoadSkillTool struct{}

type loadSkillParams struct {
	SkillName string `json:"skill_name"`
}

func (l *LoadSkillTool) Id() string {
	return "load_skill"
}

func (l *LoadSkillTool) Name() string {
	return i18n.TCurrent("tool.load_skill.name", nil)
}

func (l *LoadSkillTool) Description() string {
	return i18n.TCurrent("tool.load_skill.description", nil)
}

func (l *LoadSkillTool) RequireConfirmation() bool { return false }

func (l *LoadSkillTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "load_skill",
			Desc: i18n.TCurrent("tool.load_skill.description", nil),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"skill_name": {
					Type:     schema.String,
					Desc:     "The name of the skill to load.",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params loadSkillParams) (string, error) {
			name := strings.TrimSpace(params.SkillName)
			if name == "" {
				return "", fmt.Errorf("skill_name is required")
			}

			skill, err := skills.LoadSkill(name)
			if err != nil {
				return "", fmt.Errorf("skill not found: %s", name)
			}

			content := strings.TrimSpace(skill.Content)
			if content == "" {
				return "", fmt.Errorf("skill %s has no content", name)
			}

			return fmt.Sprintf("# Skill: %s\n\n%s", skill.Name, content), nil
		},
	)
}
```

- [ ] **Step 3: Register the tool in the router**

In `backend/pkg/llm_provider/tools/common.go`, add to the `init()` function after the existing registrations:

```go
ToolRouter.RegisterTool(&LoadSkillTool{})
```

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/llm_provider/tools/skill_tool.go backend/pkg/llm_provider/tools/common.go backend/pkg/i18n/resources_zh_cn.go backend/pkg/i18n/resources_en_us.go
git commit -m "feat(skills): create load_skill builtin tool for progressive injection"
```

---

### Task 4: Modify Chat Completions to Use Progressive Skill Injection

**Files:**
- Modify: `backend/service/chat.go:144-197,311-312`
- Modify: `backend/pkg/llm_provider/provider.go:64-77`

- [ ] **Step 1: Inject skill summaries into main agent**

In `backend/service/chat.go`, after line 311 where `directTools` is built, inject skill summaries into the main agent's prompt. The key changes:

1. Before creating the provider (~line 282), build the skill summary for the main agent:

```go
// 构建主 Agent 的 skill 摘要（渐进注入）
mainSkillSummary := skills.ResolveAllSkillSummaries()
```

2. Pass the skill summary to `NewLlmProvider`. Update the call at ~line 312:

```go
provider, err := llm_provider.NewLlmProvider(ctx, *providerModel, subAgents, directTools, toolMiddleware, localizedPrompts, mainSkillSummary)
```

- [ ] **Step 2: Update `NewLlmProvider` to accept and inject skill summary**

In `backend/pkg/llm_provider/provider.go`, update the `NewLlmProvider` function signature and pass the skill summary to the main agent instruction:

```go
func NewLlmProvider(ctx context.Context, providerModel wrapper_models.ProviderModel, subAgents []adk.Agent, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware, promptSet prompts.PromptSet, skillSummary string) (*Provider, error) {
	chatModel, err := NewToolCallingChatModel(ctx, providerModel)
	if err != nil {
		return nil, err
	}

	instruction := promptSet.MainAgentSystem
	if skillSummary != "" {
		instruction = instruction + "\n\n" + skillSummary
	}

	mainAgent, err := agents.NewMainAgent(ctx, chatModel, subAgents, tools, toolMiddleware, instruction)
	if err != nil {
		return nil, err
	}

	return &Provider{chatModel: chatModel, toolChatModel: chatModel, tools: tools, mainAgent: mainAgent, prompts: promptSet}, nil
}
```

- [ ] **Step 3: Change custom agent skill injection from full content to summary**

In `backend/service/chat.go`, modify the custom agent skill injection block (~lines 175-182). Replace the full content injection with summary injection:

```go
// 解析 skill 摘要并注入 prompt（渐进注入）
instruction := customDef.PromptText
if len(customDef.SkillIDs) > 0 {
	skillSummary := skills.ResolveSkillSummaries(customDef.SkillIDs)
	if skillSummary != "" {
		instruction = instruction + "\n\n" + skillSummary
	}
}
```

- [ ] **Step 4: Ensure `load_skill` tool is always available when skills exist**

The `load_skill` tool is now a builtin registered in `ToolRouter`. Since builtin tools are hidden from the tool selection UI but always enabled (per recent commit `0ea753d`), the `load_skill` tool will be automatically available to all agents.

Verify this by checking that `chat.go` line 145 calls `resolveSelectedTools` which fetches user-selected tools, and that the builtin `load_skill` is included in the tools passed to the agent. If builtin tools are not auto-included, we need to add the `load_skill` tool explicitly to `directTools`.

Check how builtin tools are currently added. Looking at line 311:
```go
directTools := append([]tool.BaseTool{newWorkflowHandoffTool(runner.setWorkflowHandoff)}, agentTools...)
```

The `agentTools` come from `resolveSelectedTools` which only resolves user-selected tool IDs. Builtin tools like `get_current_date` are also in the user-selected tools list. So `load_skill` needs to be explicitly added to `directTools` to ensure it's always available:

```go
directTools := append([]tool.BaseTool{newWorkflowHandoffTool(runner.setWorkflowHandoff)}, agentTools...)
// 始终添加 load_skill 工具以支持渐进式技能注入
if loadSkillTool, ok := llmtools.ToolRouter.GetToolByID("load_skill"); ok {
	directTools = append(directTools, loadSkillTool.Tool())
}
```

Similarly, for custom agents, add the `load_skill` tool to their tools if they have skills:

In the custom agent prep loop, after resolving custom tools (~line 166), if the agent has skills, append the load_skill tool:

```go
if len(customDef.SkillIDs) > 0 {
	if loadSkillTool, ok := llmtools.ToolRouter.GetToolByID("load_skill"); ok {
		customTools = append(customTools, loadSkillTool.Tool())
	}
}
```

- [ ] **Step 5: Commit**

```bash
git add backend/service/chat.go backend/pkg/llm_provider/provider.go
git commit -m "feat(skills): switch to progressive skill injection via load_skill tool"
```

---

### Task 5: Add `when` Field to Frontend Skill Settings UI

**Files:**
- Modify: `frontend/src/pages/settings/skills/index.tsx:369-403,469-503`
- Modify: `frontend/src/i18n/resources/zh-CN.ts` (skills section)
- Modify: `frontend/src/i18n/resources/en-US.ts` (skills section)

- [ ] **Step 1: Add i18n keys for `when` field**

In `frontend/src/i18n/resources/zh-CN.ts`, inside the `settings.skills.form` object, add:

```typescript
when: '触发条件',
whenPlaceholder: '描述何时应加载此技能，例如：用户要求翻译内容时',
```

In `frontend/src/i18n/resources/en-US.ts`, inside the `settings.skills.form` object, add:

```typescript
when: 'Trigger Condition',
whenPlaceholder: 'Describe when this skill should be loaded, e.g.: when user asks for translation',
```

Also update the `tip` in both files:
- zh-CN: `'技能将按需加载：系统仅向模型展示技能摘要，模型在需要时调用工具获取完整内容。'`
- en-US: `'Skills are loaded progressively: the model sees only summaries and loads full content on demand via tool call.'`

- [ ] **Step 2: Add `when` field to the create form**

In `frontend/src/pages/settings/skills/index.tsx`, in the `renderCreateModal` function, add a `when` form item after the `description` field:

```tsx
<Form.Item
  name="when"
  label={t('settings.skills.form.when')}
>
  <Input placeholder={t('settings.skills.form.whenPlaceholder')} />
</Form.Item>
```

- [ ] **Step 3: Add `when` field to the detail editor display**

In the `renderEditorBody` function, after the description paragraph (~line 385), add:

```tsx
{detail.when && (
  <Paragraph type="secondary" style={{ marginBottom: 4 }}>
    {t('settings.skills.form.when')}: {detail.when}
  </Paragraph>
)}
```

- [ ] **Step 4: Pass `when` field in create and save operations**

In `handleCreate` (~line 244), include `when` in the `SkillDetail` constructor:

```typescript
const input = new SkillDetail({
  name: values.name,
  description: values.description,
  when: values.when || '',
  version: values.version || '1.0',
  tags: values.tags || [],
  content: values.content,
});
```

In `handleSave` (~line 179), include `when`:

```typescript
const input = new SkillDetail({
  name: detail.name,
  description: detail.description,
  when: detail.when || '',
  version: detail.version,
  tags: detail.tags,
  content,
});
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/pages/settings/skills/index.tsx frontend/src/i18n/resources/zh-CN.ts frontend/src/i18n/resources/en-US.ts
git commit -m "feat(skills): add 'when' trigger condition field to skill settings UI"
```

---

### Task 6: Add Folder Import for Skills

**Files:**
- Modify: `backend/service/skill_settings.go` (add `ImportSkillsFromFolder` method)
- Modify: `backend/pkg/i18n/resources_zh_cn.go` (add dialog i18n key)
- Modify: `backend/pkg/i18n/resources_en_us.go` (add dialog i18n key)
- Modify: `frontend/src/pages/settings/skills/index.tsx` (add import button and flow)
- Modify: `frontend/src/i18n/resources/zh-CN.ts` (add import i18n keys)
- Modify: `frontend/src/i18n/resources/en-US.ts` (add import i18n keys)

- [ ] **Step 1: Add backend `SelectSkillFolder` and `ImportSkillsFromFolder` methods**

In `backend/service/skill_settings.go`, add:

```go
// SelectSkillFolder opens a folder selection dialog for importing skills.
func (s *Service) SelectSkillFolder() (string, error) {
	path, err := s.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle(i18n.TCurrent("app.dialog.select_skill_folder", nil)).
		PromptForSingleSelection()
	if err != nil {
		return "", ierror.NewError(err)
	}
	return path, nil
}

// ImportSkillsFromFolder scans a folder for .md skill files and imports them.
func (s *Service) ImportSkillsFromFolder(folderPath string) ([]view_models.SkillSummary, error) {
	if folderPath == "" {
		return nil, ierror.NewError(fmt.Errorf("folder path is empty"))
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var imported []view_models.SkillSummary
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(folderPath, entry.Name()))
		if err != nil {
			continue
		}

		// Try parsing as a skill file with frontmatter
		skill, err := skills.ParseSkillContent(raw)
		if err != nil {
			continue
		}

		// Skip if name is invalid or already exists
		if !skillNamePattern.MatchString(skill.Name) || skills.SkillExists(skill.Name) {
			continue
		}

		if err := skills.SaveSkill(*skill); err != nil {
			continue
		}

		tags := skill.Tags
		if tags == nil {
			tags = []string{}
		}
		imported = append(imported, view_models.SkillSummary{
			Name:        skill.Name,
			Description: skill.Description,
			When:        skill.When,
			Version:     skill.Version,
			Tags:        tags,
		})
	}

	return imported, nil
}
```

Add the necessary imports: `"os"`, `"path/filepath"`.

- [ ] **Step 2: Add `ParseSkillContent` helper to the skills package**

In `backend/pkg/skills/skills.go`, add a public wrapper around `parseFrontmatter`:

```go
// ParseSkillContent parses raw bytes as a skill file (frontmatter + body).
func ParseSkillContent(raw []byte) (*Skill, error) {
	meta, body, err := parseFrontmatter(raw)
	if err != nil {
		return nil, err
	}
	return &Skill{
		SkillMeta: meta,
		Content:   body,
	}, nil
}
```

- [ ] **Step 3: Add i18n keys for folder dialog**

In `backend/pkg/i18n/resources_zh_cn.go`:
```go
"app.dialog.select_skill_folder": "选择技能文件夹",
```

In `backend/pkg/i18n/resources_en_us.go`:
```go
"app.dialog.select_skill_folder": "Select Skills Folder",
```

- [ ] **Step 4: Add frontend i18n keys for import**

In `frontend/src/i18n/resources/zh-CN.ts`, inside `settings.skills`, add:

```typescript
importFromFolder: '从文件夹导入',
importSuccess: '成功导入 {{count}} 个技能',
importEmpty: '未找到可导入的技能文件',
```

In `frontend/src/i18n/resources/en-US.ts`, inside `settings.skills`, add:

```typescript
importFromFolder: 'Import from Folder',
importSuccess: 'Successfully imported {{count}} skills',
importEmpty: 'No importable skill files found',
```

- [ ] **Step 5: Update the create button to a dropdown with two options**

In `frontend/src/pages/settings/skills/index.tsx`:

Add `Dropdown` to the antd imports. Add `FolderOpenOutlined` to the icon imports.

Replace the create button in `renderSkillList` (the `<Button>` inside the card title) with a Dropdown:

```tsx
<Dropdown
  menu={{
    items: [
      {
        key: 'create',
        icon: <PlusOutlined />,
        label: t('settings.skills.actions.create'),
        onClick: () => setCreateModalOpen(true),
      },
      {
        key: 'import',
        icon: <FolderOpenOutlined />,
        label: t('settings.skills.importFromFolder'),
        onClick: () => void handleImportFromFolder(),
      },
    ],
  }}
  trigger={['click']}
>
  <Button type="text" size="small" icon={<PlusOutlined />} />
</Dropdown>
```

Add the import handler function:

```typescript
const handleImportFromFolder = async () => {
  try {
    const folderPath = await Service.SelectSkillFolder();
    if (!folderPath) return;

    const imported = await Service.ImportSkillsFromFolder(folderPath);
    if (imported && imported.length > 0) {
      message.success(t('settings.skills.importSuccess', { count: imported.length }));
      await refreshList(imported[0].name);
      setActiveName(imported[0].name);
    } else {
      message.info(t('settings.skills.importEmpty'));
    }
  } catch (error) {
    console.error('导入技能失败:', error);
    message.error(t('settings.skills.createFailed'));
  }
};
```

- [ ] **Step 6: Commit**

```bash
git add backend/service/skill_settings.go backend/pkg/skills/skills.go backend/pkg/i18n/resources_zh_cn.go backend/pkg/i18n/resources_en_us.go frontend/src/pages/settings/skills/index.tsx frontend/src/i18n/resources/zh-CN.ts frontend/src/i18n/resources/en-US.ts
git commit -m "feat(skills): add folder import support for skill files"
```

---

### Task 7: Regenerate Wails Bindings and Verify Build

**Files:**
- Modified bindings will be auto-generated

- [ ] **Step 1: Regenerate Wails bindings**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
wails3 generate bindings
```

This regenerates the TypeScript bindings so that the frontend can see the new `When` field on `SkillSummary`/`SkillDetail` and the new `SelectSkillFolder`/`ImportSkillsFromFolder` service methods.

- [ ] **Step 2: Verify Go build**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
go build ./...
```

- [ ] **Step 3: Verify frontend build**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npm run build
```

- [ ] **Step 4: Commit bindings if changed**

```bash
git add frontend/bindings/
git commit -m "chore: regenerate Wails bindings for skill changes"
```
