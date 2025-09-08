package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gorm.io/gorm"
)

// GetChats 获取对话
func (s *Storage) GetChats(ctx context.Context, offset, limit int, keyword *string) ([]view_models.Chat, int, error) {
	var chats []view_models.Chat
	var res []data_models.Chat
	var count int64
	if keyword == nil || strings.TrimSpace(*keyword) == "" {
		err := s.sqliteDB.Model(&data_models.Chat{}).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		err = s.sqliteDB.Model(&data_models.Chat{}).Order("updated_at DESC").Offset(offset).Limit(limit).Find(&res).Error
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
	err := s.sqliteDB.Model(&data_models.Chat{}).Where("title LIKE ?", "%"+keywordStr+"%").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = s.sqliteDB.Model(&data_models.Chat{}).Order("updated_at DESC").Where("title LIKE ?", "%"+keywordStr+"%").Offset(offset).Limit(limit).Find(&res).Error
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

// CreateMessage 创建消息
func (s *Storage) CreateMessage(ctx context.Context, chatUuid string, message data_models.Message) error {
	return s.sqliteDB.Create(&message).Error
}

// SaveOrUpdateDeltaMessage 创建或更新消息
func (s *Storage) SaveOrUpdateDeltaMessage(ctx context.Context, deltaMessage data_models.Message) error {
	// 如果消息ID为0，说明是新消息，直接创建
	if deltaMessage.Uuid == "" {
		return s.sqliteDB.Create(&deltaMessage).Error
	}

	// 先查询现有记录
	var existingMessage data_models.Message
	err := s.sqliteDB.Where("uuid = ?", deltaMessage.Uuid).First(&existingMessage).Error
	if err != nil {
		// 如果记录不存在，创建新消息
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.sqliteDB.Create(&deltaMessage).Error
		}
		// 其他错误直接返回
		return err
	}

	// 检查现有消息的Message字段
	schemaMsg := existingMessage.Message
	if schemaMsg == nil {
		return errors.New("existing message schema is not defined")
	}

	// 增量更新内容
	if deltaMessage.Message != nil {
		if deltaMessage.Message.Content != "" {
			schemaMsg.Content += deltaMessage.Message.Content
		}
		if deltaMessage.Message.ReasoningContent != "" {
			schemaMsg.ReasoningContent += deltaMessage.Message.ReasoningContent
		}
		if deltaMessage.Message.ResponseMeta != nil {
			schemaMsg.ResponseMeta = deltaMessage.Message.ResponseMeta
		}
	}

	// 更新现有消息的Message字段
	existingMessage.Message = schemaMsg

	// 保存更新后的消息
	return s.sqliteDB.Where("uuid = ?", deltaMessage.Uuid).Updates(&existingMessage).Error
}
