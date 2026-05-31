package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildEnvIsolatedSetsHomeAndXDG verifies isolated mode injects HOME / XDG_* pointing into the plugin data dir.
func TestBuildEnvIsolatedSetsHomeAndXDG(t *testing.T) {
	dataDir := t.TempDir()
	cliDataRoot := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli")

	env := BuildEnv(EnvParams{
		Isolation:  IsolationIsolated,
		CLIDataDir: cliDataRoot,
		PluginBin:  filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin"),
		ParentEnv:  []string{"PATH=/usr/bin", "HOME=/Users/me", "NPM_CONFIG_PREFIX=/old", "FOO=BAR"},
	})

	got := envMap(env)
	if got["HOME"] != filepath.Join(cliDataRoot, "home") {
		t.Fatalf("HOME=%q", got["HOME"])
	}
	if got["XDG_CONFIG_HOME"] != filepath.Join(cliDataRoot, "config") {
		t.Fatalf("XDG_CONFIG_HOME=%q", got["XDG_CONFIG_HOME"])
	}
	if got["XDG_DATA_HOME"] != filepath.Join(cliDataRoot, "data") {
		t.Fatalf("XDG_DATA_HOME=%q", got["XDG_DATA_HOME"])
	}
	if got["XDG_CACHE_HOME"] != filepath.Join(cliDataRoot, "cache") {
		t.Fatalf("XDG_CACHE_HOME=%q", got["XDG_CACHE_HOME"])
	}
	if got["FOO"] != "BAR" {
		t.Fatalf("FOO inheritance lost: %q", got["FOO"])
	}
	if _, ok := got["NPM_CONFIG_PREFIX"]; ok {
		t.Fatalf("NPM_CONFIG_PREFIX must be stripped in isolated mode")
	}
	if !strings.HasPrefix(got["PATH"], filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin")) {
		t.Fatalf("PATH prefix not applied: %q", got["PATH"])
	}
}

// TestBuildEnvSharedKeepsHome verifies shared mode preserves the user's real HOME and XDG_* values.
func TestBuildEnvSharedKeepsHome(t *testing.T) {
	env := BuildEnv(EnvParams{
		Isolation:  IsolationShared,
		CLIDataDir: "/ignored",
		PluginBin:  "/plugin/bin",
		ParentEnv:  []string{"HOME=/Users/me", "PATH=/usr/bin", "XDG_CONFIG_HOME=/userxdg"},
	})

	got := envMap(env)
	if got["HOME"] != "/Users/me" {
		t.Fatalf("HOME should be preserved, got %q", got["HOME"])
	}
	if got["XDG_CONFIG_HOME"] != "/userxdg" {
		t.Fatalf("XDG_CONFIG_HOME should be preserved, got %q", got["XDG_CONFIG_HOME"])
	}
	if !strings.HasPrefix(got["PATH"], "/plugin/bin") {
		t.Fatalf("PATH prefix not applied: %q", got["PATH"])
	}
}

// TestBuildEnvPrependsRuntimeBinBeforePluginBin verifies npm/npx resolve from the bundled runtime.
func TestBuildEnvPrependsRuntimeBinBeforePluginBin(t *testing.T) {
	env := BuildEnv(EnvParams{
		Isolation:  IsolationShared,
		CLIDataDir: "/ignored",
		RuntimeBin: "/runtime/bin",
		PluginBin:  "/plugin/bin",
		ParentEnv:  []string{"PATH=/usr/bin"},
	})

	got := envMap(env)
	wantPrefix := "/runtime/bin" + string(filepath.ListSeparator) + "/plugin/bin"
	if !strings.HasPrefix(got["PATH"], wantPrefix) {
		t.Fatalf("PATH should prefer runtime then plugin bin, got %q", got["PATH"])
	}
}

// envMap parses an exec.Cmd-style env slice ("KEY=VAL") into a map for assertion.
func envMap(env []string) map[string]string {
	out := make(map[string]string, len(env))
	for _, e := range env {
		idx := strings.IndexByte(e, '=')
		if idx <= 0 {
			continue
		}
		out[e[:idx]] = e[idx+1:]
	}
	return out
}
