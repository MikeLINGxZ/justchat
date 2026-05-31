package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// installFixture returns a Manager whose npm is mocked to drop a synthetic package.json + bin under TargetDir.
// The mock simulates a successful `npm install lark-cli --prefix <target>`.
func installFixture(t *testing.T) (*Manager, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("manager tests use POSIX shell mocks")
	}
	dataDir := t.TempDir()
	binDir := t.TempDir()
	nodePath := filepath.Join(binDir, "node")
	npmPath := filepath.Join(binDir, "npm")
	nodeBody := `#!/usr/bin/env bash
# Args: <npm script> install <source> --prefix <target> --no-audit --no-fund
shift
target=""
seen_prefix=0
for a in "$@"; do
  if [ "$seen_prefix" = "1" ]; then target="$a"; seen_prefix=0; fi
  if [ "$a" = "--prefix" ]; then seen_prefix=1; fi
done
mkdir -p "$target/node_modules/.bin"
mkdir -p "$target/node_modules/lark-cli"
cat > "$target/node_modules/lark-cli/package.json" <<EOF
{"name":"lark-cli","version":"1.2.3","description":"feishu","bin":"lark.js"}
EOF
cat > "$target/node_modules/.bin/lark-cli" <<EOF
#!/usr/bin/env bash
if [ "\$1" = "--help" ]; then echo "usage: lark-cli ..."; exit 0; fi
exit 1
EOF
chmod +x "$target/node_modules/.bin/lark-cli"
exit 0
`
	if err := os.WriteFile(nodePath, []byte(nodeBody), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(npmPath, []byte("// fake npm entry\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return NewManager(dataDir, nodePath, npmPath), dataDir
}

// TestManagerInstallFromNpmWritesFilesAndManifest verifies a happy-path install produces the expected layout.
func TestManagerInstallFromNpmWritesFilesAndManifest(t *testing.T) {
	mgr, dataDir := installFixture(t)
	res, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if res.Name != "lark-cli" || res.Version != "1.2.3" {
		t.Fatalf("res: %+v", res)
	}
	installRoot := filepath.Join(dataDir, "plugins", "cli", "lark-cli")
	if _, err := os.Stat(filepath.Join(installRoot, "node_modules", ".bin", "lark-cli")); err != nil {
		t.Fatalf("bin not written: %v", err)
	}
	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("manifest load: %v", err)
	}
	if manifest.Name != "lark-cli" || manifest.Version != "1.2.3" {
		t.Fatalf("manifest meta: %+v", manifest)
	}
	if manifest.Executable == "" {
		t.Fatalf("manifest executable empty")
	}
	if len(manifest.Tools) != 0 {
		t.Fatalf("expected empty tools at install time, got %d", len(manifest.Tools))
	}
}

// TestManagerInstallFromLocalUsesProvidedPath verifies InstallFromLocal forwards the source path to npm.
func TestManagerInstallFromLocalUsesProvidedPath(t *testing.T) {
	mgr, dataDir := installFixture(t)
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "package.json"), []byte(`{"name":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := mgr.InstallFromLocal(context.Background(), InstallParams{LocalPath: src, Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install local: %v", err)
	}
	if res.Name != "lark-cli" {
		t.Fatalf("res: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli", "lark-cli")); err != nil {
		t.Fatalf("install dir missing: %v", err)
	}
}

// TestManagerInstallFromNpmFallsBackToNodeModulesBin verifies install still succeeds when the package itself has no bin entry
// but npm has materialized an executable into node_modules/.bin.
func TestManagerInstallFromNpmFallsBackToNodeModulesBin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("manager tests use POSIX shell mocks")
	}
	dataDir := t.TempDir()
	binDir := t.TempDir()
	nodePath := filepath.Join(binDir, "node")
	npmPath := filepath.Join(binDir, "npm")
	nodeBody := `#!/usr/bin/env bash
shift
target=""
seen_prefix=0
for a in "$@"; do
  if [ "$seen_prefix" = "1" ]; then target="$a"; seen_prefix=0; fi
  if [ "$a" = "--prefix" ]; then seen_prefix=1; fi
done
mkdir -p "$target/node_modules/.bin"
mkdir -p "$target/node_modules/@larksuite/cli"
cat > "$target/node_modules/@larksuite/cli/package.json" <<EOF
{"name":"@larksuite/cli","version":"1.0.0","description":"feishu installer"}
EOF
cat > "$target/node_modules/.bin/lark-cli" <<EOF
#!/usr/bin/env bash
exit 0
EOF
chmod +x "$target/node_modules/.bin/lark-cli"
exit 0
`
	if err := os.WriteFile(nodePath, []byte(nodeBody), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(npmPath, []byte("// fake npm entry\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	mgr := NewManager(dataDir, nodePath, npmPath)
	res, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "@larksuite/cli", Name: "feishu-cli"})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if filepath.Base(res.Executable) != "lark-cli" {
		t.Fatalf("expected fallback executable lark-cli, got %q", res.Executable)
	}
}

// TestManagerUninstallRemovesInstallButKeepsData verifies Uninstall removes plugins/cli/{name} but preserves cli_data/{name}.
func TestManagerUninstallRemovesInstallButKeepsData(t *testing.T) {
	mgr, dataDir := installFixture(t)
	if _, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Uninstall(context.Background(), "lark-cli"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli", "lark-cli")); !os.IsNotExist(err) {
		t.Fatalf("install dir not removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")); err != nil {
		t.Fatalf("cli_data should be preserved: %v", err)
	}
}

// TestManagerResetDataClearsCLIDataDir verifies ResetData removes the cli_data/{name} subtree.
func TestManagerResetDataClearsCLIDataDir(t *testing.T) {
	mgr, dataDir := installFixture(t)
	if _, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.ResetData(context.Background(), "lark-cli"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli_data", "lark-cli")); !os.IsNotExist(err) {
		t.Fatalf("cli_data should be removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli", "lark-cli")); err != nil {
		t.Fatalf("install dir should remain: %v", err)
	}
}

// TestManagerRunCommandUsesBundledRuntime verifies arbitrary CLI commands can run through the managed env.
func TestManagerRunCommandUsesBundledRuntime(t *testing.T) {
	mgr, _ := installFixture(t)
	if _, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.RunCommand(context.Background(), "lark-cli", []string{"--help"}, OutputText, 30)
	if err != nil {
		t.Fatalf("run command: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %+v", res)
	}
	if !strings.Contains(res.Stdout, "usage: lark-cli") {
		t.Fatalf("unexpected stdout: %q", res.Stdout)
	}
}

// installToolFixture installs the standard fake plugin and writes a manifest with two tools for testing.
func installToolFixture(t *testing.T) (*Manager, string) {
	t.Helper()
	mgr, dataDir := installFixture(t)
	if _, err := mgr.InstallFromNpm(context.Background(), InstallParams{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, _ := LoadManifest(manifestPath)
	manifest.Tools = []Tool{
		{
			Name:         "echo_tool",
			InputSchema:  json.RawMessage(`{"type":"object","properties":{"msg":{"type":"string"}}}`),
			ArgvTemplate: []string{"--help"}, // our fake-bin echoes "usage: ..." for --help
			OutputMode:   OutputText,
			Enabled:      true,
		},
		{
			Name:         "missing_field",
			InputSchema:  json.RawMessage(`{"type":"object","properties":{"a":{"type":"string"}},"required":["a"]}`),
			ArgvTemplate: []string{"--a", "{a}"},
			OutputMode:   OutputText,
			Enabled:      true,
		},
		{
			Name:         "optional_flag",
			InputSchema:  json.RawMessage(`{"type":"object","properties":{"cmd":{"type":"string"},"extra":{"type":"string"}},"required":["cmd"]}`),
			ArgvTemplate: []string{"{cmd}", "--extra", "{extra}"},
			OutputMode:   OutputText,
			Enabled:      true,
		},
	}
	if err := SaveManifest(manifestPath, manifest); err != nil {
		t.Fatal(err)
	}
	return mgr, dataDir
}

// TestRunToolHappyPath verifies RunTool resolves the tool, runs the executable, and returns parsed stdout.
func TestRunToolHappyPath(t *testing.T) {
	mgr, _ := installToolFixture(t)
	res, err := mgr.RunTool(context.Background(), "lark-cli", "echo_tool", map[string]any{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%q", res.ExitCode, res.Stderr)
	}
	if res.Stdout == "" {
		t.Fatalf("expected stdout, got empty")
	}
}

// TestRunToolUnknownToolErrors verifies RunTool returns an error when the tool is not in the manifest.
func TestRunToolUnknownToolErrors(t *testing.T) {
	mgr, _ := installToolFixture(t)
	if _, err := mgr.RunTool(context.Background(), "lark-cli", "nope", nil); err == nil {
		t.Fatal("expected unknown-tool error")
	}
}

// TestRunToolMissingInputFieldErrors verifies RunTool errors when a required placeholder field is missing.
func TestRunToolMissingInputFieldErrors(t *testing.T) {
	mgr, _ := installToolFixture(t)
	if _, err := mgr.RunTool(context.Background(), "lark-cli", "missing_field", map[string]any{}); err == nil {
		t.Fatal("expected missing-input-field error")
	}
}

// TestRunToolMissingOptionalPlaceholderDropsPair verifies that an optional placeholder, when omitted by the caller,
// causes both the placeholder segment and its preceding flag segment to be dropped — matching common CLI optional-flag semantics.
func TestRunToolMissingOptionalPlaceholderDropsPair(t *testing.T) {
	mgr, dataDir := installToolFixture(t)
	binPath := filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli")
	if err := os.WriteFile(binPath, []byte("#!/usr/bin/env bash\necho \"args: $@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.RunTool(context.Background(), "lark-cli", "optional_flag", map[string]any{"cmd": "list"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(res.Stdout, "args: list") {
		t.Fatalf("expected args trimmed to just the positional, got %q", res.Stdout)
	}
	if strings.Contains(res.Stdout, "--extra") {
		t.Fatalf("expected --extra to be dropped when extra is absent, got %q", res.Stdout)
	}
}

// TestRunToolKeepsOptionalPlaceholderWhenProvided verifies the optional flag pair survives when the caller supplies its value.
func TestRunToolKeepsOptionalPlaceholderWhenProvided(t *testing.T) {
	mgr, dataDir := installToolFixture(t)
	binPath := filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli")
	if err := os.WriteFile(binPath, []byte("#!/usr/bin/env bash\necho \"args: $@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.RunTool(context.Background(), "lark-cli", "optional_flag", map[string]any{"cmd": "send", "extra": "x"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(res.Stdout, "--extra x") {
		t.Fatalf("expected --extra x present, got %q", res.Stdout)
	}
}

// TestRunToolSubstitutesNumericInput verifies non-string input values are converted to their JSON string form.
func TestRunToolSubstitutesNumericInput(t *testing.T) {
	mgr, dataDir := installToolFixture(t)
	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	m, _ := LoadManifest(manifestPath)
	// Override fake bin to echo back argv so we can assert substitution worked.
	binPath := filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli")
	if err := os.WriteFile(binPath, []byte("#!/usr/bin/env bash\necho \"args: $@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	m.Tools = []Tool{
		{
			Name:         "numtool",
			InputSchema:  json.RawMessage(`{"type":"object","properties":{"n":{"type":"number"}}}`),
			ArgvTemplate: []string{"--n", "{n}"},
			OutputMode:   OutputText,
			Enabled:      true,
		},
	}
	if err := SaveManifest(manifestPath, m); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.RunTool(context.Background(), "lark-cli", "numtool", map[string]any{"n": 42})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(res.Stdout, "--n 42") {
		t.Fatalf("expected '--n 42' in stdout, got %q", res.Stdout)
	}
}

// TestRunToolRepairsMissingExecutable verifies execution can recover from an older manifest missing executable.
func TestRunToolRepairsMissingExecutable(t *testing.T) {
	mgr, dataDir := installToolFixture(t)
	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Executable = ""
	if err := SaveManifest(manifestPath, manifest); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.RunTool(context.Background(), "lark-cli", "echo_tool", map[string]any{})
	if err != nil {
		t.Fatalf("expected RunTool to repair executable, got error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("expected successful execution after repair, exit=%d stderr=%q", res.ExitCode, res.Stderr)
	}

	repaired, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if repaired.Executable == "" {
		t.Fatalf("expected repaired executable to persist: %+v", repaired)
	}
}
