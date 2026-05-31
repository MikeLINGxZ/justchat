package cli

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeFakeNpm writes a fake bundled node + npm pair that records every invocation into log files.
// The fake node accepts "node <npmPath> <args...>", creates the --prefix dir if requested, and exits 0.
func writeFakeNpm(t *testing.T, recordDir string) (nodePath, npmPath string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("npm install tests use POSIX shell mocks")
	}
	binDir := t.TempDir()
	npmPath = filepath.Join(binDir, "npm")
	nodePath = filepath.Join(binDir, "node")
	if err := os.WriteFile(npmPath, []byte("// fake npm entry\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	nodeBody := `#!/usr/bin/env bash
echo "NODE:$@" >> "` + filepath.Join(recordDir, "node.log") + `"
echo "SCRIPT:$1" >> "` + filepath.Join(recordDir, "npm.log") + `"
shift
echo "ARGS:$@" >> "` + filepath.Join(recordDir, "npm.log") + `"
i=0
for a in "$@"; do
  i=$((i+1))
  if [ "$a" = "--prefix" ]; then
    eval "next=\${$((i+1))}"
    mkdir -p "$next/node_modules/.bin"
  fi
done
exit 0
`
	if err := os.WriteFile(nodePath, []byte(nodeBody), 0o755); err != nil {
		t.Fatal(err)
	}
	return nodePath, npmPath
}

// TestNpmInstallPassesPrefix verifies the install call includes --prefix <targetDir>.
func TestNpmInstallPassesPrefix(t *testing.T) {
	logDir := t.TempDir()
	nodePath, npmPath := writeFakeNpm(t, logDir)
	targetDir := filepath.Join(t.TempDir(), "install")

	err := NpmInstall(context.Background(), NpmInstallParams{
		NodePath:  nodePath,
		NpmPath:   npmPath,
		Source:    "lark-cli",
		TargetDir: targetDir,
		ParentEnv: []string{"PATH=/usr/bin:/bin"},
	})
	if err != nil {
		t.Fatalf("install: %v", err)
	}

	logBytes, err := os.ReadFile(filepath.Join(logDir, "npm.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logBytes), "--prefix") || !strings.Contains(string(logBytes), targetDir) {
		t.Fatalf("expected --prefix %s in log, got %q", targetDir, string(logBytes))
	}
	if !strings.Contains(string(logBytes), "lark-cli") {
		t.Fatalf("expected source 'lark-cli' in log, got %q", string(logBytes))
	}
	if !strings.Contains(string(logBytes), "--no-audit") || !strings.Contains(string(logBytes), "--no-fund") {
		t.Fatalf("expected --no-audit and --no-fund flags, got %q", string(logBytes))
	}

	nodeLogBytes, err := os.ReadFile(filepath.Join(logDir, "node.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(nodeLogBytes), npmPath) {
		t.Fatalf("expected bundled node to receive npm path %q, got %q", npmPath, string(nodeLogBytes))
	}
}

// TestNpmInstallStripsNpmConfigPrefix verifies the inherited NPM_CONFIG_PREFIX is removed before invoking npm.
func TestNpmInstallStripsNpmConfigPrefix(t *testing.T) {
	logDir := t.TempDir()
	nodePath, npmPath := writeFakeNpm(t, logDir)
	targetDir := filepath.Join(t.TempDir(), "install")

	envInspectPath := filepath.Join(filepath.Dir(npmPath), "envcheck")
	envCheckBody := `#!/usr/bin/env bash
env | grep -E '^NPM_CONFIG_PREFIX=' > "` + filepath.Join(logDir, "env.log") + `" || true
exit 0
`
	if err := os.WriteFile(envInspectPath, []byte(envCheckBody), 0o755); err != nil {
		t.Fatal(err)
	}
	nodeInspectBody := `#!/usr/bin/env bash
script="$1"
shift
exec "$script" "$@"
`
	if err := os.WriteFile(nodePath, []byte(nodeInspectBody), 0o755); err != nil {
		t.Fatal(err)
	}

	err := NpmInstall(context.Background(), NpmInstallParams{
		NodePath:  nodePath,
		NpmPath:   envInspectPath, // use env-inspecting fake instead
		Source:    "x",
		TargetDir: targetDir,
		ParentEnv: []string{"PATH=/usr/bin:/bin", "NPM_CONFIG_PREFIX=/old/leaks"},
	})
	if err != nil {
		t.Fatalf("install: %v", err)
	}

	envLog, _ := os.ReadFile(filepath.Join(logDir, "env.log"))
	if len(envLog) != 0 {
		t.Fatalf("NPM_CONFIG_PREFIX should have been stripped, got: %q", string(envLog))
	}
}
