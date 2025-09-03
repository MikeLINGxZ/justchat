package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateChat 创建新的聊天会话
func (s *Storage) CreateChat(chat *data_models.Chat) error {
	return s.sqliteDB.Create(chat).Error
}

// GetChat 根据ID获取聊天
func (s *Storage) GetChat(id uint) (*data_models.Chat, error) {
	var chat data_models.Chat
	err := s.sqliteDB.First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// GetChatByUUID 根据UUID获取聊天
func (s *Storage) GetChatByUUID(uuid string) (*data_models.Chat, error) {
	var chat data_models.Chat
	err := s.sqliteDB.Where("chat_uuid = ?", uuid).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// UpdateChat 更新聊天信息
func (s *Storage) UpdateChat(chat *data_models.Chat) error {
	return s.sqliteDB.Save(chat).Error
}

// DeleteChat 删除聊天（软删除）
func (s *Storage) DeleteChat(id uint) error {
	return s.sqliteDB.Delete(&data_models.Chat{}, id).Error
}

// ListChats 获取聊天列表
func (s *Storage) ListChats(limit, offset int) ([]data_models.Chat, error) {
	var chats []data_models.Chat
	err := s.sqliteDB.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&chats).Error
	return chats, err
}

// CreateMessage 创建新消息
func (s *Storage) CreateMessage(message *data_models.Message) error {
	return s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		// 创建消息
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		// 更新聊天统计信息
		return tx.Model(&data_models.Chat{}).Where("id = ?", message.ChatID).Updates(map[string]interface{}{
			"message_count": gorm.Expr("message_count + 1"),
			"last_activity": time.Now(),
		}).Error
	})
}

// CreateMessageFromSchema 从 schema.Message 创建消息
func (s *Storage) CreateMessageFromSchema(chatID uint, msg *schema.Message) (*data_models.Message, error) {
	message := &data_models.Message{
		ChatID: chatID,
	}

	if err := message.FromSchemaMessage(msg); err != nil {
		return nil, fmt.Errorf("failed to convert from schema.Message: %w", err)
	}

	if err := s.CreateMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessage 根据ID获取消息
func (s *Storage) GetMessage(id uint) (*data_models.Message, error) {
	var message data_models.Message
	err := s.sqliteDB.First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetMessageByUUID 根据UUID获取消息
func (s *Storage) GetMessageByUUID(uuid string) (*data_models.Message, error) {
	var message data_models.Message
	err := s.sqliteDB.Where("message_uuid = ?", uuid).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateMessage 更新消息
func (s *Storage) UpdateMessage(message *data_models.Message) error {
	return s.sqliteDB.Save(message).Error
}

// DeleteMessage 删除消息（软删除）
func (s *Storage) DeleteMessage(id uint) error {
	return s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		// 获取消息信息
		var message data_models.Message
		if err := tx.First(&message, id).Error; err != nil {
			return err
		}

		// 删除消息
		if err := tx.Delete(&message).Error; err != nil {
			return err
		}

		// 更新聊天统计信息
		return tx.Model(&data_models.Chat{}).Where("id = ?", message.ChatID).Updates(map[string]interface{}{
			"message_count": gorm.Expr("message_count - 1"),
			"last_activity": time.Now(),
		}).Error
	})
}

// GetChatMessages 获取聊天的所有消息
func (s *Storage) GetChatMessages(chatID uint, limit, offset int) ([]data_models.Message, error) {
	var messages []data_models.Message
	err := s.sqliteDB.Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// GetChatMessagesAsSchema 获取聊天的所有消息并转换为 schema.Message
func (s *Storage) GetChatMessagesAsSchema(chatID uint, limit, offset int) ([]*schema.Message, error) {
	messages, err := s.GetChatMessages(chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	var schemaMessages []*schema.Message
	for _, msg := range messages {
		schemaMsg, err := msg.ToSchemaMessage()
		if err != nil {
			return nil, fmt.Errorf("failed to convert message %d to schema: %w", msg.ID, err)
		}
		schemaMessages = append(schemaMessages, schemaMsg)
	}

	return schemaMessages, nil
}

// SearchMessagesResult 搜索结果结构
type SearchMessagesResult struct {
	Messages []data_models.Message `json:"messages"`
	Total    int64                 `json:"total"`
}

// SearchMessages 全文搜索消息
func (s *Storage) SearchMessages(query string, chatID *uint, limit, offset int) (*SearchMessagesResult, error) {
	if strings.TrimSpace(query) == "" {
		return &SearchMessagesResult{Messages: []data_models.Message{}, Total: 0}, nil
	}

	// 构建搜索查询
	searchQuery := fmt.Sprintf(`"%s"`, strings.ReplaceAll(query, `"`, `""`))

	// 基础查询条件
	baseWhere := "messages_fts MATCH ?"
	args := []interface{}{searchQuery}

	// 如果指定了 chatID，添加过滤条件
	if chatID != nil {
		baseWhere += " AND chat_id = ?"
		args = append(args, *chatID)
	}

	// 执行搜索查询
	var messages []data_models.Message
	err := s.sqliteDB.Raw(`
		SELECT m.* FROM messages m
		INNER JOIN messages_fts fts ON m.id = fts.id
		WHERE `+baseWhere+`
		ORDER BY fts.rank, m.created_at DESC
		LIMIT ? OFFSET ?
	`, append(args, limit, offset)...).Scan(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	// 获取总数
	var total int64
	err = s.sqliteDB.Raw(`
		SELECT COUNT(*) FROM messages_fts
		WHERE `+baseWhere,
		args...).Scan(&total).Error

	if err != nil {
		return nil, fmt.Errorf("failed to count search results: %w", err)
	}

	return &SearchMessagesResult{
		Messages: messages,
		Total:    total,
	}, nil
}

// SearchMessagesInChat 在指定聊天中搜索消息
func (s *Storage) SearchMessagesInChat(chatID uint, query string, limit, offset int) (*SearchMessagesResult, error) {
	return s.SearchMessages(query, &chatID, limit, offset)
}

// SearchMessagesGlobal 全局搜索消息
func (s *Storage) SearchMessagesGlobal(query string, limit, offset int) (*SearchMessagesResult, error) {
	return s.SearchMessages(query, nil, limit, offset)
}

// GetMessageStatistics 获取消息统计信息
func (s *Storage) GetMessageStatistics(chatID *uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := s.sqliteDB.Model(&data_models.Message{})
	if chatID != nil {
		query = query.Where("chat_id = ?", *chatID)
	}

	// 总消息数
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_messages"] = totalCount

	// 按角色统计
	var roleStats []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	err := query.Select("role, COUNT(*) as count").Group("role").Find(&roleStats).Error
	if err != nil {
		return nil, err
	}
	stats["by_role"] = roleStats

	// 总 token 数
	var totalTokens int64
	err = query.Select("SUM(token_count)").Scan(&totalTokens).Error
	if err != nil {
		return nil, err
	}
	stats["total_tokens"] = totalTokens

	return stats, nil
}
