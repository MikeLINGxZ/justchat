package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLoggerWritesOnlyMessagesAtOrAboveConfiguredLevel verifies level filtering and daily log file creation.
func TestLoggerWritesOnlyMessagesAtOrAboveConfiguredLevel(t *testing.T) {
	tempDir := t.TempDir()
	clock := func() time.Time {
		return time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	}

	lg, err := New(Options{
		DataDir:  tempDir,
		Level:    "warn",
		Now:      clock,
		FileMode: 0o644,
		DirMode:  0o755,
	})
	if err != nil {
		t.Fatal(err)
	}

	lg.Info("skip this")
	lg.Warn("keep this")
	lg.Error("keep this too")

	bytes, err := os.ReadFile(filepath.Join(tempDir, "logs", "2026-05-10.log"))
	if err != nil {
		t.Fatal(err)
	}

	logText := string(bytes)
	if strings.Contains(logText, "skip this") {
		t.Fatalf("info message should not have been written: %s", logText)
	}
	if !strings.Contains(logText, "keep this") || !strings.Contains(logText, "keep this too") {
		t.Fatalf("expected warn and error messages in log file: %s", logText)
	}
}
