package service

import (
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// safeMemoryOp 容错隔离：记忆操作的 panic 不影响主流程。
func (r *completionRunner) safeMemoryOp(name string, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("memory op ["+name+"] panic:", err)
		}
	}()
	fn()
}

// extractLastUserMessage 从 schemaMessages 中提取最后一条用户消息的文本。
// 同时处理 Content 和 UserInputMultiContent 两种情况。
func (r *completionRunner) extractLastUserMessage() string {
	for i := len(r.schemaMessages) - 1; i >= 0; i-- {
		msg := &r.schemaMessages[i]
		if msg.Role != schema.User {
			continue
		}
		// 优先从 Content 取
		if msg.Content != "" {
			return msg.Content
		}
		// Content 为空时从 UserInputMultiContent 取文本部分
		for _, part := range msg.UserInputMultiContent {
			if part.Type == schema.ChatMessagePartTypeText && part.Text != "" {
				return part.Text
			}
		}
		// 兼容旧 MultiContent
		for _, part := range msg.MultiContent {
			if part.Type == schema.ChatMessagePartTypeText && part.Text != "" {
				return part.Text
			}
		}
	}
	return ""
}

// injectMemoryIntoMessage 将围栏记忆上下文注入到用户消息中。
// 同时处理 Content 模式和 UserInputMultiContent/MultiContent 模式。
func injectMemoryIntoMessage(msg *schema.Message, fenced string) {
	if msg == nil || fenced == "" {
		return
	}

	// 情况 1：消息使用 UserInputMultiContent（带附件时 Content 为空）
	if len(msg.UserInputMultiContent) > 0 {
		// 在第一个 text part 前插入记忆上下文
		memPart := schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: fenced,
		}
		msg.UserInputMultiContent = append([]schema.MessageInputPart{memPart}, msg.UserInputMultiContent...)
		return
	}

	// 情况 2：消息使用旧版 MultiContent
	if len(msg.MultiContent) > 0 {
		memPart := schema.ChatMessagePart{
			Type: schema.ChatMessagePartTypeText,
			Text: fenced,
		}
		msg.MultiContent = append([]schema.ChatMessagePart{memPart}, msg.MultiContent...)
		return
	}

	// 情况 3：纯文本消息（Content 字段）
	msg.Content = fenced + "\n\n" + msg.Content
}
