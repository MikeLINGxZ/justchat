package storage

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

const (
	defaultMemoryType       = "information"
	defaultMemoryTarget     = "user"
	defaultMemorySource     = "agent"
	defaultMemoryImportance = 50
	defaultMemoryConfidence = 80
)

// MemoryListFilter describes list filters for memory management views.
type MemoryListFilter struct {
	Query            string
	Type             string
	Target           string
	IncludeForgotten bool
	Offset           int
	Limit            int
}

// MemoryUpdate carries partial memory changes.
type MemoryUpdate struct {
	Summary     *string
	Content     *string
	Type        *string
	Target      *string
	Source      *string
	Importance  *int
	Confidence  *int
	Pinned      *bool
	IsForgotten *bool
}

// MemoryStatsSummary contains aggregate memory counts for settings UI.
type MemoryStatsSummary struct {
	Total     int64            `json:"total"`
	Active    int64            `json:"active"`
	Forgotten int64            `json:"forgotten"`
	ByTarget  map[string]int64 `json:"by_target"`
	ByType    map[string]int64 `json:"by_type"`
}

// CreateMemory persists a memory after normalizing defaults and character count.
func (s *Storage) CreateMemory(memory data_models.Memory) (*data_models.Memory, error) {
	normalizeMemory(&memory)
	if err := s.sqliteDB.Create(&memory).Error; err != nil {
		return nil, err
	}
	return &memory, nil
}

// GetMemory loads a single memory by id.
func (s *Storage) GetMemory(id uint) (*data_models.Memory, error) {
	var memory data_models.Memory
	if err := s.sqliteDB.First(&memory, id).Error; err != nil {
		return nil, err
	}
	return &memory, nil
}

// ListMemories returns paginated memories and the total matching count.
func (s *Storage) ListMemories(filter MemoryListFilter) ([]data_models.Memory, int64, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	query := s.applyMemoryFilter(s.sqliteDB.Model(&data_models.Memory{}), filter)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var memories []data_models.Memory
	err := query.
		Order("pinned DESC").
		Order("importance DESC").
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&memories).Error
	if err != nil {
		return nil, 0, err
	}
	return memories, total, nil
}

// UpdateMemory applies partial changes to a memory and refreshes derived fields.
func (s *Storage) UpdateMemory(id uint, update MemoryUpdate) (*data_models.Memory, error) {
	memory, err := s.GetMemory(id)
	if err != nil {
		return nil, err
	}
	if update.Summary != nil {
		memory.Summary = strings.TrimSpace(*update.Summary)
	}
	if update.Content != nil {
		memory.Content = strings.TrimSpace(*update.Content)
	}
	if update.Type != nil {
		memory.Type = normalizeMemoryType(*update.Type)
	}
	if update.Target != nil {
		memory.Target = normalizeMemoryTarget(*update.Target)
	}
	if update.Source != nil {
		memory.Source = normalizeMemorySource(*update.Source)
	}
	if update.Importance != nil {
		memory.Importance = clampScore(*update.Importance)
	}
	if update.Confidence != nil {
		memory.Confidence = clampScore(*update.Confidence)
	}
	if update.Pinned != nil {
		memory.Pinned = *update.Pinned
	}
	if update.IsForgotten != nil {
		memory.IsForgotten = *update.IsForgotten
	}
	normalizeMemory(memory)
	if err := s.sqliteDB.Save(memory).Error; err != nil {
		return nil, err
	}
	return memory, nil
}

// ForgetMemory soft-hides a memory from injection and retrieval.
func (s *Storage) ForgetMemory(id uint) error {
	return s.sqliteDB.Model(&data_models.Memory{}).Where("id = ?", id).Update("is_forgotten", true).Error
}

// RestoreMemory returns a forgotten memory to the active set.
func (s *Storage) RestoreMemory(id uint) error {
	return s.sqliteDB.Model(&data_models.Memory{}).Where("id = ?", id).Update("is_forgotten", false).Error
}

// SearchMemories performs a lightweight LIKE search and updates recall metadata.
func (s *Storage) SearchMemories(queryText string, limit int) ([]data_models.Memory, error) {
	queryText = strings.TrimSpace(queryText)
	if limit <= 0 || limit > 20 {
		limit = 8
	}

	db := s.sqliteDB.Model(&data_models.Memory{}).Where("is_forgotten = ?", false)
	if queryText != "" {
		db = applyMemoryQuery(db, queryText)
	}

	var memories []data_models.Memory
	if err := db.
		Order("pinned DESC").
		Order("importance DESC").
		Order("updated_at DESC").
		Limit(limit).
		Find(&memories).Error; err != nil {
		return nil, err
	}
	if len(memories) == 0 && queryText != "" {
		return s.SearchMemories("", limit)
	}
	if len(memories) > 0 {
		now := time.Now()
		ids := make([]uint, 0, len(memories))
		for i := range memories {
			ids = append(ids, memories[i].ID)
			memories[i].RecallCount++
			memories[i].LastRecalledAt = &now
		}
		if err := s.sqliteDB.Model(&data_models.Memory{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"recall_count":     gorm.Expr("recall_count + 1"),
				"last_recalled_at": now,
			}).Error; err != nil {
			return nil, err
		}
	}
	return memories, nil
}

// RenderCoreMemory formats active memories into stable prompt text.
func (s *Storage) RenderCoreMemory(userCharLimit int, assistantCharLimit int) (string, error) {
	userMemories, err := s.coreMemoriesForTarget("user", userCharLimit)
	if err != nil {
		return "", err
	}
	assistantMemories, err := s.coreMemoriesForTarget("memory", assistantCharLimit)
	if err != nil {
		return "", err
	}
	if len(userMemories) == 0 && len(assistantMemories) == 0 {
		return "", nil
	}

	var builder strings.Builder
	if len(userMemories) > 0 {
		builder.WriteString("User memories:\n")
		for _, memory := range userMemories {
			builder.WriteString("- ")
			builder.WriteString(memory.Content)
			builder.WriteString("\n")
		}
	}
	if len(assistantMemories) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString("Assistant memories:\n")
		for _, memory := range assistantMemories {
			builder.WriteString("- ")
			builder.WriteString(memory.Content)
			builder.WriteString("\n")
		}
	}
	return strings.TrimSpace(builder.String()), nil
}

// MemoryStats returns aggregate memory counts.
func (s *Storage) MemoryStats() (MemoryStatsSummary, error) {
	stats := MemoryStatsSummary{
		ByTarget: map[string]int64{},
		ByType:   map[string]int64{},
	}
	if err := s.sqliteDB.Model(&data_models.Memory{}).Count(&stats.Total).Error; err != nil {
		return stats, err
	}
	if err := s.sqliteDB.Model(&data_models.Memory{}).Where("is_forgotten = ?", false).Count(&stats.Active).Error; err != nil {
		return stats, err
	}
	if err := s.sqliteDB.Model(&data_models.Memory{}).Where("is_forgotten = ?", true).Count(&stats.Forgotten).Error; err != nil {
		return stats, err
	}

	var targetRows []struct {
		Target string
		Count  int64
	}
	if err := s.sqliteDB.Model(&data_models.Memory{}).Select("target, count(*) as count").Group("target").Scan(&targetRows).Error; err != nil {
		return stats, err
	}
	for _, row := range targetRows {
		stats.ByTarget[row.Target] = row.Count
	}

	var typeRows []struct {
		Type  string
		Count int64
	}
	if err := s.sqliteDB.Model(&data_models.Memory{}).Select("type, count(*) as count").Group("type").Scan(&typeRows).Error; err != nil {
		return stats, err
	}
	for _, row := range typeRows {
		stats.ByType[row.Type] = row.Count
	}
	return stats, nil
}

func (s *Storage) applyMemoryFilter(db *gorm.DB, filter MemoryListFilter) *gorm.DB {
	if !filter.IncludeForgotten {
		db = db.Where("is_forgotten = ?", false)
	}
	if strings.TrimSpace(filter.Type) != "" {
		db = db.Where("type = ?", normalizeMemoryType(filter.Type))
	}
	if strings.TrimSpace(filter.Target) != "" {
		db = db.Where("target = ?", normalizeMemoryTarget(filter.Target))
	}
	if strings.TrimSpace(filter.Query) != "" {
		db = applyMemoryQuery(db, filter.Query)
	}
	return db
}

func applyMemoryQuery(db *gorm.DB, queryText string) *gorm.DB {
	terms := memorySearchTerms(queryText)
	if len(terms) == 0 {
		return db
	}
	conditions := make([]string, 0, len(terms))
	args := make([]any, 0, len(terms)*2)
	for _, term := range terms {
		conditions = append(conditions, "(summary LIKE ? ESCAPE '\\' OR content LIKE ? ESCAPE '\\')")
		like := "%" + escapeLike(term) + "%"
		args = append(args, like, like)
	}
	return db.Where(strings.Join(conditions, " OR "), args...)
}

func memorySearchTerms(queryText string) []string {
	fields := strings.FieldsFunc(strings.ToLower(strings.TrimSpace(queryText)), func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == ',' || r == '.' || r == ';' || r == ':' || r == '，' || r == '。' || r == '；' || r == '：'
	})
	seen := map[string]struct{}{}
	terms := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if _, ok := seen[field]; ok {
			continue
		}
		seen[field] = struct{}{}
		terms = append(terms, field)
	}
	if len(terms) == 0 && strings.TrimSpace(queryText) != "" {
		terms = append(terms, strings.TrimSpace(queryText))
	}
	return terms
}

func (s *Storage) coreMemoriesForTarget(target string, limit int) ([]data_models.Memory, error) {
	var all []data_models.Memory
	if err := s.sqliteDB.Where("target = ? AND is_forgotten = ?", target, false).
		Order("pinned DESC").
		Order("importance DESC").
		Order("updated_at DESC").
		Find(&all).Error; err != nil {
		return nil, err
	}
	if limit <= 0 {
		return all, nil
	}
	selected := make([]data_models.Memory, 0, len(all))
	used := 0
	for _, memory := range all {
		next := used + memory.CharCount
		if next > limit && len(selected) > 0 {
			continue
		}
		selected = append(selected, memory)
		used = next
	}
	return selected, nil
}

func normalizeMemory(memory *data_models.Memory) {
	memory.Summary = strings.TrimSpace(memory.Summary)
	memory.Content = strings.TrimSpace(memory.Content)
	memory.Type = normalizeMemoryType(memory.Type)
	memory.Target = normalizeMemoryTarget(memory.Target)
	memory.Source = normalizeMemorySource(memory.Source)
	memory.CharCount = utf8.RuneCountInString(memory.Content)
	if memory.Importance == 0 {
		memory.Importance = defaultMemoryImportance
	}
	memory.Importance = clampScore(memory.Importance)
	if memory.Confidence == 0 {
		memory.Confidence = defaultMemoryConfidence
	}
	memory.Confidence = clampScore(memory.Confidence)
}

func normalizeMemoryType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "fact", "information", "event":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return defaultMemoryType
	}
}

func normalizeMemoryTarget(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "user", "memory":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return defaultMemoryTarget
	}
}

func normalizeMemorySource(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return defaultMemorySource
	}
	return value
}

func clampScore(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(strings.TrimSpace(value))
}

func formatMemoryContext(memories []data_models.Memory) string {
	lines := make([]string, 0, len(memories))
	for _, memory := range memories {
		lines = append(lines, fmt.Sprintf("[%s] %s: %s", memory.UpdatedAt.Format("2006-01-02"), memory.Summary, memory.Content))
	}
	return strings.Join(lines, "\n")
}
