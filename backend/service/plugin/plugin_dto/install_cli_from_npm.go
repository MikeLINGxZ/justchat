package plugin_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

// InstallCliFromNpmInput is the request to install a CLI plugin from a published npm package.
type InstallCliFromNpmInput struct {
	NpmPackage string `json:"npm_package"`
	Name       string `json:"name"`
}

// InstallCliFromNpmOutput returns the persisted ExtensionItem after a successful install.
type InstallCliFromNpmOutput struct {
	Extension data_models.ExtensionItem `json:"extension"`
}
