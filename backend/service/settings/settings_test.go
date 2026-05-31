package settings

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/settings/settings_dto"
)

// newTestSettingsService creates a Settings instance for testing.
func newTestSettingsService(t *testing.T) *Settings {
	t.Helper()
	return &Settings{}
}

// TestApplyFileSettingsMigratesDataAndWritesLocator verifies data migration and locator persistence for custom directories.
func TestApplyFileSettingsMigratesDataAndWritesLocator(t *testing.T) {
	tempHome := t.TempDir()
	sourceDir := t.TempDir()
	targetDir := filepath.Join(t.TempDir(), "custom-target")

	t.Setenv("HOME", tempHome)
	t.Setenv("LEMONTEA_DATA_DIR", sourceDir)

	svc := newTestSettingsService(t)

	if err := os.WriteFile(filepath.Join(sourceDir, "config.json"), []byte(`{"language":"zh-CN","data_dir":"`+sourceDir+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "data.db"), []byte("sqlite"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := svc.ApplyFileSettings(context.Background(), settings_dto.ApplyFileSettingsInput{TargetDir: targetDir})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "config.json")); err != nil {
		t.Fatalf("expected migrated config file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetDir, "data.db")); err != nil {
		t.Fatalf("expected migrated database file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetDir, "logs")); err != nil {
		t.Fatalf("expected logs directory: %v", err)
	}

	locatorBytes, err := os.ReadFile(filepath.Join(tempHome, ".lemontea", "data_dir.json"))
	if err != nil {
		t.Fatal(err)
	}

	var locator map[string]string
	if err := json.Unmarshal(locatorBytes, &locator); err != nil {
		t.Fatal(err)
	}
	if locator["data_dir"] != targetDir {
		t.Fatalf("expected locator to point at %q, got %q", targetDir, locator["data_dir"])
	}
}

// TestSaveLocaleSettingsPersistsLocaleAndLanguage verifies locale settings are written back to config.json.
func TestSaveLocaleSettingsPersistsLocaleAndLanguage(t *testing.T) {
	tempDataDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDataDir)

	svc := newTestSettingsService(t)
	if err := os.WriteFile(filepath.Join(tempDataDir, "config.json"), []byte(`{"locale":"zh-CN","language":"zh-CN","font_size":"md","data_dir":"`+tempDataDir+`","log_level":"info"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := svc.SaveLocaleSettings(context.Background(), settings_dto.SaveLocaleSettingsInput{
		Locale:   "en-US",
		Language: "en",
	})
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := os.ReadFile(filepath.Join(tempDataDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}

	var cfg data_models.Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Locale != "en-US" {
		t.Fatalf("expected locale to be persisted, got %q", cfg.Locale)
	}
	if cfg.Language != "en" {
		t.Fatalf("expected language to be persisted, got %q", cfg.Language)
	}
}

// TestSaveDisplaySettingsPersistsFontSize verifies display settings are written back to config.json.
func TestSaveDisplaySettingsPersistsFontSize(t *testing.T) {
	tempDataDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDataDir)

	svc := newTestSettingsService(t)
	if err := os.WriteFile(filepath.Join(tempDataDir, "config.json"), []byte(`{"locale":"zh-CN","language":"zh-CN","font_size":"md","data_dir":"`+tempDataDir+`","log_level":"info"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := svc.SaveDisplaySettings(context.Background(), settings_dto.SaveDisplaySettingsInput{FontSize: "lg"})
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := os.ReadFile(filepath.Join(tempDataDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}

	var cfg data_models.Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.FontSize != "lg" {
		t.Fatalf("expected font_size to be persisted, got %q", cfg.FontSize)
	}
}
