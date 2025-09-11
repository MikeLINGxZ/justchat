package storage

import (
	"context"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

// GetChats 获取对话
func (s *Storage) GetChats(ctx context.Context, offset, limit int, keyword *string, isCollection bool) ([]view_models.Chat, int, error) {
	var chats []view_models.Chat
	var res []data_models.Chat
	var count int64

	queryBase := s.sqliteDB.Model(&data_models.Chat{})
	if isCollection {
		queryBase = queryBase.Where("is_collection = ?", 1)
	}

	if keyword == nil || strings.TrimSpace(*keyword) == "" {
		err := queryBase.Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		err = queryBase.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&res).Error
		if err != nil {
			return nil, 0, err
		}
		for _, item := range res {
			chats = append(chats, view_models.Chat{
				Chat: item,
			})
		}
		return chats, int(count), nil
	}

	// 使用关键字搜索包含匹配消息的聊天与聊天标题
	keywordStr := strings.TrimSpace(*keyword)
	err := queryBase.Where("title LIKE ?", "%"+keywordStr+"%").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = queryBase.Where("title LIKE ?", "%"+keywordStr+"%").Offset(offset).Limit(limit).Find(&res).Error
	if err != nil {
		return nil, 0, err
	}
	for _, item := range res {
		chats = append(chats, view_models.Chat{
			Chat: item,
		})
	}

	return chats, int(count), nil
}

// CreateChat 创建对话
func (s *Storage) CreateChat(ctx context.Context, chatUuid, title string, modelId uint) error {
	now := time.Now()
	chat := &data_models.Chat{
		OrmModel: data_models.OrmModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Uuid:    chatUuid,
		ModelID: modelId,
		Title:   title,
		Prompt:  "",
	}

	err := s.sqliteDB.Create(chat).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteChat 删除对话
func (s *Storage) DeleteChat(ctx context.Context, chatUuid string) error {
	err := s.sqliteDB.Where("uuid = ?", chatUuid).Delete(&data_models.Chat{}).Error
	if err != nil {
		return err
	}
	return nil
}

// RenameChat 重命名对话
func (s *Storage) RenameChat(ctx context.Context, chatUuid, title string) error {
	err := s.sqliteDB.Model(&data_models.Chat{}).Where("uuid = ?", chatUuid).Update("title", title).Error
	if err != nil {
		return err
	}
	return nil
}

// CollectionChat ...
func (s *Storage) CollectionChat(ctx context.Context, chatUuid string, isCollection bool) error {
	err := s.sqliteDB.Model(&data_models.Chat{}).Where("uuid = ?", chatUuid).Update("is_collection", isCollection).Error
	if err != nil {
		return err
	}
	return nil
}
