package plugin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	agentpkg "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent"
	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// InstallCliFromNpm installs a published npm package as a CLI plugin and persists it to config.json.
func (p *Plugin) InstallCliFromNpm(ctx context.Context, input plugin_dto.InstallCliFromNpmInput) (*plugin_dto.InstallCliFromNpmOutput, error) {
	name := pickInstallName(input.Name, input.NpmPackage)
	mgr, err := p.resolveCliManager()
	if err != nil {
		return nil, err
	}
	res, err := mgr.InstallFromNpm(ctx, pkgcli.InstallParams{NpmPackage: input.NpmPackage, Name: name})
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliInstallFailed, err)
	}
	item, err := p.persistCliExtension(res, true)
	if err != nil {
		return nil, err
	}
	return &plugin_dto.InstallCliFromNpmOutput{Extension: item}, nil
}

// pickInstallName returns the user-provided name when non-empty, else falls back to the npm package name.
func pickInstallName(provided string, fallback string) string {
	if v := strings.TrimSpace(provided); v != "" {
		return sanitizeName(v)
	}
	return sanitizeName(fallback)
}

// sanitizeName strips characters not allowed in directory names. We only allow [a-zA-Z0-9_.-].
func sanitizeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-', r == '.':
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "cli-plugin"
	}
	return b.String()
}

// persistCliExtension converts a cli.InstallResult into an ExtensionItem and writes it to config.json.
// enabled controls the initial Enabled flag on the persisted item.
func (p *Plugin) persistCliExtension(res pkgcli.InstallResult, enabled bool) (data_models.ExtensionItem, error) {
	config, err := p.loadConfig()
	if err != nil {
		return data_models.ExtensionItem{}, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	item := data_models.ExtensionItem{
		ID:             "cli:" + res.Name,
		Name:           res.Name,
		Description:    res.Description,
		Author:         res.Author,
		Version:        res.Version,
		Kind:           "cli",
		Enabled:        enabled,
		RuntimeStatus:  "ready",
		RuntimeMessage: "",
		RootDir:        res.InstallDir,
		SourceDir:      "",
		ConfigFilePath: res.ManifestPath,
		Tools:          []data_models.ExtensionTool{},
	}
	config.Extensions = append(removeExtension(config.Extensions, item.ID), item)
	if err := p.saveConfig(config); err != nil {
		return data_models.ExtensionItem{}, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}
	return item, nil
}

// ResetCliData removes the cli_data/<name> subtree for the given CLI plugin (login state, manifest, XDG dirs).
func (p *Plugin) ResetCliData(ctx context.Context, input plugin_dto.ResetCliDataInput) (*plugin_dto.ResetCliDataOutput, error) {
	name, err := cliNameFromID(input.ID)
	if err != nil {
		return nil, err
	}
	mgr, err := p.resolveCliManager()
	if err != nil {
		return nil, err
	}
	if err := mgr.ResetData(ctx, name); err != nil {
		return nil, ierror.Error(ierror.ErrCliResetDataFailed, err)
	}
	return &plugin_dto.ResetCliDataOutput{}, nil
}

// UpdateCliManifest validates the submitted JSON, writes it via cli.SaveManifest, and refreshes the persisted ExtensionItem.
func (p *Plugin) UpdateCliManifest(ctx context.Context, input plugin_dto.UpdateCliManifestInput) (*plugin_dto.UpdateCliManifestOutput, error) {
	_ = ctx
	config, err := p.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	item, ok := findExtension(config.Extensions, input.ID)
	if !ok || item.Kind != "cli" {
		return nil, ierror.Error(ierror.ErrSettingsReadConfig, os.ErrNotExist)
	}
	var manifest pkgcli.Manifest
	if err := json.Unmarshal([]byte(input.ManifestText), &manifest); err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestInvalid, err)
	}
	if err := pkgcli.Validate(manifest); err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestInvalid, err)
	}
	if err := pkgcli.SaveManifest(item.ConfigFilePath, manifest); err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestSaveFailed, err)
	}
	// Mirror the manifest's name/version/description back onto the persisted item so the list view stays in sync.
	item.Name = manifest.Name
	item.Version = manifest.Version
	item.Description = manifest.Description
	item.Tools = projectCliTools(item, manifest)
	if item.Enabled {
		item.RuntimeStatus = "ready"
	} else {
		item.RuntimeStatus = "idle"
	}
	item.RuntimeMessage = ""
	config.Extensions = append(removeExtension(config.Extensions, item.ID), item)
	if err := p.saveConfig(config); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}
	return &plugin_dto.UpdateCliManifestOutput{Extension: item}, nil
}

// GenerateCliManifest probes the CLI help output, asks the default chat model to draft a manifest, and persists it.
func (p *Plugin) GenerateCliManifest(ctx context.Context, input plugin_dto.GenerateCliManifestInput) (*plugin_dto.GenerateCliManifestOutput, error) {
	config, err := p.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	item, ok := findExtension(config.Extensions, input.ID)
	if !ok || item.Kind != "cli" {
		return nil, ierror.Error(ierror.ErrSettingsReadConfig, os.ErrNotExist)
	}

	existingManifest, err := pkgcli.LoadManifest(item.ConfigFilePath)
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}
	pkgMeta, err := pkgcli.ReadPackageJSON(item.RootDir)
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}
	if existingManifest.Executable == "" {
		executable, err := pkgcli.SelectExecutable(item.RootDir, pkgMeta)
		if err != nil {
			return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
		}
		existingManifest.Executable = executable
	}
	runtimeState, err := p.loadPersistedRuntime()
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}
	cliName, err := cliNameFromID(item.ID)
	if err != nil {
		return nil, err
	}
	pluginBin := filepath.Join(item.RootDir, "node_modules", ".bin")
	runtimeBin := ""
	if runtimeState.NodePath != "" {
		runtimeBin = filepath.Dir(runtimeState.NodePath)
	}
	env := pkgcli.BuildEnv(pkgcli.EnvParams{
		Isolation:  pkgcli.IsolationShared,
		CLIDataDir: filepath.Join(mustDataDir(), "plugins", "cli_data", cliName),
		RuntimeBin: runtimeBin,
		PluginBin:  pluginBin,
		ParentEnv:  os.Environ(),
	})
	helpText, err := p.probeCliHelp(ctx, existingManifest.Executable, env)
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}
	modelCfg, err := p.resolveDefaultChatModel()
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}
	nextManifest, err := p.generateManifest(ctx, pkgcli.GenerateParams{
		HelpText:    helpText,
		PackageName: pkgMeta.Name,
		PackageMeta: pkgMeta,
		Executable:  existingManifest.Executable,
		Caller: func(ctx context.Context, system, user string) (string, error) {
			resp, err := agentpkg.OneshotComplete(ctx, agentpkg.OneshotRequest{
				BaseURL:      modelCfg.BaseURL,
				APIKey:       modelCfg.APIKey,
				ModelName:    modelCfg.ModelName,
				ProviderType: modelCfg.ProviderType,
				System:       system,
				User:         user,
				Timeout:      3 * time.Minute,
			})
			if err != nil {
				return "", err
			}
			return resp.Text, nil
		},
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestGenerateFailed, err)
	}

	// Preserve user-customized login config across regenerate. The generator
	// rarely understands multi-step flows (e.g. lark-cli's config init + auth
	// login). If the previous manifest had either field set, keep both as-is so
	// careful manual edits aren't silently clobbered.
	if len(existingManifest.LoginSteps) > 0 || len(existingManifest.LoginCommand) > 0 {
		nextManifest.LoginSteps = existingManifest.LoginSteps
		nextManifest.LoginCommand = existingManifest.LoginCommand
	}

	if err := pkgcli.SaveManifest(item.ConfigFilePath, nextManifest); err != nil {
		return nil, ierror.Error(ierror.ErrCliManifestSaveFailed, err)
	}

	item.Name = nextManifest.Name
	item.Version = nextManifest.Version
	item.Description = nextManifest.Description
	item.Tools = projectCliTools(item, nextManifest)
	if item.Enabled {
		item.RuntimeStatus = "ready"
	} else {
		item.RuntimeStatus = "idle"
	}
	item.RuntimeMessage = ""
	config.Extensions = append(removeExtension(config.Extensions, item.ID), item)
	if err := p.saveConfig(config); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}
	return &plugin_dto.GenerateCliManifestOutput{Extension: item}, nil
}

// cliNameFromID parses "cli:<name>" -> "<name>"; returns an error for non-cli IDs.
func cliNameFromID(id string) (string, error) {
	const prefix = "cli:"
	if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
		return "", ierror.Error(ierror.ErrSettingsReadConfig, os.ErrInvalid)
	}
	return id[len(prefix):], nil
}
