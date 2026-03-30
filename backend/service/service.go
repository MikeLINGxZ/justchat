package service

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const (
	WindowNameHome       = "window_home"
	WindowNameOnboarding = "window_onboarding"
	WindowNameSettings   = "window_settings"
)

type Service struct {
	storage *storage.Storage
	app     *application.App
	prompts prompts.PromptSet
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) localizedPromptSet() prompts.PromptSet {
	locale := i18n.CurrentLocale()
	if locale == "" {
		locale = string(data_models.AppLanguageZhCN)
	}
	return prompts.WithResponseLanguage(s.prompts, locale)
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

	istorage, err := storage.NewStorage()
	if err != nil {
		return err
	}

	s.storage = istorage
	s.app = application.Get()
	if prefs, prefsErr := s.loadAppPreferences(ctx); prefsErr == nil {
		i18n.SetCurrentLocale(string(prefs.Language))
	}
	if err := s.reloadPromptSet(); err != nil {
		logger.Warm("load prompt set fallback:", err)
	}

	if err := s.syncCustomMCPTools(ctx); err != nil {
		return err
	}

	if err := s.recoverStaleRunningTasks(ctx); err != nil {
		return err
	}

	return nil
}
