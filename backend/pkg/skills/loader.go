package skills

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	// SkillFileName is the conventional filename for a skill's markdown manifest.
	SkillFileName = "SKILL.md"
	// SidecarFileName is the JSON sidecar stored alongside each on-disk skill.
	SidecarFileName = ".lemontea.json"
)

// LoadFromDir scans root for skill subdirectories and returns them sorted by name.
// Subdirs without a parseable SKILL.md are silently skipped so that a single
// corrupt entry never blocks the rest of the library from loading.
func LoadFromDir(root string) ([]Skill, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		// A missing root simply means "no user skills yet" — not an error.
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		skill, err := loadOneDir(dir)
		if err != nil {
			// Skip unreadable / malformed skill dirs rather than aborting the scan.
			continue
		}
		skills = append(skills, skill)
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, nil
}

// loadOneDir reads SKILL.md + optional sidecar from a single skill directory
// and returns the populated Skill. Disabled is always false at load time.
func loadOneDir(dir string) (Skill, error) {
	raw, err := os.ReadFile(filepath.Join(dir, SkillFileName))
	if err != nil {
		return Skill{}, err
	}
	skill, err := Parse(raw)
	if err != nil {
		return Skill{}, err
	}
	// Default to user source; the sidecar can override it.
	skill.Source = SourceUser
	sidecar, sErr := ReadSidecar(dir)
	if sErr == nil {
		if sidecar.Source != "" {
			skill.Source = sidecar.Source
		}
		if !sidecar.UpdatedAt.IsZero() {
			skill.UpdatedAt = sidecar.UpdatedAt
		}
	} else {
		// No sidecar — fall back to the SKILL.md file mtime.
		if info, statErr := os.Stat(filepath.Join(dir, SkillFileName)); statErr == nil {
			skill.UpdatedAt = info.ModTime()
		}
	}
	return skill, nil
}

// ReadSidecar reads the optional .lemontea.json sidecar from dir.
// Returns an error if the file is missing or malformed so callers can decide
// whether to fall back to defaults.
func ReadSidecar(dir string) (SidecarManifest, error) {
	raw, err := os.ReadFile(filepath.Join(dir, SidecarFileName))
	if err != nil {
		return SidecarManifest{}, err
	}
	var sc SidecarManifest
	if err := json.Unmarshal(raw, &sc); err != nil {
		return SidecarManifest{}, err
	}
	return sc, nil
}

// WriteSidecar persists a sidecar as pretty-printed JSON inside dir.
// UpdatedAt is auto-filled with time.Now() when zero so callers can omit it.
func WriteSidecar(dir string, sc SidecarManifest) error {
	if sc.UpdatedAt.IsZero() {
		sc.UpdatedAt = time.Now()
	}
	raw, err := json.MarshalIndent(sc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, SidecarFileName), raw, 0o644)
}

// WriteSkill writes both SKILL.md and the sidecar inside `{root}/{skill.Name}/`.
// The directory is created if it does not already exist. An existing sidecar's
// CreatedAt is preserved so repeated writes don't reset the creation timestamp.
func WriteSkill(root string, s Skill) error {
	if err := Validate(s); err != nil {
		return err
	}
	dir := filepath.Join(root, s.Name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	rendered, err := Render(s)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, SkillFileName), rendered, 0o644); err != nil {
		return err
	}
	source := s.Source
	if source == "" {
		source = SourceUser
	}
	sc := SidecarManifest{Source: source, UpdatedAt: time.Now()}
	// Preserve the original CreatedAt if a sidecar already exists.
	if existing, err := ReadSidecar(dir); err == nil {
		sc.CreatedAt = existing.CreatedAt
	}
	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = time.Now()
	}
	return WriteSidecar(dir, sc)
}

// DeleteSkill removes the entire skill directory `{root}/{name}/`.
// Returns nil if the directory does not exist.
func DeleteSkill(root, name string) error {
	return os.RemoveAll(filepath.Join(root, name))
}
