package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

type Block struct {
}

func (b *Block) Id() string {
	return "block"
}

func (b *Block) Name() string {
	return i18n.TCurrent("tool.block.name", nil)
}

func (b *Block) Tool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name: "block",
			Desc: i18n.TCurrent("tool.block.description", nil),
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
	return i18n.TCurrent("tool.block.description", nil)
}

type blockParams struct {
	Time int `json:"time"`
}
