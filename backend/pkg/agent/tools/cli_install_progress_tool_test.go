package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeCliInstallProgressReporter struct {
	sessionID uint
	item      CliInstallProgressItem
}

func (f *fakeCliInstallProgressReporter) ReportCliInstallProgress(_ context.Context, sessionID uint, item CliInstallProgressItem) error {
	f.sessionID = sessionID
	f.item = item
	return nil
}

func TestInvokeReportCliInstallProgressPersistsPayload(t *testing.T) {
	reporter := &fakeCliInstallProgressReporter{}
	out, err := InvokeReportCliInstallProgress(context.Background(), reporter, 42, json.RawMessage(`{
		"npm_package":"@lark/cli",
		"name":"lark-cli",
		"phase":"waiting_auth",
		"detail":"Open the verification URL to continue",
		"action_url":"https://example.com/verify",
		"expires_at":"10 minutes"
	}`))
	if err != nil {
		t.Fatalf("invoke report progress: %v", err)
	}
	if out != "cli install progress updated" {
		t.Fatalf("unexpected output: %q", out)
	}
	if reporter.sessionID != 42 {
		t.Fatalf("expected session id 42, got %d", reporter.sessionID)
	}
	if reporter.item.Phase != "waiting_auth" {
		t.Fatalf("expected waiting_auth phase, got %q", reporter.item.Phase)
	}
	if reporter.item.ActionURL != "https://example.com/verify" {
		t.Fatalf("expected action url to round-trip, got %q", reporter.item.ActionURL)
	}
}
