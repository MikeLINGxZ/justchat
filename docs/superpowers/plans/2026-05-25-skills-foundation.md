# Skills Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Skills system foundation (spec `docs/dev/13.plugin_cli.md` §3.4) — a Claude Code-compatible Skill format, file-based persistence, a dedicated Skills settings page, and chat-agent integration via a `Skill` meta-tool and `/skill-name` slash command. This is Plan A of three; later plans (Hidden Sessions, CLI Refactor) depend on it.

**Architecture:**
- Backend `backend/pkg/skills/` (core: types, parser, loader, manager, builtin embedding) + `backend/service/skills/` (Wails methods, DTOs).
- Skills live as directories under `{data_dir}/skills/<name>/`: a required `SKILL.md` (markdown + YAML frontmatter, Claude Code shape) plus an optional `manifest.json` sidecar that tracks origin (`user` / `ai`) and timestamps. Built-in skills are embedded into the Go binary via `//go:embed` and surface alongside on-disk skills (on-disk shadows builtin with the same name).
- Disable state is one new slice on `data_models.Config`: `DisabledSkills []string`. Everything else (content, source, sidecar metadata) lives next to the skill on disk.
- Frontend gets a new `skills` tab in `SettingsApp` mirroring the plugins page layout (list left, detail right) plus an editor and an import dialog.
- Agent layer registers a built-in `Skill` tool whose description is rebuilt every send with the live list of enabled skills (name + description). When AI calls `Skill({name})`, the manager returns the full `SKILL.md` body as the tool result. Chat input gets a TipTap suggestion popup on `/` that inserts an explicit `Skill({name})` tool call (handled by the existing agent flow).

**Tech Stack:** Go 1.22+, Wails v3, React 18, TipTap (`@tiptap/suggestion`), Vite, Vitest, Tailwind, i18next, `gopkg.in/yaml.v3` (already indirectly imported — verify in Task 0).

---

## Pre-flight (do this once before starting)

- [ ] **Confirm YAML dep:** Run `grep -R "yaml.v" /Users/linhuafeng/Work/lemon_tea_desktop/go.mod` — if `gopkg.in/yaml.v3` is missing, run `go get gopkg.in/yaml.v3` and `go mod tidy`. Commit `go.mod` / `go.sum` alone before Task 1.
- [ ] **Pin the date:** Today is 2026-05-25 — this plan filename is `docs/superpowers/plans/2026-05-25-skills-foundation.md`. Confirm before starting.

---

## Section 1 — Backend foundation

### Task 1: Data directory + ierror codes

**Files:**
- Modify: `backend/pkg/dir/dir.go` (add constant + helper at end of file)
- Modify: `backend/pkg/ierror/error.go` (add new block at end of `const ( ... )`)

- [ ] **Step 1: Write test for `dir.SkillsRoot`**

Create `backend/pkg/dir/dir_test.go` (if missing) and add:

```go
package dir

import (
	"path/filepath"
	"testing"
)

func TestSkillsRoot(t *testing.T) {
	got := SkillsRoot("/tmp/data")
	want := filepath.Join("/tmp/data", "skills")
	if got != want {
		t.Fatalf("SkillsRoot mismatch: got %q want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test — must fail with "undefined: SkillsRoot"**

```
go test ./backend/pkg/dir/...
```

- [ ] **Step 3: Add helpers to `backend/pkg/dir/dir.go`**

Append at the end of file:

```go
// SkillsSubdirName is the directory name for skills under data_dir.
const SkillsSubdirName = "skills"

// SkillsRoot returns the skills root: {dataDir}/skills.
func SkillsRoot(dataDir string) string {
	return filepath.Join(dataDir, SkillsSubdirName)
}
```

- [ ] **Step 4: Run test — must pass**

```
go test ./backend/pkg/dir/...
```

- [ ] **Step 5: Append ierror codes**

In `backend/pkg/ierror/error.go`, inside the existing `const ( ... )` block (after CLI codes), add:

```go
	// Skills error codes
	ErrSkillsLoadFailed     errorCode = "ierror.skills.load_failed"
	ErrSkillsNotFound       errorCode = "ierror.skills.not_found"
	ErrSkillsInvalidName    errorCode = "ierror.skills.invalid_name"
	ErrSkillsInvalidContent errorCode = "ierror.skills.invalid_content"
	ErrSkillsNameTaken      errorCode = "ierror.skills.name_taken"
	ErrSkillsBuiltinLocked  errorCode = "ierror.skills.builtin_locked"
	ErrSkillsWriteFailed    errorCode = "ierror.skills.write_failed"
	ErrSkillsDeleteFailed   errorCode = "ierror.skills.delete_failed"
```

- [ ] **Step 6: Add i18n strings (zh-CN + en) for new error codes**

In `frontend/src/i18n/locales/zh-CN.ts` and `en.ts`, ensure top-level `ierror.skills.*` keys exist. Pattern matches existing `ierror.cli.*`:

```ts
'ierror.skills.load_failed': '加载 skill 失败',     // en: 'Failed to load skill'
'ierror.skills.not_found': '找不到 skill',           // en: 'Skill not found'
'ierror.skills.invalid_name': 'skill 名称不合法',     // en: 'Invalid skill name'
'ierror.skills.invalid_content': 'skill 内容不合法', // en: 'Invalid skill content'
'ierror.skills.name_taken': 'skill 名称已被占用',     // en: 'Skill name already in use'
'ierror.skills.builtin_locked': '内置 skill 不可修改', // en: 'Built-in skill is read-only'
'ierror.skills.write_failed': '保存 skill 失败',     // en: 'Failed to save skill'
'ierror.skills.delete_failed': '删除 skill 失败',     // en: 'Failed to delete skill'
```

(Match the surrounding flat-key vs nested-object style of the file — read the file first.)

- [ ] **Step 7: Build to confirm Go side compiles**

```
go build ./...
```

Expected: no errors.

- [ ] **Step 8: Commit**

```bash
git add backend/pkg/dir/ backend/pkg/ierror/error.go frontend/src/i18n/
git commit -m "feat(skills): add data dir helper and ierror codes"
```

---

### Task 2: Skill domain types + frontmatter parser

**Files:**
- Create: `backend/pkg/skills/types.go`
- Create: `backend/pkg/skills/parse.go`
- Create: `backend/pkg/skills/parse_test.go`

- [ ] **Step 1: Write `parse_test.go` with failing tests**

```go
package skills

import (
	"strings"
	"testing"
)

func TestParseSkill_Valid(t *testing.T) {
	raw := strings.Join([]string{
		"---",
		"name: install-cli-from-docs",
		"description: Install a CLI plugin from official documentation",
		"---",
		"",
		"Body line 1",
		"Body line 2",
	}, "\n")

	got, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "install-cli-from-docs" {
		t.Errorf("name: got %q", got.Name)
	}
	if got.Description == "" {
		t.Errorf("description must not be empty")
	}
	if !strings.Contains(got.Body, "Body line 1") {
		t.Errorf("body lost: %q", got.Body)
	}
}

func TestParseSkill_MissingFrontmatter(t *testing.T) {
	_, err := Parse([]byte("no frontmatter here"))
	if err == nil {
		t.Fatal("expected error for missing frontmatter")
	}
}

func TestParseSkill_MissingName(t *testing.T) {
	raw := "---\ndescription: x\n---\nbody"
	_, err := Parse([]byte(raw))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParseSkill_InvalidNameChars(t *testing.T) {
	raw := "---\nname: Bad Name!\ndescription: x\n---\nbody"
	_, err := Parse([]byte(raw))
	if err == nil {
		t.Fatal("expected error for invalid name chars")
	}
}
```

- [ ] **Step 2: Run — must fail (Parse undefined)**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 3: Create `types.go`**

```go
// Package skills holds the on-disk Skill model and IO for the Skills system.
package skills

import "time"

// Source identifies where a skill originated.
type Source string

const (
	SourceBuiltin Source = "builtin"
	SourceUser    Source = "user"
	SourceAI      Source = "ai"
)

// Skill is the in-memory representation of one skill.
// Body is the markdown content beneath the frontmatter.
type Skill struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Source      Source    `json:"source"`
	Disabled    bool      `json:"disabled"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SidecarManifest is the .lemontea sidecar stored next to a SKILL.md on disk.
// Built-in skills do not have a sidecar.
type SidecarManifest struct {
	Source    Source    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- [ ] **Step 4: Create `parse.go`**

```go
package skills

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var skillNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

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
	body := trimmed[len(delim):]
	// Skip optional newline after first delim.
	body = bytes.TrimLeft(body, "\r\n")
	end := bytes.Index(body, []byte("\n"+delim))
	if end < 0 {
		return Skill{}, errors.New("frontmatter not closed")
	}
	headerBytes := body[:end]
	rest := body[end+len("\n"+delim):]
	rest = bytes.TrimLeft(rest, "\r\n")

	var fm Frontmatter
	if err := yaml.Unmarshal(headerBytes, &fm); err != nil {
		return Skill{}, err
	}
	fm.Name = strings.TrimSpace(fm.Name)
	fm.Description = strings.TrimSpace(fm.Description)
	if fm.Name == "" {
		return Skill{}, errors.New("frontmatter.name is required")
	}
	if fm.Description == "" {
		return Skill{}, errors.New("frontmatter.description is required")
	}
	if !skillNameRe.MatchString(fm.Name) {
		return Skill{}, errors.New("frontmatter.name must be kebab-case (a-z0-9- only, max 64 chars)")
	}
	return Skill{
		Name:        fm.Name,
		Description: fm.Description,
		Body:        string(rest),
	}, nil
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
```

- [ ] **Step 5: Run tests — must pass**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 6: Commit**

```bash
git add backend/pkg/skills/
git commit -m "feat(skills): add domain types and frontmatter parser"
```

---

### Task 3: On-disk loader + sidecar IO

**Files:**
- Create: `backend/pkg/skills/loader.go`
- Create: `backend/pkg/skills/loader_test.go`

- [ ] **Step 1: Write failing test**

```go
package skills

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeSkill(t *testing.T, root, name, description, body string, sidecar *SidecarManifest) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	skill := Skill{Name: name, Description: description, Body: body}
	out, err := Render(skill)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), out, 0o644); err != nil {
		t.Fatal(err)
	}
	if sidecar != nil {
		if err := WriteSidecar(dir, *sidecar); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLoadFromDir_UserDefault(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "foo", "Foo skill", "Hello", nil)

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Source != SourceUser {
		t.Errorf("expected SourceUser, got %q", skills[0].Source)
	}
}

func TestLoadFromDir_AISourceFromSidecar(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "ai-thing", "AI skill", "x", &SidecarManifest{
		Source:    SourceAI,
		CreatedAt: time.Now(),
	})

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if skills[0].Source != SourceAI {
		t.Errorf("expected SourceAI, got %q", skills[0].Source)
	}
}

func TestLoadFromDir_SkipsInvalid(t *testing.T) {
	root := t.TempDir()
	// Valid skill
	writeSkill(t, root, "good", "Good", "x", nil)
	// Invalid: no SKILL.md
	if err := os.MkdirAll(filepath.Join(root, "broken"), 0o755); err != nil {
		t.Fatal(err)
	}

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "good" {
		t.Fatalf("expected only 'good', got %+v", skills)
	}
}
```

- [ ] **Step 2: Run — must fail (LoadFromDir undefined)**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 3: Create `loader.go`**

```go
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
	SkillFileName    = "SKILL.md"
	SidecarFileName  = ".lemontea.json"
)

// LoadFromDir scans root for skill subdirectories and returns them sorted by name.
// Subdirs without a parseable SKILL.md are silently skipped.
func LoadFromDir(root string) ([]Skill, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
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
			continue
		}
		skills = append(skills, skill)
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, nil
}

// loadOneDir reads SKILL.md + optional sidecar into a Skill (Disabled is false here).
func loadOneDir(dir string) (Skill, error) {
	raw, err := os.ReadFile(filepath.Join(dir, SkillFileName))
	if err != nil {
		return Skill{}, err
	}
	skill, err := Parse(raw)
	if err != nil {
		return Skill{}, err
	}
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
		if info, statErr := os.Stat(filepath.Join(dir, SkillFileName)); statErr == nil {
			skill.UpdatedAt = info.ModTime()
		}
	}
	return skill, nil
}

// ReadSidecar reads the optional .lemontea.json sidecar.
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

// WriteSidecar persists a sidecar atomically.
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

// WriteSkill writes both SKILL.md and sidecar inside `{root}/{skill.Name}/`. Creates the dir if needed.
func WriteSkill(root string, s Skill) error {
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
	if existing, err := ReadSidecar(dir); err == nil {
		sc.CreatedAt = existing.CreatedAt
	}
	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = time.Now()
	}
	return WriteSidecar(dir, sc)
}

// DeleteSkill removes the directory `{root}/{name}/`.
func DeleteSkill(root, name string) error {
	return os.RemoveAll(filepath.Join(root, name))
}
```

- [ ] **Step 4: Run tests — must pass**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/skills/loader.go backend/pkg/skills/loader_test.go
git commit -m "feat(skills): add filesystem loader and sidecar IO"
```

---

### Task 4: Built-in skill embedding

**Files:**
- Create: `backend/pkg/skills/builtin.go`
- Create: `backend/pkg/skills/_builtin/.gitkeep` (empty file; populated in Plan C)
- Create: `backend/pkg/skills/builtin_test.go`

- [ ] **Step 1: Create empty embed dir**

```bash
mkdir -p backend/pkg/skills/_builtin
touch backend/pkg/skills/_builtin/.gitkeep
```

- [ ] **Step 2: Write failing test for `LoadBuiltin`**

```go
package skills

import "testing"

func TestLoadBuiltin_ReturnsEmptySliceWhenNoSkills(t *testing.T) {
	got, err := LoadBuiltin()
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		got = []Skill{}
	}
	if len(got) != 0 {
		t.Fatalf("expected zero builtin skills in Plan A, got %d", len(got))
	}
}
```

- [ ] **Step 3: Run — must fail (LoadBuiltin undefined)**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 4: Create `builtin.go`**

```go
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
```

- [ ] **Step 5: Run test — must pass**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 6: Commit**

```bash
git add backend/pkg/skills/builtin.go backend/pkg/skills/builtin_test.go backend/pkg/skills/_builtin/
git commit -m "feat(skills): add builtin skill embedding"
```

---

### Task 5: Manager (merge builtin + disk + disabled state)

**Files:**
- Create: `backend/pkg/skills/manager.go`
- Create: `backend/pkg/skills/manager_test.go`

- [ ] **Step 1: Write failing test**

```go
package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_DiskShadowsBuiltinIsNoOpWhenBuiltinAbsent(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "user-a", "User A", "body-a", nil)
	writeSkill(t, root, "user-b", "User B", "body-b", nil)

	m := NewManager(root)
	if err := m.Refresh(nil); err != nil {
		t.Fatal(err)
	}
	skills := m.List()
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestManager_Disabled(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "alpha", "Alpha", "x", nil)

	m := NewManager(root)
	if err := m.Refresh([]string{"alpha"}); err != nil {
		t.Fatal(err)
	}
	got, ok := m.Get("alpha")
	if !ok {
		t.Fatal("missing skill alpha")
	}
	if !got.Disabled {
		t.Fatal("expected alpha to be disabled")
	}
	if len(m.Enabled()) != 0 {
		t.Fatalf("expected zero enabled, got %d", len(m.Enabled()))
	}
}

func TestManager_Create_Update_Delete(t *testing.T) {
	root := t.TempDir()
	m := NewManager(root)
	if err := m.Refresh(nil); err != nil {
		t.Fatal(err)
	}

	created, err := m.Create(Skill{
		Name: "x-one", Description: "desc", Body: "body", Source: SourceUser,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "x-one" {
		t.Fatalf("unexpected name: %q", created.Name)
	}
	if _, err := os.Stat(filepath.Join(root, "x-one", SkillFileName)); err != nil {
		t.Fatalf("SKILL.md missing: %v", err)
	}

	if _, err := m.Update("x-one", Skill{Name: "x-one", Description: "new", Body: "new-body"}); err != nil {
		t.Fatal(err)
	}
	got, _ := m.Get("x-one")
	if got.Description != "new" {
		t.Fatalf("update did not persist description: %q", got.Description)
	}

	if err := m.Delete("x-one"); err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Get("x-one"); ok {
		t.Fatal("delete did not remove from cache")
	}
}
```

- [ ] **Step 2: Run — must fail**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 3: Create `manager.go`**

```go
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
	mu       sync.RWMutex
	rootDir  string
	skills   map[string]Skill // keyed by name
}

func NewManager(rootDir string) *Manager {
	return &Manager{
		rootDir: rootDir,
		skills:  make(map[string]Skill),
	}
}

// Refresh reloads the in-memory state from builtin embed + on-disk dir, applying disabled names.
func (m *Manager) Refresh(disabled []string) error {
	disabledSet := make(map[string]bool, len(disabled))
	for _, name := range disabled {
		disabledSet[name] = true
	}

	merged := make(map[string]Skill)

	builtin, err := LoadBuiltin()
	if err != nil {
		return err
	}
	for _, s := range builtin {
		s.Disabled = disabledSet[s.Name]
		merged[s.Name] = s
	}

	onDisk, err := LoadFromDir(m.rootDir)
	if err != nil {
		return err
	}
	for _, s := range onDisk {
		s.Disabled = disabledSet[s.Name]
		merged[s.Name] = s // disk shadows builtin
	}

	m.mu.Lock()
	m.skills = merged
	m.mu.Unlock()
	return nil
}

// List returns all skills (including disabled) sorted by name.
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

// Enabled returns only enabled skills.
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

// Create persists a new on-disk skill. Fails if name already exists on disk
// (builtin skills with the same name are shadowed, not blocked).
func (m *Manager) Create(s Skill) (Skill, error) {
	if _, err := Parse(mustRender(s)); err != nil {
		return Skill{}, err
	}
	m.mu.RLock()
	existing, ok := m.skills[s.Name]
	m.mu.RUnlock()
	if ok && existing.Source != SourceBuiltin {
		return Skill{}, errors.New("skill name already exists")
	}
	if s.Source == "" {
		s.Source = SourceUser
	}
	if err := WriteSkill(m.rootDir, s); err != nil {
		return Skill{}, err
	}
	s.UpdatedAt = time.Now()
	m.mu.Lock()
	m.skills[s.Name] = s
	m.mu.Unlock()
	return s, nil
}

// Update overwrites an existing on-disk skill. Cannot update builtin skills.
func (m *Manager) Update(name string, s Skill) (Skill, error) {
	m.mu.RLock()
	existing, ok := m.skills[name]
	m.mu.RUnlock()
	if !ok {
		return Skill{}, errors.New("skill not found")
	}
	if existing.Source == SourceBuiltin {
		return Skill{}, errors.New("builtin skill is read-only")
	}
	if s.Name == "" {
		s.Name = name
	}
	if s.Name != name {
		// Rename: write new, delete old.
		if err := WriteSkill(m.rootDir, s); err != nil {
			return Skill{}, err
		}
		if err := DeleteSkill(m.rootDir, name); err != nil {
			return Skill{}, err
		}
		m.mu.Lock()
		delete(m.skills, name)
		s.Source = existing.Source
		s.UpdatedAt = time.Now()
		m.skills[s.Name] = s
		m.mu.Unlock()
		return s, nil
	}
	s.Source = existing.Source
	if err := WriteSkill(m.rootDir, s); err != nil {
		return Skill{}, err
	}
	s.UpdatedAt = time.Now()
	m.mu.Lock()
	m.skills[s.Name] = s
	m.mu.Unlock()
	return s, nil
}

// Delete removes an on-disk skill. Builtin skills cannot be deleted.
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
	m.mu.Lock()
	delete(m.skills, name)
	m.mu.Unlock()
	return nil
}

func mustRender(s Skill) []byte {
	b, err := Render(s)
	if err != nil {
		return nil
	}
	return b
}
```

- [ ] **Step 4: Run tests — must pass**

```
go test ./backend/pkg/skills/...
```

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/skills/manager.go backend/pkg/skills/manager_test.go
git commit -m "feat(skills): add manager merging builtin and on-disk skills"
```

---

### Task 6: Persist `DisabledSkills` on Config

**Files:**
- Modify: `backend/models/data_models/config.go` (add field at end of `Config`)

- [ ] **Step 1: Add field to Config**

In `backend/models/data_models/config.go`, append to the `Config` struct:

```go
	DisabledSkills    []string        `json:"disabled_skills"`
```

So `Config` ends with `Extensions` then `DisabledSkills`.

- [ ] **Step 2: Build to confirm**

```
go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/models/data_models/config.go
git commit -m "feat(skills): persist disabled skill names on Config"
```

---

### Task 7: Wails service `backend/service/skills/`

**Files:**
- Create: `backend/service/skills/skills.go`
- Create: `backend/service/skills/skills_implement.go`
- Create: `backend/service/skills/skills_internal.go`
- Create: `backend/service/skills/skills_dto/list_skills.go`
- Create: `backend/service/skills/skills_dto/get_skill.go`
- Create: `backend/service/skills/skills_dto/create_skill.go`
- Create: `backend/service/skills/skills_dto/update_skill.go`
- Create: `backend/service/skills/skills_dto/delete_skill.go`
- Create: `backend/service/skills/skills_dto/toggle_skill.go`
- Create: `backend/service/skills/skills_dto/import_skill.go`
- Create: `backend/service/skills/skills_test.go`

- [ ] **Step 1: Create DTO files**

Each DTO file matches the existing pattern. `list_skills.go`:

```go
package skills_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"

type ListSkillsInput struct{}

type SkillItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Disabled    bool   `json:"disabled"`
	HasBody     bool   `json:"has_body"`
}

type ListSkillsOutput struct {
	Skills []SkillItem `json:"skills"`
}

func ToItem(s skills.Skill) SkillItem {
	return SkillItem{
		Name:        s.Name,
		Description: s.Description,
		Source:      string(s.Source),
		Disabled:    s.Disabled,
		HasBody:     s.Body != "",
	}
}
```

`get_skill.go`:

```go
package skills_dto

type GetSkillInput struct {
	Name string `json:"name"`
}

type GetSkillOutput struct {
	Skill SkillFull `json:"skill"`
}

type SkillFull struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Source      string `json:"source"`
	Disabled    bool   `json:"disabled"`
}
```

`create_skill.go`:

```go
package skills_dto

type CreateSkillInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Source      string `json:"source"` // "user" | "ai"; empty defaults to "user"
}

type CreateSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
```

`update_skill.go`:

```go
package skills_dto

type UpdateSkillInput struct {
	OriginalName string `json:"original_name"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Body         string `json:"body"`
}

type UpdateSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
```

`delete_skill.go`:

```go
package skills_dto

type DeleteSkillInput struct {
	Name string `json:"name"`
}

type DeleteSkillOutput struct{}
```

`toggle_skill.go`:

```go
package skills_dto

type ToggleSkillInput struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type ToggleSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
```

`import_skill.go`:

```go
package skills_dto

// ImportSkillInput accepts raw SKILL.md content. Optionally name override (otherwise pulled from frontmatter).
type ImportSkillInput struct {
	Content string `json:"content"`
}

type ImportSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
```

- [ ] **Step 2: Create `skills_implement.go`**

```go
package skills

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup is invoked by Wails after the app is constructed.
// It captures the application handle for later event emission.
func (s *Skills) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.wailsApp = application.Get()
	return s.refresh()
}
```

- [ ] **Step 3: Create `skills_internal.go`**

```go
package skills

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// refresh reloads the in-memory skills manager from disk + disabled state in config.
func (s *Skills) refresh() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	return s.manager.Refresh(cfg.DisabledSkills)
}

// loadConfig reads config.json (returns zero-value on missing file).
func loadConfig() (data_models.Config, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return data_models.Config{}, err
	}
	path := filepath.Join(dataDir, dir.ConfigFileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return data_models.Config{}, nil
		}
		return data_models.Config{}, err
	}
	var cfg data_models.Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return data_models.Config{}, err
	}
	return cfg, nil
}

// saveConfig writes config.json with merged updates.
func saveConfig(cfg data_models.Config) error {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dataDir, dir.ConfigFileName)
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// setDisabled updates DisabledSkills in config.json (idempotent).
func setDisabled(name string, disabled bool) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	existing := cfg.DisabledSkills
	cfg.DisabledSkills = cfg.DisabledSkills[:0]
	for _, n := range existing {
		if n == name {
			continue
		}
		cfg.DisabledSkills = append(cfg.DisabledSkills, n)
	}
	if disabled {
		cfg.DisabledSkills = append(cfg.DisabledSkills, name)
	}
	return saveConfig(cfg)
}
```

- [ ] **Step 4: Create `skills.go` (Wails methods)**

```go
package skills

import (
	"context"
	"errors"

	"github.com/wailsapp/wails/v3/pkg/application"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto"
)

// Skills is the Wails service for managing skills.
type Skills struct {
	manager  *pkgskills.Manager
	wailsApp *application.App
}

// NewSkills wires the manager to {data_dir}/skills.
func NewSkills() *Skills {
	dataDir, _ := dir.GetDataDir()
	return &Skills{
		manager: pkgskills.NewManager(dir.SkillsRoot(dataDir)),
	}
}

// ListSkills returns all known skills (built-in, user, ai).
func (s *Skills) ListSkills(ctx context.Context, _ skills_dto.ListSkillsInput) (*skills_dto.ListSkillsOutput, error) {
	if err := s.refresh(); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsLoadFailed, err)
	}
	out := skills_dto.ListSkillsOutput{}
	for _, sk := range s.manager.List() {
		out.Skills = append(out.Skills, skills_dto.ToItem(sk))
	}
	return &out, nil
}

// GetSkill returns one skill including its full body.
func (s *Skills) GetSkill(ctx context.Context, input skills_dto.GetSkillInput) (*skills_dto.GetSkillOutput, error) {
	sk, ok := s.manager.Get(input.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, errors.New(input.Name))
	}
	return &skills_dto.GetSkillOutput{
		Skill: skills_dto.SkillFull{
			Name:        sk.Name,
			Description: sk.Description,
			Body:        sk.Body,
			Source:      string(sk.Source),
			Disabled:    sk.Disabled,
		},
	}, nil
}

// CreateSkill persists a new user / ai skill.
func (s *Skills) CreateSkill(ctx context.Context, input skills_dto.CreateSkillInput) (*skills_dto.CreateSkillOutput, error) {
	source := pkgskills.SourceUser
	if input.Source == "ai" {
		source = pkgskills.SourceAI
	}
	skill := pkgskills.Skill{
		Name: input.Name, Description: input.Description, Body: input.Body, Source: source,
	}
	created, err := s.manager.Create(skill)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	return &skills_dto.CreateSkillOutput{Skill: skills_dto.ToItem(created)}, nil
}

// UpdateSkill replaces an existing user / ai skill (builtin is read-only).
func (s *Skills) UpdateSkill(ctx context.Context, input skills_dto.UpdateSkillInput) (*skills_dto.UpdateSkillOutput, error) {
	skill := pkgskills.Skill{Name: input.Name, Description: input.Description, Body: input.Body}
	updated, err := s.manager.Update(input.OriginalName, skill)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	return &skills_dto.UpdateSkillOutput{Skill: skills_dto.ToItem(updated)}, nil
}

// DeleteSkill removes a non-builtin skill.
func (s *Skills) DeleteSkill(ctx context.Context, input skills_dto.DeleteSkillInput) (*skills_dto.DeleteSkillOutput, error) {
	if err := s.manager.Delete(input.Name); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsDeleteFailed, err)
	}
	return &skills_dto.DeleteSkillOutput{}, nil
}

// ToggleSkill enables / disables a skill by updating config.DisabledSkills, then refreshing.
func (s *Skills) ToggleSkill(ctx context.Context, input skills_dto.ToggleSkillInput) (*skills_dto.ToggleSkillOutput, error) {
	if _, ok := s.manager.Get(input.Name); !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, errors.New(input.Name))
	}
	if err := setDisabled(input.Name, !input.Enabled); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	if err := s.refresh(); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsLoadFailed, err)
	}
	sk, _ := s.manager.Get(input.Name)
	return &skills_dto.ToggleSkillOutput{Skill: skills_dto.ToItem(sk)}, nil
}

// ImportSkill parses raw SKILL.md content and creates a user skill.
func (s *Skills) ImportSkill(ctx context.Context, input skills_dto.ImportSkillInput) (*skills_dto.ImportSkillOutput, error) {
	parsed, err := pkgskills.Parse([]byte(input.Content))
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsInvalidContent, err)
	}
	parsed.Source = pkgskills.SourceUser
	created, err := s.manager.Create(parsed)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	return &skills_dto.ImportSkillOutput{Skill: skills_dto.ToItem(created)}, nil
}
```

- [ ] **Step 5: Write a service-level test**

`skills_test.go`:

```go
package skills

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto"
)

func newTestService(t *testing.T) (*Skills, string) {
	t.Helper()
	tempDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDir)
	skillsRoot := filepath.Join(tempDir, "skills")
	if err := os.MkdirAll(skillsRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	svc := &Skills{manager: pkgskills.NewManager(skillsRoot)}
	if err := svc.refresh(); err != nil {
		t.Fatal(err)
	}
	return svc, tempDir
}

func TestSkillsService_CreateListGet(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	if _, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name: "demo", Description: "d", Body: "hi",
	}); err != nil {
		t.Fatal(err)
	}
	list, err := svc.ListSkills(ctx, skills_dto.ListSkillsInput{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Skills) != 1 || list.Skills[0].Name != "demo" {
		t.Fatalf("list: %+v", list)
	}
	got, err := svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "demo"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Skill.Body != "hi\n" && got.Skill.Body != "hi" {
		t.Fatalf("body: %q", got.Skill.Body)
	}
}
```

- [ ] **Step 6: Run tests — must pass**

```
go test ./backend/service/skills/...
```

- [ ] **Step 7: Commit**

```bash
git add backend/service/skills/
git commit -m "feat(skills): add Wails service for skill CRUD"
```

---

### Task 8: Register service in main.go

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add import**

Inside the import block of `main.go`, add (alphabetically near other `service/...`):

```go
	skillsSvc "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills"
```

- [ ] **Step 2: Register in Services list**

Insert in the `Services: []application.Service{ ... }` slice, after the plugin service line:

```go
		application.NewService(skillsSvc.NewSkills()),
```

- [ ] **Step 3: Build to confirm**

```
go build ./...
```

- [ ] **Step 4: Regenerate Wails bindings**

```
wails generate bindings
```

Confirm `frontend/bindings/.../service/skills/index.ts` now exists.

- [ ] **Step 5: Commit**

```bash
git add main.go frontend/bindings/
git commit -m "feat(skills): register skills service and generate bindings"
```

---

## Section 2 — Frontend Skills settings page

### Task 9: i18n strings + types

**Files:**
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`
- Create: `frontend/src/types/skills.ts`

- [ ] **Step 1: Add `settingsPage.skills.*` strings (zh-CN)**

```ts
'settingsPage.primary.skills': 'Skills',
'settingsPage.skills.title': 'Skills',
'settingsPage.skills.empty': '还没有 Skills',
'settingsPage.skills.addNew': '新建 Skill',
'settingsPage.skills.import': '导入 Skill',
'settingsPage.skills.sourceBuiltin': '内置',
'settingsPage.skills.sourceUser': '用户',
'settingsPage.skills.sourceAI': 'AI 生成',
'settingsPage.skills.enabled': '启用',
'settingsPage.skills.delete': '删除',
'settingsPage.skills.deleteConfirm': '确定要删除该 skill 吗？',
'settingsPage.skills.editor.name': '名称',
'settingsPage.skills.editor.description': '描述',
'settingsPage.skills.editor.body': '正文 (Markdown)',
'settingsPage.skills.editor.save': '保存',
'settingsPage.skills.editor.cancel': '取消',
'settingsPage.skills.importDialog.title': '导入 Skill',
'settingsPage.skills.importDialog.placeholder': '粘贴一份 SKILL.md（含 --- frontmatter ---）',
'settingsPage.skills.importDialog.submit': '导入',
'settingsPage.skills.builtinLockedMsg': '内置 skill 仅可启用 / 禁用，不可编辑或删除',
```

- [ ] **Step 2: Mirror to `en.ts`** with English values.

- [ ] **Step 3: Create `frontend/src/types/skills.ts`**

```ts
export type SkillSource = 'builtin' | 'user' | 'ai'

export type SkillItem = {
  name: string
  description: string
  source: SkillSource
  disabled: boolean
  has_body: boolean
}

export type SkillFull = SkillItem & {
  body: string
}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/i18n/ frontend/src/types/skills.ts
git commit -m "feat(skills): add i18n strings and frontend types"
```

---

### Task 10: Add `skills` tab to SettingsApp and primary menu

**Files:**
- Modify: `frontend/src/components/settings/SettingsPrimaryMenu.tsx`
- Modify: `frontend/src/store/settingsStore.ts`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`

- [ ] **Step 1: Add `skills` to the activeTab union**

In `frontend/src/store/settingsStore.ts`, locate the `activeTab` type and extend:

```ts
activeTab: 'general' | 'providers' | 'plugins' | 'skills' | 'about'
```

(Search the file first — exact path may live alongside `setActiveTab`.)

- [ ] **Step 2: Add menu item**

In `SettingsPrimaryMenu.tsx`, add an item alongside the other 4 (use a `Sparkles` icon from `lucide-react`):

```tsx
{ key: 'skills', label: t('settingsPage.primary.skills'), icon: <Sparkles size={16} /> },
```

Adjust imports to include `Sparkles`.

- [ ] **Step 3: Route the tab in `SettingsApp.tsx`**

Find the chain `if (activeTab === 'about')` ... `else if (activeTab === 'plugins')`. Insert a new `else if (activeTab === 'skills')` branch *before* the trailing `else` (general fallback):

```tsx
  } else if (activeTab === 'skills') {
    content = (
      <SettingsSectionLayout
        sidebarClassName="w-80"
        sidebar={<SkillsList />}
      >
        <SkillsDetailPane />
      </SettingsSectionLayout>
    )
  }
```

Import the new components at top of file:

```tsx
import { SkillsList } from '@/components/settings/skills/SkillsList'
import { SkillsDetailPane } from '@/components/settings/skills/SkillsDetailPane'
```

Also update the `URLSearchParams` allow-list line:

```tsx
if (tab === 'general' || tab === 'providers' || tab === 'plugins' || tab === 'skills' || tab === 'about') {
```

- [ ] **Step 4: Build front-end check**

```
cd frontend && npm run build
```

Expect compile errors only from missing components — those land in Tasks 11–14.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/settings/SettingsPrimaryMenu.tsx frontend/src/store/settingsStore.ts frontend/src/components/settings/SettingsApp.tsx
git commit -m "feat(skills): wire Skills tab into settings navigation"
```

---

### Task 11: `SkillsList` + store

**Files:**
- Create: `frontend/src/components/settings/skills/SkillsList.tsx`
- Create: `frontend/src/components/settings/skills/SkillsListItem.tsx`
- Modify: `frontend/src/store/settingsStore.ts` (add `skills` slice)

- [ ] **Step 1: Add skills state to store**

In `settingsStore.ts`, add fields:

```ts
skills: SkillItem[]
selectedSkillName: string | null
setSkills: (items: SkillItem[]) => void
setSelectedSkillName: (name: string | null) => void
updateSkill: (item: SkillItem) => void
removeSkill: (name: string) => void
```

Implement matching reducers (mirror existing `extensions` reducer style). Add `import type { SkillItem } from '@/types/skills'`.

- [ ] **Step 2: Create `SkillsListItem.tsx`**

```tsx
import { useTranslation } from 'react-i18next'
import type { SkillItem } from '@/types/skills'

type Props = {
  item: SkillItem
  selected: boolean
  onSelect: (name: string) => void
  onToggle: (item: SkillItem) => void
}

export function SkillsListItem({ item, selected, onSelect, onToggle }: Props) {
  const { t } = useTranslation()
  const badge = t(`settingsPage.skills.source${capitalize(item.source)}`)
  return (
    <button
      type="button"
      onClick={() => onSelect(item.name)}
      className={`flex w-full flex-col items-start gap-1 rounded-md px-3 py-2 text-left transition ${
        selected ? 'bg-accent text-accent-foreground' : 'hover:bg-muted'
      }`}
    >
      <div className="flex w-full items-center justify-between">
        <span className="text-sm font-medium">{item.name}</span>
        <span className="text-xs uppercase tracking-wide opacity-60">{badge}</span>
      </div>
      <span className="text-xs opacity-70 line-clamp-2">{item.description}</span>
      <label className="mt-1 flex items-center gap-1 text-xs">
        <input
          type="checkbox"
          checked={!item.disabled}
          onChange={(e) => {
            e.stopPropagation()
            onToggle(item)
          }}
        />
        {t('settingsPage.skills.enabled')}
      </label>
    </button>
  )
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1)
}
```

- [ ] **Step 3: Create `SkillsList.tsx`**

```tsx
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, Upload } from 'lucide-react'
import { useSettingsStore } from '@/store/settingsStore'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { ToggleSkillInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'
import { SkillsListItem } from './SkillsListItem'
import { SkillImportDialog } from './SkillImportDialog'
import type { SkillItem } from '@/types/skills'

export function SkillsList() {
  const { t } = useTranslation()
  const skills = useSettingsStore((s) => s.skills)
  const setSkills = useSettingsStore((s) => s.setSkills)
  const selectedSkillName = useSettingsStore((s) => s.selectedSkillName)
  const setSelectedSkillName = useSettingsStore((s) => s.setSelectedSkillName)
  const updateSkill = useSettingsStore((s) => s.updateSkill)
  const [importOpen, setImportOpen] = useState(false)

  const refresh = useCallback(async () => {
    const result = await SkillsBinding.ListSkills({})
    setSkills(((result?.skills ?? []) as unknown as SkillItem[]))
  }, [setSkills])

  useEffect(() => { void refresh() }, [refresh])

  const handleToggle = useCallback(async (item: SkillItem) => {
    const nextEnabled = item.disabled
    const optimistic: SkillItem = { ...item, disabled: !nextEnabled }
    updateSkill(optimistic)
    try {
      const result = await SkillsBinding.ToggleSkill(new ToggleSkillInput({ name: item.name, enabled: nextEnabled }))
      if (result?.skill) {
        updateSkill(result.skill as unknown as SkillItem)
      }
    } catch {
      updateSkill(item)
    }
  }, [updateSkill])

  const handleStartNew = () => setSelectedSkillName('__new__')

  return (
    <div className="flex h-full flex-col gap-2 p-3">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold">{t('settingsPage.skills.title')}</h3>
        <div className="flex gap-1">
          <button type="button" onClick={handleStartNew} className="rounded-md p-1 hover:bg-muted">
            <Plus size={16} />
          </button>
          <button type="button" onClick={() => setImportOpen(true)} className="rounded-md p-1 hover:bg-muted">
            <Upload size={16} />
          </button>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto">
        {skills.length === 0 ? (
          <p className="px-2 py-6 text-center text-sm opacity-60">{t('settingsPage.skills.empty')}</p>
        ) : (
          <ul className="flex flex-col gap-1">
            {skills.map((s) => (
              <li key={s.name}>
                <SkillsListItem
                  item={s}
                  selected={selectedSkillName === s.name}
                  onSelect={setSelectedSkillName}
                  onToggle={handleToggle}
                />
              </li>
            ))}
          </ul>
        )}
      </div>
      <SkillImportDialog open={importOpen} onClose={() => setImportOpen(false)} onImported={() => { void refresh() }} />
    </div>
  )
}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/settings/skills/SkillsList.tsx frontend/src/components/settings/skills/SkillsListItem.tsx frontend/src/store/settingsStore.ts
git commit -m "feat(skills): add Skills list view and store slice"
```

---

### Task 12: `SkillsDetailPane` (read + edit + delete)

**Files:**
- Create: `frontend/src/components/settings/skills/SkillsDetailPane.tsx`
- Create: `frontend/src/components/settings/skills/SkillEditor.tsx`

- [ ] **Step 1: Create `SkillEditor.tsx`** — used by both "new" and "edit" modes.

```tsx
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { SkillFull } from '@/types/skills'

type Props = {
  initial: Partial<SkillFull>
  readOnly?: boolean
  onSubmit: (draft: { name: string; description: string; body: string }) => Promise<void> | void
  onCancel: () => void
}

export function SkillEditor({ initial, readOnly = false, onSubmit, onCancel }: Props) {
  const { t } = useTranslation()
  const [name, setName] = useState(initial.name ?? '')
  const [description, setDescription] = useState(initial.description ?? '')
  const [body, setBody] = useState(initial.body ?? '')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    setName(initial.name ?? '')
    setDescription(initial.description ?? '')
    setBody(initial.body ?? '')
  }, [initial.name, initial.description, initial.body])

  const handleSave = async () => {
    setSaving(true)
    try {
      await onSubmit({ name, description, body })
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex h-full flex-col gap-3 p-4">
      <label className="text-sm">
        {t('settingsPage.skills.editor.name')}
        <input className="mt-1 w-full rounded-md border px-2 py-1 text-sm"
          value={name} disabled={readOnly} onChange={(e) => setName(e.target.value)} />
      </label>
      <label className="text-sm">
        {t('settingsPage.skills.editor.description')}
        <textarea className="mt-1 w-full rounded-md border px-2 py-1 text-sm" rows={2}
          value={description} disabled={readOnly} onChange={(e) => setDescription(e.target.value)} />
      </label>
      <label className="flex flex-1 flex-col text-sm">
        {t('settingsPage.skills.editor.body')}
        <textarea className="mt-1 min-h-[200px] flex-1 rounded-md border p-2 font-mono text-xs"
          value={body} disabled={readOnly} onChange={(e) => setBody(e.target.value)} />
      </label>
      {!readOnly && (
        <div className="flex justify-end gap-2">
          <button type="button" className="rounded-md border px-3 py-1 text-sm" onClick={onCancel}>
            {t('settingsPage.skills.editor.cancel')}
          </button>
          <button type="button" className="rounded-md bg-primary px-3 py-1 text-sm text-primary-foreground disabled:opacity-50"
            onClick={() => { void handleSave() }} disabled={saving || !name || !description}>
            {t('settingsPage.skills.editor.save')}
          </button>
        </div>
      )}
    </div>
  )
}
```

- [ ] **Step 2: Create `SkillsDetailPane.tsx`**

```tsx
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Trash2 } from 'lucide-react'
import { useSettingsStore } from '@/store/settingsStore'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import {
  CreateSkillInput,
  DeleteSkillInput,
  GetSkillInput,
  UpdateSkillInput,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'
import { SkillEditor } from './SkillEditor'
import type { SkillFull, SkillItem } from '@/types/skills'

export function SkillsDetailPane() {
  const { t } = useTranslation()
  const selectedSkillName = useSettingsStore((s) => s.selectedSkillName)
  const setSelectedSkillName = useSettingsStore((s) => s.setSelectedSkillName)
  const setSkills = useSettingsStore((s) => s.setSkills)
  const removeSkill = useSettingsStore((s) => s.removeSkill)
  const updateSkill = useSettingsStore((s) => s.updateSkill)
  const [detail, setDetail] = useState<SkillFull | null>(null)
  const [loading, setLoading] = useState(false)

  const isNew = selectedSkillName === '__new__'

  const refreshList = useCallback(async () => {
    const result = await SkillsBinding.ListSkills({})
    setSkills(((result?.skills ?? []) as unknown as SkillItem[]))
  }, [setSkills])

  useEffect(() => {
    if (isNew || !selectedSkillName) {
      setDetail(null)
      return
    }
    setLoading(true)
    SkillsBinding.GetSkill(new GetSkillInput({ name: selectedSkillName }))
      .then((res) => setDetail((res?.skill ?? null) as unknown as SkillFull))
      .catch(() => setDetail(null))
      .finally(() => setLoading(false))
  }, [isNew, selectedSkillName])

  if (!selectedSkillName) {
    return <p className="p-6 text-sm opacity-60">{t('settingsPage.skills.empty')}</p>
  }

  if (isNew) {
    return (
      <SkillEditor
        initial={{}}
        onCancel={() => setSelectedSkillName(null)}
        onSubmit={async (draft) => {
          await SkillsBinding.CreateSkill(new CreateSkillInput({ ...draft, source: 'user' }))
          await refreshList()
          setSelectedSkillName(draft.name)
        }}
      />
    )
  }

  if (loading || !detail) {
    return <p className="p-6 text-sm opacity-60">…</p>
  }

  const readOnly = detail.source === 'builtin'

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b px-4 py-2">
        <h3 className="text-sm font-semibold">{detail.name}</h3>
        {!readOnly && (
          <button
            type="button"
            className="flex items-center gap-1 rounded-md p-1 text-red-600 hover:bg-red-50"
            onClick={async () => {
              if (!window.confirm(t('settingsPage.skills.deleteConfirm'))) return
              await SkillsBinding.DeleteSkill(new DeleteSkillInput({ name: detail.name }))
              removeSkill(detail.name)
              setSelectedSkillName(null)
            }}
          >
            <Trash2 size={14} /> {t('settingsPage.skills.delete')}
          </button>
        )}
      </div>
      {readOnly && (
        <p className="px-4 py-2 text-xs opacity-60">{t('settingsPage.skills.builtinLockedMsg')}</p>
      )}
      <SkillEditor
        initial={detail}
        readOnly={readOnly}
        onCancel={() => setSelectedSkillName(null)}
        onSubmit={async (draft) => {
          const result = await SkillsBinding.UpdateSkill(new UpdateSkillInput({
            original_name: detail.name,
            name: draft.name,
            description: draft.description,
            body: draft.body,
          }))
          if (result?.skill) {
            updateSkill(result.skill as unknown as SkillItem)
            setSelectedSkillName(result.skill.name)
          }
        }}
      />
    </div>
  )
}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/settings/skills/SkillsDetailPane.tsx frontend/src/components/settings/skills/SkillEditor.tsx
git commit -m "feat(skills): add skill detail/edit pane"
```

---

### Task 13: `SkillImportDialog`

**Files:**
- Create: `frontend/src/components/settings/skills/SkillImportDialog.tsx`

- [ ] **Step 1: Create the dialog component**

```tsx
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { ImportSkillInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto/models'

type Props = {
  open: boolean
  onClose: () => void
  onImported: () => void
}

export function SkillImportDialog({ open, onClose, onImported }: Props) {
  const { t } = useTranslation()
  const [content, setContent] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  if (!open) return null

  const handleSubmit = async () => {
    setSubmitting(true)
    setError(null)
    try {
      await SkillsBinding.ImportSkill(new ImportSkillInput({ content }))
      setContent('')
      onImported()
      onClose()
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-[600px] max-w-[90vw] rounded-md bg-card p-4 shadow-xl">
        <h2 className="mb-3 text-sm font-semibold">{t('settingsPage.skills.importDialog.title')}</h2>
        <textarea
          className="h-64 w-full rounded-md border p-2 font-mono text-xs"
          placeholder={t('settingsPage.skills.importDialog.placeholder')}
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
        {error && <p className="mt-2 text-xs text-red-600">{error}</p>}
        <div className="mt-3 flex justify-end gap-2">
          <button type="button" className="rounded-md border px-3 py-1 text-sm" onClick={onClose}>
            {t('settingsPage.skills.editor.cancel')}
          </button>
          <button
            type="button"
            className="rounded-md bg-primary px-3 py-1 text-sm text-primary-foreground disabled:opacity-50"
            onClick={() => { void handleSubmit() }}
            disabled={submitting || !content}
          >
            {t('settingsPage.skills.importDialog.submit')}
          </button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Build front-end check**

```
cd frontend && npm run build
```

Expect success now (all referenced components exist).

- [ ] **Step 3: Smoke test in dev**

Run the app, open Settings → Skills, verify: empty list, "+", "↑" buttons render, import dialog opens / closes, can create a new user skill and see it appear, toggle disable works, delete works.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/settings/skills/SkillImportDialog.tsx
git commit -m "feat(skills): add skill import dialog"
```

---

## Section 3 — Agent integration

### Task 14: Register `Skill` meta-tool in agent

**Files:**
- Modify: `backend/pkg/agent/manager.go` (or wherever built-in tools are registered)
- Create: `backend/pkg/agent/tools/skill_tool.go`

- [ ] **Step 1: Verify where built-in tools register**

```
grep -nR "RegisterBuiltin\|tools.NewRegistry\|registry.Register" backend/pkg/agent/
```

This identifies the registration site. Quote one line in your commit message for the reviewer.

- [ ] **Step 2: Create `skill_tool.go`**

```go
package tools

import (
	"context"
	"encoding/json"
	"errors"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

// SkillProvider abstracts the manager so agent tests can stub it.
type SkillProvider interface {
	Enabled() []pkgskills.Skill
	Get(name string) (pkgskills.Skill, bool)
}

// SkillToolName is the canonical name of the meta-tool used by the agent.
const SkillToolName = "Skill"

// BuildSkillTool constructs the Skill meta-tool. The description is rebuilt at send time
// (the agent should call BuildSkillTool fresh on each turn rather than caching the meta).
func BuildSkillTool(provider SkillProvider) ToolMeta {
	enabled := provider.Enabled()
	desc := "Load a skill (a set of natural-language instructions for a specific task). Call with {\"name\": \"<skill-name>\"}. Available skills:\n"
	if len(enabled) == 0 {
		desc += "(none)\n"
	}
	for _, s := range enabled {
		desc += "- " + s.Name + ": " + s.Description + "\n"
	}
	return ToolMeta{
		Name:        SkillToolName,
		Description: desc,
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				Name string `json:"name"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Loading skill: " + parsed.Name
		},
	}
}

// InvokeSkill returns the SKILL.md body of the requested skill, used by the agent's
// tool-call handler when the model invokes Skill({name}).
func InvokeSkill(_ context.Context, provider SkillProvider, args json.RawMessage) (string, error) {
	var parsed struct {
		Name string `json:"name"`
	}
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
```

- [ ] **Step 3: Inject the provider into the agent manager**

In the agent's `manager.go` (or wherever the registry is built), accept a `SkillProvider` (use the live `Skills` service / `pkgskills.Manager`). At the place that constructs the per-request tool list, append `BuildSkillTool(provider)`.

If the agent already takes a closure for per-turn tool lookup, pass `BuildSkillTool` there. If not — add one field to the agent struct:

```go
skillProvider tools.SkillProvider
```

And add a setter `SetSkillProvider` plus call it from `main.go` after wiring services (since both `agent` and `skills` services exist by then).

- [ ] **Step 4: Hook tool dispatch**

In the file that handles tool-call dispatch (search `case "Shell":` or similar), add a branch:

```go
case tools.SkillToolName:
    result, err := tools.InvokeSkill(ctx, a.skillProvider, call.Args)
    // … return as a normal tool result
```

- [ ] **Step 5: Wire into main.go**

After both services are constructed but before `app.Run()`, call the setter:

```go
skillsService := skillsSvc.NewSkills()
agentService := agentSvc.NewAgent(istorage)
agentService.SetSkillProvider(skillsService.Manager())
```

This requires adding `Manager()` accessor on `service/skills.Skills` and using `agentService` in the `Services` slice (instead of inline construction).

- [ ] **Step 6: Add a unit test on `BuildSkillTool`**

`backend/pkg/agent/tools/skill_tool_test.go`:

```go
package tools

import (
	"strings"
	"testing"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

type stubProvider struct{ items []pkgskills.Skill }

func (s stubProvider) Enabled() []pkgskills.Skill { return s.items }
func (s stubProvider) Get(name string) (pkgskills.Skill, bool) {
	for _, sk := range s.items {
		if sk.Name == name {
			return sk, true
		}
	}
	return pkgskills.Skill{}, false
}

func TestBuildSkillTool_ListsEnabled(t *testing.T) {
	p := stubProvider{items: []pkgskills.Skill{{Name: "foo", Description: "Foo desc"}}}
	meta := BuildSkillTool(p)
	if meta.Name != SkillToolName {
		t.Fatalf("name: %q", meta.Name)
	}
	if !strings.Contains(meta.Description, "foo: Foo desc") {
		t.Fatalf("description missing skill: %s", meta.Description)
	}
}
```

- [ ] **Step 7: Run all backend tests**

```
go test ./...
```

- [ ] **Step 8: Commit**

```bash
git add backend/pkg/agent/tools/skill_tool.go backend/pkg/agent/tools/skill_tool_test.go backend/pkg/agent/manager.go backend/service/skills/skills.go main.go
git commit -m "feat(skills): register Skill meta-tool in agent"
```

---

### Task 15: TipTap `/<skill-name>` slash command in chat input

**Files:**
- Modify: `frontend/src/components/chat/ChatInput.tsx`
- Create: `frontend/src/components/chat/slashSuggestion.ts`
- Create: `frontend/src/components/chat/SlashSuggestionPopup.tsx`

- [ ] **Step 1: Add suggestion config**

Create `slashSuggestion.ts`:

```ts
import { ReactRenderer } from '@tiptap/react'
import tippy, { type Instance as TippyInstance } from 'tippy.js'
import type { SuggestionOptions } from '@tiptap/suggestion'
import { Skills as SkillsBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills'
import { SlashSuggestionPopup, type SlashSuggestionItem } from './SlashSuggestionPopup'

export function buildSlashSuggestion(): Omit<SuggestionOptions, 'editor'> {
  return {
    char: '/',
    startOfLine: false,
    items: async ({ query }: { query: string }) => {
      const list = await SkillsBinding.ListSkills({})
      const enabled = (list?.skills ?? []).filter((s) => !s.disabled)
      return enabled
        .filter((s) => s.name.toLowerCase().startsWith(query.toLowerCase()))
        .slice(0, 10)
        .map((s) => ({ name: s.name, description: s.description })) as SlashSuggestionItem[]
    },
    render: () => {
      let component: ReactRenderer
      let popup: TippyInstance[] = []
      return {
        onStart: (props: { clientRect?: () => DOMRect | null }) => {
          component = new ReactRenderer(SlashSuggestionPopup, { props, editor: undefined as never })
          popup = tippy('body', {
            getReferenceClientRect: props.clientRect as () => DOMRect,
            appendTo: () => document.body,
            content: component.element,
            showOnCreate: true,
            interactive: true,
            trigger: 'manual',
            placement: 'top-start',
          })
        },
        onUpdate: (props) => {
          component.updateProps(props)
          popup[0]?.setProps({ getReferenceClientRect: props.clientRect as () => DOMRect })
        },
        onKeyDown: (props) => (component.ref as { onKeyDown?: (p: unknown) => boolean })?.onKeyDown?.(props) ?? false,
        onExit: () => {
          popup[0]?.destroy()
          component.destroy()
        },
      }
    },
  }
}
```

- [ ] **Step 2: Create `SlashSuggestionPopup.tsx`**

```tsx
import { forwardRef, useImperativeHandle, useState } from 'react'

export type SlashSuggestionItem = {
  name: string
  description: string
}

type Props = {
  items: SlashSuggestionItem[]
  command: (item: SlashSuggestionItem) => void
}

export const SlashSuggestionPopup = forwardRef<unknown, Props>(({ items, command }, ref) => {
  const [index, setIndex] = useState(0)

  useImperativeHandle(ref, () => ({
    onKeyDown: ({ event }: { event: KeyboardEvent }) => {
      if (event.key === 'ArrowDown') {
        setIndex((i) => (i + 1) % items.length)
        return true
      }
      if (event.key === 'ArrowUp') {
        setIndex((i) => (i + items.length - 1) % items.length)
        return true
      }
      if (event.key === 'Enter') {
        const target = items[index]
        if (target) command(target)
        return true
      }
      return false
    },
  }), [items, index, command])

  if (items.length === 0) return null

  return (
    <div className="max-h-60 w-72 overflow-y-auto rounded-md border bg-popover p-1 shadow-lg">
      {items.map((it, i) => (
        <button
          key={it.name}
          type="button"
          onClick={() => command(it)}
          className={`flex w-full flex-col items-start gap-0.5 rounded-md px-2 py-1 text-left text-xs ${
            i === index ? 'bg-accent text-accent-foreground' : 'hover:bg-muted'
          }`}
        >
          <span className="font-medium">/{it.name}</span>
          <span className="line-clamp-2 opacity-70">{it.description}</span>
        </button>
      ))}
    </div>
  )
})
```

- [ ] **Step 3: Wire into ChatInput**

In `frontend/src/components/chat/ChatInput.tsx` find the TipTap editor setup. Add a Suggestion extension:

```tsx
import Suggestion from '@tiptap/suggestion'
import { buildSlashSuggestion } from './slashSuggestion'

// inside useEditor's `extensions` list, append:
Suggestion.configure({
  ...buildSlashSuggestion(),
  command: ({ editor, range, props }: { editor: { chain: () => any }; range: { from: number; to: number }; props: { name: string } }) => {
    editor.chain().focus().deleteRange(range).insertContent(`/${props.name} `).run()
  },
}),
```

Confirm `@tiptap/suggestion` and `tippy.js` are in `frontend/package.json`. If not:

```bash
cd frontend && npm install @tiptap/suggestion tippy.js
```

- [ ] **Step 4: Smoke test in dev**

Run the app, open chat, type `/`. Verify popup shows enabled skills. Arrow keys + Enter select; the chosen skill name is inserted with leading `/`. (How the agent interprets `/name` is intentionally light in Plan A — Plan C strengthens the dispatch.)

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/chat/ChatInput.tsx frontend/src/components/chat/slashSuggestion.ts frontend/src/components/chat/SlashSuggestionPopup.tsx frontend/package.json frontend/package-lock.json
git commit -m "feat(skills): add /skill-name slash suggestion in chat input"
```

---

## Section 4 — End-to-end smoke test + close-out

### Task 16: Manual end-to-end test

- [ ] **Step 1: Restart app in dev mode**

```
task dev
```

- [ ] **Step 2: Settings → Skills**

Verify:
- Page renders with empty list (no built-in skill yet — that lands in Plan C).
- "+" opens new-skill editor.
- Create a skill `test-skill` with description `Test skill` and body `Hello world`. Save → it appears in the list.
- Click `test-skill` → editor loads body.
- Edit description → Save → list reflects new description.
- Toggle off → list shows disabled; toggle on → enabled.
- Click Delete → confirm → skill disappears + dir removed under `{data_dir}/skills/test-skill/`.

- [ ] **Step 3: Skill import**

Click ↑. Paste a valid SKILL.md (frontmatter + body). Submit → it appears as a `user` source skill.

- [ ] **Step 4: Chat agent**

- Create / select a skill `quick-greet`, description `Greet briefly`, body `Just say hi to the user.`.
- Open chat, type `/qui` → popup shows `quick-greet`. Press Enter → `/quick-greet` inserted.
- Send a normal message and ask the model `用 quick-greet skill 给我打个招呼`. Verify the model calls the `Skill` tool with `{"name":"quick-greet"}` (visible in tool-call panel) and the response is short / greeting-ish.

- [ ] **Step 5: Persistence**

Restart the app. Verify:
- Created skills still appear.
- Disabled state persists.

- [ ] **Step 6: Commit any small fixes uncovered**

If smoke test reveals issues, fix and commit each fix as its own commit referencing the task it amends (e.g. `fix(skills): correct list refresh after delete`).

---

## Self-review (do this before declaring Plan A done)

Run the checklist below; fix inline.

1. **Spec coverage** — confirm each spec item maps to a task above:
   - §3.4.1 (Claude Code-compatible form) → Tasks 2, 4.
   - §3.4.2 (data dir layout) → Tasks 1, 3.
   - §3.4.3 (3 sources, MVP) → Task 4 (builtin), Task 11/12 (user-import/create), Task 7 (`CreateSkill` accepts `source: "ai"`). User-confirmation for AI source intentionally deferred to Plan C.
   - §3.4.4 (triggers: AI auto / `/skill-name` / app-internal silent) → Task 14 (AI auto via Skill tool list), Task 15 (`/skill-name`). App-internal silent trigger is the Plan C entry point.
   - §3.4.5 (independent Skills settings page) → Tasks 10–13.
   - §5 (no marketplace, no version mgmt) → confirmed absent.

2. **Placeholder scan** — search the plan for TBD / TODO / "appropriate" / "similar to". Should find zero outside of "tested manually" prose.

3. **Type consistency** — check `Skill` struct fields (Name, Description, Body, Source, Disabled, UpdatedAt) match across Go, DTO, and TypeScript types. The TS `SkillItem` does NOT include body (only `has_body`); `SkillFull` adds body. Verify both sides agree.

4. **Naming consistency** — `SkillToolName = "Skill"` (capitalized), DTO methods are `ListSkills`/`GetSkill`/etc, store fields use `skills`/`selectedSkillName`. Confirm.

If any check fails, fix in the plan before handing off.

---

## Execution handoff

Plan complete at `docs/superpowers/plans/2026-05-25-skills-foundation.md`. Two execution options:

1. **Subagent-Driven (recommended)** — dispatch a fresh subagent per task, review between tasks, fast iteration.
2. **Inline Execution** — execute tasks in this session using executing-plans, batch with checkpoints.

Pick one to start Plan A.
