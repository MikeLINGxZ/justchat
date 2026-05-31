package window

import (
	"context"
	"fmt"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/window_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/window_options"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window/window_dto"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Window management service
type Window struct {
	wailsApp *application.App
}

// OpenSettings open the settings window
func (p *Window) OpenSettings(ctx context.Context, input window_dto.OpenSettingsInput) (*window_dto.OpenSettingsOutput, error) {
	tab := "settings"
	if input.Tab != "" {
		tab = input.Tab
	}
	pageUrl := fmt.Sprintf("/?entry=settings&tab=%s", tab)

	settingsWindow, ok := p.wailsApp.Window.GetByName(window_id.Settings)
	if ok {
		settingsWindow.SetURL(pageUrl)
		settingsWindow.Focus()
		p.centerWindowOnHomeScreen(settingsWindow)
		settingsWindow.Show()
		return &window_dto.OpenSettingsOutput{}, nil
	}
	settingsWindow = p.wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  window_id.Settings,
		Title: i18n.TCurrent("app.window.settings_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 48,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              pageUrl,
		Width:            1200,
		Height:           800,
		MinWidth:         550,
		MinHeight:        550,
	})
	settingsWindow.Focus()
	p.centerWindowOnHomeScreen(settingsWindow)
	settingsWindow.Show()
	return &window_dto.OpenSettingsOutput{}, nil
}

// OpenAddProvider open the add provider window
func (p *Window) OpenAddProvider(ctx context.Context, input window_dto.OpenAddProviderInput) (*window_dto.OpenAddProviderOutput, error) {

	addProviderWindow, ok := p.wailsApp.Window.GetByName(window_id.AddProvider)
	if ok {
		addProviderWindow.Focus()
		p.centerWindowOnHomeScreen(addProviderWindow)
		addProviderWindow.Show()
		return &window_dto.OpenAddProviderOutput{}, nil
	}
	addProviderWindow = p.wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  window_id.AddProvider,
		Title: i18n.TCurrent("app.window.add_provider_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 48,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/?entry=add_provider",
		Width:            750,
		Height:           800,
		MinWidth:         550,
		MinHeight:        550,
	})
	addProviderWindow.Focus()
	p.centerWindowOnHomeScreen(addProviderWindow)
	addProviderWindow.Show()
	return &window_dto.OpenAddProviderOutput{}, nil
}

// CloseAddProvider close the add provider window
func (p *Window) CloseAddProvider(ctx context.Context, input window_dto.CloseAddProviderInput) (*window_dto.CloseAddProviderOutput, error) {
	addProviderWindow, ok := p.wailsApp.Window.GetByName(window_id.AddProvider)
	if ok {
		addProviderWindow.Close()
		return &window_dto.CloseAddProviderOutput{}, nil
	}
	return &window_dto.CloseAddProviderOutput{}, nil
}

// OpenAddSkill opens or focuses the add skill window.
func (p *Window) OpenAddSkill(ctx context.Context, input window_dto.OpenAddSkillInput) (*window_dto.OpenAddSkillOutput, error) {
	addSkillWindow, ok := p.wailsApp.Window.GetByName(window_id.AddSkill)
	if ok {
		addSkillWindow.Focus()
		p.centerWindowOnHomeScreen(addSkillWindow)
		addSkillWindow.Show()
		return &window_dto.OpenAddSkillOutput{}, nil
	}

	addSkillWindow = p.wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  window_id.AddSkill,
		Title: i18n.TCurrent("app.window.add_skill_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 48,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/?entry=add_skill",
		Width:            900,
		Height:           820,
		MinWidth:         700,
		MinHeight:        620,
	})
	addSkillWindow.Focus()
	p.centerWindowOnHomeScreen(addSkillWindow)
	addSkillWindow.Show()
	return &window_dto.OpenAddSkillOutput{}, nil
}

// CloseAddSkill closes the add skill window.
func (p *Window) CloseAddSkill(ctx context.Context, input window_dto.CloseAddSkillInput) (*window_dto.CloseAddSkillOutput, error) {
	addSkillWindow, ok := p.wailsApp.Window.GetByName(window_id.AddSkill)
	if ok {
		addSkillWindow.Close()
		return &window_dto.CloseAddSkillOutput{}, nil
	}
	return &window_dto.CloseAddSkillOutput{}, nil
}

// OpenAddMemory opens or focuses the add memory window.
func (p *Window) OpenAddMemory(ctx context.Context, input window_dto.OpenAddMemoryInput) (*window_dto.OpenAddMemoryOutput, error) {
	addMemoryWindow, ok := p.wailsApp.Window.GetByName(window_id.AddMemory)
	if ok {
		addMemoryWindow.Focus()
		p.centerWindowOnHomeScreen(addMemoryWindow)
		addMemoryWindow.Show()
		return &window_dto.OpenAddMemoryOutput{}, nil
	}

	addMemoryWindow = p.wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  window_id.AddMemory,
		Title: i18n.TCurrent("app.window.add_memory_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 48,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/?entry=add_memory",
		Width:            900,
		Height:           820,
		MinWidth:         550,
		MinHeight:        620,
	})
	addMemoryWindow.Focus()
	p.centerWindowOnHomeScreen(addMemoryWindow)
	addMemoryWindow.Show()
	return &window_dto.OpenAddMemoryOutput{}, nil
}

// CloseAddMemory closes the add memory window.
func (p *Window) CloseAddMemory(ctx context.Context, input window_dto.CloseAddMemoryInput) (*window_dto.CloseAddMemoryOutput, error) {
	addMemoryWindow, ok := p.wailsApp.Window.GetByName(window_id.AddMemory)
	if ok {
		addMemoryWindow.Close()
		return &window_dto.CloseAddMemoryOutput{}, nil
	}
	return &window_dto.CloseAddMemoryOutput{}, nil
}

// OpenOnboarding opens or focuses the first-launch onboarding window.
func (p *Window) OpenOnboarding(ctx context.Context, input window_dto.OpenOnboardingInput) (*window_dto.OpenOnboardingOutput, error) {
	if onboardingWindow, ok := p.wailsApp.Window.GetByName(window_id.Onboarding); ok {
		onboardingWindow.Focus()
		p.centerWindowOnHomeScreen(onboardingWindow)
		onboardingWindow.Show()
		return &window_dto.OpenOnboardingOutput{}, nil
	}

	onboardingWindow := p.wailsApp.Window.NewWithOptions(window_options.DefaultOnboarding())
	onboardingWindow.Focus()
	p.centerWindowOnHomeScreen(onboardingWindow)
	onboardingWindow.Show()
	return &window_dto.OpenOnboardingOutput{}, nil
}
