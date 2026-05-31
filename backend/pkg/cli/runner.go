package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// DefaultStdoutLimit is the per-stream byte cap when RunParams.StdoutLimit/StderrLimit is zero.
const DefaultStdoutLimit = 32 * 1024

// DefaultTimeoutSeconds is the per-call wall-clock cap when RunParams.TimeoutSec is zero.
const DefaultTimeoutSeconds = 60

// MaxTimeoutSeconds is the hard ceiling regardless of manifest overrides.
const MaxTimeoutSeconds = 600

// RunParams describes one CLI subprocess invocation.
type RunParams struct {
	Executable  string
	Argv        []string
	Env         []string
	WorkingDir  string
	OutputMode  OutputMode
	TimeoutSec  int
	StdoutLimit int
	StderrLimit int
	UsePTY      bool
}

// RunResult captures the outcome of a CLI invocation in a format suitable for forwarding to the LLM.
type RunResult struct {
	ExitCode        int
	Stdout          string
	Stderr          string
	Parsed          json.RawMessage // populated when OutputMode==OutputJSON and stdout is valid JSON
	ParsedLines     []string        // populated when OutputMode==OutputLines
	TruncatedStdout bool
	TruncatedStderr bool
	DurationMS      int64
}

// RunProgress is a cumulative snapshot of streaming subprocess output.
type RunProgress struct {
	Stdout          string
	Stderr          string
	TruncatedStdout bool
	TruncatedStderr bool
}

// ErrRunTimeout is returned by Run when the wall-clock timeout fires before the process exits.
var ErrRunTimeout = errors.New("cli: run timed out")

// Run executes one CLI subprocess applying timeout, output truncation, and OutputMode parsing.
// A non-zero exit code is reported via RunResult.ExitCode, not as a Go error.
func Run(ctx context.Context, p RunParams) (RunResult, error) {
	return RunWithProgress(ctx, p, nil)
}

// RunWithProgress executes one CLI subprocess and emits cumulative output snapshots as stdout/stderr arrive.
func RunWithProgress(ctx context.Context, p RunParams, onProgress func(RunProgress)) (RunResult, error) {
	if p.UsePTY {
		return runWithProgressPTY(ctx, p, onProgress)
	}
	return runWithProgressPipe(ctx, p, onProgress)
}

func runWithProgressPipe(ctx context.Context, p RunParams, onProgress func(RunProgress)) (RunResult, error) {
	stdoutLimit := p.StdoutLimit
	if stdoutLimit <= 0 {
		stdoutLimit = DefaultStdoutLimit
	}
	stderrLimit := p.StderrLimit
	if stderrLimit <= 0 {
		stderrLimit = DefaultStdoutLimit
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
	// Put the child in its own process group and, on cancel, signal the whole group so
	// shell children (e.g. `bash -c "sleep 5"`) do not outlive the immediate child.
	applyProcessGroup(cmd)
	// WaitDelay is a safety net: even if a grandchild keeps the stdout pipe open after
	// the cancel signal, Wait will return at most this long after Cancel ran.
	cmd.WaitDelay = 2 * time.Second

	stdoutBuf := &cappedBuffer{limit: stdoutLimit}
	stderrBuf := &cappedBuffer{limit: stderrLimit}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return RunResult{}, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return RunResult{}, err
	}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		return RunResult{}, err
	}

	chunks := make(chan cliOutputChunk, 32)
	var wg sync.WaitGroup
	wg.Add(2)
	go readCLIOutputPipe(&wg, stdoutPipe, stdoutBuf, true, chunks)
	go readCLIOutputPipe(&wg, stderrPipe, stderrBuf, false, chunks)
	go func() {
		wg.Wait()
		close(chunks)
	}()

	for range chunks {
		if onProgress != nil {
			onProgress(RunProgress{
				Stdout:          stdoutBuf.String(),
				Stderr:          stderrBuf.String(),
				TruncatedStdout: stdoutBuf.truncated,
				TruncatedStderr: stderrBuf.truncated,
			})
		}
	}

	err = cmd.Wait()
	res := RunResult{
		Stdout:          stdoutBuf.String(),
		Stderr:          stderrBuf.String(),
		TruncatedStdout: stdoutBuf.truncated,
		TruncatedStderr: stderrBuf.truncated,
		DurationMS:      time.Since(start).Milliseconds(),
		ExitCode:        cmd.ProcessState.ExitCode(),
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

type cliOutputChunk struct{}

func readCLIOutputPipe(wg *sync.WaitGroup, r io.ReadCloser, buf *cappedBuffer, emit bool, chunks chan<- cliOutputChunk) {
	defer wg.Done()
	defer r.Close()

	tmp := make([]byte, 2048)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			_, _ = buf.Write(tmp[:n])
			if emit {
				chunks <- cliOutputChunk{}
			} else {
				chunks <- cliOutputChunk{}
			}
		}
		if err != nil {
			return
		}
	}
}

// parseOutputs populates Parsed / ParsedLines based on OutputMode and stdout content.
func parseOutputs(res *RunResult, mode OutputMode) {
	switch mode {
	case OutputJSON:
		trimmed := strings.TrimSpace(res.Stdout)
		if trimmed == "" {
			return
		}
		var probe interface{}
		if err := json.Unmarshal([]byte(trimmed), &probe); err == nil {
			res.Parsed = json.RawMessage(trimmed)
		}
	case OutputLines:
		if res.Stdout == "" {
			res.ParsedLines = []string{}
			return
		}
		res.ParsedLines = strings.Split(strings.TrimRight(res.Stdout, "\n"), "\n")
	}
}

// cappedBuffer is an io.Writer that stops appending past limit and remembers it was capped.
type cappedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

// Write appends up to limit bytes; further writes are silently dropped and truncated is set.
func (c *cappedBuffer) Write(b []byte) (int, error) {
	if c.buf.Len() >= c.limit {
		c.truncated = true
		return len(b), nil
	}
	remaining := c.limit - c.buf.Len()
	if len(b) > remaining {
		c.buf.Write(b[:remaining])
		c.truncated = true
		return len(b), nil
	}
	return c.buf.Write(b)
}

// String returns the accumulated bytes as a string.
func (c *cappedBuffer) String() string { return c.buf.String() }
