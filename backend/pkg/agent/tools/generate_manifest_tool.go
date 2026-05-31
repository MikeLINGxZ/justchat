package tools

import (
	"context"
	"encoding/json"
	"errors"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// CliManifestGenerator is satisfied by the plugin service.
type CliManifestGenerator interface {
	GenerateCliManifestSync(ctx context.Context, extensionID string) (extensionJSON string, err error)
}

const GenerateCliManifestToolName = "GenerateCliManifest"

func BuildGenerateCliManifestTool() ToolMeta {
	return ToolMeta{
		Name:        GenerateCliManifestToolName,
		Description: "Generate (or regenerate) the manifest for an installed CLI plugin. Call with {\"id\":\"cli:<name>\"}. Returns the updated extension item JSON.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				ID string `json:"id"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Generating manifest for: " + parsed.ID
		},
	}
}

type generateCliManifestInput struct {
	ID string `json:"id" jsonschema:"description=extension ID in the form cli:<name>,required"`
}

func InvokeGenerateCliManifest(ctx context.Context, generator CliManifestGenerator, args json.RawMessage) (string, error) {
	var parsed generateCliManifestInput
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.ID == "" {
		return "", errors.New("id is required")
	}
	return generator.GenerateCliManifestSync(ctx, parsed.ID)
}

// NewGenerateCliManifestTool creates a function tool that generates a CLI plugin manifest.
func NewGenerateCliManifestTool(generator CliManifestGenerator) *function.FunctionTool[generateCliManifestInput, string] {
	meta := BuildGenerateCliManifestTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input generateCliManifestInput) (string, error) {
			payload, _ := json.Marshal(input)
			return InvokeGenerateCliManifest(ctx, generator, payload)
		},
		function.WithName(GenerateCliManifestToolName),
		function.WithDescription(meta.Description),
	)
}
