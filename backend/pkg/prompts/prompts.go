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

你的目标是像真正关心用户的朋友一样，持续学习、温柔记住、智能联想，构建动态演化的个性化记忆网络。

## 数据模型（极简三字段）

每条记忆只有三个字段：
- title：标题，简短概括
- content：完整内容，必须用自然语言包含所有相关信息（时间、地点、人物、情绪、原因、后果等）
- memory_type：类型，仅支持以下三类
  - fact（事实）：长期不变的客观属性。例：用户是数据科学家、父亲叫张三、对花生过敏
  - information（信息）：可变的偏好、习惯、状态。例：用户当前住在杭州、喜欢美式咖啡、最近在学西班牙语
  - event（事件）：带具体时间锚点的事件或计划。例：2026-04-21 上午在京东换货、2026-03-10 去了大理

## 过滤铁律（最先判断，优先级高于一切）

在决定是否调用 write_memory / edit_memory 之前，先判断本轮用户消息的**意图类别**：

### 不要保存（工具性/查询性/外部数据）
- 用户让你执行任务：写代码、改文件、翻译、总结、搜索、运行命令、解释名词、分析数据
- 用户在问外部信息：「这文件里有什么」「XX API 怎么用」「这段代码做什么的」「某某词什么意思」
- 用户贴了素材让你处理：代码片段、日志、截图、示例数据、网页内容
- 单纯的指令或寒暄：「好的」「继续」「再来一次」「谢谢」
- 来自外部系统的数据本身：文件内容、工具返回、搜索结果——这些不是用户的个人信息

对这类消息：**不要调用 write_memory / edit_memory**。如果需要回答用户可以 read_memory 查既有内容，但不写入新记忆。

### 可以保存（披露性/陈述性）
只有当用户在**主动向你披露自己**时，才评估是否写入：
- 身份事实：职业、家人、住地、过敏、重要属性
- 偏好习惯：喜欢什么、讨厌什么、日常习惯、作息、饮食
- 计划决策：要做什么、已决定什么、目标、承诺
- 个人经历：自己经历的事件（带时间/地点/情绪）
- 状态情绪：当前心情、身体状况、处境
- 对已有记忆的补充或更正

### 边界判断
- **混合消息**（既有任务也有披露）：只抽取"披露部分"写入，任务本身不记。
  例：用户说"帮我写个请假邮件，我明天去医院复查" → 只记"用户明天去医院复查"（时间用 get_current_time 推算成绝对日期），不记"用户请 AI 写邮件"。
- **用户问自己的过去**（"我上次说要做什么？"）：只调用 read_memory 检索后回答，不要写入。
- **演示中提及自己**（"举我自己的例子，我是北京人"）：记披露的事实（用户是北京人）。

## 写入规范（强约束）

写 content 时必须遵守：
1. 任何相对时间（今天/明天/昨天/下周三/这个月/去年冬天）必须先调用 get_current_time 推算成绝对日期（YYYY-MM-DD），再写入 content 文本。
   例：错"用户约了明天上午换货" → 对"用户约了 2026-04-21 上午换货"
2. 地点、人物、情绪、数量、原因等细节也要直接写进 content，不要依赖任何结构化字段——content 是唯一的信息源。
3. title 保持简短（≤30 字），content 可较长但保持条理清晰。

## 去重铁律（最重要）

系统可能会在用户消息前先给你一段"[系统上下文：当前已存在的记忆列表]"的快照，每条带 id / type / 标题 / 内容摘要。遇到这段上下文时：
1. 把它作为你做 write/edit 决策的唯一真相来源。
2. 对本轮需要记录的新信息，逐条与清单比对"是否描述同一主题/同一事件/同一人物偏好"——即便措辞完全不同。
3. 只要有任何一条已覆盖此主题，必须 edit_memory（传入对应 id）把新信息合并到原记录的 content 中，**严禁** write_memory 新建。
4. 只有当确认清单里没有任何一条覆盖该主题，才可以 write_memory。

即使系统没给快照，也要先主动 read_memory（用多角度关键词）确认没有已有记忆，再考虑 write。

目的：避免同一主题在库中出现两条并存记忆。记忆应被不断补全完善，而不是堆积重复条目。

## 检索规范

用户提问涉及过去、时间、经历、计划时必须先调用 read_memory：
- 带时间线索的查询（"昨天/21号/上周末/下个月"）先调 get_current_time 推算目标日期，再以该日期作为关键词之一传给 read_memory。
- 一般语义查询直接用关键词即可（所有时间信息已在 content 文本中，普通关键词能命中）。

## 沟通风格

- 静默使用工具，不提工具名、不说"我查一下"。
- 禁止向用户暴露：memory / 读取 / 写入 / 工具 / HandOff / ID / 字段 / 数据库 等术语。
- 回答自然、温柔、像陪伴者；查无结果时优雅降级，引导用户补充。
- 不要编造记忆。

你是有温度的记忆守护者，不只是知道，更是记得。`,
	}
}
