package settings

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/version"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/settings/settings_dto"
)

// Settings manages config persistence and file-system level settings actions.
type Settings struct {
}

// LoadBootstrap loads the config file and provider list required by the settings frontend entry.
func (s *Settings) LoadBootstrap(ctx context.Context, input settings_dto.LoadBootstrapInput) (*settings_dto.LoadBootstrapOutput, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}

	return &settings_dto.LoadBootstrapOutput{
		Bootstrap: view_model.SettingsBootstrap{
			Locale:            config.Locale,
			Language:          config.Language,
			FontSize:          config.FontSize,
			DataDir:           config.DataDir,
			LogLevel:          config.LogLevel,
			DefaultProviderID: config.DefaultProviderID,
			Version:           version.ApplicationVersion,
		},
	}, nil
}

// ApplyFileSettings migrates config and database artifacts into a new data directory and updates the locator file.
func (s *Settings) ApplyFileSettings(ctx context.Context, input settings_dto.ApplyFileSettingsInput) (*settings_dto.ApplyFileSettingsOutput, error) {
	sourceDir, err := dir.GetDataDir()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	if input.TargetDir == "" {
		return nil, ierror.Error(ierror.ErrSettingsTargetDir, fmt.Errorf("target data dir is required"))
	}
	if err := os.MkdirAll(input.TargetDir, 0o755); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsCreateDir, err)
	}
	if err := os.MkdirAll(filepath.Join(input.TargetDir, "logs"), 0o755); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsCreateDir, err)
	}

	config, err := s.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	config.DataDir = input.TargetDir

	if err := copyFile(filepath.Join(sourceDir, dir.DataBaseFileName), filepath.Join(input.TargetDir, dir.DataBaseFileName)); err != nil && !os.IsNotExist(err) {
		return nil, ierror.Error(ierror.ErrSettingsCopyFile, err)
	}
	if err := s.saveConfigToDir(input.TargetDir, config); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}
	if err := dir.WriteLocatorDataDir(input.TargetDir); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsWriteLocator, err)
	}

	return &settings_dto.ApplyFileSettingsOutput{}, nil
}

// SaveLocaleSettings updates locale and language in the persisted config file.
func (s *Settings) SaveLocaleSettings(ctx context.Context, input settings_dto.SaveLocaleSettingsInput) (*settings_dto.SaveLocaleSettingsOutput, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}

	config.Locale = input.Locale
	config.Language = input.Language

	err = s.saveConfigToDir(config.DataDir, config)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}

	return &settings_dto.SaveLocaleSettingsOutput{}, nil
}

// SaveDisplaySettings updates the persisted application font-size preference.
func (s *Settings) SaveDisplaySettings(ctx context.Context, input settings_dto.SaveDisplaySettingsInput) (*settings_dto.SaveDisplaySettingsOutput, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}

	config.FontSize = input.FontSize

	err = s.saveConfigToDir(config.DataDir, config)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}

	return &settings_dto.SaveDisplaySettingsOutput{}, nil
}
