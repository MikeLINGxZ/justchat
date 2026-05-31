package skills

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// loadConfig reads the app config from config.json in the data directory.
// Returns a sensible default when the file does not yet exist.
func (s *Skills) loadConfig() (*data_models.Config, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(dataDir, dir.ConfigFileName)
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &data_models.Config{DataDir: dataDir}, nil
		}
		return nil, err
	}
	cfg := &data_models.Config{}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		return nil, err
	}
	if cfg.DataDir == "" {
		cfg.DataDir = dataDir
	}
	return cfg, nil
}

// saveConfig writes the app config to config.json in the data directory.
func (s *Skills) saveConfig(cfg *data_models.Config) error {
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cfg.DataDir, dir.ConfigFileName), content, 0o644)
}

// renameDisabledSkill rewrites a disabled skill entry when a user skill is renamed.
func (s *Skills) renameDisabledSkill(oldName string, newName string) error {
	cfg, err := s.loadConfig()
	if err != nil {
		return err
	}
	changed := false
	for i, name := range cfg.DisabledSkills {
		if name == oldName {
			cfg.DisabledSkills[i] = newName
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return s.saveConfig(cfg)
}
