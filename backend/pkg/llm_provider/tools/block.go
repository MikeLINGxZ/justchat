package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Block struct {
}

func (b *Block) Id() string {
	return "block"
}

func (b *Block) Name() string {
	return "阻塞器"
}

func (b *Block) Tool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name: "block",
			Desc: "获取当前日期，返回格式为 YYYY-MM-DD 以及星期几。在解析「明天」「后天」等相对日期时请先调用此工具",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"time": &schema.ParameterInfo{
					Type:      "integer",
					SubParams: nil,
					Desc:      "需要阻塞的秒数",
					Required:  false,
				},
			}),
		},
		func(ctx context.Context, params blockParams) (string, error) {
			if params.Time > 0 {
				time.Sleep(time.Second * time.Duration(params.Time))
			} else {
				time.Sleep(time.Second * 30)
			}
			return "", nil
		},
	)
	return t
}

func (b *Block) Description() string {
	return "这个一个阻塞器，默认阻塞30s，可以传入所需的阻塞时间"
}

type blockParams struct {
	Time int `json:"time"`
}
