package skills

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// skillNameRe enforces kebab-case names: lowercase alphanumerics and hyphens, max 64 chars.
var skillNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

var (
	// ErrInvalidSkillName indicates that a skill name is missing or not kebab-case.
	ErrInvalidSkillName = errors.New("skill name must be kebab-case (a-z0-9- only, max 64 chars)")
	// ErrInvalidSkillContent indicates that required skill content is missing.
	ErrInvalidSkillContent = errors.New("skill description is required")
	// ErrSkillNameTaken indicates another skill already uses the requested name.
	ErrSkillNameTaken = errors.New("skill name is already taken")
)

// Frontmatter is the YAML header at the top of every SKILL.md.
type Frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Parse parses a SKILL.md byte slice into a Skill (without Source / Disabled / UpdatedAt set).
// The caller is responsible for filling those fields from the loader context.
func Parse(raw []byte) (Skill, error) {
	// Expect leading "---\n" delimiter.
	const delim = "---"
	trimmed := bytes.TrimLeft(raw, "\r\n ")
	if !bytes.HasPrefix(trimmed, []byte(delim)) {
		return Skill{}, errors.New("missing YAML frontmatter")
	}
	// Skip past the opening delimiter.
	body := trimmed[len(delim):]
	// Skip optional newline after first delim.
	body = bytes.TrimLeft(body, "\r\n")
	// Locate the closing delimiter on its own line.
	end := bytes.Index(body, []byte("\n"+delim))
	if end < 0 {
		return Skill{}, errors.New("frontmatter not closed")
	}
	headerBytes := body[:end]
	rest := body[end+len("\n"+delim):]
	rest = bytes.TrimLeft(rest, "\r\n")

	// Decode the YAML frontmatter block.
	var fm Frontmatter
	if err := yaml.Unmarshal(headerBytes, &fm); err != nil {
		return Skill{}, err
	}
	fm.Name = strings.TrimSpace(fm.Name)
	fm.Description = strings.TrimSpace(fm.Description)
	if fm.Name == "" {
		return Skill{}, ErrInvalidSkillName
	}
	if fm.Description == "" {
		return Skill{}, ErrInvalidSkillContent
	}
	if !skillNameRe.MatchString(fm.Name) {
		return Skill{}, ErrInvalidSkillName
	}
	return Skill{
		Name:        fm.Name,
		Description: fm.Description,
		Body:        string(rest),
	}, nil
}

// IsValidName reports whether name satisfies the Skill kebab-case naming rules.
func IsValidName(name string) bool {
	return skillNameRe.MatchString(strings.TrimSpace(name))
}

// Validate checks the fields required for a persisted Skill.
func Validate(s Skill) error {
	if !IsValidName(s.Name) {
		return ErrInvalidSkillName
	}
	if strings.TrimSpace(s.Description) == "" {
		return ErrInvalidSkillContent
	}
	return nil
}

// Render produces the on-disk SKILL.md bytes (frontmatter + body).
func Render(s Skill) ([]byte, error) {
	header, err := yaml.Marshal(Frontmatter{Name: s.Name, Description: s.Description})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(header)
	buf.WriteString("---\n\n")
	buf.WriteString(strings.TrimRight(s.Body, "\n"))
	buf.WriteString("\n")
	return buf.Bytes(), nil
}
