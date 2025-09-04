package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func (s *Service) ChatList(ctx context.Context, offset, limit int, keyword *string) ([]view_models.Chat, error) {
	chats, err := s.storage.GetChats(ctx, offset, limit, keyword)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
