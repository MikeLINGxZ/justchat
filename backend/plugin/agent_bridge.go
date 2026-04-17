package plugin

import (
	"sync"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
)

// PluginAgentDef implements agents.IAgent for plugin-contributed agents.
type PluginAgentDef struct {
	pluginID     string
	agentID      string // full: "plugin:<pluginName>:<localId>"
	displayName  string
	description  string
	systemPrompt string
	toolIDs      []string
	role         agents.AgentRole
}

func (d *PluginAgentDef) Name() string           { return d.agentID }
func (d *PluginAgentDef) Desc() string           { return d.description }
func (d *PluginAgentDef) Prompt() string         { return d.systemPrompt }
func (d *PluginAgentDef) Type() agents.AgentType { return agents.AgentTypePlugin }

func (d *PluginAgentDef) Role() agents.AgentRole {
	if d.role == "" {
		return agents.AgentRoleWorker
	}
	return d.role
}

func (d *PluginAgentDef) PromptNames() []string                 { return []string{} }
func (d *PluginAgentDef) PromptMetas() []agents.AgentPromptMeta { return []agents.AgentPromptMeta{} }
func (d *PluginAgentDef) DefaultPrompts() map[string]string     { return map[string]string{} }

// AgentBridge manages plugin-contributed agents and bridges them to the global agent registry.
type AgentBridge struct {
	mu     sync.RWMutex
	agents map[string]*PluginAgentDef // full agentID → def
}

// NewAgentBridge creates a new AgentBridge.
func NewAgentBridge() *AgentBridge {
	return &AgentBridge{
		agents: make(map[string]*PluginAgentDef),
	}
}

// Register creates a PluginAgentDef and registers it with both the bridge and the global agent registry.
func (b *AgentBridge) Register(pluginID, pluginName string, contrib AgentContrib, systemPrompt string, toolIDs []string, role agents.AgentRole) {
	fullID := FullAgentID(pluginName, contrib.ID)
	def := &PluginAgentDef{
		pluginID:     pluginID,
		agentID:      fullID,
		displayName:  contrib.Name,
		description:  contrib.Description,
		systemPrompt: systemPrompt,
		toolIDs:      toolIDs,
		role:         role,
	}

	b.mu.Lock()
	b.agents[fullID] = def
	b.mu.Unlock()

	agents.RegisterAgent(def)
}

// UnregisterByPlugin removes all agents belonging to the given pluginID from both the bridge and the global registry.
func (b *AgentBridge) UnregisterByPlugin(pluginID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for id, def := range b.agents {
		if def.pluginID == pluginID {
			agents.UnregisterAgentByName(id)
			delete(b.agents, id)
		}
	}
}

// GetAgentsByPlugin returns all agents belonging to the given pluginID.
func (b *AgentBridge) GetAgentsByPlugin(pluginID string) []*PluginAgentDef {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var result []*PluginAgentDef
	for _, def := range b.agents {
		if def.pluginID == pluginID {
			result = append(result, def)
		}
	}
	return result
}
