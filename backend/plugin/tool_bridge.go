package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
)

// toolExecParams is the JSON-RPC parameter for tool/execute calls.
type toolExecParams struct {
	PluginID string          `json:"pluginId"`
	ToolID   string          `json:"toolId"`
	Input    json.RawMessage `json:"input"`
}

// toolExecResult is the JSON-RPC result from tool/execute calls.
type toolExecResult struct {
	Content string `json:"content"`
}

// PluginTool bridges a plugin-contributed tool to the Eino tool system.
// It implements both tools.ITool and einotool.InvokableTool (BaseTool).
type PluginTool struct {
	pluginID          string
	pluginDisplayName string
	toolID            string // full namespaced ID: "plugin:<pluginName>:<localId>"
	localToolID       string
	name              string
	description       string
	parameters        map[string]any // JSON Schema for tool parameters
	host              *ExtensionHost
}

func (p *PluginTool) Id() string                { return p.toolID }
func (p *PluginTool) Name() string              { return p.name }
func (p *PluginTool) Description() string       { return p.description }
func (p *PluginTool) PluginDisplayName() string { return p.pluginDisplayName }

// BuildApprovalPrompt implements tool_approval.ApprovalAwareTool.
func (p *PluginTool) BuildApprovalPrompt(_ context.Context, argumentsJSON string) (*tool_approval.ApprovalPrompt, error) {
	return &tool_approval.ApprovalPrompt{
		Title:   fmt.Sprintf("插件工具: %s", p.name),
		Message: fmt.Sprintf("插件「%s」请求执行工具「%s」\n参数: %s", p.pluginDisplayName, p.name, argumentsJSON),
		Scope:   "plugin",
	}, nil
}

// RequireConfirmation returns true because plugin tools always need user approval.
func (p *PluginTool) RequireConfirmation() bool { return true }

// Tool returns self since PluginTool implements einotool.BaseTool.
func (p *PluginTool) Tool() einotool.BaseTool { return p }

// Info returns the Eino ToolInfo for this plugin tool.
func (p *PluginTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	info := &schema.ToolInfo{
		Name: p.toolID,
		Desc: p.description,
	}

	if p.parameters != nil {
		paramBytes, err := json.Marshal(p.parameters)
		if err != nil {
			return nil, fmt.Errorf("marshal tool parameters: %w", err)
		}
		var js jsonschema.Schema
		if err := json.Unmarshal(paramBytes, &js); err != nil {
			return nil, fmt.Errorf("parse tool parameters schema: %w", err)
		}
		info.ParamsOneOf = schema.NewParamsOneOfByJSONSchema(&js)
	}

	return info, nil
}

// InvokableRun executes the plugin tool via JSON-RPC.
func (p *PluginTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	raw, err := p.host.RPC().Call("tool/execute", toolExecParams{
		PluginID: p.pluginID,
		ToolID:   p.localToolID,
		Input:    json.RawMessage(argumentsInJSON),
	})
	if err != nil {
		return "", fmt.Errorf("plugin tool %s execution failed: %w", p.toolID, err)
	}

	var result toolExecResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("unmarshal tool result for %s: %w", p.toolID, err)
	}

	return result.Content, nil
}

// ToolBridge manages the lifecycle of plugin-contributed tools.
type ToolBridge struct {
	mu    sync.RWMutex
	tools map[string]*PluginTool // full toolID → PluginTool
}

// NewToolBridge creates a new ToolBridge instance.
func NewToolBridge() *ToolBridge {
	return &ToolBridge{
		tools: make(map[string]*PluginTool),
	}
}

// Register creates a PluginTool from a ToolContrib and registers it with the global tool router.
func (b *ToolBridge) Register(pluginID, pluginName, pluginDisplayName string, contrib ToolContrib, host *ExtensionHost) {
	fullID := FullToolID(pluginName, contrib.ID)

	pt := &PluginTool{
		pluginID:          pluginID,
		pluginDisplayName: pluginDisplayName,
		toolID:            fullID,
		localToolID:       contrib.ID,
		name:              contrib.Name,
		description:       contrib.Description,
		parameters:        contrib.Parameters,
		host:              host,
	}

	b.mu.Lock()
	b.tools[fullID] = pt
	b.mu.Unlock()

	tools.ToolRouter.UpsertDynamicTool(pt)
}

// UnregisterByPlugin removes all tools belonging to the given plugin.
func (b *ToolBridge) UnregisterByPlugin(pluginID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for toolID, pt := range b.tools {
		if pt.pluginID == pluginID {
			tools.ToolRouter.RemoveDynamicTool(toolID)
			delete(b.tools, toolID)
		}
	}
}

// GetToolsByPlugin returns all tools registered by the given plugin.
func (b *ToolBridge) GetToolsByPlugin(pluginID string) []*PluginTool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var result []*PluginTool
	for _, pt := range b.tools {
		if pt.pluginID == pluginID {
			result = append(result, pt)
		}
	}
	return result
}
