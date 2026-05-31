// backend/pkg/prompt/prompt.go
package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

var registry = make(map[string]string)

func Register(promptId string, defaultContent string) {
	registry[promptId] = defaultContent
}

func Load(promptId string) (string, error) {
	defaultContent, ok := registry[promptId]
	if !ok {
		return "", fmt.Errorf("prompt %q is not registered", promptId)
	}

	dataDir, err := dir.GetDataDir()
	if err != nil {
		return defaultContent, nil
	}

	customPath := filepath.Join(dataDir, "prompt", promptId, "index.md")
	content, err := os.ReadFile(customPath)
	if err != nil {
		return defaultContent, nil
	}

	return string(content), nil
}
