package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"fmt"
	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
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

// UpdateMemory 编辑记忆并在可用时自动重建 embedding。
func (s *Service) UpdateMemory(id uint, input view_models.MemoryUpdateInput) (*view_models.Memory, error) {
	if s.memoryStorage == nil {
		return nil, fmt.Errorf("memory system not initialized")
	}

	summary := strings.TrimSpace(input.Summary)
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, ierror.NewError(errors.New("memory content cannot be empty"))
	}
	if input.Importance < 0 || input.Importance > 1 {
		return nil, ierror.NewError(errors.New("importance must be between 0 and 1"))
	}
	if input.EmotionalValence < -1 || input.EmotionalValence > 1 {
		return nil, ierror.NewError(errors.New("emotional valence must be between -1 and 1"))
	}

	startAt, err := parseMemoryDate(input.TimeRangeStart)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	endAt, err := parseMemoryDate(input.TimeRangeEnd)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if startAt != nil && endAt != nil && startAt.After(*endAt) {
		return nil, ierror.NewError(errors.New("start time cannot be after end time"))
	}

	update := memory_models.Memory{
		Summary:          summary,
		Content:          content,
		Type:             memory_models.MemoryType(strings.TrimSpace(input.Type)),
		TimeRangStart:    startAt,
		TimeRangeEnd:     endAt,
		Location:         normalizeOptionalString(input.Location),
		Characters:       normalizeOptionalString(input.Characters),
		Importance:       input.Importance,
		EmotionalValence: input.EmotionalValence,
	}

	if err := s.memoryStorage.ReplaceMemoryEditableFields(context.Background(), id, update); err != nil {
		return nil, ierror.NewError(err)
	}

	if s.memorySearcher != nil && s.memorySearcher.HasEmbedder() {
		go func(memoryID uint, text string) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if embedErr := s.memorySearcher.EmbedAndStore(bgCtx, memoryID, text); embedErr != nil {
				logger.Error("re-embed memory error:", embedErr)
				return
			}
			s.memorySearcher.RefreshCache()
		}(id, strings.TrimSpace(summary+" "+content))
	}

	return s.GetMemoryDetail(id)
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
		HasEmbedding:     m.EmbeddingID != nil,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func parseMemoryDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	layouts := []string{
		time.DateOnly,
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return &parsed, nil
		}
	}
	return nil, fmt.Errorf("invalid date format: %s", trimmed)
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
