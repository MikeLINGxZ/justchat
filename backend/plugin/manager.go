package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

// PluginState represents the current lifecycle state of a plugin.
type PluginState string

const (
	PluginStateInstalled PluginState = "installed"
	PluginStateActive    PluginState = "active"
	PluginStateError     PluginState = "error"
)

// PluginInfo holds runtime information about an installed plugin.
type PluginInfo struct {
	ID       string      `json:"id"`
	Dir      string      `json:"dir"`
	Manifest *Manifest   `json:"manifest"`
	State    PluginState `json:"state"`
	Enabled  bool        `json:"enabled"`
	Error    string      `json:"error,omitempty"`
}

// Manager is the central coordinator for all plugin operations.
type Manager struct {
	mu            sync.RWMutex
	pluginsDir    string
	plugins       map[string]*PluginInfo
	host          *ExtensionHost
	hookChain     *HookChain
	toolBridge    *ToolBridge
	agentBridge   *AgentBridge
	storage       *PluginStorage
	app           *application.App
	lastCrashTime time.Time // for crash loop detection
	crashCount    int       // consecutive crash count
}

// NewManager creates a new plugin Manager. It sets up the plugins directory and
// initialises the extension host, tool bridge, agent bridge and storage.
func NewManager(app *application.App) (*Manager, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return nil, fmt.Errorf("get data path: %w", err)
	}

	pluginsDir := filepath.Join(dataPath, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return nil, fmt.Errorf("create plugins directory: %w", err)
	}

	// Resolve extension host script path.
	// Priority: env var > project dir (for development) > data dir
	hostScriptPath := os.Getenv("LEMONTEA_EXT_HOST_PATH")
	if hostScriptPath == "" {
		// Try the project directory first (development mode).
		if wd, wdErr := os.Getwd(); wdErr == nil {
			candidate := filepath.Join(wd, "extension-host", "dist", "index.js")
			if _, statErr := os.Stat(candidate); statErr == nil {
				hostScriptPath = candidate
			}
		}
	}
	if hostScriptPath == "" {
		hostScriptPath = filepath.Join(dataPath, "extension-host", "dist", "index.js")
	}

	host := NewExtensionHost(hostScriptPath, pluginsDir)
	toolBridge := NewToolBridge()
	agentBridge := NewAgentBridge()

	storage, err := NewPluginStorage()
	if err != nil {
		return nil, fmt.Errorf("create plugin storage: %w", err)
	}

	m := &Manager{
		pluginsDir:  pluginsDir,
		plugins:     make(map[string]*PluginInfo),
		host:        host,
		toolBridge:  toolBridge,
		agentBridge: agentBridge,
		storage:     storage,
		app:         app,
	}

	return m, nil
}

// Init starts the extension host, registers RPC handlers, and loads all installed plugins.
func (m *Manager) Init() error {
	// Verify Node.js is available.
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("node.js is required but not found in PATH: %w", err)
	}

	// Verify the extension host script exists.
	if _, err := os.Stat(m.host.hostScriptPath); err != nil {
		return fmt.Errorf("extension host script not found at %s: %w", m.host.hostScriptPath, err)
	}

	// Start the Extension Host.
	if err := m.host.Start(); err != nil {
		return fmt.Errorf("start extension host: %w", err)
	}

	// Create the HookChain now that the host is running.
	m.hookChain = NewHookChain(m.host)

	// Register RPC handlers for storage and events.
	m.registerRPCHandlers()

	// Set crash callback.
	m.host.SetOnCrash(func() {
		m.handleHostCrash()
	})

	// Load enabled states from persistent config.
	enabledStates := m.loadEnabledStates()

	// Scan plugins directory.
	entries, err := os.ReadDir(m.pluginsDir)
	if err != nil {
		return fmt.Errorf("scan plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(m.pluginsDir, entry.Name())
		manifest, err := ParseManifest(pluginDir)
		if err != nil {
			logger.Warmf("skipping directory %s: %v", entry.Name(), err)
			continue
		}

		pluginID := manifest.Name
		enabled := true
		if v, ok := enabledStates[pluginID]; ok {
			enabled = v
		}

		info := &PluginInfo{
			ID:       pluginID,
			Dir:      pluginDir,
			Manifest: manifest,
			State:    PluginStateInstalled,
			Enabled:  enabled,
		}

		m.mu.Lock()
		m.plugins[pluginID] = info
		m.mu.Unlock()
	}

	// Activate all enabled plugins that have the "onStartup" activation event.
	m.mu.RLock()
	var toActivate []string
	for id, info := range m.plugins {
		if !info.Enabled {
			continue
		}
		for _, event := range info.Manifest.LemonTea.ActivationEvents {
			if event == "onStartup" {
				toActivate = append(toActivate, id)
				break
			}
		}
	}
	m.mu.RUnlock()

	for _, id := range toActivate {
		if err := m.Activate(id); err != nil {
			logger.Errorf("failed to activate plugin %s: %v", id, err)
		}
	}

	logger.Infof("plugin manager initialised: %d plugins loaded, %d activated", len(m.plugins), len(toActivate))
	return nil
}

// Shutdown deactivates all active plugins, stops the extension host and closes storage.
func (m *Manager) Shutdown() {
	m.mu.RLock()
	var activeIDs []string
	for id, info := range m.plugins {
		if info.State == PluginStateActive {
			activeIDs = append(activeIDs, id)
		}
	}
	m.mu.RUnlock()

	for _, id := range activeIDs {
		if err := m.Deactivate(id); err != nil {
			logger.Errorf("failed to deactivate plugin %s during shutdown: %v", id, err)
		}
	}

	_ = m.host.Stop()
	_ = m.storage.Close()
}

// Install copies a plugin from folderPath into the managed plugins directory,
// runs npm install, and activates it.
func (m *Manager) Install(folderPath string) error {
	manifest, err := ParseManifest(folderPath)
	if err != nil {
		return fmt.Errorf("parse manifest: %w", err)
	}
	if err := ValidateManifest(manifest); err != nil {
		return fmt.Errorf("validate manifest: %w", err)
	}

	pluginID := manifest.Name

	// Check for duplicate.
	m.mu.RLock()
	_, exists := m.plugins[pluginID]
	m.mu.RUnlock()
	if exists {
		return fmt.Errorf("plugin %s is already installed", pluginID)
	}

	// Copy to plugins directory.
	targetDir := filepath.Join(m.pluginsDir, manifest.Name)
	if err := copyDir(folderPath, targetDir); err != nil {
		return fmt.Errorf("copy plugin directory: %w", err)
	}

	// Run cnpm install with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cnpm", "install", "--production")
	cmd.Dir = targetDir
	if output, err := cmd.CombinedOutput(); err != nil {
		// Clean up on failure.
		_ = os.RemoveAll(targetDir)
		return fmt.Errorf("cnpm install failed: %w\noutput: %s", err, string(output))
	}

	// Re-parse manifest from installed location.
	manifest, err = ParseManifest(targetDir)
	if err != nil {
		_ = os.RemoveAll(targetDir)
		return fmt.Errorf("re-parse manifest after install: %w", err)
	}

	info := &PluginInfo{
		ID:       pluginID,
		Dir:      targetDir,
		Manifest: manifest,
		State:    PluginStateInstalled,
		Enabled:  true,
	}

	m.mu.Lock()
	m.plugins[pluginID] = info
	m.mu.Unlock()

	m.saveEnabledStates()

	if err := m.Activate(pluginID); err != nil {
		logger.Errorf("failed to activate plugin %s after install: %v", pluginID, err)
	}

	m.app.Event.Emit("plugin:installed", pluginID)
	m.app.Event.Emit("plugin:changed", pluginID)
	return nil
}

// Uninstall deactivates and removes a plugin entirely.
func (m *Manager) Uninstall(pluginId string) error {
	m.mu.RLock()
	info, ok := m.plugins[pluginId]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginId)
	}

	// Deactivate if active.
	if info.State == PluginStateActive {
		if err := m.Deactivate(pluginId); err != nil {
			logger.Errorf("failed to deactivate plugin %s during uninstall: %v", pluginId, err)
		}
	}

	// Delete storage.
	if err := m.storage.DeleteAll(pluginId); err != nil {
		logger.Errorf("failed to delete storage for plugin %s: %v", pluginId, err)
	}

	// Remove plugin directory.
	if err := os.RemoveAll(info.Dir); err != nil {
		return fmt.Errorf("remove plugin directory: %w", err)
	}

	m.mu.Lock()
	delete(m.plugins, pluginId)
	m.mu.Unlock()

	m.saveEnabledStates()

	m.app.Event.Emit("plugin:uninstalled", pluginId)
	m.app.Event.Emit("plugin:changed", pluginId)
	return nil
}

// Activate sends an activation RPC to the extension host and registers the
// plugin's contributed tools, hooks and agents.
func (m *Manager) Activate(pluginId string) error {
	m.mu.RLock()
	info, ok := m.plugins[pluginId]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginId)
	}

	if !m.host.IsRunning() {
		return fmt.Errorf("extension host is not running")
	}

	// Send activate RPC.
	_, err := m.host.RPC().Call("plugin/activate", map[string]any{
		"pluginId":  pluginId,
		"pluginDir": info.Dir,
	})
	if err != nil {
		m.mu.Lock()
		info.State = PluginStateError
		info.Error = err.Error()
		m.mu.Unlock()
		return fmt.Errorf("activate plugin %s: %w", pluginId, err)
	}

	// Register contributed tools.
	for _, tool := range info.Manifest.LemonTea.Contributes.Tools {
		m.toolBridge.Register(pluginId, info.Manifest.Name, info.Manifest.DisplayName, tool, m.host)
	}

	// Register contributed hooks.
	hooks := info.Manifest.LemonTea.Contributes.Hooks
	if hooks.OnBeforeChat {
		m.hookChain.Register(pluginId, "onBeforeChat")
	}
	if hooks.OnAfterChat {
		m.hookChain.Register(pluginId, "onAfterChat")
	}

	// Register contributed agents.
	for _, agent := range info.Manifest.LemonTea.Contributes.Agents {
		m.agentBridge.Register(pluginId, info.Manifest.Name, agent, "", nil, agents.AgentRoleWorker)
	}

	m.mu.Lock()
	info.State = PluginStateActive
	info.Error = ""
	m.mu.Unlock()

	return nil
}

// Deactivate sends a deactivation RPC and unregisters all contributed extensions.
func (m *Manager) Deactivate(pluginId string) error {
	m.mu.RLock()
	info, ok := m.plugins[pluginId]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginId)
	}

	if info.State != PluginStateActive {
		return nil
	}

	// Send deactivate RPC (best-effort).
	if m.host.IsRunning() {
		_, _ = m.host.RPC().Call("plugin/deactivate", map[string]any{
			"pluginId": pluginId,
		})
	}

	// Unregister all extensions.
	m.toolBridge.UnregisterByPlugin(pluginId)
	m.hookChain.Unregister(pluginId)
	m.agentBridge.UnregisterByPlugin(pluginId)

	m.mu.Lock()
	info.State = PluginStateInstalled
	info.Error = ""
	m.mu.Unlock()

	return nil
}

// Enable marks a plugin as enabled, activates it, and persists the state.
func (m *Manager) Enable(pluginId string) error {
	m.mu.Lock()
	info, ok := m.plugins[pluginId]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %s not found", pluginId)
	}
	info.Enabled = true
	m.mu.Unlock()

	m.saveEnabledStates()

	err := m.Activate(pluginId)
	m.app.Event.Emit("plugin:changed", pluginId)
	return err
}

// Disable marks a plugin as disabled, deactivates it, and persists the state.
func (m *Manager) Disable(pluginId string) error {
	m.mu.Lock()
	info, ok := m.plugins[pluginId]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %s not found", pluginId)
	}
	info.Enabled = false
	m.mu.Unlock()

	m.saveEnabledStates()

	err := m.Deactivate(pluginId)
	m.app.Event.Emit("plugin:changed", pluginId)
	return err
}

// ListPlugins returns a snapshot of all installed plugin infos.
func (m *Manager) ListPlugins() []*PluginInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*PluginInfo, 0, len(m.plugins))
	for _, info := range m.plugins {
		result = append(result, info)
	}
	return result
}

// GetPlugin returns the PluginInfo for the given ID, or false if not found.
func (m *Manager) GetPlugin(pluginId string) (*PluginInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	info, ok := m.plugins[pluginId]
	return info, ok
}

// HookChain returns the hook chain instance.
func (m *Manager) HookChain() *HookChain {
	return m.hookChain
}

// GetPluginTools returns all tools from all active plugins.
func (m *Manager) GetPluginTools() []*PluginTool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PluginTool
	for id, info := range m.plugins {
		if info.State == PluginStateActive {
			result = append(result, m.toolBridge.GetToolsByPlugin(id)...)
		}
	}
	return result
}

// registerRPCHandlers sets up the JSON-RPC handlers for storage and event emission.
func (m *Manager) registerRPCHandlers() {
	rpc := m.host.RPC()

	// Storage: get
	rpc.RegisterHandler("app/storage/get", func(params json.RawMessage) (any, error) {
		var req struct {
			PluginID string `json:"pluginId"`
			Key      string `json:"key"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal storage/get params: %w", err)
		}
		val, err := m.storage.Get(req.PluginID, req.Key)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		return json.RawMessage(val), nil
	})

	// Storage: set
	rpc.RegisterHandler("app/storage/set", func(params json.RawMessage) (any, error) {
		var req struct {
			PluginID string          `json:"pluginId"`
			Key      string          `json:"key"`
			Value    json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal storage/set params: %w", err)
		}
		return nil, m.storage.Set(req.PluginID, req.Key, []byte(req.Value))
	})

	// Storage: delete
	rpc.RegisterHandler("app/storage/delete", func(params json.RawMessage) (any, error) {
		var req struct {
			PluginID string `json:"pluginId"`
			Key      string `json:"key"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal storage/delete params: %w", err)
		}
		return nil, m.storage.Delete(req.PluginID, req.Key)
	})

	// Dynamic tool registration from plugins at runtime
	rpc.RegisterHandler("app/registerTool", func(params json.RawMessage) (any, error) {
		var req struct {
			PluginID    string         `json:"pluginId"`
			ToolID      string         `json:"toolId"`
			Description string         `json:"description"`
			Parameters  map[string]any `json:"parameters"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal registerTool params: %w", err)
		}
		m.mu.RLock()
		info, ok := m.plugins[req.PluginID]
		m.mu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("plugin not found: %s", req.PluginID)
		}
		contrib := ToolContrib{
			ID:          req.ToolID,
			Name:        req.ToolID,
			Description: req.Description,
			Parameters:  req.Parameters,
		}
		m.toolBridge.Register(req.PluginID, info.Manifest.Name, info.Manifest.DisplayName, contrib, m.host)
		return nil, nil
	})

	// Dynamic agent registration from plugins at runtime
	rpc.RegisterHandler("app/registerAgent", func(params json.RawMessage) (any, error) {
		var req struct {
			PluginID     string   `json:"pluginId"`
			AgentID      string   `json:"agentId"`
			Name         string   `json:"name"`
			Description  string   `json:"description"`
			SystemPrompt string   `json:"systemPrompt"`
			Tools        []string `json:"tools"`
			Role         string   `json:"role"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal registerAgent params: %w", err)
		}
		m.mu.RLock()
		info, ok := m.plugins[req.PluginID]
		m.mu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("plugin not found: %s", req.PluginID)
		}
		role := agents.AgentRoleWorker
		if req.Role == "main" {
			role = agents.AgentRoleMain
		}
		contrib := AgentContrib{
			ID:          req.AgentID,
			Name:        req.Name,
			Description: req.Description,
		}
		m.agentBridge.Register(req.PluginID, info.Manifest.Name, contrib, req.SystemPrompt, req.Tools, role)
		return nil, nil
	})

	// App config access
	rpc.RegisterHandler("app/getConfig", func(params json.RawMessage) (any, error) {
		var req struct {
			Key string `json:"key"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal getConfig params: %w", err)
		}
		// For now, return nil for all config keys - can be extended later
		return nil, nil
	})

	// Event emission
	rpc.RegisterHandler("app/emitEvent", func(params json.RawMessage) (any, error) {
		var req struct {
			Event string `json:"event"`
			Data  any    `json:"data"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("unmarshal emitEvent params: %w", err)
		}
		m.app.Event.Emit("plugin:"+req.Event, req.Data)
		return nil, nil
	})
}

// saveEnabledStates persists the enabled/disabled state of all plugins to disk.
func (m *Manager) saveEnabledStates() {
	m.mu.RLock()
	states := make(map[string]bool, len(m.plugins))
	for id, info := range m.plugins {
		states[id] = info.Enabled
	}
	m.mu.RUnlock()

	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal plugin states: %v", err)
		return
	}

	path := filepath.Join(m.pluginsDir, "plugins-state.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		logger.Errorf("failed to write plugin states: %v", err)
	}
}

// loadEnabledStates reads the persisted enabled/disabled states from disk.
func (m *Manager) loadEnabledStates() map[string]bool {
	path := filepath.Join(m.pluginsDir, "plugins-state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]bool)
	}

	var states map[string]bool
	if err := json.Unmarshal(data, &states); err != nil {
		logger.Errorf("failed to parse plugin states: %v", err)
		return make(map[string]bool)
	}
	return states
}

// handleHostCrash is called when the extension host process exits unexpectedly.
// It marks all active plugins as errored and attempts to restart the host.
// Includes crash loop protection: gives up after 3 consecutive crashes within 30 seconds.
func (m *Manager) handleHostCrash() {
	logger.Error("extension host crashed, attempting recovery")

	// Crash loop protection: track consecutive crashes.
	now := time.Now()
	m.mu.Lock()
	if now.Sub(m.lastCrashTime) < 30*time.Second {
		m.crashCount++
	} else {
		m.crashCount = 1
	}
	m.lastCrashTime = now
	crashCount := m.crashCount
	m.mu.Unlock()

	if crashCount >= 3 {
		logger.Error("extension host crash loop detected (3 crashes in 30s), giving up")
		m.mu.Lock()
		for _, info := range m.plugins {
			if info.State == PluginStateActive {
				info.State = PluginStateError
				info.Error = "extension host crash loop"
			}
		}
		m.mu.Unlock()
		return
	}

	// Mark all active plugins as errored.
	m.mu.Lock()
	var wasActive []string
	for id, info := range m.plugins {
		if info.State == PluginStateActive {
			info.State = PluginStateError
			info.Error = "extension host crashed"
			wasActive = append(wasActive, id)
		}
	}
	m.mu.Unlock()

	// Try to restart the host.
	if err := m.host.Restart(); err != nil {
		logger.Errorf("failed to restart extension host: %v", err)
		return
	}

	// Recreate hook chain with the new host connection.
	m.hookChain = NewHookChain(m.host)
	m.registerRPCHandlers()
	m.host.SetOnCrash(func() {
		m.handleHostCrash()
	})

	// Re-activate plugins that were active before the crash.
	for _, id := range wasActive {
		if err := m.Activate(id); err != nil {
			logger.Errorf("failed to re-activate plugin %s after host crash: %v", id, err)
		}
	}

	logger.Infof("extension host recovered, %d plugins re-activated", len(wasActive))
}

// copyDir recursively copies the directory at src to dst, skipping node_modules.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules to avoid copying potentially huge dependency trees.
		// cnpm install will recreate them in the target directory.
		if d.IsDir() && d.Name() == "node_modules" {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
