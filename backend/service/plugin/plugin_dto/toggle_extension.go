package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type ToggleExtensionInput struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type ToggleExtensionOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
