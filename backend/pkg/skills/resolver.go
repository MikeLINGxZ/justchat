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

// ResolveSkillSummaries builds a prompt section listing available skills
// with their name, description, and trigger condition, instructing the LLM
// to call load_skill when needed.
func ResolveSkillSummaries(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}

	var lines []string
	for _, name := range skillNames {
		skill, err := LoadSkill(name)
		if err != nil {
			continue
		}

		desc := strings.TrimSpace(skill.Description)
		when := strings.TrimSpace(skill.When)
		if desc == "" {
			continue
		}

		line := "- **" + skill.Name + "**: " + desc
		if when != "" {
			line += " (when: " + when + ")"
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return ""
	}

	return "## Available Skills\n\n" +
		"You have the following skills available. When a user's request matches a skill's trigger condition, " +
		"call the `load_skill` tool with the skill name to load its full content before responding.\n\n" +
		strings.Join(lines, "\n")
}

// ResolveAllSkillSummaries loads all skills from disk and builds a summary prompt section.
func ResolveAllSkillSummaries() string {
	metas, err := ListSkills()
	if err != nil || len(metas) == 0 {
		return ""
	}

	var names []string
	for _, m := range metas {
		names = append(names, m.Name)
	}
	return ResolveSkillSummaries(names)
}
