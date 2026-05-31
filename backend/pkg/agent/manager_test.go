package agent

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/prompt_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompt"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeCliInstallProgressReporter struct{}

func (f *fakeCliInstallProgressReporter) ReportCliInstallProgress(_ context.Context, _ uint, _ tools.CliInstallProgressItem) error {
	return nil
}

func newTestStorage(t *testing.T) *storage.Storage {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stor, err := storage.NewStorageFromDB(db)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return stor
}

func writeAgentConfig(t *testing.T, cfg data_models.Config) {
	t.Helper()

	dataDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", dataDir)

	bytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, dir.ConfigFileName), bytes, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func writeMCPConfigFile(t *testing.T) string {
	t.Helper()

	rootDir := t.TempDir()
	configPath := filepath.Join(rootDir, "mcp.json")
	if err := os.WriteFile(configPath, []byte(`{
  "transport": "stdio",
  "command": "node",
  "args": ["server.js"],
  "description": "Browser tools"
}`), 0o644); err != nil {
		t.Fatalf("write mcp config: %v", err)
	}
	return configPath
}

func TestNewManagerRegistersPromptAndBuiltInTools(t *testing.T) {
	prompt.Register(prompt_id.MainAgent, "")
	prompt.Register(prompt_id.GenChatTitle, "")

	manager := NewManager(newTestStorage(t))

	loadedPrompt, err := prompt.Load(prompt_id.MainAgent)
	if err != nil {
		t.Fatalf("load prompt: %v", err)
	}
	if loadedPrompt == "" {
		t.Fatal("expected main agent prompt to be registered")
	}

	titlePrompt, err := prompt.Load(prompt_id.GenChatTitle)
	if err != nil {
		t.Fatalf("load title prompt: %v", err)
	}
	if titlePrompt == "" {
		t.Fatal("expected title prompt to be registered")
	}

	if manager.ToolRegistry() == nil {
		t.Fatal("expected tool registry to be initialized")
	}
	if manager.Streams() == nil {
		t.Fatal("expected stream manager to be initialized")
	}

	builtinTools := manager.ToolRegistry().BuiltinTools()
	if len(builtinTools) != 19 {
		t.Fatalf("expected 19 builtin tools, got %d", len(builtinTools))
	}

	userTools := manager.ToolRegistry().UserTools()
	if len(userTools) != 1 {
		t.Fatalf("expected 1 user tool, got %d", len(userTools))
	}
}

func TestBuildAgentToolsIncludesBuiltinAndEnabledUserTools(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	agentTools := manager.buildAgentTools([]string{"code_exec"}, "user", 1)
	if len(agentTools) != 11 {
		t.Fatalf("expected 11 tools (10 active builtin + 1 user), got %d", len(agentTools))
	}
}

func TestBuildAgentToolsIncludesWebToolsByDefault(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	agentTools := manager.buildAgentTools(nil, "user", 1)
	names := make(map[string]bool, len(agentTools))
	for _, current := range agentTools {
		names[current.Declaration().Name] = true
	}

	for _, want := range []string{"web_fetch", tools.WebSearchToolName} {
		if !names[want] {
			t.Fatalf("expected builtin web tool %q to be available by default, got %v", want, names)
		}
	}
}

func TestBuildAgentToolsIncludesTaskStateToolsForTaskSessions(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	manager.SetCliInstallProgressReporter(&fakeCliInstallProgressReporter{})

	agentTools := manager.buildAgentTools(nil, "task", 1)
	names := make(map[string]bool, len(agentTools))
	for _, current := range agentTools {
		names[current.Declaration().Name] = true
	}
	if !names[tools.SaveTaskStateToolName] || !names[tools.LoadTaskStateToolName] {
		t.Fatalf("expected task state tools in task session, got %v", names)
	}
	if !names[tools.ReportCliInstallProgressToolName] {
		t.Fatalf("expected cli install progress tool in task session, got %v", names)
	}
}

func TestBuildAgentToolsIncludesOrdinaryUserTools(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	agentTools := manager.buildAgentTools([]string{"web_search"}, "user", 1)
	names := make(map[string]bool, len(agentTools))
	for _, current := range agentTools {
		names[current.Declaration().Name] = true
	}

	for _, want := range []string{tools.QuestionToolName, tools.TodoWriteToolName, tools.WebSearchToolName} {
		if !names[want] {
			t.Fatalf("expected ordinary user tool %q to be available, got %v", want, names)
		}
	}
}

func TestBuildMCPToolSetsReturnsNoneWhenNoExtensionToolIsSelected(t *testing.T) {
	configPath := writeMCPConfigFile(t)
	writeAgentConfig(t, data_models.Config{
		Extensions: []data_models.ExtensionItem{
			{
				ID:             "mcp:browser:1.0.0",
				Name:           "browser",
				Description:    "Browser tools",
				Kind:           "mcp",
				Enabled:        true,
				ConfigFilePath: configPath,
			},
		},
	})

	manager := NewManager(newTestStorage(t))

	toolSets := manager.buildMCPToolSets([]string{})
	if len(toolSets) != 0 {
		t.Fatalf("expected 0 mcp tool sets when no extension tool is selected, got %d", len(toolSets))
	}
}

func TestGetOrCreateRunnerCachesRunnerByBaseURLAndModel(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	first, err := manager.GetOrCreateRunner(
		"https://api.example.com/v1",
		"key-1",
		"gpt-test",
		pkgProvider.OpenAiCompatibility,
		[]string{"code_exec"},
		"user",
		1,
	)
	if err != nil {
		t.Fatalf("first runner: %v", err)
	}

	second, err := manager.GetOrCreateRunner(
		"https://api.example.com/v1",
		"key-2",
		"gpt-test",
		pkgProvider.OpenAiCompatibility,
		[]string{"web_fetch"},
		"user",
		1,
	)
	if err != nil {
		t.Fatalf("second runner: %v", err)
	}

	if first == second {
		t.Fatal("expected different runners when enabled tools differ")
	}
}

func TestInferOpenAIVariantFromProviderType(t *testing.T) {
	tests := []struct {
		name         string
		providerType pkgProvider.Type
		expected     string
	}{
		{
			name:         "aliyun maps to qwen variant",
			providerType: pkgProvider.Aliyun,
			expected:     "qwen",
		},
		{
			name:         "deepseek maps to deepseek variant",
			providerType: pkgProvider.Deepseek,
			expected:     "deepseek",
		},
		{
			name:         "openai compatibility defaults to openai variant",
			providerType: pkgProvider.OpenAiCompatibility,
			expected:     "openai",
		},
		{
			name:         "unknown provider falls back to openai variant",
			providerType: pkgProvider.Type("custom"),
			expected:     "openai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inferOpenAIVariant(tt.providerType)
			if string(got) != tt.expected {
				t.Fatalf("inferOpenAIVariant(%q) = %q, want %q", tt.providerType, got, tt.expected)
			}
		})
	}
}
