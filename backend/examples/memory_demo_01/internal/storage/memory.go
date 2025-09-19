package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/models"
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
		ftsArgs := sanitizeFTSQuery(keyword) // 防止特殊字符导致语法错误

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
func sanitizeFTSQuery(keyword string) string {
	// 移除危险字符，保留字母数字和基本符号
	re := regexp.MustCompile(`[^a-zA-Z0-9\u4e00-\u9fa5\s\-_\*]+`)
	cleaned := re.ReplaceAllString(keyword, " ")

	// 分词后加上双引号进行短语匹配？或者用 OR 联合
	words := strings.Fields(cleaned)
	if len(words) == 0 {
		return ""
	}

	// 使用 NEAR 或 OR 取决于需求；这里用 OR 实现宽松匹配
	return strings.Join(words, " OR ")
}

type MemoryQuery struct {
	Keyword        string
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
	query := db.Model(&models.Memory{}).Where("is_forget = ?", false)

	// 1. 关键词搜索（summary & content）→ 使用 FTS5
	if q.Keyword != "" {
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
