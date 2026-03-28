package storage

import (
	"context"
	"errors"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateMessage 创建消息
func (s *Storage) CreateMessage(ctx context.Context, chatUuid string, message data_models.Message) (uint, error) {
	err := s.sqliteDB.Create(&message).Error
	if err != nil {
		return message.ID, err
	}
	return message.ID, s.touchChat(chatUuid)
}

func (s *Storage) SaveOrUpdateMessage(ctx context.Context, message data_models.Message) error {
	// 先查询现有记录
	var existingMessage data_models.Message
	err := s.sqliteDB.Where("message_uuid = ?", message.MessageUuid).First(&existingMessage).Error
	if err != nil {
		// 如果记录不存在，创建新消息
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = s.sqliteDB.Create(&message).Error; err != nil {
				return err
			}
			return s.touchChat(message.ChatUuid)
		}
		// 其他错误直接返回
		return err
	}
	if err = s.sqliteDB.Where("message_uuid = ?", message.MessageUuid).Updates(&message).Error; err != nil {
		return err
	}
	return s.touchChat(message.ChatUuid)
}

// GetMessage 获取 chat 消息
func (s *Storage) GetMessage(ctx context.Context, chatUuid string, offset, limit int) ([]data_models.Message, int, error) {
	var messages []data_models.Message
	err := s.sqliteDB.Model(&data_models.Message{}).Where("chat_uuid = ?", chatUuid).Order("created_at asc").Offset(offset).Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = s.sqliteDB.Model(&data_models.Message{}).Where("chat_uuid = ?", chatUuid).Order("created_at asc").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return messages, int(total), nil
}

func (s *Storage) touchChat(chatUuid string) error {
	if chatUuid == "" {
		return nil
	}
	return s.sqliteDB.Model(&data_models.Chat{}).
		Where("uuid = ?", chatUuid).
		Update("updated_at", time.Now()).
		Error
}
