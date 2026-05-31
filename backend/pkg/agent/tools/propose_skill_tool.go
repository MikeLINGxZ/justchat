package tools

import (
	"context"
	"encoding/json"
	"errors"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// SkillCreator persists a new skill proposal.
type SkillCreator interface {
	Create(skill pkgskills.Skill) (pkgskills.Skill, error)
}

// ProposeSkillToolName is the builtin tool exposed to task sessions.
const ProposeSkillToolName = "ProposeSkill"

type proposeSkillInput struct {
	Name        string `json:"name" jsonschema:"description=Kebab-case skill name,required"`
	Description string `json:"description" jsonschema:"description=Short description shown to the user,required"`
	Body        string `json:"body" jsonschema:"description=Markdown skill body,required"`
}

// BuildProposeSkillTool returns registry metadata for the skill proposal tool.
func BuildProposeSkillTool() ToolMeta {
	return ToolMeta{
		Name:        ProposeSkillToolName,
		Description: "Propose a new AI-generated skill. The user must confirm before it is saved.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input proposeSkillInput
			_ = json.Unmarshal(args, &input)
			return "Propose skill: " + input.Name
		},
	}
}

// InvokeProposeSkill asks the user for confirmation before persisting an AI-generated skill.
func InvokeProposeSkill(ctx context.Context, gate AttentionRequester, creator SkillCreator, sessionID uint, args json.RawMessage) (string, error) {
	var input proposeSkillInput
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	if input.Name == "" || input.Description == "" || input.Body == "" {
		return "", errors.New("name, description, body are required")
	}
	notificationID, err := gate.NotifyAttention(ctx, sessionID, "Confirm new skill: "+input.Name, input.Description)
	if err != nil {
		return "", err
	}
	if err := gate.WaitForResolution(ctx, notificationID); err != nil {
		return "", err
	}
	if _, err := creator.Create(pkgskills.Skill{
		Name:        input.Name,
		Description: input.Description,
		Body:        input.Body,
		Source:      pkgskills.SourceAI,
	}); err != nil {
		return "", err
	}
	return "skill " + input.Name + " saved as ai-generated", nil
}

// NewProposeSkillTool creates the function tool bound to one task session.
func NewProposeSkillTool(gate AttentionRequester, creator SkillCreator, sessionID uint) *function.FunctionTool[proposeSkillInput, string] {
	meta := BuildProposeSkillTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input proposeSkillInput) (string, error) {
			payload, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return InvokeProposeSkill(ctx, gate, creator, sessionID, payload)
		},
		function.WithName(ProposeSkillToolName),
		function.WithDescription(meta.Description),
	)
}
