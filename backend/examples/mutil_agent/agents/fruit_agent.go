// 水果价格 Agent - 可查询各种水果的价格（模拟数据）

package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
)

const fruitAgentInstruction = `你是一个水果价格查询助手。
当用户询问水果价格时，请使用 get_fruit_price 工具查询后回答。
支持查询：苹果、香蕉、橙子、葡萄、草莓、西瓜、芒果、樱桃、荔枝、榴莲、火龙果、猕猴桃等。
用简洁清晰的中文回复，价格单位为元/斤。`

// 水果价格模拟数据（元/斤）
var fruitPrices = map[string]float64{
	"苹果": 5.5, "香蕉": 3.0, "橙子": 4.5, "葡萄": 12.0, "草莓": 18.0,
	"西瓜": 2.5, "芒果": 8.0, "樱桃": 35.0, "荔枝": 15.0, "榴莲": 28.0,
	"火龙果": 6.0, "猕猴桃": 10.0, "梨": 4.0, "桃子": 6.5, "李子": 7.0,
	"柚子": 5.0, "柠檬": 8.0, "蓝莓": 25.0, "车厘子": 40.0, "山竹": 20.0,
	"杨梅": 15.0, "龙眼": 12.0, "枇杷": 18.0, "石榴": 9.0, "无花果": 22.0,
}

type fruitPriceParams struct {
	FruitName string `json:"fruit_name" jsonschema:"description=水果名称，如：苹果、香蕉、橙子"`
}

func getFruitPriceTools() []tool.BaseTool {
	getPriceTool, _ := utils.InferTool(
		"get_fruit_price",
		"查询指定水果的价格，返回每斤的单价（元）。参数为水果名称。",
		func(ctx context.Context, params fruitPriceParams) (string, error) {
			name := strings.TrimSpace(params.FruitName)
			if name == "" {
				return "请提供要查询的水果名称。", nil
			}
			// 尝试精确匹配或模糊匹配
			for k, v := range fruitPrices {
				if k == name || strings.Contains(name, k) || strings.Contains(k, name) {
					return fmt.Sprintf("%s：%.1f 元/斤", k, v), nil
				}
			}
			// 未找到时列出已有水果
			available := make([]string, 0, len(fruitPrices))
			for k := range fruitPrices {
				available = append(available, k)
			}
			return fmt.Sprintf("暂无「%s」的价格数据。当前支持查询：%v", name, available), nil
		},
	)
	return []tool.BaseTool{getPriceTool}
}

// NewFruitPriceAgent 创建水果价格 Agent，作为子 Agent 供 ChatAgent 通过 AgentTool 调用
func NewFruitPriceAgent(ctx context.Context, chatModel model.ToolCallingChatModel) (adk.Agent, error) {
	tools := getFruitPriceTools()

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "FruitPriceAgent",
		Description: "可以查询各种水果价格的助手。当用户询问苹果多少钱、草莓价格、榴莲贵不贵等水果价格相关问题时，转交给此 Agent 处理。",
		Instruction: fruitAgentInstruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
