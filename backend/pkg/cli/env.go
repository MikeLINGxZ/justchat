package cli

import (
	"path/filepath"
	"runtime"
	"strings"
)

// EnvParams gathers the inputs needed to assemble an isolated or shared environment.
type EnvParams struct {
	Isolation  IsolationMode
	CLIDataDir string   // absolute path to {data_dir}/plugins/cli_data/{name}
	RuntimeBin string   // absolute path to bundled runtime bin dir; npm/npx must resolve here first
	PluginBin  string   // absolute path to {install_dir}/node_modules/.bin
	ParentEnv  []string // typically os.Environ()
}

// strippedEnvKeys are removed from the inherited environment to prevent the host's npm config
// or NODE_PATH from leaking into the isolated subprocess.
var strippedEnvKeys = []string{"NPM_CONFIG_PREFIX", "NODE_PATH"}

// overriddenKeysIsolated are env vars that BuildEnv unconditionally replaces in isolated mode.
var overriddenKeysIsolated = []string{
	"HOME",
	"USERPROFILE",
	"XDG_CONFIG_HOME",
	"XDG_DATA_HOME",
	"XDG_CACHE_HOME",
}

// BuildEnv returns an exec.Cmd-style environment slice for a CLI subprocess.
// In isolated mode, HOME / USERPROFILE / XDG_* are pointed at subdirs of CLIDataDir;
// in shared mode the parent's values for those keys are preserved.
// In both modes RuntimeBin and PluginBin are prepended to PATH and strippedEnvKeys are removed.
func BuildEnv(p EnvParams) []string {
	out := make([]string, 0, len(p.ParentEnv)+8)

	skipKeys := map[string]struct{}{}
	for _, k := range strippedEnvKeys {
		skipKeys[k] = struct{}{}
	}
	if p.Isolation == IsolationIsolated {
		for _, k := range overriddenKeysIsolated {
			skipKeys[k] = struct{}{}
		}
	}

	for _, entry := range p.ParentEnv {
		idx := strings.IndexByte(entry, '=')
		if idx <= 0 {
			continue
		}
		key := entry[:idx]
		if key == "PATH" {
			out = append(out, "PATH="+prependPath(pathPrefix(p.RuntimeBin, p.PluginBin), entry[idx+1:]))
			continue
		}
		if _, skip := skipKeys[key]; skip {
			continue
		}
		out = append(out, entry)
	}

	if hasNoPath(p.ParentEnv) {
		out = append(out, "PATH="+pathPrefix(p.RuntimeBin, p.PluginBin))
	}

	if p.Isolation == IsolationIsolated {
		out = append(out,
			"HOME="+filepath.Join(p.CLIDataDir, "home"),
			"USERPROFILE="+filepath.Join(p.CLIDataDir, "home"),
			"XDG_CONFIG_HOME="+filepath.Join(p.CLIDataDir, "config"),
			"XDG_DATA_HOME="+filepath.Join(p.CLIDataDir, "data"),
			"XDG_CACHE_HOME="+filepath.Join(p.CLIDataDir, "cache"),
		)
	}
	return out
}

func pathPrefix(parts ...string) string {
	nonEmpty := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, pathListSeparator())
}

// prependPath returns "<bin>{sep}<existing>", skipping empty bin.
func prependPath(bin string, existing string) string {
	if bin == "" {
		return existing
	}
	return bin + pathListSeparator() + existing
}

func pathListSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	}
	return ":"
}

// hasNoPath reports whether the parent env lacks a PATH entry.
func hasNoPath(env []string) bool {
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			return false
		}
	}
	return true
}
