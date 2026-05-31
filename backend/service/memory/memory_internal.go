package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory/memory_dto"
)

const memoryTimeFormat = time.RFC3339

func toMemoryItems(memories []data_models.Memory) []memory_dto.MemoryItem {
	items := make([]memory_dto.MemoryItem, 0, len(memories))
	for _, memory := range memories {
		items = append(items, toMemoryItem(memory))
	}
	return items
}

func toMemoryItem(memory data_models.Memory) memory_dto.MemoryItem {
	return memory_dto.MemoryItem{
		ID:             memory.ID,
		Summary:        memory.Summary,
		Content:        memory.Content,
		Type:           memory.Type,
		Target:         memory.Target,
		Source:         memory.Source,
		CharCount:      memory.CharCount,
		EmbeddingID:    memory.EmbeddingID,
		IsForgotten:    memory.IsForgotten,
		RecallCount:    memory.RecallCount,
		LastRecalledAt: formatOptionalTime(memory.LastRecalledAt),
		LastUsedAt:     formatOptionalTime(memory.LastUsedAt),
		Importance:     memory.Importance,
		Confidence:     memory.Confidence,
		Pinned:         memory.Pinned,
		Created:        memory.CreatedAt.Format(memoryTimeFormat),
		Updated:        memory.UpdatedAt.Format(memoryTimeFormat),
	}
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(memoryTimeFormat)
	return &formatted
}

func (m *Memory) loadConfig() (*data_models.Config, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}
	config := &data_models.Config{}
	bytes, err := os.ReadFile(filepath.Join(dataDir, dir.ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			config.DataDir = dataDir
			config.Memory.Enabled = true
			return config, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	if !strings.Contains(string(bytes), `"memory"`) {
		config.Memory.Enabled = true
	}
	if config.DataDir == "" {
		config.DataDir = dataDir
	}
	return config, nil
}

func (m *Memory) saveConfig(config *data_models.Config) error {
	dataDir := config.DataDir
	if dataDir == "" {
		var err error
		dataDir, err = dir.GetDataDir()
		if err != nil {
			return err
		}
		config.DataDir = dataDir
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dataDir, dir.ConfigFileName), content, 0o644)
}
