package agent_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"

type GenerateTitleInput struct {
	SessionID    uint          `json:"session_id"`
	BaseURL      string        `json:"base_url"`
	ApiKey       string        `json:"api_key"`
	ModelName    string        `json:"model_name"`
	ProviderType provider.Type `json:"provider_type"`
}

type GenerateTitleOutput struct {
	Title string `json:"title"`
}
