package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type GetExtensionDetailInput struct {
	ID string `json:"id"`
}

type GetExtensionDetailOutput struct {
	Extension  data_models.ExtensionItem `json:"extension"`
	ConfigText string                    `json:"config_text"`
}
