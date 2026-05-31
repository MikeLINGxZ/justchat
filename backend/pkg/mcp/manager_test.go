package mcp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testMCPBundle struct {
	Name         string
	Version      string
	Description  string
	Author       string
	ConfigFile   string
	Command      string
	Args         []string
	Env          map[string]string
	WithManifest bool
	ServerName   string
}

func writeTestMCPBundle(t *testing.T, bundle testMCPBundle) string {
	t.Helper()

	rootDir := t.TempDir()
	configPath := filepath.Join(rootDir, bundle.ConfigFile)

	if bundle.WithManifest {
		manifestPath := filepath.Join(rootDir, "manifest.json")
		manifest := `{
  "name": "` + bundle.Name + `",
  "version": "` + bundle.Version + `",
  "description": "` + bundle.Description + `",
  "author": "` + bundle.Author + `",
  "config_file": "` + bundle.ConfigFile + `"
}`
		if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}
	}

	config := `{
  "transport": "stdio",
  "command": "` + bundle.Command + `",
  "args": ["` + strings.Join(bundle.Args, `","`) + `"],
  "description": "` + bundle.Description + `"
}`
	if len(bundle.Args) == 0 {
		config = `{
  "transport": "stdio",
  "command": "` + bundle.Command + `",
  "args": [],
  "description": "` + bundle.Description + `"
}`
	}
	if bundle.ServerName != "" {
		envPairs := []string{}
		for key, value := range bundle.Env {
			envPairs = append(envPairs, `"`+key+`":"`+value+`"`)
		}
		config = `{
  "mcpServers": {
    "` + bundle.ServerName + `": {
      "command": "` + bundle.Command + `",
      "args": ["` + strings.Join(bundle.Args, `","`) + `"],
      "env": {` + strings.Join(envPairs, ",") + `}
    }
  }
}`
		if len(bundle.Args) == 0 {
			config = `{
  "mcpServers": {
    "` + bundle.ServerName + `": {
      "command": "` + bundle.Command + `",
      "args": [],
      "env": {` + strings.Join(envPairs, ",") + `}
    }
  }
}`
		}
	}
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return rootDir
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	return NewManager(filepath.Join(t.TempDir(), "data"))
}

func TestImportMCPDirectoryCopiesIntoVersionedDataDir(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Name:         "filesystem",
		Version:      "1.2.3",
		Command:      "node",
		Args:         []string{"server.js"},
		ConfigFile:   "mcp.json",
		Description:  "Filesystem tools",
		Author:       "Lemontea",
		WithManifest: true,
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{
		SourceDir: sourceDir,
		Enabled:   true,
	})
	if err != nil {
		t.Fatalf("import mcp: %v", err)
	}
	if !strings.Contains(item.RootDir, filepath.Join("plugins", "mcp", "filesystem", "1.2.3")) {
		t.Fatalf("expected versioned mcp path, got %q", item.RootDir)
	}
	if _, err := os.Stat(filepath.Join(item.RootDir, "manifest.json")); err != nil {
		t.Fatalf("expected copied manifest: %v", err)
	}
}

func TestImportMCPDirectoryFallsBackToDefaultVersionDir(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Name:         "clock",
		Command:      "node",
		Args:         []string{"server.js"},
		ConfigFile:   "mcp.json",
		WithManifest: true,
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{SourceDir: sourceDir})
	if err != nil {
		t.Fatalf("import mcp: %v", err)
	}
	if !strings.Contains(item.RootDir, filepath.Join("plugins", "mcp", "clock", "default")) {
		t.Fatalf("expected default version dir, got %q", item.RootDir)
	}
}

func TestImportMCPDirectoryWithoutManifestUsesFolderDefaults(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Command:    "uv",
		Args:       []string{"run", "mysql_mcp_server"},
		ConfigFile: "mcp.json",
		ServerName: "mysql",
		Env: map[string]string{
			"MYSQL_HOST": "127.0.0.1",
		},
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{SourceDir: sourceDir})
	if err != nil {
		t.Fatalf("import mcp without manifest: %v", err)
	}
	if item.Name != filepath.Base(sourceDir) {
		t.Fatalf("expected folder name %q, got %q", filepath.Base(sourceDir), item.Name)
	}
	if _, err := os.Stat(filepath.Join(item.RootDir, "mcp.json")); err != nil {
		t.Fatalf("expected copied config file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(item.RootDir, "manifest.json")); err != nil {
		t.Fatalf("expected generated manifest file: %v", err)
	}
}

func TestLoadServerConfigSupportsStandardMCPServersShape(t *testing.T) {
	rootDir := writeTestMCPBundle(t, testMCPBundle{
		Command:    "uv",
		Args:       []string{"run", "mysql_mcp_server"},
		ConfigFile: "mcp.json",
		ServerName: "mysql",
		Env: map[string]string{
			"MYSQL_HOST": "127.0.0.1",
			"MYSQL_PORT": "3306",
		},
	})

	cfg, err := loadServerConfig(filepath.Join(rootDir, "mcp.json"))
	if err != nil {
		t.Fatalf("load server config: %v", err)
	}
	if cfg.Command != "uv" {
		t.Fatalf("expected command uv, got %q", cfg.Command)
	}
	if cfg.Env["MYSQL_HOST"] != "127.0.0.1" {
		t.Fatalf("expected MYSQL_HOST env to be parsed")
	}
}

func TestBuildToolSetAcceptsWholeMCPSelection(t *testing.T) {
	manager := newTestManager(t)
	sourceDir := writeTestMCPBundle(t, testMCPBundle{
		Name:         "filesystem",
		Version:      "1.0.0",
		Command:      "node",
		Args:         []string{"server.js"},
		ConfigFile:   "mcp.json",
		Description:  "Filesystem tools",
		WithManifest: true,
	})

	item, err := manager.ImportMCP(context.Background(), ImportInput{
		SourceDir: sourceDir,
		Enabled:   true,
	})
	if err != nil {
		t.Fatalf("import mcp: %v", err)
	}

	toolSet, err := manager.BuildToolSet(item, []string{item.ID})
	if err != nil {
		t.Fatalf("build tool set: %v", err)
	}

	if _, ok := toolSet.(emptyToolSet); ok {
		t.Fatalf("expected whole-mcp selection to keep the tool set enabled")
	}
}

func TestMCPToolSetNameSanitizesInvalidModelToolNameCharacters(t *testing.T) {
	name := mcpToolSetName("mcp:filesystem:1.0.0")

	if strings.Contains(name, ":") {
		t.Fatalf("expected sanitized toolset name without colon, got %q", name)
	}
	if strings.Contains(name, ".") {
		t.Fatalf("expected sanitized toolset name without dot, got %q", name)
	}
	if name != "mcp_filesystem_1_0_0" {
		t.Fatalf("expected sanitized toolset name mcp_filesystem_1_0_0, got %q", name)
	}
}
