package utils

import (
	"os"
	"path/filepath"
)

func GetDataPath() (string, error) {
	if os.Getenv("LEMONTEA_DATA_PATH") != "" {
		if err := os.MkdirAll(os.Getenv("LEMONTEA_DATA_PATH"), 0755); err != nil {
			return "", err
		}
		return os.Getenv("LEMONTEA_DB_PATH"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataPath := filepath.Join(homeDir, ".lemon_tea")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return "", err
	}
	return dataPath, nil
}
