package plugin

import (
	"context"
	"encoding/json"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// InstallCliSync satisfies tools.CliInstaller. It wraps the existing InstallCliFromNpm and emits progress events.
func (p *Plugin) InstallCliSync(ctx context.Context, npmPackage string, name string, onProgress func(string, string)) (string, error) {
	emit := func(phase, detail string) {
		if onProgress != nil {
			onProgress(phase, detail)
		}
		if p.wailsApp != nil {
			p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
				"npm_package": npmPackage,
				"name":        name,
				"phase":       phase,
				"detail":      detail,
			})
		}
	}

	emit("downloading", "Fetching "+npmPackage)
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{
		NpmPackage: npmPackage,
		Name:       name,
	})
	if err != nil {
		emit("failed", err.Error())
		return "", err
	}
	emit("installed", "Package installed to "+out.Extension.RootDir)

	resultBytes, _ := json.Marshal(out.Extension)
	return string(resultBytes), nil
}

// ReportCliInstallProgress satisfies tools.CliInstallProgressReporter.
func (p *Plugin) ReportCliInstallProgress(ctx context.Context, sessionID uint, item tools.CliInstallProgressItem) error {
	_ = ctx
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.install.progress", map[string]any{
			"session_id":   sessionID,
			"npm_package":  item.NpmPackage,
			"name":         item.Name,
			"extension_id": item.ExtensionID,
			"phase":        item.Phase,
			"detail":       item.Detail,
			"action_url":   item.ActionURL,
			"action_label": item.ActionLabel,
			"expires_at":   item.ExpiresAt,
		})
	}
	return nil
}

// GenerateCliManifestSync satisfies tools.CliManifestGenerator.
func (p *Plugin) GenerateCliManifestSync(ctx context.Context, extensionID string) (string, error) {
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
			"id":    extensionID,
			"phase": "generating",
		})
	}
	out, err := p.GenerateCliManifest(ctx, plugin_dto.GenerateCliManifestInput{ID: extensionID})
	if err != nil {
		if p.wailsApp != nil {
			p.wailsApp.Event.Emit("cli.install.progress", map[string]string{
				"id":    extensionID,
				"phase": "failed",
			})
		}
		return "", err
	}
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.install.done", map[string]interface{}{
			"extension": out.Extension,
		})
	}
	resultBytes, _ := json.Marshal(out.Extension)
	return string(resultBytes), nil
}

// RunCliCommandSync satisfies tools.CliCommandRunner.
func (p *Plugin) RunCliCommandSync(ctx context.Context, extensionID string, argv []string, outputMode string, timeoutSeconds int, onProgress func(result string)) (string, error) {
	name, err := cliNameFromID(extensionID)
	if err != nil {
		return "", err
	}
	mgr, err := p.resolveCliManager()
	if err != nil {
		return "", err
	}

	mode := pkgcli.OutputText
	switch strings.TrimSpace(outputMode) {
	case "", string(pkgcli.OutputText):
		mode = pkgcli.OutputText
	case string(pkgcli.OutputJSON):
		mode = pkgcli.OutputJSON
	case string(pkgcli.OutputLines):
		mode = pkgcli.OutputLines
	}

	res, err := mgr.RunCommandStreaming(ctx, name, argv, mode, timeoutSeconds, func(progress pkgcli.RunProgress) {
		if onProgress == nil {
			return
		}
		onProgress(formatCLIProgress(progress))
	})
	if err != nil {
		return "", err
	}
	resultBytes, _ := json.Marshal(res)
	return string(resultBytes), nil
}

func formatCLIProgress(progress pkgcli.RunProgress) string {
	var terminalOutput string
	switch {
	case progress.Stdout == "" && progress.Stderr == "":
		return ""
	case progress.Stderr == "":
		terminalOutput = progress.Stdout
	case progress.Stdout == "":
		terminalOutput = "[stderr]\n" + progress.Stderr
	default:
		terminalOutput = progress.Stdout + "\n[stderr]\n" + progress.Stderr
	}
	payload := map[string]any{
		"interactive_terminal": true,
		"terminal_status":      "active",
		"terminal_output":      terminalOutput,
	}
	out, _ := json.Marshal(payload)
	return string(out)
}
