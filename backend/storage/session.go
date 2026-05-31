package storage

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateSession inserts a new session row.
func (s *Storage) CreateSession(session data_models.Session) (*data_models.Session, error) {
	if err := s.sqliteDB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSession loads one session by primary key.
func (s *Storage) GetSession(id uint) (*data_models.Session, error) {
	var session data_models.Session
	if err := s.sqliteDB.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessions returns sessions ordered from newest to oldest, optionally including hidden legacy background sessions.
func (s *Storage) ListSessions(cursor uint, limit int, starredOnly bool, includeHidden bool) ([]data_models.Session, error) {
	var sessions []data_models.Session
	q := s.sqliteDB.Order("updated_at DESC")
	if cursor > 0 {
		q = q.Where("id < ?", cursor)
	}
	if starredOnly {
		q = q.Where("starred = ?", true)
	}
	if !includeHidden {
		q = q.Where("kind <> ? OR kind = ''", "background")
	}
	if err := q.Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// UpdateSessionTitle updates the persisted session title.
func (s *Storage) UpdateSessionTitle(id uint, title string) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("title", title).Error
}

// UpdateSessionStatus updates the persisted session status.
func (s *Storage) UpdateSessionStatus(id uint, status string) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateSessionStarred updates the persisted favorite flag.
func (s *Storage) UpdateSessionStarred(id uint, starred bool) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("starred", starred).Error
}

// DeleteSession removes the session and all of its messages in one transaction.
func (s *Storage) DeleteSession(id uint) error {
	return s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = ?", id).Delete(&data_models.Message{}).Error; err != nil {
			return err
		}
		return tx.Delete(&data_models.Session{}, id).Error
	})
}

// TouchSession bumps updated_at so recent sessions sort correctly.
func (s *Storage) TouchSession(id uint) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
