package agent_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"

// SpawnTaskSessionInput describes the initial prompt for an automated task session.
type SpawnTaskSessionInput struct {
	Title        string        `json:"title"`
	SystemPrompt string        `json:"system_prompt"`
	UserMessage  string        `json:"user_message"`
	SkillName    string        `json:"skill_name"`
	BaseURL      string        `json:"base_url"`
	ApiKey       string        `json:"api_key"`
	ModelName    string        `json:"model_name"`
	ProviderType provider.Type `json:"provider_type"`
}

// SpawnTaskSessionOutput returns the created task session identity.
type SpawnTaskSessionOutput struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
}
