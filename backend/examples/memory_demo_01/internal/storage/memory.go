package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/models"
)

func (s *Storage) WriterMemory(ctx context.Context, memory models.Memory) error {
	result := s.sqliteDb.WithContext(ctx).Create(&memory)
	return result.Error
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
