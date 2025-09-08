package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

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

func (s *Storage) GetMessage(ctx context.Context, chatUuid string, offset, limit int) ([]data_models.Message, int, error) {
	var messages []data_models.Message
	err := s.sqliteDB.Model(&data_models.Message{}).Where("chat_uuid = ?", chatUuid).Offset(offset).Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = s.sqliteDB.Model(&data_models.Message{}).Where("chat_uuid = ?", chatUuid).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return messages, int(total), nil
}
