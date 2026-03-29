package service

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const (
	WindowNameHome     = "window_home"
	WindowNameSettings = "window_settings"
)

type Service struct {
	storage *storage.Storage
	app     *application.App
	prompts prompts.PromptSet
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

	istorage, err := storage.NewStorage()
	if err != nil {
		return err
	}

	s.storage = istorage
	s.app = application.Get()
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
