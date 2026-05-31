package storage

import (
	"errors"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateTerminal persists a terminal session row with default start metadata.
func (s *Storage) CreateTerminal(term data_models.Terminal) (*data_models.Terminal, error) {
	if term.StartedAt.IsZero() {
		term.StartedAt = time.Now()
	}
	if term.Status == "" {
		term.Status = "active"
	}
	if err := s.sqliteDB.Create(&term).Error; err != nil {
		return nil, err
	}
	return &term, nil
}

// GetTerminal returns a terminal session by its stable terminal ID.
func (s *Storage) GetTerminal(terminalID string) (*data_models.Terminal, error) {
	var term data_models.Terminal
	if err := s.sqliteDB.Where("terminal_id = ?", terminalID).First(&term).Error; err != nil {
		return nil, err
	}
	return &term, nil
}

// ListTerminalsForSession returns terminal sessions belonging to a chat session.
func (s *Storage) ListTerminalsForSession(sessionID uint) ([]data_models.Terminal, error) {
	var terms []data_models.Terminal
	if err := s.sqliteDB.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&terms).Error; err != nil {
		return nil, err
	}
	return terms, nil
}

// UpdateTerminalVisibility changes whether a persisted terminal should be shown in chat.
func (s *Storage) UpdateTerminalVisibility(terminalID string, visible bool, title string) error {
	updates := map[string]any{"visible": visible}
	if title != "" {
		updates["title"] = title
	}
	return s.sqliteDB.Model(&data_models.Terminal{}).
		Where("terminal_id = ?", terminalID).
		Updates(updates).Error
}

// FinishTerminal records terminal process completion state.
func (s *Storage) FinishTerminal(terminalID string, status string, exitCode int) error {
	now := time.Now()
	return s.sqliteDB.Model(&data_models.Terminal{}).
		Where("terminal_id = ?", terminalID).
		Updates(map[string]any{
			"status":    status,
			"exit_code": exitCode,
			"ended_at":  &now,
		}).Error
}

// AppendTerminalOutput stores one terminal output chunk and advances the byte cursor.
func (s *Storage) AppendTerminalOutput(terminalID string, data string) (*data_models.TerminalOutputChunk, error) {
	var chunk data_models.TerminalOutputChunk
	err := s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		var term data_models.Terminal
		if err := tx.Where("terminal_id = ?", terminalID).First(&term).Error; err != nil {
			return err
		}

		var last data_models.TerminalOutputChunk
		err := tx.Where("terminal_id = ?", terminalID).Order("seq DESC").First(&last).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// Cursor values are byte offsets so frontend polling can resume from
		// an exact chunk boundary even when chunks contain multi-byte text.
		nextSeq := last.Seq + 1
		start := term.CurrentCursor
		end := start + int64(len([]byte(data)))
		chunk = data_models.TerminalOutputChunk{
			TerminalID:  terminalID,
			Seq:         nextSeq,
			CursorStart: start,
			CursorEnd:   end,
			Data:        data,
		}
		if err := tx.Create(&chunk).Error; err != nil {
			return err
		}
		return tx.Model(&term).Update("current_cursor", end).Error
	})
	if err != nil {
		return nil, err
	}
	return &chunk, nil
}

// ReadTerminalOutput returns chunks after cursor, trimming the first partial chunk if needed.
func (s *Storage) ReadTerminalOutput(terminalID string, cursor int64) ([]data_models.TerminalOutputChunk, error) {
	var chunks []data_models.TerminalOutputChunk
	if err := s.sqliteDB.Where("terminal_id = ? AND cursor_end > ?", terminalID, cursor).
		Order("seq ASC").
		Find(&chunks).Error; err != nil {
		return nil, err
	}
	for i := range chunks {
		if cursor > chunks[i].CursorStart && cursor < chunks[i].CursorEnd {
			// Trim only the already-read bytes from the first overlapping chunk;
			// later chunks are returned unchanged.
			offset := cursor - chunks[i].CursorStart
			data := []byte(chunks[i].Data)
			if offset < int64(len(data)) {
				chunks[i].Data = string(data[offset:])
				chunks[i].CursorStart = cursor
			}
		}
	}
	return chunks, nil
}
