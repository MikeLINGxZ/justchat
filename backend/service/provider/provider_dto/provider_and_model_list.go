package provider_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type ProviderAndModelListInput struct {
}

type ProviderAndModelListOutput struct {
	ProviderModels []view_model.ProviderModel `json:"provider_models"`
}
