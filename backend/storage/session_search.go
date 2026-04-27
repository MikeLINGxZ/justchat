package storage

import (
	"context"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type SessionSearchResult struct {
	ChatUuid    string `json:"chat_uuid"`
	MessageUuid string `json:"message_uuid"`
	Role        string `json:"role"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
}

func (s *Storage) AutoMigrateSessionSearch(ctx context.Context) error {
	db := s.sqliteDB.WithContext(ctx)
	if err := db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
		message_uuid UNINDEXED,
		chat_uuid UNINDEXED,
		role UNINDEXED,
		content,
		created_at UNINDEXED
	)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`INSERT INTO messages_fts(message_uuid, chat_uuid, role, content, created_at)
		SELECT message_uuid, chat_uuid, role, content, created_at
		FROM messages
		WHERE COALESCE(content, '') <> ''
		  AND message_uuid NOT IN (SELECT message_uuid FROM messages_fts)`).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) SearchSessions(ctx context.Context, query string, limit int) ([]SessionSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}
	if err := s.AutoMigrateSessionSearch(ctx); err == nil {
		results, searchErr := s.searchSessionsFTS(ctx, query, limit)
		if searchErr == nil {
			return results, nil
		}
	}
	return s.searchSessionsLike(ctx, query, limit)
}

func (s *Storage) searchSessionsFTS(ctx context.Context, query string, limit int) ([]SessionSearchResult, error) {
	terms := strings.Fields(query)
	if len(terms) == 0 {
		terms = []string{query}
	}
	cleaned := make([]string, 0, len(terms))
	for _, term := range terms {
		term = strings.TrimSpace(strings.ReplaceAll(term, `"`, `""`))
		if term != "" {
			cleaned = append(cleaned, `"`+term+`"`)
		}
	}
	if len(cleaned) == 0 {
		return nil, nil
	}
	match := strings.Join(cleaned, " OR ")
	var rows []SessionSearchResult
	err := s.sqliteDB.WithContext(ctx).
		Raw(`SELECT chat_uuid, message_uuid, role, content, created_at
			FROM messages_fts
			WHERE messages_fts MATCH ?
			ORDER BY rank
			LIMIT ?`, match, limit).
		Scan(&rows).Error
	return rows, err
}

func (s *Storage) searchSessionsLike(ctx context.Context, query string, limit int) ([]SessionSearchResult, error) {
	pattern := "%" + query + "%"
	var messages []data_models.Message
	if err := s.sqliteDB.WithContext(ctx).
		Where("content LIKE ?", pattern).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	results := make([]SessionSearchResult, 0, len(messages))
	for _, msg := range messages {
		results = append(results, SessionSearchResult{
			ChatUuid:    msg.ChatUuid,
			MessageUuid: msg.MessageUuid,
			Role:        string(msg.Role),
			Content:     msg.Content,
			CreatedAt:   msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return results, nil
}
