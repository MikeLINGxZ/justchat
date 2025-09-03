package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type Service struct {
	ctx     context.Context
	storage *storage.Storage
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}
