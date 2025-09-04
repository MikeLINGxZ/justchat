package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func (s *Service) ChatList(offset, limit int, keyword *string) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(s.ctx, offset, limit, keyword)
	if err != nil {
		return nil, err
	}

	res := view_models.ChatList{
		Lists: chats,
		Total: total,
	}

	return &res, nil
}
