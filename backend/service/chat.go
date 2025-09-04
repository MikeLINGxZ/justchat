package service

import (
	"github.com/cloudwego/eino/schema"
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

func (s *Service) Completions(chatUuid, model string, message schema.Message) (string, error) {

	//cancle := runtime.EventsOn(s.ctx, chatUuid, func(optionalData ...interface{}) {
	//
	//})
	return "", nil
}
