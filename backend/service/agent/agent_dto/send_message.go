package agent_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"

type SendMessageInput struct {
	SessionID        uint              `json:"session_id"`
	Content          string            `json:"content"`
	SystemPrompt     string            `json:"system_prompt,omitempty"`
	SkillName        string            `json:"skill_name,omitempty"`
	BaseURL          string            `json:"base_url"`
	ApiKey           string            `json:"api_key"`
	ModelName        string            `json:"model_name"`
	ProviderType     provider.Type     `json:"provider_type"`
	EnabledUserTools []string          `json:"enabled_user_tools"`
	Attachments      []AttachmentInput `json:"attachments,omitempty"`
}

type SendMessageOutput struct {
}
