package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsAppInitializedAndMarkInitialized(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_PATH", tempDir)

	initialized, err := IsAppInitialized()
	if err != nil {
		t.Fatalf("IsAppInitialized() error = %v", err)
	}
	if initialized {
		t.Fatal("IsAppInitialized() = true, want false before init file exists")
	}

	if err := markAppInitialized(); err != nil {
		t.Fatalf("markAppInitialized() error = %v", err)
	}

	initialized, err = IsAppInitialized()
	if err != nil {
		t.Fatalf("IsAppInitialized() error = %v", err)
	}
	if !initialized {
		t.Fatal("IsAppInitialized() = false, want true after init file exists")
	}

	initFilePath := filepath.Join(tempDir, initFileName)
	if _, err := os.Stat(initFilePath); err != nil {
		t.Fatalf("expected init file to exist: %v", err)
	}
}
