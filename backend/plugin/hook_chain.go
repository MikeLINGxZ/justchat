package plugin

import (
	"encoding/json"
	"sync"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// ChatContext holds the data passed through the hook pipeline.
type ChatContext struct {
	Messages []map[string]any `json:"messages"`
	AgentID  string           `json:"agentId"`
	Tools    []string         `json:"tools,omitempty"`
	Response string           `json:"response,omitempty"`
}

// HookEntry represents a single plugin registration for a hook.
type HookEntry struct {
	PluginID string
}

// HookChain manages before/after chat hook pipelines.
type HookChain struct {
	mu         sync.RWMutex
	beforeChat []HookEntry
	afterChat  []HookEntry
	host       *ExtensionHost
}

type hookCallParams struct {
	PluginID string       `json:"pluginId"`
	Context  *ChatContext `json:"context"`
}

// NewHookChain creates a new HookChain bound to the given ExtensionHost.
func NewHookChain(host *ExtensionHost) *HookChain {
	return &HookChain{
		host: host,
	}
}

// Register adds a hook entry for the given plugin. hookType must be
// "onBeforeChat" or "onAfterChat". Duplicate registrations are ignored.
func (c *HookChain) Register(pluginId string, hookType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch hookType {
	case "onBeforeChat":
		for _, e := range c.beforeChat {
			if e.PluginID == pluginId {
				return
			}
		}
		c.beforeChat = append(c.beforeChat, HookEntry{PluginID: pluginId})
	case "onAfterChat":
		for _, e := range c.afterChat {
			if e.PluginID == pluginId {
				return
			}
		}
		c.afterChat = append(c.afterChat, HookEntry{PluginID: pluginId})
	}
}

// Unregister removes all hook entries for the given plugin from both slices.
func (c *HookChain) Unregister(pluginId string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.beforeChat = removePlugin(c.beforeChat, pluginId)
	c.afterChat = removePlugin(c.afterChat, pluginId)
}

func removePlugin(entries []HookEntry, pluginId string) []HookEntry {
	n := 0
	for _, e := range entries {
		if e.PluginID != pluginId {
			entries[n] = e
			n++
		}
	}
	return entries[:n]
}

// RunBeforeChat executes the before-chat hook pipeline. Each registered hook
// may transform the ChatContext. Individual hook failures are logged and skipped.
func (c *HookChain) RunBeforeChat(ctx *ChatContext) (*ChatContext, error) {
	c.mu.RLock()
	hooks := make([]HookEntry, len(c.beforeChat))
	copy(hooks, c.beforeChat)
	c.mu.RUnlock()

	if !c.host.IsRunning() || len(hooks) == 0 {
		return ctx, nil
	}

	for _, entry := range hooks {
		raw, err := c.host.RPC().Call("hook/onBeforeChat", hookCallParams{
			PluginID: entry.PluginID,
			Context:  ctx,
		})
		if err != nil {
			logger.Warmf("hook/onBeforeChat failed for plugin %s: %v", entry.PluginID, err)
			continue
		}

		var next ChatContext
		if err := json.Unmarshal(raw, &next); err != nil {
			logger.Warmf("hook/onBeforeChat unmarshal failed for plugin %s: %v", entry.PluginID, err)
			continue
		}
		ctx = &next
	}

	return ctx, nil
}

// RunAfterChat executes the after-chat hook pipeline. Each registered hook
// may transform the ChatContext. Individual hook failures are logged and skipped.
func (c *HookChain) RunAfterChat(ctx *ChatContext) (*ChatContext, error) {
	c.mu.RLock()
	hooks := make([]HookEntry, len(c.afterChat))
	copy(hooks, c.afterChat)
	c.mu.RUnlock()

	if !c.host.IsRunning() || len(hooks) == 0 {
		return ctx, nil
	}

	for _, entry := range hooks {
		raw, err := c.host.RPC().Call("hook/onAfterChat", hookCallParams{
			PluginID: entry.PluginID,
			Context:  ctx,
		})
		if err != nil {
			logger.Warmf("hook/onAfterChat failed for plugin %s: %v", entry.PluginID, err)
			continue
		}

		var next ChatContext
		if err := json.Unmarshal(raw, &next); err != nil {
			logger.Warmf("hook/onAfterChat unmarshal failed for plugin %s: %v", entry.PluginID, err)
			continue
		}
		ctx = &next
	}

	return ctx, nil
}
