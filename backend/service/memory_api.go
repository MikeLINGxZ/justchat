package service

import (
	"context"

	"fmt"
	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// GetMemories 分页查询记忆列表。
func (s *Service) GetMemories(query view_models.MemoryListQuery) (*view_models.MemoryListResponse, error) {
	if s.memoryStorage == nil {
		return &view_models.MemoryListResponse{}, nil
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	ctx := context.Background()
	memories, total, err := s.memoryStorage.ListMemories(ctx, query.Offset, query.Limit, query.Keyword, query.Type, query.IsForgotten)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	items := make([]view_models.Memory, 0, len(memories))
	for _, m := range memories {
		items = append(items, toMemoryViewModel(m))
	}

	return &view_models.MemoryListResponse{
		Memories: items,
		Total:    total,
	}, nil
}

// GetMemoryDetail 获取单条记忆详情。
func (s *Service) GetMemoryDetail(id uint) (*view_models.Memory, error) {
	if s.memoryStorage == nil {
		return nil, fmt.Errorf("memory system not initialized")
	}
	ctx := context.Background()
	m, err := s.memoryStorage.GetMemoryByID(ctx, id)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	vm := toMemoryViewModel(*m)
	return &vm, nil
}

// DeleteMemory 软删除记忆。
func (s *Service) DeleteMemory(id uint) error {
	if s.memoryStorage == nil {
		return fmt.Errorf("memory system not initialized")
	}
	return s.memoryStorage.SoftDeleteMemory(context.Background(), id)
}

// RestoreMemory 恢复已遗忘的记忆。
func (s *Service) RestoreMemory(id uint) error {
	if s.memoryStorage == nil {
		return fmt.Errorf("memory system not initialized")
	}
	return s.memoryStorage.RestoreMemory(context.Background(), id)
}

// GetMemoryStats 获取记忆统计。
func (s *Service) GetMemoryStats() (*view_models.MemoryStats, error) {
	if s.memoryStorage == nil {
		return &view_models.MemoryStats{}, nil
	}
	total, weekNew, forgotten, err := s.memoryStorage.GetMemoryStats(context.Background())
	if err != nil {
		return nil, ierror.NewError(err)
	}
	return &view_models.MemoryStats{
		Total:     total,
		WeekNew:   weekNew,
		Forgotten: forgotten,
	}, nil
}

func toMemoryViewModel(m memory_models.Memory) view_models.Memory {
	return view_models.Memory{
		ID:               m.ID,
		Summary:          m.Summary,
		Content:          m.Content,
		Type:             string(m.Type),
		TimeRangeStart:   m.TimeRangStart,
		TimeRangeEnd:     m.TimeRangeEnd,
		Location:         m.Location,
		Characters:       m.Characters,
		Importance:       m.Importance,
		EmotionalValence: m.EmotionalValence,
		TrustScore:       m.TrustScore,
		IsForgotten:      m.IsForgotten,
		RecallCount:      m.RecallCount,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
