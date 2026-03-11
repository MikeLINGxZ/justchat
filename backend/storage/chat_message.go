package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateMessage 创建消息
func (s *Storage) CreateMessage(ctx context.Context, chatUuid string, message data_models.Message) (uint, error) {
	err := s.sqliteDB.Create(&message).Error
	return message.ID, err
}

func (s *Storage) SaveOrUpdateMessage(ctx context.Context, message data_models.Message) error {
	// 先查询现有记录
	var existingMessage data_models.Message
	err := s.sqliteDB.Where("message_uuid = ?", message.MessageUuid).First(&existingMessage).Error
	if err != nil {
		// 如果记录不存在，创建新消息
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.sqliteDB.Create(&message).Error
		}
		// 其他错误直接返回
		return err
	}
	return s.sqliteDB.Where("message_uuid = ?", message.MessageUuid).Updates(&message).Error
}

// GetMessage 获取 chat 消息
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
