package onboarding

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/window_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/window_options"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/onboarding/onboarding_dto"
	providerSvc "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto"
)

// Onboarding manages the first-launch initialization flow.
type Onboarding struct {
	wailsApp    *application.App
	providerSvc *providerSvc.Provider
}

// NewOnboarding constructs a new Onboarding service bound to the provider service.
func NewOnboarding(p *providerSvc.Provider) *Onboarding {
	return &Onboarding{providerSvc: p}
}

// IsInitialized reports whether the application has completed first-launch onboarding.
func (s *Onboarding) IsInitialized(ctx context.Context, input onboarding_dto.IsInitializedInput) (*onboarding_dto.IsInitializedOutput, error) {
	initialized, err := isInitialized()
	if err != nil {
		return nil, ierror.Error(ierror.ErrOnboardingReadInit, err)
	}
	return &onboarding_dto.IsInitializedOutput{Initialized: initialized}, nil
}

// SaveProviderAndMarkInitialized creates the first provider and writes the init.json marker.
// Returns an error if the provider creation fails; the init marker is only written on success.
func (s *Onboarding) SaveProviderAndMarkInitialized(ctx context.Context, input onboarding_dto.SaveProviderAndMarkInitializedInput) (*onboarding_dto.SaveProviderAndMarkInitializedOutput, error) {
	if _, err := s.providerSvc.CreateProvider(ctx, provider_dto.CreateProviderInput{
		ProviderName: input.ProviderName,
		ProviderType: input.ProviderType,
		BaseUrl:      input.BaseUrl,
		ApiKey:       input.ApiKey,
		Enable:       input.Enable,
		DefaultModel: input.DefaultModel,
		Models:       input.Models,
	}); err != nil {
		return nil, err
	}

	if err := markInitialized(); err != nil {
		return nil, ierror.Error(ierror.ErrOnboardingWriteInit, err)
	}

	return &onboarding_dto.SaveProviderAndMarkInitializedOutput{}, nil
}

// ExitApp terminates the Wails application; the OS-level process exits.
func (s *Onboarding) ExitApp(ctx context.Context, input onboarding_dto.ExitAppInput) (*onboarding_dto.ExitAppOutput, error) {
	if s.wailsApp != nil {
		s.wailsApp.Quit()
	}
	return &onboarding_dto.ExitAppOutput{}, nil
}

// EnterHome opens the main window and closes the onboarding window.
func (s *Onboarding) EnterHome(ctx context.Context, input onboarding_dto.EnterHomeInput) (*onboarding_dto.EnterHomeOutput, error) {
	if s.wailsApp == nil {
		return &onboarding_dto.EnterHomeOutput{}, nil
	}

	homeWindow, ok := s.wailsApp.Window.GetByName(window_id.Home)
	if !ok {
		homeWindow = s.wailsApp.Window.NewWithOptions(window_options.DefaultHome())
	}
	homeWindow.Show()
	homeWindow.Focus()

	if onboardWin, ok := s.wailsApp.Window.GetByName(window_id.Onboarding); ok {
		onboardWin.Close()
	}

	return &onboarding_dto.EnterHomeOutput{}, nil
}
