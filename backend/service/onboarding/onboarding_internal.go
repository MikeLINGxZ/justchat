package onboarding

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

const initFileName = "init.json"

type initState struct {
	Initialized bool      `json:"initialized"`
	CompletedAt time.Time `json:"completed_at"`
}

// initFilePath returns the absolute path to the init marker file under the data dir.
func initFilePath() (string, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, initFileName), nil
}

// isInitialized reports whether the init.json marker file exists.
func isInitialized() (bool, error) {
	path, err := initFilePath()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// markInitialized writes the init.json marker file containing the completion timestamp.
func markInitialized() error {
	path, err := initFilePath()
	if err != nil {
		return err
	}
	payload, err := json.MarshalIndent(initState{
		Initialized: true,
		CompletedAt: time.Now(),
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}
