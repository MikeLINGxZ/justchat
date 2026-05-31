package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type ListExtensionsInput struct{}

type ListExtensionsOutput struct {
	Extensions []data_models.ExtensionItem `json:"extensions"`
}
