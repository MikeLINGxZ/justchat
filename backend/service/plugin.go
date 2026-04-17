package service

import (
	"fmt"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func (s *Service) GetInstalledPlugins() []view_models.PluginSummary {
	if s.pluginManager == nil {
		return nil
	}
	plugins := s.pluginManager.ListPlugins()
	result := make([]view_models.PluginSummary, 0, len(plugins))
	for _, p := range plugins {
		viewCount := 0
		hookCount := 0
		var toolItems []view_models.PluginContribItem
		var agentItems []view_models.PluginContribItem
		var hooks []string
		if p.Manifest != nil {
			c := p.Manifest.LemonTea.Contributes
			for _, t := range c.Tools {
				toolItems = append(toolItems, view_models.PluginContribItem{ID: t.ID, Name: t.Name})
			}
			for _, a := range c.Agents {
				agentItems = append(agentItems, view_models.PluginContribItem{ID: a.ID, Name: a.Name})
			}
			viewCount = len(c.Views.Sidebar) + len(c.Views.ChatCard) + len(c.Views.Settings) + len(c.Views.Page)
			if c.Hooks.OnBeforeChat {
				hookCount++
				hooks = append(hooks, "onBeforeChat")
			}
			if c.Hooks.OnAfterChat {
				hookCount++
				hooks = append(hooks, "onAfterChat")
			}
		}
		if toolItems == nil {
			toolItems = []view_models.PluginContribItem{}
		}
		if agentItems == nil {
			agentItems = []view_models.PluginContribItem{}
		}
		if hooks == nil {
			hooks = []string{}
		}
		result = append(result, view_models.PluginSummary{
			ID:          p.ID,
			DisplayName: p.Manifest.DisplayName,
			Version:     p.Manifest.Version,
			Description: p.Manifest.Description,
			Enabled:     p.Enabled,
			State:       string(p.State),
			ToolCount:   len(toolItems),
			ViewCount:   viewCount,
			HookCount:   hookCount,
			AgentCount:  len(agentItems),
			Tools:       toolItems,
			Agents:      agentItems,
			Hooks:       hooks,
		})
	}
	return result
}

func (s *Service) InstallPlugin() error {
	if s.pluginManager == nil {
		return fmt.Errorf("plugin system not initialized")
	}
	path, err := s.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("Select Plugin Folder").
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	return s.pluginManager.Install(path)
}

func (s *Service) UninstallPlugin(pluginId string) error {
	if s.pluginManager == nil {
		return fmt.Errorf("plugin system not initialized")
	}
	return s.pluginManager.Uninstall(pluginId)
}

func (s *Service) EnablePlugin(pluginId string) error {
	if s.pluginManager == nil {
		return fmt.Errorf("plugin system not initialized")
	}
	return s.pluginManager.Enable(pluginId)
}

func (s *Service) DisablePlugin(pluginId string) error {
	if s.pluginManager == nil {
		return fmt.Errorf("plugin system not initialized")
	}
	return s.pluginManager.Disable(pluginId)
}

// GetPluginAssetPath returns the absolute file path for a plugin UI asset.
// This allows the frontend to construct file:// URLs for plugin iframes.
func (s *Service) GetPluginAssetPath(pluginId, assetPath string) (string, error) {
	if s.pluginManager == nil {
		return "", fmt.Errorf("plugin system not initialized")
	}
	info, ok := s.pluginManager.GetPlugin(pluginId)
	if !ok {
		return "", fmt.Errorf("plugin not found: %s", pluginId)
	}
	return filepath.Join(info.Dir, assetPath), nil
}

// GetPluginViews returns the view contributions for all active plugins.
func (s *Service) GetPluginViews() map[string]any {
	if s.pluginManager == nil {
		return nil
	}
	result := map[string]any{}
	for _, p := range s.pluginManager.ListPlugins() {
		if p.State != "active" || p.Manifest == nil {
			continue
		}
		views := p.Manifest.LemonTea.Contributes.Views
		if len(views.Sidebar) > 0 || len(views.ChatCard) > 0 || len(views.Settings) > 0 || len(views.Page) > 0 {
			result[p.ID] = map[string]any{
				"pluginDir": p.Dir,
				"sidebar":   views.Sidebar,
				"chatCard":  views.ChatCard,
				"settings":  views.Settings,
				"page":      views.Page,
			}
		}
	}
	return result
}
