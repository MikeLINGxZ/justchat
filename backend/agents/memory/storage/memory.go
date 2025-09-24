package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gorm.io/gorm"
)

func (s *Storage) WriterMemory(ctx context.Context, memory models.Memory) (uint, error) {
	result := s.sqliteDb.WithContext(ctx).Create(&memory)
	return memory.ID, result.Error
}

func (s *Storage) ReadMemory(ctx context.Context, keyword string, startAt, endAt *time.Time) ([]models.Memory, error) {
	var memories []models.Memory

	db := s.sqliteDb.WithContext(ctx)

	// 构建基础查询
	query := db.Model(&models.Memory{}).Where("is_forbidden = ?", false) // 排除已遗忘的记忆（可选）

	// 1. 处理关键词搜索：使用 FTS5 搜索 summary 和 content
	if keyword != "" {
		// 使用 FTS5 匹配
		ftsCondition := "memory_fts MATCH ?"
		ftsArgs := sanitizeFTSQuery([]string{keyword}) // 防止特殊字符导致语法错误

		// 子查询获取匹配的 rowid
		subQuery := db.Table("memory_fts").Select("rowid").Where(ftsCondition, ftsArgs)

		// 主查询限制 ID 在 FTS 结果中
		query = query.Where("id IN (?)", subQuery)
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

// sanitizeFTSQuery 清理并格式化用于 FTS5 查询的关键词
func sanitizeFTSQuery(keywords []string) string {
	if len(keywords) == 0 {
		return ""
	}

	// 使用 \p{Han} 匹配所有汉字，支持中文；使用原始字符串 `` `...` `` 没问题
	re := regexp.MustCompile(`[^a-zA-Z0-9\p{Han}\s\-_*]+`)

	var allWords []string
	for _, keyword := range keywords {
		cleaned := re.ReplaceAllString(keyword, " ")
		words := strings.Fields(cleaned)
		allWords = append(allWords, words...)
	}

	if len(allWords) == 0 {
		return ""
	}

	return strings.Join(allWords, " OR ")
}

type MemoryQuery struct {
	Keyword        []string
	Location       *string
	Characters     *string
	EmotionalMin   *float64
	EmotionalMax   *float64
	ImportanceMin  *float64
	Type           *string
	TimeRangeStart *time.Time
	TimeRangeEnd   *time.Time
}

func (s *Storage) QueryMemories(ctx context.Context, q MemoryQuery) ([]models.Memory, error) {
	var memories []models.Memory
	db := s.sqliteDb.WithContext(ctx)

	// 基础查询：排除已遗忘的记忆
	query := db.Model(&models.Memory{}).Where("is_forgotten = ?", false)

	// 1. 关键词搜索（summary & content）→ 使用 FTS5
	if len(q.Keyword) > 0 {
		ftsArgs := sanitizeFTSQuery(q.Keyword)
		if ftsArgs != "" {
			subQuery := db.Table("memory_fts").Select("rowid").Where("memory_fts MATCH ?", ftsArgs)
			query = query.Where("id IN (?)", subQuery)
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
		Order("importance DESC").                           // 首要：按重要性降序
		Order("recall_count DESC").                         // 其次：常被回忆的优先
		Order("COALESCE(time_rang_start, created_at) DESC") // 最近发生的靠前

	// 执行最终查询
	err := query.Find(&memories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}

	return memories, nil
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

	// 更新记忆，只更新非零值字段
	updateData := make(map[string]interface{})

	if memory.Summary != "" {
		updateData["summary"] = memory.Summary
	}
	if memory.Content != "" {
		updateData["content"] = memory.Content
	}
	if memory.Type != "" {
		updateData["type"] = memory.Type
	}
	if memory.TimeRangStart != nil {
		updateData["time_rang_start"] = memory.TimeRangStart
	}
	if memory.TimeRangeEnd != nil {
		updateData["time_range_end"] = memory.TimeRangeEnd
	}
	if memory.Location != nil {
		updateData["location"] = memory.Location
	}
	if memory.Characters != nil {
		updateData["characters"] = memory.Characters
	}
	if memory.Context != nil {
		updateData["context"] = memory.Context
	}
	if memory.Importance != 0 {
		updateData["importance"] = memory.Importance
	}
	if memory.EmotionalValence != 0 {
		updateData["emotional_valence"] = memory.EmotionalValence
	}

	// 更新修改时间
	updateData["updated_at"] = time.Now()

	// 执行更新
	result = s.sqliteDb.WithContext(ctx).Model(&models.Memory{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("更新记忆失败: %w", result.Error)
	}

	// 如果更新了 content 或 summary，需要更新 FTS 索引
	if memory.Summary != "" || memory.Content != "" {
		// 重新构建 FTS 索引（这里简化处理，实际可能需要更复杂的逻辑）
		err := s.sqliteDb.Exec((&models.Memory{}).Fts()).Error
		if err != nil {
			// FTS 更新失败不应该影响主更新操作，只记录错误
			fmt.Printf("警告：FTS 索引更新失败: %v\n", err)
		}
	}

	return nil
}
