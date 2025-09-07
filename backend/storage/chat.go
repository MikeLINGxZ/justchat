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

	// 收集所有匹配的聊天ID
	matchedChatIDs := make(map[uint]bool)

	// 搜索聊天标题包含关键字的聊天ID
	var titleMatchChatIDs []uint
	err := s.sqliteDB.Model(&data_models.Chat{}).
		Where("title LIKE ?", "%"+keywordStr+"%").
		Pluck("id", &titleMatchChatIDs).Error
	if err != nil {
		return chats, 0, nil
	}
	for _, id := range titleMatchChatIDs {
		matchedChatIDs[id] = true
	}

	// 搜索消息内容包含关键字的聊天ID
	var messageMatchChatIDs []uint
	err = s.sqliteDB.Model(&data_models.Message{}).
		Where("searchable_content LIKE ? OR searchable_reasoning_content LIKE ?", "%"+keywordStr+"%", "%"+keywordStr+"%").
		Distinct("chat_id").
		Pluck("chat_id", &messageMatchChatIDs).Error
	if err != nil {
		return chats, 0, nil
	}
	for _, id := range messageMatchChatIDs {
		matchedChatIDs[id] = true
	}

	// 将匹配的聊天ID转换为切片
	var allMatchedIDs []uint
	for id := range matchedChatIDs {
		allMatchedIDs = append(allMatchedIDs, id)
	}
	count = int64(len(allMatchedIDs))

	// 如果没有匹配的聊天，直接返回空结果
	if len(allMatchedIDs) == 0 {
		return chats, 0, nil
	}

	// 对匹配的聊天应用分页并获取详细信息
	err = s.sqliteDB.Model(&data_models.Chat{}).
		Where("id IN ?", allMatchedIDs).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&res).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为view_models.Chat并填充匹配的消息内容
	for _, item := range res {
		chat := view_models.Chat{
			Chat: item,
		}

		// 获取该聊天中匹配关键字的消息
		var matchedMessages []data_models.Message
		err = s.sqliteDB.Model(&data_models.Message{}).
			Where("chat_id = ? AND (searchable_content LIKE ? OR searchable_reasoning_content LIKE ?)",
				item.ID, "%"+keywordStr+"%", "%"+keywordStr+"%").
			Find(&matchedMessages).Error
		if err != nil {
			return nil, 0, err
		}

		// 填充匹配的消息内容
		for _, msg := range matchedMessages {
			if msg.Message != nil {
				// 检查普通内容是否匹配
				if strings.Contains(strings.ToLower(msg.SearchableContent), strings.ToLower(keywordStr)) {
					chat.Content = append(chat.Content, view_models.MatchMessage{
						Role:    string(msg.Message.Role),
						Content: msg.Message.Content,
					})
				}
				// 检查推理内容是否匹配
				if strings.Contains(strings.ToLower(msg.SearchableReasoningContent), strings.ToLower(keywordStr)) {
					chat.ReasoningContent = append(chat.ReasoningContent, view_models.MatchMessage{
						Role:    string(msg.Message.Role),
						Content: msg.Message.ReasoningContent,
					})
				}
			}
		}

		chats = append(chats, chat)
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

// SaveOrUpdateDeltaMessage 创建或更新消息
func (s *Storage) SaveOrUpdateDeltaMessage(ctx context.Context, deltaMessage data_models.Message) error {
	// 如果消息ID为0，说明是新消息，直接创建
	if deltaMessage.ID == 0 {
		return s.sqliteDB.Create(deltaMessage).Error
	}

	// 先查询现有记录
	var existingMessages []data_models.Message
	var existingMessage data_models.Message
	err := s.sqliteDB.Model(&data_models.Message{}).Where("id = ?", deltaMessage.ID).Find(&existingMessages).Error
	if err != nil {
		// 如果记录不存在，创建新消息
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = s.sqliteDB.Create(deltaMessage).Error
			if err != nil {
				return err
			}
			existingMessage = deltaMessage
		}
		return err
	}
	if len(existingMessages) > 0 {
		existingMessage = existingMessages[0]
	}

	schemaMsg := existingMessage.Message
	if schemaMsg == nil {
		return errors.New("msg is not defined")
	}

	if deltaMessage.Message != nil && deltaMessage.Message.Content != "" {
		schemaMsg.Content += deltaMessage.Message.Content
	}
	if deltaMessage.Message != nil && deltaMessage.Message.ReasoningContent != "" {
		schemaMsg.ReasoningContent += deltaMessage.Message.ReasoningContent
	}
	if deltaMessage.Message.ResponseMeta != nil {
		schemaMsg.ResponseMeta = deltaMessage.Message.ResponseMeta
	}
	existingMessage.Message = schemaMsg

	return s.sqliteDB.Model(&data_models.Message{}).Where("uuid = ?", deltaMessage.Uuid).Save(&existingMessage).Error
}
