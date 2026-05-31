//go:build windows

package tools

import "context"

func shellFuncPTY(ctx context.Context, emitter ShellProgressEmitter, sessionID uint, input shellInput) (shellOutput, error) {
	return shellFuncPipe(ctx, emitter, sessionID, input)
}
