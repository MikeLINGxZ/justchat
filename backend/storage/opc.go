package storage

import (
	"context"
	"errors"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// ==================== OPCPerson ====================

func (s *Storage) CreateOPCPerson(ctx context.Context, person *data_models.OPCPerson) error {
	return s.sqliteDB.WithContext(ctx).Create(person).Error
}

func (s *Storage) GetOPCPersons(ctx context.Context) ([]data_models.OPCPerson, error) {
	var persons []data_models.OPCPerson
	err := s.sqliteDB.WithContext(ctx).
		Order("is_pinned DESC, updated_at DESC").
		Find(&persons).Error
	return persons, err
}

func (s *Storage) GetOPCPerson(ctx context.Context, uuid string) (*data_models.OPCPerson, error) {
	var person data_models.OPCPerson
	err := s.sqliteDB.WithContext(ctx).Where("uuid = ?", uuid).First(&person).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &person, err
}

func (s *Storage) GetOPCPersonByAgentID(ctx context.Context, agentID string) (*data_models.OPCPerson, error) {
	var person data_models.OPCPerson
	err := s.sqliteDB.WithContext(ctx).Where("agent_id = ?", agentID).First(&person).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &person, err
}

func (s *Storage) UpdateOPCPerson(ctx context.Context, person *data_models.OPCPerson) error {
	return s.sqliteDB.WithContext(ctx).Where("uuid = ?", person.Uuid).Updates(person).Error
}

func (s *Storage) DeleteOPCPerson(ctx context.Context, uuid string) error {
	return s.sqliteDB.WithContext(ctx).Where("uuid = ?", uuid).Delete(&data_models.OPCPerson{}).Error
}

func (s *Storage) TogglePinOPCPerson(ctx context.Context, uuid string, pinned bool) error {
	return s.sqliteDB.WithContext(ctx).Model(&data_models.OPCPerson{}).
		Where("uuid = ?", uuid).Update("is_pinned", pinned).Error
}

func (s *Storage) SearchOPCPersons(ctx context.Context, keyword string) ([]data_models.OPCPerson, error) {
	var persons []data_models.OPCPerson
	err := s.sqliteDB.WithContext(ctx).
		Where("name LIKE ?", "%"+keyword+"%").
		Order("is_pinned DESC, updated_at DESC").
		Find(&persons).Error
	return persons, err
}

// ==================== OPCGroup ====================

func (s *Storage) CreateOPCGroup(ctx context.Context, group *data_models.OPCGroup) error {
	return s.sqliteDB.WithContext(ctx).Create(group).Error
}

func (s *Storage) GetOPCGroups(ctx context.Context) ([]data_models.OPCGroup, error) {
	var groups []data_models.OPCGroup
	err := s.sqliteDB.WithContext(ctx).
		Order("is_pinned DESC, updated_at DESC").
		Find(&groups).Error
	return groups, err
}

func (s *Storage) GetOPCGroup(ctx context.Context, uuid string) (*data_models.OPCGroup, error) {
	var group data_models.OPCGroup
	err := s.sqliteDB.WithContext(ctx).Where("uuid = ?", uuid).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &group, err
}

func (s *Storage) UpdateOPCGroup(ctx context.Context, group *data_models.OPCGroup) error {
	return s.sqliteDB.WithContext(ctx).Where("uuid = ?", group.Uuid).Updates(group).Error
}

func (s *Storage) DeleteOPCGroup(ctx context.Context, uuid string) error {
	return s.sqliteDB.WithContext(ctx).Where("uuid = ?", uuid).Delete(&data_models.OPCGroup{}).Error
}

func (s *Storage) TogglePinOPCGroup(ctx context.Context, uuid string, pinned bool) error {
	return s.sqliteDB.WithContext(ctx).Model(&data_models.OPCGroup{}).
		Where("uuid = ?", uuid).Update("is_pinned", pinned).Error
}

func (s *Storage) SearchOPCGroups(ctx context.Context, keyword string) ([]data_models.OPCGroup, error) {
	var groups []data_models.OPCGroup
	err := s.sqliteDB.WithContext(ctx).
		Where("name LIKE ?", "%"+keyword+"%").
		Order("is_pinned DESC, updated_at DESC").
		Find(&groups).Error
	return groups, err
}

// ==================== OPCGroupMember ====================

func (s *Storage) CreateOPCGroupMembers(ctx context.Context, members []data_models.OPCGroupMember) error {
	if len(members) == 0 {
		return nil
	}
	return s.sqliteDB.WithContext(ctx).Create(&members).Error
}

func (s *Storage) GetOPCGroupMembers(ctx context.Context, groupUuid string) ([]data_models.OPCGroupMember, error) {
	var members []data_models.OPCGroupMember
	err := s.sqliteDB.WithContext(ctx).Where("group_uuid = ?", groupUuid).Find(&members).Error
	return members, err
}

func (s *Storage) DeleteOPCGroupMembersByGroup(ctx context.Context, groupUuid string) error {
	return s.sqliteDB.WithContext(ctx).Where("group_uuid = ?", groupUuid).Delete(&data_models.OPCGroupMember{}).Error
}

func (s *Storage) DeleteOPCGroupMembersByPerson(ctx context.Context, personUuid string) error {
	return s.sqliteDB.WithContext(ctx).Where("person_uuid = ?", personUuid).Delete(&data_models.OPCGroupMember{}).Error
}

// ==================== Chat 扩展 ====================

// ClearChatMessages 清除某个聊天的所有消息（保留 Chat 记录）
func (s *Storage) ClearChatMessages(ctx context.Context, chatUuid string) error {
	return s.sqliteDB.WithContext(ctx).Where("chat_uuid = ?", chatUuid).Delete(&data_models.Message{}).Error
}

// CreateChatWithType 创建指定类型的对话
func (s *Storage) CreateChatWithType(ctx context.Context, chatUuid, title, chatType string) error {
	now := time.Now()
	chat := &data_models.Chat{
		OrmModel: data_models.OrmModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Uuid:     chatUuid,
		Title:    title,
		ChatType: chatType,
	}
	return s.sqliteDB.WithContext(ctx).Create(chat).Error
}

// DeleteChatAndMessages 删除对话及其所有消息
func (s *Storage) DeleteChatAndMessages(ctx context.Context, chatUuid string) error {
	if err := s.sqliteDB.WithContext(ctx).Where("chat_uuid = ?", chatUuid).Delete(&data_models.Message{}).Error; err != nil {
		return err
	}
	return s.sqliteDB.WithContext(ctx).Where("uuid = ?", chatUuid).Delete(&data_models.Chat{}).Error
}

// GetLastMessage 获取某个聊天的最后一条消息
func (s *Storage) GetLastMessage(ctx context.Context, chatUuid string) (*data_models.Message, error) {
	var msg data_models.Message
	err := s.sqliteDB.WithContext(ctx).
		Where("chat_uuid = ?", chatUuid).
		Order("created_at DESC").
		First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &msg, err
}
