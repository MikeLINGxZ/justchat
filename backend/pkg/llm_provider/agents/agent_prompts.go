package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

// AgentPromptDir 返回某个 Agent 的提示词目录：{dataPath}/agents/{agentName}/
func AgentPromptDir(agentName string) (string, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(dataPath, "agents", agentName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// LoadAgentPrompt 从 Agent 目录加载提示词文件。
// 如果文件不存在或为空，写入 defaultContent 并返回。
func LoadAgentPrompt(agentName, promptFileName, defaultContent string) (string, error) {
	dir, err := AgentPromptDir(agentName)
	if err != nil {
		return defaultContent, err
	}
	path := filepath.Join(dir, promptFileName)

	content, readErr := os.ReadFile(path)
	if readErr == nil {
		text := strings.TrimSpace(string(content))
		if text != "" {
			return text, nil
		}
	}

	// 文件不存在或为空：写入默认内容
	if writeErr := writePromptFile(path, defaultContent); writeErr != nil {
		return defaultContent, writeErr
	}
	return defaultContent, nil
}

// SaveAgentPrompt 保存 Agent 提示词到文件。
func SaveAgentPrompt(agentName, promptFileName, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("prompt content cannot be empty")
	}
	dir, err := AgentPromptDir(agentName)
	if err != nil {
		return err
	}
	return writePromptFile(filepath.Join(dir, promptFileName), content)
}

// DefaultAgentPromptContent 从已注册的 Agent 定义中获取某个提示词的默认内容。
func DefaultAgentPromptContent(agentName, promptFileName string) (string, bool) {
	agentDef, ok := FindAgent(agentName)
	if !ok {
		return "", false
	}
	defaults := agentDef.DefaultPrompts()
	content, found := defaults[promptFileName]
	return content, found
}

func writePromptFile(path string, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content+"\n"), 0o644)
}
