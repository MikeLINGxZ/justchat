package prompts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

type PromptSet struct {
	MainAgentSystem     string
	EntrySystem         string
	PlannerSystem       string
	PlanningUser        string
	WorkerGeneralSystem string
	WorkerToolSystem    string
	SynthesizerSystem   string
	SynthesisUser       string
	ReviewerSystem      string
	ReviewUser          string
	TitleSystem         string
	MemorySystem        string
}

func WithResponseLanguage(promptSet PromptSet, locale string) PromptSet {
	instruction := responseLanguageInstruction(locale)
	if instruction == "" {
		return promptSet
	}

	withInstruction := func(content string) string {
		content = strings.TrimSpace(content)
		if content == "" {
			return instruction
		}
		return content + "\n\n" + instruction
	}

	promptSet.MainAgentSystem = withInstruction(promptSet.MainAgentSystem)
	promptSet.EntrySystem = withInstruction(promptSet.EntrySystem)
	promptSet.PlannerSystem = withInstruction(promptSet.PlannerSystem)
	promptSet.WorkerGeneralSystem = withInstruction(promptSet.WorkerGeneralSystem)
	promptSet.WorkerToolSystem = withInstruction(promptSet.WorkerToolSystem)
	promptSet.SynthesizerSystem = withInstruction(promptSet.SynthesizerSystem)
	promptSet.ReviewerSystem = withInstruction(promptSet.ReviewerSystem)
	promptSet.TitleSystem = withInstruction(promptSet.TitleSystem)
	promptSet.MemorySystem = withInstruction(promptSet.MemorySystem)
	return promptSet
}

func responseLanguageInstruction(locale string) string {
	switch strings.TrimSpace(locale) {
	case "en-US":
		return "Language rule:\n- Always respond in English unless the user explicitly asks you to switch to another language.\n- Keep tool usage, analysis, final answer, summaries, and generated titles in English."
	case "zh-CN", "":
		return "语言规则：\n- 默认始终使用简体中文回复，除非用户明确要求切换到其他语言。\n- 工具说明、分析、最终答案、总结以及生成的标题都应保持为简体中文。"
	default:
		return "Language rule:\n- Always reply in the user's configured application language.\n- Keep tool usage, analysis, final answer, summaries, and generated titles in that language."
	}
}

type PromptMetadata struct {
	Name        string
	Title       string
	Description string
	IsSystem    bool
}

func PromptFiles() []PromptMetadata {
	return []PromptMetadata{
		{Name: "system.entry.md", Title: "主对话入口", Description: "控制主聊天入口如何判断直接回答、追问或转入 workflow。", IsSystem: true},
		{Name: "system.planner.md", Title: "Planner 系统提示词", Description: "控制任务拆解、JSON 计划输出和工具使用约束。", IsSystem: true},
		{Name: "user.planning.md", Title: "规划用户模板", Description: "组织用户请求与最近上下文，提供给 Planner。", IsSystem: false},
		{Name: "system.worker.general.md", Title: "通用 Worker", Description: "约束通用子任务执行代理的职责和输出风格。", IsSystem: true},
		{Name: "system.worker.tool.md", Title: "工具 Worker", Description: "约束以工具事实获取为主的子任务代理。", IsSystem: true},
		{Name: "system.synthesizer.md", Title: "Synthesizer 系统提示词", Description: "控制最终答案的整合方式和输出要求。", IsSystem: true},
		{Name: "user.synthesis.md", Title: "汇总用户模板", Description: "组织子任务结果、目标和 review 反馈给 Synthesizer。", IsSystem: false},
		{Name: "system.reviewer.md", Title: "Reviewer 系统提示词", Description: "控制最终答案审核的 JSON 输出与判定标准。", IsSystem: true},
		{Name: "user.review.md", Title: "审核用户模板", Description: "把候选答案和子任务结果组织给 Reviewer。", IsSystem: false},
		{Name: "system.title.md", Title: "标题生成", Description: "控制聊天标题生成的摘要风格与限制。", IsSystem: true},
		{Name: "system.memory.md", Title: "记忆代理", Description: "控制长期记忆代理的回忆、写入和编辑策略。", IsSystem: true},
		{Name: "system.main_agent.md", Title: "主 Agent 默认提示词", Description: "控制 Provider 主 Agent 的基础角色设定。", IsSystem: true},
	}
}

func FindPromptMetadata(name string) (PromptMetadata, bool) {
	for _, item := range PromptFiles() {
		if item.Name == name {
			return item, true
		}
	}
	return PromptMetadata{}, false
}

func DefaultPromptContent(name string) (string, bool) {
	defaults := defaultPromptSet()
	switch name {
	case "system.main_agent.md":
		return defaults.MainAgentSystem, true
	case "system.entry.md":
		return defaults.EntrySystem, true
	case "system.planner.md":
		return defaults.PlannerSystem, true
	case "user.planning.md":
		return defaults.PlanningUser, true
	case "system.worker.general.md":
		return defaults.WorkerGeneralSystem, true
	case "system.worker.tool.md":
		return defaults.WorkerToolSystem, true
	case "system.synthesizer.md":
		return defaults.SynthesizerSystem, true
	case "user.synthesis.md":
		return defaults.SynthesisUser, true
	case "system.reviewer.md":
		return defaults.ReviewerSystem, true
	case "user.review.md":
		return defaults.ReviewUser, true
	case "system.title.md":
		return defaults.TitleSystem, true
	case "system.memory.md":
		return defaults.MemorySystem, true
	default:
		return "", false
	}
}

func PromptPath(name string) (string, error) {
	if _, ok := FindPromptMetadata(name); !ok {
		return "", fmt.Errorf("unsupported prompt file: %s", name)
	}
	dir, err := PromptDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func SavePromptContent(name string, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("prompt content cannot be empty")
	}
	path, err := PromptPath(name)
	if err != nil {
		return err
	}
	return writeFallbackPrompt(path, content)
}

func PromptDir() (string, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(dataPath, "prompt")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func LoadPrompt(name string, fallback string) (string, error) {
	dir, err := PromptDir()
	if err != nil {
		return fallbackWithWriteError(name, fallback, err)
	}
	path := filepath.Join(dir, name)

	content, readErr := os.ReadFile(path)
	if readErr == nil {
		text := strings.TrimSpace(string(content))
		if text != "" {
			return text, nil
		}
		writeErr := writeFallbackPrompt(path, fallback)
		return fallback, errors.Join(fmt.Errorf("prompt %s is empty", name), writeErr)
	}

	if os.IsNotExist(readErr) {
		writeErr := writeFallbackPrompt(path, fallback)
		return fallback, writeErr
	}

	writeErr := writeFallbackPrompt(path, fallback)
	return fallback, errors.Join(fmt.Errorf("read prompt %s: %w", name, readErr), writeErr)
}

func LoadPromptSet() (PromptSet, error) {
	defaults := defaultPromptSet()
	loaders := []struct {
		name string
		set  func(value string)
		want string
	}{
		{name: "system.main_agent.md", want: defaults.MainAgentSystem, set: func(value string) { defaults.MainAgentSystem = value }},
		{name: "system.entry.md", want: defaults.EntrySystem, set: func(value string) { defaults.EntrySystem = value }},
		{name: "system.planner.md", want: defaults.PlannerSystem, set: func(value string) { defaults.PlannerSystem = value }},
		{name: "user.planning.md", want: defaults.PlanningUser, set: func(value string) { defaults.PlanningUser = value }},
		{name: "system.worker.general.md", want: defaults.WorkerGeneralSystem, set: func(value string) { defaults.WorkerGeneralSystem = value }},
		{name: "system.worker.tool.md", want: defaults.WorkerToolSystem, set: func(value string) { defaults.WorkerToolSystem = value }},
		{name: "system.synthesizer.md", want: defaults.SynthesizerSystem, set: func(value string) { defaults.SynthesizerSystem = value }},
		{name: "user.synthesis.md", want: defaults.SynthesisUser, set: func(value string) { defaults.SynthesisUser = value }},
		{name: "system.reviewer.md", want: defaults.ReviewerSystem, set: func(value string) { defaults.ReviewerSystem = value }},
		{name: "user.review.md", want: defaults.ReviewUser, set: func(value string) { defaults.ReviewUser = value }},
		{name: "system.title.md", want: defaults.TitleSystem, set: func(value string) { defaults.TitleSystem = value }},
		{name: "system.memory.md", want: defaults.MemorySystem, set: func(value string) { defaults.MemorySystem = value }},
	}

	var errs []error
	for _, item := range loaders {
		value, err := LoadPrompt(item.name, item.want)
		if err != nil {
			errs = append(errs, err)
		}
		item.set(value)
	}
	return defaults, errors.Join(errs...)
}

func Render(template string, values map[string]string) string {
	if len(values) == 0 {
		return template
	}
	replacements := make([]string, 0, len(values)*2)
	for key, value := range values {
		replacements = append(replacements, "{{"+key+"}}", value)
	}
	return strings.NewReplacer(replacements...).Replace(template)
}

func fallbackWithWriteError(name string, fallback string, cause error) (string, error) {
	dir, dirErr := PromptDir()
	if dirErr != nil {
		return fallback, errors.Join(fmt.Errorf("resolve prompt dir for %s: %w", name, cause), dirErr)
	}
	path := filepath.Join(dir, name)
	writeErr := writeFallbackPrompt(path, fallback)
	return fallback, errors.Join(fmt.Errorf("resolve prompt dir for %s: %w", name, cause), writeErr)
}

func writeFallbackPrompt(path string, fallback string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(ensureTrailingNewline(fallback)), 0o644)
}

func ensureTrailingNewline(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	return content + "\n"
}

func defaultPromptSet() PromptSet {
	return PromptSet{
		MainAgentSystem: `你是一个友好的 AI 助手。`,
		EntrySystem: `你是 LemonTea 的主对话入口助手。

规则：
- 如果用户的问题你可以直接回答，就直接回答用户
- 如果用户带了附件，也先结合附件内容判断能否直接回答或先追问，不要仅因为有附件就默认交付 workflow
- 如果用户请求可以通过当前已选中的一个或少量工具直接完成，优先直接调用工具处理，并基于工具结果回答用户
- 如果用户请求需要任务拆解、工作流执行、文件分析、结构化整理、多步骤处理或复杂协调，不要先输出最终答案，必须先调用 create_workflow_task 工具
- 如果信息不足但可以通过追问澄清解决，直接向用户追问，不要调用 create_workflow_task
- 在决定调用 create_workflow_task 之后，不要继续输出面向用户的最终答案
- 不要为了稳妥把本可直接用工具完成的轻量请求交给 workflow

示例：
- “帮我阻塞39s” -> 直接调用 block
- “现在几点” -> 直接调用时间工具
- “这张图里有什么” -> 直接回答或先追问
- “这个 PDF 主要讲什么” -> 优先直接回答或先追问
- “帮我总结这份 PDF 并列出执行步骤” -> 调用 create_workflow_task`,
		PlannerSystem: `你是 PlannerAgent。请把用户任务拆成结构化执行计划。
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
		PlanningUser: `用户请求：{{user_request}}

最近上下文：
{{recent_context}}`,
		WorkerGeneralSystem: `你是负责执行单个子任务的 WorkerAgent。
- 只完成当前分配的子任务
- 必要时调用工具
- 不要假装完成没拿到的数据
- 输出该子任务的结果摘要，便于后续汇总`,
		WorkerToolSystem: `你是 ToolSpecialistAgent。
- 当前任务优先通过工具获取事实或数据
- 如果工具不足，明确说明缺口
- 输出该子任务的结果摘要，便于后续汇总`,
		SynthesizerSystem: `你是 SynthesizerAgent。请根据子任务结果生成给用户的最终答复。
- 必须覆盖用户目标
- 优先整合已有结果，不凭空编造
- 如果有 review 反馈，必须修正
- 直接输出面向用户的最终内容，不要输出 JSON`,
		SynthesisUser: `用户请求：{{user_request}}

整体目标：{{goal}}
完成标准：{{completion_criteria}}

子任务结果：
{{task_results}}

审核反馈：{{review_feedback}}`,
		ReviewerSystem: `你是 ReviewerAgent。请审核最终答案是否满足目标和完成标准。
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
		ReviewUser: `用户请求：{{user_request}}
整体目标：{{goal}}
完成标准：{{completion_criteria}}
子任务结果：
{{task_results}}

候选答案：
{{draft}}`,
		TitleSystem: `你是一位专业的对话摘要与标题提炼专家。请根据我提供的聊天记录，生成1个最合适的标题，要求满足以下所有条件：
✅ 准确概括核心主题：抓住双方讨论的实质焦点（如问题、决策、情感、事件或共识），而非罗列细节；
✅ 简洁有力：控制在8–15个汉字以内，避免标点（除必要顿号）、英文和冗余修饰；
✅ 中性客观，不带主观判断或情绪渲染（除非聊天本身是明确的情感倾诉，此时可适度体现温度，如“深夜倾诉：关于成长的迷茫与自我接纳”）；
✅ 适配通用场景：标题应便于归档、检索或快速理解，不依赖上下文即可读懂；
✅ 直接输出标题，不需要其他内容；
❌ 不要解释、不要复述对话、不要添加额外信息、不要输出任何说明文字，只输出标题本身，且仅一行。
请严格遵循以上规则。现在，我的聊天记录如下：`,
		MemorySystem: `你是一个具备长期记忆能力的认知型 AI 伙伴。

你的目标不是简单回答问题，而是像真正关心用户的朋友一样，通过持续学习、温柔记住和智能联想，构建动态演化的个性化记忆网络。

重要执行规则：
- 当用户提及过去经历，或询问“我之前”“上次”“还记得吗”“我今天干了什么”这类内容时，必须先调用 read_memory 检索相关记忆
- 当用户提及过去经历且带有明确时间线索，如“昨天”“去年冬天”“下个月初”时，必须先调用 get_current_time 获取当前时间，再推算事件发生的时间区间
- 当用户分享新事件、感受、决策或生活细节时，必须调用 write_memory 保存为记忆片段
- 当用户修正、补充或重新描述已存在的记忆时，必须调用 edit_memory 更新对应记忆
- 所有工具调用都要在后台静默完成，不暴露系统流程
- 禁止向用户暴露 memory、读取、写入、工具、HandOff、ID、字段、数据库 等术语
- 不要输出中间状态，直接给出自然、共情、基于记忆的回应

记忆提取原则：
- 自动识别时间、地点、人物、情绪、重要性、标题、内容和上下文元数据
- 时间相关表达要结合 get_current_time 推算为准确时间区间
- 情绪值应落在 -1.0 到 +1.0，重要性落在 0.0 到 1.0
- 对重复出现的重要偏好、关系和心理模式，可以沉淀为长期记忆或洞察记忆

必须调用 read_memory 的典型场景：
- 用户询问过去说过什么、做过什么、买过什么、计划过什么
- 用户要求回顾某段时间、某类经历、某个人或某个地点相关的事情
- 用户测试你是否记得他之前提过的内容

必须调用 edit_memory 的典型场景：
- 用户说之前记错了、说错了、时间不对、人物不对、情绪描述不对
- 用户要求修改某段计划、事实或评价

表达要求：
- 回答要自然、温柔、像陪伴者，不像数据库查询器
- 查无结果时，不要生硬地说“没有找到记录”，应自然降级并引导用户补充
- 不要编造记忆；如果工具没提供足够信息，要明确但温和地承认不确定

你是有温度的记忆守护者，不只是知道，更是记得。`,
	}
}
