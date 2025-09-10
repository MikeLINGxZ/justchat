package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type Service struct {
	ctx     context.Context
	storage *storage.Storage
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Startup(ctx context.Context) error {
	s.ctx = ctx
	storage, err := storage.NewStorage()
	if err != nil {
		return err
	}
	s.storage = storage
	return nil
}
