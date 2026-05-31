package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// writeFakeBin writes an executable shell script under dir and returns its absolute path.
// On Windows the helper skips the test (we rely on POSIX shell for fake CLIs in unit tests).
func writeFakeBin(t *testing.T, dir, name, script string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("runner tests use POSIX shell scripts")
	}
	path := filepath.Join(dir, name)
	body := "#!/usr/bin/env bash\nset -e\n" + script + "\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write fake bin: %v", err)
	}
	return path
}

// TestRunReturnsStdoutAndExitCode verifies a successful run captures stdout and exit_code=0.
func TestRunReturnsStdoutAndExitCode(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "echoer", `echo "hello $1"`)
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		Argv:       []string{"world"},
		OutputMode: OutputText,
		TimeoutSec: 5,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit code: %d", res.ExitCode)
	}
	if strings.TrimSpace(res.Stdout) != "hello world" {
		t.Fatalf("stdout=%q", res.Stdout)
	}
}

// TestRunHonorsTimeout verifies that exceeding TimeoutSec kills the process and reports a timeout result.
func TestRunHonorsTimeout(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "sleeper", `sleep 5`)
	start := time.Now()
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputText,
		TimeoutSec: 1,
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("expected timeout error, got nil; res=%+v", res)
	}
	if elapsed >= 4*time.Second {
		t.Fatalf("timeout did not fire promptly: %v", elapsed)
	}
}

// TestRunTruncatesOversizedStdout verifies stdout is capped at StdoutLimit bytes with TruncatedStdout=true.
func TestRunTruncatesOversizedStdout(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "spammer", `head -c 200 /dev/zero | tr '\0' 'x'`)
	res, err := Run(context.Background(), RunParams{
		Executable:  bin,
		OutputMode:  OutputText,
		TimeoutSec:  5,
		StdoutLimit: 64,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(res.Stdout) != 64 {
		t.Fatalf("stdout len=%d (want 64): %q", len(res.Stdout), res.Stdout)
	}
	if !res.TruncatedStdout {
		t.Fatalf("expected TruncatedStdout=true")
	}
}

// TestRunOutputModeJSONParsesValidJSON verifies OutputJSON populates Parsed on valid JSON stdout.
func TestRunOutputModeJSONParsesValidJSON(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "jsonbin", `printf '{"k":"v"}'`)
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputJSON,
		TimeoutSec: 5,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Parsed == nil {
		t.Fatalf("expected Parsed != nil for valid JSON output")
	}
	var doc map[string]string
	if err := json.Unmarshal(res.Parsed, &doc); err != nil {
		t.Fatalf("Parsed not valid JSON: %v", err)
	}
	if doc["k"] != "v" {
		t.Fatalf("Parsed=%s", string(res.Parsed))
	}
}

// TestRunOutputModeJSONFallsBackOnInvalidJSON verifies invalid JSON leaves Parsed nil and keeps Stdout.
func TestRunOutputModeJSONFallsBackOnInvalidJSON(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "txtbin", `echo "not json"`)
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputJSON,
		TimeoutSec: 5,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Parsed != nil {
		t.Fatalf("Parsed should be nil for invalid JSON")
	}
	if !strings.Contains(res.Stdout, "not json") {
		t.Fatalf("Stdout=%q", res.Stdout)
	}
}

// TestRunOutputModeLinesSplitsOnNewline verifies OutputLines populates ParsedLines.
func TestRunOutputModeLinesSplitsOnNewline(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "lines", `printf "a\nb\nc"`)
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputLines,
		TimeoutSec: 5,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(res.ParsedLines) != 3 || res.ParsedLines[0] != "a" || res.ParsedLines[2] != "c" {
		t.Fatalf("ParsedLines=%v", res.ParsedLines)
	}
}

// TestRunPropagatesNonZeroExit verifies a non-zero exit_code is returned without error.
func TestRunPropagatesNonZeroExit(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "fail", `echo "oops" 1>&2; exit 7`)
	res, err := Run(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputText,
		TimeoutSec: 5,
	})
	if err != nil {
		t.Fatalf("run: %v (a non-zero exit should not be a Go error)", err)
	}
	if res.ExitCode != 7 {
		t.Fatalf("exit_code=%d", res.ExitCode)
	}
	if !strings.Contains(res.Stderr, "oops") {
		t.Fatalf("stderr=%q", res.Stderr)
	}
}

func TestRunWithProgressStreamsCumulativeOutput(t *testing.T) {
	bin := writeFakeBin(t, t.TempDir(), "streamer", `printf "hello"; sleep 0.1; printf " world"; printf " warn" 1>&2`)
	var snapshots []RunProgress
	res, err := RunWithProgress(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputText,
		TimeoutSec: 5,
	}, func(progress RunProgress) {
		snapshots = append(snapshots, progress)
	})
	if err != nil {
		t.Fatalf("run with progress: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %d", res.ExitCode)
	}
	if len(snapshots) == 0 {
		t.Fatal("expected progress snapshots")
	}
	last := snapshots[len(snapshots)-1]
	if !strings.Contains(last.Stdout, "hello world") {
		t.Fatalf("expected cumulative stdout, got %+v", last)
	}
	if !strings.Contains(last.Stderr, " warn") {
		t.Fatalf("expected stderr snapshot, got %+v", last)
	}
}

func TestRunWithProgressPTYKeepsInteractiveTTYSessionOpen(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PTY runner is only implemented on POSIX platforms")
	}
	bin := writeFakeBin(t, t.TempDir(), "interactive-login", `
if [ ! -t 0 ]; then
  echo "QR"
  exit 0
fi
echo "QR"
sleep 0.1
echo "AUTH_OK"
`)
	var snapshots []RunProgress
	res, err := RunWithProgress(context.Background(), RunParams{
		Executable: bin,
		OutputMode: OutputText,
		TimeoutSec: 5,
		UsePTY:     true,
	}, func(progress RunProgress) {
		snapshots = append(snapshots, progress)
	})
	if err != nil {
		t.Fatalf("run with pty progress: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %d", res.ExitCode)
	}
	if !strings.Contains(res.Stdout, "QR") || !strings.Contains(res.Stdout, "AUTH_OK") {
		t.Fatalf("expected PTY stdout to keep session through auth, got %q", res.Stdout)
	}
	if len(snapshots) == 0 || !strings.Contains(snapshots[len(snapshots)-1].Stdout, "AUTH_OK") {
		t.Fatalf("expected cumulative PTY progress with auth completion, got %+v", snapshots)
	}
}
