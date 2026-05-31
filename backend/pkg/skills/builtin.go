package skills

import (
	"embed"
	"io/fs"
	"path"
	"sort"
)

//go:embed all:_builtin
var builtinFS embed.FS

// LoadBuiltin returns all skills embedded in the binary, marked as SourceBuiltin.
// Each top-level subdir of _builtin/ that contains a SKILL.md becomes one skill.
// Subdirs without a parseable SKILL.md are silently skipped so that a single
// corrupt entry never blocks the rest of the library from loading.
func LoadBuiltin() ([]Skill, error) {
	entries, err := fs.ReadDir(builtinFS, "_builtin")
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Only directories containing a SKILL.md qualify as skills.
		mdPath := path.Join("_builtin", entry.Name(), SkillFileName)
		raw, err := fs.ReadFile(builtinFS, mdPath)
		if err != nil {
			continue
		}
		skill, err := Parse(raw)
		if err != nil {
			continue
		}
		skill.Source = SourceBuiltin
		skills = append(skills, skill)
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, nil
}
