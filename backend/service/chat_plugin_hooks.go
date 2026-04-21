package service

import (
	"encoding/json"

	"github.com/cloudwego/eino/schema"
)

func schemaMessagesToHookMessages(msgs []schema.Message) []map[string]any {
	result := make([]map[string]any, len(msgs))
	for i, msg := range msgs {
		raw, err := json.Marshal(msg)
		if err != nil {
			result[i] = map[string]any{
				"role":    string(msg.Role),
				"content": msg.Content,
			}
			continue
		}

		var hookMsg map[string]any
		if err := json.Unmarshal(raw, &hookMsg); err != nil {
			result[i] = map[string]any{
				"role":    string(msg.Role),
				"content": msg.Content,
			}
			continue
		}
		result[i] = hookMsg
	}
	return result
}

func hookMessagesToSchemaMessages(msgs []map[string]any) []schema.Message {
	result := make([]schema.Message, len(msgs))
	for i, m := range msgs {
		raw, err := json.Marshal(m)
		if err != nil {
			role, _ := m["role"].(string)
			content, _ := m["content"].(string)
			result[i] = schema.Message{
				Role:    schema.RoleType(role),
				Content: content,
			}
			continue
		}

		var msg schema.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			role, _ := m["role"].(string)
			content, _ := m["content"].(string)
			result[i] = schema.Message{
				Role:    schema.RoleType(role),
				Content: content,
			}
			continue
		}
		result[i] = msg
	}
	return result
}
