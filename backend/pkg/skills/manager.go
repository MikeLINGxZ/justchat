package skills

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// Manager owns the merged in-memory view of skills (builtin + on-disk + disabled flags)
// and exposes mutation methods that persist changes to disk.
type Manager struct {
	mu      sync.RWMutex
	rootDir string
	skills  map[string]Skill // keyed by name
}

// NewManager creates a Manager that reads user skills from rootDir.
// Call Refresh to populate the in-memory cache.
func NewManager(rootDir string) *Manager {
	return &Manager{
		rootDir: rootDir,
		skills:  make(map[string]Skill),
	}
}

// Refresh reloads the in-memory state from builtin embed + on-disk dir, applying disabled names.
// Disk skills shadow builtins with the same name.
func (m *Manager) Refresh(disabledNames []string) error {
	disabled := make(map[string]bool, len(disabledNames))
	for _, n := range disabledNames {
		disabled[n] = true
	}

	builtin, err := LoadBuiltin()
	if err != nil {
		return err
	}
	disk, err := LoadFromDir(m.rootDir)
	if err != nil {
		return err
	}

	// Merge: builtins first, then disk overwrites (disk shadows builtin).
	merged := make(map[string]Skill, len(builtin)+len(disk))
	for _, s := range builtin {
		s.Disabled = disabled[s.Name]
		merged[s.Name] = s
	}
	for _, s := range disk {
		s.Disabled = disabled[s.Name]
		merged[s.Name] = s
	}

	m.mu.Lock()
	m.skills = merged
	m.mu.Unlock()
	return nil
}

// List returns all skills sorted by name.
func (m *Manager) List() []Skill {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Skill, 0, len(m.skills))
	for _, s := range m.skills {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Enabled returns only non-disabled skills.
func (m *Manager) Enabled() []Skill {
	all := m.List()
	out := make([]Skill, 0, len(all))
	for _, s := range all {
		if !s.Disabled {
			out = append(out, s)
		}
	}
	return out
}

// Get returns one skill by name.
func (m *Manager) Get(name string) (Skill, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.skills[name]
	return s, ok
}

// Create persists a new on-disk skill and refreshes the cache.
// Cannot create a skill that shadows a builtin.
func (m *Manager) Create(s Skill) (Skill, error) {
	if err := Validate(s); err != nil {
		return Skill{}, err
	}
	if s.Source == SourceBuiltin {
		return Skill{}, errors.New("cannot create builtin skill")
	}
	m.mu.RLock()
	if existing, ok := m.skills[s.Name]; ok && existing.Source == SourceBuiltin {
		m.mu.RUnlock()
		return Skill{}, errors.New("builtin skill with this name already exists")
	}
	m.mu.RUnlock()

	if s.Source == "" {
		s.Source = SourceUser
	}
	if err := WriteSkill(m.rootDir, s); err != nil {
		return Skill{}, err
	}
	// Reload from disk to get the sidecar-populated version.
	_ = m.Refresh(nil)
	got, ok := m.Get(s.Name)
	if !ok {
		return Skill{}, errors.New("skill not found after create")
	}
	return got, nil
}

// Update overwrites an existing on-disk skill. Cannot update builtin skills.
func (m *Manager) Update(name string, s Skill) (Skill, error) {
	m.mu.RLock()
	existing, ok := m.skills[name]
	if !ok {
		m.mu.RUnlock()
		return Skill{}, errors.New("skill not found")
	}
	targetName := s.Name
	if targetName == "" {
		targetName = name
	}
	if targetName != name {
		if _, exists := m.skills[targetName]; exists {
			m.mu.RUnlock()
			return Skill{}, ErrSkillNameTaken
		}
	}
	m.mu.RUnlock()
	if existing.Source == SourceBuiltin {
		return Skill{}, errors.New("builtin skill is read-only")
	}
	s.Name = targetName
	s.Source = existing.Source
	if err := WriteSkill(m.rootDir, s); err != nil {
		return Skill{}, err
	}
	if targetName != name {
		if err := DeleteSkill(m.rootDir, name); err != nil {
			return Skill{}, err
		}
	}
	_ = m.Refresh(nil)
	got, ok := m.Get(targetName)
	if !ok {
		return Skill{}, errors.New("skill not found after update")
	}
	return got, nil
}

// Delete removes a non-builtin skill from disk and cache.
func (m *Manager) Delete(name string) error {
	m.mu.RLock()
	existing, ok := m.skills[name]
	m.mu.RUnlock()
	if !ok {
		return errors.New("skill not found")
	}
	if existing.Source == SourceBuiltin {
		return errors.New("builtin skill cannot be deleted")
	}
	if err := DeleteSkill(m.rootDir, name); err != nil {
		return err
	}
	_ = m.Refresh(nil)
	return nil
}

// SetDisabled updates the in-memory disabled flag for a skill.
// This does not persist — the caller must persist DisabledNames separately.
func (m *Manager) SetDisabled(name string, disabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.skills[name]
	if !ok {
		return errors.New("skill not found")
	}
	s.Disabled = disabled
	s.UpdatedAt = time.Now()
	m.skills[name] = s
	return nil
}

// RootDir returns the on-disk skills root directory.
func (m *Manager) RootDir() string {
	return m.rootDir
}
