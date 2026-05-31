package dir

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// defaultBaseDirName defines the default name of the base directory used for storing application data.
const defaultBaseDirName = ".lemontea"

// ConfigFileName defines the name of the configuration file.
const ConfigFileName = "config.json"

// DataBaseFileName defines the name of the database file.
const DataBaseFileName = "data.db"

type dataDirLocator struct {
	DataDir string `json:"data_dir"`
}

// GetDefaultBaseDir returns the stable home-based metadata directory used by Lemontea.
func GetDefaultBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataPath := filepath.Join(homeDir, defaultBaseDirName)
	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return "", err
	}
	return dataPath, nil
}

// GetLocatorFilePath returns the path of the custom data-directory locator file.
func GetLocatorFilePath() (string, error) {
	baseDir, err := GetDefaultBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "data_dir.json"), nil
}

// ReadLocatorDataDir reads the custom data-directory locator and returns an empty string when it does not exist.
func ReadLocatorDataDir() (string, error) {
	locatorPath, err := GetLocatorFilePath()
	if err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(locatorPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var locator dataDirLocator
	if err := json.Unmarshal(bytes, &locator); err != nil {
		return "", err
	}

	return locator.DataDir, nil
}

// WriteLocatorDataDir persists the active custom data directory inside the stable home-based metadata folder.
func WriteLocatorDataDir(dataDir string) error {
	locatorPath, err := GetLocatorFilePath()
	if err != nil {
		return err
	}

	content, err := json.Marshal(dataDirLocator{DataDir: dataDir})
	if err != nil {
		return err
	}

	return os.WriteFile(locatorPath, content, 0o644)
}

// GetDataDir get the data storage directory
func GetDataDir() (string, error) {
	if os.Getenv("LEMONTEA_DATA_DIR") != "" {
		if err := os.MkdirAll(os.Getenv("LEMONTEA_DATA_DIR"), 0o755); err != nil {
			return "", err
		}
		return os.Getenv("LEMONTEA_DATA_DIR"), nil
	}

	locatorDir, err := ReadLocatorDataDir()
	if err != nil {
		return "", err
	}
	if locatorDir != "" {
		if err := os.MkdirAll(locatorDir, 0o755); err != nil {
			return "", err
		}
		return locatorDir, nil
	}

	dataPath, err := GetDefaultBaseDir()
	if err != nil {
		return "", err
	}
	return dataPath, nil
}

// HomeDir get the home directory
func HomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDir
}

// PluginsSubdirName is the unified parent directory that holds all extension subdirs under data_dir.
const PluginsSubdirName = "plugins"

// MCPSubdirName is the directory name for MCP bundles under the plugins root.
const MCPSubdirName = "mcp"

// PluginSubdirName is the directory name for Node.js plugins under the plugins root.
const PluginSubdirName = "plugin"

// CLISubdirName is the directory name for CLI plugin installs under the plugins root.
const CLISubdirName = "cli"

// CLIDataSubdirName is the directory name for CLI plugin persistent data under the plugins root.
const CLIDataSubdirName = "cli_data"

// ExtensionsRoot returns the unified extensions parent dir: {dataDir}/plugins.
func ExtensionsRoot(dataDir string) string {
	return filepath.Join(dataDir, PluginsSubdirName)
}

// MCPRoot returns the MCP bundles root: {dataDir}/plugins/mcp.
func MCPRoot(dataDir string) string {
	return filepath.Join(ExtensionsRoot(dataDir), MCPSubdirName)
}

// PluginRoot returns the Node.js plugins root: {dataDir}/plugins/plugin.
func PluginRoot(dataDir string) string {
	return filepath.Join(ExtensionsRoot(dataDir), PluginSubdirName)
}

// CLIRoot returns the CLI plugin install root: {dataDir}/plugins/cli.
func CLIRoot(dataDir string) string {
	return filepath.Join(ExtensionsRoot(dataDir), CLISubdirName)
}

// CLIDataRoot returns the CLI plugin persistent data root: {dataDir}/plugins/cli_data.
func CLIDataRoot(dataDir string) string {
	return filepath.Join(ExtensionsRoot(dataDir), CLIDataSubdirName)
}

// LegacyMCPRoot returns the pre-migration MCP root used by older builds.
func LegacyMCPRoot(dataDir string) string {
	return filepath.Join(dataDir, MCPSubdirName)
}

// LegacyPluginRoot returns the pre-migration Node.js plugin root used by older builds.
func LegacyPluginRoot(dataDir string) string {
	return filepath.Join(dataDir, PluginSubdirName)
}

// SkillsSubdirName is the directory name for skills under data_dir.
const SkillsSubdirName = "skills"

// SkillsRoot returns the skills root: {dataDir}/skills.
func SkillsRoot(dataDir string) string {
	return filepath.Join(dataDir, SkillsSubdirName)
}
