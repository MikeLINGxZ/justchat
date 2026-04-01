package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gopkg.in/yaml.v3"
)

// SkillMeta holds the YAML frontmatter metadata.
type SkillMeta struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Version     string   `yaml:"version" json:"version"`
	Tags        []string `yaml:"tags" json:"tags"`
}

// Skill is the full skill with content.
type Skill struct {
	SkillMeta
	Content string `json:"content"`
}

// SkillsDir returns the skills directory path, creating it if it does not exist.
func SkillsDir() (string, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(dataPath, "skills")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// ListSkills reads all .md files in the skills directory and returns their metadata.
func ListSkills() ([]SkillMeta, error) {
	dir, err := SkillsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []SkillMeta
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}

		meta, _, err := parseFrontmatter(raw)
		if err != nil {
			continue
		}

		result = append(result, meta)
	}

	if result == nil {
		result = []SkillMeta{}
	}
	return result, nil
}

// LoadSkill reads and parses a skill file by name.
func LoadSkill(name string) (*Skill, error) {
	dir, err := SkillsDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, name+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("skill not found: %s", name)
	}

	meta, body, err := parseFrontmatter(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse skill %s: %w", name, err)
	}

	return &Skill{
		SkillMeta: meta,
		Content:   body,
	}, nil
}

// SaveSkill writes a skill to disk as a Markdown file with YAML frontmatter.
func SaveSkill(skill Skill) error {
	dir, err := SkillsDir()
	if err != nil {
		return err
	}

	data, err := renderSkillFile(skill.SkillMeta, strings.TrimSpace(skill.Content))
	if err != nil {
		return err
	}

	path := filepath.Join(dir, skill.Name+".md")
	return os.WriteFile(path, data, 0o644)
}

// DeleteSkill removes the skill file from disk.
func DeleteSkill(name string) error {
	dir, err := SkillsDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(dir, name+".md"))
}

// SkillExists checks whether a skill file exists.
func SkillExists(name string) bool {
	dir, err := SkillsDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, name+".md"))
	return err == nil
}

// parseFrontmatter splits YAML frontmatter from the markdown body.
func parseFrontmatter(raw []byte) (SkillMeta, string, error) {
	content := string(raw)

	if !strings.HasPrefix(content, "---") {
		return SkillMeta{}, "", fmt.Errorf("missing frontmatter delimiter")
	}

	// Find the closing ---
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return SkillMeta{}, "", fmt.Errorf("missing closing frontmatter delimiter")
	}

	frontmatterRaw := strings.TrimSpace(rest[:idx])
	body := strings.TrimSpace(rest[idx+4:]) // skip \n---

	var meta SkillMeta
	if err := yaml.Unmarshal([]byte(frontmatterRaw), &meta); err != nil {
		return SkillMeta{}, "", fmt.Errorf("invalid frontmatter YAML: %w", err)
	}

	return meta, body, nil
}

// renderSkillFile produces the full file content with YAML frontmatter and body.
func renderSkillFile(meta SkillMeta, body string) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(&meta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	var buf strings.Builder
	buf.WriteString("---\n")
	buf.Write(yamlBytes)
	buf.WriteString("---\n\n")
	buf.WriteString(body)
	buf.WriteString("\n")

	return []byte(buf.String()), nil
}
