package agents

// ---- Agent 元数据定义 ----

type SynthesizerAgentDef struct{}

func init() { RegisterAgent(&SynthesizerAgentDef{}) }

func (a *SynthesizerAgentDef) Name() string    { return "synthesizer_agent" }
func (a *SynthesizerAgentDef) Desc() string    { return "综合 Agent，负责整合所有子任务输出为最终答案" }
func (a *SynthesizerAgentDef) Type() AgentType { return AgentTypeSystem }
func (a *SynthesizerAgentDef) Role() AgentRole { return AgentRoleSynthesizer }
func (a *SynthesizerAgentDef) PromptNames() []string {
	return []string{"system.synthesizer.md", "user.synthesis.md"}
}
func (a *SynthesizerAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{FileName: "system.synthesizer.md", Title: "Synthesizer 系统提示词", Description: "控制最终答案的整合方式和输出要求。", IsSystem: true},
		{FileName: "user.synthesis.md", Title: "汇总用户模板", Description: "组织子任务结果、目标和 review 反馈给 Synthesizer。", IsSystem: false},
	}
}
func (a *SynthesizerAgentDef) DefaultPrompts() map[string]string {
	return map[string]string{
		"system.synthesizer.md": `你是 SynthesizerAgent。请根据子任务结果生成给用户的最终答复。
- 必须覆盖用户目标
- 优先整合已有结果，不凭空编造
- 如果有 review 反馈，必须修正
- 直接输出面向用户的最终内容，不要输出 JSON`,
		"user.synthesis.md": `用户请求：{{user_request}}

整体目标：{{goal}}
完成标准：{{completion_criteria}}

子任务结果：
{{task_results}}

审核反馈：{{review_feedback}}`,
	}
}
func (a *SynthesizerAgentDef) Prompt() string {
	content, _ := LoadAgentPrompt(a.Name(), "system.synthesizer.md", a.DefaultPrompts()["system.synthesizer.md"])
	return content
}
