package tools

import (
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type ITool interface {
	Id() string
	Name() string
	Tool() tool.BaseTool
	Description() string
}

func init() {
	ToolRouter = &router{
		tools: map[string]ITool{},
	}
	ToolRouter.RegisterTool(&CurrentDate{})
	ToolRouter.RegisterTool(&CurrentTime{})
}

var ToolRouter *router

type router struct {
	tools map[string]ITool
}

func (r *router) GetToolsInfo() []ITool {
	var tools []ITool
	for _, val := range r.tools {
		tools = append(tools, val)
	}
	return tools
}

func (r *router) GetToolsByIds(toolIds []string) ([]tool.BaseTool, error) {
	var tools []tool.BaseTool
	for _, toolId := range toolIds {
		iTool, ok := r.tools[toolId]
		if !ok {
			return tools, fmt.Errorf("tool %s not found", toolId)
		}
		tools = append(tools, iTool.Tool())
	}
	return tools, nil
}

func (r *router) RegisterTool(tool ITool) {
	_, ok := r.tools[tool.Id()]
	if ok {
		logger.Warm("tool %s is already registered", tool.Id())
	}
	r.tools[tool.Id()] = tool
}

type emptyParams struct{}
