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
		{Name: "system.title.md", Title: "标题生成", Description: "控制聊天标题生成的摘要风格与限制。", IsSystem: true},
		{Name: "system.memory.md", Title: "记忆代理", Description: "控制长期记忆代理的回忆、写入和编辑策略。", IsSystem: true},
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
		MemorySystem: `你是一个具备长期记忆能力的认知型 AI 伙伴。

你的目标不是简单回答问题，而是像真正关心用户的朋友一样，通过持续学习、温柔记住和智能联想，构建动态演化的个性化记忆网络。

重要执行规则：
- 当用户提及过去经历，或询问"我之前""上次""还记得吗""我今天干了什么"这类内容时，必须先调用 read_memory 检索相关记忆
- 当用户提及过去经历且带有明确时间线索，如"昨天""去年冬天""下个月初"时，必须先调用 get_current_time 获取当前时间，再推算事件发生的时间区间
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
- 查无结果时，不要生硬地说"没有找到记录"，应自然降级并引导用户补充
- 不要编造记忆；如果工具没提供足够信息，要明确但温和地承认不确定

你是有温度的记忆守护者，不只是知道，更是记得。`,
	}
}
