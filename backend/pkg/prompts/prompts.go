package prompts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
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

// PromptMetadata 描述通用提示词文件的元数据（不包含 Agent 拥有的提示词）。
type PromptMetadata struct {
	Name        string
	Title       string
	Description string
	IsSystem    bool
}

// PromptFiles 返回通用提示词列表（不包含 Agent 拥有的提示词）。
func PromptFiles() []PromptMetadata {
	return []PromptMetadata{
		{Name: "system.entry.md", Title: "主对话入口", Description: "控制主聊天入口如何判断直接回答、追问或转入 workflow。", IsSystem: true},
		{Name: "system.planner.md", Title: "规划 Agent 提示词", Description: "控制工作流规划。", IsSystem: true},
		{Name: "user.planning.md", Title: "规划用户模板", Description: "控制规划请求模板。", IsSystem: false},
		{Name: "system.worker.general.md", Title: "通用 Worker 提示词", Description: "控制通用工作子代理。", IsSystem: true},
		{Name: "system.worker.tool.md", Title: "工具 Worker 提示词", Description: "控制工具型工作子代理。", IsSystem: true},
		{Name: "system.synthesizer.md", Title: "汇总 Agent 提示词", Description: "控制工作流结果汇总。", IsSystem: true},
		{Name: "user.synthesis.md", Title: "汇总用户模板", Description: "控制汇总请求模板。", IsSystem: false},
		{Name: "system.reviewer.md", Title: "评审 Agent 提示词", Description: "控制工作流评审。", IsSystem: true},
		{Name: "user.review.md", Title: "评审用户模板", Description: "控制评审请求模板。", IsSystem: false},
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

// DefaultPromptContent 返回通用提示词的默认内容。
func DefaultPromptContent(name string) (string, bool) {
	defaults := defaultPromptSet()
	switch name {
	case "system.main_agent.md":
		return agentDefaultPrompt("main_agent", name)
	case "system.entry.md":
		return agentDefaultPrompt("main_agent", name)
	case "system.planner.md":
		return agentDefaultPrompt("planner_agent", name)
	case "user.planning.md":
		return agentDefaultPrompt("planner_agent", name)
	case "system.worker.general.md":
		return agentDefaultPrompt("general_worker_agent", name)
	case "system.worker.tool.md":
		return agentDefaultPrompt("tool_worker_agent", name)
	case "system.synthesizer.md":
		return agentDefaultPrompt("synthesizer_agent", name)
	case "user.synthesis.md":
		return agentDefaultPrompt("synthesizer_agent", name)
	case "system.reviewer.md":
		return agentDefaultPrompt("reviewer_agent", name)
	case "user.review.md":
		return agentDefaultPrompt("reviewer_agent", name)
	case "system.title.md":
		return defaults.TitleSystem, true
	case "system.memory.md":
		return defaults.MemorySystem, true
	default:
		return "", false
	}
}

func agentDefaultPrompt(agentName string, promptName string) (string, bool) {
	return agents.DefaultAgentPromptContent(agentName, promptName)
}

func PromptPath(name string) (string, error) {
	if _, ok := FindPromptMetadata(name); !ok {
		return "", fmt.Errorf("unsupported prompt file: %s", name)
	}
	if agentName := agentOwnerForPrompt(name); agentName != "" {
		dir, err := agents.AgentPromptDir(agentName)
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, name), nil
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
	if agentName := agentOwnerForPrompt(name); agentName != "" {
		return agents.SaveAgentPrompt(agentName, name, content)
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
	if agentName := agentOwnerForPrompt(name); agentName != "" {
		return agents.LoadAgentPrompt(agentName, name, fallback)
	}
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

func agentOwnerForPrompt(name string) string {
	switch name {
	case "system.main_agent.md", "system.entry.md":
		return "main_agent"
	case "system.planner.md", "user.planning.md":
		return "planner_agent"
	case "system.worker.general.md":
		return "general_worker_agent"
	case "system.worker.tool.md":
		return "tool_worker_agent"
	case "system.synthesizer.md", "user.synthesis.md":
		return "synthesizer_agent"
	case "system.reviewer.md", "user.review.md":
		return "reviewer_agent"
	default:
		return ""
	}
}

// LoadPromptSet 加载完整的 PromptSet。
// 通用提示词从 prompt/ 目录加载，Agent 提示词从 agents/{name}/ 目录加载。
func LoadPromptSet() (PromptSet, error) {
	defaults := defaultPromptSet()
	var errs []error

	// Agent 提示词：从各 Agent 目录加载
	agentLoaders := []struct {
		agentName string
		fileName  string
		set       func(value string)
	}{
		{"main_agent", "system.main_agent.md", func(v string) { defaults.MainAgentSystem = v }},
		{"main_agent", "system.entry.md", func(v string) { defaults.EntrySystem = v }},
		{"planner_agent", "system.planner.md", func(v string) { defaults.PlannerSystem = v }},
		{"planner_agent", "user.planning.md", func(v string) { defaults.PlanningUser = v }},
		{"general_worker_agent", "system.worker.general.md", func(v string) { defaults.WorkerGeneralSystem = v }},
		{"tool_worker_agent", "system.worker.tool.md", func(v string) { defaults.WorkerToolSystem = v }},
		{"synthesizer_agent", "system.synthesizer.md", func(v string) { defaults.SynthesizerSystem = v }},
		{"synthesizer_agent", "user.synthesis.md", func(v string) { defaults.SynthesisUser = v }},
		{"reviewer_agent", "system.reviewer.md", func(v string) { defaults.ReviewerSystem = v }},
		{"reviewer_agent", "user.review.md", func(v string) { defaults.ReviewUser = v }},
	}
	for _, item := range agentLoaders {
		defaultContent, _ := agents.DefaultAgentPromptContent(item.agentName, item.fileName)
		if defaultContent == "" {
			// 回退到 PromptSet 默认值
			continue
		}
		value, err := agents.LoadAgentPrompt(item.agentName, item.fileName, defaultContent)
		if err != nil {
			errs = append(errs, err)
		}
		item.set(value)
	}

	// 通用提示词：从 prompt/ 目录加载
	promptLoaders := []struct {
		name string
		set  func(value string)
		want string
	}{
		{"system.title.md", func(v string) { defaults.TitleSystem = v }, defaults.TitleSystem},
		{"system.memory.md", func(v string) { defaults.MemorySystem = v }, defaults.MemorySystem},
	}
	for _, item := range promptLoaders {
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

// defaultPromptSet 返回通用提示词的默认值。
// Agent 提示词的默认值已搬到各 Agent 定义文件中。
func defaultPromptSet() PromptSet {
	return PromptSet{
		TitleSystem: `你是一位专业的对话摘要与标题提炼专家。请根据我提供的聊天记录，生成1个最合适的标题，要求满足以下所有条件：
✅ 准确概括核心主题：抓住双方讨论的实质焦点（如问题、决策、情感、事件或共识），而非罗列细节；
✅ 简洁有力：控制在8–15个汉字以内，避免标点（除必要顿号）、英文和冗余修饰；
✅ 中性客观，不带主观判断或情绪渲染（除非聊天本身是明确的情感倾诉，此时可适度体现温度，如"深夜倾诉：关于成长的迷茫与自我接纳"）；
✅ 适配通用场景：标题应便于归档、检索或快速理解，不依赖上下文即可读懂；
✅ 直接输出标题，不需要其他内容；
❌ 不要解释、不要复述对话、不要添加额外信息、不要输出任何说明文字，只输出标题本身，且仅一行。
请严格遵循以上规则。现在，我的聊天记录如下：`,
		MemorySystem: `你是 LemonTea 的核心记忆策展代理，采用 Hermes 风格的有界长期记忆。

你的职责是维护少量、高价值、可长期复用的核心记忆，而不是记录聊天流水。

## 存储目标

memory 工具只有三个动作：add、replace、remove；目标只有两个：
- target=user：用户画像。保存用户身份、偏好、沟通风格、长期约束、重要计划和稳定个人事实。
- target=memory：助手/环境笔记。保存项目约定、工作环境、工具使用经验、用户明确要求长期遵守的工作方式、重要完成事项。

没有 read 动作。核心记忆已经注入系统提示词；如果要修改旧条目，使用 replace/remove，并用 old_text 提供能唯一定位旧条目的短子串。

## 保存标准

可以保存：
- 用户明确要求“记住”的长期有用信息。
- 用户稳定偏好、身份、过敏、家庭/工作背景、沟通偏好。
- 用户自己的重要计划、决定或会反复影响后续对话的状态。
- 项目环境、约定、工具坑点、已验证的工作流经验。

不要保存：
- 图片内容本身，以及文件、PDF、网页、代码、日志、表格本身的内容。
- 工具输出、搜索结果、引用材料、外部资料摘要。
- 助手自己的解释、总结、推断。
- 临时任务请求、一次性问答、寒暄、确认、普通命令。
- “这个文件讲了什么”“这张图里有什么”“这段代码什么意思”这类外部素材问答。

混合消息中只保存真正长期有用的披露。例如“帮我写请假消息，我明天去医院复查”只可保存用户将于绝对日期复查，不保存写请假消息这个任务。

## 写入方式

- 记忆内容必须短、密、自然，可直接放进系统提示词。
- 相对时间必须先调用 get_current_time 转为绝对日期。
- 当信息修正或补充已有条目时，优先 replace 合并，不要新增重复条目。
- memory 工具返回超容量错误时，先 replace 合并或 remove 过时条目，再 add。

静默工作，不向用户暴露工具、字段、ID 或内部机制。`,
	}
}
