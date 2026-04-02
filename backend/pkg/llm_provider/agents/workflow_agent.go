package agents

// ---- Agent 元数据定义 ----

// WorkflowAgentDef 是工作流编排的逻辑容器，本身不拥有提示词，
// 子 Agent（Planner、Worker、Synthesizer、Reviewer）各自管理提示词。
type WorkflowAgentDef struct{}

func init() { RegisterAgent(&WorkflowAgentDef{}) }

func (a *WorkflowAgentDef) Name() string    { return "workflow_agent" }
func (a *WorkflowAgentDef) Desc() string    { return "工作流编排 Agent，协调规划、执行、综合、审核等子 Agent" }
func (a *WorkflowAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *WorkflowAgentDef) Role() AgentRole { return AgentRoleWorkflow }
func (a *WorkflowAgentDef) PromptNames() []string         { return nil }
func (a *WorkflowAgentDef) PromptMetas() []AgentPromptMeta { return nil }
func (a *WorkflowAgentDef) DefaultPrompts() map[string]string {
	return nil
}
func (a *WorkflowAgentDef) Prompt() string { return "" }
