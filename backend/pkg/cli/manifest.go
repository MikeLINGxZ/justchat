// Package cli provides primitives for installing, executing and managing CLI plugin bundles.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OutputMode controls how RunTool parses a CLI's stdout before returning it to the caller.
type OutputMode string

const (
	OutputJSON  OutputMode = "json"
	OutputText  OutputMode = "text"
	OutputLines OutputMode = "lines"
)

// IsolationMode controls how the runner constructs the environment for the CLI subprocess.
type IsolationMode string

const (
	IsolationIsolated IsolationMode = "isolated"
	IsolationShared   IsolationMode = "shared"
)

// Tool describes one AI-callable tool exposed by a CLI plugin.
type Tool struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	InputSchema     json.RawMessage `json:"input_schema"`
	ArgvTemplate    []string        `json:"argv_template"`
	OutputMode      OutputMode      `json:"output_mode"`
	TimeoutSeconds  int             `json:"timeout_seconds,omitempty"`
	RequiresConfirm bool            `json:"requires_confirm"`
	Enabled         bool            `json:"enabled"`
}

// Manifest describes a fully installed CLI plugin and its exposed tools.
type Manifest struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Description  string        `json:"description"`
	Executable   string        `json:"executable"`
	LoginCommand []string      `json:"login_command,omitempty"`
	LoginSteps   [][]string    `json:"login_steps,omitempty"`
	Isolation    IsolationMode `json:"isolation"`
	Tools        []Tool        `json:"tools"`
}

// ResolveLoginSteps returns the ordered argv vectors to run for an interactive login.
// LoginSteps wins when non-empty; otherwise LoginCommand is wrapped as a single step.
// An empty result means the manifest has no login flow configured.
func (m Manifest) ResolveLoginSteps() [][]string {
	if len(m.LoginSteps) > 0 {
		out := make([][]string, 0, len(m.LoginSteps))
		for _, step := range m.LoginSteps {
			if len(step) == 0 {
				continue
			}
			out = append(out, step)
		}
		if len(out) > 0 {
			return out
		}
	}
	if len(m.LoginCommand) > 0 {
		return [][]string{m.LoginCommand}
	}
	return nil
}

// LoadManifest reads a manifest JSON file. A missing file returns an empty (zero-value) manifest with nil error.
func LoadManifest(path string) (Manifest, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Manifest{}, nil
		}
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(bytes, &m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// SaveManifest writes the manifest JSON to path atomically via tmp file + rename.
func SaveManifest(path string, m Manifest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, bytes, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// RepairManifestExecutable backfills a missing executable from the installed package metadata and persists the repair.
func RepairManifestExecutable(path string, installDir string, manifest Manifest) (Manifest, error) {
	if strings.TrimSpace(manifest.Executable) != "" {
		return manifest, nil
	}

	pkg, err := ReadPackageJSON(installDir)
	if err != nil {
		return Manifest{}, err
	}
	executable, err := SelectExecutable(installDir, pkg)
	if err != nil {
		return Manifest{}, err
	}
	manifest.Executable = executable
	if err := SaveManifest(path, manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

var placeholderRegex = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

// Validate checks the manifest for structural issues, currently:
// every {placeholder} in any tool's argv_template must reference a property declared in input_schema.
func Validate(m Manifest) error {
	if strings.TrimSpace(m.Executable) == "" {
		return errors.New("manifest: executable is empty")
	}
	for _, t := range m.Tools {
		props, err := schemaPropertyNames(t.InputSchema)
		if err != nil {
			return fmt.Errorf("manifest tool %q: parse input_schema: %w", t.Name, err)
		}
		for _, segment := range t.ArgvTemplate {
			matches := placeholderRegex.FindAllStringSubmatch(segment, -1)
			for _, m := range matches {
				if _, ok := props[m[1]]; !ok {
					return fmt.Errorf("manifest tool %q: argv_template uses {%s} but input_schema has no such property", t.Name, m[1])
				}
			}
		}
	}
	return nil
}

// schemaPropertyNames extracts the top-level "properties" keys from a JSON Schema object.
func schemaPropertyNames(raw json.RawMessage) (map[string]struct{}, error) {
	if len(raw) == 0 {
		return map[string]struct{}{}, nil
	}
	var doc struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	out := make(map[string]struct{}, len(doc.Properties))
	for k := range doc.Properties {
		out[k] = struct{}{}
	}
	return out, nil
}
