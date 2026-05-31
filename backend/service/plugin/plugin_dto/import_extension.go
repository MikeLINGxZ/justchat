package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type ImportExtensionInput struct {
	Kind string `json:"kind"`
	Path string `json:"path"`
}

type ImportExtensionOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
