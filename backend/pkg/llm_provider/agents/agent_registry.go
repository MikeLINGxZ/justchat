package agents

import "sync"

var (
	mu       sync.RWMutex
	registry []IAgent
)

// RegisterAgent 注册一个 Agent 到全局注册表（通常在 init() 中调用）。
func RegisterAgent(a IAgent) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, a)
}

// AllAgents 返回所有已注册的 Agent。
func AllAgents() []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]IAgent, len(registry))
	copy(result, registry)
	return result
}

// FindAgent 按名称查找 Agent。
func FindAgent(name string) (IAgent, bool) {
	mu.RLock()
	defer mu.RUnlock()
	for _, a := range registry {
		if a.Name() == name {
			return a, true
		}
	}
	return nil, false
}

// AgentOwnedPromptNames 返回所有 Agent 拥有的提示词名称集合。
func AgentOwnedPromptNames() map[string]struct{} {
	mu.RLock()
	defer mu.RUnlock()
	owned := make(map[string]struct{})
	for _, a := range registry {
		for _, name := range a.PromptNames() {
			owned[name] = struct{}{}
		}
	}
	return owned
}

// UnregisterAgentsByType removes all agents of the given type from the registry.
func UnregisterAgentsByType(agentType AgentType) {
	mu.Lock()
	defer mu.Unlock()
	filtered := make([]IAgent, 0, len(registry))
	for _, a := range registry {
		if a.Type() != agentType {
			filtered = append(filtered, a)
		}
	}
	registry = filtered
}

// AgentsByType 按类型过滤 Agent。
func AgentsByType(agentType AgentType) []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	var result []IAgent
	for _, a := range registry {
		if a.Type() == agentType {
			result = append(result, a)
		}
	}
	return result
}

// AgentsByRole 按角色过滤 Agent。
func AgentsByRole(role AgentRole) []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	var result []IAgent
	for _, a := range registry {
		if a.Role() == role {
			result = append(result, a)
		}
	}
	return result
}
