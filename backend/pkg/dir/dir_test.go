package dir

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestGetDataDirUsesLocatorFileBeforeDefaultDir ensures a custom directory locator wins over the default home path.
func TestGetDataDirUsesLocatorFileBeforeDefaultDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("LEMONTEA_DATA_DIR", "")

	targetDir := filepath.Join(tempHome, "custom-data")
	metaDir := filepath.Join(tempHome, defaultBaseDirName)
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content, err := json.Marshal(map[string]string{"data_dir": targetDir})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(metaDir, "data_dir.json"), content, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := GetDataDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != targetDir {
		t.Fatalf("expected custom data dir %q, got %q", targetDir, got)
	}
}

// TestExtensionPathHelpers 验证 plugins 子目录 helper 返回的拼装路径。
func TestExtensionPathHelpers(t *testing.T) {
	base := filepath.Join("/", "tmp", "lemontea-test")

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ExtensionsRoot", ExtensionsRoot(base), filepath.Join(base, "plugins")},
		{"MCPRoot", MCPRoot(base), filepath.Join(base, "plugins", "mcp")},
		{"PluginRoot", PluginRoot(base), filepath.Join(base, "plugins", "plugin")},
		{"CLIRoot", CLIRoot(base), filepath.Join(base, "plugins", "cli")},
		{"CLIDataRoot", CLIDataRoot(base), filepath.Join(base, "plugins", "cli_data")},
		{"LegacyMCPRoot", LegacyMCPRoot(base), filepath.Join(base, "mcp")},
		{"LegacyPluginRoot", LegacyPluginRoot(base), filepath.Join(base, "plugin")},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Fatalf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestSkillsRoot(t *testing.T) {
	got := SkillsRoot("/tmp/data")
	want := filepath.Join("/tmp/data", "skills")
	if got != want {
		t.Fatalf("SkillsRoot mismatch: got %q want %q", got, want)
	}
}
