package storage

import (
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// SaveTaskState upserts one key/value pair for a task session.
func (s *Storage) SaveTaskState(sessionID uint, key, value string) error {
	var state data_models.TaskState
	err := s.sqliteDB.Where("session_id = ? AND key = ?", sessionID, key).First(&state).Error
	if err == nil {
		return s.sqliteDB.Model(&state).Update("value", value).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.sqliteDB.Create(&data_models.TaskState{
		SessionID: sessionID,
		Key:       key,
		Value:     value,
	}).Error
}

// LoadTaskState returns one saved key/value pair for a task session.
func (s *Storage) LoadTaskState(sessionID uint, key string) (string, bool, error) {
	var state data_models.TaskState
	if err := s.sqliteDB.Where("session_id = ? AND key = ?", sessionID, key).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, nil
		}
		return "", false, err
	}
	return state.Value, true, nil
}

// DeleteTaskState removes one saved key for a task session.
func (s *Storage) DeleteTaskState(sessionID uint, key string) error {
	return s.sqliteDB.Where("session_id = ? AND key = ?", sessionID, key).Delete(&data_models.TaskState{}).Error
}
