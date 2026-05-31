//go:build !windows

package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// StartLoginCommand loads the named plugin's manifest, builds env per manifest.Isolation,
// and spawns a LoginSession running `manifest.Executable + argv`. The caller supplies
// argv (typically one step from manifest.ResolveLoginSteps()), keeping multi-step
// orchestration out of this layer. Concurrency-conflict detection is the service layer's job.
func (m *Manager) StartLoginCommand(ctx context.Context, name string, argv []string) (*LoginSession, error) {
	if len(argv) == 0 {
		return nil, errors.New("cli: login step has empty argv")
	}
	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	installDir := filepath.Join(dir.CLIRoot(m.dataDir), name)
	manifestPath := filepath.Join(cliDataDir, "manifest.json")

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	if manifest.Executable == "" {
		manifest, err = RepairManifestExecutable(manifestPath, installDir, manifest)
		if err != nil {
			return nil, errors.New("cli: manifest has no executable; reinstall the plugin")
		}
	}

	env := BuildEnv(EnvParams{
		Isolation:  manifest.Isolation,
		CLIDataDir: cliDataDir,
		RuntimeBin: filepath.Dir(m.nodePath),
		PluginBin:  filepath.Join(installDir, "node_modules", ".bin"),
		ParentEnv:  os.Environ(),
	})

	// Pin Cwd to installDir so login scripts that resolve files relative to themselves
	// (config templates, embedded node modules) work regardless of where the app was launched from.
	// RunTool intentionally leaves Cwd unset because tool argv is fully resolved by the manifest.
	return StartLoginSession(ctx, LoginSessionParams{
		Executable: manifest.Executable,
		Argv:       argv,
		Env:        env,
		Cwd:        installDir,
	})
}

// StartLogin runs the first resolved login step. Retained for callers that have
// always treated login as a single command; multi-step flows should call
// StartLoginCommand once per step from manifest.ResolveLoginSteps().
// Returns errors.New("cli: manifest has no login_command") when no steps are configured.
func (m *Manager) StartLogin(ctx context.Context, name string) (*LoginSession, error) {
	cliDataDir := filepath.Join(dir.CLIDataRoot(m.dataDir), name)
	manifestPath := filepath.Join(cliDataDir, "manifest.json")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	steps := manifest.ResolveLoginSteps()
	if len(steps) == 0 {
		return nil, errors.New("cli: manifest has no login_command")
	}
	return m.StartLoginCommand(ctx, name, steps[0])
}
