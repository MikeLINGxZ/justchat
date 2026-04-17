package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manifest represents a plugin's package.json with the lemontea extension field.
type Manifest struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	DisplayName string         `json:"displayName"`
	Description string         `json:"description"`
	Main        string         `json:"main"`
	LemonTea    LemonTeaConfig `json:"lemontea"`
}

// LemonTeaConfig holds the plugin-specific configuration under the "lemontea" key.
type LemonTeaConfig struct {
	Engine           string      `json:"engine"`
	ActivationEvents []string    `json:"activationEvents"`
	Contributes      Contributes `json:"contributes"`
}

// Contributes describes the extension points a plugin provides.
type Contributes struct {
	Tools  []ToolContrib  `json:"tools"`
	Agents []AgentContrib `json:"agents"`
	Views  ViewContribs   `json:"views"`
	Hooks  HookContrib    `json:"hooks"`
}

// ToolContrib describes a tool contributed by a plugin.
type ToolContrib struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// AgentContrib describes an agent contributed by a plugin.
type AgentContrib struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ViewContribs groups view contributions by location.
type ViewContribs struct {
	Sidebar  []ViewEntry `json:"sidebar"`
	ChatCard []ViewEntry `json:"chatCard"`
	Settings []ViewEntry `json:"settings"`
	Page     []ViewEntry `json:"page"`
}

// ViewEntry describes a single view contribution.
type ViewEntry struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Entry string `json:"entry"`
}

// HookContrib indicates which lifecycle hooks the plugin wants to participate in.
type HookContrib struct {
	OnBeforeChat bool `json:"onBeforeChat"`
	OnAfterChat  bool `json:"onAfterChat"`
}

// ParseManifest reads package.json from pluginDir, unmarshals it into a Manifest,
// and returns an error if the file is missing or the "lemontea" field is absent.
func ParseManifest(pluginDir string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(pluginDir, "package.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin manifest: %w", err)
	}

	// First, do a raw decode to check for the lemontea field.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse plugin manifest: %w", err)
	}

	lt, ok := raw["lemontea"]
	if !ok || string(lt) == "null" || string(lt) == "{}" {
		return nil, fmt.Errorf("plugin manifest missing required 'lemontea' field")
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin manifest: %w", err)
	}

	return &m, nil
}

// ValidateManifest checks that required fields (name, version, main, engine) are non-empty.
func ValidateManifest(m *Manifest) error {
	if m.Name == "" {
		return fmt.Errorf("manifest validation failed: 'name' is required")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest validation failed: 'version' is required")
	}
	if m.Main == "" {
		return fmt.Errorf("manifest validation failed: 'main' is required")
	}
	if m.LemonTea.Engine == "" {
		return fmt.Errorf("manifest validation failed: 'lemontea.engine' is required")
	}
	return nil
}

// FullToolID returns a namespaced tool identifier in the form "plugin:<pluginName>:<localToolID>".
func FullToolID(pluginName, localToolID string) string {
	return "plugin:" + pluginName + ":" + localToolID
}

// FullAgentID returns a namespaced agent identifier in the form "plugin:<pluginName>:<localAgentID>".
func FullAgentID(pluginName, localAgentID string) string {
	return "plugin:" + pluginName + ":" + localAgentID
}
