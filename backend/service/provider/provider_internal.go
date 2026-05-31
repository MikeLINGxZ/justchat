package provider

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

// toProviderViewModel maps a data_models.Provider row to its view_model representation.
func toProviderViewModel(p data_models.Provider, isDefault bool) view_model.Provider {
	return view_model.Provider{
		ID:        p.ID,
		Icon:      providerIconPath(p.ProviderType),
		Type:      p.ProviderType,
		Name:      p.ProviderName,
		BaseURL:   p.BaseUrl,
		Enabled:   p.Enable,
		ApiKey:    p.ApiKey,
		IsDefault: isDefault,
	}
}

// toModelViewModel maps a data_models.Model row to its view_model representation.
func toModelViewModel(m data_models.Model, isDefault bool) view_model.Model {
	return view_model.Model{
		ID:         m.ID,
		ProviderId: m.ProviderId,
		Model:      m.Model,
		OwnedBy:    m.OwnedBy,
		Object:     m.Object,
		Enable:     m.Enable,
		Alias:      m.Alias,
		IsCustom:   m.IsCustom,
		IsDefault:  isDefault,
	}
}

// providerIconPath returns the bundled icon asset path for a provider type.
func providerIconPath(t pkgProvider.Type) string {
	icons := map[pkgProvider.Type]string{
		pkgProvider.Deepseek:            "/providers/deepseek_icon.png",
		pkgProvider.Aliyun:              "/providers/qwen_icon.png",
		pkgProvider.Ollama:              "/providers/ollama_icon.png",
		pkgProvider.OpenAiCompatibility: "/providers/openai_icon.png",
		pkgProvider.Openrouter:          "/providers/openrouter_icon.png",
	}
	if icon, ok := icons[t]; ok {
		return icon
	}
	return ""
}

// readDefaultProviderID reads DefaultProviderID from the config file, returning 0 on any error.
func readDefaultProviderID() uint {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return 0
	}
	bytes, err := os.ReadFile(filepath.Join(dataDir, dir.ConfigFileName))
	if err != nil {
		return 0
	}
	var cfg data_models.Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return 0
	}
	return cfg.DefaultProviderID
}

// writeDefaultProviderID persists DefaultProviderID in the config file; preserves all other config fields.
func writeDefaultProviderID(providerID uint) error {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return ierror.Error(ierror.ErrProviderSetDefault, err)
	}
	configPath := filepath.Join(dataDir, dir.ConfigFileName)

	var cfg data_models.Config
	if raw, readErr := os.ReadFile(configPath); readErr == nil {
		_ = json.Unmarshal(raw, &cfg)
	}
	cfg.DefaultProviderID = providerID

	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return ierror.Error(ierror.ErrProviderSetDefault, err)
	}
	return ierror.Error(ierror.ErrProviderSetDefault, os.WriteFile(configPath, content, 0o644))
}
