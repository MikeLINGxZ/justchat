package cli

import (
	"context"
	"errors"
	"path/filepath"
)

// NpmInstallParams describes a single npm install invocation.
type NpmInstallParams struct {
	NodePath  string   // absolute path to the bundled node binary; reserved for future direct invocation
	NpmPath   string   // absolute path to the bundled npm executable
	Source    string   // npm package name (e.g. "lark-cli", "@scope/pkg") or absolute local dir
	TargetDir string   // {data_dir}/plugins/cli/{name}; --prefix lands here
	ParentEnv []string // typically os.Environ(); BuildEnv-style scrubbing is applied
}

// NpmInstall runs `npm install <source> --prefix <targetDir> --no-audit --no-fund` using the bundled npm.
// The subprocess inherits ParentEnv with NPM_CONFIG_PREFIX / NODE_PATH stripped, so user-level npm config
// cannot redirect the install.
func NpmInstall(ctx context.Context, p NpmInstallParams) error {
	if p.NodePath == "" {
		return errors.New("npm: NodePath required")
	}
	if p.NpmPath == "" {
		return errors.New("npm: NpmPath required")
	}
	if p.Source == "" {
		return errors.New("npm: Source required")
	}
	if p.TargetDir == "" {
		return errors.New("npm: TargetDir required")
	}

	env := BuildEnv(EnvParams{
		Isolation:  IsolationShared, // npm itself needs the user's network / cert env; only key scrubbing applies
		CLIDataDir: "",
		RuntimeBin: filepath.Dir(p.NpmPath),
		PluginBin:  filepath.Join(p.TargetDir, "node_modules", ".bin"),
		ParentEnv:  p.ParentEnv,
	})

	args := []string{
		"install",
		p.Source,
		"--prefix", p.TargetDir,
		"--no-audit",
		"--no-fund",
	}

	res, err := Run(ctx, RunParams{
		Executable: p.NodePath,
		Argv:       append([]string{p.NpmPath}, args...),
		Env:        env,
		OutputMode: OutputText,
		TimeoutSec: MaxTimeoutSeconds, // npm installs of larger packages can take minutes
	})
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return errors.New("npm install failed: " + res.Stderr)
	}
	return nil
}
