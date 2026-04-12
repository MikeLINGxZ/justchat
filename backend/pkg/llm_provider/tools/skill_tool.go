package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

type LoadSkillTool struct{}

type loadSkillParams struct {
	SkillName string `json:"skill_name"`
}

func (l *LoadSkillTool) Id() string {
	return "load_skill"
}

func (l *LoadSkillTool) Name() string {
	return i18n.TCurrent("tool.load_skill.name", nil)
}

func (l *LoadSkillTool) Description() string {
	return i18n.TCurrent("tool.load_skill.description", nil)
}

func (l *LoadSkillTool) RequireConfirmation() bool { return false }

func (l *LoadSkillTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "load_skill",
			Desc: i18n.TCurrent("tool.load_skill.description", nil),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"skill_name": {
					Type:     schema.String,
					Desc:     "The name of the skill to load.",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params loadSkillParams) (string, error) {
			name := strings.TrimSpace(params.SkillName)
			if name == "" {
				return "", fmt.Errorf("skill_name is required")
			}

			skill, err := skills.LoadSkill(name)
			if err != nil {
				return "", fmt.Errorf("skill not found: %s", name)
			}

			content := strings.TrimSpace(skill.Content)
			if content == "" {
				return "", fmt.Errorf("skill %s has no content", name)
			}

			return fmt.Sprintf("# Skill: %s\n\n%s", skill.Name, content), nil
		},
	)
}
