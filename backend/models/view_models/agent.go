package view_models

type AgentSummary struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	AgentType   string   `json:"agent_type"`
	AgentRole   string   `json:"agent_role"`
	PromptNames []string `json:"prompt_names"`
	IsDeletable bool     `json:"is_deletable"`
	Tools       []string `json:"tools"`
	Skills      []string `json:"skills"`
}

type AgentDetail struct {
	AgentSummary
	Prompts []AgentPrompt `json:"prompts"`
}

type CustomAgentInput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Prompt      string   `json:"prompt"`
	Tools       []string `json:"tools"`
	Skills      []string `json:"skills"`
}

type AgentPrompt struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	IsSystem    bool   `json:"is_system"`
}
