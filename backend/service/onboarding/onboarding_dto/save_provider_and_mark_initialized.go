package onboarding_dto

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

// SaveProviderAndMarkInitializedInput mirrors provider_dto.CreateProviderInput plus onboarding intent.
type SaveProviderAndMarkInitializedInput struct {
	ProviderName string             `json:"provider_name"`
	ProviderType provider.Type      `json:"provider_type"`
	BaseUrl      string             `json:"base_url"`
	ApiKey       string             `json:"api_key"`
	Enable       bool               `json:"enable"`
	DefaultModel *string            `json:"default_model"`
	Models       []view_model.Model `json:"models"`
}

// SaveProviderAndMarkInitializedOutput is an empty acknowledgement.
type SaveProviderAndMarkInitializedOutput struct{}
