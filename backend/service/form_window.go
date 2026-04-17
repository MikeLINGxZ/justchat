package service

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

func (s *Service) openFormWindow(name, title, url string, width, height int) {
	if s.app == nil {
		return
	}
	if existing, ok := s.app.Window.GetByName(name); ok {
		existing.SetURL(url)
		existing.Focus()
		existing.Center()
		existing.Show()
		return
	}
	window := s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  name,
		Title: title,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              url,
		Width:            width,
		Height:           height,
		MinWidth:         480,
		MinHeight:        520,
	})
	window.Focus()
	window.Center()
	window.Show()
}

func (s *Service) OpenAddProviderWindow() {
	s.openFormWindow(
		WindowNameFormProvider,
		i18n.TCurrent("app.window.form_provider_title", nil),
		"/?entry=form_provider",
		720, 780,
	)
}

func (s *Service) OpenAddAgentWindow() {
	s.openFormWindow(
		WindowNameFormAgent,
		i18n.TCurrent("app.window.form_agent_title", nil),
		"/?entry=form_agent",
		680, 720,
	)
}

func (s *Service) OpenAddSkillWindow() {
	s.openFormWindow(
		WindowNameFormSkill,
		i18n.TCurrent("app.window.form_skill_title", nil),
		"/?entry=form_skill",
		720, 760,
	)
}

// OpenEditMemoryWindow opens a per-ID singleton window for editing a memory record.
// Multiple memories can be edited in parallel; reopening the same id focuses the existing window.
func (s *Service) OpenEditMemoryWindow(id uint) {
	name := fmt.Sprintf("%s_%d", WindowNameFormMemory, id)
	url := fmt.Sprintf("/?entry=form_memory&id=%d", id)
	s.openFormWindow(
		name,
		i18n.TCurrent("app.window.form_memory_title", nil),
		url,
		680, 820,
	)
}

// CloseFormWindow closes a form window by name. Used by form pages to self-close after submit.
func (s *Service) CloseFormWindow(name string) {
	if s.app == nil || name == "" {
		return
	}
	if window, ok := s.app.Window.GetByName(name); ok {
		window.Close()
	}
}
