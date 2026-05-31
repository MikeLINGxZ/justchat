package storage

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

// CreateNotification inserts a new notification row.
func (s *Storage) CreateNotification(notification data_models.Notification) (*data_models.Notification, error) {
	if err := s.sqliteDB.Create(&notification).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// ListNotifications returns recent notifications, optionally including resolved rows.
func (s *Storage) ListNotifications(includeResolved bool) ([]data_models.Notification, error) {
	var items []data_models.Notification
	query := s.sqliteDB.Order("created_at DESC")
	if !includeResolved {
		query = query.Where("resolved_at IS NULL")
	}
	if err := query.Limit(200).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ResolveNotification sets the notification's resolved timestamp.
func (s *Storage) ResolveNotification(id uint) error {
	now := time.Now()
	return s.sqliteDB.Model(&data_models.Notification{}).Where("id = ?", id).Update("resolved_at", &now).Error
}

// DeleteNotification permanently removes a notification row.
func (s *Storage) DeleteNotification(id uint) error {
	return s.sqliteDB.Delete(&data_models.Notification{}, id).Error
}
