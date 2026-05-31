package agent

import (
	"encoding/json"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
)

const timeFormat = time.RFC3339

func toSessionItem(s data_models.Session) agent_dto.SessionItem {
	tags := sessionTags(s)
	return agent_dto.SessionItem{
		ID:      s.ID,
		Title:   s.Title,
		Kind:    s.Kind,
		Tags:    tags,
		Starred: s.Starred,
		Status:  s.Status,
		Created: s.CreatedAt.Format(timeFormat),
		Updated: s.UpdatedAt.Format(timeFormat),
	}
}

func sessionTags(s data_models.Session) []string {
	tags := parseSessionTags(s.Tags)
	if len(tags) == 0 && s.Kind == "task" {
		return []string{"task"}
	}
	return tags
}

func normalizeSessionTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	return out
}

func parseSessionTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return nil
	}
	return normalizeSessionTags(tags)
}

func marshalSessionTags(tags []string) string {
	tags = normalizeSessionTags(tags)
	if len(tags) == 0 {
		return ""
	}
	content, _ := json.Marshal(tags)
	return string(content)
}

func tagsContain(tags []string, target string) bool {
	target = strings.ToLower(strings.TrimSpace(target))
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

func toMessageItem(m data_models.Message) agent_dto.MessageItem {
	return agent_dto.MessageItem{
		ID:          m.ID,
		SessionID:   m.SessionID,
		ParentID:    m.ParentID,
		Role:        m.Role,
		ContentType: m.ContentType,
		Content:     m.Content,
		ModelName:   m.ModelName,
		AgentName:   m.AgentName,
		TokensIn:    m.TokensIn,
		TokensOut:   m.TokensOut,
		Extra:       m.Extra,
		Attachments: m.Attachments,
		Created:     m.CreatedAt.Format(timeFormat),
	}
}
