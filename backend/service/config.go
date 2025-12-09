package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

// GetSupportProviders 获取支持的供应商列表
func (s *Service) GetSupportProviders() ([]view_models.SupportProvider, error) {
	return []view_models.SupportProvider{
		{
			Name:              "深度求索",
			BaseUrl:           "https://api.deepseek.com/v1",
			FileUploadBaseUrl: nil,
		}, {
			Name:              "阿里云百炼",
			BaseUrl:           "https://dashscope.aliyuncs.com/compatible-mode/v1",
			FileUploadBaseUrl: utils.StringPointer("https://dashscope.aliyuncs.com/api/v1/uploads"),
		}, {
			Name:              "OpenRouter",
			BaseUrl:           "https://openrouter.ai/api/v1",
			FileUploadBaseUrl: nil,
		}, {
			Name:              "Openai 标准接口",
			BaseUrl:           "",
			FileUploadBaseUrl: utils.StringPointer(""),
		},
	}, nil
}
