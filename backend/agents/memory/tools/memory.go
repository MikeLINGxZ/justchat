package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	istorage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
)

type CoreMemoryToolRequest struct {
	Action  string `json:"action" jsonschema:"title=操作,description=add / replace / remove 之一。required=true"`
	Target  string `json:"target" jsonschema:"title=目标,description=user 表示用户画像；memory 表示助手/环境笔记。required=true"`
	Content string `json:"content,omitempty" jsonschema:"title=内容,description=add/replace 时的新记忆内容，要求精炼、信息密度高"`
	OldText string `json:"old_text,omitempty" jsonschema:"title=旧文本片段,description=replace/remove 时用于唯一定位旧条目的短子串"`
}

type CoreMemoryToolResponse struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message,omitempty"`
	MemoryID       uint     `json:"memory_id,omitempty"`
	Usage          string   `json:"usage,omitempty"`
	CurrentEntries []string `json:"current_entries,omitempty"`
}

func NewCoreMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"memory",
		"Hermes-style bounded core memory tool. Use add/replace/remove on target=user or target=memory. There is no read action: core memory is already injected into the system prompt.",
		func(ctx context.Context, in CoreMemoryToolRequest) (CoreMemoryToolResponse, error) {
			response := CoreMemoryToolResponse{Success: false}
			target, err := istorage.NormalizeMemoryTarget(in.Target)
			if err != nil {
				response.Message = err.Error()
				return response, nil
			}

			var result *istorage.CoreMemoryMutationResult
			switch strings.ToLower(strings.TrimSpace(in.Action)) {
			case "add":
				result, err = storage.AddCoreMemory(ctx, target, in.Content, "agent")
			case "replace":
				result, err = storage.ReplaceCoreMemory(ctx, target, in.OldText, in.Content)
			case "remove":
				result, err = storage.RemoveCoreMemory(ctx, target, in.OldText)
			default:
				response.Message = "unsupported action: use add, replace, or remove"
				return response, nil
			}
			if err != nil {
				response.Message = err.Error()
				return response, nil
			}
			if result == nil {
				response.Message = "memory operation returned no result"
				return response, nil
			}
			response.MemoryID = result.MemoryID
			response.Message = result.Message
			response.Usage = result.Usage
			response.CurrentEntries = result.CurrentEntries
			response.Success = !strings.Contains(strings.ToLower(result.Message), "exceed") &&
				!strings.Contains(strings.ToLower(result.Message), "matched") &&
				!strings.Contains(strings.ToLower(result.Message), "no matching")
			return response, nil
		})
}

const memoryTypeDesc = "记忆类型，必须是以下之一：" +
	"fact（事实，长期不变的客观属性）、" +
	"information（信息，可变的偏好/状态/习惯）、" +
	"event（事件，带时间锚点的具体事件或计划）。" +
	"无法判断时可留空。"

// ---- write_memory ----

type WriteMemoryToolRequest struct {
	Title      string `json:"title" jsonschema:"title=标题,description=记忆的简短标题，例如'2026-04-21 京东快递员上门换货'。required=true"`
	Content    string `json:"content" jsonschema:"title=内容,description=完整记忆内容。必须把时间（写成绝对日期 YYYY-MM-DD）、地点、人物、情绪等所有相关信息以自然语言融入这段文本，不要让信息散落在结构化字段里。required=true"`
	MemoryType string `json:"memory_type,omitempty" jsonschema:"title=类型,description=fact / information / event 之一，可留空"`
}

type WriteMemoryToolResponse struct {
	Success  bool   `json:"success" jsonschema:"title=是否成功"`
	Message  string `json:"message,omitempty" jsonschema:"title=提示信息"`
	MemoryID uint   `json:"memory_id,omitempty" jsonschema:"title=记忆ID"`
}

func NewWriteMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"write_memory",
		"记录一条新的长期记忆。仅有 title、content、memory_type 三个字段；时间/地点/人物/情绪等信息必须直接写入 content 文本（时间必须为绝对日期）。"+memoryTypeDesc,
		func(ctx context.Context, in WriteMemoryToolRequest) (WriteMemoryToolResponse, error) {
			response := WriteMemoryToolResponse{Success: false}

			if strings.TrimSpace(in.Title) == "" {
				response.Message = "无法记录记忆：缺少标题"
				return response, nil
			}
			if strings.TrimSpace(in.Content) == "" {
				response.Message = "无法记录记忆：缺少内容"
				return response, nil
			}

			memType, err := normalizeMemoryType(in.MemoryType)
			if err != nil {
				response.Message = err.Error()
				return response, nil
			}

			memoryData := models.Memory{
				Summary:     strings.TrimSpace(in.Title),
				Content:     strings.TrimSpace(in.Content),
				Type:        memType,
				IsForgotten: false,
				RecallCount: 0,
			}

			id, err := storage.WriterMemory(ctx, memoryData)
			if err != nil {
				response.Message = "保存失败：" + err.Error()
				return response, nil
			}

			response.Success = true
			response.Message = "记忆已成功记录"
			response.MemoryID = id
			return response, nil
		})
}

// ---- read_memory ----

type ReadMemoryToolRequest struct {
	Keyword    *string `json:"keyword,omitempty" jsonschema:"title=关键词,description=用于在标题或内容中模糊匹配，多个关键词使用英文逗号分隔，例如：'杭州,旅行'"`
	MemoryType *string `json:"memory_type,omitempty" jsonschema:"title=类型过滤,description=仅返回指定类型的记忆（fact / information / event），可留空"`
	Limit      *int    `json:"limit,omitempty" jsonschema:"title=返回上限,description=最多返回多少条，默认 20，最大 100"`
}

type ReadMemoryToolResponse struct {
	Success  bool                `json:"success" jsonschema:"title=是否成功"`
	Message  string              `json:"message,omitempty" jsonschema:"title=提示信息"`
	Memories []ReadMemorySummary `json:"memories" jsonschema:"title=记忆列表"`
	Total    int                 `json:"total" jsonschema:"title=总数"`
}

type ReadMemorySummary struct {
	ID         uint   `json:"id" jsonschema:"title=记忆ID"`
	Title      string `json:"title" jsonschema:"title=标题"`
	Content    string `json:"content" jsonschema:"title=完整内容"`
	MemoryType string `json:"memory_type,omitempty" jsonschema:"title=类型"`
	CreatedAt  string `json:"created_at" jsonschema:"title=创建时间"`
}

func NewReadMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"read_memory",
		"按关键词或类型查询用户的长期记忆。所有时间/地点/人物信息都已写入 content 字段，请直接用自然语言关键词检索。"+memoryTypeDesc,
		func(ctx context.Context, in *ReadMemoryToolRequest) (*ReadMemoryToolResponse, error) {
			response := &ReadMemoryToolResponse{Success: false, Memories: []ReadMemorySummary{}}

			limit := 20
			var keywords []string
			var typeFilter *string

			if in != nil {
				if in.Limit != nil && *in.Limit > 0 {
					limit = *in.Limit
					if limit > 100 {
						limit = 100
					}
				}
				if in.Keyword != nil {
					for _, kw := range strings.Split(*in.Keyword, ",") {
						kw = strings.TrimSpace(kw)
						if kw != "" {
							keywords = append(keywords, kw)
						}
					}
				}
				if in.MemoryType != nil {
					t := strings.TrimSpace(*in.MemoryType)
					if t != "" {
						normalized, err := normalizeMemoryType(t)
						if err != nil {
							response.Message = err.Error()
							return response, nil
						}
						s := string(normalized)
						typeFilter = &s
					}
				}
			}

			query := istorage.MemoryQuery{
				Keyword: keywords,
				Type:    typeFilter,
				Limit:   limit,
			}

			memories, err := storage.QueryMemories(ctx, query)
			if err != nil {
				response.Message = "查询失败：" + err.Error()
				return response, nil
			}

			summaries := make([]ReadMemorySummary, 0, len(memories))
			for _, m := range memories {
				title := m.Summary
				if title == "" {
					title = "(无标题)"
				}
				summaries = append(summaries, ReadMemorySummary{
					ID:         m.ID,
					Title:      title,
					Content:    m.Content,
					MemoryType: string(m.Type),
					CreatedAt:  m.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}

			response.Success = true
			response.Memories = summaries
			response.Total = len(summaries)
			switch {
			case len(summaries) == 0:
				response.Message = "未找到匹配记忆"
			case len(summaries) == 1:
				response.Message = "找到 1 条相关记忆"
			default:
				response.Message = fmt.Sprintf("找到 %d 条相关记忆", len(summaries))
			}
			return response, nil
		},
	)
}

// ---- edit_memory ----

type EditMemoryIn struct {
	MemoryID   uint    `json:"memory_id" jsonschema:"title=记忆ID,description=需要编辑的记忆 ID。required=true"`
	Title      *string `json:"title,omitempty" jsonschema:"title=标题,description=只有提供时才更新"`
	Content    *string `json:"content,omitempty" jsonschema:"title=内容,description=只有提供时才更新；同样要求把时间/地点/人物等信息融入文本"`
	MemoryType *string `json:"memory_type,omitempty" jsonschema:"title=类型,description=只有提供时才更新（fact / information / event）"`
	IsForget   bool    `json:"is_forget,omitempty" jsonschema:"title=是否标记为已遗忘,description=true 表示软删除该记忆"`
}

type EditMemoryResponse struct {
	WriteMemoryToolResponse
}

func NewEditMemoryTool(storage *istorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"edit_memory",
		"编辑已有记忆。优先用此工具补全/合并旧记忆，而不是新建。仅提供的字段会更新。",
		func(ctx context.Context, input *EditMemoryIn) (*EditMemoryResponse, error) {
			response := &EditMemoryResponse{
				WriteMemoryToolResponse: WriteMemoryToolResponse{Success: false},
			}

			if input == nil || input.MemoryID == 0 {
				response.Message = "无法编辑记忆：缺少 memory_id"
				return response, nil
			}

			if input.IsForget {
				if err := storage.SoftDeleteMemory(ctx, input.MemoryID); err != nil {
					response.Message = "标记遗忘失败：" + err.Error()
					return response, nil
				}
				response.Success = true
				response.Message = "记忆已标记为遗忘"
				response.MemoryID = input.MemoryID
				return response, nil
			}

			update := models.Memory{}
			if input.Title != nil {
				update.Summary = strings.TrimSpace(*input.Title)
			}
			if input.Content != nil {
				update.Content = strings.TrimSpace(*input.Content)
			}
			if input.MemoryType != nil {
				memType, err := normalizeMemoryType(*input.MemoryType)
				if err != nil {
					response.Message = err.Error()
					return response, nil
				}
				update.Type = memType
			}

			if err := storage.UpdateMemory(ctx, input.MemoryID, update); err != nil {
				response.Message = "更新失败：" + err.Error()
				return response, nil
			}

			response.Success = true
			response.Message = "记忆已成功更新"
			response.MemoryID = input.MemoryID
			return response, nil
		},
	)
}

// ---- helpers ----

func normalizeMemoryType(raw string) (models.MemoryType, error) {
	t := strings.ToLower(strings.TrimSpace(raw))
	switch t {
	case "":
		return models.MemoryTypeNone, nil
	case "fact":
		return models.MemoryTypeFact, nil
	case "information", "info":
		return models.MemoryTypeInfo, nil
	case "event":
		return models.MemoryTypeEvent, nil
	default:
		return models.MemoryTypeNone, fmt.Errorf("不支持的记忆类型: %q（可选: fact / information / event）", raw)
	}
}
