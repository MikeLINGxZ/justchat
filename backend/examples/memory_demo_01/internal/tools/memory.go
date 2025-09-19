package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/storage"
)

type WriteMemoryToolRequest struct {
	Title   string  `json:"title" jsonschema:"title=记忆标题;description=记忆或事件的简短标题，例如“童年第一次骑自行车”；required=true"`
	Content string  `json:"content" jsonschema:"title=记忆内容;description=详细的记忆描述，可以包括人物、地点、感受等；required=true"`
	Date    *string `json:"date" jsonschema:"title=发生日期;description=该记忆实际发生的日期（格式：YYYY-MM-DD），如果不确定可为空；required=false"`
}

type WriteMemoryToolResponse struct {
	Success  bool   `json:"success" jsonschema:"title=是否成功;description=表示记忆是否成功写入"`
	Message  string `json:"message,omitempty" jsonschema:"title=提示信息;description=成功或失败的详细说明，例如 '记忆已保存' 或 '数据库错误'"`
	MemoryID uint   `json:"memory_id,omitempty" jsonschema:"title=记忆ID;description=成功时返回新创建的记忆唯一标识符"`
}

func NewWriteMemoryTool(storage *storage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"write_memory",
		"Record a personal memory or life event with title, content, and optional date",
		func(ctx context.Context, in *WriteMemoryToolRequest) (output *WriteMemoryToolResponse, err error) {
			// 初始化响应
			response := &WriteMemoryToolResponse{
				Success: false,
			}

			// 处理 nil 请求
			if in == nil {
				response.Message = "无法记录记忆：未提供任何数据"
				return response, nil
			}

			// 基础验证
			if in.Title == "" {
				response.Message = "无法记录记忆：缺少标题（title）"
				return response, nil
			}
			if in.Content == "" {
				response.Message = "无法记录记忆：缺少内容（content）"
				return response, nil
			}

			var date *time.Time
			if in.Date != nil {
				parseDate, err := time.Parse(time.DateOnly, *in.Date)
				if err == nil {
					date = &parseDate
				}
			}

			// 调用 Storage 写入记忆
			err = storage.WriterMemory(ctx, in.Title, in.Content, date)
			if err != nil {
				response.Message = "保存失败：" + err.Error()
				return response, nil // 返回 response 给 AI，而不是抛出 error
			}

			// 成功响应
			response.Success = true
			response.Message = "记忆已成功记录"

			return response, nil
		})
}

type ReadMemoryToolRequest struct {
	Keyword string `json:"keyword,omitempty" jsonschema:"title=关键词;description=用于模糊搜索记忆的标题或内容，例如‘杭州旅行’；required=false"`
	StartAt string `json:"start_at,omitempty" jsonschema:"title=起始日期;description=只返回此日期之后（含）发生的记忆，格式 YYYY-MM-DD；required=false"`
	EndAt   string `json:"end_at,omitempty" jsonschema:"title=结束日期;description=只返回此日期之前（含）发生的记忆，格式 YYYY-MM-DD；required=false"`
}

type ReadMemoryToolResponse struct {
	Success  bool                `json:"success" jsonschema:"title=是否成功;description=查询是否成功执行"`
	Message  string              `json:"message,omitempty" jsonschema:"title=提示信息;description=错误或状态说明，如'未找到匹配记忆'"`
	Memories []ReadMemorySummary `json:"memories" jsonschema:"title=记忆列表;description=匹配的回忆摘要列表"`
	Total    int                 `json:"total" jsonschema:"title=总数;description=返回的记忆条数"`
}

func NewReadMemoryTool(storage *storage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"read_memory",
		"根据语义关键词、时间范围、记忆类型或情感特征，检索用户过往的记忆片段。用于帮助AI回忆共同经历，实现个性化共情对话。",
		func(ctx context.Context, in *ReadMemoryToolRequest) (output *ReadMemoryToolResponse, err error) {
			// 初始化响应
			response := &ReadMemoryToolResponse{
				Success:  false,
				Memories: []ReadMemorySummary{},
				Total:    0,
			}

			// 处理 nil 请求：视为无条件查询
			if in == nil {
				response.Message = "未提供查询条件，默认返回空结果"
				response.Success = true
				return response, nil
			}

			// 参数验证和清理
			keyword := strings.TrimSpace(in.Keyword)
			if keyword == "" {
				keyword = "" // 空字符串表示查询所有记忆
			}

			// 解析和验证日期
			var startAt, endAt *time.Time
			if in.StartAt != "" && strings.TrimSpace(in.StartAt) != "" {
				if parsedDate, err := time.Parse(time.DateOnly, strings.TrimSpace(in.StartAt)); err == nil {
					startAt = &parsedDate
				} else {
					response.Message = fmt.Sprintf("起始日期格式错误: %s，应为 YYYY-MM-DD 格式", in.StartAt)
					return response, nil
				}
			}

			if in.EndAt != "" && strings.TrimSpace(in.EndAt) != "" {
				if parsedDate, err := time.Parse(time.DateOnly, strings.TrimSpace(in.EndAt)); err == nil {
					endAt = &parsedDate
				} else {
					response.Message = fmt.Sprintf("结束日期格式错误: %s，应为 YYYY-MM-DD 格式", in.EndAt)
					return response, nil
				}
			}

			// 验证日期范围逻辑
			if startAt != nil && endAt != nil && startAt.After(*endAt) {
				response.Message = "起始日期不能晚于结束日期"
				return response, nil
			}

			// 调用 Storage 层查询数据
			memories, err := storage.ReadMemory(ctx, keyword, startAt, endAt)
			if err != nil {
				response.Message = "查询记忆失败: " + err.Error()
				return response, nil // 返回给 AI 友好的结构，不抛出 error
			}

			// 转换为摘要列表
			summaries := make([]ReadMemorySummary, 0, len(memories))
			for _, m := range memories {
				excerpt := strings.TrimSpace(m.Content)
				if len(excerpt) > 100 {
					excerpt = excerpt[:97] + "..."
				} else if excerpt == "" {
					excerpt = "(无内容)"
				}

				// 确保标题不为空
				title := strings.TrimSpace(m.Title)
				if title == "" {
					title = "(无标题)"
				}

				summaries = append(summaries, ReadMemorySummary{
					ID:           m.ID,
					Title:        title,
					Excerpt:      excerpt,
					DateOccurred: m.DateOccurred,
					CreatedAt:    m.CreatedAt,
				})
			}

			// 填充成功响应
			response.Success = true
			response.Memories = summaries
			response.Total = len(summaries)

			// 生成更友好的消息
			if len(summaries) == 0 {
				if keyword != "" {
					response.Message = fmt.Sprintf("没有找到包含'%s'的记忆记录", keyword)
				} else {
					response.Message = "没有找到匹配的记忆记录"
				}
			} else {
				response.Message = fmt.Sprintf("成功找到 %d 条相关记忆", len(summaries))
			}

			return response, nil
		},
	)
}

type ReadMemorySummary struct {
	ID           uint       `json:"id" jsonschema:"title=记忆ID;description=该记忆的唯一编号"`
	Title        string     `json:"title" jsonschema:"title=标题;description=记忆的简短标题"`
	Excerpt      string     `json:"excerpt" jsonschema:"title=摘要;description=内容前50个字符作为预览"`
	DateOccurred *time.Time `json:"date_occurred,omitempty" jsonschema:"title=发生日期;description=事件实际发生的日期，可能为空"`
	CreatedAt    time.Time  `json:"created_at" jsonschema:"title=创建时间;description=该记忆被记录到系统的时间"`
}
