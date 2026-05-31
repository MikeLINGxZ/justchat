package mcp

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Manifest describes the import metadata for one MCP bundle directory.
type Manifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	ConfigFile  string `json:"config_file"`
}

// loadManifest reads the bundle manifest and falls back to directory defaults when fields are missing.
func loadManifest(rootDir string) (Manifest, error) {
	manifestPath := filepath.Join(rootDir, "manifest.json")
	bytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return Manifest{}, err
	}

	manifest := Manifest{}
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		return Manifest{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(rootDir)
	}
	if manifest.ConfigFile == "" {
		manifest.ConfigFile = "mcp.json"
	}
	return manifest, nil
}

// loadManifestOptional reads manifest.json when present and otherwise falls back to directory defaults.
func loadManifestOptional(rootDir string) (Manifest, error) {
	manifest, err := loadManifest(rootDir)
	if err == nil {
		return manifest, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return Manifest{}, err
	}
	return Manifest{
		Name: filepath.Base(rootDir),
	}, nil
}

// writeManifest persists a normalized manifest into the copied extension directory.
func writeManifest(rootDir string, manifest Manifest) error {
	bytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(rootDir, "manifest.json"), bytes, 0o644)
}
