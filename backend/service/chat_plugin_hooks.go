package service

import "github.com/cloudwego/eino/schema"

func schemaMessagesToHookMessages(msgs []schema.Message) []map[string]any {
	result := make([]map[string]any, len(msgs))
	for i, msg := range msgs {
		result[i] = map[string]any{
			"role":    string(msg.Role),
			"content": msg.Content,
		}
	}
	return result
}

func hookMessagesToSchemaMessages(msgs []map[string]any) []schema.Message {
	result := make([]schema.Message, len(msgs))
	for i, m := range msgs {
		role, _ := m["role"].(string)
		content, _ := m["content"].(string)
		result[i] = schema.Message{
			Role:    schema.RoleType(role),
			Content: content,
		}
	}
	return result
}
