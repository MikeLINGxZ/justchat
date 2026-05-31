package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type ReloadExtensionInput struct {
	ID string `json:"id"`
}

type ReloadExtensionOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
