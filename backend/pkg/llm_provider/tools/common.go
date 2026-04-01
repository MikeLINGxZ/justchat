package tools

import (
	"fmt"
	"sort"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type ITool interface {
	Id() string
	Name() string
	Description() string
	Tool() tool.BaseTool
	RequireConfirmation() bool
}

func init() {
	ToolRouter = &router{
		tools:        map[string]ITool{},
		dynamicTools: map[string]ITool{},
	}
	ToolRouter.RegisterTool(&CurrentDate{})
	ToolRouter.RegisterTool(&CurrentTime{})
	ToolRouter.RegisterTool(&Block{})
	ToolRouter.RegisterTool(&FileTool{})
	ToolRouter.RegisterTool(&ShellTool{})
}

var ToolRouter *router

type router struct {
	mu           sync.RWMutex
	tools        map[string]ITool
	dynamicTools map[string]ITool
}

func (r *router) GetToolsInfo() []ITool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tools []ITool
	for _, val := range r.tools {
		tools = append(tools, val)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Id() < tools[j].Id()
	})
	return tools
}

func (r *router) GetBuiltinToolsInfo() []ITool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tools []ITool
	for toolID, val := range r.tools {
		if _, ok := r.dynamicTools[toolID]; ok {
			continue
		}
		tools = append(tools, val)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Id() < tools[j].Id()
	})
	return tools
}

func (r *router) GetToolsByIds(toolIds []string) ([]tool.BaseTool, error) {
	var tools []tool.BaseTool
	for _, toolId := range toolIds {
		iTool, ok := r.GetToolByID(toolId)
		if !ok {
			return tools, fmt.Errorf("tool %s not found", toolId)
		}
		tools = append(tools, iTool.Tool())
	}
	return tools, nil
}

func (r *router) GetToolByID(toolID string) (ITool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	iTool, ok := r.tools[toolID]
	return iTool, ok
}

func (r *router) RegisterTool(tool ITool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.tools[tool.Id()]
	if ok {
		logger.Warm("tool %s is already registered", tool.Id())
	}
	r.tools[tool.Id()] = tool
}

func (r *router) UpsertDynamicTool(tool ITool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.dynamicTools[tool.Id()] = tool
	r.tools[tool.Id()] = tool
}

func (r *router) RemoveDynamicTool(toolID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.dynamicTools, toolID)
	delete(r.tools, toolID)
}

func (r *router) ResetDynamicTools(tools []ITool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for toolID := range r.dynamicTools {
		delete(r.tools, toolID)
	}
	r.dynamicTools = map[string]ITool{}
	for _, item := range tools {
		r.dynamicTools[item.Id()] = item
		r.tools[item.Id()] = item
	}
}

type emptyParams struct{}
