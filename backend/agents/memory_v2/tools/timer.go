package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// GetCurrentTimeToolRequest 是获取当前时间工具的输入参数
// 目前不需要输入参数，但为了接口统一和未来扩展，保留空结构
type GetCurrentTimeToolRequest struct{}

// GetCurrentTimeToolResponse 是获取当前时间工具的输出结果
type GetCurrentTimeToolResponse struct {
	Success   bool   `json:"success" jsonschema:"title=是否成功;description=表示是否成功获取当前时间"`
	Message   string `json:"message,omitempty" jsonschema:"title=提示信息;description=操作结果描述，例如 '获取时间成功'"`
	Timestamp string `json:"timestamp" jsonschema:"title=时间戳;description=当前时间的 ISO 8601 格式字符串，例如 '2025-04-05T12:34:56Z'"` // RFC3339
	Timezone  string `json:"timezone,omitempty" jsonschema:"title=时区;description=系统返回时间所处的时区名称，例如 'UTC' 或 'Asia/Shanghai'"`
}

// NewGetCurrentTimeTool 返回一个标准化的可调用工具，用于获取当前时间（北京时间）
func NewGetCurrentTimeTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"get_current_time",
		"获取当前的北京时间（CST, UTC+8）和时区信息",
		func(ctx context.Context, in *GetCurrentTimeToolRequest) (output *GetCurrentTimeToolResponse, err error) {
			// 初始化响应
			response := &GetCurrentTimeToolResponse{
				Success: false,
			}

			// 获取北京时间（UTC+8）
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				response.Message = "无法加载时区信息: " + err.Error()
				return response, nil
			}
			now := time.Now().In(loc)

			// 填充响应数据
			response.Success = true
			response.Message = "当前时间获取成功"
			response.Timestamp = now.Format(time.RFC3339) // ISO 8601 标准格式
			response.Timezone = "Asia/Shanghai"           // 明确表示时区

			return response, nil
		},
	)
}
