package terminal

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"
)

// ServiceStartup wires terminal manager events into Wails runtime events.
func (t *Terminal) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	_ = ctx
	_ = options
	t.wailsApp = application.Get()
	t.manager.SetOutputEmitter(func(event pkgterminal.OutputEvent) {
		if t.wailsApp != nil {
			t.wailsApp.Event.Emit(eventTerminalOutput, event)
		}
	})
	t.manager.SetStatusEmitter(func(event pkgterminal.StatusEvent) {
		if t.wailsApp != nil {
			t.wailsApp.Event.Emit(eventTerminalExited, event)
		}
	})
	return nil
}
