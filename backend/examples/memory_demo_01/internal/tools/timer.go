package tools

import (
	"context"
	"time"
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

// NewGetCurrentTimeTool 返回一个可被调用的工具函数，用于获取当前时间
func NewGetCurrentTimeTool() func(ctx context.Context, in *GetCurrentTimeToolRequest) (*GetCurrentTimeToolResponse, error) {
	return func(ctx context.Context, in *GetCurrentTimeToolRequest) (*GetCurrentTimeToolResponse, error) {
		// 初始化响应
		response := &GetCurrentTimeToolResponse{
			Success: false,
		}

		// 获取当前时间（UTC 时间，确保标准化）
		now := time.Now().UTC()

		// 设置成功响应
		response.Success = true
		response.Message = "当前时间获取成功"
		response.Timestamp = now.Format(time.RFC3339) // 标准 ISO 8601 / RFC3339 格式
		response.Timezone = "UTC"                     // 明确标注为 UTC，也可根据需求改为本地时区或传参指定

		return response, nil
	}
}
