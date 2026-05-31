package cli

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// writePackageJSON helper writes a package.json with the given content to dir/package.json.
func writePackageJSON(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestReadPackageJSONStringBin verifies bin-as-string is honored.
func TestReadPackageJSONStringBin(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, `{"name":"lark-cli","version":"1.0.0","bin":"lark.js"}`)

	pkg, err := ReadPackageJSON(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name != "lark-cli" || pkg.Version != "1.0.0" {
		t.Fatalf("meta: %+v", pkg)
	}
	if len(pkg.Bin) != 1 || pkg.Bin["lark-cli"] != "lark.js" {
		t.Fatalf("bin map: %+v", pkg.Bin)
	}
}

// TestReadPackageJSONObjectBin verifies bin-as-object preserves all entries.
func TestReadPackageJSONObjectBin(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, `{"name":"toolset","bin":{"foo":"a.js","bar":"b.js"}}`)

	pkg, err := ReadPackageJSON(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Bin["foo"] != "a.js" || pkg.Bin["bar"] != "b.js" {
		t.Fatalf("bin map: %+v", pkg.Bin)
	}
}

// TestReadPackageJSONScopedPackage verifies scoped npm packages are discovered under node_modules/@scope/pkg/package.json.
func TestReadPackageJSONScopedPackage(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, filepath.Join(dir, "node_modules", "@larksuite", "cli"), `{"name":"@larksuite/cli","version":"1.0.0"}`)

	pkg, err := ReadPackageJSON(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name != "@larksuite/cli" {
		t.Fatalf("unexpected package meta: %+v", pkg)
	}
}

// TestSelectExecutablePrefersNameMatch verifies the bin entry matching the package name wins.
func TestSelectExecutablePrefersNameMatch(t *testing.T) {
	pkg := PackageMeta{Name: "lark-cli", Bin: map[string]string{"foo": "a.js", "lark-cli": "lark.js"}}
	got, err := SelectExecutable("/some/install", pkg)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "lark-cli" {
		t.Fatalf("expected base 'lark-cli', got %q", got)
	}
}

// TestSelectExecutableSingleBin verifies single-entry bin map is picked even when name differs.
func TestSelectExecutableSingleBin(t *testing.T) {
	pkg := PackageMeta{Name: "lark-cli", Bin: map[string]string{"lark": "lark.js"}}
	got, err := SelectExecutable("/some/install", pkg)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "lark" {
		t.Fatalf("expected base 'lark', got %q", got)
	}
}

// TestSelectExecutableErrorsWhenNoBin verifies empty bin map yields an explicit error.
func TestSelectExecutableErrorsWhenNoBin(t *testing.T) {
	if _, err := SelectExecutable("/x", PackageMeta{Name: "x", Bin: map[string]string{}}); err == nil {
		t.Fatal("expected error on missing bin")
	}
}

// TestProbeHelpReturnsOutput verifies a CLI with --help support yields its stdout.
func TestProbeHelpReturnsOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("probe tests use POSIX shell scripts")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "helpcli")
	body := "#!/usr/bin/env bash\nif [ \"$1\" = \"--help\" ]; then echo 'usage info'; exit 0; fi\nexit 1\n"
	if err := os.WriteFile(bin, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}

	out, err := ProbeHelp(context.Background(), bin, nil)
	if err != nil {
		t.Fatalf("probe: %v", err)
	}
	if out == "" {
		t.Fatal("expected non-empty help output")
	}
}

// TestProbeHelpFallsBackThroughVariants verifies the helper tries -h and help when --help fails.
func TestProbeHelpFallsBackThroughVariants(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("probe tests use POSIX shell scripts")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "helpcli")
	body := "#!/usr/bin/env bash\nif [ \"$1\" = \"help\" ]; then echo 'via help subcmd'; exit 0; fi\nexit 1\n"
	if err := os.WriteFile(bin, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}

	out, err := ProbeHelp(context.Background(), bin, nil)
	if err != nil {
		t.Fatalf("probe: %v", err)
	}
	if out == "" {
		t.Fatal("expected non-empty help output via help subcmd")
	}
}
