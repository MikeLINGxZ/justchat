package config_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type SupportedProviderListInput struct {
}

type SupportedProviderListOutput struct {
	SupportedProviders []view_model.SupportedProvider `json:"supported_providers"`
}
