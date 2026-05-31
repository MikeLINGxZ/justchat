//go:build !windows

package tools

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
)

const defaultShellPTYWaitDelay = time.Second

func shellFuncPTY(ctx context.Context, emitter ShellProgressEmitter, sessionID uint, input shellInput) (shellOutput, error) {
	shell, args := resolveShellCommand(input.Command)
	cmd := exec.CommandContext(ctx, shell, args...)
	cmd.Dir = resolveShellWorkDir(input.WorkDir)
	cmd.WaitDelay = 2 * defaultShellPTYWaitDelay

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 40, Cols: 120})
	if err != nil {
		return shellOutput{}, fmt.Errorf("start pty command: %w", err)
	}
	defer ptmx.Close()

	var stdout strings.Builder
	buf := make([]byte, 4096)
	for {
		n, readErr := ptmx.Read(buf)
		if n > 0 {
			text := string(buf[:n])
			stdout.WriteString(text)
			if emitter != nil && sessionID != 0 {
				emitter.EmitToolResult(sessionID, "shell", stdout.String())
			}
		}
		if readErr != nil {
			break
		}
	}

	waitErr := cmd.Wait()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	out := shellOutput{
		Stdout:   stdout.String(),
		Stderr:   "",
		ExitCode: exitCode,
	}
	if ctx.Err() != nil {
		return out, fmt.Errorf("exec: %w", ctx.Err())
	}
	if waitErr != nil {
		var exitErr *exec.ExitError
		if !errors.As(waitErr, &exitErr) {
			return out, fmt.Errorf("exec: %w", waitErr)
		}
	}
	return out, nil
}
