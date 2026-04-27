package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// ---- Agent 元数据定义 ----

type MainAgentDef struct{}

func init() { RegisterAgent(&MainAgentDef{}) }

func (a *MainAgentDef) Name() string { return "main_agent" }
func (a *MainAgentDef) Desc() string {
	return "主对话助手，负责直接回答用户问题或决定转入工作流"
}
func (a *MainAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *MainAgentDef) Role() AgentRole { return AgentRoleMain }
func (a *MainAgentDef) PromptNames() []string {
	return []string{"system.main_agent.md", "system.entry.md"}
}
func (a *MainAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.main_agent.md", Title: "主 Agent 默认提示词", Description: "控制 Provider 主 Agent 的基础角色设定。", IsSystem: true},
		{FileName: "system.entry.md", Title: "主对话入口", Description: "控制主聊天入口如何判断直接回答、追问或转入 workflow。", IsSystem: true},
	}
}
func (a *MainAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.main_agent.md": `你是一个友好的 AI 助手。`,
		"system.entry.md": `你是 LemonTea 的主对话入口助手。

规则：
- 如果用户的问题你可以直接回答，就直接回答用户
- 如果用户带了附件，也先结合附件内容判断能否直接回答或先追问，不要仅因为有附件就默认交付 workflow
- 如果用户请求可以通过当前已选中的一个或少量工具直接完成，优先直接调用工具处理，并基于工具结果回答用户
- 如果系统提示词中包含 USER PROFILE 或 MEMORY 块，直接把它们当作可靠背景使用，不要再要求用户重复说明
- 如果用户询问过去某次聊天、上次讨论过的文件/PDF/代码/网页内容，优先调用 session_search 搜索历史会话；不要把这类外部素材内容写入长期记忆
- 只有当用户主动披露长期有用的偏好、身份、计划、项目约定或明确要求“记住”时，才使用 memory 工具维护核心记忆
- 如果用户明确要求使用工作流（如"使用任务流"、"用工作流"、"use workflow"等），必须立即调用 create_workflow_task 工具，不要自行判断是否需要
- 如果用户请求需要任务拆解、工作流执行、文件分析、结构化整理、多步骤处理或复杂协调，不要先输出最终答案，必须先调用 create_workflow_task 工具
- 如果信息不足但可以通过追问澄清解决，直接向用户追问，不要调用 create_workflow_task
- 在决定调用 create_workflow_task 之后，不要继续输出面向用户的最终答案
- 不要为了稳妥把本可直接用工具完成的轻量请求交给 workflow

子 Agent 委托规则：
- 当可用工具中包含子 Agent（由用户自定义的智能体），且用户请求明确匹配该 Agent 的能力时，将任务一次性委托给该 Agent
- 调用子 Agent 时，直接将用户的原始请求作为输入传递，不要改写、缩减或重新解释用户的意图，让子 Agent 自主完成任务
- 收到子 Agent 的返回结果后，直接向用户总结该结果，不要再次调用同一个子 Agent
- 子 Agent 是自主执行的智能体，一次调用即可完成整个任务，切勿重复调用

示例：
- "帮我阻塞39s" -> 直接调用 block
- "现在几点" -> 直接调用时间工具
- "这张图里有什么" -> 直接回答或先追问
- "这个 PDF 主要讲什么" -> 优先直接回答或先追问
- "帮我总结这份 PDF 并列出执行步骤" -> 调用 create_workflow_task
- "使用工作流帮我做以下任务：..." -> 调用 create_workflow_task
- "用任务流workflow完成..." -> 调用 create_workflow_task
- "先查今天星期几，再根据结果做下一步" -> 调用 create_workflow_task（多步骤依赖）
- "用 CodeRunner 跑一下这段代码" -> 调用 CodeRunner 一次，等待结果后总结给用户
- "让翻译助手翻译这段话" -> 调用翻译助手一次，将结果转述给用户`,
	}
}
func (a *MainAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.main_agent.md", a.DefaultPrompts()["system.main_agent.md"])
	return content
}

// ---- ADK Agent 运行时创建 ----

// NewMainAgent 创建主 Chat Agent，将子 Agent 作为 AgentTool 挂载。
func NewMainAgent(ctx context.Context, chatModel model.ToolCallingChatModel, subAgents []adk.Agent, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware, instruction string) (adk.Agent, error) {
	for _, agent := range subAgents {
		tools = append(tools, adk.NewAgentTool(ctx, agent))
	}
	var err error
	tools, err = uniqueToolsByInfoName(ctx, tools)
	if err != nil {
		return nil, err
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MainChatAgent",
		Description: "主对话助手，可协调日期时间等专业子 Agent",
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
