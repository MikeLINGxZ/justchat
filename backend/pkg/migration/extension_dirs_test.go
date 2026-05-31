package migration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeConfigJSON helper writes a minimal config.json with the given Extensions slice for migration tests.
func writeConfigJSON(t *testing.T, dataDir string, extensions []map[string]any) {
	t.Helper()
	cfg := map[string]any{
		"data_dir":   dataDir,
		"extensions": extensions,
	}
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "config.json"), bytes, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

// readConfigJSON helper loads back config.json as a generic map.
func readConfigJSON(t *testing.T, dataDir string) map[string]any {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Join(dataDir, "config.json"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	out := map[string]any{}
	if err := json.Unmarshal(bytes, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

// TestMigrateExtensionDirs_NoLegacyNoOp verifies migration is a no-op when neither legacy dir exists.
func TestMigrateExtensionDirs_NoLegacyNoOp(t *testing.T) {
	dataDir := t.TempDir()
	writeConfigJSON(t, dataDir, nil)

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if result.MovedMCP || result.MovedPlugin {
		t.Fatalf("expected no moves, got %+v", result)
	}

	markerPath := filepath.Join(dataDir, "plugins", ".migrated_v1")
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("expected marker written even on no-op, got %v", err)
	}
}

// TestMigrateExtensionDirs_MovesLegacyMCP verifies a legacy mcp/ dir is moved to plugins/mcp/.
func TestMigrateExtensionDirs_MovesLegacyMCP(t *testing.T) {
	dataDir := t.TempDir()
	writeConfigJSON(t, dataDir, []map[string]any{
		{
			"id":               "mcp:demo:default",
			"kind":             "mcp",
			"root_dir":         filepath.Join(dataDir, "mcp", "demo", "default"),
			"config_file_path": filepath.Join(dataDir, "mcp", "demo", "default", "mcp.json"),
		},
	})

	legacyDir := filepath.Join(dataDir, "mcp", "demo", "default")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyDir, "marker.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if !result.MovedMCP || !result.RewroteConfig {
		t.Fatalf("expected MovedMCP and RewroteConfig, got %+v", result)
	}

	newPath := filepath.Join(dataDir, "plugins", "mcp", "demo", "default", "marker.txt")
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected file moved to %s, got %v", newPath, err)
	}

	cfg := readConfigJSON(t, dataDir)
	extensions := cfg["extensions"].([]any)
	first := extensions[0].(map[string]any)
	expectedRoot := filepath.Join(dataDir, "plugins", "mcp", "demo", "default")
	if first["root_dir"] != expectedRoot {
		t.Fatalf("root_dir not rewritten, got %v", first["root_dir"])
	}
	expectedCfg := filepath.Join(expectedRoot, "mcp.json")
	if first["config_file_path"] != expectedCfg {
		t.Fatalf("config_file_path not rewritten, got %v", first["config_file_path"])
	}
}

// TestMigrateExtensionDirs_IdempotentSkipsOnMarker verifies that a second run is a no-op.
func TestMigrateExtensionDirs_IdempotentSkipsOnMarker(t *testing.T) {
	dataDir := t.TempDir()
	writeConfigJSON(t, dataDir, nil)

	if _, err := MigrateExtensionDirs(dataDir); err != nil {
		t.Fatalf("first migrate: %v", err)
	}

	legacyDir := filepath.Join(dataDir, "mcp", "demo", "default")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	if !result.AlreadyMigrated {
		t.Fatalf("expected AlreadyMigrated, got %+v", result)
	}
	if _, err := os.Stat(legacyDir); err != nil {
		t.Fatalf("legacy dir should be untouched on second run, got %v", err)
	}
}

// TestMigrateExtensionDirs_BothExistKeepsLegacy verifies both legacy and new co-existing produces a warning.
func TestMigrateExtensionDirs_BothExistKeepsLegacy(t *testing.T) {
	dataDir := t.TempDir()
	writeConfigJSON(t, dataDir, nil)

	if err := os.MkdirAll(filepath.Join(dataDir, "mcp", "x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dataDir, "plugins", "mcp", "y"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if result.MovedMCP {
		t.Fatalf("should not move when both exist")
	}
	if len(result.Warnings) == 0 {
		t.Fatalf("expected at least one warning")
	}
	if _, err := os.Stat(filepath.Join(dataDir, "mcp", "x")); err != nil {
		t.Fatalf("legacy x should remain, got %v", err)
	}
}

// TestMigrateExtensionDirs_MovesLegacyPlugin verifies a legacy plugin/ dir is moved to plugins/plugin/.
func TestMigrateExtensionDirs_MovesLegacyPlugin(t *testing.T) {
	dataDir := t.TempDir()
	writeConfigJSON(t, dataDir, []map[string]any{
		{
			"id":               "plugin:demo:default",
			"kind":             "plugin",
			"root_dir":         filepath.Join(dataDir, "plugin", "demo", "default"),
			"config_file_path": filepath.Join(dataDir, "plugin", "demo", "default", "plugin.json"),
		},
	})

	legacyDir := filepath.Join(dataDir, "plugin", "demo", "default")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyDir, "marker.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if !result.MovedPlugin || !result.RewroteConfig {
		t.Fatalf("expected MovedPlugin and RewroteConfig, got %+v", result)
	}

	newPath := filepath.Join(dataDir, "plugins", "plugin", "demo", "default", "marker.txt")
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected file moved to %s, got %v", newPath, err)
	}

	cfg := readConfigJSON(t, dataDir)
	extensions := cfg["extensions"].([]any)
	first := extensions[0].(map[string]any)
	expectedRoot := filepath.Join(dataDir, "plugins", "plugin", "demo", "default")
	if first["root_dir"] != expectedRoot {
		t.Fatalf("root_dir not rewritten, got %v", first["root_dir"])
	}
	expectedCfg := filepath.Join(expectedRoot, "plugin.json")
	if first["config_file_path"] != expectedCfg {
		t.Fatalf("config_file_path not rewritten, got %v", first["config_file_path"])
	}
}

// TestMigrateExtensionDirs_DoesNotRewriteSimilarPrefix verifies entries that merely share
// a substring with a legacy root (e.g. "mcp_old/...") are left untouched while real
// legacy entries in the same config still get rewritten.
func TestMigrateExtensionDirs_DoesNotRewriteSimilarPrefix(t *testing.T) {
	dataDir := t.TempDir()
	similarPath := filepath.Join(dataDir, "mcp_old", "something")
	similarCfgPath := filepath.Join(similarPath, "mcp.json")
	realLegacyRoot := filepath.Join(dataDir, "mcp", "foo", "default")
	realLegacyCfg := filepath.Join(realLegacyRoot, "mcp.json")

	writeConfigJSON(t, dataDir, []map[string]any{
		{
			"id":               "mcp:similar:default",
			"kind":             "mcp",
			"root_dir":         similarPath,
			"config_file_path": similarCfgPath,
		},
		{
			"id":               "mcp:foo:default",
			"kind":             "mcp",
			"root_dir":         realLegacyRoot,
			"config_file_path": realLegacyCfg,
		},
	})

	if err := os.MkdirAll(realLegacyRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(realLegacyRoot, "marker.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateExtensionDirs(dataDir)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if !result.MovedMCP || !result.RewroteConfig {
		t.Fatalf("expected MovedMCP and RewroteConfig, got %+v", result)
	}

	cfg := readConfigJSON(t, dataDir)
	extensions := cfg["extensions"].([]any)

	similar := extensions[0].(map[string]any)
	if similar["root_dir"] != similarPath {
		t.Fatalf("mcp_old root_dir should be unchanged, got %v", similar["root_dir"])
	}
	if similar["config_file_path"] != similarCfgPath {
		t.Fatalf("mcp_old config_file_path should be unchanged, got %v", similar["config_file_path"])
	}

	real := extensions[1].(map[string]any)
	expectedRoot := filepath.Join(dataDir, "plugins", "mcp", "foo", "default")
	if real["root_dir"] != expectedRoot {
		t.Fatalf("real legacy root_dir should be rewritten, got %v", real["root_dir"])
	}
	expectedCfg := filepath.Join(expectedRoot, "mcp.json")
	if real["config_file_path"] != expectedCfg {
		t.Fatalf("real legacy config_file_path should be rewritten, got %v", real["config_file_path"])
	}
}
