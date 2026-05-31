package memory

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory/memory_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

// Memory exposes long-term memory management operations to the frontend.
type Memory struct {
	stor *storage.Storage
}

// NewMemory creates a memory service backed by application storage.
func NewMemory(stor *storage.Storage) *Memory {
	return &Memory{stor: stor}
}

// ListMemories returns paginated long-term memories for the management page.
func (m *Memory) ListMemories(ctx context.Context, input memory_dto.ListMemoriesInput) (*memory_dto.ListMemoriesOutput, error) {
	items, total, err := m.stor.ListMemories(storage.MemoryListFilter{
		Query:            input.Query,
		Type:             input.Type,
		Target:           input.Target,
		IncludeForgotten: input.IncludeForgotten,
		Offset:           input.Offset,
		Limit:            input.Limit,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemoryList, err)
	}
	return &memory_dto.ListMemoriesOutput{
		Items: toMemoryItems(items),
		Total: total,
	}, nil
}

// GetMemory returns one memory by id.
func (m *Memory) GetMemory(ctx context.Context, input memory_dto.GetMemoryInput) (*memory_dto.GetMemoryOutput, error) {
	memory, err := m.stor.GetMemory(input.ID)
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemoryGet, err)
	}
	return &memory_dto.GetMemoryOutput{Memory: toMemoryItem(*memory)}, nil
}

// CreateMemory creates a manual long-term memory.
func (m *Memory) CreateMemory(ctx context.Context, input memory_dto.CreateMemoryInput) (*memory_dto.CreateMemoryOutput, error) {
	created, err := m.stor.CreateMemory(data_models.Memory{
		Summary:    input.Summary,
		Content:    input.Content,
		Type:       input.Type,
		Target:     input.Target,
		Source:     input.Source,
		Importance: input.Importance,
		Confidence: input.Confidence,
		Pinned:     input.Pinned,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemoryCreate, err)
	}
	return &memory_dto.CreateMemoryOutput{Memory: toMemoryItem(*created)}, nil
}

// UpdateMemory updates editable fields on a long-term memory.
func (m *Memory) UpdateMemory(ctx context.Context, input memory_dto.UpdateMemoryInput) (*memory_dto.UpdateMemoryOutput, error) {
	updated, err := m.stor.UpdateMemory(input.ID, storage.MemoryUpdate{
		Summary:    &input.Summary,
		Content:    &input.Content,
		Type:       &input.Type,
		Target:     &input.Target,
		Source:     &input.Source,
		Importance: &input.Importance,
		Confidence: &input.Confidence,
		Pinned:     &input.Pinned,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemoryUpdate, err)
	}
	return &memory_dto.UpdateMemoryOutput{Memory: toMemoryItem(*updated)}, nil
}

// ForgetMemory hides a memory from future retrieval without deleting the row.
func (m *Memory) ForgetMemory(ctx context.Context, input memory_dto.ForgetMemoryInput) (*memory_dto.ForgetMemoryOutput, error) {
	if err := m.stor.ForgetMemory(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrMemoryForget, err)
	}
	return &memory_dto.ForgetMemoryOutput{}, nil
}

// RestoreMemory reactivates a forgotten memory.
func (m *Memory) RestoreMemory(ctx context.Context, input memory_dto.RestoreMemoryInput) (*memory_dto.RestoreMemoryOutput, error) {
	if err := m.stor.RestoreMemory(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrMemoryRestore, err)
	}
	return &memory_dto.RestoreMemoryOutput{}, nil
}

// GetMemoryStats returns aggregate counts for the memory settings page.
func (m *Memory) GetMemoryStats(ctx context.Context, input memory_dto.GetMemoryStatsInput) (*memory_dto.GetMemoryStatsOutput, error) {
	stats, err := m.stor.MemoryStats()
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemoryStats, err)
	}
	return &memory_dto.GetMemoryStatsOutput{
		Stats: memory_dto.MemoryStatsItem{
			Total:     stats.Total,
			Active:    stats.Active,
			Forgotten: stats.Forgotten,
			ByTarget:  stats.ByTarget,
			ByType:    stats.ByType,
		},
	}, nil
}

// GetMemorySettings loads the persisted memory feature switch.
func (m *Memory) GetMemorySettings(ctx context.Context, input memory_dto.GetMemorySettingsInput) (*memory_dto.GetMemorySettingsOutput, error) {
	config, err := m.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemorySettings, err)
	}
	return &memory_dto.GetMemorySettingsOutput{Enabled: config.Memory.Enabled}, nil
}

// SaveMemorySettings persists the memory feature switch.
func (m *Memory) SaveMemorySettings(ctx context.Context, input memory_dto.SaveMemorySettingsInput) (*memory_dto.SaveMemorySettingsOutput, error) {
	config, err := m.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrMemorySettings, err)
	}
	config.Memory.Enabled = input.Enabled
	if err := m.saveConfig(config); err != nil {
		return nil, ierror.Error(ierror.ErrMemorySettings, err)
	}
	return &memory_dto.SaveMemorySettingsOutput{}, nil
}
