package plugin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

func writeTestMCPBundleDir(t *testing.T) string {
	t.Helper()

	rootDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(rootDir, "manifest.json"), []byte(`{
  "name": "filesystem",
  "version": "1.0.0",
  "description": "Filesystem tools",
  "author": "Lemontea",
  "config_file": "mcp.json"
}`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rootDir, "mcp.json"), []byte(`{
  "transport": "stdio",
  "command": "node",
  "args": ["server.js"],
  "description": "Filesystem tools"
}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return rootDir
}

func newTestPluginService(t *testing.T) *Plugin {
	t.Helper()

	tempDataDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDataDir)
	return NewPlugin()
}

func TestListExtensionsReturnsImportedItems(t *testing.T) {
	service := newTestPluginService(t)
	ctx := context.Background()

	_, err := service.ImportExtension(ctx, plugin_dto.ImportExtensionInput{
		Kind: "mcp",
		Path: writeTestMCPBundleDir(t),
	})
	if err != nil {
		t.Fatalf("import extension: %v", err)
	}

	out, err := service.ListExtensions(ctx, plugin_dto.ListExtensionsInput{})
	if err != nil {
		t.Fatalf("list extensions: %v", err)
	}
	if len(out.Extensions) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(out.Extensions))
	}
	if out.Extensions[0].Name != "filesystem" {
		t.Fatalf("expected filesystem extension, got %q", out.Extensions[0].Name)
	}
}

func TestListExtensionsPopulatesRuntimeErrorWhenDiscoveryFails(t *testing.T) {
	service := newTestPluginService(t)
	ctx := context.Background()

	rootDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(rootDir, "mcp.json"), []byte(`{
  "mcpServers": {
    "broken": {
      "command": "definitely-missing-command",
      "args": []
    }
  }
}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := service.ImportExtension(ctx, plugin_dto.ImportExtensionInput{
		Kind: "mcp",
		Path: rootDir,
	})
	if err != nil {
		t.Fatalf("import extension: %v", err)
	}

	out, err := service.ListExtensions(ctx, plugin_dto.ListExtensionsInput{})
	if err != nil {
		t.Fatalf("list extensions: %v", err)
	}
	if len(out.Extensions) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(out.Extensions))
	}
	if out.Extensions[0].RuntimeStatus != "error" {
		t.Fatalf("expected runtime error status, got %q", out.Extensions[0].RuntimeStatus)
	}
	if out.Extensions[0].RuntimeMessage == "" {
		t.Fatalf("expected runtime error message to be populated")
	}
}

func TestListAvailableToolsAggregatesEnabledMCPByServer(t *testing.T) {
	service := newTestPluginService(t)
	ctx := context.Background()

	cfg, err := service.loadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Extensions = []data_models.ExtensionItem{
		{
			ID:          "mcp:filesystem:1.0.0",
			Name:        "filesystem",
			Description: "Filesystem tools",
			Kind:        "mcp",
			Enabled:     true,
			Tools: []data_models.ExtensionTool{
				{ToolID: "mcp:filesystem:1.0.0_read_file", Name: "read_file", Description: "Read a file"},
				{ToolID: "mcp:filesystem:1.0.0_write_file", Name: "write_file", Description: "Write a file"},
			},
		},
	}
	if err := service.saveConfig(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	out, err := service.ListAvailableTools(ctx, plugin_dto.ListAvailableToolsInput{})
	if err != nil {
		t.Fatalf("list available tools: %v", err)
	}

	var mcpTools []plugin_dto.ToolItem
	for _, item := range out.Tools {
		if item.Category == "mcp" {
			mcpTools = append(mcpTools, item)
		}
	}

	if len(mcpTools) != 1 {
		t.Fatalf("expected 1 aggregated mcp entry, got %d", len(mcpTools))
	}
	if mcpTools[0].ID != "mcp:filesystem:1.0.0" {
		t.Fatalf("expected aggregated mcp id, got %q", mcpTools[0].ID)
	}
	if mcpTools[0].Name != "filesystem" {
		t.Fatalf("expected mcp name filesystem, got %q", mcpTools[0].Name)
	}
}

func TestListAvailableToolsOnlyShowsPracticalBuiltinsForOrdinaryUsers(t *testing.T) {
	service := newTestPluginService(t)

	out, err := service.ListAvailableTools(context.Background(), plugin_dto.ListAvailableToolsInput{})
	if err != nil {
		t.Fatalf("list available tools: %v", err)
	}

	got := map[string]bool{}
	for _, item := range out.Tools {
		got[item.ID] = true
	}

	for _, want := range []string{"shell", "question", "web_fetch", "web_search", "Skill", "todo_write"} {
		if !got[want] {
			t.Fatalf("expected practical builtin %q to be visible, got %v", want, got)
		}
	}
	for _, hidden := range []string{"file_read", "file_write", "code_exec"} {
		if got[hidden] {
			t.Fatalf("expected developer builtin %q to stay hidden, got %v", hidden, got)
		}
	}
}
