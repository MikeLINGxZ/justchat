package onboarding

import (
	"os"
	"path/filepath"
	"testing"
)

// withTempDataDir overrides LEMONTEA_DATA_DIR for the test scope.
func withTempDataDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tmp)
	return tmp
}

func TestIsInitialized_FalseWhenNoFile(t *testing.T) {
	withTempDataDir(t)

	got, err := isInitialized()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Fatalf("expected false, got true")
	}
}

func TestMarkInitialized_WritesFile(t *testing.T) {
	tmp := withTempDataDir(t)

	if err := markInitialized(); err != nil {
		t.Fatalf("markInitialized error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "init.json")); err != nil {
		t.Fatalf("init.json missing: %v", err)
	}

	got, err := isInitialized()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Fatalf("expected true after mark, got false")
	}
}

func TestMarkInitialized_Idempotent(t *testing.T) {
	withTempDataDir(t)
	if err := markInitialized(); err != nil {
		t.Fatalf("first mark: %v", err)
	}
	if err := markInitialized(); err != nil {
		t.Fatalf("second mark: %v", err)
	}
	got, err := isInitialized()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Fatalf("expected true, got false")
	}
}
