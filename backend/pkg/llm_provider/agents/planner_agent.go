package agents

// ---- Agent 元数据定义 ----

type PlannerAgentDef struct{}

func init() { RegisterAgent(&PlannerAgentDef{}) }

func (a *PlannerAgentDef) Name() string { return "planner_agent" }
func (a *PlannerAgentDef) Desc() string {
	return "规划 Agent，负责将用户请求拆解为可执行的子任务计划"
}
func (a *PlannerAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *PlannerAgentDef) Role() AgentRole { return AgentRolePlanner }
func (a *PlannerAgentDef) PromptNames() []string {
	return []string{"system.planner.md", "user.planning.md"}
}
func (a *PlannerAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.planner.md", Title: "Planner 系统提示词", Description: "控制任务拆解、JSON 计划输出和工具使用约束。", IsSystem: true},
		{FileName: "user.planning.md", Title: "规划用户模板", Description: "组织用户请求与最近上下文，提供给 Planner。", IsSystem: false},
	}
}
func (a *PlannerAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.planner.md": `你是 PlannerAgent。请把用户任务拆成结构化执行计划。
要求：
- 只输出 JSON
- tasks 至少 1 个，最多 4 个
- dependencies 使用 task.id
- suggested_agent 只允许: GeneralWorkerAgent, ToolSpecialistAgent
- required_tools 填写你认为需要的工具名，可为空
- completion_criteria 写成字符串数组

可用工具：{{tool_names}}

输出格式：
{
  "goal":"...",
  "completion_criteria":["..."],
  "tasks":[
    {
      "id":"task_1",
      "title":"...",
      "description":"...",
      "dependencies":["task_0"],
      "suggested_agent":"GeneralWorkerAgent",
      "required_tools":["tool_a"],
      "expected_output":"..."
    }
  ]
}`,
		"user.planning.md": `用户请求：{{user_request}}

最近上下文：
{{recent_context}}`,
	}
}
func (a *PlannerAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.planner.md", a.DefaultPrompts()["system.planner.md"])
	return content
}
