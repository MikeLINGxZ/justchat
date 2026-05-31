package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

// GenerateCliManifestInput identifies which CLI plugin to regenerate the manifest for.
type GenerateCliManifestInput struct {
	ID string `json:"id"`
}

// GenerateCliManifestOutput returns the refreshed extension after a successful generation.
type GenerateCliManifestOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
