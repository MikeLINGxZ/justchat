package llm_opc

import (
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

// ResolveInstruction 解析指令文本（提示词 + 技能内容）
func ResolveInstruction(person IOpcPerson) string {
	instruction := person.Prompt()
	if len(person.Skills()) > 0 {
		skillContent := skills.ResolveSkillContents(person.Skills())
		if skillContent != "" {
			instruction = instruction + "\n\n" + skillContent
		}
	}
	return instruction
}

// BuildPersonSystemPrompt 构建个人聊天的系统提示词
func BuildPersonSystemPrompt(person IOpcPerson) string {
	instruction := ResolveInstruction(person)
	return fmt.Sprintf(
		"你是%s，你的职责是：%s\n\n%s\n\n请直接回复，不要在回复中提到你是AI。保持专业且简洁。",
		person.Name(), person.Duties(), instruction,
	)
}

// BuildGroupSystemPrompt 构建群聊回复的系统提示词（不含成员名单，改用工具获取）
func BuildGroupSystemPrompt(person IOpcPerson) string {
	instruction := ResolveInstruction(person)
	return fmt.Sprintf(
		"你是%s，职责是：%s。你正在一个群聊中。\n\n%s\n\n"+
			"你可以通过工具获取群聊信息和成员信息。"+
			"请直接回复，不要在回复中提到你是AI。保持专业且简洁。",
		person.Name(), person.Duties(), instruction,
	)
}

// BuildDecisionPrompt 构建群聊决策轮提示词（仍包含成员名单文本，用于单次 JSON 决策）
func BuildDecisionPrompt(person IOpcPerson, memberRoster, contextSummary string) string {
	return fmt.Sprintf(
		`你是%s，职责是：%s。

你所在的群聊中有以下成员：
%s
以下是群聊的最近消息：
%s

请判断你是否需要回复最新的消息。用JSON格式回答：
{"should_respond": true/false, "priority": 1-10, "reason": "简短原因"}

判断规则：
1. 如果消息明确提到你的名字或用 @%s 指定了你，你必须回复
2. 如果消息涉及的问题与你的职责直接相关，且群里没有比你更合适的人，应该回复
3. 如果群里有其他成员的职责比你更匹配这个问题，你不需要回复
4. 如果消息是闲聊、通知、或与你无关的话题，不需要回复
5. 不要为了表现存在感而回复，只在你能提供有价值信息时才回复
6. priority 越高表示越需要你回复（10最高）

只返回JSON，不要添加其他内容。`,
		person.Name(), person.Duties(), memberRoster, contextSummary, person.Name(),
	)
}

// BuildMemberRoster 从成员列表生成文本名单
func BuildMemberRoster(members []IOpcPerson) string {
	var roster strings.Builder
	for _, m := range members {
		roster.WriteString(fmt.Sprintf("- %s：%s\n", m.Name(), m.Duties()))
	}
	return roster.String()
}

// BuildContextSummary 从消息历史生成上下文摘要
// memberLookup 根据 SenderPersonUuid 查找发送者名称，找不到时返回空字符串
func BuildContextSummary(messages []data_models.Message, memberLookup func(uuid string) string) string {
	var summary strings.Builder
	for _, msg := range messages {
		var sender string
		if msg.Role == schema.User {
			sender = "用户"
		} else if msg.SenderPersonUuid != "" {
			sender = memberLookup(msg.SenderPersonUuid)
			if sender == "" {
				sender = "某成员"
			}
		} else {
			sender = "助手"
		}
		summary.WriteString(fmt.Sprintf("%s: %s\n", sender, truncateString(msg.Content, 200)))
	}
	return summary.String()
}

// BuildSchemaMessages 将历史消息转换为 schema.Message 切片（带系统提示词）
func BuildSchemaMessages(systemPrompt string, historyMessages []data_models.Message, skipMessageUuid string) []*schema.Message {
	var msgs []*schema.Message
	sysMsg := &schema.Message{
		Role:    schema.System,
		Content: systemPrompt,
	}
	msgs = append(msgs, sysMsg)

	for _, msg := range historyMessages {
		if skipMessageUuid != "" && msg.MessageUuid == skipMessageUuid {
			continue
		}
		schemaMsg, err := msg.ToSchemaMessage()
		if err != nil {
			continue
		}
		msgs = append(msgs, schemaMsg)
	}
	return msgs
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
