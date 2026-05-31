package mcp

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	agenttool "trpc.group/trpc-go/trpc-agent-go/tool"
	toolmcp "trpc.group/trpc-go/trpc-agent-go/tool/mcp"
)

// ImportInput describes a local MCP directory import request.
type ImportInput struct {
	SourceDir string
	Enabled   bool
}

// Manager handles MCP bundle importing and runtime helper operations.
type Manager struct {
	dataDir string
}

// NewManager creates a new MCP manager rooted at the given app data directory.
func NewManager(dataDir string) *Manager {
	return &Manager{dataDir: dataDir}
}

// ImportMCP copies a local MCP bundle into the app data directory and returns its persisted metadata.
func (m *Manager) ImportMCP(ctx context.Context, input ImportInput) (data_models.ExtensionItem, error) {
	_ = ctx

	manifest, err := loadManifestOptional(input.SourceDir)
	if err != nil {
		return data_models.ExtensionItem{}, err
	}
	if manifest.ConfigFile == "" {
		manifest.ConfigFile = "mcp.json"
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(input.SourceDir)
	}

	versionDir := strings.TrimSpace(manifest.Version)
	if versionDir == "" {
		versionDir = "default"
	}

	targetDir := filepath.Join(dir.MCPRoot(m.dataDir), manifest.Name, versionDir)
	if err := copyDirectory(input.SourceDir, targetDir); err != nil {
		return data_models.ExtensionItem{}, err
	}
	if err := writeManifest(targetDir, manifest); err != nil {
		return data_models.ExtensionItem{}, err
	}

	configPath := filepath.Join(targetDir, manifest.ConfigFile)
	if _, err := loadServerConfig(configPath); err != nil {
		return data_models.ExtensionItem{}, err
	}

	return data_models.ExtensionItem{
		ID:             buildExtensionID("mcp", manifest.Name, versionDir),
		Name:           manifest.Name,
		Description:    manifest.Description,
		Author:         manifest.Author,
		Version:        manifest.Version,
		Kind:           "mcp",
		Enabled:        input.Enabled,
		RuntimeStatus:  "idle",
		RuntimeMessage: "",
		RootDir:        targetDir,
		SourceDir:      input.SourceDir,
		ConfigFilePath: configPath,
		Tools:          []data_models.ExtensionTool{},
	}, nil
}

// ImportPlugin copies a local plugin bundle into the app data directory and returns its persisted metadata.
func (m *Manager) ImportPlugin(ctx context.Context, input ImportInput) (data_models.ExtensionItem, error) {
	_ = ctx

	manifest, err := loadManifestOptional(input.SourceDir)
	if err != nil {
		return data_models.ExtensionItem{}, err
	}

	versionDir := strings.TrimSpace(manifest.Version)
	if versionDir == "" {
		versionDir = "default"
	}

	targetDir := filepath.Join(dir.PluginRoot(m.dataDir), manifest.Name, versionDir)
	if err := copyDirectory(input.SourceDir, targetDir); err != nil {
		return data_models.ExtensionItem{}, err
	}

	return data_models.ExtensionItem{
		ID:             buildExtensionID("plugin", manifest.Name, versionDir),
		Name:           manifest.Name,
		Description:    manifest.Description,
		Author:         manifest.Author,
		Version:        manifest.Version,
		Kind:           "plugin",
		Enabled:        input.Enabled,
		RuntimeStatus:  "idle",
		RuntimeMessage: "",
		RootDir:        targetDir,
		SourceDir:      input.SourceDir,
		ConfigFilePath: "",
		Tools:          []data_models.ExtensionTool{},
	}, nil
}

func buildExtensionID(kind string, name string, version string) string {
	return fmt.Sprintf("%s:%s:%s", kind, name, version)
}

func mcpToolSetName(extensionID string) string {
	var builder strings.Builder
	builder.Grow(len(extensionID))

	lastUnderscore := false
	for _, r := range extensionID {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			builder.WriteByte('_')
			lastUnderscore = true
		}
	}

	name := strings.Trim(builder.String(), "_")
	if name == "" {
		return "mcp"
	}
	return name
}

func mcpToolNamePrefix(item data_models.ExtensionItem) string {
	return mcpToolSetName(item.ID) + "_"
}

// DiscoverTools initializes one MCP extension and returns the prefixed tool ids exposed to the agent.
func (m *Manager) DiscoverTools(ctx context.Context, item data_models.ExtensionItem) ([]data_models.ExtensionTool, error) {
	toolSet, err := m.BuildToolSet(item, nil)
	if err != nil {
		return nil, err
	}
	defer toolSet.Close()

	if initToolSet, ok := toolSet.(*toolmcp.ToolSet); ok {
		if err := initToolSet.Init(ctx); err != nil {
			return nil, err
		}
	}

	tools := toolSet.Tools(ctx)
	result := make([]data_models.ExtensionTool, 0, len(tools))
	safePrefix := mcpToolNamePrefix(item)
	legacyPrefix := item.ID + "_"
	for _, current := range tools {
		decl := current.Declaration()
		name := decl.Name
		switch {
		case strings.HasPrefix(name, safePrefix):
			name = strings.TrimPrefix(name, safePrefix)
		case strings.HasPrefix(name, legacyPrefix):
			name = strings.TrimPrefix(name, legacyPrefix)
		}
		result = append(result, data_models.ExtensionTool{
			ToolID:          decl.Name,
			ServerID:        item.ID,
			Name:            name,
			Description:     decl.Description,
			Enabled:         true,
			RequiresConfirm: false,
		})
	}
	return result, nil
}

// BuildToolSet creates an MCP toolset for one extension and optionally filters it to selected tool ids.
func (m *Manager) BuildToolSet(item data_models.ExtensionItem, enabledToolIDs []string) (agenttool.ToolSet, error) {
	config, err := loadServerConfig(item.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	safeToolSetName := mcpToolSetName(item.ID)
	options := []toolmcp.ToolSetOption{toolmcp.WithName(safeToolSetName)}
	if len(enabledToolIDs) > 0 {
		for _, toolID := range enabledToolIDs {
			if toolID == item.ID {
				return toolmcp.NewMCPToolSet(config.toToolSetConfig(), options...), nil
			}
		}

		rawNames := make([]string, 0, len(enabledToolIDs))
		safePrefix := safeToolSetName + "_"
		legacyPrefix := item.ID + "_"
		for _, toolID := range enabledToolIDs {
			switch {
			case strings.HasPrefix(toolID, safePrefix):
				rawNames = append(rawNames, strings.TrimPrefix(toolID, safePrefix))
			case strings.HasPrefix(toolID, legacyPrefix):
				rawNames = append(rawNames, strings.TrimPrefix(toolID, legacyPrefix))
			}
		}
		if len(rawNames) == 0 {
			return emptyToolSet{name: item.ID}, nil
		}
		options = append(options, toolmcp.WithToolFilterFunc(agenttool.NewIncludeToolNamesFilter(rawNames...)))
	}

	return toolmcp.NewMCPToolSet(config.toToolSetConfig(), options...), nil
}

type emptyToolSet struct {
	name string
}

func (e emptyToolSet) Tools(context.Context) []agenttool.Tool {
	return []agenttool.Tool{}
}

func (e emptyToolSet) Close() error {
	return nil
}

func (e emptyToolSet) Name() string {
	return e.name
}
