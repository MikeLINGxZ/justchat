package tools

import (
	"context"
	"encoding/json"
	"errors"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// CliInstaller is satisfied by the plugin service (avoids import cycle).
type CliInstaller interface {
	InstallCliSync(ctx context.Context, npmPackage string, name string, onProgress func(phase string, detail string)) (installResultJSON string, err error)
}

const InstallCliToolName = "InstallCli"

func BuildInstallCliTool() ToolMeta {
	return ToolMeta{
		Name:        InstallCliToolName,
		Description: "Install a CLI plugin from npm. Call with {\"npm_package\":\"@scope/pkg\",\"name\":\"my-cli\"}. Returns the installed extension item JSON.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				NpmPackage string `json:"npm_package"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Installing CLI: " + parsed.NpmPackage
		},
	}
}

type installCliInput struct {
	NpmPackage string `json:"npm_package" jsonschema:"description=npm package name to install,required"`
	Name       string `json:"name" jsonschema:"description=short kebab-case name for the installed plugin"`
}

func InvokeInstallCli(ctx context.Context, installer CliInstaller, args json.RawMessage) (string, error) {
	var parsed installCliInput
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.NpmPackage == "" {
		return "", errors.New("npm_package is required")
	}
	return installer.InstallCliSync(ctx, parsed.NpmPackage, parsed.Name, nil)
}

// NewInstallCliTool creates a function tool that installs a CLI plugin from npm.
func NewInstallCliTool(installer CliInstaller) *function.FunctionTool[installCliInput, string] {
	meta := BuildInstallCliTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input installCliInput) (string, error) {
			payload, _ := json.Marshal(input)
			return InvokeInstallCli(ctx, installer, payload)
		},
		function.WithName(InstallCliToolName),
		function.WithDescription(meta.Description),
	)
}
