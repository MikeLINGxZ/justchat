package storage

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

func (s *Storage) CreateMessage(msg data_models.Message) (*data_models.Message, error) {
	if err := s.sqliteDB.Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *Storage) CreateMessages(msgs *[]data_models.Message) error {
	if len(*msgs) == 0 {
		return nil
	}
	return s.sqliteDB.Create(msgs).Error
}

func (s *Storage) ListMessagesForSession(sessionID uint, offset int, limit int) ([]data_models.Message, error) {
	var msgs []data_models.Message
	if err := s.sqliteDB.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *Storage) CountMessagesForSession(sessionID uint) (int64, error) {
	var count int64
	if err := s.sqliteDB.Model(&data_models.Message{}).Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Storage) DeleteMessagesForSession(sessionID uint) error {
	return s.sqliteDB.Where("session_id = ?", sessionID).Delete(&data_models.Message{}).Error
}
