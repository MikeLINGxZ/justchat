package provider_dto

type SetDefaultInput struct {
	ProviderId int64  `json:"provider_id"` // 供应商ID
	ModelId    *int64 `json:"model_id"`    // 默认模型id
}

type SetDefaultOutput struct {
}
