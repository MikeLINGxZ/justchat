package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

// UpdateCliManifestInput carries the edited manifest JSON text for a specific CLI plugin.
type UpdateCliManifestInput struct {
	ID           string `json:"id"`
	ManifestText string `json:"manifest_text"`
}

// UpdateCliManifestOutput returns the updated ExtensionItem after manifest validation + save.
type UpdateCliManifestOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
