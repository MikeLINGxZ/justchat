package view_models

type PluginContribItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PluginSummary struct {
	ID          string              `json:"id"`
	DisplayName string              `json:"display_name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	Enabled     bool                `json:"enabled"`
	State       string              `json:"state"`
	ToolCount   int                 `json:"tool_count"`
	ViewCount   int                 `json:"view_count"`
	HookCount   int                 `json:"hook_count"`
	AgentCount  int                 `json:"agent_count"`
	Tools       []PluginContribItem `json:"tools"`
	Agents      []PluginContribItem `json:"agents"`
	Hooks       []string            `json:"hooks"`
}
