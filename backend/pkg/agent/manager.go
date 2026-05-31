package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/prompt_id"
	pkgMCP "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/mcp"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompt"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	pkgRuntimeState "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/runtime_state"
	pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
	agenttool "trpc.group/trpc-go/trpc-agent-go/tool"
	toolpkg "trpc.group/trpc-go/trpc-agent-go/tool"
)

const defaultMainAgentPrompt = `You are Lemontea, a helpful AI assistant. You can help users with various tasks including answering questions, writing, coding, and more. Be concise, accurate, and helpful. When using tools, explain what you are doing and why.`
const defaultGenChatTitlePrompt = `Generate a concise title under 20 characters for this conversation. Return only the title.`

// Manager manages runners, prompts, and tool registration for the main agent.
type Manager struct {
	mu                   sync.RWMutex
	istorage             *storage.Storage
	wailsApp             *application.App
	toolRegistry         *tools.Registry
	streamManager        *StreamManager
	runners              map[string]runner.Runner
	mcpManager           *pkgMCP.Manager
	skillProvider        tools.SkillProvider
	attention            tools.AttentionRequester
	skillCreator         tools.SkillCreator
	cliInstaller         tools.CliInstaller
	cliInstallProgress   tools.CliInstallProgressReporter
	cliManifestGenerator tools.CliManifestGenerator
	cliCommandRunner     tools.CliCommandRunner
	terminalRunner       *pkgterminal.Manager
	toolConfirmDecider   ToolConfirmDecider
	memoryEncodeSem      chan struct{}
}

type toolResultEmitter struct {
	app *application.App
}

func (e toolResultEmitter) EmitToolResult(sessionID uint, toolName, result string) {
	if e.app == nil {
		return
	}
	e.app.Event.Emit(event_id.AgentStreamToolResult, map[string]any{
		"sessionId": sessionID,
		"toolName":  toolName,
		"result":    result,
	})
}

// NewManager creates a new manager and registers the default main-agent prompt/tools.
func NewManager(istorage *storage.Storage) *Manager {
	prompt.Register(prompt_id.MainAgent, defaultMainAgentPrompt)
	prompt.Register(prompt_id.GenChatTitle, defaultGenChatTitlePrompt)

	registry := tools.NewRegistry()
	registry.Register(tools.DateTimeMeta())
	registry.Register(tools.FileReadMeta())
	registry.Register(tools.FileWriteMeta())
	registry.Register(tools.ShellMeta())
	registry.Register(tools.BuildQuestionTool())
	registry.Register(tools.WebFetchMeta())
	registry.Register(tools.BuildWebSearchTool())
	registry.Register(tools.CodeExecMeta())
	registry.Register(tools.BuildQRCodeTool())
	registry.Register(tools.BuildInteractiveTerminalTool())
	registry.Register(tools.BuildRequestAttentionTool())
	registry.Register(tools.BuildProposeSkillTool())
	registry.Register(tools.BuildInstallCliTool())
	registry.Register(tools.BuildReportCliInstallProgressTool())
	registry.Register(tools.BuildGenerateCliManifestTool())
	registry.Register(tools.BuildRunCliCommandTool())
	registry.Register(tools.BuildSaveTaskStateTool())
	registry.Register(tools.BuildLoadTaskStateTool())
	registry.Register(tools.BuildTodoWriteTool())
	// Register a static Skill meta for discovery; buildAgentTools() builds
	// the dynamic description from the live provider on each turn.
	registry.Register(tools.ToolMeta{
		Name:        tools.SkillToolName,
		Description: "Load a skill's instructions by name",
		Category:    tools.CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				Name string `json:"name"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Loading skill: " + parsed.Name
		},
	})

	return &Manager{
		istorage:           istorage,
		toolRegistry:       registry,
		streamManager:      NewStreamManager(),
		runners:            make(map[string]runner.Runner),
		mcpManager:         pkgMCP.NewManager(mustDataDir()),
		toolConfirmDecider: defaultToolConfirmDecider,
		memoryEncodeSem:    make(chan struct{}, 2),
	}
}

// SetApp stores the Wails app reference during service startup.
func (m *Manager) SetApp(app *application.App) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wailsApp = app
}

// ToolRegistry exposes the registered tool metadata.
func (m *Manager) ToolRegistry() *tools.Registry {
	return m.toolRegistry
}

// Streams exposes the active stream manager.
func (m *Manager) Streams() *StreamManager {
	return m.streamManager
}

// Storage exposes the backing storage.
func (m *Manager) Storage() *storage.Storage {
	return m.istorage
}

// App returns the Wails application handle.
func (m *Manager) App() *application.App {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wailsApp
}

// SetSkillProvider stores the skill provider used by the Skill meta-tool.
// It must be called before any chat messages are processed.
func (m *Manager) SetSkillProvider(provider tools.SkillProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.skillProvider = provider
}

// SkillProvider returns the currently configured skill provider, or nil.
func (m *Manager) SkillProvider() tools.SkillProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.skillProvider
}

// SetAttentionRequester stores the notification-backed attention requester.
func (m *Manager) SetAttentionRequester(requester tools.AttentionRequester) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attention = requester
}

// SetSkillCreator stores the persistence layer used by ProposeSkill.
func (m *Manager) SetSkillCreator(creator tools.SkillCreator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.skillCreator = creator
}

// SetCliInstaller stores the CLI installer used by the InstallCli tool.
func (m *Manager) SetCliInstaller(i tools.CliInstaller) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cliInstaller = i
}

// SetCliInstallProgressReporter stores the progress reporter used by task install flows.
func (m *Manager) SetCliInstallProgressReporter(r tools.CliInstallProgressReporter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cliInstallProgress = r
}

// SetCliManifestGenerator stores the manifest generator used by the GenerateCliManifest tool.
func (m *Manager) SetCliManifestGenerator(g tools.CliManifestGenerator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cliManifestGenerator = g
}

// SetCliCommandRunner stores the arbitrary CLI command runner used by the RunCliCommand tool.
func (m *Manager) SetCliCommandRunner(r tools.CliCommandRunner) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cliCommandRunner = r
}

func (m *Manager) SetTerminalRunner(r *pkgterminal.Manager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.terminalRunner = r
}

// SetToolConfirmDecider overrides the runtime approval decider used before tool execution.
func (m *Manager) SetToolConfirmDecider(decider ToolConfirmDecider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.toolConfirmDecider = decider
}

func (m *Manager) buildAgentTools(enabledUserTools []string, sessionKind string, sessionID uint) []toolpkg.Tool {
	metas := m.toolRegistry.EnabledTools(enabledUserTools)
	agentTools := make([]toolpkg.Tool, 0, len(metas))

	for _, meta := range metas {
		switch meta.Name {
		case "datetime":
			agentTools = append(agentTools, tools.NewDateTimeTool())
		case "file_read":
			agentTools = append(agentTools, tools.NewFileReadTool())
		case "file_write":
			agentTools = append(agentTools, tools.NewFileWriteTool())
		case "shell":
			if m.terminalRunner != nil {
				agentTools = append(agentTools, tools.NewShellTool(toolResultEmitter{app: m.wailsApp}, sessionID, m.terminalRunner))
			} else {
				agentTools = append(agentTools, tools.NewShellTool(toolResultEmitter{app: m.wailsApp}, sessionID))
			}
		case tools.QuestionToolName:
			agentTools = append(agentTools, tools.NewQuestionTool())
		case "web_fetch":
			agentTools = append(agentTools, tools.NewWebFetchTool())
		case tools.WebSearchToolName:
			agentTools = append(agentTools, tools.NewWebSearchTool())
		case "code_exec":
			agentTools = append(agentTools, tools.NewCodeExecTool())
		case tools.QRCodeToolName:
			agentTools = append(agentTools, tools.NewQRCodeTool())
		case tools.InteractiveTerminalToolName:
			agentTools = append(agentTools, tools.NewInteractiveTerminalTool())
		case tools.SkillToolName:
			// The Skill meta-tool is only included when a provider has been
			// configured so the dynamic description can be built from the
			// current set of enabled skills.
			// Note: m.skillProvider is accessed directly here because
			// buildAgentTools is always called while m.mu is held by the caller.
			if m.skillProvider != nil {
				agentTools = append(agentTools, tools.NewSkillTool(m.skillProvider))
			}
		case tools.RequestAttentionToolName:
			if sessionKind == "task" && m.attention != nil {
				agentTools = append(agentTools, tools.NewRequestAttentionTool(m.attention, sessionID))
			}
		case tools.ProposeSkillToolName:
			if sessionKind == "task" && m.attention != nil && m.skillCreator != nil {
				agentTools = append(agentTools, tools.NewProposeSkillTool(m.attention, m.skillCreator, sessionID))
			}
		case tools.SaveTaskStateToolName:
			if sessionKind == "task" {
				agentTools = append(agentTools, tools.NewSaveTaskStateTool(m.istorage, sessionID))
			}
		case tools.LoadTaskStateToolName:
			if sessionKind == "task" {
				agentTools = append(agentTools, tools.NewLoadTaskStateTool(m.istorage, sessionID))
			}
		case tools.TodoWriteToolName:
			agentTools = append(agentTools, tools.NewTodoWriteTool())
		case tools.InstallCliToolName:
			if m.cliInstaller != nil {
				agentTools = append(agentTools, tools.NewInstallCliTool(m.cliInstaller))
			}
		case tools.ReportCliInstallProgressToolName:
			if sessionKind == "task" && m.cliInstallProgress != nil {
				agentTools = append(agentTools, tools.NewReportCliInstallProgressTool(m.cliInstallProgress, sessionID))
			}
		case tools.GenerateCliManifestToolName:
			if m.cliManifestGenerator != nil {
				agentTools = append(agentTools, tools.NewGenerateCliManifestTool(m.cliManifestGenerator))
			}
		case tools.RunCliCommandToolName:
			if m.cliCommandRunner != nil {
				agentTools = append(agentTools, tools.NewRunCliCommandTool(m.cliCommandRunner, toolResultEmitter{app: m.wailsApp}, sessionID))
			}
		}
	}

	return agentTools
}

func (m *Manager) buildMCPToolSets(enabledUserTools []string) []agenttool.ToolSet {
	if len(enabledUserTools) == 0 {
		return nil
	}

	config, err := loadAgentConfig()
	if err != nil || len(config.Extensions) == 0 {
		return nil
	}

	toolSets := make([]agenttool.ToolSet, 0)
	for _, item := range config.Extensions {
		if item.Kind != "mcp" || !item.Enabled {
			continue
		}
		toolSet, buildErr := m.mcpManager.BuildToolSet(item, enabledUserTools)
		if buildErr != nil {
			continue
		}
		if toolSet.Name() != "" {
			toolSets = append(toolSets, toolSet)
		}
	}
	return toolSets
}

func (m *Manager) buildCliToolSets(enabledUserTools []string) []agenttool.ToolSet {
	if len(enabledUserTools) == 0 {
		return nil
	}

	config, err := loadAgentConfig()
	if err != nil || len(config.Extensions) == 0 {
		return nil
	}
	state, err := pkgRuntimeState.LoadPersistedState()
	if err != nil || state.State != "ready" || state.NodePath == "" || state.NpmPath == "" {
		return nil
	}
	mgr := pkgcli.NewManager(mustDataDir(), state.NodePath, state.NpmPath)

	toolSets := make([]agenttool.ToolSet, 0)
	for _, item := range config.Extensions {
		if item.Kind != "cli" || !item.Enabled {
			continue
		}
		manifest, loadErr := pkgcli.LoadManifest(item.ConfigFilePath)
		if loadErr != nil || manifest.Executable == "" {
			continue
		}
		toolSet, buildErr := pkgcli.BuildToolSet(mgr, item, manifest, enabledUserTools)
		if buildErr != nil {
			continue
		}
		if toolSet.Name() != "" {
			toolSets = append(toolSets, toolSet)
		}
	}
	return toolSets
}

// inferOpenAIVariant maps the persisted provider type to the matching OpenAI variant.
func inferOpenAIVariant(providerType pkgProvider.Type) openai.Variant {
	switch providerType {
	case pkgProvider.Aliyun:
		return openai.VariantQwen
	case pkgProvider.Deepseek:
		return openai.VariantDeepSeek
	default:
		return openai.VariantOpenAI
	}
}

// GetOrCreateRunner returns a cached runner for the provider/model pair, or creates one.
func (m *Manager) GetOrCreateRunner(
	baseURL string,
	apiKey string,
	modelName string,
	providerType pkgProvider.Type,
	enabledUserTools []string,
	sessionKind string,
	sessionID uint,
) (runner.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := runnerCacheKey(baseURL, modelName, providerType, enabledUserTools, sessionKind)
	if existing, ok := m.runners[key]; ok {
		return existing, nil
	}
	legacyKey := fmt.Sprintf("%s/%s", baseURL, modelName)
	if len(enabledUserTools) == 0 {
		if existing, ok := m.runners[legacyKey]; ok {
			return existing, nil
		}
	}

	instruction, err := prompt.Load(prompt_id.MainAgent)
	if err != nil {
		instruction = defaultMainAgentPrompt
	}

	variant := inferOpenAIVariant(providerType)

	mdl := openai.New(
		modelName,
		openai.WithBaseURL(baseURL),
		openai.WithAPIKey(apiKey),
		openai.WithVariant(variant),
	)

	ag := llmagent.New(
		"main",
		llmagent.WithModel(mdl),
		llmagent.WithInstruction(instruction),
		llmagent.WithGenerationConfig(model.GenerationConfig{
			Stream: true,
		}),
		llmagent.WithTools(m.buildAgentTools(enabledUserTools, sessionKind, sessionID)),
		llmagent.WithToolSets(append(m.buildMCPToolSets(enabledUserTools), m.buildCliToolSets(enabledUserTools)...)),
		llmagent.WithRefreshToolSetsOnRun(true),
		llmagent.WithToolCallbacks(m.newToolCallbacks()),
	)

	sessionSvc := inmemory.NewSessionService()
	created := runner.NewRunner(
		"lemontea",
		ag,
		runner.WithSessionService(sessionSvc),
	)

	m.runners[key] = created
	if len(enabledUserTools) == 0 {
		m.runners[legacyKey] = created
	}
	return created, nil
}

func runnerCacheKey(baseURL string, modelName string, providerType pkgProvider.Type, enabledUserTools []string, sessionKind string) string {
	names := append([]string(nil), enabledUserTools...)
	slices.Sort(names)
	return fmt.Sprintf("%s|%s|%s|%s|%s", baseURL, modelName, providerType, strings.Join(names, ","), sessionKind)
}

func loadAgentConfig() (*data_models.Config, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(filepath.Join(dataDir, dir.ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return &data_models.Config{
				Extensions: []data_models.ExtensionItem{},
				Memory:     data_models.MemoryConfig{Enabled: true},
			}, nil
		}
		return nil, err
	}

	config := &data_models.Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	if !strings.Contains(string(bytes), `"memory"`) {
		config.Memory.Enabled = true
	}
	if config.Extensions == nil {
		config.Extensions = []data_models.ExtensionItem{}
	}
	return config, nil
}

func mustDataDir() string {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return dir.HomeDir()
	}
	return dataDir
}

// NewChatHandler creates a chat handler bound to the manager.
func (m *Manager) NewChatHandler() *ChatHandler {
	return NewChatHandler(m)
}
