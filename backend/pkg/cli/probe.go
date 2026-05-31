package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PackageMeta is the subset of package.json fields used by the CLI installer / probe.
type PackageMeta struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Bin         map[string]string `json:"-"` // normalized from string|object bin field
}

// ReadPackageJSON loads {installDir}/package.json (or .../node_modules/<name>/package.json fallback)
// and returns its key fields, normalizing the bin field into a name->script map.
func ReadPackageJSON(installDir string) (PackageMeta, error) {
	candidates := []string{
		filepath.Join(installDir, "package.json"),
	}
	// npm install --prefix puts the package under node_modules/<name>/.
	entries, _ := os.ReadDir(filepath.Join(installDir, "node_modules"))
	for _, e := range entries {
		if !e.IsDir() || e.Name() == ".bin" {
			continue
		}
		if strings.HasPrefix(e.Name(), "@") {
			scopedEntries, _ := os.ReadDir(filepath.Join(installDir, "node_modules", e.Name()))
			for _, scoped := range scopedEntries {
				if !scoped.IsDir() {
					continue
				}
				candidates = append(candidates, filepath.Join(installDir, "node_modules", e.Name(), scoped.Name(), "package.json"))
			}
			continue
		}
		candidates = append(candidates, filepath.Join(installDir, "node_modules", e.Name(), "package.json"))
	}

	var lastErr error
	for _, path := range candidates {
		meta, err := readOnePackageJSON(path)
		if err == nil {
			return meta, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no package.json found")
	}
	return PackageMeta{}, lastErr
}

// readOnePackageJSON parses one package.json file and normalizes the bin field.
func readOnePackageJSON(path string) (PackageMeta, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return PackageMeta{}, err
	}
	var raw struct {
		Name        string          `json:"name"`
		Version     string          `json:"version"`
		Description string          `json:"description"`
		Author      json.RawMessage `json:"author"`
		Bin         json.RawMessage `json:"bin"`
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return PackageMeta{}, err
	}
	meta := PackageMeta{
		Name:        raw.Name,
		Version:     raw.Version,
		Description: raw.Description,
		Bin:         map[string]string{},
	}
	if len(raw.Author) > 0 {
		var s string
		if err := json.Unmarshal(raw.Author, &s); err == nil {
			meta.Author = s
		} else {
			var obj struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw.Author, &obj); err == nil {
				meta.Author = obj.Name
			}
		}
	}
	if len(raw.Bin) > 0 {
		var s string
		if err := json.Unmarshal(raw.Bin, &s); err == nil {
			meta.Bin[meta.Name] = s
		} else {
			var obj map[string]string
			if err := json.Unmarshal(raw.Bin, &obj); err == nil {
				meta.Bin = obj
			}
		}
	}
	return meta, nil
}

// SelectExecutable picks the most appropriate bin entry for the package and returns the absolute path
// to its symlink under {installDir}/node_modules/.bin/{binName}.
// Selection rules: prefer a bin entry whose key matches the package name; otherwise the only entry; else error.
func SelectExecutable(installDir string, pkg PackageMeta) (string, error) {
	if len(pkg.Bin) == 0 {
		return discoverExecutableFromBinDir(installDir, pkg)
	}
	if _, ok := pkg.Bin[pkg.Name]; ok {
		return filepath.Join(installDir, "node_modules", ".bin", pkg.Name), nil
	}
	if len(pkg.Bin) == 1 {
		for binName := range pkg.Bin {
			return filepath.Join(installDir, "node_modules", ".bin", binName), nil
		}
	}
	return discoverExecutableFromBinDir(installDir, pkg)
}

func discoverExecutableFromBinDir(installDir string, pkg PackageMeta) (string, error) {
	binDir := filepath.Join(installDir, "node_modules", ".bin")
	entries, err := os.ReadDir(binDir)
	if err != nil {
		if len(pkg.Bin) == 0 {
			return "", errors.New("package has no bin entry; not a CLI")
		}
		return "", errors.New("package has multiple bin entries and none matches the package name")
	}

	candidates := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		candidates = append(candidates, name)
	}
	sort.Strings(candidates)
	if len(candidates) == 1 {
		return filepath.Join(binDir, candidates[0]), nil
	}

	leaf := packageLeafName(pkg.Name)
	for _, name := range candidates {
		if name == pkg.Name || name == leaf {
			return filepath.Join(binDir, name), nil
		}
	}
	if len(candidates) == 0 && len(pkg.Bin) == 0 {
		return "", errors.New("package has no bin entry; not a CLI")
	}
	if len(candidates) == 0 {
		return "", errors.New("package has multiple bin entries and none matches the package name")
	}
	return "", errors.New("unable to infer executable from node_modules/.bin")
}

func packageLeafName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		return name[idx+1:]
	}
	return name
}

// HelpVariants is the ordered list of help-fetch attempts ProbeHelp will try.
var HelpVariants = [][]string{
	{"--help"},
	{"-h"},
	{"help"},
}

// ProbeHelp runs <executable> with each HelpVariants entry in order and returns the first non-empty stdout.
// A 10-second timeout is applied to each attempt.
func ProbeHelp(ctx context.Context, executable string, env []string) (string, error) {
	var lastErr error
	for _, args := range HelpVariants {
		res, err := Run(ctx, RunParams{
			Executable: executable,
			Argv:       args,
			Env:        env,
			OutputMode: OutputText,
			TimeoutSec: 10,
		})
		if err != nil {
			lastErr = err
			continue
		}
		if res.ExitCode == 0 && res.Stdout != "" {
			return res.Stdout, nil
		}
		lastErr = errors.New("help command exited non-zero")
	}
	if lastErr == nil {
		lastErr = errors.New("all help variants failed")
	}
	return "", lastErr
}
