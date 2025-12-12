package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

// GetSupportProviders 获取支持的供应商列表
func (s *Service) GetSupportProviders() ([]view_models.SupportProvider, error) {
	return []view_models.SupportProvider{
		{
			ProviderType:      data_models.ProviderTypeDeepseek,
			Icon:              "/providers/deepseek_icon.png",
			Name:              "深度求索",
			BaseUrl:           "https://api.deepseek.com/v1",
			FileUploadBaseUrl: nil,
			Description:       "成立于2023年，专注于研究世界领先的通用人工智能底层模型与技术，挑战人工智能前沿性难题。",
		}, {
			ProviderType:      data_models.ProviderTypeAliyuns,
			Icon:              "/providers/qwen_icon.png",
			Name:              "阿里云百炼",
			BaseUrl:           "https://dashscope.aliyuncs.com/compatible-mode/v1",
			FileUploadBaseUrl: utils.StringPointer("https://dashscope.aliyuncs.com/api/v1/uploads"),
			Description:       "一键部署大模型 — 阿里云百炼-极强推理能力极高性价比.百炼支持多种模态的大模型调用服务。",
		}, {
			ProviderType:      data_models.ProviderTypeOpenrouter,
			Icon:              "/providers/openrouter_icon.png",
			Name:              "OpenRouter",
			BaseUrl:           "https://openrouter.ai/api/v1",
			FileUploadBaseUrl: nil,
			Description:       "模型集合供应商",
		}, {
			ProviderType:      data_models.ProviderTypeOther,
			Icon:              "/providers/openai_icon.png",
			Name:              "Openai 标准接口",
			BaseUrl:           "",
			FileUploadBaseUrl: utils.StringPointer(""),
			Description:       "",
		},
	}, nil
}
