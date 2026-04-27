package storage

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gorm.io/gorm"
)

const (
	UserMemoryCharLimit  = 1375
	AgentMemoryCharLimit = 2200
)

func (s *Storage) WriterMemory(ctx context.Context, memory models.Memory) (uint, error) {
	memory = normalizeMemory(memory)
	// 精确 summary 匹配：标题完全相同的视为重复，不新建。
	// 更细粒度的语义去重交由 Memory Agent（在写入前已拿到相关记忆列表并自行决定 write/edit）。
	if memory.Summary != "" {
		var existing models.Memory
		if err := s.sqliteDb.WithContext(ctx).
			Where("summary = ? AND is_forgotten = ?", memory.Summary, false).
			First(&existing).Error; err == nil {
			return existing.ID, nil
		}
	}
	result := s.sqliteDb.WithContext(ctx).Create(&memory)
	return memory.ID, result.Error
}

func normalizeMemory(memory models.Memory) models.Memory {
	memory.Summary = strings.TrimSpace(memory.Summary)
	memory.Content = strings.TrimSpace(memory.Content)
	if memory.Target == "" {
		memory.Target = models.MemoryTargetUser
	}
	if memory.Source == "" {
		memory.Source = "agent"
	}
	memory.CharCount = utf8.RuneCountInString(memory.Content)
	return memory
}

func NormalizeMemoryTarget(raw string) (models.MemoryTarget, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "user", "profile":
		return models.MemoryTargetUser, nil
	case "memory", "agent":
		return models.MemoryTargetAgent, nil
	default:
		return "", fmt.Errorf("unsupported memory target: %q (use user or memory)", raw)
	}
}

func TargetCharLimit(target models.MemoryTarget) int {
	if target == models.MemoryTargetAgent {
		return AgentMemoryCharLimit
	}
	return UserMemoryCharLimit
}

func targetTitle(target models.MemoryTarget) string {
	if target == models.MemoryTargetAgent {
		return "MEMORY (assistant notes)"
	}
	return "USER PROFILE"
}

// RenderCoreMemorySnapshot renders bounded Hermes-style memory blocks for prompt injection.
func (s *Storage) RenderCoreMemorySnapshot(ctx context.Context) (string, error) {
	userBlock, err := s.renderTargetSnapshot(ctx, models.MemoryTargetUser)
	if err != nil {
		return "", err
	}
	agentBlock, err := s.renderTargetSnapshot(ctx, models.MemoryTargetAgent)
	if err != nil {
		return "", err
	}
	blocks := make([]string, 0, 2)
	if userBlock != "" {
		blocks = append(blocks, userBlock)
	}
	if agentBlock != "" {
		blocks = append(blocks, agentBlock)
	}
	return strings.Join(blocks, "\n\n"), nil
}

func (s *Storage) renderTargetSnapshot(ctx context.Context, target models.MemoryTarget) (string, error) {
	entries, used, err := s.ActiveCoreMemories(ctx, target)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", nil
	}
	limit := TargetCharLimit(target)
	percent := 0
	if limit > 0 {
		percent = int(float64(used) / float64(limit) * 100)
	}
	var lines []string
	for _, entry := range entries {
		content := strings.TrimSpace(entry.Content)
		if content != "" {
			lines = append(lines, content)
		}
	}
	if len(lines) == 0 {
		return "", nil
	}
	return fmt.Sprintf("══════════════════════════════════════════════\n%s [%d%% — %d/%d chars]\n══════════════════════════════════════════════\n%s",
		targetTitle(target), percent, used, limit, strings.Join(lines, "\n§\n")), nil
}

func (s *Storage) ActiveCoreMemories(ctx context.Context, target models.MemoryTarget) ([]models.Memory, int, error) {
	var memories []models.Memory
	if err := s.sqliteDb.WithContext(ctx).
		Where("is_forgotten = ? AND target = ?", false, target).
		Order("created_at ASC").
		Find(&memories).Error; err != nil {
		return nil, 0, err
	}
	used := 0
	for _, m := range memories {
		used += utf8.RuneCountInString(strings.TrimSpace(m.Content))
	}
	return memories, used, nil
}

type CoreMemoryMutationResult struct {
	MemoryID       uint
	Message        string
	Usage          string
	CurrentEntries []string
}

func (s *Storage) AddCoreMemory(ctx context.Context, target models.MemoryTarget, content, source string) (*CoreMemoryMutationResult, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("memory content cannot be empty")
	}
	if source == "" {
		source = "agent"
	}
	entries, used, err := s.ActiveCoreMemories(ctx, target)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if strings.TrimSpace(entry.Content) == content {
			return &CoreMemoryMutationResult{
				MemoryID: entry.ID,
				Message:  "duplicate memory already exists; no new entry added",
				Usage:    fmt.Sprintf("%d/%d", used, TargetCharLimit(target)),
			}, nil
		}
	}
	added := utf8.RuneCountInString(content)
	limit := TargetCharLimit(target)
	if used+added > limit {
		return &CoreMemoryMutationResult{
			Message:        fmt.Sprintf("Memory at %d/%d chars. Adding this entry (%d chars) would exceed the limit. Replace or remove existing entries first.", used, limit, added),
			Usage:          fmt.Sprintf("%d/%d", used, limit),
			CurrentEntries: memoryEntryTexts(entries),
		}, nil
	}
	memory := normalizeMemory(models.Memory{
		Summary: makeCoreMemorySummary(content),
		Content: content,
		Target:  target,
		Source:  source,
	})
	id, err := s.WriterMemory(ctx, memory)
	if err != nil {
		return nil, err
	}
	return &CoreMemoryMutationResult{
		MemoryID: id,
		Message:  "memory entry added",
		Usage:    fmt.Sprintf("%d/%d", used+added, limit),
	}, nil
}

func (s *Storage) ReplaceCoreMemory(ctx context.Context, target models.MemoryTarget, oldText, content string) (*CoreMemoryMutationResult, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("replacement content cannot be empty")
	}
	matches, entries, used, err := s.findCoreMemoryMatches(ctx, target, oldText)
	if err != nil {
		return nil, err
	}
	limit := TargetCharLimit(target)
	if len(matches) != 1 {
		return ambiguousMatchResult(matches, entries, used, limit, "replace"), nil
	}
	nextUsed := used - utf8.RuneCountInString(strings.TrimSpace(matches[0].Content)) + utf8.RuneCountInString(content)
	if nextUsed > limit {
		return &CoreMemoryMutationResult{
			Message:        fmt.Sprintf("Memory at %d/%d chars. Replacement would use %d chars. Shorten or remove entries first.", used, limit, nextUsed),
			Usage:          fmt.Sprintf("%d/%d", used, limit),
			CurrentEntries: memoryEntryTexts(entries),
		}, nil
	}
	update := normalizeMemory(models.Memory{
		Summary: makeCoreMemorySummary(content),
		Content: content,
		Target:  target,
		Source:  matches[0].Source,
	})
	if err := s.ReplaceMemoryEditableFields(ctx, matches[0].ID, update); err != nil {
		return nil, err
	}
	return &CoreMemoryMutationResult{MemoryID: matches[0].ID, Message: "memory entry replaced", Usage: fmt.Sprintf("%d/%d", nextUsed, limit)}, nil
}

func (s *Storage) RemoveCoreMemory(ctx context.Context, target models.MemoryTarget, oldText string) (*CoreMemoryMutationResult, error) {
	matches, entries, used, err := s.findCoreMemoryMatches(ctx, target, oldText)
	if err != nil {
		return nil, err
	}
	limit := TargetCharLimit(target)
	if len(matches) != 1 {
		return ambiguousMatchResult(matches, entries, used, limit, "remove"), nil
	}
	if err := s.SoftDeleteMemory(ctx, matches[0].ID); err != nil {
		return nil, err
	}
	nextUsed := used - utf8.RuneCountInString(strings.TrimSpace(matches[0].Content))
	return &CoreMemoryMutationResult{MemoryID: matches[0].ID, Message: "memory entry removed", Usage: fmt.Sprintf("%d/%d", nextUsed, limit)}, nil
}

func (s *Storage) findCoreMemoryMatches(ctx context.Context, target models.MemoryTarget, oldText string) ([]models.Memory, []models.Memory, int, error) {
	oldText = strings.TrimSpace(oldText)
	if oldText == "" {
		return nil, nil, 0, fmt.Errorf("old_text is required")
	}
	entries, used, err := s.ActiveCoreMemories(ctx, target)
	if err != nil {
		return nil, nil, 0, err
	}
	var matches []models.Memory
	for _, entry := range entries {
		if strings.Contains(entry.Content, oldText) {
			matches = append(matches, entry)
		}
	}
	return matches, entries, used, nil
}

func ambiguousMatchResult(matches, entries []models.Memory, used, limit int, action string) *CoreMemoryMutationResult {
	if len(matches) == 0 {
		return &CoreMemoryMutationResult{Message: "no matching memory entry found for " + action, Usage: fmt.Sprintf("%d/%d", used, limit), CurrentEntries: memoryEntryTexts(entries)}
	}
	return &CoreMemoryMutationResult{Message: fmt.Sprintf("old_text matched %d entries; provide a more specific substring", len(matches)), Usage: fmt.Sprintf("%d/%d", used, limit), CurrentEntries: memoryEntryTexts(matches)}
}

func memoryEntryTexts(entries []models.Memory) []string {
	texts := make([]string, 0, len(entries))
	for _, entry := range entries {
		texts = append(texts, strings.TrimSpace(entry.Content))
	}
	return texts
}

func makeCoreMemorySummary(content string) string {
	runes := []rune(strings.TrimSpace(content))
	if len(runes) <= 40 {
		return string(runes)
	}
	return string(runes[:40])
}

// RecentMemoriesByType 返回最近 N 条活跃记忆，用于在 Memory Agent 编码前注入上下文。
// typeFilter 为空则不按类型过滤。
func (s *Storage) RecentMemoriesByType(ctx context.Context, typeFilter string, limit int) ([]models.Memory, error) {
	query := s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("is_forgotten = ?", false)
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	var memories []models.Memory
	err := query.Order("created_at DESC").Limit(limit).Find(&memories).Error
	return memories, err
}

// FTSSearch 使用 FTS5 全文检索搜索记忆，返回按相关性排序的结果。
// keywords 之间使用 OR 语义，limit 为返回上限。
func (s *Storage) FTSSearch(ctx context.Context, keywords []string, limit int) ([]models.Memory, error) {
	var cleaned []string
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if kw != "" {
			// 转义 FTS5 特殊字符
			kw = strings.ReplaceAll(kw, "\"", "\"\"")
			cleaned = append(cleaned, "\""+kw+"\"")
		}
	}
	if len(cleaned) == 0 {
		return nil, nil
	}

	matchExpr := strings.Join(cleaned, " OR ")
	var memories []models.Memory
	err := s.sqliteDb.WithContext(ctx).
		Raw(`SELECT m.* FROM memories m
			 JOIN memories_fts ON memories_fts.rowid = m.id
			 WHERE memories_fts MATCH ? AND m.is_forgotten = 0
			 ORDER BY rank
			 LIMIT ?`, matchExpr, limit).
		Scan(&memories).Error
	if err != nil {
		// FTS5 不可用时降级为 LIKE
		return s.fallbackLIKESearch(ctx, keywords, limit)
	}
	// FTS5 查询成功但无结果时，也尝试 LIKE 兜底
	// （FTS5 的 unicode61 分词器对中文短词可能不够敏感）
	if len(memories) == 0 {
		return s.fallbackLIKESearch(ctx, keywords, limit)
	}
	return memories, nil
}

// fallbackLIKESearch FTS5 不可用时的降级检索。
func (s *Storage) fallbackLIKESearch(ctx context.Context, keywords []string, limit int) ([]models.Memory, error) {
	query := s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("is_forgotten = ?", false)

	var conditions []string
	var args []interface{}
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if kw != "" {
			pattern := "%" + kw + "%"
			conditions = append(conditions, "(summary LIKE ? OR content LIKE ?)")
			args = append(args, pattern, pattern)
		}
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	var memories []models.Memory
	err := query.Order("importance DESC").Limit(limit).Find(&memories).Error
	return memories, err
}

// TopImportantMemories 返回最重要的 N 条记忆（按 importance 降序）。
// 用于关键词检索无结果时的兜底，确保身份/概况类查询仍能获取用户信息。
func (s *Storage) TopImportantMemories(ctx context.Context, limit int) ([]models.Memory, error) {
	var memories []models.Memory
	err := s.sqliteDb.WithContext(ctx).
		Model(&models.Memory{}).
		Where("is_forgotten = ?", false).
		Order("importance DESC, recall_count DESC, created_at DESC").
		Limit(limit).
		Find(&memories).Error
	return memories, err
}

// IncrementRecallCount 更新记忆的召回计数和最后召回时间。
func (s *Storage) IncrementRecallCount(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return s.sqliteDb.WithContext(ctx).
		Model(&models.Memory{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"recall_count":     gorm.Expr("recall_count + 1"),
			"last_recalled_at": now,
			"updated_at":       now,
		}).Error
}

// AdjustTrustScore 非对称信任反馈。delta > 0 为正面反馈，delta < 0 为负面反馈。
func (s *Storage) AdjustTrustScore(ctx context.Context, id uint, delta float64) error {
	return s.sqliteDb.WithContext(ctx).
		Model(&models.Memory{}).
		Where("id = ?", id).
		Update("trust_score",
			gorm.Expr("MIN(1.0, MAX(0.0, trust_score + ?))", delta),
		).Error
}

// ---- 前端 CRUD ----

// ListMemories 分页查询记忆列表。
func (s *Storage) ListMemories(ctx context.Context, offset, limit int, keyword string, memType string, isForgotten bool) ([]models.Memory, int64, error) {
	query := s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("is_forgotten = ?", isForgotten)
	if memType != "" {
		query = query.Where("type = ?", memType)
	}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("summary LIKE ? OR content LIKE ?", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var memories []models.Memory
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&memories).Error
	return memories, total, err
}

// ListMemoriesByTarget paginates memories with optional type and target filters.
func (s *Storage) ListMemoriesByTarget(ctx context.Context, offset, limit int, keyword string, memType string, target string, isForgotten bool) ([]models.Memory, int64, error) {
	query := s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("is_forgotten = ?", isForgotten)
	if memType != "" {
		query = query.Where("type = ?", memType)
	}
	if target != "" {
		query = query.Where("target = ?", target)
	}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("summary LIKE ? OR content LIKE ?", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var memories []models.Memory
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&memories).Error
	return memories, total, err
}

// GetMemoryByID 获取单条记忆详情。
func (s *Storage) GetMemoryByID(ctx context.Context, id uint) (*models.Memory, error) {
	var m models.Memory
	err := s.sqliteDb.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// SoftDeleteMemory 软删除（标记 IsForgotten）。
func (s *Storage) SoftDeleteMemory(ctx context.Context, id uint) error {
	return s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("id = ?", id).
		Update("is_forgotten", true).Error
}

// RestoreMemory 恢复已遗忘的记忆。
func (s *Storage) RestoreMemory(ctx context.Context, id uint) error {
	return s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("id = ?", id).
		Update("is_forgotten", false).Error
}

// GetMemoryStats 返回统计数据。
func (s *Storage) GetMemoryStats(ctx context.Context) (total int64, weekNew int64, forgotten int64, err error) {
	db := s.sqliteDb.WithContext(ctx)
	if err = db.Model(&models.Memory{}).Where("is_forgotten = ?", false).Count(&total).Error; err != nil {
		return
	}
	weekAgo := time.Now().AddDate(0, 0, -7)
	if err = db.Model(&models.Memory{}).Where("is_forgotten = ? AND created_at >= ?", false, weekAgo).Count(&weekNew).Error; err != nil {
		return
	}
	err = db.Model(&models.Memory{}).Where("is_forgotten = ?", true).Count(&forgotten).Error
	return
}

func (s *Storage) ReadMemory(ctx context.Context, keyword string, startAt, endAt *time.Time) ([]models.Memory, error) {
	var memories []models.Memory

	db := s.sqliteDb.WithContext(ctx)

	// 构建基础查询
	query := db.Model(&models.Memory{}).Where("is_forgotten = ?", false) // 排除已遗忘的记忆（可选）

	// 1. 处理关键词搜索：使用 LIKE 搜索 summary 和 content
	if keyword != "" {
		// 使用 LIKE 模糊匹配 summary 或 content
		likePattern := "%" + keyword + "%"
		query = query.Where("summary LIKE ? OR content LIKE ?", likePattern, likePattern)
	}

	// 2. 处理时间范围过滤：要求 TimeRangStart 在 [startAt, endAt] 区间内有交集
	if startAt != nil {
		// 记忆的开始时间小于等于查询结束时间（endAt）
		query = query.Where("time_rang_start <= ?", *endAt)
	}
	if endAt != nil {
		// 记忆的结束时间大于等于查询起始时间（startAt），如果没有 TimeRangeEnd，则视为持续事件？
		// 如果 TimeRangeEnd 为 nil，我们假设它是一个瞬时事件 or 持续有效的事件
		// 这里可以灵活处理：比如如果 End 是 nil，只要 Start <= endAt 即可
		query = query.Where("time_range_end >= ? OR time_range_end IS NULL", *startAt)
	}

	// 3. 排序：优先显示重要性和最近活动的
	query = query.Order("importance DESC").Order("recall_count DESC")

	// 执行查询
	err := query.Find(&memories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to read memories: %w", err)
	}

	return memories, nil
}

type MemoryQuery struct {
	Keyword []string
	Type    *string
	Limit   int
	// 以下字段为历史兼容字段（三字段精简后新路径已不再使用），保留以支撑旧代码。
	Location       *string
	Characters     *string
	EmotionalMin   *float64
	EmotionalMax   *float64
	ImportanceMin  *float64
	TimeRangeStart *time.Time
	TimeRangeEnd   *time.Time
}

func (s *Storage) QueryMemories(ctx context.Context, q MemoryQuery) ([]models.Memory, error) {
	var memories []models.Memory
	db := s.sqliteDb.WithContext(ctx)

	// 基础查询：排除已遗忘的记忆
	query := db.Model(&models.Memory{}).Where("is_forgotten = ?", false)

	// 1. 关键词搜索（summary & content）→ 使用 LIKE
	if len(q.Keyword) > 0 {
		// 构建LIKE查询条件
		var conditions []string
		var args []interface{}

		for _, keyword := range q.Keyword {
			if strings.TrimSpace(keyword) != "" {
				likePattern := "%" + strings.TrimSpace(keyword) + "%"
				conditions = append(conditions, "(summary LIKE ? OR content LIKE ?)")
				args = append(args, likePattern, likePattern)
			}
		}

		if len(conditions) > 0 {
			// 使用 OR 连接多个关键词条件
			combinedCondition := strings.Join(conditions, " OR ")
			query = query.Where(combinedCondition, args...)
		}
	}

	// 2. 地点匹配（模糊包含）
	if q.Location != nil && *q.Location != "" {
		likePattern := "%" + *q.Location + "%"
		query = query.Where("location LIKE ?", likePattern)
	}

	// 3. 人物匹配（模糊包含）
	if q.Characters != nil && *q.Characters != "" {
		likePattern := "%" + *q.Characters + "%"
		query = query.Where("characters LIKE ?", likePattern)
	}

	// 4. 情感极性范围过滤
	if q.EmotionalMin != nil {
		query = query.Where("emotional_valence >= ?", *q.EmotionalMin)
	}
	if q.EmotionalMax != nil {
		query = query.Where("emotional_valence <= ?", *q.EmotionalMax)
	}

	// 5. 重要性阈值过滤
	if q.ImportanceMin != nil {
		query = query.Where("importance >= ?", *q.ImportanceMin)
	}

	// 6. 记忆类型匹配
	if q.Type != nil && *q.Type != "" {
		query = query.Where("type = ?", *q.Type)
	}

	// 7. 时间范围交集判断（关键逻辑）
	// 我们希望查询区间 [q.Start, q.End] 与记忆时间 [Start, End] 存在交集
	//
	// 区间相交条件：!(A_end < B_start || A_start > B_end)
	// 即：A 和 B 相交 ⇔ A_start <= B_end && B_start <= A_end
	//
	// 这里 A 是 memory 的时间区间，B 是查询区间
	if q.TimeRangeStart != nil || q.TimeRangeEnd != nil {
		// 构建动态 WHERE 条件
		conditions := ""

		// 如果查询有结束时间，则 memory 开始时间不能晚于它
		if q.TimeRangeEnd != nil {
			conditions += "time_rang_start <= ? OR time_rang_start IS NULL"
		}

		// 如果查询有开始时间，则 memory 结束时间不能早于它
		// 注意：如果 memory.TimeRangeEnd 为 NULL，我们认为它是“持续中”或“瞬时事件”，应视为满足条件
		if q.TimeRangeStart != nil {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += "(time_range_end >= ? OR time_range_end IS NULL)"
		}

		// 准备参数
		var args []interface{}
		if q.TimeRangeEnd != nil {
			args = append(args, *q.TimeRangeEnd)
		}
		if q.TimeRangeStart != nil {
			args = append(args, *q.TimeRangeStart)
		}

		if conditions != "" {
			query = query.Where(conditions, args...)
		}
	}

	// 8. 排序：优先级顺序
	query = query.
		Order("recall_count DESC").
		Order("COALESCE(time_rang_start, created_at) DESC")

	if q.Limit > 0 {
		query = query.Limit(q.Limit)
	}

	// 执行最终查询
	err := query.Find(&memories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}

	return memories, nil
}

// migrationLog 记录一次性数据迁移执行记录（确保迁移幂等）。
type migrationLog struct {
	Name       string    `gorm:"primaryKey;type:varchar(128)"`
	ExecutedAt time.Time `gorm:"autoCreateTime"`
}

func (migrationLog) TableName() string { return "memory_migration_log" }

const migrationLegacyFieldsToContent = "legacy_fields_to_content_v1"
const migrationCoreMemoryMetadata = "core_memory_metadata_v1"

// MigrateLegacyFieldsToContent 把旧的结构化字段（时间、地点、人物）拼到 content 末尾，
// 并把旧类型值映射到新类型（skill→information，plan/event→event，flow→event）。
// 幂等：通过 memory_migration_log 表记录已执行。
func (s *Storage) MigrateLegacyFieldsToContent(ctx context.Context) error {
	db := s.sqliteDb.WithContext(ctx)
	if err := db.AutoMigrate(&migrationLog{}); err != nil {
		return fmt.Errorf("auto migrate log table: %w", err)
	}

	var log migrationLog
	if err := db.Where("name = ?", migrationLegacyFieldsToContent).First(&log).Error; err == nil {
		return nil // 已迁移
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check migration log: %w", err)
	}

	var memories []models.Memory
	if err := db.Find(&memories).Error; err != nil {
		return fmt.Errorf("load memories for migration: %w", err)
	}

	for _, m := range memories {
		newType := mapLegacyType(string(m.Type))
		extras := buildLegacyExtras(m)
		newContent := m.Content
		if extras != "" {
			newContent = strings.TrimRight(m.Content, "\n")
			if newContent != "" {
				newContent += "\n\n"
			}
			newContent += extras
		}

		updates := map[string]any{}
		if m.Target == "" {
			updates["target"] = string(models.MemoryTargetUser)
		}
		if strings.TrimSpace(m.Source) == "" {
			updates["source"] = "legacy"
		}
		if m.CharCount == 0 && strings.TrimSpace(newContent) != "" {
			updates["char_count"] = utf8.RuneCountInString(strings.TrimSpace(newContent))
		}
		if newType != string(m.Type) {
			updates["type"] = newType
		}
		if newContent != m.Content {
			updates["content"] = newContent
		}
		if len(updates) == 0 {
			continue
		}
		if err := db.Model(&models.Memory{}).Where("id = ?", m.ID).Updates(updates).Error; err != nil {
			return fmt.Errorf("update memory %d: %w", m.ID, err)
		}
	}

	if err := db.Create(&migrationLog{Name: migrationLegacyFieldsToContent}).Error; err != nil {
		return fmt.Errorf("record migration log: %w", err)
	}
	return nil
}

func (s *Storage) MigrateCoreMemoryMetadata(ctx context.Context) error {
	db := s.sqliteDb.WithContext(ctx)
	if err := db.AutoMigrate(&migrationLog{}); err != nil {
		return fmt.Errorf("auto migrate log table: %w", err)
	}
	var log migrationLog
	if err := db.Where("name = ?", migrationCoreMemoryMetadata).First(&log).Error; err == nil {
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check migration log: %w", err)
	}
	if err := db.Model(&models.Memory{}).Where("target = '' OR target IS NULL").Update("target", string(models.MemoryTargetUser)).Error; err != nil {
		return fmt.Errorf("backfill target: %w", err)
	}
	if err := db.Model(&models.Memory{}).Where("source = '' OR source IS NULL").Update("source", "legacy").Error; err != nil {
		return fmt.Errorf("backfill source: %w", err)
	}
	var memories []models.Memory
	if err := db.Find(&memories).Error; err != nil {
		return fmt.Errorf("load memories for metadata migration: %w", err)
	}
	for _, m := range memories {
		count := utf8.RuneCountInString(strings.TrimSpace(m.Content))
		if m.CharCount == count {
			continue
		}
		if err := db.Model(&models.Memory{}).Where("id = ?", m.ID).Update("char_count", count).Error; err != nil {
			return fmt.Errorf("update char_count for memory %d: %w", m.ID, err)
		}
	}
	if err := db.Create(&migrationLog{Name: migrationCoreMemoryMetadata}).Error; err != nil {
		return fmt.Errorf("record migration log: %w", err)
	}
	return nil
}

// mapLegacyType 映射旧类型值到新的 fact/information/event 语义。
func mapLegacyType(old string) string {
	switch strings.TrimSpace(old) {
	case "skill":
		return "information"
	case "event", "plan", "plan ", "flow":
		return "event"
	default:
		return old
	}
}

// buildLegacyExtras 把旧结构化字段渲染成要追加到 content 末尾的自然语言片段。
func buildLegacyExtras(m models.Memory) string {
	var parts []string
	if m.TimeRangStart != nil || m.TimeRangeEnd != nil {
		var timeStr string
		if m.TimeRangStart != nil && m.TimeRangeEnd != nil {
			s := m.TimeRangStart.Format("2006-01-02")
			e := m.TimeRangeEnd.Format("2006-01-02")
			if s == e {
				timeStr = s
			} else {
				timeStr = s + " 至 " + e
			}
		} else if m.TimeRangStart != nil {
			timeStr = m.TimeRangStart.Format("2006-01-02")
		} else {
			timeStr = m.TimeRangeEnd.Format("2006-01-02")
		}
		parts = append(parts, "[时间] "+timeStr)
	}
	if m.Location != nil && strings.TrimSpace(*m.Location) != "" {
		parts = append(parts, "[地点] "+strings.TrimSpace(*m.Location))
	}
	if m.Characters != nil && strings.TrimSpace(*m.Characters) != "" {
		parts = append(parts, "[人物] "+strings.TrimSpace(*m.Characters))
	}
	return strings.Join(parts, "\n")
}

// UpdateMemory 更新指定的记忆
func (s *Storage) UpdateMemory(ctx context.Context, id uint, memory models.Memory) error {
	// 首先检查记忆是否存在
	var existingMemory models.Memory
	result := s.sqliteDb.WithContext(ctx).First(&existingMemory, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("记忆 ID %d 不存在", id)
		}
		return fmt.Errorf("查询记忆失败: %w", result.Error)
	}

	// 仅更新三字段中非零值的部分（保持 Agent 增量补全语义）
	updateData := make(map[string]any)

	if memory.Summary != "" {
		updateData["summary"] = memory.Summary
	}
	if memory.Content != "" {
		updateData["content"] = memory.Content
	}
	if memory.Type != "" {
		updateData["type"] = memory.Type
	}
	if memory.Target != "" {
		updateData["target"] = memory.Target
	}
	if memory.Source != "" {
		updateData["source"] = memory.Source
	}
	if memory.Content != "" {
		updateData["char_count"] = utf8.RuneCountInString(memory.Content)
	}

	updateData["updated_at"] = time.Now()

	// 执行更新
	result = s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("更新记忆失败: %w", result.Error)
	}

	// 如果更新了 content 或 summary，不需要更新 FTS 索引（已移除FTS）
	// 直接使用LIKE查询，无需FTS索引维护

	return nil
}

// ReplaceMemoryEditableFields 覆盖更新记忆三字段（summary/content/type）。
func (s *Storage) ReplaceMemoryEditableFields(ctx context.Context, id uint, memory models.Memory) error {
	var existingMemory models.Memory
	result := s.sqliteDb.WithContext(ctx).First(&existingMemory, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("记忆 ID %d 不存在", id)
		}
		return fmt.Errorf("查询记忆失败: %w", result.Error)
	}

	updateData := map[string]any{
		"summary":    memory.Summary,
		"content":    memory.Content,
		"type":       memory.Type,
		"target":     memory.Target,
		"source":     memory.Source,
		"char_count": utf8.RuneCountInString(strings.TrimSpace(memory.Content)),
		"updated_at": time.Now(),
	}
	if updateData["target"] == "" {
		updateData["target"] = existingMemory.Target
	}
	if updateData["source"] == "" {
		updateData["source"] = existingMemory.Source
	}

	result = s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("更新记忆失败: %w", result.Error)
	}
	return nil
}
