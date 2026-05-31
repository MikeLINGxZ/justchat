//go:build !windows

package cli

import (
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func runWithProgressPTY(ctx context.Context, p RunParams, onProgress func(RunProgress)) (RunResult, error) {
	stdoutLimit := p.StdoutLimit
	if stdoutLimit <= 0 {
		stdoutLimit = DefaultStdoutLimit
	}
	timeoutSec := p.TimeoutSec
	if timeoutSec <= 0 {
		timeoutSec = DefaultTimeoutSeconds
	}
	if timeoutSec > MaxTimeoutSeconds {
		timeoutSec = MaxTimeoutSeconds
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, p.Executable, p.Argv...)
	if p.Env != nil {
		cmd.Env = p.Env
	}
	if p.WorkingDir != "" {
		cmd.Dir = p.WorkingDir
	}
	cmd.WaitDelay = 2 * time.Second

	start := time.Now()
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 40, Cols: 120})
	if err != nil {
		return RunResult{}, err
	}
	defer ptmx.Close()

	stdoutBuf := &cappedBuffer{limit: stdoutLimit}
	tmp := make([]byte, 4096)
	for {
		n, readErr := ptmx.Read(tmp)
		if n > 0 {
			_, _ = stdoutBuf.Write(tmp[:n])
			if onProgress != nil {
				onProgress(RunProgress{
					Stdout:          stdoutBuf.String(),
					TruncatedStdout: stdoutBuf.truncated,
				})
			}
		}
		if readErr != nil {
			break
		}
	}

	err = cmd.Wait()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	res := RunResult{
		Stdout:          stdoutBuf.String(),
		TruncatedStdout: stdoutBuf.truncated,
		DurationMS:      time.Since(start).Milliseconds(),
		ExitCode:        exitCode,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return res, ErrRunTimeout
	}
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return res, err
		}
	}

	parseOutputs(&res, p.OutputMode)
	return res, nil
}
