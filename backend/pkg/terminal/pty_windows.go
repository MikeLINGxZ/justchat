//go:build windows

package terminal

import (
	"context"
	"errors"
)

// startPTYProcess reports that PTY sessions are not yet available on Windows.
func startPTYProcess(ctx context.Context, params CreateParams) (ptyProcess, int, error) {
	_ = ctx
	_ = params
	return nil, 0, errors.New("terminal: PTY sessions are not implemented on Windows yet")
}
