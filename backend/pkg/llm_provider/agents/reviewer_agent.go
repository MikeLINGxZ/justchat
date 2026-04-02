package agents

// ---- Agent 元数据定义 ----

type ReviewerAgentDef struct{}

func init() { RegisterAgent(&ReviewerAgentDef{}) }

func (a *ReviewerAgentDef) Name() string    { return "reviewer_agent" }
func (a *ReviewerAgentDef) Desc() string    { return "审核 Agent，负责检查工作流最终答案是否满足目标" }
func (a *ReviewerAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *ReviewerAgentDef) Role() AgentRole { return AgentRoleReviewer }
func (a *ReviewerAgentDef) PromptNames() []string {
	return []string{"system.reviewer.md", "user.review.md"}
}
func (a *ReviewerAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.reviewer.md", Title: "Reviewer 系统提示词", Description: "控制最终答案审核的 JSON 输出与判定标准。", IsSystem: true},
		{FileName: "user.review.md", Title: "审核用户模板", Description: "把候选答案和子任务结果组织给 Reviewer。", IsSystem: false},
	}
}
func (a *ReviewerAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.reviewer.md": `你是 ReviewerAgent。请审核最终答案是否满足目标和完成标准。
- 只输出 JSON
- 如果通过，approved=true，issues 可为空
- 如果不通过，给出 issues、retry_instructions，并尽量指出 affected_task_ids

格式：
{
  "approved": true,
  "issues": [],
  "retry_instructions": "",
  "affected_task_ids": ["task_1"]
}`,
		"user.review.md": `用户请求：{{user_request}}
整体目标：{{goal}}
完成标准：{{completion_criteria}}
子任务结果：
{{task_results}}

候选答案：
{{draft}}`,
	}
}
func (a *ReviewerAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.reviewer.md", a.DefaultPrompts()["system.reviewer.md"])
	return content
}
