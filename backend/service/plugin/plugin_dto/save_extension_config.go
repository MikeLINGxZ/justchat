package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type SaveExtensionConfigInput struct {
	ID         string `json:"id"`
	ConfigText string `json:"config_text"`
}

type SaveExtensionConfigOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
