package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// InstallParams describes a single install request from either npm or a local directory.
// Exactly one of NpmPackage or LocalPath must be non-empty.
type InstallParams struct {
	NpmPackage string // e.g. "lark-cli", "@scope/pkg"
	LocalPath  string // absolute path to a dir containing package.json
	Name       string // user-chosen plugin name; install lands at plugins/cli/<Name>/
}

// InstallResult summarizes the on-disk artifacts after a successful install.
type InstallResult struct {
	Name         string
	Version      string
	Description  string
	Author       string
	InstallDir   string
	CLIDataDir   string
	ManifestPath string
	Executable   string
}

// Manager is the external entry point for CLI plugin lifecycle (install / uninstall / reset / run).
// It is constructed with the data dir and the absolute paths to the bundled node/npm binaries,
// so it has no compile-time dependency on the runtime service.
type Manager struct {
	dataDir  string
	nodePath string
	npmPath  string
}

// NewManager constructs a Manager bound to a data dir and the bundled node/npm binaries.
func NewManager(dataDir, nodePath, npmPath string) *Manager {
	return &Manager{dataDir: dataDir, nodePath: nodePath, npmPath: npmPath}
}

// InstallFromNpm installs a published npm package as a CLI plugin under plugins/cli/<Name>/.
func (m *Manager) InstallFromNpm(ctx context.Context, p InstallParams) (InstallResult, error) {
	if strings.TrimSpace(p.NpmPackage) == "" {
		return InstallResult{}, errors.New("cli: NpmPackage required")
	}
	return m.install(ctx, p.Name, p.NpmPackage)
}

// InstallFromLocal installs a local directory (containing package.json) as a CLI plugin.
func (m *Manager) InstallFromLocal(ctx context.Context, p InstallParams) (InstallResult, error) {
	if strings.TrimSpace(p.LocalPath) == "" {
		return InstallResult{}, errors.New("cli: LocalPath required")
	}
	if _, err := os.Stat(filepath.Join(p.LocalPath, "package.json")); err != nil {
		return InstallResult{}, errors.New("cli: source dir has no package.json")
	}
	return m.install(ctx, p.Name, p.LocalPath)
}

// install runs the common install pipeline: npm -> probe -> empty manifest -> result.
func (m *Manager) install(ctx context.Context, name, source string) (InstallResult, error) {
	if strings.TrimSpace(name) == "" {
		return InstallResult{}, errors.New("cli: Name required")
	}
	installDir := filepath.Join(dir.CLIRoot(m.dataDir), name)
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return InstallResult{}, err
	}
	if err := NpmInstall(ctx, NpmInstallParams{
		NodePath:  m.nodePath,
		NpmPath:   m.npmPath,
		Source:    source,
		TargetDir: installDir,
		ParentEnv: os.Environ(),
	}); err != nil {
		return InstallResult{}, err
	}

	pkg, err := ReadPackageJSON(installDir)
	if err != nil {
		return InstallResult{}, err
	}
	executable, err := SelectExecutable(installDir, pkg)
	if err != nil {
		return InstallResult{}, err
	}

	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	manifestPath := filepath.Join(cliDataDir, "manifest.json")

	existing, _ := LoadManifest(manifestPath)
	manifest := Manifest{
		Name:        pkg.Name,
		Version:     pkg.Version,
		Description: pkg.Description,
		Executable:  executable,
		Isolation:   existing.Isolation, // preserve user's isolated/shared choice across reinstall
		Tools:       existing.Tools,     // preserve any AI/user-edited tools across reinstall
	}
	if manifest.Isolation == "" {
		manifest.Isolation = IsolationIsolated
	}
	if err := SaveManifest(manifestPath, manifest); err != nil {
		return InstallResult{}, err
	}

	return InstallResult{
		Name:         name,
		Version:      pkg.Version,
		Description:  pkg.Description,
		Author:       pkg.Author,
		InstallDir:   installDir,
		CLIDataDir:   cliDataDir,
		ManifestPath: manifestPath,
		Executable:   executable,
	}, nil
}

// Uninstall removes plugins/cli/<name>/ but preserves plugins/cli_data/<name>/.
func (m *Manager) Uninstall(ctx context.Context, name string) error {
	_ = ctx
	if strings.TrimSpace(name) == "" {
		return errors.New("cli: name required")
	}
	installDir := filepath.Join(dir.CLIRoot(m.dataDir), name)
	return os.RemoveAll(installDir)
}

// ResetData removes plugins/cli_data/<name>/ (login tokens, generated manifest, XDG state).
// The install dir itself is left untouched.
func (m *Manager) ResetData(ctx context.Context, name string) error {
	_ = ctx
	if strings.TrimSpace(name) == "" {
		return errors.New("cli: name required")
	}
	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	return os.RemoveAll(cliDataDir)
}

// RunTool resolves the named tool from the plugin's manifest, renders argv with input substitution,
// and executes the CLI subprocess under the manifest's isolation mode.
func (m *Manager) RunTool(ctx context.Context, name string, toolName string, input map[string]any) (RunResult, error) {
	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	installDir := filepath.Join(dir.CLIRoot(m.dataDir), name)
	manifestPath := filepath.Join(cliDataDir, "manifest.json")

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return RunResult{}, err
	}
	if manifest.Executable == "" {
		manifest, err = RepairManifestExecutable(manifestPath, installDir, manifest)
		if err != nil {
			return RunResult{}, errors.New("cli: manifest has no executable; reinstall the plugin")
		}
	}

	var tool Tool
	found := false
	for _, t := range manifest.Tools {
		if t.Name == toolName {
			tool = t
			found = true
			break
		}
	}
	if !found {
		return RunResult{}, errors.New("cli: tool not found in manifest: " + toolName)
	}

	argv, err := renderArgv(tool, input)
	if err != nil {
		return RunResult{}, err
	}

	env := BuildEnv(EnvParams{
		Isolation:  manifest.Isolation,
		CLIDataDir: cliDataDir,
		RuntimeBin: filepath.Dir(m.nodePath),
		PluginBin:  filepath.Join(installDir, "node_modules", ".bin"),
		ParentEnv:  os.Environ(),
	})

	return Run(ctx, RunParams{
		Executable: manifest.Executable,
		Argv:       argv,
		Env:        env,
		OutputMode: tool.OutputMode,
		TimeoutSec: tool.TimeoutSeconds,
		UsePTY:     true,
	})
}

// RunCommand executes an arbitrary argv against one installed CLI using the managed runtime env.
// This is intended for setup/login/init flows before dedicated manifest tools are available.
func (m *Manager) RunCommand(ctx context.Context, name string, argv []string, outputMode OutputMode, timeoutSeconds int) (RunResult, error) {
	return m.runCommand(ctx, name, argv, outputMode, timeoutSeconds, nil)
}

// RunCommandStreaming executes an arbitrary argv and emits cumulative stdout/stderr snapshots while it runs.
func (m *Manager) RunCommandStreaming(ctx context.Context, name string, argv []string, outputMode OutputMode, timeoutSeconds int, onProgress func(RunProgress)) (RunResult, error) {
	return m.runCommand(ctx, name, argv, outputMode, timeoutSeconds, onProgress)
}

func (m *Manager) runCommand(ctx context.Context, name string, argv []string, outputMode OutputMode, timeoutSeconds int, onProgress func(RunProgress)) (RunResult, error) {
	if len(argv) == 0 {
		return RunResult{}, errors.New("cli: argv required")
	}

	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	installDir := filepath.Join(dir.CLIRoot(m.dataDir), name)
	manifestPath := filepath.Join(cliDataDir, "manifest.json")

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return RunResult{}, err
	}
	if manifest.Executable == "" {
		manifest, err = RepairManifestExecutable(manifestPath, installDir, manifest)
		if err != nil {
			return RunResult{}, errors.New("cli: manifest has no executable; reinstall the plugin")
		}
	}

	env := BuildEnv(EnvParams{
		Isolation:  manifest.Isolation,
		CLIDataDir: cliDataDir,
		RuntimeBin: filepath.Dir(m.nodePath),
		PluginBin:  filepath.Join(installDir, "node_modules", ".bin"),
		ParentEnv:  os.Environ(),
	})

	return RunWithProgress(ctx, RunParams{
		Executable: manifest.Executable,
		Argv:       argv,
		Env:        env,
		OutputMode: outputMode,
		TimeoutSec: timeoutSeconds,
		WorkingDir: installDir,
		UsePTY:     onProgress != nil,
	}, onProgress)
}

// renderArgv replaces every {placeholder} segment in argv_template with the corresponding input field's stringified value.
// Behavior on missing keys:
//   - If the key is listed in input_schema.required, returns an error (caller must supply it).
//   - Otherwise the segment is dropped; if it was a pure {placeholder} preceded by a flag-looking segment
//     (begins with '-'), that flag is also dropped, so an optional `--flag {value}` pair vanishes cleanly.
func renderArgv(t Tool, input map[string]any) ([]string, error) {
	required, err := schemaRequired(t.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("cli: tool %q: parse input_schema.required: %w", t.Name, err)
	}

	out := make([]string, 0, len(t.ArgvTemplate))
	for _, segment := range t.ArgvTemplate {
		matches := placeholderRegex.FindAllStringSubmatchIndex(segment, -1)
		if len(matches) == 0 {
			out = append(out, segment)
			continue
		}

		var missingOptional bool
		for _, m := range matches {
			key := segment[m[2]:m[3]]
			if _, present := input[key]; present {
				continue
			}
			if _, req := required[key]; req {
				return nil, errors.New("cli: argv_template references missing input field: " + key)
			}
			missingOptional = true
		}

		if missingOptional {
			if isPureValuePlaceholder(segment) && len(out) > 0 && isFlagSegment(out[len(out)-1]) {
				out = out[:len(out)-1]
			}
			continue
		}

		var builder strings.Builder
		cursor := 0
		for _, m := range matches {
			start, end := m[0], m[1]
			builder.WriteString(segment[cursor:start])
			builder.WriteString(stringifyInput(input[segment[m[2]:m[3]]]))
			cursor = end
		}
		builder.WriteString(segment[cursor:])
		out = append(out, builder.String())
	}
	return out, nil
}

// schemaRequired extracts the "required" array of an input_schema, returning a set for O(1) lookup.
func schemaRequired(raw json.RawMessage) (map[string]struct{}, error) {
	if len(raw) == 0 {
		return map[string]struct{}{}, nil
	}
	var doc struct {
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	out := make(map[string]struct{}, len(doc.Required))
	for _, name := range doc.Required {
		out[name] = struct{}{}
	}
	return out, nil
}

// isPureValuePlaceholder reports whether the segment is exactly one {placeholder} with no surrounding literal text.
// Only pure-placeholder segments pop their preceding flag when dropped, so mixed segments like "--key={val}" do not.
func isPureValuePlaceholder(segment string) bool {
	loc := placeholderRegex.FindStringIndex(segment)
	return loc != nil && loc[0] == 0 && loc[1] == len(segment)
}

// isFlagSegment reports whether the segment looks like a CLI option flag, e.g. "-v" or "--data".
func isFlagSegment(segment string) bool {
	return len(segment) >= 2 && segment[0] == '-'
}

// stringifyInput converts an input value (string, number, bool, etc.) to its string form for argv use.
func stringifyInput(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
}
