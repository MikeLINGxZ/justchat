package skills

import "strings"

// ResolveSkillContents loads the specified skills and concatenates their content
// into a structured prompt section for injection into agent instructions.
func ResolveSkillContents(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}

	var sections []string
	for _, name := range skillNames {
		skill, err := LoadSkill(name)
		if err != nil {
			continue
		}

		content := strings.TrimSpace(skill.Content)
		if content == "" {
			continue
		}

		sections = append(sections, "### "+skill.Name+"\n"+content)
	}

	if len(sections) == 0 {
		return ""
	}

	return "## Skills\n\n" + strings.Join(sections, "\n\n")
}
