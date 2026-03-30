package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

const initFileName = "init.json"

type initState struct {
	Initialized bool      `json:"initialized"`
	CompletedAt time.Time `json:"completed_at"`
}

func getInitFilePath() (string, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataPath, initFileName), nil
}

func IsAppInitialized() (bool, error) {
	initFilePath, err := getInitFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(initFilePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func markAppInitialized() error {
	initFilePath, err := getInitFilePath()
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

	return os.WriteFile(initFilePath, payload, 0o644)
}

func (s *Service) IsInitialized() (bool, error) {
	initialized, err := IsAppInitialized()
	if err != nil {
		return false, ierror.NewError(err)
	}
	return initialized, nil
}

func (s *Service) CompleteOnboarding(provider view_models.Provider) error {
	if err := s.AddProvider(context.Background(), provider); err != nil {
		return err
	}

	if err := markAppInitialized(); err != nil {
		return ierror.NewError(err)
	}

	s.showHomeWindow()

	if s.app != nil {
		if onboardingWindow, ok := s.app.Window.GetByName(WindowNameOnboarding); ok {
			onboardingWindow.Close()
		}
	}

	return nil
}

func (s *Service) ExitApp() {
	if s.app != nil {
		s.app.Quit()
	}
}
