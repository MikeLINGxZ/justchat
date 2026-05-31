package provider_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type RequestProviderModelListInput struct {
	BaseUrl string `json:"base_url"` // 供应商基础URL
	ApiKey  string `json:"api_key"`  // 供应商API密钥
}

type RequestProviderModelListOutput struct {
	Models []view_model.Model `json:"models"`
}
