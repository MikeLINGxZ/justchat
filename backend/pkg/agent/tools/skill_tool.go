package tools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// SkillProvider abstracts the skills manager so the agent can list and load skills
// without importing the concrete skills package directly.
type SkillProvider interface {
	// Enabled returns all non-disabled skills available to the agent.
	Enabled() []pkgskills.Skill
	// Get returns a single skill by name, or false if not found.
	Get(name string) (pkgskills.Skill, bool)
}

// SkillToolName is the canonical name of the meta-tool used by the agent.
const SkillToolName = "Skill"

// skillInput defines the JSON arguments accepted by the Skill tool.
type skillInput struct {
	Name string `json:"name" jsonschema:"description=Name of the skill to load,required"`
}

// BuildSkillTool constructs ToolMeta for the Skill meta-tool. The description is
// rebuilt each time this function is called so it always reflects the current set
// of enabled skills.
func BuildSkillTool(provider SkillProvider) ToolMeta {
	enabled := provider.Enabled()
	var sb strings.Builder
	sb.WriteString("Load a skill (a set of natural-language instructions for a specific task). Call with {\"name\": \"<skill-name>\"}. Available skills:\n")
	if len(enabled) == 0 {
		sb.WriteString("(none)\n")
	}
	for _, s := range enabled {
		sb.WriteString("- ")
		sb.WriteString(s.Name)
		sb.WriteString(": ")
		sb.WriteString(s.Description)
		sb.WriteString("\n")
	}
	return ToolMeta{
		Name:        SkillToolName,
		Description: sb.String(),
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed skillInput
			_ = json.Unmarshal(args, &parsed)
			return "Loading skill: " + parsed.Name
		},
	}
}

// InvokeSkill returns the SKILL.md body of the requested skill. It is used by the
// agent's tool-call handler when the model invokes Skill({name}).
func InvokeSkill(_ context.Context, provider SkillProvider, args json.RawMessage) (string, error) {
	var parsed skillInput
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.Name == "" {
		return "", errors.New("skill name is required")
	}
	sk, ok := provider.Get(parsed.Name)
	if !ok {
		return "", errors.New("skill not found: " + parsed.Name)
	}
	if sk.Disabled {
		return "", errors.New("skill is disabled: " + parsed.Name)
	}
	return sk.Body, nil
}

// NewSkillTool creates the Skill meta-tool as a function tool for the LLM agent.
// The tool's description is built from the provider's current enabled skills so the
// model always sees an up-to-date list.
func NewSkillTool(provider SkillProvider) *function.FunctionTool[skillInput, string] {
	meta := BuildSkillTool(provider)
	return function.NewFunctionTool(
		func(_ context.Context, input skillInput) (string, error) {
			if input.Name == "" {
				return "", errors.New("skill name is required")
			}
			sk, ok := provider.Get(input.Name)
			if !ok {
				return "", errors.New("skill not found: " + input.Name)
			}
			if sk.Disabled {
				return "", errors.New("skill is disabled: " + input.Name)
			}
			return sk.Body, nil
		},
		function.WithName(SkillToolName),
		function.WithDescription(meta.Description),
	)
}
