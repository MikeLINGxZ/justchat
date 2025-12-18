package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory_v2/tools"
)

func main() {
	ctx := context.Background()
	qwenModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		APIKey:  os.Getenv("ALIYUN_API_KEY"),
		Model:   "qwen-plus-2025-07-14",
	})
	if err != nil {
		panic(err)
	}

	timeTool, _ := tools.NewGetCurrentTimeTool()
	dataTool, _ := NewFakeDataTool()

	chatAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "智能助理",
		Description: "能够使用多种工具解决复杂问题的智能助手",
		Instruction: "您是一名专业助理，可以使用提供的工具帮助用户解决问题",
		Model:       qwenModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					timeTool,
					dataTool,
				},
			},
			ReturnDirectly: nil,
		},
	})
	if err != nil {
		panic(err)
	}

	run := chatAgent.Run(ctx, &adk.AgentInput{
		Messages: []adk.Message{
			{
				Role:    schema.User,
				Content: "今天的产值是多少？",
			},
		},
		EnableStreaming: true,
	})

	for {
		next, hasNext := run.Next()
		if !hasNext {
			break
		}
		marshal, _ := json.Marshal(next)
		fmt.Println(string(marshal))
	}

}

type NewFakeDataToolRequest struct {
	Date string `json:"date"  jsonschema:"title=需要获取数据的日前;description=格式为yyyy-MM-dd"`
}

type NewFakeDataToolResponse struct {
	Success bool   `json:"success" jsonschema:"title=是否成功;description=表示是否成功获取当前时间"`
	Message string `json:"message,omitempty" jsonschema:"title=提示信息;description=操作结果描述，例如 '获取时间成功'"`
	Value   int64  `json:"value" jsonschema:"title=产值;description=产值（亿元）"`
}

func NewFakeDataTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"get_data",
		"获取某天产值",
		func(ctx context.Context, in *NewFakeDataToolRequest) (output *NewFakeDataToolResponse, err error) {
			ymdInt, err := ymdDashToYmdInt(in.Date)
			if err != nil {
				return nil, err
			}
			output = &NewFakeDataToolResponse{
				Success: true,
				Value:   ymdInt,
				Message: "获取产值成功",
			}
			return
		},
	)
}

func ymdDashToYmdInt(s string) (int64, error) {
	// 简单校验长度和分隔符位置（可选，增强健壮性）
	if len(s) != 10 || s[4] != '-' || s[7] != '-' {
		return 0, fmt.Errorf("invalid format: expected yyyy-MM-dd, got %q", s)
	}
	// 提取年、月、日数字部分（避免内存分配，用切片+strconv）
	year := s[0:4]
	month := s[5:7]
	day := s[8:10]
	// 拼接为 yyyymmdd 字符串（注意：这里用字符串拼接，但可优化为 []byte 构造避免中间字符串）
	ymdStr := year + month + day
	return strconv.ParseInt(ymdStr, 10, 64)
}
