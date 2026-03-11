package service

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const (
	WindowNameHome     = "window_home"
	WindowNameSettings = "window_settings"
)

type Service struct {
	storage           *storage.Storage
	app               *application.App
	completionsStopCh map[string]chan struct{} // 用于用户手动结束的ch
}

func NewService() *Service {
	return &Service{
		completionsStopCh: make(map[string]chan struct{}),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

	istorage, err := storage.NewStorage()
	if err != nil {
		return err
	}

	s.storage = istorage
	s.app = application.Get()

	return nil
}
