// 日期时间 Agent - 可获取当前日期、时间

package agents

import (
	"context"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

const (
	agentInstruction = `你是一个可以查询当前日期和时间的助手。
当用户询问"今天几号"、"现在几点"、"今天是星期几"等问题时，请使用工具获取准确的日期和时间信息后回答。
用简洁清晰的中文回复用户。`
)

// 无参数工具使用的空结构
type emptyParams struct{}

func getDateTimeTools() []tool.BaseTool {
	getDateTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_date",
			Desc:        "获取当前日期，返回格式为 YYYY-MM-DD，以及星期几",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			now := time.Now()
			weekdays := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
			return now.Format("2006-01-02") + " " + weekdays[now.Weekday()], nil
		},
	)

	getTimeTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        "获取当前时间，返回格式为 HH:MM:SS",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			return time.Now().Format("15:04:05"), nil
		},
	)

	return []tool.BaseTool{getDateTool, getTimeTool}
}

// NewDateTimeAgent 创建日期时间 Agent，作为子 Agent 供 ChatAgent 通过 AgentTool 调用
func NewDateTimeAgent(ctx context.Context, chatModel model.ToolCallingChatModel) (adk.Agent, error) {
	tools := getDateTimeTools()

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "DateTimeAgent",
		Description: "可以查询当前日期和时间的助手。当用户询问今天星期几、现在几点、今天几号等日期时间相关问题时，转交给此 Agent 处理。",
		Instruction: agentInstruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
