// Package migration handles one-time data directory layout migrations between Lemontea releases.
package migration

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// markerFileName is the idempotency marker written under the plugins root after a successful migration.
const markerFileName = ".migrated_v1"

// MigrationResult summarizes what happened during a single migration run.
type MigrationResult struct {
	AlreadyMigrated bool
	MovedMCP        bool
	MovedPlugin     bool
	RewroteConfig   bool
	Warnings        []string
}

// MigrateExtensionDirs performs a one-time migration of legacy mcp/ and plugin/ dirs into plugins/.
// It is safe to call on every startup: an idempotency marker prevents repeat work.
func MigrateExtensionDirs(dataDir string) (MigrationResult, error) {
	result := MigrationResult{}

	if err := os.MkdirAll(dir.ExtensionsRoot(dataDir), 0o755); err != nil {
		return result, err
	}

	markerPath := filepath.Join(dir.ExtensionsRoot(dataDir), markerFileName)
	if _, err := os.Stat(markerPath); err == nil {
		result.AlreadyMigrated = true
		return result, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	// Step 1: move legacy dirs if present.
	if moved, warn, err := moveLegacy(dir.LegacyMCPRoot(dataDir), dir.MCPRoot(dataDir)); err != nil {
		return result, err
	} else {
		result.MovedMCP = moved
		if warn != "" {
			result.Warnings = append(result.Warnings, warn)
		}
	}
	if moved, warn, err := moveLegacy(dir.LegacyPluginRoot(dataDir), dir.PluginRoot(dataDir)); err != nil {
		return result, err
	} else {
		result.MovedPlugin = moved
		if warn != "" {
			result.Warnings = append(result.Warnings, warn)
		}
	}

	// Step 2: always attempt to rewrite config paths. The helper is idempotent and only
	// touches entries that match a legacy prefix, so this self-heals previous runs that
	// moved files but failed to update config.json before crashing.
	rewrote, err := rewriteConfigPaths(dataDir)
	if err != nil {
		return result, err
	}
	result.RewroteConfig = rewrote

	// Step 3: write idempotency marker.
	if err := os.WriteFile(markerPath, []byte("v1"), 0o644); err != nil {
		return result, err
	}
	return result, nil
}

// moveLegacy moves srcDir into dstDir when srcDir exists and dstDir does not.
// It returns (moved, warning, error). When both exist we keep srcDir and log via warning.
func moveLegacy(srcDir string, dstDir string) (bool, string, error) {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, "", nil
		}
		return false, "", err
	}
	if !srcInfo.IsDir() {
		return false, "legacy path is not a directory: " + srcDir, nil
	}

	if _, err := os.Stat(dstDir); err == nil {
		return false, "both legacy and new dirs exist, leaving as-is: " + srcDir, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, "", err
	}

	if err := os.MkdirAll(filepath.Dir(dstDir), 0o755); err != nil {
		return false, "", err
	}
	if err := os.Rename(srcDir, dstDir); err != nil {
		return false, "", err
	}
	return true, "", nil
}

// rewriteConfigPaths loads config.json into a typed data_models.Config, replaces extension
// dir prefixes on each ExtensionItem, and writes the file back atomically via tmp+rename.
// It returns true when at least one ExtensionItem path was rewritten.
//
// Note: the round-trip through data_models.Config silently drops any unknown top-level keys
// from config.json. This matches the canonical schema used everywhere else in the app and
// is acceptable today; if the on-disk schema ever diverges from the Go struct, this helper
// must switch to a surgical token-level rewrite.
func rewriteConfigPaths(dataDir string) (bool, error) {
	configPath := filepath.Join(dataDir, dir.ConfigFileName)
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	var cfg data_models.Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return false, err
	}

	rewrote := false
	for i := range cfg.Extensions {
		if newValue, changed := rewriteLegacyPathPrefix(cfg.Extensions[i].RootDir, dataDir); changed {
			cfg.Extensions[i].RootDir = newValue
			rewrote = true
		}
		if newValue, changed := rewriteLegacyPathPrefix(cfg.Extensions[i].ConfigFilePath, dataDir); changed {
			cfg.Extensions[i].ConfigFilePath = newValue
			rewrote = true
		}
	}
	if !rewrote {
		return false, nil
	}

	updated, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return false, err
	}

	tmpPath := configPath + ".tmp"
	if err := os.WriteFile(tmpPath, updated, 0o644); err != nil {
		return false, err
	}
	if err := os.Rename(tmpPath, configPath); err != nil {
		// Best-effort cleanup; ignore error since the real failure is the rename.
		_ = os.Remove(tmpPath)
		return false, err
	}
	return true, nil
}

// rewriteLegacyPathPrefix returns a rewritten path when value starts with the legacy
// "{dataDir}/mcp/" or "{dataDir}/plugin/" prefix, plus a boolean indicating whether the
// value changed. The trailing path separator on the prefix guards against substring
// matches like "{dataDir}/mcp_old/...".
func rewriteLegacyPathPrefix(value string, dataDir string) (string, bool) {
	if value == "" {
		return value, false
	}
	sep := string(filepath.Separator)
	legacyMCP := dir.LegacyMCPRoot(dataDir) + sep
	legacyPlugin := dir.LegacyPluginRoot(dataDir) + sep

	switch {
	case strings.HasPrefix(value, legacyMCP):
		return dir.MCPRoot(dataDir) + sep + value[len(legacyMCP):], true
	case strings.HasPrefix(value, legacyPlugin):
		return dir.PluginRoot(dataDir) + sep + value[len(legacyPlugin):], true
	}
	return value, false
}
