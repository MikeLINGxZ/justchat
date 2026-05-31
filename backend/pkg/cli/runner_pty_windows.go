//go:build windows

package cli

import "context"

func runWithProgressPTY(ctx context.Context, p RunParams, onProgress func(RunProgress)) (RunResult, error) {
	p.UsePTY = false
	return runWithProgressPipe(ctx, p, onProgress)
}
