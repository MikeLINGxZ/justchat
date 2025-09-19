package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	istorage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
)

type WriteMemoryToolRequest struct {
	Title            string   `json:"title" jsonschema:"title=记忆标题,description=记忆或事件的简短标题，例如“童年第一次骑自行车”，required=true"`
	Content          string   `json:"content" jsonschema:"title=记忆内容,description=详细的记忆描述，可以包括人物、地点、感受等，required=true"`
	TimeRangeStart   *string  `json:"time_range_start,omitempty" jsonschema:"title=开始日期,description=记忆发生的起始日期（格式：YYYY-MM-DD 或 YYYY-MM-DD HH:MM:SS），可为空；若只有单日事件，请只填此项"`
	TimeRangeEnd     *string  `json:"time_range_end,omitempty" jsonschema:"title=结束日期,description=记忆结束的日期（格式同上），适用于持续性事件，如旅行、恋情等，可为空"`
	Location         *string  `json:"location,omitempty" jsonschema:"title=发生地点,description=该记忆发生的地点，多个地点可用逗号分隔，例如“北京,朝阳公园”"`
	Characters       *string  `json:"characters,omitempty" jsonschema:"title=相关人物,description=涉及的人物名称，多个用逗号分隔，例如“父亲,小明”"`
	MemoryType       string   `json:"memory_type" jsonschema:"title=记忆类型,description=这个记忆的类型，可以不填，如果填则需要在以下范围内：skill、event、plan"`
	EmotionalValence *float64 `json:"emotional_valence,omitempty" jsonschema:"title=情感极性,description=情绪倾向评分，范围 -1.0（极度负面）到 +1.0（极度正面），默认为 0.0，可为空"`
	Importance       *float64 `json:"importance,omitempty" jsonschema:"title=重要性评分,description=此记忆的重要性，范围 0.0 ~ 1.0，默认为 0.5，可为空"`
	Context          *string  `json:"context,omitempty" jsonschema:"title=上下文元数据,description=额外的 JSON 格式元数据，例如设备来源、天气、心情标签等，可为空，格式必须是合法 JSON"`
}

type WriteMemoryToolResponse struct {
	Success  bool   `json:"success" jsonschema:"title=是否成功,description=表示记忆是否成功写入"`
	Message  string `json:"message,omitempty" jsonschema:"title=提示信息,description=成功或失败的详细说明，例如 '记忆已保存' 或 '数据库错误'"`
	MemoryID uint   `json:"memory_id,omitempty" jsonschema:"title=记忆ID,description=成功时返回新创建的记忆唯一标识符"`
}

func NewWriteMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"write_memory",
		"Record a personal memory or life event with detailed metadata including time, location, people, emotion, and importance",
		func(ctx context.Context, in WriteMemoryToolRequest) (output WriteMemoryToolResponse, err error) {
			response := WriteMemoryToolResponse{Success: false}

			// 基础验证
			if in.Title == "" {
				response.Message = "无法记录记忆：缺少标题（title）"
				return response, nil
			}
			if in.Content == "" {
				response.Message = "无法记录记忆：缺少内容（content）"
				return response, nil
			}

			// 解析时间范围
			var startTime, endTime *time.Time
			if in.TimeRangeStart != nil {
				t, err := parseTime(*in.TimeRangeStart)
				if err != nil {
					response.Message = "无法解析开始时间：" + err.Error()
					return response, nil
				}
				startTime = &t
			}
			if in.TimeRangeEnd != nil {
				t, err := parseTime(*in.TimeRangeEnd)
				if err != nil {
					response.Message = "无法解析结束时间：" + err.Error()
					return response, nil
				}
				endTime = &t
			}

			// 构造 Memory 对象的部分字段
			memoryData := models.Memory{
				Summary:          in.Title,
				Content:          in.Content,
				Type:             "", // 可通过 AI 自动推断类型，或留空由系统处理
				TimeRangStart:    startTime,
				TimeRangeEnd:     endTime,
				Location:         in.Location,
				Characters:       in.Characters,
				Context:          in.Context,
				Importance:       0.5,
				EmotionalValence: 0.0,
				IsForgotten:      false,
				RecallCount:      0,
			}

			// 设置可选浮点值（带默认）
			if in.Importance != nil {
				if *in.Importance < 0.0 || *in.Importance > 1.0 {
					response.Message = "重要性评分必须在 0.0 ~ 1.0 之间"
					return response, nil
				}
				memoryData.Importance = *in.Importance
			}
			if in.EmotionalValence != nil {
				if *in.EmotionalValence < -1.0 || *in.EmotionalValence > 1.0 {
					response.Message = "情感极性必须在 -1.0 ~ +1.0 之间"
					return response, nil
				}
				memoryData.EmotionalValence = *in.EmotionalValence
			}

			// 调用 Storage 写入记忆（假设方法已更新）
			id, err := storage.WriterMemory(ctx, memoryData)
			if err != nil {
				response.Message = "保存失败：" + err.Error()
				return response, nil
			}

			// 成功响应
			response.Success = true
			response.Message = "记忆已成功记录"
			response.MemoryID = id
			return response, nil
		})
}

// 辅助函数：尝试多种时间格式解析
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		time.DateOnly,         // YYYY-MM-DD
		"2006-01-02 15:04:05", // YYYY-MM-DD HH:MM:SS
		time.RFC3339,          // ISO8601
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("不支持的时间格式: %s", s)
}

type ReadMemoryToolRequest struct {
	Keyword        *string  `json:"keyword,omitempty" jsonschema:"title=关键词,description=在标题或内容中进行模糊匹配的词语，例如‘杭州旅行’"`
	Location       *string  `json:"location,omitempty" jsonschema:"title=发生地点,description=精确或模糊匹配记忆发生的地点，例如‘北京’或‘西湖’"`
	Characters     *string  `json:"characters,omitempty" jsonschema:"title=相关人物,description=包含某个人物的记忆，多个名字可用逗号分隔，系统会匹配任一出现的人物"`
	EmotionalMin   *float64 `json:"emotional_min,omitempty" jsonschema:"title=最小情感值,description=只返回情感极性大于等于该值的记忆，范围 -1.0 ~ +1.0，例如 0.5 表示只看正面情绪"`
	EmotionalMax   *float64 `json:"emotional_max,omitempty" jsonschema:"title=最大情感值,description=只返回情感极性小于等于该值的记忆，例如 -0.3 可用于查找负面经历"`
	ImportanceMin  *float64 `json:"importance_min,omitempty" jsonschema:"title=最低重要性,description=只返回重要性评分高于此值的记忆，例如 0.7 表示只关注高价值记忆"`
	Type           *string  `json:"type,omitempty" jsonschema:"title=记忆类型,description=按预定义类型过滤，例如 'childhood', 'travel', 'work' 等"`
	TimeRangeStart *string  `json:"time_range_start,omitempty" jsonschema:"title=起始时间,description=格式 YYYY-MM-DD 或 YYYY-MM-DD HH:MM:SS，包含当天/时刻"`
	TimeRangeEnd   *string  `json:"time_range_end,omitempty" jsonschema:"title=结束时间,description=同上，查询区间为闭区间 [start, end]"`
}

type ReadMemoryToolResponse struct {
	Success  bool                `json:"success" jsonschema:"title=是否成功,description=查询是否成功执行"`
	Message  string              `json:"message,omitempty" jsonschema:"title=提示信息,description=错误或状态说明，如'未找到匹配记忆'"`
	Memories []ReadMemorySummary `json:"memories" jsonschema:"title=记忆列表,description=匹配的回忆摘要列表"`
	Total    int                 `json:"total" jsonschema:"title=总数,description=返回的记忆条数"`
}

type ReadMemorySummary struct {
	ID               uint       `json:"id" jsonschema:"title=记忆ID,description=该记忆的唯一编号"`
	Title            string     `json:"title" jsonschema:"title=标题,description=记忆的简短标题"`
	Excerpt          string     `json:"excerpt" jsonschema:"title=内容摘要,description=内容前100字符作为预览"`
	TimeRangeStart   *time.Time `json:"time_range_start,omitempty" jsonschema:"title=开始时间,description=事件发生起始时间"`
	TimeRangeEnd     *time.Time `json:"time_range_end,omitempty" jsonschema:"title=结束时间,description=事件结束时间（如有）"`
	Location         *string    `json:"location,omitempty" jsonschema:"title=地点,description=记忆发生地"`
	Characters       *string    `json:"characters,omitempty" jsonschema:"title=人物,description=涉及的人物"`
	EmotionalValence float64    `json:"emotional_valence" jsonschema:"title=情感极性,description=情绪倾向：-1.0(悲伤)~+1.0(快乐)"`
	Importance       float64    `json:"importance" jsonschema:"title=重要性,description=用户标记的重要性评分"`
	CreatedAt        time.Time  `json:"created_at" jsonschema:"title=记录时间,description=该记忆被录入系统的日期时间"`
}

func NewReadMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"read_memory",
		"提供多维度复合查询能力，支持基于关键词模糊匹配、时空范围限定（YYYY-MM-DD HH:MM:SS）、人物关联、地理位置、情感倾向（emotional_valence）与重要性权重（importance）的记忆检索。适用于构建上下文连贯性、追溯用户决策路径、识别情绪变化趋势及触发个性化回应。",
		func(ctx context.Context, in *ReadMemoryToolRequest) (output *ReadMemoryToolResponse, err error) {
			response := &ReadMemoryToolResponse{
				Success:  false,
				Memories: nil,
				Total:    0,
			}

			// 如果请求为空，则视为获取所有记忆（可限制数量）
			if in == nil {
				response.Success = true
				response.Message = "未提供查询条件，返回空结果集"
				response.Memories = []ReadMemorySummary{}
				return response, nil
			}

			// 解析时间范围
			var startAt, endAt *time.Time
			if in.TimeRangeStart != nil && *in.TimeRangeStart != "" {
				t, err := parseTime(*in.TimeRangeStart)
				if err != nil {
					response.Message = "无法解析起始时间：" + err.Error()
					return response, nil
				}
				startAt = &t
			}
			if in.TimeRangeEnd != nil && *in.TimeRangeEnd != "" {
				t, err := parseTime(*in.TimeRangeEnd)
				if err != nil {
					response.Message = "无法解析结束时间：" + err.Error()
					return response, nil
				}
				endAt = &t
			}
			if startAt != nil && endAt != nil && startAt.After(*endAt) {
				response.Message = "起始时间不能晚于结束时间"
				return response, nil
			}

			// 校验情感值范围
			var emotionalMin, emotionalMax *float64
			if in.EmotionalMin != nil {
				v := *in.EmotionalMin
				if v < -1.0 || v > 1.0 {
					response.Message = "情感最小值必须在 [-1.0, 1.0] 范围内"
					return response, nil
				}
				emotionalMin = &v
			}
			if in.EmotionalMax != nil {
				v := *in.EmotionalMax
				if v < -1.0 || v > 1.0 {
					response.Message = "情感最大值必须在 [-1.0, 1.0] 范围内"
					return response, nil
				}
				emotionalMax = &v
			}
			if emotionalMin != nil && emotionalMax != nil && *emotionalMin > *emotionalMax {
				response.Message = "情感最小值不能大于最大值"
				return response, nil
			}

			// 校验重要性
			var importanceMin *float64
			if in.ImportanceMin != nil {
				v := *in.ImportanceMin
				if v < 0.0 || v > 1.0 {
					response.Message = "重要性阈值必须在 [0.0, 1.0] 范围内"
					return response, nil
				}
				importanceMin = &v
			}

			// 构建查询参数
			query := istorage.MemoryQuery{
				Keyword:        derefOrEmpty(in.Keyword),
				Location:       in.Location,
				Characters:     in.Characters,
				EmotionalMin:   emotionalMin,
				EmotionalMax:   emotionalMax,
				ImportanceMin:  importanceMin,
				Type:           in.Type,
				TimeRangeStart: startAt,
				TimeRangeEnd:   endAt,
			}

			// 查询存储层
			memories, err := storage.QueryMemories(ctx, query)
			if err != nil {
				response.Message = "查询失败：" + err.Error()
				return response, nil
			}

			// 转换为响应摘要
			summaries := make([]ReadMemorySummary, 0, len(memories))
			for _, m := range memories {
				excerpt := trimContent(m.Content, 100)
				title := m.Summary
				if title == "" {
					title = "(无标题)"
				}

				summaries = append(summaries, ReadMemorySummary{
					ID:               m.ID,
					Title:            title,
					Excerpt:          excerpt,
					TimeRangeStart:   m.TimeRangStart,
					TimeRangeEnd:     m.TimeRangeEnd,
					Location:         m.Location,
					Characters:       m.Characters,
					EmotionalValence: m.EmotionalValence,
					Importance:       m.Importance,
					CreatedAt:        m.CreatedAt,
				})
			}

			// 返回结果
			response.Success = true
			response.Memories = summaries
			response.Total = len(summaries)

			switch {
			case len(summaries) == 0:
				response.Message = "没有找到符合条件的记忆记录"
			case len(summaries) == 1:
				response.Message = "找到 1 条相关记忆"
			default:
				response.Message = fmt.Sprintf("成功找到 %d 条相关记忆", len(summaries))
			}

			return response, nil
		},
	)
}

func trimContent(content string, maxLen int) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return "(无内容)"
	}
	if len(content) > maxLen {
		return content[:maxLen-3] + "..."
	}
	return content
}

func derefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
