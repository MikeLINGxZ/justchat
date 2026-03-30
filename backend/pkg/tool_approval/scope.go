package tool_approval

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WorkspaceRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func ResolvePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path is empty")
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Join(WorkspaceRoot(), path), nil
}

func ResolveWorkingDirectory(path string) string {
	if strings.TrimSpace(path) == "" {
		return WorkspaceRoot()
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(WorkspaceRoot(), path)
}

func DescribeScope(target string) string {
	root := filepath.Clean(WorkspaceRoot())
	target = filepath.Clean(target)
	rel, err := filepath.Rel(root, target)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "工作区内"
	}
	return "工作区外"
}
