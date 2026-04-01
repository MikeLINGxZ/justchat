package agents

import "fmt"

// CustomAgentDef represents a user-defined custom agent.
type CustomAgentDef struct {
	ID_         string   `json:"id"`
	DisplayName string   `json:"name"`
	Description string   `json:"description"`
	PromptText  string   `json:"prompt"`
	ToolIDs     []string `json:"tools"`
	SkillIDs    []string `json:"skills"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func (c *CustomAgentDef) Name() string    { return c.ID_ }
func (c *CustomAgentDef) Desc() string    { return c.DisplayName }
func (c *CustomAgentDef) Prompt() string  { return c.PromptText }
func (c *CustomAgentDef) Type() AgentType { return AgentTypeCustom }
func (c *CustomAgentDef) Role() AgentRole { return AgentRoleWorker }

func (c *CustomAgentDef) PromptNames() []string {
	return []string{fmt.Sprintf("custom.%s.md", c.ID_)}
}

func (c *CustomAgentDef) PromptMetas() []AgentPromptMeta {
	return []AgentPromptMeta{
		{
			FileName:    fmt.Sprintf("custom.%s.md", c.ID_),
			Title:       c.DisplayName,
			Description: c.Description,
			IsSystem:    false,
		},
	}
}

func (c *CustomAgentDef) DefaultPrompts() map[string]string {
	promptName := fmt.Sprintf("custom.%s.md", c.ID_)
	return map[string]string{promptName: c.PromptText}
}

// GetSkillNames implements ISkillCapableAgent.
func (c *CustomAgentDef) GetSkillNames() []string {
	return c.SkillIDs
}
