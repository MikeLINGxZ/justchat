package storage

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func (s *Storage) GetChats(ctx context.Context, offset, limit int, keyword *string) ([]data_models.Chat, error) {
	var res []data_models.Chat
	err := s.sqliteDB.Model(&data_models.Chat{}).Offset(offset).Limit(limit).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
