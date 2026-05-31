//go:build !windows

package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// makeBashBin writes a bash script at a temp path, makes it executable, and returns the path.
func makeBashBin(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "bin.sh")
	if err := os.WriteFile(p, []byte("#!/usr/bin/env bash\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

// TestLoginSessionEchoesStdin starts a bash script that reads a line and echoes it back,
// writes "hello\n" to the session, and verifies the output contains "got: hello".
func TestLoginSessionEchoesStdin(t *testing.T) {
	bin := makeBashBin(t, "read line\necho \"got: $line\"\nexit 0\n")
	p := LoginSessionParams{
		Executable: "/usr/bin/env",
		Argv:       []string{"bash", bin},
	}
	s, err := StartLoginSession(context.Background(), p)
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := s.Write([]byte("hello\n")); err != nil {
		t.Fatalf("write: %v", err)
	}

	timeout := time.After(5 * time.Second)
	var collected []byte
	for {
		select {
		case chunk, ok := <-s.Output():
			if !ok {
				goto done
			}
			collected = append(collected, chunk...)
			if bytes.Contains(collected, []byte("got: hello")) {
				goto done
			}
		case <-timeout:
			t.Fatalf("timed out waiting for echo output; got %q", collected)
		}
	}
done:
	code, err := s.Wait()
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !bytes.Contains(collected, []byte("got: hello")) {
		t.Fatalf("expected 'got: hello' in output; got %q", collected)
	}
}

// TestLoginSessionCancelTerminatesChild starts a long-running bash sleep and verifies
// Cancel causes Wait to return within 2 seconds with a nonzero exit code.
func TestLoginSessionCancelTerminatesChild(t *testing.T) {
	bin := makeBashBin(t, "sleep 30\n")
	p := LoginSessionParams{
		Executable: "/usr/bin/env",
		Argv:       []string{"bash", bin},
	}
	s, err := StartLoginSession(context.Background(), p)
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := s.Cancel(); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	done := make(chan struct{})
	var exitCode int
	go func() {
		exitCode, _ = s.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Wait did not return within 5s after Cancel")
	}
	if exitCode == 0 {
		t.Fatalf("expected nonzero exit code after cancel, got 0")
	}
}

// TestLoginSessionWriteAfterCloseErrors verifies that calling Write after the session has
// exited returns an error containing "already closed".
func TestLoginSessionWriteAfterCloseErrors(t *testing.T) {
	bin := makeBashBin(t, "exit 0\n")
	p := LoginSessionParams{
		Executable: "/usr/bin/env",
		Argv:       []string{"bash", bin},
	}
	s, err := StartLoginSession(context.Background(), p)
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	// Drain output until closed.
	for range s.Output() {
	}

	if _, err := s.Wait(); err != nil {
		t.Fatalf("wait: %v", err)
	}

	err = s.Write([]byte("hello\n"))
	if err == nil {
		t.Fatal("expected error writing to closed session, got nil")
	}
	if !strings.Contains(err.Error(), "already closed") {
		t.Fatalf("expected 'already closed' in error, got %q", err.Error())
	}
}

// TestLoginSessionResize verifies that Resize does not return an error.
func TestLoginSessionResize(t *testing.T) {
	bin := makeBashBin(t, "sleep 2\n")
	p := LoginSessionParams{
		Executable: "/usr/bin/env",
		Argv:       []string{"bash", bin},
	}
	s, err := StartLoginSession(context.Background(), p)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() {
		_ = s.Cancel()
		_, _ = s.Wait()
	}()

	if err := s.Resize(40, 120); err != nil {
		t.Fatalf("resize: %v", err)
	}
}

// TestLoginSessionSurvivesCallerContextCancel verifies that the spawned login process
// is not tied to the lifetime of the RPC/request context that started it.
func TestLoginSessionSurvivesCallerContextCancel(t *testing.T) {
	bin := makeBashBin(t, "sleep 1\necho ready\nexit 0\n")
	ctx, cancel := context.WithCancel(context.Background())
	p := LoginSessionParams{
		Executable: "/usr/bin/env",
		Argv:       []string{"bash", bin},
	}
	s, err := StartLoginSession(ctx, p)
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	cancel()

	timeout := time.After(5 * time.Second)
	var collected []byte
	for {
		select {
		case chunk, ok := <-s.Output():
			if !ok {
				goto done
			}
			collected = append(collected, chunk...)
			if bytes.Contains(collected, []byte("ready")) {
				goto done
			}
		case <-timeout:
			t.Fatalf("timed out waiting for ready output after caller context cancel; got %q", collected)
		}
	}

done:
	code, err := s.Wait()
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit 0 after caller context cancel, got %d", code)
	}
	if !bytes.Contains(collected, []byte("ready")) {
		t.Fatalf("expected 'ready' in output; got %q", collected)
	}
}
