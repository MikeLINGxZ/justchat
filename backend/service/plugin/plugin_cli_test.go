package plugin

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// withFakeCliManager swaps in a Manager whose npm is mocked to install a synthetic lark-cli bundle.
// It also points the data dir at a tempdir via LEMONTEA_DATA_DIR. Returns the dataDir for assertions.
func withFakeCliManager(t *testing.T, p *Plugin) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("plugin cli tests use POSIX shell mocks")
	}
	dataDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", dataDir)

	binDir := t.TempDir()
	npmPath := filepath.Join(binDir, "npm")
	nodePath := filepath.Join(binDir, "node")
	nodeBody := `#!/usr/bin/env bash
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
	p.cliManager = pkgcli.NewManager(dataDir, nodePath, npmPath)
	return dataDir
}

// TestInstallCliFromNpmPersistsExtension verifies a successful install creates an ExtensionItem with kind="cli".
func TestInstallCliFromNpmPersistsExtension(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	out, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if out.Extension.Kind != "cli" || out.Extension.ID != "cli:lark-cli" {
		t.Fatalf("ext: %+v", out.Extension)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli")); err != nil {
		t.Fatalf("bin missing: %v", err)
	}
	// Config should have one extension persisted.
	cfg, err := p.loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Extensions) != 1 {
		t.Fatalf("extensions: %+v", cfg.Extensions)
	}
}

// TestResetCliDataClearsDirectory verifies ResetCliData removes the cli_data/<name>/ subtree.
func TestResetCliDataClearsDirectory(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	if _, err := p.ResetCliData(context.Background(), plugin_dto.ResetCliDataInput{ID: "cli:lark-cli"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "plugins", "cli_data", "lark-cli")); !os.IsNotExist(err) {
		t.Fatalf("cli_data should be removed: %v", err)
	}
}

// TestUpdateCliManifestPersistsAndSyncs verifies UpdateCliManifest writes the file and updates the persisted item.
func TestUpdateCliManifestPersistsAndSyncs(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	newManifest := `{
		"name":"lark-cli",
		"version":"9.9.9",
		"description":"edited",
		"executable":"` + filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli") + `",
		"isolation":"isolated",
		"tools":[]
	}`
	out, err := p.UpdateCliManifest(context.Background(), plugin_dto.UpdateCliManifestInput{ID: "cli:lark-cli", ManifestText: newManifest})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if out.Extension.Version != "9.9.9" || out.Extension.Description != "edited" {
		t.Fatalf("ext: %+v", out.Extension)
	}
}

// TestSaveExtensionConfigRoutesCliThroughManifestValidation verifies invalid JSON for a cli kind fails fast via the cli manifest validator.
func TestSaveExtensionConfigRoutesCliThroughManifestValidation(t *testing.T) {
	p := NewPlugin()
	withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}
	_, err := p.SaveExtensionConfig(context.Background(), plugin_dto.SaveExtensionConfigInput{
		ID:         "cli:lark-cli",
		ConfigText: `{ invalid json`,
	})
	if err == nil {
		t.Fatal("expected validation error for invalid manifest")
	}
}
