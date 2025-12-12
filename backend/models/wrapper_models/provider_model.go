package wrapper_models

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type ProviderModel struct {
	ProviderType      data_models.ProviderType
	BaseUrl           string
	FileUploadBaseUrl *string
	ApiKey            string
	Model             string
	ModelId           uint
}
