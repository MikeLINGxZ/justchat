package agents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

// CustomAgentsDir returns the directory for storing custom agent JSON files.
// Creates the directory if it does not exist.
func CustomAgentsDir() (string, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(dataPath, "agents", "custom")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// LoadCustomAgents reads all .json files from the custom agents directory
// and returns the parsed agent definitions.
func LoadCustomAgents() ([]CustomAgentDef, error) {
	dir, err := CustomAgentsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var agents []CustomAgentDef
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var def CustomAgentDef
		if err := json.Unmarshal(data, &def); err != nil {
			continue
		}
		agents = append(agents, def)
	}
	return agents, nil
}

// LoadCustomAgent reads a single custom agent definition by ID.
func LoadCustomAgent(id string) (*CustomAgentDef, error) {
	dir, err := CustomAgentsDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, id+".json"))
	if err != nil {
		return nil, err
	}
	var def CustomAgentDef
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, err
	}
	return &def, nil
}

// SaveCustomAgent writes a custom agent definition to disk as {id}.json.
func SaveCustomAgent(agent CustomAgentDef) error {
	dir, err := CustomAgentsDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(agent, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, agent.ID_+".json"), data, 0o644)
}

// DeleteCustomAgent removes a custom agent JSON file by ID.
func DeleteCustomAgent(id string) error {
	dir, err := CustomAgentsDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(dir, id+".json"))
}

// CustomAgentExists checks whether a custom agent file exists on disk.
func CustomAgentExists(id string) bool {
	dir, err := CustomAgentsDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, id+".json"))
	return err == nil
}

// SyncCustomAgentsToRegistry loads all custom agents from disk and
// synchronizes them into the global agent registry.
func SyncCustomAgentsToRegistry() {
	UnregisterAgentsByType(AgentTypeCustom)

	customAgents, err := LoadCustomAgents()
	if err != nil {
		return
	}
	for i := range customAgents {
		RegisterAgent(&customAgents[i])
	}
}
