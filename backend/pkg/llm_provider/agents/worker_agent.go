package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// ---- Agent 元数据定义 ----

// GeneralWorkerAgentDef 通用执行 Agent
type GeneralWorkerAgentDef struct{}

func init() {
	RegisterAgent(&GeneralWorkerAgentDef{})
	RegisterAgent(&ToolWorkerAgentDef{})
}

func (a *GeneralWorkerAgentDef) Name() string    { return "general_worker_agent" }
func (a *GeneralWorkerAgentDef) Desc() string    { return "通用 Worker Agent，负责执行子任务并产出结果" }
func (a *GeneralWorkerAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *GeneralWorkerAgentDef) Role() AgentRole { return AgentRoleWorker }
func (a *GeneralWorkerAgentDef) PromptNames() []string {
	return []string{"system.worker.general.md"}
}
func (a *GeneralWorkerAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.worker.general.md", Title: "通用 Worker", Description: "约束通用子任务执行代理的职责和输出风格。", IsSystem: true},
	}
}
func (a *GeneralWorkerAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.worker.general.md": `你是负责执行单个子任务的 WorkerAgent。
- 只完成当前分配的子任务
- 必要时调用工具
- 不要假装完成没拿到的数据
- 输出该子任务的结果摘要，便于后续汇总`,
	}
}
func (a *GeneralWorkerAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.worker.general.md", a.DefaultPrompts()["system.worker.general.md"])
	return content
}

// ToolWorkerAgentDef 工具专家 Agent
type ToolWorkerAgentDef struct{}

func (a *ToolWorkerAgentDef) Name() string    { return "tool_worker_agent" }
func (a *ToolWorkerAgentDef) Desc() string    { return "工具 Worker Agent，以工具调用为主的子任务执行" }
func (a *ToolWorkerAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *ToolWorkerAgentDef) Role() AgentRole { return AgentRoleWorker }
func (a *ToolWorkerAgentDef) PromptNames() []string {
	return []string{"system.worker.tool.md"}
}
func (a *ToolWorkerAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.worker.tool.md", Title: "工具 Worker", Description: "约束以工具事实获取为主的子任务代理。", IsSystem: true},
	}
}
func (a *ToolWorkerAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.worker.tool.md": `你是 ToolSpecialistAgent。
- 当前任务优先通过工具获取事实或数据
- 如果工具不足，明确说明缺口
- 输出该子任务的结果摘要，便于后续汇总`,
	}
}
func (a *ToolWorkerAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.worker.tool.md", a.DefaultPrompts()["system.worker.tool.md"])
	return content
}

// ---- ADK Agent 运行时创建 ----

// NewRoleAgent 创建角色 Agent（用于工作流子任务执行）。
func NewRoleAgent(ctx context.Context, chatModel model.ToolCallingChatModel, name, description, instruction string, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        name,
		Description: description,
		Instruction: instruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Middlewares: []adk.AgentMiddleware{
			{
				WrapToolCall: toolMiddleware,
			},
		},
	})
}
