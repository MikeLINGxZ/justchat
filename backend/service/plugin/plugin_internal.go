package plugin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

func mustDataDir() string {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return dir.HomeDir()
	}
	return dataDir
}

func (p *Plugin) loadConfig() (*data_models.Config, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(dataDir, dir.ConfigFileName)
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &data_models.Config{
				Locale:     "zh-CN",
				Language:   "zh-CN",
				FontSize:   "md",
				DataDir:    dataDir,
				LogLevel:   "info",
				Extensions: []data_models.ExtensionItem{},
			}
			if saveErr := p.saveConfig(cfg); saveErr != nil {
				return nil, saveErr
			}
			return cfg, nil
		}
		return nil, err
	}

	cfg := &data_models.Config{}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		return nil, err
	}
	if cfg.Extensions == nil {
		cfg.Extensions = []data_models.ExtensionItem{}
	}
	if cfg.DataDir == "" {
		cfg.DataDir = dataDir
	}
	return cfg, nil
}

func (p *Plugin) saveConfig(config *data_models.Config) error {
	if err := os.MkdirAll(config.DataDir, 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(config.DataDir, dir.ConfigFileName), content, 0o644)
}

func findExtension(items []data_models.ExtensionItem, id string) (data_models.ExtensionItem, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return data_models.ExtensionItem{}, false
}

func removeExtension(items []data_models.ExtensionItem, id string) []data_models.ExtensionItem {
	result := make([]data_models.ExtensionItem, 0, len(items))
	for _, item := range items {
		if item.ID != id {
			result = append(result, item)
		}
	}
	return result
}

// syncExtensionRuntime resolves runtime status and tool list for a single extension item.
// For cli plugins it projects tools from manifest.json; for mcp plugins it invokes the MCP manager to discover tools.
func (p *Plugin) syncExtensionRuntime(ctx context.Context, item data_models.ExtensionItem) data_models.ExtensionItem {
	if item.Kind == "cli" {
		manifest, err := pkgcli.LoadManifest(item.ConfigFilePath)
		if err != nil || manifest.Executable == "" {
			item.RuntimeStatus = "error"
			item.RuntimeMessage = "manifest missing; reinstall"
			item.Tools = []data_models.ExtensionTool{}
			return item
		}
		item.Tools = projectCliTools(item, manifest)
		if item.Enabled {
			item.RuntimeStatus = "ready"
			item.RuntimeMessage = ""
		} else {
			item.RuntimeStatus = "idle"
			item.RuntimeMessage = ""
		}
		return item
	}
	if item.Kind != "mcp" {
		item.RuntimeStatus = "ready"
		item.RuntimeMessage = ""
		return item
	}
	if !item.Enabled {
		item.Tools = []data_models.ExtensionTool{}
		item.RuntimeStatus = "idle"
		item.RuntimeMessage = ""
		return item
	}

	tools, err := p.mcpManager.DiscoverTools(ctx, item)
	if err != nil {
		item.Tools = []data_models.ExtensionTool{}
		item.RuntimeStatus = "error"
		item.RuntimeMessage = err.Error()
		return item
	}

	item.Tools = tools
	item.RuntimeStatus = "ready"
	item.RuntimeMessage = ""
	return item
}

func projectCliTools(item data_models.ExtensionItem, manifest pkgcli.Manifest) []data_models.ExtensionTool {
	result := make([]data_models.ExtensionTool, 0, len(manifest.Tools))
	for _, current := range manifest.Tools {
		result = append(result, data_models.ExtensionTool{
			ToolID:          pkgcli.CliToolID(item.Name, current.Name),
			ServerID:        item.ID,
			Name:            current.Name,
			Description:     current.Description,
			Enabled:         current.Enabled,
			RequiresConfirm: current.RequiresConfirm,
		})
	}
	return result
}
