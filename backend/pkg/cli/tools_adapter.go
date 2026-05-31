package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type toolRunner interface {
	RunTool(ctx context.Context, name string, toolName string, input map[string]any) (RunResult, error)
}

type cliToolSet struct {
	name  string
	tools []tool.Tool
}

// Tools returns the currently available CLI-backed tools.
func (s cliToolSet) Tools(context.Context) []tool.Tool {
	return s.tools
}

// Close releases resources owned by the tool set.
func (s cliToolSet) Close() error {
	return nil
}

// Name returns the unique tool set name.
func (s cliToolSet) Name() string {
	return s.name
}

// BuildToolSet creates an agent ToolSet for one enabled CLI extension.
func BuildToolSet(mgr *Manager, item data_models.ExtensionItem, manifest Manifest, enabledUserTools []string) (tool.ToolSet, error) {
	return buildToolSetFromRunner(mgr, item, manifest, enabledUserTools)
}

func buildToolSetFromRunner(runner toolRunner, item data_models.ExtensionItem, manifest Manifest, enabledUserTools []string) (tool.ToolSet, error) {
	setName := cliToolSetName(item.Name)
	if len(manifest.Tools) == 0 {
		return cliToolSet{name: setName, tools: []tool.Tool{}}, nil
	}

	includeAll := false
	selected := make(map[string]struct{}, len(enabledUserTools))
	for _, id := range enabledUserTools {
		selected[id] = struct{}{}
		if id == item.ID {
			includeAll = true
		}
	}

	tools := make([]tool.Tool, 0, len(manifest.Tools))
	runtimeName := cliRuntimeName(item)
	for _, current := range manifest.Tools {
		if !current.Enabled {
			continue
		}
		toolID := CliToolID(item.Name, current.Name)
		if len(selected) > 0 && !includeAll {
			if _, ok := selected[toolID]; !ok {
				continue
			}
		}

		schema, err := decodeInputSchema(current.InputSchema)
		if err != nil {
			return nil, err
		}

		toolName := current.Name
		extensionName := runtimeName
		tools = append(tools, function.NewFunctionTool(
			func(ctx context.Context, input map[string]any) (RunResult, error) {
				return runner.RunTool(ctx, extensionName, toolName, input)
			},
			function.WithName(current.Name),
			function.WithDescription(current.Description),
			function.WithInputSchema(schema),
		))
	}
	return cliToolSet{name: setName, tools: tools}, nil
}

func cliRuntimeName(item data_models.ExtensionItem) string {
	if base := strings.TrimSpace(filepath.Base(item.RootDir)); base != "" && base != "." && base != string(filepath.Separator) {
		return base
	}
	return item.Name
}

func decodeInputSchema(raw json.RawMessage) (*tool.Schema, error) {
	if len(raw) == 0 {
		return &tool.Schema{
			Type:       "object",
			Properties: map[string]*tool.Schema{},
		}, nil
	}
	var schema tool.Schema
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil, fmt.Errorf("decode cli input schema: %w", err)
	}
	if schema.Type == "" {
		schema.Type = "object"
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*tool.Schema{}
	}
	return &schema, nil
}

func cliToolSetName(name string) string {
	return sanitizeToolToken("cli_" + name)
}

// CliToolID builds the stable tool ID exposed to the chat agent.
func CliToolID(extensionName string, toolName string) string {
	return sanitizeToolToken("cli_" + extensionName + "_" + toolName)
}

func sanitizeToolToken(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))
	lastUnderscore := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
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
		return "cli"
	}
	return name
}
