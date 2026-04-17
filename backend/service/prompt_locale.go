package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
)

func (s *Service) localizedPromptSet() prompts.PromptSet {
	locale := i18n.CurrentLocale()
	if locale == "" {
		locale = string(data_models.AppLanguageZhCN)
	}
	return prompts.WithResponseLanguage(s.prompts, locale)
}
