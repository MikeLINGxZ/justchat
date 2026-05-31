package skills

import (
	"testing"
)

func TestLoadBuiltin_HasInstallCliFromDocs(t *testing.T) {
	got, err := LoadBuiltin()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 1 {
		t.Fatalf("expected at least 1 builtin skill, got %d", len(got))
	}
	found := false
	for _, s := range got {
		if s.Name == "install-cli-from-docs" {
			found = true
			if s.Source != SourceBuiltin {
				t.Errorf("expected SourceBuiltin, got %q", s.Source)
			}
		}
	}
	if !found {
		t.Fatal("install-cli-from-docs builtin skill not found")
	}
}
