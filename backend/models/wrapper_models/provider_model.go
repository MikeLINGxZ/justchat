package wrapper_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/llm"

type ProviderModel struct {
	ProviderType llm.ProviderType
	BaseUrl      string
	ApiKey       string
	Model        string
	ModelId      uint
}
