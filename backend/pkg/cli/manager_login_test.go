//go:build !windows

package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestManagerStartLoginUsesLoginCommand verifies StartLogin spawns a login session with the manifest's LoginCommand.
func TestManagerStartLoginUsesLoginCommand(t *testing.T) {
	mgr, dataDir := installToolFixture(t)
	binPath := filepath.Join(dataDir, "plugins", "cli", "lark-cli", "node_modules", ".bin", "lark-cli")
	if err := os.WriteFile(binPath, []byte("#!/usr/bin/env bash\nif [ \"$1\" = \"auth\" ]; then echo \"ready\"; exit 0; fi\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest.LoginCommand = []string{"auth"}
	if err := SaveManifest(manifestPath, manifest); err != nil {
		t.Fatal(err)
	}

	session, err := mgr.StartLogin(context.Background(), "lark-cli")
	if err != nil {
		t.Fatalf("StartLogin: %v", err)
	}

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	var accumulated []byte
drain:
	for {
		select {
		case chunk, ok := <-session.Output():
			if !ok {
				t.Fatalf("output channel closed before 'ready'; got: %q", string(accumulated))
			}
			accumulated = append(accumulated, chunk...)
			if strings.Contains(string(accumulated), "ready") {
				break drain
			}
		case <-timeout.C:
			t.Fatalf("timeout waiting for 'ready' in output; got: %q", string(accumulated))
		}
	}

	exitCode, err := session.Wait()
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}

// TestManagerStartLoginErrorsWhenNoCommand verifies StartLogin returns an error when LoginCommand is empty.
func TestManagerStartLoginErrorsWhenNoCommand(t *testing.T) {
	mgr, _ := installToolFixture(t)
	_, err := mgr.StartLogin(context.Background(), "lark-cli")
	if err == nil {
		t.Fatal("expected error when LoginCommand is empty")
	}
	if !strings.Contains(err.Error(), "no login_command") {
		t.Fatalf("expected error message to contain 'no login_command', got: %v", err)
	}
}
