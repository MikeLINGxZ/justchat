package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type emptyParams struct{}

// GetCurrentDateTool 获取当前日期（YYYY-MM-DD 及星期几）
func GetCurrentDateTool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_date",
			Desc:        "获取当前日期，返回格式为 YYYY-MM-DD 以及星期几。在解析「明天」「后天」等相对日期时请先调用此工具",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			now := time.Now()
			weekdays := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
			return now.Format("2006-01-02") + " " + weekdays[now.Weekday()], nil
		},
	)
	return t
}

// GetCurrentTimeTool 获取当前时间（HH:MM:SS）
func GetCurrentTimeTool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        "获取当前时间，返回格式为 HH:MM:SS",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			return time.Now().Format("15:04:05"), nil
		},
	)
	return t
}
